# Source Process 端口与适配器设计（contracts + internal/adapters/sourceprocess）

## 1. 目标

- 将 **source-process 相关能力** 从“上层直接感知 Python worker / gRPC job 细节”收拢为：**明确的端口（port）+ 可替换的适配器（adapter）**。
- **上层（actions / flows / app）** 只依赖 **Go 领域侧类型** 与 source-process contract，**不直接持有** `*kmemov1.*`、gRPC 状态码、Python 内部目录结构等传输细节。
- 与 `docs/contracts-design.md` 对齐：source-process 在 contracts 层应被建模为**异步导入任务能力**，而不是旧的同步 `PrepareImportMaterial` 包装。
- 与 `docs/python/source-process-implementation-design.md` 对齐：Go 侧 contract 以 **Submit/Get/List/Cancel/Capabilities** 五类稳定能力为主边界，便于 Python worker、未来本地实现或其他 transport 替换。

本文档关注两件事：

1. Go 上层到底需要 source-process 提供哪些稳定能力；
2. 这些能力应如何映射到 `internal/adapters/sourceprocess` 中的具体实现。

---

## 2. 与整体分层的关系

```text
Actions / Flows
    ↓
contracts/sourceprocess
    ↓
adapters/sourceprocess
    ↓
外部实现（Python gRPC worker / future local worker）
```

source-process 与 repository 的关系仍然是并列：

```text
ImportDocumentFlow
   ├── repository        # 导入业务记录、知识库对象、文章/卡片/资产持久化
   └── contracts         # source-process 异步导入能力
```

一句话理解：

- `repository` 负责“导入后的业务数据怎么落库”；
- `contracts/sourceprocess` 负责“如何让外部处理器把源文件转成标准输出目录 + manifest”。

---

## 3. 设计前提

### 3.1 主路径是 job 化导入，不再沿用旧同步接口

本设计以 `docs/python/source-process-implementation-design.md` 的新前提为准：

- 旧 `PrepareImportMaterial` 已废弃，不再作为主设计的一部分；
- source-process 对 Go 的主边界是 **异步 job contract**；
- `CleanHtml` 作为 source-process contract 中的同步轻量能力保留，用于导入前预处理与 HTML 清洗；其余导入主流程仍以 job 能力为主。

因此，Go contracts 层不应再出现：

- `PrepareImportMaterialRequest/Response`
- `prepared_json`
- 以原始大字节同步传输整包导入内容的接口设计

### 3.2 Go 是业务真相源，Python 是处理执行者

source-process 的职责边界如下：

Python worker 负责：

- 接收导入 job 请求
- 运行 converter / cleaner / pipeline
- 写出标准输出目录与 manifest
- 提供 job 状态与事件查询
- 返回能力发现信息

Go 主控层负责：

- 决定何时发起导入
- 维护业务侧导入记录与最终持久化状态
- 读取 manifest / 输出目录
- 创建 `SourceDocument`、`Article`、`Card`、`Asset` 等业务对象

因此，contract 的返回值应聚焦于：**job 状态、结果摘要、输出位置、可判断错误语义**。

---

## 4. 推荐包结构

```text
internal/contracts/sourceprocess/
├── doc.go              # 包说明：source-process 是导入任务 port，不暴露 proto/gRPC 细节
├── client.go           # SourceProcessClient 接口定义
├── types.go            # Input/Output/Job/Event/Capabilities 等 contract 类型
└── errors.go           # contract 级错误语义

internal/adapters/grpcworker/
└── client.go           # 共享 gRPC 连接与生命周期管理

internal/adapters/sourceprocess/
├── grpc_python/
│   ├── client.go       # 对接 Python worker 的 job 化 gRPC 接口
│   ├── conv.go         # proto ↔ contracts/sourceprocess 类型映射
│   └── errors.go       # gRPC/status → contract 错误语义映射
├── noop/
└── inmemory/
```

说明：

- `contracts/sourceprocess` 只包含接口与纯 Go 类型，不 import proto。
- `grpcworker` 仅负责连接复用，不承载 source-process 语义。
- `adapters/sourceprocess/grpc_python` 是当前主实现，负责把 typed gRPC 请求/响应转换为 Go contract 类型。
- 若以后增加本地 worker 或测试用实现，只需实现同一 `SourceProcessClient`。

---

## 5. 端口（contracts）设计

### 5.1 核心接口

```go
type SourceProcessClient interface {
    CleanHtml(ctx context.Context, in CleanHTMLInput) (*CleanHTMLOutput, error)
    SubmitImportJob(ctx context.Context, in SubmitImportJobInput) (*SubmitImportJobOutput, error)
    GetJob(ctx context.Context, jobID string) (*Job, error)
    ListJobEvents(ctx context.Context, in ListJobEventsInput) ([]JobEvent, error)
    CancelJob(ctx context.Context, jobID string) (*CancelJobOutput, error)
    GetCapabilities(ctx context.Context) (*Capabilities, error)
}
```

说明：

- `CleanHtml`：执行同步 HTML 清洗，适用于导入前预处理、轻量 HTML 规范化与正文清理；
- `SubmitImportJob`：提交导入任务，返回是否受理与初始状态；
- `GetJob`：查询当前任务状态、阶段、结果摘要、错误信息；
- `ListJobEvents`：按顺序拉取事件流，供 Go 侧轮询、日志展示、增量同步；
- `CancelJob`：请求取消任务；
- `GetCapabilities`：查询支持的 source type / conversion mode / converter / cleaner。

这个接口粒度与 `docs/contracts-design.md:173` 的建议保持一致，但命名与类型进一步收敛为 Go contract 风格，便于上层直接使用。

### 5.2 输入输出类型

建议使用显式的 Input / Output 结构体，而不是在方法签名中散落大量 string 或 transport 对象。

#### CleanHTMLInput / CleanHTMLOutput

```go
type CleanHTMLInput struct {
    HTML            string
    BaseURL         *string
    SourceURL       *string
    ExtractMain     bool
    Sanitize        bool
    PreserveSemanticTags bool
    EnabledCleaners []string
    Metadata        map[string]string
}

type CleanHTMLOutput struct {
    HTML              string
    Title             *string
    TextContent       *string
    ExtractedMetadata map[string]string
    Warnings          []string
}
```

#### SubmitImportJobInput

```go
type SubmitImportJobInput struct {
    JobID          string
    SourceType     string
    SourcePath     *string
    SourceURI      *string
    SourceURL      *string
    RawHTML        *string
    WorkspaceDir   string
    OutputDir      string
    TempDir        string
    Options        ImportOptions
    Metadata       map[string]string
    IdempotencyKey *string
}
```

#### SubmitImportJobOutput

```go
type SubmitImportJobOutput struct {
    JobID    string
    Status   JobStatus
    Accepted bool
}
```

#### Job

```go
type Job struct {
    JobID        string
    Status       JobStatus
    Stage        string
    Progress     float32
    Result       *ImportResult
    ErrorCode    *string
    ErrorMessage *string
}
```

#### ListJobEventsInput / JobEvent

```go
type ListJobEventsInput struct {
    JobID         string
    AfterSequence *int64
}

type JobEvent struct {
    JobID         string
    Sequence      int64
    Stage         string
    Message       string
    CreatedAtUnix int64
}
```

#### CancelJobOutput

```go
type CancelJobOutput struct {
    JobID  string
    Status JobStatus
}
```

#### Capabilities

```go
type Capabilities struct {
    SourceTypes     []string
    ConversionModes []string
    Converters      []ConverterCapability
    Cleaners        []CleanerCapability
}
```

#### ImportResult

```go
type ImportResult struct {
    EntryHTMLPath           string
    CleanedHTMLPath         string
    RawTextPath             string
    ManifestPath            string
    Assets                  []string
    ExtractedMetadata       map[string]string
    ContentHash             string
    EffectiveConversionMode string
    ConverterName           string
    ConverterVersion        string
    CleanerVersion          string
}
```

#### ImportOptions

```go
type ImportOptions struct {
    ConversionMode       string
    FallbackModes        []string
    ExtractMainContent   bool
    SanitizeHTML         bool
    PreserveSemanticTags bool
    DownloadRemoteAssets bool
    InlineSmallImages    bool
    GenerateTOC          bool
    AnalyzeStructure     bool
    KeepSourceCopy       bool
    EnabledCleaners      []string
    ConverterParamsJSON  *string
}
```

### 5.3 状态类型建议

建议在 contract 层把常见 job 状态收敛为明确枚举型字符串，避免上层散落硬编码：

```go
type JobStatus string

const (
    JobStatusPending   JobStatus = "pending"
    JobStatusRunning   JobStatus = "running"
    JobStatusSucceeded JobStatus = "succeeded"
    JobStatusFailed    JobStatus = "failed"
    JobStatusCanceled  JobStatus = "canceled"
)
```

这样上层 flow 可以直接围绕业务语义判断，而不是依赖 proto enum 或 Python 常量名。

---

## 6. 设计约束与边界

### 6.1 大文件路径传递，不走 gRPC 大字节主路径

与 Python 设计文档一致，主路径通过这些字段表达输入与工作目录：

- `source_path`
- `source_uri`
- `workspace_dir`
- `output_dir`
- `temp_dir`

原因：

- 导入对象可能是 EPUB / PDF / 本地文件；
- 避免大文件经 gRPC 直接传输；
- 保持 Python worker 的处理方式与 Go 的工作目录组织解耦但兼容。

`RawHTML` 只作为小型 HTML 文本直传入口，不应扩展为通用大内容传输方案。

### 6.2 contract 返回“结果摘要”，不是直接返回业务实体

source-process 不应直接返回：

- `models.SourceDocument`
- `models.Article`
- `models.Card`
- repository 层的任何事务性结果

它应该返回的是：

- 标准输出目录相关路径
- manifest 路径
- 抽取出的元数据摘要
- 实际使用的 converter / cleaner / mode
- 任务状态与事件

这样 contracts 仍然保持“外部能力边界”，不越界到业务持久化。

### 6.3 capabilities 是正式能力发现接口

`GetCapabilities` 不是调试接口，而是稳定的能力发现入口，用于：

- 上层决定可选 `source_type` 与 `conversion_mode`
- 导入 UI 或 action 做参数校验
- 后续根据 converter / cleaner 名称做兼容判断

因此 contract 层需要为 converter / cleaner 设计清晰的能力结构，而不是返回临时 JSON blob。

---

## 7. 错误语义设计

建议在 `internal/contracts/sourceprocess/errors.go` 中统一定义 contract 级错误：

```go
var (
    ErrUnavailable  = errors.New("sourceprocess: processor unavailable")
    ErrInvalidInput = errors.New("sourceprocess: invalid input")
    ErrNotFound     = errors.New("sourceprocess: job not found")
    ErrConflict     = errors.New("sourceprocess: conflict")
    ErrCanceled     = errors.New("sourceprocess: canceled")
)
```

设计原则：

- 上层只判断 contract 错误语义；
- adapter 负责把 gRPC `status`、Python 领域错误、底层 IO 错误映射为这些语义；
- 不让 actions / flows 直接依赖 `codes.NotFound`、`codes.FailedPrecondition` 等 transport 细节。

可选地，也可以补一个错误包装类型：

```go
type ContractError struct {
    Op  string
    Err error
}
```

用于附带 `SubmitImportJob` / `GetJob` / `CancelJob` 等操作上下文。

---

## 8. 适配器（adapters/sourceprocess）职责

### 8.1 `grpc_python` 实现职责

当前推荐主实现为 `internal/adapters/sourceprocess/grpc_python`，职责如下：

1. 通过共享 `grpcworker` 连接调用 Python worker 的：
   - `CleanHtml`
   - `SubmitImportJob`
   - `GetJob`
   - `ListJobEvents`
   - `CancelJob`
   - `GetCapabilities`
2. 在 adapter 内完成：
   - request：contract types → `*kmemov1.*Request`
   - response：`*kmemov1.*Response` → contract types
3. 把 gRPC 错误映射为 contract 错误语义。

### 8.2 转换代码放置建议

| 方案 | 转换代码位置 | contracts 是否 import proto | 说明 |
|------|----------------|----------------------------|------|
| A（推荐） | `internal/adapters/sourceprocess/grpc_python/conv.go` | 否 | 与 FSRS 现有风格一致；contracts 只定义 Go 语义，adapter 负责 transport 映射 |

这样可以确保：

- 上层调用简洁；
- transport 协议变化时只影响 adapter；
- source-process 与 FSRS 的 contracts 风格保持一致。

### 8.3 与 `internal/contracts/fsrs` 的风格对齐

参考 `internal/contracts/fsrs/doc.go:1`、`internal/contracts/fsrs/scheduler.go:8`、`internal/contracts/fsrs/types.go:8`，source-process contract 也应遵守同样原则：

- 包注释明确“本包是 port，不暴露 protobuf / gRPC 类型”；
- 接口名直接表达业务能力；
- 输入输出使用纯 Go struct；
- 错误语义在本包内统一。

---

## 9. 上层调用约定（actions / flows）

典型导入调用链建议如下：

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
  -> repository.Asset.Create
```

约定：

- flow / action 只依赖 `contracts/sourceprocess`；
- 轮询策略、状态同步、最终入库编排由 Go 主控层决定；
- adapter 不负责业务事务、数据库回写、知识对象构建；
- 取消导入时，上层可同时更新业务侧 job 状态，但 Python worker 取消请求仍通过 `CancelJob(...)` 发送。

---

## 10. 测试与可替换性

建议测试分两层：

- **conv 单测**：仅验证 proto ↔ contract 类型映射；
- **集成测试**：启动真实 Python worker，验证 job 提交、状态查询、事件拉取、能力发现链路。

可替换性目标：

- 未来若增加 `adapters/sourceprocess/local` 或其他执行后端，只需实现同一 `SourceProcessClient`；
- 上层导入 flow 不需要因 transport 改变而改写。

---

## 11. 小结

- source-process 在 contracts 层的正确建模是：**异步导入任务 port**，而不是旧同步导入 RPC 的 Go 包装。
- `internal/contracts/sourceprocess` 应只暴露五类核心能力：`SubmitImportJob`、`GetJob`、`ListJobEvents`、`CancelJob`、`GetCapabilities`。
- 输入输出使用纯 Go 类型，围绕 job / event / manifest result / capabilities 建模，不暴露 proto 与 gRPC 细节。
- `internal/adapters/sourceprocess/grpc_python` 负责 typed gRPC 对接与错误、类型映射。
- Go 负责业务真相源与最终落库，Python 负责导入执行与标准输出生成；两者通过该 contract 解耦。
