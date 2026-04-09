
* **同类任务支持切换不同转换方法**
* **Go 负责任务编排 / 主业务落库 / 文件与 sqlite 管理**
* **Python 负责导入、解析、清洗、转换**
* **Go 与 Python 通过 gRPC 通信**
* **内部代码结构清晰，遵从社区开源项目最佳实践，可读性高，可修改性强**

---
# 一、总体目标

* Go 作为主控服务，负责任务编排、文件系统管理、sqlite 主业务持久化、状态管理、错误归档。
* Python 作为内容处理服务，负责 HTML/EPUB/PDF 的导入、解析、转换、清洗、资源重写、结构分析。
* Go 与 Python 通过 gRPC 通信。
* 所有输入最终统一输出为标准 HTML 结果目录。
* 最终实现一个面向渐进阅读系统的“导入 / 清洗 / 转换”子系统：

1. 输入来源支持：

   * HTML
   * EPUB
   * PDF
   * 后续可扩展 URL

2. 统一输出为系统接受的 HTML 文件及其资源文件：

   * `article/1.html, article/1.1.html, article/2.html`
   * `manifest.json`
   * `assets/`

3. Go 与 Python 通过 gRPC 协作：

   * Go：任务编排、目录管理、sqlite 主业务持久化、错误归档、状态管理
   * Python：格式识别、导入转换、正文抽取、HTML 清洗、资源重写、结构分析

4. 每个任务支持指定“转换模式 / 转换器实现”：

   * 同一个 PDF 任务可切换不同处理路线
   * 同一个 HTML 任务可切换不同正文提取策略

---

# 二、核心设计原则

## 1. 明确边界

* **Go 是主控层**
* **Python 是内容处理引擎**
* Python **不直接写主业务 sqlite**
* Python 只负责把输入处理成标准输出目录和结构化结果

## 2. 面向扩展

* 导入器 / 转换器 / 清洗器均采用**可注册、可替换**设计
* 支持以后增加：

  * 新的 source type
  * 新的 converter
  * 新的 cleaner
  * 新的输出分析步骤

## 3. 易维护

* 模块职责单一
* 目录清晰
* 协议层、领域层、基础设施层分离
* 严格避免“大而全的 god object / god package / god file”

## 4. 易测试

* proto 层、服务层、转换器层、cleaner 层都应可独立测试
* 核心逻辑尽量避免强耦合 IO
* 复杂处理逻辑抽象为纯函数或小型 service，便于单测

## 5. 可读性优先

* 命名清晰
* 接口小而稳定
* 避免过度抽象
* 优先使用社区熟悉的模式和项目结构

---

# 三、非目标

以下内容**不作为本阶段实现目标**：

1. 不实现前端 UI
2. 不实现分布式队列
3. 不实现远程对象存储
4. 不实现 URL 抓取的完整浏览器渲染链路
5. 不实现跨版本稳定 block id
6. 不实现复杂的摘录定位回溯机制
7. 不让 Python 直接操作主业务数据库

---

# 四、关键约束

## 1. 输出 HTML 不保留稳定 block id

* 不要求在 HTML 节点中写入长期稳定的 `data-block-id`
* 可保留结构分析结果，但仅作为本次导入辅助信息
* `DocumentStructure.blocks` 中可包含顺序、类型、文本、来源引用，但不要求长期稳定 ID

## 2. 同类任务支持切换转换方式

必须支持每个任务通过参数指定转换方法，例如：

* HTML:

  * raw
  * readability
  * trafilatura

* EPUB:

  * spine
  * chapters

* PDF:

  * pymupdf_flow
  * pymupdf_layout
  * ocr

## 3. 大文件不通过 gRPC 直接传输

优先使用：

* `source_path`
* `source_uri`
* `workspace_dir / output_dir / temp_dir`

仅对小型 HTML 文本允许直接传 `raw_html`

---

# 五、系统整体架构

---

## 5.1 高层架构

### Go 侧职责

* 接收业务层导入请求
* 生成 `job_id`
* 创建任务目录
* 保存或组织输入文件
* 调用 Python gRPC 服务
* 跟踪任务状态
* 订阅任务事件
* 校验输出目录
* 读取结果
* 写入 sqlite 主业务表
* 记录错误与告警

### Python 侧职责

* 校验输入与 options
* 根据 `source_type + conversion_mode` 选择 converter
* 执行导入、解析、转换、清洗
* 输出标准 HTML 文件和 manifest
* 上报 job 状态与事件
* 提供 capabilities 查询

---

## 5.2 通信方式

采用 gRPC：

### 异步任务式接口

* `SubmitImportJob`
* `GetJob`
* `ListJobEvents`
* `CancelJob`

### 同步轻量接口

* `CleanHtml`

说明：`GetCapabilities` 虽然调用成本低，但语义上属于 source-process 领域的能力发现，更适合与任务式导入接口同属一个领域服务，而不是与 `CleanHtml` 归为同一类同步处理能力。

---

# 六、标准输出协议

每个任务一个工作目录：

```text
jobs/<job_id>/
  input/
    source.pdf | source.epub | source.html
  work/
    intermediate/
    logs/
  output/
    article/
    cleaned.html
    manifest.json
    assets/
      ...
```

## 输出语义

### `article/`

* 原始转换并进行章节切分后的 HTML目录,html文件按章节编号进行命名.
* 尽量保留语义内容
* 可包含必要结构整理

### `cleaned.html`

* 经过清洗后的 HTML
* 更适合系统内部阅读与渐进处理

### `manifest.json`

至少应包含：

* `job_id`
* `source_type`
* `entry_html_path`
* `cleaned_html_path`
* `raw_text_path`
* `assets`
* `metadata`
* `structure`
* `content_hash`
* `effective_conversion_mode`
* `converter_name`
* `converter_version`
* `cleaner_version`

### `assets/`

* 本地化后的图片等资源文件

---

# 七、gRPC 协议设计要求

使用上一轮讨论确定的 v1 proto 草案，要求满足：

## 7.1 必须包含的服务

* `SubmitImportJob`
* `GetJob`
* `ListJobEvents`
* `CancelJob`
* `CleanHtml`
* `GetCapabilities`

## 7.2 必须包含的核心字段

### `SubmitImportJobRequest`

* `job_id`
* `source_type`
* `source_path | source_uri | raw_html | source_url`
* `workspace_dir`
* `output_dir`
* `temp_dir`
* `ImportOptions`
* `metadata`
* `idempotency_key`

### `ImportOptions`

* `conversion_mode`
* `fallback_modes`
* `extract_main_content`
* `sanitize_html`
* `preserve_semantic_tags`
* `download_remote_assets`
* `inline_small_images`
* `generate_toc`
* `analyze_structure`
* `keep_source_copy`
* `enabled_cleaners`
* `html / epub / pdf options`
* `converter_params`

### `ImportResult`

* `entry_html_path`
* `cleaned_html_path`
* `raw_text_path`
* `manifest_path`
* `assets`
* `extracted_metadata`
* `structure`
* `content_hash`
* `effective_conversion_mode`
* `converter_name`
* `converter_version`
* `cleaner_version`

## 7.3 必须支持的转换模式

要求定义统一 `ConversionMode`，至少包括：

* `AUTO`
* `HTML_RAW`
* `HTML_READABILITY`
* `HTML_TRAFILATURA`
* `EPUB_SPINE`
* `EPUB_CHAPTERS`
* `PDF_PYMUPDF_FLOW`
* `PDF_PYMUPDF_LAYOUT`
* `PDF_OCR`

## 7.4 能力发现接口

`GetCapabilities` 必须返回：

* 支持的 source types
* 支持的 conversion modes
* converter 元信息
* cleaner 元信息
* converter 参数说明

---

# 八、Go 侧设计方案

---

## 8.1 Go 侧项目目标

实现一个清晰分层的服务端模块，负责：

* 导入任务创建
* 状态管理
* 调用 Python pipeline
* 结果落库
* 失败处理
* 可扩展的业务集成

---

## 8.4 Go 侧核心接口建议

### JobRepository

```go
type JobRepository interface {
    Create(ctx context.Context, job *Job) error
    UpdateStatus(ctx context.Context, jobID string, status JobStatus, stage Stage, progress float32) error
    GetByID(ctx context.Context, jobID string) (*Job, error)
    SaveResult(ctx context.Context, result *JobResult) error
}
```

### PipelineClient

```go
type PipelineClient interface {
    SubmitImportJob(ctx context.Context, req *SubmitImportRequest) (*SubmitImportResponse, error)
    GetJob(ctx context.Context, jobID string) (*PipelineJobStatus, error)
    ListJobEvents(ctx context.Context, jobID string) (<-chan PipelineJobEvent, <-chan error)
    CancelJob(ctx context.Context, jobID string) error
    CleanHtml(ctx context.Context, req *CleanHTMLRequest) (*CleanHTMLResponse, error)
    GetCapabilities(ctx context.Context) (*Capabilities, error)
}
```

### WorkspaceManager

```go
type WorkspaceManager interface {
    PrepareJobWorkspace(ctx context.Context, jobID string) (*WorkspacePaths, error)
    ValidateOutput(ctx context.Context, jobID string) error
}
```

---

## 8.5 Go 侧实现要求

1. 使用 `context.Context`
2. 错误必须 wrap，保留上下文
3. 领域错误与基础设施错误分开
4. 业务事务边界清晰
5. 日志结构化
6. 不允许在 handler 中写复杂业务逻辑
7. 不允许 repository 中混入业务规则
8. 每个 package 文件数与职责保持合理，不堆砌到单文件

---

# 九、Python 侧设计方案

---

## 9.1 Python 侧目标

实现一个内容处理服务，负责：

* 接收 Go 提交的任务
* 选择转换器
* 执行导入 / 转换 / 清洗 / 输出
* 维护任务状态
* 返回结果
* 提供能力发现接口

---

## 9.3 Python 侧核心抽象

### Converter

每个转换器只负责一种路线：

```python
class Converter(Protocol):
    mode: ConversionMode
    name: str
    version: str
    supported_source_types: set[SourceType]

    def convert(self, ctx: JobContext) -> ConversionOutput:
        ...
```

### Cleaner

每个 cleaner 只负责一个小步骤：

```python
class Cleaner(Protocol):
    name: str
    version: str

    def clean(self, doc: HtmlDocument, ctx: CleanContext) -> HtmlDocument:
        ...
```

### Registry

不要在业务代码里写大段 `if source_type == ... and mode == ...`

而是使用：

* `ConverterRegistry`
* `CleanerRegistry`

---

## 9.4 Python 侧处理流程

### 提交任务后

1. 校验输入
2. 构建 job context
3. 选择 converter
4. 执行转换
5. 执行 cleaner pipeline
6. 提取 metadata / structure / assets
7. 写出 HTML 与 manifest
8. 更新任务状态
9. 返回结果

### 选择 converter 的规则

1. 若 `conversion_mode != AUTO`，直接选择指定 mode
2. 若 `AUTO`，由系统按 source type 选择默认 converter
3. 若失败，依次尝试 `fallback_modes`
4. 最终记录 `effective_conversion_mode`

---

## 9.5 Python 侧实现要求

1. 使用类型注解
2. 业务数据模型明确，优先 dataclass / pydantic
3. 不在 handler 中写复杂逻辑
4. 避免单个文件过大
5. 复杂清洗流程拆成多个 cleaner
6. 所有转换器都通过统一协议暴露
7. 异常分层：

   * ValidationError
   * UnsupportedModeError
   * ParseError
   * ConvertError
   * CleanError
   * OutputWriteError

---

# 十、导入与转换实现策略

---

## 10.1 HTML

支持三种模式：

### HTML_RAW

* 只做基础解析与规范化
* 不做正文抽取

### HTML_READABILITY

* 使用 readability 类算法提取正文

### HTML_TRAFILATURA

* 使用 trafilatura 提取正文

### 推荐技术

* `lxml`
* `beautifulsoup4`
* `readability-lxml`
* `trafilatura`
* `bleach`

---

## 10.2 EPUB

支持两种模式：

### EPUB_SPINE

* 按 spine 顺序合并内容
* 输出单一主 HTML

### EPUB_CHAPTERS

* 保留章节边界后输出单 HTML
* 可在 HTML 中用 `section` 反映章节结构

### 推荐技术

* `ebooklib`
* `lxml`
* `beautifulsoup4`

---

## 10.3 PDF

支持三种模式：

### PDF_PYMUPDF_FLOW

* 优先面向阅读
* 将文本块重组为流式 HTML

### PDF_PYMUPDF_LAYOUT

* 尽量保留布局信息
* 适合版面重要的文档

---

# 十一、清洗流水线设计

必须实现 **pipeline + 小步骤 cleaner** 结构，不允许写成一个巨大函数。

## 推荐 cleaner 顺序

1. parse_html
2. extract_main_content（可选）
3. remove_scripts
4. remove_forms
5. remove_unsafe_attrs
6. rewrite_assets
7. normalize_headings
8. normalize_whitespace
9. sanitize_html
10. finalize_document

## 约束

* 每个 cleaner 单一职责
* cleaner 可配置启用或禁用
* cleaner 注册表支持能力发现

---

# 十二、结构分析要求

虽然不保留稳定 block id，但仍应支持基本结构分析：

## `DocumentStructure`

* headings
* blocks

## `Heading`

* level
* text
* anchor
* order

## `ContentBlock`

* order
* block_type
* text
* source_ref

## 说明

* `source_ref` 可用于保留弱引用，例如：

  * `page:3`
  * `chapter:2`
  * `xpath:/html/body/...`
* 不要求长期稳定

---

# 十三、sqlite 与持久化约束

主业务数据库由 Go 管理。

## 建议表

### import_jobs

* job_id
* source_type
* source_path
* status
* stage
* progress
* error_code
* error_message
* created_at
* updated_at

### sources

* source_id
* job_id
* original_path
* content_hash
* mime_type
* metadata_json

### articles

* article_id
* source_id
* title
* entry_html_path
* cleaned_html_path
* raw_text_path
* language
* author
* structure_json
* created_at

### assets

* asset_id
* article_id
* local_path
* mime_type
* sha256
* size

---

# 十四、错误处理设计

要求实现清晰的错误模型。

## Go 侧

* transport error
* application error
* domain error
* persistence error

## Python 侧

* validation error
* unsupported converter
* parse failure
* clean failure
* output failure

## proto 层

通过 `JobError` 返回：

* code
* stage
* message
* retryable
* details

## 要求

* 不允许只返回一个裸字符串错误
* 错误信息要可审计、可定位、可日志检索

---

# 十五、日志与可观测性

## 基本要求

* 使用结构化日志
* 日志中必须带：

  * job_id
  * source_type
  * conversion_mode
  * stage

## Python 任务事件

必须通过 `ListJobEvents` 输出阶段事件，例如：

* task accepted
* converter selected
* parse started
* clean started
* write output completed

## 禁止

* 在业务日志中打印过大原始 HTML 正文
* 在错误日志中泄露敏感本地路径以外的不必要信息

---

# 十八、代码风格与工程规范

这是提交给 AI 生成代码时必须强调的部分。

## Go 代码要求

1. 遵循 Go 社区最佳实践
2. 小接口优先
3. 组合优于继承风格抽象
4. 依赖倒置
5. context 贯穿服务层
6. 显式错误处理
7. 避免 util 大杂烩
8. 包职责聚合，不循环依赖
9. 不生成过度复杂的泛型结构
10. public API 要有注释

## Python 代码要求

1. 遵循 Python 社区最佳实践
2. 使用类型注解
3. 优先 dataclass / pydantic 表达结构
4. 模块拆分合理
5. 不写超长函数
6. 不写大而杂的 manager 类
7. 不滥用动态特性
8. 公共类和函数应有 docstring
9. 异常类型明确
10. 避免隐式全局状态

## 通用要求

1. 可读性优先于炫技
2. 便于后续添加新 converter / cleaner
3. 命名一致
4. 文件组织稳定
5. 生成代码时附带必要注释，但不要冗长废话

---

# 十九、AI 生成代码时的强约束

把下面这些要求直接作为生成指令的一部分：

1. **先生成目录结构与核心接口，再生成实现**
2. **先生成 proto 与模型，再生成服务**
3. **严格分层，不要把业务逻辑堆到 handler**
4. **不要把所有逻辑塞进一个文件**
5. **所有 converter 都必须通过 registry 注册**
6. **所有 cleaner 都必须通过 pipeline 编排**
7. **Python 不允许直接写主 sqlite**
8. **Go 负责最终结果落库**
9. **任务必须支持 conversion_mode 与 fallback_modes**
10. **不实现稳定 block id**
11. **输出目录协议固定**
12. **需要提供基础单元测试骨架**
13. **需要提供 README 或模块注释说明扩展方式**
14. **需要保留可扩展点，但不要过度工程化**

---

# 二十一、可直接提交给 AI 的最终指令文本

下面这段你可以直接作为提示词发给 AI 生成代码。

---

## AI 代码生成指令

请为一个“渐进阅读系统的导入/清洗/转换子系统”生成工程代码，采用 **Go + Python + gRPC** 架构，并严格遵守以下设计要求。

### 二、关键约束

* **不在输出 HTML 中保留稳定 block id。**
* **每个任务都必须支持指定 conversion_mode，用于切换相同任务的不同转换方法。**
* 支持 `fallback_modes`，在主转换方法失败时可依次尝试备用方法。
* Python **不得直接写主业务 sqlite**。
* Go 负责最终结果落库。
* 大文件通过文件路径/URI 传递，不通过 gRPC 直接传输大文件字节。
* 输出目录协议固定为：

  * `article.html`
  * `cleaned.html`
  * `manifest.json`
  * `assets/`

### 三、必须支持的 source type

* HTML
* EPUB
* PDF

### 四、必须支持的 conversion mode

统一定义 `ConversionMode`，至少包括：

* AUTO
* HTML_RAW
* HTML_READABILITY
* HTML_TRAFILATURA
* EPUB_SPINE
* EPUB_CHAPTERS
* PDF_PYMUPDF_FLOW
* PDF_PYMUPDF_LAYOUT
* PDF_OCR

### 五、gRPC 服务

必须实现以下接口：

* SubmitImportJob
* GetJob
* ListJobEvents
* CancelJob
* CleanHtml
* GetCapabilities

### 六、实现要求

#### Go

* 使用清晰分层结构：domain / application / infrastructure / interfaces
* repository interface 与实现分离
* grpc client 单独封装
* workspace/file system 抽象单独封装
* 所有业务流程使用 context.Context
* 错误处理清晰、可 wrap
* 日志采用结构化方式

#### Python

* 使用清晰分层结构：api / core / services / converters / cleaners / extractors / io / utils
* 所有 converter 必须通过 registry 注册
* 所有 cleaner 必须通过 pipeline 组合
* 使用类型注解
* 公共接口和类应有 docstring
* 错误类型明确分层

### 七、推荐技术

#### Go

* grpc-go
* sqlite
* 标准库 + 少量必要依赖
* 不要引入过重框架

#### Python

* grpcio / grpc.aio
* lxml
* beautifulsoup4
* bleach
* readability-lxml
* trafilatura
* ebooklib
* pymupdf
* OCR 路径预留为可扩展模式

### 八、结构分析

* 允许生成 `DocumentStructure`
* 包括 headings 与 blocks
* blocks 可包含：

  * order
  * block_type
  * text
  * source_ref
* **不需要稳定 block id**

### 九、代码质量要求

* 内部代码结构清晰
* 遵从社区开源项目最佳实践
* 可读性高
* 可修改性强
* 避免单文件过大
* 避免 god object / god package
* 优先简单、清晰、可测试的实现
* 保留扩展点，但不要过度工程化

### 十、生成顺序

请分阶段生成代码：

1. 先生成目录结构、proto、核心模型与接口
2. 再生成 Go/Python 服务骨架
3. 再生成 registry、pipeline、基础 converter/cleaner
4. 再生成 EPUB/PDF converter
5. 最后补测试骨架与 README

### 十一、输出内容

请输出：

1. 项目目录结构
2. proto 文件
3. Go 侧核心代码骨架
4. Python 侧核心代码骨架
5. 基础 HTML converter 与 cleaner pipeline 示例
6. README 中的扩展说明
7. 测试骨架

