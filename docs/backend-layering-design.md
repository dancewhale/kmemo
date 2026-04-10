# 后端分层设计文档

## 1. 目标

本文档用于明确 `Wails UI API` 到 `database API` 之间的后端代码分层，帮助项目在后续接入以下能力时保持结构清晰、可理解、可修改：

- 数据库读写
- Bleve 全文索引
- HTML 清洗与标准化处理
- 文件存储
- Python gRPC 能力（FSRS、HTML、导入处理）
- 阅读队列
- 卡片管理和复习
- 后续可能的异步任务、重建索引、导入流程

核心目标不是“分层越多越好”，而是：

- 每一层职责明确
- 读代码时能快速判断某段逻辑应该放哪里
- 业务流程容易追踪
- 基础设施可替换
- 以后重构时改动范围可控

---

## 2. 总体分层

建议采用 **5 层结构**：

```text
Wails UI API
    ↓
Actions 层
    ↓
Rules / Flows 层
    ↓
Contracts + Adapters 层
    ↓
Repository / Database API 层
```

对应当前项目可以理解为：

```text
desktop/app.go + internal/app/desktop.go
    ↓
internal/actions/
    ↓
internal/rules/ + internal/flows/
    ↓
internal/contracts/ + internal/adapters/
    ↓
internal/storage/repository/
```

其中，`UI API` 到 `database API` 中间重点保留 **3 层**：

1. `actions`
2. `rules/flows`
3. `contracts/adapters`

这是兼顾可读性与实际复杂度后的平衡方案。

---

## 3. 各层职责

## 3.1 UI API 层

**建议位置**：

- `desktop/app.go`
- `internal/app/desktop.go`
- `internal/app/dto/`

**职责**：

- 向 Wails 暴露方法
- 接收前端参数，返回前端 DTO
- 做最薄的一层参数转换和输入校验
- 调用 actions
- 不直接操作 repository
- 不直接处理 Bleve、文件、HTML、Python gRPC

**这一层应该像 Controller / Facade。**

### 适合放在这里的内容

- `CreateCard(req)`
- `GetDueCards(...)`
- `SubmitReview(...)`
- `ListKnowledge(...)`

### 不应该放在这里的内容

- 卡片创建事务逻辑
- HTML 清洗细节
- 索引写入
- 文件路径分配
- 复习状态机逻辑

---

## 3.2 Actions 层

**建议位置**：

- `internal/actions/`

**职责**：

- 表达“用户要完成什么事”
- 编排多个 repository / contract
- 决定一次业务动作的执行顺序
- 控制事务边界
- 统一输入输出结构

**这一层是未来理解系统最重要的一层。**

读 action 应该能直接回答：

- “创建卡片到底做了哪些事？”
- “提交复习时依次发生了什么？”
- “导入文档会更新哪些子系统？”

### 典型 action 示例

#### CreateCardAction

负责：
1. 校验 knowledge 是否存在
2. 调用 SourceProcessor 清理内容
3. 调用 FileStore 保存 HTML
4. 写入 Card 数据
5. 绑定标签
6. 更新 Bleve 索引
7. 返回 DTO

#### SubmitReviewAction

负责：
1. 获取卡片与当前 SRS
2. 调用 FSRSClient 计算新调度结果
3. 更新 `card_srs`
4. 写入 `review_log`
5. 必要时更新阅读队列或统计
6. 返回结果

### 这一层的规则

- actions 可以依赖 repository 和 contracts
- actions 不应该直接依赖 Bleve、os.File、grpc 细节实现
- actions 不承载具体业务规则细节时，可下沉给 rules/flows

---

## 3.3 Rules / Flows 层

建议拆成两个子方向：

### A. `internal/rules/`

放 **稳定的业务规则、约束、状态转换**。

### B. `internal/flows/`

放 **跨多个子系统的完整流程**。

---

### 3.3.1 Rules 层职责

**适合放这里的内容**：

- CardType 的定义和合法值
- Review rating 的语义
- 复习状态转换规则
- 阅读队列的排序/过滤规则
- 搜索文档的规则化表示
- 导入后卡片生成的基本规则

**不适合放这里的内容**：

- 具体数据库查询
- 具体 Bleve 调用
- 具体文件系统读写
- 具体 grpc 调用

### 3.3.2 Flows 层职责

**适合放这里的内容**：

- 导入文档流程
- 提交复习流程
- 撤销复习流程
- 从摘录生成卡片流程
- 重建搜索索引流程

### 为什么 rules/flows 需要存在

因为这个项目后续不是简单 CRUD，而是明显会出现大量“动作型功能”：

- 导入一本书
- 清理 HTML 并拆分卡片
- 复习一张卡
- 撤销一次复习
- 生成阅读队列
- 重建全文索引

如果这些逻辑全部塞进 actions，前期可行，后期会快速膨胀。

因此建议：

- **简单动作**：actions 直接完成
- **复杂动作**：actions 调 flows
- **稳定规则**：放 rules

---

## 3.4 Contracts + Adapters 层

**建议位置**：

- `internal/contracts/`
- `internal/adapters/`

这层用于隔离所有“外部能力约定”和“具体适配实现”。

---

### 3.4.1 Contracts 层职责

定义系统依赖的能力约定，而不是具体实现。

暂时举例,不作为后续约束和设计基础
例如：

```go
type SearchIndexer interface {
    IndexCard(ctx context.Context, doc SearchDocument) error
    DeleteCard(ctx context.Context, cardID string) error
    Search(ctx context.Context, query SearchQuery) ([]SearchHit, error)
}

type FileStore interface {
    Save(ctx context.Context, path string, content []byte) error
    Read(ctx context.Context, path string) ([]byte, error)
    Delete(ctx context.Context, path string) error
}

type FSRSClient interface {
    Calculate(ctx context.Context, req FSRSRequest) (*FSRSResult, error)
}

type SourceProcessClient interface {
    SubmitImportJob(ctx context.Context, req SubmitImportJobRequest) (*SubmitImportJobResult, error)
    GetJob(ctx context.Context, jobID string) (*SourceProcessJob, error)
    ListJobEvents(ctx context.Context, jobID string) ([]SourceProcessJobEvent, error)
    CancelJob(ctx context.Context, jobID string) error
    GetCapabilities(ctx context.Context) (*SourceProcessCapabilities, error)
}
```

### 3.4.2 Adapters 层职责

提供 contracts 的具体适配实现。

建议实现目录：

```text
internal/adapters/
├── bleve/
├── filestore/
├── htmlproc/
├── grpcworker/
├── fsrs/
├── sourceprocess/
└── clock/
```

### 这层的收益

通过 `contracts/adapters` 分离，未来可以替换：

- Bleve → SQLite FTS / 其他搜索实现
- 本地文件 → 对象存储 / 云同步目录
- Python FSRS → Go 原生实现
- Python HTML 清洗 → Go 实现

而 actions 不需要重写。

---

## 3.5 Repository / Database API 层

**建议位置**：

- `internal/storage/repository/`

**职责**：

- 面向数据库实体读写
- 封装 GORM / DAO
- 提供事务内外的统一数据访问接口
- 不承载 UI 语义
- 不直接更新索引
- 不直接写文件
- 不承载复杂业务流程

**判断标准**：

如果一段逻辑的核心问题是“这个表怎么查/怎么写”，就属于 repository。
如果核心问题是“为了完成某个用户动作，需要协调哪些子系统”，就不该放 repository。

---

## 4. 推荐目录结构

建议逐步整理为如下结构：

```text
internal/
├── app/                        # Wails/CLI 绑定层
│   ├── desktop.go
│   └── dto/
│       ├── knowledge.go
│       ├── card.go
│       ├── review.go
│       └── tag.go
│
├── actions/                    # 动作层：面向前端动作
│   ├── knowledge/
│   │   ├── create.go
│   │   ├── list.go
│   │   └── move.go
│   ├── card/
│   │   ├── create.go
│   │   ├── update.go
│   │   ├── delete.go
│   │   └── list.go
│   ├── review/
│   │   ├── submit.go
│   │   ├── undo.go
│   │   └── due_cards.go
│   ├── importing/
│   │   └── import_document.go
│   └── search/
│       └── search_cards.go
│
├── flows/                      # 完整流程
│   ├── import_document.go
│   ├── submit_review.go
│   ├── undo_review.go
│   ├── create_card_from_excerpt.go
│   └── rebuild_index.go
│
├── rules/                      # 稳定业务规则
│   ├── card/
│   │   ├── types.go
│   │   └── rules.go
│   ├── review/
│   │   ├── rating.go
│   │   └── state_machine.go
│   ├── queue/
│   │   └── policy.go
│   └── search/
│       └── document.go
│
├── contracts/                  # 外部能力约定
│   ├── search.go
│   ├── html.go
│   ├── filestore.go
│   ├── fsrs.go
│   ├── source_process.go
│   └── clock.go
│
├── adapters/                   # contracts 的具体适配实现
│   ├── bleve/
│   ├── filestore/
│   ├── htmlproc/
│   ├── grpcworker/
│   ├── fsrs/
│   ├── sourceprocess/
│   └── clock/
│
├── storage/
│   ├── models/
│   ├── dao/
│   ├── repository/
│   ├── db.go
│   └── storage.go
│
├── config/
└── bootstrap/
```

---

## 5. 几个关键能力分别放哪

## 5.1 Bleve 全文索引

### 建议放置

- 能力约定：`internal/contracts/search.go`
- 适配实现：`internal/adapters/bleve/`
- 搜索文档结构：`internal/rules/search/`
- 索引编排：`internal/actions/search/` 或 `internal/flows/rebuild_index.go`

### 不建议放置

- 不放在 repository
- 不放在 Wails API
- 不让 CardRepository 直接更新索引

### 原因

Bleve 是搜索基础设施，不是数据库本体。
如果 repository 直接更新 Bleve，会导致数据库层掺入基础设施耦合，后续替换和调试都困难。

### 当前阶段建议

先采用 **actions 内同步更新索引**：

1. DB 写入成功
2. 文件保存成功
3. 调用 SearchIndexer 更新索引

后期如果规模变大，再考虑：

- outbox
- 异步重建
- 定时修复索引

---

## 5.2 HTML 处理

### 建议放置

- 抽象接口：`internal/contracts/html.go`
- Go 实现：`internal/adapters/htmlproc/`
- Python 调用实现：`internal/adapters/htmlproc/` 或通过共享 `grpcworker` 连接的远端适配器

### 设计建议

上层只依赖：

```go
type SourceProcessor interface {
}
```

这样 actions 不需要知道：

- 是 Go 本地清理
- 还是 Python gRPC 清理
- 还是两者混合

### 原则

Python 端即便由 source-process 相关服务实现，也不改变上层 contract 仍是 `SourceProcessor` 这一事实。
HTML 的业务规则可以由 actions/flows 决定，但 HTML 的技术处理能力应收敛到 `SourceProcessor` 接口后面。

---

## 5.3 文件存储

### 建议放置

- 抽象接口：`internal/contracts/filestore.go`
- 具体实现：`internal/adapters/filestore/`

### 原因

正式业务资产（卡片 HTML、长期保留的来源文件、最终资产文件）应统一经 `FileStore` 管理；
但 source-process job 的 `workspace_dir / output_dir / temp_dir` 属于流程工作目录，不等同于正式 `FileStore`。
建议边界如下：

- job workspace：供导入流程中间产物和标准输出目录使用
- FileStore：仅承接最终需要长期保留并被业务对象引用的文件

这样可避免把临时流程目录误建模成正式业务存储。

### 建议能力

- 保存 HTML
- 保存附件
- 读取文件
- 删除文件
- 移动文件
- 检查是否存在

### 推荐数据目录结构

```text
data/
├── cards/
│   ├── <card-id>/content.html
│   └── <card-id>/answer.html
├── assets/
│   ├── images/
│   ├── audio/
│   └── files/
├── sources/
│   └── <source-id>/original.xxx
└── search/
    └── bleve/
```

### 原则

- DB 里存相对路径
- 文件读写统一经 FileStore
- 业务层不直接拼接物理路径

---

## 5.4 Python gRPC 能力

### 建议放置

- 接口：`internal/contracts/fsrs.go`、`internal/contracts/html.go`、`internal/contracts/source_process.go`
- 实现：`internal/adapters/grpcworker/` + 领域适配器（如 `internal/adapters/fsrs/`、`internal/adapters/sourceprocess/`）

### 建议拆分

不要把 Python 视作一个大而全的“远程服务对象”，而是拆成能力接口：

- `FSRSClient`
- `SourceProcessor`
- `SourceProcessClient`

连接层可通过共享 `grpcworker` 复用同一条 gRPC 连接，但 proto service 与 adapter 应按领域拆分，而不是继续保留历史上的单一 Python 网关大包。

---

## 5.5 阅读队列

### 设计建议

阅读队列不要一开始塞进 SRS repository 或 Card repository。
它更像一个“策略模块”，而不是一个单纯表操作模块。

### 建议分层

- 队列规则：`internal/rules/queue/`
- 队列动作：`internal/actions/queue/`
- 如果未来落库存储：再增加对应 repository

### 前期建议

先采用 **计算视图**：

- 根据规则动态计算阅读队列
- 不单独落表

适合当前阶段，简单、可变性高。

### 后期何时考虑落表

当你出现这些需求时再考虑新增 `reading_queue` 表：

- 人工重排顺序
- 置顶/延后
- 队列快照
- 队列变更日志
- 稳定持久化当前位置

### 原则

- 前期把“规则”放 domain
- 把“调度动作”放 actions
- 不要过早数据库化

---

## 6. 推荐的命名方式

为提高长期可读性，建议命名直接体现层次和职责。

### UI API 层

- `CreateCard`
- `SubmitReview`
- `GetDueCards`

### Actions 层

- `CreateCardAction`
- `SubmitReviewAction`
- `SearchCardsAction`

### Flows 层

- `ImportDocumentFlow`
- `UndoReviewFlow`
- `RebuildIndexFlow`

### Rules 层

- `CardRules`
- `ReviewStateRules`
- `QueuePolicy`

### Contracts 层

- `SearchIndexer`
- `FileStore`
- `SourceProcessor`
- `FSRSClient`

### Adapters 层

- `BleveIndexer`
- `LocalFileStore`
- `PythonSourceProcessor`
- `PythonFSRSClient`

这样一眼就能看出：

- 是“业务动作”
- 还是“领域规则”
- 还是“基础设施实现”

---

## 7. 不建议的做法

## 7.1 不建议增加含糊的 manager 层

例如：

- `CardManager`
- `ReviewManager`
- `IndexManager`

这些命名通常边界不清，时间久了容易演变成大杂烩。

---

## 7.2 不建议形成过深调用链

例如：

```text
UI -> facade -> service -> manager -> helper -> repository
```

层级太深会让项目变得“看起来分层很多，实际上很难理解”。

推荐保持：

```text
UI -> actions -> flows/rules -> contracts/repository
```

---

## 7.3 不建议把基础设施逻辑塞进 domain

例如：

- domain 直接调用 Bleve
- domain 直接写文件
- domain 直接发 grpc

domain 只负责规则，不负责外部依赖。

---

## 7.4 不建议让 repository 承担跨系统业务流程

例如：

- `CardRepository.CreateAndIndex()`
- `ReviewRepository.SubmitReviewAndUpdateSRS()`

这类接口会让数据库层承担本不属于它的职责。

---

## 8. 两个典型调用链

## 8.1 创建卡片

```text
Wails UI
  -> Desktop.CreateCard(req)
  -> actions/card/create.go
      -> repository.Knowledge.GetByID
      -> contracts.SourceProcessor.Clean
      -> contracts.FileStore.SaveHTML
      -> repository.Card.Create
      -> repository.Card.AddTags
      -> contracts.SearchIndexer.IndexCard
  -> 返回 cardID / DTO
```

### 理解收益

读 `CreateCardAction` 时就能完整看见：

- 数据如何进入系统
- HTML 何时处理
- 文件何时落盘
- 数据何时写库
- 索引何时更新

---

## 8.2 提交复习

```text
Wails UI
  -> Desktop.SubmitReview(req)
  -> actions/review/submit.go
      -> repository.Card.GetByID
      -> repository.SRS.GetByCardID
      -> contracts.FSRSClient.Calculate
      -> repository.SRS.UpdateAfterReview
      -> repository.ReviewLog.Create / 或在事务内统一处理
      -> actions/queue 或 flows 更新阅读状态（如需要）
  -> 返回新的 SRS 状态
```

### 理解收益

复习逻辑不会散在：

- Wails 方法里
- repository 里
- gRPC 领域适配器里

而是集中在一个明确的 actions / flows 中。

---

## 9. 当前阶段的最小落地方案

为了兼顾当前项目状态与未来演进，建议分两步推进。

### 第一阶段：先建立最小清晰结构

优先引入：

```text
internal/
├── app/
├── actions/
├── contracts/
├── adapters/
├── storage/repository/
├── flows/   # 先少量使用
└── domain/     # 先只放稳定规则
```

### 第二阶段：随着功能增长再增强

再逐步补强：

- 队列策略
- 搜索文档建模
- 导入流程 flows
- 索引重建 flows
- 更完整的 domain state machine

### 原则

- 先把结构搭清楚
- 不一开始做过重的 DDD
- 复杂性增长时再扩 domain/flows

---

## 10. 结论

对于 kmemo 这个 `Wails + Golang + Python` 项目，推荐采用以下结构：

### 总体 5 层

1. UI API
2. Actions
3. Domain / Flows
4. Contracts / Adapters
5. Repository / Database API

### UI API 到 DB API 中间保留 3 层

1. `actions`
2. `domain/flows`
3. `contracts/adapters`

### 设计收益

这种结构最适合本项目未来的发展方向：

- DB 读写、Bleve、文件、HTML、Python 分工清楚
- 阅读队列可以自然接入
- 业务动作的主流程容易阅读
- 基础设施替换影响范围小
- 后续重构不会牵一发而动全身

### 一句话总结

**让 UI 只表达入口，让 actions 表达动作，让 domain/flows 表达规则和流程，让 contracts/adapters 承担外部依赖，让 repository 只做数据持久化。**
