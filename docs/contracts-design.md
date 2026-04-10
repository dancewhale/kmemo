# Contracts 层设计文档

## 1. 目标

本文档用于明确 kmemo 后端中的 `internal/contracts/` 应该承担什么职责、定义哪些接口、如何与 `actions`、`rules/flows`、`adapters`、`repository` 协作。

目标不是把所有东西都抽象成 interface，而是把**非数据库外部能力**收敛成一组稳定、可理解、可替换的能力约定。

contracts 层主要解决这几个问题：

- actions 不直接依赖 gRPC、Bleve、文件系统等技术实现
- 后续替换 Python / Bleve / 本地文件存储时，不需要重写业务动作
- 业务代码读起来能快速区分：
  - 哪些是数据库持久化
  - 哪些是外部能力调用
  - 哪些是业务动作编排
- 接口命名和输入输出尽量贴近业务意图，而不是贴近底层 SDK

---

## 2. contracts 在整体分层中的位置

```text
Wails UI API
    ↓
Actions
    ↓
Rules / Flows
    ↓
Contracts
    ↓
Adapters
    ↓
外部能力实现（Python gRPC / Bleve / 文件系统 / 时钟等）
```

同时它与 repository 的关系是并列而非替代：

```text
Actions / Flows
   ├── Repository        # 数据库存取
   └── Contracts         # 非数据库外部能力
```

### 一句话理解

- `repository` 解决“数据怎么查、怎么写”
- `contracts` 解决“系统要向外部能力要什么结果”

---

## 3. contracts 层的职责边界

## 3.1 应该放在 contracts 的内容

适合放进 contracts 的，通常具有这几个特征：

- 不是数据库 CRUD
- 背后依赖外部系统、外部进程、文件系统或基础设施能力
- actions 只关心“要什么能力”，不关心“怎么实现”
- 后续存在替换实现的可能

当前项目中最适合放入 contracts 的能力：

1. 搜索索引
2. HTML 处理
3. 文件存储
4. FSRS 调度计算
5. source-process 导入任务编排
6. 时钟/时间来源
7. 可选：任务调度、事件发布（未来）

## 3.2 不应该放在 contracts 的内容

以下内容不建议塞进 contracts：

- `KnowledgeRepository`、`CardRepository` 这类数据库接口
- 纯业务规则判断
- 前端 DTO
- 只在单个 action 内部使用的一次性 helper
- 直接暴露第三方 SDK 原始对象

## 3.3 contracts 不应该做的事

contracts 只定义约定，不承载这些内容：

- 不写 Bleve 具体逻辑
- 不写 grpc 请求细节
- 不写本地文件路径拼接细节
- 不做数据库事务控制
- 不承载业务动作编排

这些都应分别放在 adapters、repository、actions/flows 中。

---

## 4. 设计原则

## 4.1 面向动作，而不是面向技术协议

contracts 的接口应该让 actions 看起来像是在调用业务能力，而不是在拼底层协议。

例如：

```go
type PythonClient interface {
    CallCleanHtmlRPC(ctx context.Context, protoReq *pb.CleanHtmlRequest) (*pb.CleanHtmlResponse, error)
}
```

原因：后者让上层直接知道 gRPC proto 和具体传输方式，耦合过深。

## 4.2 接口粒度按能力拆，不按实现来源拆

不要按“Python 服务”“Bleve 服务”“本地文件服务”这种实现来源拆 contracts。
应按能力拆：

- `FSRSClient`
- `SourceProcessClient`
- `SearchIndexer`
- `FileStore`
- `Clock`

这样以后即使 `FileStore` 从 Python 改成 Go，本层接口仍然稳定。

## 4.3 输入输出要稳定、可测试、可理解

contracts 不应直接暴露：

- proto message
- `os.File`
- Bleve 原始 document/query 对象
- gRPC client

应该定义自己的 request/result 结构，作为 actions 与 adapters 间的稳定边界。

## 4.4 错误语义要统一到业务可判断层级

actions 不应该关心 gRPC status code、Bleve 错误类型、文件系统 syscall 错误。
adapters 应把底层错误转换为 contracts 层可识别的错误语义。

例如：

- `ErrUnavailable`：外部能力暂时不可用
- `ErrInvalidInput`：输入非法
- `ErrNotFound`：资源不存在
- `ErrConflict`：冲突

这样 actions 才能基于统一语义做处理。

---

## 6.5 SourceProcessClient

### 适用场景

- 提交 HTML/EPUB/PDF 等来源导入任务
- 跟踪 Python worker 的异步处理状态与事件
- 查询 converter / cleaner / mode 等能力发现信息

### 建议接口

```go
type SourceProcessClient interface {
    SubmitImportJob(ctx context.Context, req SubmitImportJobRequest) (*SubmitImportJobResult, error)
    GetJob(ctx context.Context, jobID string) (*SourceProcessJob, error)
    ListJobEvents(ctx context.Context, jobID string) ([]SourceProcessJobEvent, error)
    CancelJob(ctx context.Context, jobID string) error
    GetCapabilities(ctx context.Context) (*SourceProcessCapabilities, error)
}
```

### 设计说明

- `SourceProcessClient` 是 source-process 领域的**任务式 transport contract**，承接 Python worker 的异步导入能力。
- 为了保持 action 侧语义清晰，可以在更高层保留一个面向业务的导入 flow/Action 封装，但 contracts 主边界应以异步 job 为准，而不是旧的同步 `PrepareImportMaterial`。

### 建议数据结构（最小集合）

```go
type SubmitImportJobRequest struct {
    JobID         string
    SourceType    string
    SourcePath    *string
    SourceURI     *string
    SourceURL     *string
    RawHTML       *string
    WorkspaceDir  string
    OutputDir     string
    TempDir       string
    Options       SourceProcessOptions
    Metadata      map[string]string
    IdempotencyKey *string
}

type SubmitImportJobResult struct {
    JobID   string
    Status  string
    Accepted bool
}

type SourceProcessJob struct {
    JobID      string
    Status     string
    Stage      string
    Progress   float32
    ResultPath *string
    ErrorCode  *string
    ErrorMessage *string
}

type SourceProcessJobEvent struct {
    JobID      string
    Sequence   int64
    Stage      string
    Message    string
    CreatedAtUnix int64
}

type SourceProcessCapabilities struct {
    SourceTypes      []string
    ConversionModes  []string
    ConverterNames   []string
    CleanerNames     []string
}

type SourceProcessOptions struct {
    ConversionMode        string
    FallbackModes         []string
    ExtractMainContent    bool
    SanitizeHTML          bool
    PreserveSemanticTags  bool
    DownloadRemoteAssets  bool
    InlineSmallImages     bool
    GenerateTOC           bool
    AnalyzeStructure      bool
    KeepSourceCopy        bool
    EnabledCleaners       []string
    ConverterParamsJSON   *string
}
```

### 设计要点

- 大文件主路径通过 `source_path/source_uri/workspace_dir/output_dir/temp_dir` 传递，不通过 gRPC 直接传大字节。
- Python 输出的是**标准输出目录 + manifest**；Go 在任务成功后读取结果并决定如何创建 `SourceDocument`、`Article`、`Card`、`Asset`。
- Go 是导入 job 状态的持久化真相源；Python 负责执行与事件上报。

---

## 6.6 Clock

### 适用场景

- 复习时间、导入时间、统计窗口
- 避免在 actions/flows 中直接到处调用 `time.Now()`

### 建议接口

```go
type Clock interface {
    Now() time.Time
}
```

### 设计要点

- 接口很小，但很有价值
- 便于测试复习调度、到期判断、统计逻辑
- 不建议再额外设计复杂时间工具层

---

## 7. 错误设计建议

建议在 `internal/contracts/errors.go` 统一定义错误语义：

```go
var (
    ErrUnavailable = errors.New("contract unavailable")
    ErrInvalidInput = errors.New("contract invalid input")
    ErrNotFound = errors.New("contract resource not found")
    ErrConflict = errors.New("contract conflict")
)
```

如果希望保留更多上下文，可加错误包装类型：

```go
type ContractError struct {
    Capability string
    Op         string
    Err        error
}
```

设计原则：

- actions 只判断 contracts 错误语义
- adapters 负责把 grpc / bleve / fs / 第三方错误映射进来
- 不要让上层直接依赖 `status.Code(err)` 之类技术细节

---

## 8. contracts 与 actions 的协作方式

## 8.1 CreateCardAction

推荐调用链：

```text
CreateCardAction
  -> repository.Knowledge.GetByID
  -> contracts.SourceProcess.Clean
  -> contracts.FileStore.Save
  -> repository.Card.Create
  -> repository.Card.AddTags
  -> contracts.SearchIndexer.IndexCard
```

这里的关键点：

- action 负责顺序和事务边界
- contracts 负责提供外部能力
- adapters 负责具体实现

## 8.2 SubmitReviewAction

```text
SubmitReviewAction
  -> repository.Card.GetByID
  -> repository.SRS.GetByCardID
  -> repository.FSRSParameter.GetDefault
  -> contracts.FSRSClient.SetScheduler
  -> contracts.FSRSClient.Calculate
  -> repository.SRS.UpdateAfterReview
  -> repository.ReviewLog.Create
```

## 8.3 ImportDocumentFlow

```text
ImportDocumentFlow
  -> repository.Knowledge.GetByID
  -> contracts.SourceProcessClient.SubmitImportJob
  -> repository.ImportJob.Create / UpdateStatus
  -> contracts.SourceProcessClient.GetJob / ListJobEvents
  -> 读取 manifest 与标准输出目录
  -> repository.SourceDocument.Create
  -> repository.Article.Create
  -> repository.Card.Create
  -> contracts.SearchIndexer.IndexCard
```

---

## 9. contracts 与 adapters 的对应关系

推荐保持“一组 contract，可有多个 adapter 实现”：

```text
contracts.FSRSClient
  ├── adapters/fsrs/grpc_python/client.go
  └── adapters/fsrs/local.go               # 未来可选

contracts.SourceProcessClient
  └── adapters/sourceprocess/client.go

contracts.SearchIndexer
  ├── adapters/bleve/indexer.go
  └── adapters/search/sqlite_fts.go        # 未来可选

contracts.FileStore
  ├── adapters/filestore/local.go
  └── adapters/filestore/object_store.go   # 未来可选

contracts.Clock
  └── adapters/clock/system.go
```

原则：

- actions 只依赖 contracts
- bootstrap 负责注入具体 adapter 实现
- 不把 adapter 暴露回 actions

---

## 10. 是否需要 ContractFactory

当前阶段**不建议**额外引入 `ContractFactory`。

原因：

- contract 数量还不多
- 过早增加 factory 会让依赖跳转更多
- 直接在 action struct 中显式注入更清晰

推荐先采用：

```go
type CreateCardAction struct {
    cardRepo      repository.CardRepository
    knowledgeRepo repository.KnowledgeRepository
    fileStore     contracts.FileStore
    searchIndexer contracts.SearchIndexer
}
```

只有当构造参数明显膨胀、且多个 action 的依赖组高度重复时，再考虑增加聚合结构，例如：

```go
type ExternalContracts struct {
    FileStore         contracts.FileStore
    SearchIndexer     contracts.SearchIndexer
    FSRSClient        contracts.FSRSClient
    SourceProcess     contracts.SourceProcessClient
    Clock             contracts.Clock
}
```

但这一步建议后置。

---

## 11. 当前阶段最小可落地方案

如果按当前仓库状态推进，建议第一批先落这 6 个 contract：

1. `FileStore`
2. `FSRSClient`
3. `SourceProcessClient`
4. `Clock`
5. `SearchIndexer`（接口先行，实现可后置）

### 推荐原因

- `FSRSClient` 与 `SourceProcessClient` 能直接对应当前 Python 侧能力拆分
- `Clock` 很轻，但对测试和后续规则实现很有帮助
- `FileStore` 会很快在卡片 HTML 与资源落盘中变成刚需
- `SearchIndexer` 先定义边界，后接 Bleve 时不会再反向侵入 action

---

## 13. 一句话结论

**contracts 层的本质，是把“业务动作依赖的外部能力”收敛成稳定接口，让 actions 只表达业务意图，让 adapters 承担技术实现细节。**
