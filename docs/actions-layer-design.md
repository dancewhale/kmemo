# Actions 层设计文档

## 1. 文档目的

本文档在 [`backend-layering-design.md`](./backend-layering-design.md) 的总体分层之上，专门定义 **Actions 层** 的职责边界、目录结构、与 SuperMemo 风格操作需求的映射，以及实现依赖关系。

**读者**：实现 Wails 绑定、业务编排与仓储联动的开发者。

**不重复**：数据库表结构与 Repository 接口细节以 [`database-api-design.md`](./database-api-design.md) 为准；前端调用形状以 [`wails-api-design.md`](./wails-api-design.md) 为准。

---

## 2. 在架构中的位置

```text
Wails UI API (desktop/app.go, internal/app/)
    ↓
Actions 层 (internal/actions/)          ← 本文档
    ↓
Rules / Flows (internal/rules/, internal/flows/)
    ↓
Contracts + Adapters (internal/contracts/, internal/adapters/)
    ↓
Repository (internal/storage/repository/)
```

**一句话**：Actions 表达「用户要完成什么事」，编排多个 Repository 与 Contract，控制事务边界；**不**直接操作 Bleve、gRPC 细节、裸文件句柄。

---

## 3. SuperMemo 操作需求与 kmemo 的映射

SuperMemo 的核心是 **增量阅读** + **间隔重复** + **以知识库/条目组织材料**。kmemo 用现有领域模型对齐如下（便于从「操作」反查应落在哪一类 Action）：

| SuperMemo 风格需求 | kmemo 概念 | 主要 Action 包 |
|-------------------|------------|------------------|
| 知识库树、文件夹 | `Knowledge` 树 | `actions/knowledge` |
| 阅读材料 / 主题 vs 摘录 / 问答条目 | `Card` + `cardType`（article/excerpt/qa/cloze/note 等） | `actions/card` |
| 复习队列、评分、遗忘曲线 | FSRS + `CardSRS` + `ReviewLog` | `actions/review` |
| 增量阅读「下一步读什么」 | 阅读队列策略（前期可为计算视图） | `actions/queue` |
| 导入外部源、拆分为卡片 | `SourceDocument`（原始来源）+ `Article`（导入后的标准化阅读内容）+ 导入流水线 | `actions/importing` |
| 全文检索 | 索引器（规划中） | `actions/search` |
| 标签与过滤 | `Tag` 与关联表 | `actions/tag` |

若某条产品线需求与上表不一致，应在实现前更新本表（见第 9 节待确认项）。

---

## 4. Actions 层职责

### 4.1 必须承担

- 将「一个用户意图」落实为 **有序步骤**（先读哪张表、再调哪个算法、再写哪条日志）。
- **事务边界**：同一业务动作内需要原子一致的多表写入，应在 Action 内开启/提交/回滚（通过 `repository.Transaction` 或封装好的 `WithTx`）。
- **统一依赖注入**：通过构造函数接收 `Repository` 工厂、`contracts.*` 接口，便于测试替换。
- **与 UI 的输入输出契约**：接收来自 `internal/app` 已校验的 DTO 或领域输入结构，返回 DTO 或 `(id, error)`；不在 Action 内解析 JSON 或关心 Wails。

### 4.2 必须避免

- 不直接 `import` gRPC generated 包或 `grpc.ClientConn`（应经 `contracts` + `adapters/grpcworker`）。
- 不直接写 Bleve / 索引路径；只依赖 `contracts.SearchIndexer`（或等价抽象）。
- 不把 **长期稳定的业务恒等式** 写在 Action 里；应下沉到 `internal/rules`（例如评分语义、状态合法迁移）。
- 不把 **跨多子系统、步骤超过 ~10 步且会持续增长** 的流程整段写在单个 Action 文件；应抽取 `internal/flows`，由 Action 一行调用。

### 4.3 与 Rules / Flows 的分工

| 复杂度 | 建议 |
|--------|------|
| 简单 CRUD + 单次外部调用 | Action 内直接完成 |
| 多表一致 + FSRS + 日志 | Action 内事务；规则用 `rules/review` |
| 导入、撤销复习、重建索引 | `flows/*`，Action 仅做入口与参数组装 |

---

## 5. 建议目录结构

与 [`backend-layering-design.md`](./backend-layering-design.md) 第 4 节一致，按领域分子包；每个文件一个「动词」或「用例」：

```text
internal/actions/
├── knowledge/
│   ├── create.go          // CreateKnowledgeAction
│   ├── get.go
│   ├── list.go
│   ├── tree.go
│   ├── update.go
│   ├── delete.go
│   ├── move.go
│   ├── archive.go
│   └── unarchive.go
├── card/
│   ├── create.go
│   ├── get.go
│   ├── list.go
│   ├── update.go
│   ├── delete.go
│   ├── tags_add.go
│   ├── tags_remove.go
│   ├── suspend.go
│   └── resume.go
├── review/
│   ├── due_cards.go
│   ├── submit.go
│   ├── undo.go
│   ├── statistics.go
│   ├── history.go
│   └── stats_by_date.go
├── queue/
│   ├── next_for_reading.go   // 增量阅读「下一项」等（策略见 rules/queue）
│   └── list.go               // 可选：列表视图
├── importing/
│   └── import_document.go    // 或仅调 flows/import_document
├── search/
│   └── search_cards.go       // 全文检索入口（依赖 SearchIndexer）
├── tag/
│   ├── create.go
│   ├── list.go
│   ├── update.go
│   ├── delete.go
│   └── search_cards_by_tags.go
└── system/
    ├── system_info.go
    └── fsrs_parameters.go     // 列表/默认预设（只读或简单更新）
```

**命名约定**：导出类型建议后缀 `Action`，例如 `CreateCardAction`，与后端分层文档一致。

---

## 6. 各类 Action 的典型编排步骤

以下步骤描述「逻辑顺序」，实际表名与 Repository 方法名以实现为准。

### 6.1 知识库（`actions/knowledge`）

- **Create / Update / Delete / Move / Archive**：以 `KnowledgeRepository` 为主；删除前需业务规则检查（子节点、卡片引用等）时调用 `rules` 或轻量校验函数。
- **GetTree / List**：只读；可聚合 `Card`/`SRS` 计数（若 DTO 需要）时组合多个 Repository 查询。

### 6.2 卡片（`actions/card`）

- **CreateCard**（对齐增量阅读「从摘录生成条目」）：
  1. 校验 `Knowledge` 存在且可用。
  2. 通过 `contracts.FileStore` 写入正文（路径写入 `Card` 模型字段）。
  3. `CardRepository.Create`；按需 `AddTags`。
  4. 若卡片参与复习：初始化 `CardSRS`（规则见 `rules`）。
  5. 若启用搜索：`SearchIndexer.IndexCard`（与 DB 提交顺序见第 7 节）。
- **UpdateCard**：更新元数据；若 HTML 变更则覆盖文件 + 重索引。
- **DeleteCard**：软删除 DB；文件移入 trash 或按 FileStore 策略；索引删除。

### 6.3 复习（`actions/review`）

- **GetDueCards**：`SRSRepository.GetDueCards` + 联接 `Card` 展示字段；可过滤 `cardType`（若规则要求 article 不参与队列）。
- **SubmitReview**：
  1. 加载 `Card` + `CardSRS`。
  2. 通过 `contracts.FSRSClient`（例如 Python gRPC 的 FSRS 领域适配器）计算新调度。
  3. 事务内：`UpdateAfterReview`（含 `ReviewLog`）。
  4. 可选：更新阅读队列相关派生状态（若未来落表）。
- **UndoLastReview**：优先委托 `flows/undo_review` 或 `SRSRepository.UndoLastReview` 封装，保证与日志一致。

### 6.4 阅读队列（`actions/queue`）

- **NextForReading / List**：按 [`backend-layering-design.md`](./backend-layering-design.md) 第 5.5 节，**前期推荐计算视图**：依据 `rules/queue`（排序、过滤 article、每日上限等）查询 `Card` + 可选 SRS 字段，**不一定**单独建 `reading_queue` 表。

### 6.5 导入（`actions/importing`）

- **ImportDocument**：推荐由 Action 发起并编排异步导入任务，而不是将整份材料同步塞进单次 `PrepareImportMaterial` RPC。
- 推荐边界：
  1. Go 侧创建并持久化 `import_jobs` 状态，生成 `job_id`。
  2. Go 侧准备 `workspace_dir / output_dir / temp_dir` 与原始输入文件位置。
  3. 通过 **任务式 source-process contract** 提交异步 job（`SubmitImportJob`），并按需轮询/订阅 `GetJob` / `ListJobEvents`。
  4. Python worker 负责在工作目录下产出标准输出目录（`article/`、`cleaned.html`、`manifest.json`、`assets/`）。
  5. Go 侧在任务成功后读取 `manifest` 与输出目录，决定如何创建/更新 `SourceDocument`、`Article`、`Card`、`Asset` 等业务对象。
- 约束：
  - `CleanHtml` 仍保留为同步轻量能力；长流程导入不再复用旧的同步 `PrepareImportMaterial` 主路径。
  - job workspace 属于流程工作目录，不等同于正式 `FileStore`；仅最终需要长期保留的来源文件、文章 HTML、资产文件再迁入 `FileStore`。

### 6.6 搜索（`actions/search`）

- **SearchCards**：调用 `SearchIndexer.Search`，再按 ID 批量加载 `Card`（或仅存索引字段时的补齐策略）。

### 6.7 标签（`actions/tag`）

- CRUD 以 `TagRepository` 为主；`SearchCardsByTags` 可走 DB 过滤或索引，与 [`wails-api-design.md`](./wails-api-design.md) 一致。

---

## 7. 事务、索引与文件的一致性

推荐顺序（与 [`backend-layering-design.md`](./backend-layering-design.md) 第 5.1、8.1 节一致）：

1. 业务校验。
2. 文件写入成功（或进入可回滚的临时区后再 promote）。
3. 数据库提交。
4. 搜索索引更新（当前阶段可同步；失败时应记录可重试日志，后续由 `flows/rebuild_index` 修复）。

若某一步失败，Action 应明确是 **整单失败** 还是 **部分成功 + 补偿**；MVP 建议整单失败并依赖事务与文件操作的幂等设计。

---

## 8. 与 Wails UI API 的对应关系

`internal/app` 中的方法应保持 **薄封装**：参数校验 → 调用对应 `*Action` → 映射 `dto`。

| Wails 方法（节选） | 建议调用的 Action |
|-------------------|-------------------|
| `CreateKnowledge` | `knowledge.Create` |
| `CreateCard` | `card.Create` |
| `GetDueCards` / `SubmitReview` / `UndoLastReview` | `review.DueCards` / `review.Submit` / `review.Undo` |
| `ImportDocument` | `importing.Import` / `flows.ImportDocument`（提交异步 source-process job、跟踪状态、消费 manifest） |

完整列表可在实现阶段用表格补全。

---

## 9. 待与产品确认的事项

以下内容会影响 Action 的输入参数、是否落表、以及是否单独拆 `flows`，建议在编码前确认：

1. **阅读队列**：必须支持「人工拖拽排序 / 置顶 / 延后」等 **持久化顺序**.需要 `reading_queue` 表或等价结构，Actions 将多出 `queue/reorder`、`queue/postpone` 等；请帮我同步修改 **计算视图**（[`backend-layering-design.md`](./backend-layering-design.md) 第 5.5 节）。
2. **卡片类型与 SRS**：`article` 需要进入标准复习队列，并它和`excerpt/qa/cloze` 的默认 SRS 初始化完全一致.
3. **导入**：已确认采用 **异步任务 + 进度推送**；Go 持久化 job 状态为真相源，Python 负责执行与事件上报。旧的同步 `PrepareImportMaterial` 不再作为主路径。
4. **FSRS 参数作用域**：全局一套默认参数。
5. **全文搜索 MVP**：首版只做 **标题/标签/SQLite**，Bleve 作为 Phase 2.`actions/search` 可先实现 DB 路径，接口仍预留 `SearchIndexer`。

---

## 10. 底层依赖未实现项 — 建议优先任务

结合当前仓库状态（`internal/actions` 尚未建立；`contracts` 除 `FileStore` 外多数能力未以接口固化；`internal/indexing` 为占位；`htmlproc` 未接线；Desktop 尚未注入 Repository）：

| 优先级 | 任务 | 说明 |
|--------|------|------|
| P0 | 建立 `internal/actions` 包与构造函数签名 | 定义 `Action` 依赖 `repository.*` 与 `contracts.*` 的注入方式，与 `bootstrap.NewDesktop` 扩展衔接。 |
| P0 | `bootstrap` / `app.Desktop` 注入 Repository 与事务 | 使 Action 可运行；与 [`database-api-design.md`](./database-api-design.md) 中的 `Repository` 工厂一致。 |
| P1 | `contracts.SearchIndexer` + Bleve（或 SQLite FTS）适配 | `internal/indexing/indexing.go` 仅为 TODO；与 `actions/search`、`card` 删除/更新索引联动。 |
| P1 | 核心 `flows`：`SubmitReview`、`ImportDocument`、`UndoReview` | 步骤稳定后从 Action 中沉淀，避免 `submit.go` 膨胀。 |
| P1 | `rules/review`、`rules/queue` | 评分语义、到期筛选、阅读队列排序策略。 |
| P2 | 导入后台任务（若确认需求） | 任务持久化、进度回调与 Action 的边界。 |
| P2 | Outbox 或异步重建索引（规模上来后） | 与 [`backend-layering-design.md`](./backend-layering-design.md) 第 5.1 节后期方案一致。 |

---

## 11. 结论

Actions 层是 **理解 kmemo 业务的主线**：读一个 Action 文件应能回答「完成该 SuperMemo 风格操作时，系统依次做了什么」。通过严格遵守对 Repository 与 Contracts 的依赖方向，可以在不重写业务流程的前提下替换 FSRS 实现、搜索引擎与文件后端。

**与分层总述的关系**：[`backend-layering-design.md`](./backend-layering-design.md) 回答「层与层之间放什么」；本文档回答「Actions 这一层有哪些包、哪些用例、依赖谁、哪些还没实现」。
