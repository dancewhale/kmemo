# Python Source Process 详细实现设计文档

## 1. 文档目标

本文档用于把 `docs/python/source-process-module-design.md` 的领域设计，进一步细化为**可直接指导 Python worker 落地实现**的代码结构设计。

本文件以一个新的前提为准：

- **仓库中的旧 `PrepareImportMaterial` proto 已被完全废弃和删除**
- Python worker 的 source-process 主路径只围绕 **job 化 gRPC 接口** 设计
- 不再保留任何“旧同步导入接口兼容层”

因此，本文档重点回答：

- 在当前 `python/app/` 代码结构下，job 化 source-process 应如何落地
- `server.py`、`services/import_service.py`、`source_process/*` 子包之间如何分层
- source-process 相关 proto 应如何按 job 模型组织
- Python worker 内部如何组织 converter、cleaner、registry、workspace、manifest、job 状态与事件
- 如何与现有 logging / request_id 基础设施对齐

本文档作为以下内容的 Python 落地依据：

- `proto/kmemo/v1/kmemo.proto` 中 source-process 相关 service / message
- `python/app/server.py`
- `python/app/services/import_service.py`
- Python worker 中未来的 `source_process/` 子包
- source-process 相关日志、request_id、错误语义、workspace 输出约定

---

## 2. 当前代码现状与新的设计前提

在设计详细实现之前，必须先基于仓库里的真实代码状态，同时明确新的设计前提。

### 2.1 当前 worker 入口结构仍然可保留

`python/app/server.py` 当前已经形成了清晰的薄入口模式：

- gRPC servicer 为 `KmemoProcessor`
- 每个 RPC 方法只负责：
  - 调用 `with_rpc_logging(...)`
  - 将请求委托给 `app.services.*`
- 业务逻辑不直接写在 servicer 中

这个结构性模式应继续保留。

source-process 的重构不应破坏这一点，而应把它演进为：

- servicer 方法对应 job 化的 source-process RPC
- handler 仍由 `services/import_service.py` 承接
- 真正业务逻辑继续下沉到 `source_process/*`

### 2.2 当前 import service 仍是占位实现

`python/app/services/import_service.py` 当前仍是 skeleton，尚未形成真实导入 pipeline。

这意味着：

- 可以直接按新前提重构，不需要迁就已经存在的复杂旧逻辑
- `services/import_service.py` 很适合作为 job 化 source-process 的 transport facade
- 当前是重建 source-process 边界的合适时机

### 2.3 当前日志与 request_id 基础设施可以直接复用

当前仓库已有两块可直接复用的基础设施：

#### request_id 上下文

`python/app/request_context.py`：

- `set_request_id(...)`
- `get_request_id()`
- `clear_request_id(...)`

#### RPC 日志包装

`python/app/grpc_logging.py`：

- 从 metadata 读取 `x-request-id`
- 绑定 `service=python-worker`
- 绑定 `component=grpc_server`
- 自动记录：
  - `rpc started`
  - `rpc finished`
  - `rpc failed`
  - `duration_ms`

source-process 的 job 化实现应直接复用这套机制，不重新实现另一套请求上下文或 RPC 日志系统。

### 2.4 新前提：旧 `PrepareImportMaterial` 已完全删除

本文档不再接受以下历史设计继续存在：

- `PrepareImportMaterialRequest`
- `PrepareImportMaterialResponse`
- `prepared_json`
- 基于 `raw_bytes` 的同步整包导入主路径

新的唯一前提是：

source-process 主路径只通过以下 job 化接口暴露：

- `SubmitImportJob`
- `GetJob`
- `ListJobEvents`
- `CancelJob`
- `GetCapabilities`

而 `CleanHtml` 仍保留为独立同步轻量能力，不并入 source-process job service 主路径。

---

## 3. 设计总原则

## 3.1 继续保持 thin server + thick service/facade

推荐继续采用现有模式：

- `server.py`：只做 gRPC 绑定、日志包装、service 分发
- `services/import_service.py`：承担 proto request/result → Python 内部模型之间的 facade 职责
- `source_process/*` 子包：承担真正的领域实现、converter 选择、cleaner pipeline、workspace 输出、manifest 写入、job 状态管理

不推荐：

- 在 `server.py` 中直接写导入流程
- 在 `services/import_service.py` 里塞入所有文件 IO、HTML 处理、转换器分派、manifest 生成、job 状态变更
- 让 proto message 成为长期内部模型

## 3.2 内部结构直接按 job 模型设计，不保留兼容层

由于旧 `PrepareImportMaterial` 已删除，Python 内部代码不需要任何过渡兼容层。

应直接围绕这些稳定领域概念组织：

- submit job
- job context
- job status
- job event
- source type
- conversion mode
- converter
- cleaner pipeline
- output writer
- manifest
- capability discovery

## 3.3 复杂逻辑拆成小模块，不做 god file

当前 worker 代码风格比较简洁，source-process 的落地也要保持这一点：

- 一个文件只负责一个清晰主题
- registry 单独放
- cleaner pipeline 单独放
- workspace/path 逻辑单独放
- manifest 生成单独放
- job store / event store 单独放
- converter 实现按 source type / mode 分拆

## 3.4 先固定边界，再逐步填实现

第一阶段最重要的是先把这些边界搭好：

- submit/get/list/cancel handler
- request / result model
- job runtime context
- converter protocol
- cleaner protocol
- registry
- output/manifest writer
- job status / event store
- 错误类型

再逐步补：

- HTML converter
- EPUB converter
- PDF converter
- metadata / structure extractor
- cleaner 具体实现
- 更完整的 job 生命周期与取消控制

## 3.5 Python 不直接接管 Go 的业务持久化

即使 source-process 已转为 job 化，Python 仍只负责：

- 接收 job 请求
- 执行导入和转换
- 产生 output 目录与 manifest
- 提供 worker 侧 job 状态与事件查询
- 返回结果摘要与 capability 信息

Python 不负责：

- 主业务 sqlite
- SourceDocument / Article / Card 的最终落库
- import_jobs 作为业务真相源的持久化

这些仍由 Go 主控层负责。

---

## 4. 推荐 proto 重设计

由于旧接口已经删除，source-process 相关 proto 应直接重设计为 job 化模型。

## 4.1 服务划分建议

建议继续保留单一 gRPC service 容器 `KmemoProcessor` 作为当前 proto 容器，但 source-process 相关 RPC 必须改为以下接口：

```proto
service KmemoProcessor {
  rpc SchedulerSetSetting(SchedulerSetSettingRequest) returns (SchedulerSetSettingResponse);
  rpc GetCardRetrievability(GetCardRetrievabilityRequest) returns (GetCardRetrievabilityResponse);
  rpc ReviewCard(ReviewCardRequest) returns (ReviewCardResponse);
  rpc RescheduleCard(RescheduleCardRequest) returns (RescheduleCardResponse);
  rpc OptimizeParameters(OptimizeParametersRequest) returns (OptimizeParametersResponse);
  rpc CleanHtml(CleanHtmlRequest) returns (CleanHtmlResponse);

  rpc SubmitImportJob(SubmitImportJobRequest) returns (SubmitImportJobResponse);
  rpc GetJob(GetJobRequest) returns (GetJobResponse);
  rpc ListJobEvents(ListJobEventsRequest) returns (ListJobEventsResponse);
  rpc CancelJob(CancelJobRequest) returns (CancelJobResponse);
  rpc GetCapabilities(GetCapabilitiesRequest) returns (GetCapabilitiesResponse);
}
```

后续若 proto 变复杂，再考虑独立 `SourceProcessService`；但 Python 内部实现不应依赖是否拆 service name。

## 4.2 核心 message 建议

### SubmitImportJobRequest

```proto
message SubmitImportJobRequest {
  string job_id = 1;
  string source_type = 2;
  optional string source_path = 3;
  optional string source_uri = 4;
  optional string source_url = 5;
  optional string raw_html = 6;
  string workspace_dir = 7;
  string output_dir = 8;
  string temp_dir = 9;
  ImportOptions options = 10;
  map<string, string> metadata = 11;
  optional string idempotency_key = 12;
}
```

### SubmitImportJobResponse

```proto
message SubmitImportJobResponse {
  string job_id = 1;
  string status = 2;
  bool accepted = 3;
}
```

### GetJobRequest / GetJobResponse

```proto
message GetJobRequest {
  string job_id = 1;
}

message GetJobResponse {
  SourceProcessJob job = 1;
}
```

### ListJobEventsRequest / Response

```proto
message ListJobEventsRequest {
  string job_id = 1;
  optional int64 after_sequence = 2;
}

message ListJobEventsResponse {
  repeated SourceProcessJobEvent events = 1;
}
```

### CancelJobRequest / Response

```proto
message CancelJobRequest {
  string job_id = 1;
}

message CancelJobResponse {
  string job_id = 1;
  string status = 2;
}
```

### GetCapabilitiesRequest / Response

```proto
message GetCapabilitiesRequest {}

message GetCapabilitiesResponse {
  repeated string source_types = 1;
  repeated string conversion_modes = 2;
  repeated ConverterCapability converters = 3;
  repeated CleanerCapability cleaners = 4;
}
```

### Job / Event / Result

```proto
message SourceProcessJob {
  string job_id = 1;
  string status = 2;
  string stage = 3;
  float progress = 4;
  ImportResult result = 5;
  optional string error_code = 6;
  optional string error_message = 7;
}

message SourceProcessJobEvent {
  string job_id = 1;
  int64 sequence = 2;
  string stage = 3;
  string message = 4;
  int64 created_at_unix = 5;
}

message ImportResult {
  string entry_html_path = 1;
  string cleaned_html_path = 2;
  string raw_text_path = 3;
  string manifest_path = 4;
  repeated string assets = 5;
  map<string, string> extracted_metadata = 6;
  string content_hash = 7;
  string effective_conversion_mode = 8;
  string converter_name = 9;
  string converter_version = 10;
  string cleaner_version = 11;
}
```

## 4.3 proto 设计原则

- 大文件主路径通过 `source_path/source_uri/workspace_dir/output_dir/temp_dir` 传递
- 只允许小型 HTML 文本走 `raw_html`
- 返回 job status / event / result 时使用 typed fields，不再嵌 JSON blob
- `GetCapabilities` 视为 source-process 领域能力发现接口，而不是临时调试接口

---

## 5. 推荐目录结构

建议在 `python/app/` 下新增 source-process 子包，而不是继续把所有逻辑塞进 `services/import_service.py`。

推荐结构如下：

```text
python/app/
├── server.py
├── logging_setup.py
├── grpc_logging.py
├── request_context.py
├── errors.py
├── services/
│   ├── fsrs_service.py
│   ├── html_service.py
│   └── import_service.py
└── source_process/
    ├── __init__.py
    ├── facade.py
    ├── models.py
    ├── enums.py
    ├── errors.py
    ├── capabilities.py
    ├── registry.py
    ├── workspace.py
    ├── manifest.py
    ├── pipeline.py
    ├── structure.py
    ├── metadata.py
    ├── job_store.py
    ├── event_store.py
    ├── cancellation.py
    ├── proto_conv.py
    ├── converters/
    │   ├── __init__.py
    │   ├── base.py
    │   ├── html_raw.py
    │   ├── html_readability.py
    │   ├── html_trafilatura.py
    │   ├── epub_spine.py
    │   ├── epub_chapters.py
    │   ├── pdf_pymupdf_flow.py
    │   ├── pdf_pymupdf_layout.py
    │   └── pdf_ocr.py
    └── cleaners/
        ├── __init__.py
        ├── base.py
        ├── parse_html.py
        ├── extract_main_content.py
        ├── remove_scripts.py
        ├── remove_forms.py
        ├── remove_unsafe_attrs.py
        ├── rewrite_assets.py
        ├── normalize_headings.py
        ├── normalize_whitespace.py
        ├── sanitize_html.py
        └── finalize_document.py
```

### 结构说明

#### `services/import_service.py`

只做 transport facade：

- 接收 proto request
- 调用 `source_process.facade`
- 返回 proto response

#### `source_process/facade.py`

负责 source-process 对外统一入口：

- submit/get/list/cancel/capabilities 五类主入口
- 调用 registry / pipeline / workspace / manifest / job store / event store
- 返回内部 result model

#### `source_process/job_store.py`

负责 Python worker 侧 job 运行态存储：

- job 状态
- 当前 stage
- progress
- result 摘要
- error_code / error_message

注意：这是 **worker 进程内运行态**，不是 Go 主业务真相源。

#### `source_process/event_store.py`

负责 worker 侧事件追加与查询：

- append event
- list events
- sequence 分配

#### `source_process/cancellation.py`

负责取消信号：

- job 是否被标记取消
- pipeline 中断检查
- converter / cleaner 长流程中途退出

---

## 6. 分层职责设计

## 6.1 `server.py`

职责保持最薄：

- 注册 gRPC servicer
- 绑定 source-process RPC 到 `services/import_service.py`
- 使用 `with_rpc_logging(...)`

source-process 不应在这一层做任何：

- 业务参数校验
- converter 选择
- 文件写入
- HTML 清洗
- manifest 生成
- job 状态机更新

## 6.2 `services/import_service.py`

这是 Python worker 的 transport facade。

职责：

- 接收 `kmemo_pb2.*` request
- 做最薄的一层 proto → internal model 转换
- 调用 `source_process.facade`
- 把内部结果转换回 proto response

建议直接暴露这些 handler：

```python
def submit_import_job(request: kmemo_pb2.SubmitImportJobRequest) -> kmemo_pb2.SubmitImportJobResponse:
    ...


def get_job(request: kmemo_pb2.GetJobRequest) -> kmemo_pb2.GetJobResponse:
    ...


def list_job_events(request: kmemo_pb2.ListJobEventsRequest) -> kmemo_pb2.ListJobEventsResponse:
    ...


def cancel_job(request: kmemo_pb2.CancelJobRequest) -> kmemo_pb2.CancelJobResponse:
    ...


def get_capabilities(request: kmemo_pb2.GetCapabilitiesRequest) -> kmemo_pb2.GetCapabilitiesResponse:
    ...
```

## 6.3 `source_process.facade`

这是真正的领域入口。

建议提供类似接口：

```python
class SourceProcessFacade:
    def submit_job(self, req: SubmitJobModel) -> SubmitJobResultModel:
        ...

    def get_job(self, job_id: str) -> JobStatusModel:
        ...

    def list_job_events(self, job_id: str, after_sequence: int | None = None) -> list[JobEventModel]:
        ...

    def cancel_job(self, job_id: str) -> CancelJobResultModel:
        ...

    def get_capabilities(self) -> CapabilityModel:
        ...
```

关键点是：

- service 层不直接知道 converter / cleaner 的实现细节
- facade 层组织 job 生命周期
- 具体处理细节交给 registry / pipeline / writer / extractor
- job store / event store 由 facade 驱动

## 6.4 job 执行方式

第一阶段建议仍采用 **进程内后台执行**：

- `SubmitImportJob` 接受请求并登记 job
- Python worker 在本进程线程池或后台执行单元中开始处理
- `GetJob` / `ListJobEvents` 查询运行态
- `CancelJob` 仅做 cooperative cancellation

这样可以先满足 job 化语义，同时避免一开始引入外部队列系统。

---

## 7. 内部数据模型建议

推荐优先使用 `dataclass`，与当前仓库风格保持一致。

### 7.1 输入模型

```python
@dataclass(frozen=True)
class ImportOptionsModel:
    conversion_mode: str
    fallback_modes: tuple[str, ...] = ()
    extract_main_content: bool = True
    sanitize_html: bool = True
    preserve_semantic_tags: bool = True
    download_remote_assets: bool = False
    inline_small_images: bool = False
    generate_toc: bool = False
    analyze_structure: bool = True
    keep_source_copy: bool = False
    enabled_cleaners: tuple[str, ...] = ()
    converter_params_json: str | None = None
```

```python
@dataclass(frozen=True)
class SubmitJobModel:
    job_id: str
    source_type: str
    source_path: str | None
    source_uri: str | None
    source_url: str | None
    raw_html: str | None
    workspace_dir: str
    output_dir: str
    temp_dir: str
    metadata: dict[str, str]
    options: ImportOptionsModel
    idempotency_key: str | None = None
```

### 7.2 运行时上下文

```python
@dataclass(frozen=True)
class JobContext:
    job_id: str
    request_id: str | None
    workspace_dir: Path
    output_dir: Path
    temp_dir: Path
    source_type: str
    logger: object
```

### 7.3 转换输出模型

```python
@dataclass(frozen=True)
class ConversionOutput:
    article_html: tuple[str, ...]
    cleaned_html: str
    raw_text: str | None
    asset_paths: tuple[str, ...]
    extracted_metadata: dict[str, str]
    effective_conversion_mode: str
    converter_name: str
    converter_version: str
```

### 7.4 job / event / result 模型

```python
@dataclass(frozen=True)
class ImportResultModel:
    entry_html_path: str
    cleaned_html_path: str
    raw_text_path: str | None
    manifest_path: str
    assets: tuple[str, ...]
    extracted_metadata: dict[str, str]
    content_hash: str
    effective_conversion_mode: str
    converter_name: str
    converter_version: str
    cleaner_version: str | None
```

```python
@dataclass(frozen=True)
class JobEventModel:
    job_id: str
    sequence: int
    stage: str
    message: str
    created_at_unix: int
```

```python
@dataclass(frozen=True)
class JobStatusModel:
    job_id: str
    status: str
    stage: str
    progress: float
    result: ImportResultModel | None
    error_code: str | None
    error_message: str | None
```

### 7.5 为什么推荐 dataclass

原因：

- 现有 `fsrs_service.py` 已经使用 dataclass 建模运行时与结果对象
- 风格统一
- 类型清晰
- 测试和构造成本低
- 当前阶段足够，不必过早引入更重校验框架

---

## 8. converter 设计

## 8.1 Converter 协议

建议定义统一 protocol：

```python
class Converter(Protocol):
    mode: str
    name: str
    version: str
    supported_source_types: set[str]

    def convert(self, ctx: JobContext, req: SubmitJobModel) -> ConversionOutput:
        ...
```

### 约束

- 一个 converter 只负责一种明确路线
- 不要在一个 converter 中兼容太多 source type
- `AUTO` 不由 converter 自己处理，而由 registry 决策后映射到具体 converter
- converter 内部要支持取消检查

## 8.2 第一批 converter 划分

建议按 source type + mode 拆分为小实现：

### HTML
- `html_raw.py`
- `html_readability.py`
- `html_trafilatura.py`

### EPUB
- `epub_spine.py`
- `epub_chapters.py`

### PDF
- `pdf_pymupdf_flow.py`
- `pdf_pymupdf_layout.py`
- `pdf_ocr.py`

### 原则

- 模式切换通过 registry 选择，不在上层写 if/else 链
- 每个实现返回统一 `ConversionOutput`
- 各 converter 只负责转换，不负责最终 manifest 写入

---

## 9. cleaner pipeline 设计

## 9.1 Cleaner 协议

```python
class Cleaner(Protocol):
    name: str
    version: str

    def clean(self, html: str, ctx: JobContext) -> str:
        ...
```

如果未来要保留更丰富中间态，可扩展为 `HtmlDocument` 模型，但第一阶段直接以 `str` 或轻量 document model 传递即可。

## 9.2 推荐 cleaner 顺序

与主设计文档保持一致：

1. `parse_html`
2. `extract_main_content`（可选）
3. `remove_scripts`
4. `remove_forms`
5. `remove_unsafe_attrs`
6. `rewrite_assets`
7. `normalize_headings`
8. `normalize_whitespace`
9. `sanitize_html`
10. `finalize_document`

## 9.3 Pipeline 组织方式

建议由 `source_process/pipeline.py` 统一组织：

```python
class CleanerPipeline:
    def __init__(self, cleaners: list[Cleaner]) -> None:
        self._cleaners = cleaners

    def run(self, html: str, ctx: JobContext) -> str:
        current = html
        for cleaner in self._cleaners:
            current = cleaner.clean(current, ctx)
        return current
```

### 原则

- 每个 cleaner 单一职责
- cleaner 顺序显式可读
- 启用/禁用由 options 或 registry 决定
- capability discovery 能反映有哪些 cleaner 可用
- 长流程步骤需要检查取消标记

---

## 10. registry 与 capabilities 设计

## 10.1 Registry 职责

`source_process/registry.py` 负责：

- 注册所有 converter
- 注册所有 cleaner
- 根据 `source_type + conversion_mode` 选择 converter
- 根据 `enabled_cleaners` 组装 cleaner pipeline
- 汇总 capability 信息

建议接口：

```python
class ConverterRegistry:
    def register(self, converter: Converter) -> None:
        ...

    def resolve(self, source_type: str, conversion_mode: str, fallback_modes: tuple[str, ...]) -> Converter:
        ...
```

```python
class CleanerRegistry:
    def register(self, cleaner: Cleaner) -> None:
        ...

    def build_pipeline(self, enabled_names: tuple[str, ...] | None = None) -> list[Cleaner]:
        ...
```

## 10.2 Capabilities

`source_process/capabilities.py` 负责：

- 枚举支持的 source types
- 枚举支持的 conversion modes
- 输出 converter 元信息
- 输出 cleaner 元信息
- 输出 converter 参数说明

它直接服务于 `GetCapabilities` RPC，不再作为任何旧接口的附属信息。

---

## 11. workspace 与输出写入设计

## 11.1 目录职责

建议由 `source_process/workspace.py` 统一管理：

- `workspace_dir`
- `output_dir`
- `temp_dir`
- `article/`
- `assets/`
- `manifest.json`
- `cleaned.html`
- `raw_text.txt`

不要把路径拼接散在 converter、cleaner、service 中。

## 11.2 推荐输出结构

与主设计文档保持一致：

```text
jobs/<job_id>/
  input/
  work/
  output/
    article/
    cleaned.html
    raw_text.txt
    manifest.json
    assets/
```

## 11.3 OutputWriter

建议增加统一 writer：

```python
class OutputWriter:
    def write_articles(self, ctx: JobContext, article_html: tuple[str, ...]) -> tuple[str, ...]:
        ...

    def write_cleaned_html(self, ctx: JobContext, html: str) -> str:
        ...

    def write_raw_text(self, ctx: JobContext, raw_text: str | None) -> str | None:
        ...

    def write_manifest(self, ctx: JobContext, manifest: ManifestModel) -> str:
        ...
```

这样可以保证：

- 文件命名规则集中
- output layout 集中
- manifest 中的相对路径统一

---

## 12. manifest 与结构分析设计

## 12.1 manifest 生成

建议由 `source_process/manifest.py` 统一负责：

- 计算 `content_hash`
- 记录 `effective_conversion_mode`
- 记录 `converter_name` / `converter_version`
- 记录 `cleaner_version`
- 记录资源路径
- 记录 metadata / structure 摘要

## 12.2 结构分析

建议由 `source_process/structure.py` 负责：

- headings 提取
- block 顺序提取
- source_ref 生成

这部分不要求长期稳定 block id，只要求：

- 可用于本次导入辅助分析
- 与主文档一致地提供 `Heading` / `ContentBlock` 风格结构

## 12.3 metadata 提取

建议单独放在 `source_process/metadata.py`：

- title
- author
- language
- source hints
- converter 衍生 metadata

不要和 manifest writer 或 cleaner 混写在同一文件。

---

## 13. job 生命周期设计

建议至少定义以下状态：

- `queued`
- `running`
- `succeeded`
- `failed`
- `cancel_requested`
- `cancelled`

建议至少定义以下 stage：

- `validate_input`
- `prepare_workspace`
- `select_converter`
- `convert_source`
- `run_cleaners`
- `extract_metadata`
- `analyze_structure`
- `write_output`
- `write_manifest`
- `finalize`

### 生命周期规则

- `SubmitImportJob` 创建 `queued` job
- 后台执行开始后切到 `running`
- 每个关键阶段都写事件并更新 progress
- 失败时切到 `failed` 并记录错误摘要
- 收到取消后切到 `cancel_requested`，在安全中断点切到 `cancelled`
- 成功完成后切到 `succeeded`

---

## 14. 错误设计

source-process 应单独定义错误，而不是全部抛裸 `ValueError`。

推荐错误类型：

```python
class SourceProcessError(Exception):
    pass

class ValidationError(SourceProcessError):
    pass

class UnsupportedModeError(SourceProcessError):
    pass

class ParseError(SourceProcessError):
    pass

class ConvertError(SourceProcessError):
    pass

class CleanError(SourceProcessError):
    pass

class OutputWriteError(SourceProcessError):
    pass

class JobNotFoundError(SourceProcessError):
    pass

class JobCancelledError(SourceProcessError):
    pass
```

### 错误映射原则

- `services/import_service.py` 负责把领域错误映射到 job 化 proto response / gRPC 错误
- `GetJob` / `ListJobEvents` 的不存在错误应稳定映射
- `SubmitImportJob` 的输入错误与运行期失败错误应分开表达
- 不再使用 `ok/message/prepared_json` 这种旧返回模式

### 日志原则

- 错误日志打在 facade / pipeline 边界
- 记录：
  - `service=python-worker`
  - `layer=worker` 或 `layer=adapter`
  - `component=source_process`
  - `job_id`
  - `source_type`
  - `conversion_mode`
  - `stage`
  - `error`

---

## 15. 与现有日志体系的集成

当前仓库已经有 `structlog + contextvars` 基础设施，source-process 应直接接入。

## 15.1 logger 获取方式

统一通过：

```python
from app.logging_setup import get_logger
```

在 facade / pipeline 中绑定：

- `service="python-worker"`
- `component="source_process"`
- `job_id`
- `source_type`
- `conversion_mode`
- `request_id=request_context.get_request_id()`

## 15.2 推荐日志事件

第一阶段优先覆盖：

- `job accepted`
- `job started`
- `job stage changed`
- `converter selected`
- `conversion started`
- `conversion finished`
- `cleaner pipeline started`
- `cleaner pipeline finished`
- `manifest written`
- `job cancel requested`
- `job cancelled`
- `job failed`
- `job succeeded`

### 注意

- 由 `grpc_logging.py` 负责 RPC 级 started/finished/failed
- source-process 内部只记录领域边界事件，不要每个小函数都打日志

---

## 16. 实现阶段建议顺序

为了匹配当前仓库状态，建议按最小闭环推进。

### 第一批（必须先落）

1. 新增 `python/app/source_process/` 基础目录
2. 增加：
   - `models.py`
   - `errors.py`
   - `registry.py`
   - `workspace.py`
   - `pipeline.py`
   - `manifest.py`
   - `job_store.py`
   - `event_store.py`
   - `cancellation.py`
   - `facade.py`
3. 把 `services/import_service.py` 改造成真正的 job 化 facade adapter
4. 在 `server.py` 上接入：
   - `SubmitImportJob`
   - `GetJob`
   - `ListJobEvents`
   - `CancelJob`
   - `GetCapabilities`
5. 先实现最小 HTML_RAW 路径
6. 先打通单机内 job 提交、状态查询、事件查询闭环

### 第二批（随后补）

1. HTML_READABILITY
2. HTML_TRAFILATURA
3. EPUB_SPINE
4. EPUB_CHAPTERS
5. Cleaner pipeline 细化
6. structure / metadata 提取
7. cooperative cancellation 完善

### 第三批（后续扩展）

1. PDF_PYMUPDF_FLOW
2. PDF_PYMUPDF_LAYOUT
3. PDF_OCR
4. 更细粒度 job event
5. 更稳定的后台执行与资源回收策略

---

## 17. 与现有 Python 风格的一致性要求

为了让这个实现看起来像当前仓库自然长出来的代码，而不是外来的一整套框架，必须遵守以下风格：

1. **继续使用函数式 service 入口**
   - 当前 `fsrs_service.py` / `html_service.py` 都是 module-level function
   - source-process 第一阶段也保持这一点

2. **内部核心模型优先用 dataclass**
   - 与 `fsrs_service.py` 保持一致

3. **不要引入过重框架**
   - 当前仓库风格偏轻量，不应突然引入大型 DI 容器或复杂插件系统

4. **按文件主题拆分，不做单文件巨物**
   - 这与现有 `server.py`、`grpc_logging.py`、`request_context.py` 的小文件风格一致

5. **复用现有 logging / request_context 基础设施**
   - 不重新实现另一套上下文系统

6. **内部按 job 化语义命名**
   - 不再保留 `prepare_import`、`prepared_json` 这类旧时代命名

---

## 18. 一句话结论

**在“`PrepareImportMaterial` 已完全删除”的前提下，Python 侧 source-process 应直接按 job 化 gRPC 能力落地：由 `server.py -> services/import_service.py` 维持薄 transport 入口，把 submit/get/list/cancel/capabilities 五类能力下沉到独立的 `source_process/` 子包中，并以 converter、cleaner、registry、workspace、manifest、job_store、event_store、facade 为核心分层组织，实现一个面向长期演进的任务式导入处理系统。**
