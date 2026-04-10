# 日志系统设计说明文档

## 1. 目标

本文档用于为 kmemo 设计一套适合当前阶段、同时支持后续演进的日志系统。

kmemo 当前是一个 `Wails + Go + Python gRPC + SQLite` 的桌面应用骨架，因此日志系统的目标不是做成分布式观测平台，而是优先解决下面这些问题：

- 启动、关闭、配置加载、外部依赖连接过程可观察
- Go Host 与 Python Worker 之间的问题能快速定位
- action / flow 的主业务链路可追踪
- 数据库、文件、索引、gRPC 这几类外部能力错误能区分清楚
- 开发环境日志易读，后续生产环境日志可结构化收集
- 不把日志逻辑散落在各层，避免“想到哪里打到哪里”

---

## 2. 日志系统在整体架构中的位置

建议把日志视为一个**横切能力**，但使用上仍然遵守分层边界：

```text
UI API / Actions / Flows / Adapters / Repository
                ↓
             Logger
```

日志系统不是单独业务层，也不是 repository 的一部分。
它的作用是：

- 记录系统状态
- 记录关键动作
- 记录错误上下文
- 帮助定位跨层问题

但它**不负责**：

- 业务规则判断
- 错误恢复
- 指标统计
- 审计替代

---

## 3. 设计原则

## 3.1 先统一日志规范，再决定具体库

当前阶段最重要的是：

- 统一记录什么
- 在哪一层记录
- 字段怎么命名
- 错误如何附带上下文

而不是先纠结选哪一个日志库。

因此建议：

- Go：对外暴露统一 `Logger` 接口，底层默认采用成熟开源日志库 `uber-go/zap`
- Python：使用成熟第三方结构化日志库 `structlog`，并与标准库 `logging` 集成输出

这样上层业务代码只依赖统一接口和字段规范；底层实现则使用成熟、稳定、扩展性强的开源方案。

## 3.2 结构化日志优先，文本易读为辅

建议日志事件始终按“消息 + 字段”的形式组织，而不是只拼长字符串。

推荐：

```text
msg="submit review failed" card_id=123 action=submit_review capability=fsrs_scheduler error="timeout"
```

不推荐：

```text
submit review failed because fsrs client timeout for card 123
```

原因：

- 前者适合过滤、检索、后续导出
- 前者更容易跨 Go/Python 保持一致
- 前者更适合未来写文件、调试面板或上传

## 3.3 一条日志只表达一个事件

不要把多个状态、多个分支、多个错误揉进一条日志。

例如：

- `action started`
- `action finished`
- `grpc call failed`
- `file saved`

这些应该各自独立。

## 3.4 在边界打日志，不在细枝末节滥打日志

优先记录这些边界：

- 应用启动/关闭
- action 开始/结束/失败
- 调用外部能力前后
- 数据库事务失败
- Python worker 启动、RPC 失败
- 配置加载结果

不要默认给每个小函数都加日志，否则后续噪音会很大。

## 3.5 错误日志必须带上下文字段

错误日志不能只打 `err`。
至少应带：

- 发生在哪一层
- 当前动作是什么
- 关键业务对象 ID
- 外部能力名称（如果有）
- 可区分的操作名

---

## 4. 推荐日志分层策略

## 4.1 UI API 层

### 该层应该记录什么

- Wails 生命周期：startup / shutdown
- 重要 UI API 调用入口
- 参数校验失败
- 返回给前端的关键错误

### 不应该记录什么

- repository 细节
- 文件系统写入细节
- gRPC request/response 原始内容

### 建议字段

- `layer=ui_api`
- `method`
- `request_id`
- `knowledge_id` / `card_id` 等
- `duration_ms`

### 示例

```text
level=INFO msg="ui api request started" layer=ui_api method=CreateCard request_id=... knowledge_id=...
level=ERROR msg="ui api request failed" layer=ui_api method=CreateCard request_id=... error="invalid input"
```

---

## 4.2 Actions / Flows 层

这是最重要的日志层。

### 该层应该记录什么

- 动作开始 / 成功 / 失败
- 关键业务分支
- 是否进入事务
- 是否调用外部 contract
- 主要结果数量或结果 ID

### 建议字段

- `layer=action` 或 `layer=flow`
- `action` / `flow`
- `request_id`
- `knowledge_id` / `card_id` / `source_document_id`
- `duration_ms`
- `result_count` / `created_id`

### 示例

```text
level=INFO msg="action started" layer=action action=create_card request_id=... knowledge_id=...
level=INFO msg="action finished" layer=action action=create_card request_id=... card_id=... duration_ms=12
level=ERROR msg="action failed" layer=action action=submit_review request_id=... card_id=... error="fsrs unavailable"
```

### 原则

- action/flow 日志要能帮助你看清“这次业务动作做了哪些事”
- 不要把所有细节都放在 adapter 或 repository 层，否则主流程难追

---

## 4.3 Contracts / Adapters 层

### 该层应该记录什么

- 外部能力调用开始/结束/失败
- 调用目标信息（能力名、实现名、远端地址、路径）
- 能区分问题来源的上下文

### 建议字段

- `layer=adapter`
- `capability`
- `adapter`
- `op`
- `request_id`
- `target`
- `duration_ms`

### 能力示例

#### Python gRPC

- `capability=fsrs_scheduler`
- `capability=fsrs_optimizer`
- `capability=html_processor`
- `capability=source_process`
- `adapter=grpcworker`
- `target=127.0.0.1:50051`

#### FileStore

- `capability=file_store`
- `adapter=local_fs`
- `path`

#### SearchIndexer

- `capability=search_indexer`
- `adapter=bleve`
- `card_id`

### 示例

```text
level=INFO msg="contract call started" layer=adapter capability=fsrs_scheduler adapter=grpcworker op=set_scheduler request_id=...
level=INFO msg="contract call started" layer=adapter capability=fsrs_scheduler adapter=grpcworker op=calculate request_id=...
level=ERROR msg="contract call failed" layer=adapter capability=file_store adapter=local_fs op=save path="cards/1/content.html" error="permission denied"
```

### 原则

- adapter 日志用于定位“外部能力为什么出问题”
- 不要在这里记录高层业务结论，那是 action 的职责

---

## 4.4 Repository 层

### 该层应该记录什么

repository 层默认只建议记录：

- 慢查询
- 事务失败
- 明显异常的数据访问错误
- 批量操作的关键统计

### Debug 模式下额外记录什么

当显式打开 repository debug 开关时，可以额外记录：

- 普通 CRUD 成功日志
- 查询条件摘要
- 受影响行数
- 关键 SQL 执行耗时

但这类详细日志只应用于开发排查，不应作为日常默认输出。

### 不建议默认记录什么

- 每个普通 CRUD 成功日志
- 每次查询参数细节全量输出
- 正常路径上的高频低价值日志

### 建议字段

- `layer=repository`
- `repo`
- `op`
- `request_id`
- `duration_ms`
- `rows`
- `error`
- `debug_enabled`

### 原则

repository 不是业务主流程展示层。
默认模式下，repository 只保留慢查询和失败信息；详细日志通过单独 debug 参数开启，避免污染日常输出。

---

## 4.5 Bootstrap / Runtime 层

### 该层应该记录什么

- 配置加载
- Python worker 连接结果
- 数据库初始化结果
- Wails app 生命周期
- 关键依赖是否启用/跳过

### 建议字段

- `layer=bootstrap`
- `component`
- `python_grpc`
- `skip_python`
- `dial_timeout_ms`

### 示例

```text
level=INFO msg="bootstrap config loaded" layer=bootstrap python_grpc="127.0.0.1:50051" skip_python=false
level=INFO msg="python grpc connected" layer=bootstrap component=grpcworker target="127.0.0.1:50051"
level=WARN msg="python grpc skipped" layer=bootstrap component=grpcworker reason="KMEMO_SKIP_PYTHON=1"
```

---

## 5. 推荐日志级别规范

建议统一使用 4 个级别即可：

- `DEBUG`
- `INFO`
- `WARN`
- `ERROR`

当前阶段不建议引入更复杂的级别体系。

## 5.1 DEBUG

用于：

- 开发期细节追踪
- 输入输出摘要
- 非默认开启的排查信息

例如：

- grpc payload 大小
- 搜索重建批次大小
- action 内部关键分支

## 5.2 INFO

用于：

- 启动/关闭
- 关键业务动作开始/结束
- 外部依赖连接成功
- 明确的重要状态变化

## 5.3 WARN

用于：

- 系统仍可继续运行，但状态异常
- 外部依赖被跳过
- 回退到降级路径
- 非预期但已处理的异常

例如：

- Python worker 未连接但当前模式允许跳过
- 搜索索引暂时不可用，但主流程继续

## 5.4 ERROR

用于：

- 当前操作失败
- 外部能力调用失败且影响结果
- 数据不一致风险
- 未处理异常

原则：

- `ERROR` 要能对应“用户动作失败”或“系统能力失败”
- 不要把正常分支打成 ERROR

---

## 6. 推荐字段规范

建议 Go 和 Python 统一使用这些核心字段名：

### 通用字段

- `ts`：时间戳
- `level`
- `msg`
- `service`：`go-host` / `python-worker`
- `layer`
- `component`
- `request_id`

### 动作字段

- `action`
- `flow`
- `method`
- `op`

### 业务对象字段

- `knowledge_id`
- `card_id`
- `tag_id`
- `source_document_id`
- `review_log_id`
- `import_job_id`

### 外部能力字段

- `capability`
- `adapter`
- `target`
- `path`

### 性能字段

- `duration_ms`
- `count`
- `rows`
- `size_bytes`

### 错误字段

- `error`
- `error_kind`
- `grpc_code`

---

## 7. request_id / trace 设计建议

`request_id` 不再只是建议项，而应从现在开始作为 `context.Context` 的标准字段。

## 7.1 为什么需要

因为 kmemo 很多问题会跨越：

- Wails UI API
- Go action
- repository
- pyclient
- Python worker

没有统一 request_id 时，很难把一次用户动作串起来。

## 7.2 标准做法

每次 UI API 入口都生成一个 `request_id`，并立即写入 `context.Context`，后续所有层统一从 context 获取，不再各自单独生成。

```text
Desktop.CreateCard
  -> CreateCardAction
  -> SourceProcessor
  -> FileStore
  -> CardRepository
```

传递方式建议：

- Go：`request_id` 作为 `context.Context` 标准字段
- gRPC：通过 metadata 传给 Python worker
- Python：从 metadata 取出后写入日志上下文

## 7.3 最小实现要求

第一阶段就应做到：

- UI API 入口生成 request_id 并写入 context
- action / adapter / repository 统一从 context 中读取 request_id
- Python worker 收到 request_id 后写入日志
- 不允许在下层重新生成另一套 request_id


---

## 8. Go 侧日志接口设计建议

建议在 Go 侧定义一个轻量统一接口，而不是让所有层直接依赖某个具体日志库。

### 建议位置

- `internal/zaplog/`

### 建议接口

```go
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)

    With(fields ...Field) Logger
}
```

### Field 建议

```go
type Field struct {
    Key   string
    Value any
}
```

### 为什么建议这样做

- action/repository/adapter 不直接耦合具体库
- Go 底层默认采用 `uber-go/zap`，功能成熟、性能稳定、生态完善
- 测试里可以注入 fake logger
- 可以在 `With(...)` 中预绑定 `request_id`、`action`、`component`

---

## 9. Python 侧日志设计建议

Python 不建议继续只停留在标准库 `logging` 的默认用法，建议采用成熟第三方日志库 `structlog`，并与标准库 `logging` 集成输出。

这样做的原因：

- `structlog` 更适合结构化字段日志
- 与 Go 侧结构化日志风格更容易对齐
- 可以平滑接入标准输出、文件、JSON formatter
- 后续扩展上下文绑定（如 `request_id`）更自然

### 建议

- `server.py` 初始化统一的 `structlog` + `logging` 配置
- 为 `kmemo.worker`、`app.services.*` 使用固定 logger name
- 从 gRPC metadata 中提取 `request_id`
- 将 `request_id` 绑定到 Python 日志上下文
- 错误日志统一附带：
  - `service=python-worker`
  - `layer=adapter` 或 `layer=worker`
  - `capability`
  - `op`
  - `grpc_code`

### 示例

```text
level=INFO service=python-worker layer=worker component=grpc_server msg="worker started" target="[::]:50051"
level=ERROR service=python-worker layer=adapter capability=html_processor op=clean request_id=... error="invalid html"
```

---

## 10. 日志输出目标建议

当前阶段建议支持两种输出目标：

1. 标准输出（默认）
2. 本地日志文件（可选）

## 10.1 开发环境

默认：标准输出。

原因：

- 最简单
- 配合 `task run:go` / `task run:python` 最方便看
- 不会引入额外文件管理复杂度

## 10.2 桌面应用运行期

后续建议增加本地日志文件落盘：

```text
data/logs/
├── host.log
└── worker.log
```

建议：

- 按大小轮转
- 保留最近 N 个文件
- 不需要一开始就做远程上传

## 10.3 生产/发布阶段

如果未来需要问题上报或诊断导出，再考虑：

- 导出最近日志文件
- 用户主动附带日志提交 bug

不建议当前阶段就做自动上传。

---

## 11. 敏感信息与日志边界

日志里不应直接记录这些内容：

- 全量 HTML 正文
- 原始导入文件字节内容
- 未来若有 token / 密钥 / 本地隐私路径
- 大段 SQL / 大段 JSON payload

建议：

- 只记录摘要或长度
- 路径尽量记录相对路径
- HTML 只记录长度或片段摘要
- payload 记录 `size_bytes`

---

## 12. 慢操作日志建议

后续很容易变慢的环节：

- Python gRPC 调用
- 文件写入
- 搜索重建
- 批量导入
- 复杂数据库查询

建议对这些操作加 `duration_ms`，并定义慢操作阈值：

- gRPC 调用：> 300ms 记 WARN
- 文件写入：> 100ms 记 WARN
- DB 查询：> 200ms 记 WARN
- 导入/重建流程：记录总耗时 INFO

这类日志非常适合后面定位“功能没错但很慢”的问题。

---

## 13. 错误日志与用户错误的关系

建议明确区分两类错误：

## 13.1 用户可预期错误

例如：

- 输入不合法
- 目标知识库不存在
- 标签重复

建议：

- 记录 `WARN` 或较轻量 `INFO`
- 不要把所有这类错误都打成 `ERROR`

## 13.2 系统错误

例如：

- Python worker 连接失败
- 文件保存失败
- gRPC 返回 INTERNAL
- SQLite 操作失败

建议：

- 记录 `ERROR`
- 附带足够上下文

原则：

- “用户输错了”不等于系统故障
- “系统能力挂了”才是重点 ERROR

---

## 14. 推荐的日志事件清单

建议第一阶段优先覆盖这些事件：

### Host / Bootstrap

- `bootstrap config loaded`
- `python client connect started`
- `python client connected`
- `python client connect failed`
- `app startup`
- `app shutdown`

### UI API

- `ui api request started`
- `ui api request finished`
- `ui api request failed`

### Actions

- `action started`
- `action finished`
- `action failed`

### Adapters

- `contract call started`
- `contract call finished`
- `contract call failed`

### Repository

- `transaction failed`
- `slow query`
- `repository operation failed`

### Python Worker

- `worker started`
- `rpc handled`
- `rpc failed`

---

## 15. 当前阶段最小可落地方案

如果按项目当前骨架推进，建议分两步：

## 第一阶段

先建立最小统一规范：

- Go 定义 `internal/zaplog` 统一 logger 接口，底层实现采用 `uber-go/zap`
- bootstrap / app / pyclient 先接入日志
- UI API 入口生成 `request_id` 并写入 context 标准字段
- Python worker 统一初始化 `structlog`
- repository 默认只输出慢查询和失败日志，详细日志通过独立 debug 参数开启
- 统一字段：`service`、`layer`、`component`、`request_id`、`duration_ms`、`error`

## 第二阶段

随着 actions / contracts / adapters 落地，再补充：

- action / flow 标准起止日志
- gRPC metadata 透传 request_id
- 慢操作 WARN
- 本地日志文件轮转
- SearchIndexer / FileStore / SourceProcessClient / FSRSClient 的详细边界日志

---

## 16. 推荐目录结构

```text
internal/
├── zaplog/
│   ├── logger.go       # Logger 接口
│   ├── field.go        # Field 定义
│   ├── context.go      # request_id 标准字段 / logger from context
│   ├── zap_logger.go   # 基于 uber-go/zap 的默认实现
│   └── noop.go         # 测试或空实现
```

Python 侧可先保持：

```text
python/app/
├── logging_setup.py    # 初始化 structlog + logging
├── server.py
└── services/
```

---

## 17. 一句话结论

**kmemo 的日志系统应以“结构化、边界清晰、跨 Go/Python 可串联”为核心，让 action 能看主流程，让 adapter 能看外部能力，让 bootstrap 能看运行状态。**
