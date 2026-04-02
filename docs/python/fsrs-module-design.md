# Python FSRS 模块设计文档

## 1. 目标

本文档用于设计 kmemo 在 Python worker 中的 FSRS 能力实现方案，明确 Python 模块应如何集成 `py-fsrs` 与 `fsrs-optimizer`，并通过 gRPC 暴露稳定、可扩展、可理解的算法能力。

本次设计要回答的问题包括：

- Python worker 中 FSRS Scheduler 与 Optimizer 的职责边界是什么
- 为什么需要在 Scheduler 计算前显式执行 `Scheduler.setting(...)`
- 应该如何通过 gRPC 暴露：
  - `Scheduler.setting()`
  - `Scheduler.get_card_retrievability()`
  - `Scheduler.review_card()`
  - `Scheduler.reschedule_card()`
- Optimizer 能力应如何作为独立 gRPC 服务暴露
- Python worker 内部应如何组织 server / services / facade / converters
- 如何统一 Go contracts 与 Python worker 两侧的 FSRS 字段模型
- 如何删除未启用的旧 `CalculateFsrs` 接口，避免继续维护 JSON blob 协议

本文档的目标不是：

- 立即实现 FSRS 算法代码
- 重新设计完整的 review 业务流
- 让 Python worker 自己接管数据库或持久化职责
- 直接把第三方库 API 原样暴露给 gRPC 使用方

这份文档将作为后续修改以下内容的专项依据：

- `proto/kmemo/v1/kmemo.proto`
- `python/app/server.py`
- `python/app/services/fsrs_service.py`
- Python worker 中未来的 FSRS 子模块
- Go 侧后续的 adapter / facade 对接

---

## 2. 在整体架构中的位置

Python FSRS 模块位于当前 kmemo 分层中的算法服务一侧，而不是业务编排或持久化一侧：

```text
Wails UI API / Headless Go Host
    ↓
Actions / Flows / Contracts
    ↓
Go adapter / pyclient
    ↓ gRPC
Python worker server
    ↓
Python FSRS services / facade
    ├── py-fsrs Scheduler
    └── fsrs-optimizer Optimizer
```

### 一句话理解

- Go 侧负责业务编排、持久化、review log 管理、preset 选择
- Python 侧负责提供可复用的调度计算与参数优化能力

### 2.1 Python FSRS 模块负责什么

Python FSRS 模块负责：

- 将 gRPC typed request 转换成 Python 内部 FSRS 领域输入
- 调用 `py-fsrs` 执行 retrievability、review、reschedule 计算
- 调用 `fsrs-optimizer` 进行参数优化
- 将库返回结果映射为稳定的 gRPC typed response
- 将第三方库异常折叠为稳定的错误语义

### 2.2 Python FSRS 模块不负责什么

Python FSRS 模块不负责：

- 数据库存取
- `CardSRS` / `ReviewLog` 的持久化
- 事务提交
- 默认参数 preset 的选择策略
- review action 的业务编排
- Go contracts 的最终接口定义

这些职责仍然保留在 Go 侧的 repository / action / contract / adapter 中。

---

## 3. 当前代码现状

在决定 Python FSRS 模块如何设计之前，必须先面对当前仓库里的真实状态。

## 3.1 当前 proto 仍有待从旧占位接口迁移到 typed FSRS service

`proto/kmemo/v1/kmemo.proto` 历史上曾有一个粗粒度占位式 FSRS RPC：

```text
CalculateFsrs(item_id, payload_json) -> ok, message, result_json
```

但结合当前需求，这个接口应被视为需要直接移除的旧设计，问题包括：

- 输入输出 schema 被埋进 JSON blob 中
- scheduler 初始化参数无法作为独立能力表达
- Go/Python 两侧字段模型无法稳定对齐
- retrievability / review / reschedule 三种能力没有被拆开
- 未启用接口继续保留只会放大迁移噪音

因此，本文档不再将 `CalculateFsrs` 视为兼容层，而是将其定义为应被删除的旧接口；后续 Python FSRS 主路径只讨论新的 typed Scheduler/Optimizer service。

## 3.2 当前 Python server 是 thin dispatcher

`python/app/server.py:26` 当前实现了 `KmemoProcessorServicer`，并采用很薄的一层分发模式：

- servicer 方法只负责调用 `with_rpc_logging(...)`
- 业务逻辑委托给 `app.services.*`

这是一种应该保留的结构性优点，因为它天然适合继续扩展新的 FSRS service，而不需要把算法逻辑塞进 gRPC 层。

## 3.3 当前 FSRS service 还是旧 stub

`python/app/services/fsrs_service.py:8` 当前只有：

- `calculate_fsrs(request)`
- 固定返回 `ok=False`
- `message="TODO: FSRS not implemented"`

这说明：

- 当前 Python worker 还没有真实调度能力
- 也还没有形成内部领域模型与第三方库封装层
- 旧 `calculate_fsrs` 既未启用，也不应再作为演进目标继续补完
- 本次设计文档应优先解决新的 Scheduler/Optimizer 接口边界，并为后续删除旧 stub 提供依据

## 3.4 当前 Python 依赖还未引入 FSRS 相关库

`python/requirements.txt:3-7` 当前仅包含：

- `grpcio`
- `grpcio-tools`
- `protobuf`
- `pytest`
- `structlog`

当前还没有：

- `py-fsrs`
- `fsrs-optimizer`

这意味着文档必须显式讨论依赖引入策略、版本 pinning 和运行时边界。

## 3.5 当前 Go 侧 transport 仍是 raw generated client 暴露

`internal/pyclient/client.go` 当前暴露 raw generated client，这说明：

- 当前 Go/Python 之间仍是 transport-first 的过渡状态
- 后续若要让 Go contracts 真正依赖 typed capability，必须先在 Python 一侧把 typed RPC 设计稳定下来

---

## 4. 设计原则

## 4.1 按能力设计 gRPC，而不是按当前实现残留设计

推荐：

- 为 Scheduler 与 Optimizer 设计明确的 typed capability
- 让 RPC 名称表达业务能力

不推荐：

- 围绕 `CalculateFsrs(payload_json)` 继续堆更多 JSON schema

原因：后者会把过渡实现误当作长期接口，导致上层持续依赖 transport 细节。

## 4.2 Scheduler 与 Optimizer 必须拆成两个能力域

虽然二者都属于 FSRS 范畴，但职责不同：

- Scheduler：在线、单次、低延迟、围绕单卡状态计算
- Optimizer：离线、批处理、围绕 review history 拟合参数

若把二者混成一个 service：

- 接口会迅速失焦
- 运行时约束会相互污染
- 后续 Go adapter 也难以按能力封装

## 4.3 gRPC 暴露的是 kmemo 的稳定协议，而不是第三方库 API

文档应明确：

- `py-fsrs` 与 `fsrs-optimizer` 只是实现依赖
- gRPC request/result 必须由 kmemo 自己定义
- Python 内部可以调整映射方式，但不应破坏 gRPC 协议

这能保证：

- 未来换库时上层接口不必重写
- Python 库内部对象不会泄漏到 Go 或 proto
- 文档能够真正描述项目自己的能力边界

## 4.4 transport typed 化优先，旧 JSON blob 接口直接移除

当前占位式 `payload_json` / `result_json` 已不再作为迁移保留目标。

文档应锁定：

- 正式主路径使用 typed proto fields
- `Scheduler.setting(...)` 作为独立初始化能力暴露
- 旧 `CalculateFsrs` 不再兼容保留，而是直接删除
- 新增能力不再以 JSON blob 形式表达

## 4.5 Python worker 继续保持薄 server、厚 service/facade 的结构

推荐保持：

- `server.py` 只绑定 gRPC 与 logging wrapper
- `services/*` 负责 proto request/result 映射
- `fsrs/*` 子包负责领域模型、转换、校验与第三方库封装

不推荐：

- 在 servicer 方法中直接实例化 Scheduler / Optimizer 并写全部逻辑
- 让 proto generated message 直接成为长期内部模型

## 4.6 Python worker 只做算法能力，不直接读数据库

Optimizer 虽然需要大量 review history，但它的输入仍应来自 gRPC 显式传入，而不是让 Python worker 自己打开 SQLite。

原因：

- 保持当前 repo 的分层一致性
- 避免 Python worker 与 storage schema 深耦合
- 便于未来本地测试与替换实现来源

## 4.7 错误语义应稳定可判断，不直接透传库异常文本

调用方不应依赖：

- `py-fsrs` 内部异常文案
- `fsrs-optimizer` 的原始 traceback
- 任意 Python 报错字符串

应由 Python service/facade 把异常归类为稳定的错误语义，再映射到 gRPC status 与 response details。

---

## 5. 服务划分方案

## 5.1 当前可选方向

对于 proto 层，当前至少存在两种方向：

### 方案 A：继续把所有能力挂在 `KmemoProcessor`

优点：

- 改动入口少
- 与当前 proto 保持表面兼容

缺点：

- `KmemoProcessor` 会继续变成“杂项算法箱”
- Scheduler 与 Optimizer 的边界不清楚
- 后续 Go 侧 client/facade 也难以清晰命名

### 方案 B：新增专用 FSRS service

例如：

- `FsrsSchedulerService`
- `FsrsOptimizerService`

优点：

- 能力边界清晰
- 便于生成更聚焦的 Go/Python client 封装
- 文档、实现、测试都更容易围绕单一目标组织

缺点：

- 需要调整 proto 结构
- 迁移期需要同时考虑旧接口与新接口

## 5.2 当前阶段推荐结论

推荐采用 **方案 B**：

```text
service FsrsSchedulerService
service FsrsOptimizerService
```

原因：

- 用户需求本身已经明确区分了 Scheduler 方法和 Optimizer 方法
- 当前 `CalculateFsrs` 明显只是过渡接口
- 继续把 FSRS 复杂能力堆在 `KmemoProcessor` 上，只会延续 placeholder 设计

## 5.3 与旧接口的关系

推荐策略：

- 删除未启用的 `CalculateFsrs`
- 新的 typed Scheduler/Optimizer service 直接成为唯一主路径
- 后续 Go 新代码只对接新 service
- 不再为旧 JSON blob 协议保留兼容设计

---

## 6. Scheduler 能力设计

## 6.1 目标能力

结合 py-fsrs 的实际使用约束，用户当前希望 Scheduler 暴露以下能力：

- `Scheduler.setting()`
- `Scheduler.get_card_retrievability()`
- `Scheduler.review_card()`
- `Scheduler.reschedule_card()`

文档建议将它们在 proto 层对应为：

- `SettingScheduler`
- `GetCardRetrievability`
- `ReviewCard`
- `RescheduleCard`

其中 `SettingScheduler` 是本次新增的前置能力，不应再被内嵌进 `ReviewCard` 的隐式 JSON 参数里。

## 6.2 SettingScheduler

### 作用

在执行任意 scheduler 计算前，先基于当前 preset 初始化 Python 侧 scheduler。

这是因为 py-fsrs 的 scheduler 在计算前需要先完成初始化，典型设置项包括：

- `parameters`
- `desired_retention`
- `maximum_interval`

因此文档必须显式承认：scheduler 不是“零配置的纯函数入口”，而是“先 setting，再计算”的能力模型。

### 输入至少应讨论的字段

- `parameters`
- `desired_retention`
- `maximum_interval`

### 输出至少应讨论的字段

- 是否设置成功
- 当前生效配置摘要（可选）
- 可选的 `warnings[]`

### 边界说明

- `SettingScheduler` 只负责设置 scheduler，不负责任何单卡调度计算
- preset 的选择仍由 Go action / flow 决定
- Python service/facade 负责把 request 映射到 `Scheduler.setting(...)` 或等价初始化动作
- 如果相同参数重复下发，worker 内部可以做幂等优化，但这不改变接口语义

## 6.3 GetCardRetrievability

### 作用

给定当前卡片调度状态和当前时点，计算该卡片当前可回忆概率（retrievability）。

### 输入至少应讨论的字段

- `card_id`（可选，用于追踪与日志，不应作为算法必须字段）
- `state`
- `stability`
- `difficulty`
- `due`
- `last_review`
- `elapsed_days`
- `scheduled_days`
- `reps`
- `lapses`
- `now`

### 输出至少应讨论的字段

- `retrievability`
- `evaluated_at`
- 可选的 `warnings[]`

### 边界说明

- 这是纯计算能力，不负责持久化
- 计算依赖已完成的 scheduler setting
- 若输入状态不足以计算，应返回明确的 invalid/precondition 错误，而不是静默给默认值

## 6.4 ReviewCard

### 作用

给定当前卡状态、评分与 review 时间，计算下一次调度结果。

这与 `docs/contracts/fsrs-design.md` 中 `FSRSClient.Calculate` 的语义保持一致：

- 负责计算
- 不负责持久化

### 输入至少应讨论的字段

- `card_id`
- `state`
- `stability`
- `difficulty`
- `due`
- `last_review`
- `elapsed_days`
- `scheduled_days`
- `reps`
- `lapses`
- `rating`
- `reviewed_at`

### 输出至少应讨论的字段

- `card.state`
- `card.stability`
- `card.difficulty`
- `card.due`
- `review_log.rating`
- `review_log.review`
- `review_log.elapsed_days`
- `review_log.scheduled_days`
- 可选的 `retrievability`
- 可选的 `warnings[]`

### 边界说明

- 不创建 `ReviewLog`
- 不写数据库
- 不决定使用哪个 preset
- 只返回“给定输入后的调度计算结果”
- 返回字段命名与 Python 侧结果结构保持一致，Go adapter 不再自行发明另一套结果字段名

## 6.5 RescheduleCard

### 作用

对已有卡片状态进行重算，返回新的调度结果。

### 当前阶段建议

文档应先把 `RescheduleCard` 定义为**单卡重算能力**，而不是一开始就扩展到批量重排。

原因：

- 用户目前明确列出的 API 是单个 Scheduler 方法
- 批量重排会额外引入任务编排、吞吐量、响应大小等问题
- 当前设计文档应先固定最小稳定边界

### 输入至少应讨论的字段

- `card_id`
- 当前 card state snapshot
- `reschedule_at`
- 可选 reschedule mode

### 输出至少应讨论的字段

- `card.state`
- `card.stability`
- `card.difficulty`
- `card.due`
- `review_log.elapsed_days`
- `review_log.scheduled_days`
- 可选 diagnostics

### 边界说明

- v1 仅讨论单卡
- 计算依赖已完成的 scheduler setting
- 若未来确实需要批量 reschedule，再单独扩展新的 RPC，而不是让当前接口过度膨胀

---

## 7. Optimizer 能力设计

## 7.1 Optimizer 的定位

Optimizer 的职责不是代替 Scheduler，而是根据 review history 产出更合适的 FSRS 参数。

它应被视为：

- 离线能力
- 批处理能力
- 面向分析与训练的能力

而不是：

- 每次 review 时都会调用的在线低延迟路径

## 7.2 当前阶段推荐接口形状

当前阶段建议先收敛为一个主 RPC，例如：

- `OptimizeParameters`

不建议一开始就设计过多变体，例如：

- `EvaluateParameters`
- `CompareParameters`
- `TrainPresetForKnowledge`
- `OptimizeAndApply`

这些需求未来可以继续扩展，但目前文档应先把主路径锁住。

## 7.3 输入至少应讨论的字段

- `dataset_id` / `request_id`（可选）
- `review_logs[]`
- 当前 baseline parameters（可选）
- 优化配置（如最大迭代数、过滤选项、目标函数等）
- 数据集上下文（如 knowledge_id / preset_id，仅用于追踪，不应成为算法硬耦合字段）

### 关于 `review_logs[]`

文档必须明确 review log 的最小结构化字段，而不能只写“传原始 JSON”：

- `card_id`
- `reviewed_at`
- `rating`
- review 前卡状态相关字段
- review 后调度结果相关字段（若 optimizer 需要）
- 可选真实 recall outcome / labels（若库需要）

## 7.4 输出至少应讨论的字段

- `optimized_parameters`
- `metrics`
- `warnings[]`
- `sample_count`
- `dropped_sample_count`
- 可选 `diagnostics`

文档应强调：

- Optimizer 的主要产物是参数与指标
- 是否应用这些参数，仍属于 Go 侧业务决策

## 7.5 Optimizer 不应做什么

Optimizer 不应：

- 直接写回数据库中的 `FSRSParameter`
- 自己决定哪个 preset 成为默认值
- 直接绑定某个 knowledge 的持久化逻辑
- 在训练时自行读取 SQLite

---

## 8. typed 数据模型建议

## 8.1 枚举与合法值

为保持与 `docs/contracts/fsrs-design.md` 一致，文档建议统一以下稳定词汇：

### FSRSState

- `new`
- `learning`
- `review`
- `relearning`

### Rating

- `1 = Again`
- `2 = Hard`
- `3 = Good`
- `4 = Easy`

当前阶段可以在设计上继续允许 Go/Python 内部用整数或字符串映射，但 proto typed schema 应逐步收敛到明确的 enum/合法值集合，而不是继续依赖隐式 JSON 约定。

## 8.2 Go / Python 字段统一原则

本次设计要求 Go contracts、gRPC proto、Python worker 三层在 FSRS 结果模型上使用同一套语义字段，而不是各自维护一层“转换后命名”。

统一原则如下：

- Python 侧返回什么结构，Go 侧 contract 就应对齐什么结构
- Go adapter 负责 transport 映射，不负责重新发明字段语义
- `state`、`stability`、`difficulty`、`due`、`elapsed_days`、`scheduled_days` 等核心字段统一命名
- review 结果若天然包含 `card` 与 `review_log` 两部分，Go 侧也应按相同结构建模
- 不再继续使用 `next_fsrs_state`、`due_at` 这类仅存在于单侧文档中的过渡命名

建议 review 结果统一收敛为类似下面的结构语义：

```text
review_result {
  card {
    state
    stability
    difficulty
    due
  }
  review_log {
    rating
    review
    elapsed_days
    scheduled_days
  }
  retrievability
  warnings[]
}
```

这里的重点不是要求 proto 必须一字不差照搬某个 Python 库对象，而是要求 kmemo 自己定义的 typed schema 在 Go 与 Python 两侧保持同构，避免一边输出 `card.due`，另一边再包装成 `due_at_unix` 之类的二次模型。

## 8.3 时间字段

文档应明确统一原则：

- 正式 typed proto 主路径优先使用 `google.protobuf.Timestamp`
- 所有时间按 UTC 语义传输
- Python 内部统一转换为 timezone-aware datetime
- 空时间字段必须有显式 absent 语义，而不是使用 magic value

不推荐：

- 在 typed RPC 中继续把时间埋成 JSON string
- 让不同 RPC 自己决定时区语义

## 8.4 参数模型

文档应把 scheduler setting 视为独立 typed 结构，而不是继续把参数混入单卡调度请求。

当前阶段应明确：

- `SettingSchedulerRequest` 是 scheduler 初始化参数的正式主路径
- 其字段至少包括 `parameters`、`desired_retention`、`maximum_interval`
- Go 侧 `FSRSParameter` 需要先映射成该请求，再调用 scheduler setting
- `FSRSParameter.ParametersJSON` 只可作为持久化来源，不应继续成为 Go/Python RPC 之间的公开协议格式

## 8.5 warnings 与 diagnostics

文档建议对所有主要 RPC 预留：

- `warnings[]`
- 可选 `diagnostics`

原因：

- 这比继续使用 `ok/message` 更适合 typed 接口
- 可在不把错误升级为失败的情况下提示边缘情况
- 便于后续调试和 observability

---

## 9. Python 模块内部组织

## 9.1 推荐目录关系

建议后续实现按以下结构演进：

```text
python/app/
├── server.py
├── services/
│   ├── fsrs_scheduler_service.py
│   └── fsrs_optimizer_service.py
└── fsrs/
    ├── models.py
    ├── scheduler_facade.py
    ├── optimizer_facade.py
    ├── converters.py
    ├── validation.py
    └── errors.py
```

## 9.2 server.py 的职责

继续沿用 `python/app/server.py:26-55` 当前的薄层模式：

- 注册 gRPC servicer
- 将 RPC 调用委托给 service 层
- 通过 `with_rpc_logging(...)` 做统一的开始/结束/异常日志

不应在 `server.py` 中：

- 写业务转换逻辑
- 直接操作 `py-fsrs`
- 直接操作 `fsrs-optimizer`

## 9.3 services 层的职责

`services/*` 负责：

- 从 proto request 提取 typed 输入
- 调用 facade
- 把 facade 结果映射回 proto response
- 将领域错误映射到 gRPC 错误

这层是 transport 与领域逻辑之间的边界层。

## 9.4 facade 层的职责

`scheduler_facade.py` 与 `optimizer_facade.py` 负责：

- 封装第三方库 API
- 对外提供项目内部稳定方法
- 管理对象构建、参数映射、版本差异与返回结果整理

文档应明确：

- facade 是项目的内部稳定层
- gRPC 不应直接暴露第三方库对象

## 9.5 converters 与 validation 的职责

`converters.py` 负责：

- proto <-> 内部模型
- 内部模型 <-> 第三方库对象

`validation.py` 负责：

- rating/state 合法值校验
- 时间字段合法性校验
- optimizer dataset 最小约束校验

## 9.6 errors.py 的职责

`errors.py` 应定义 Python worker 内部可判断的错误类型，例如：

- invalid input
- insufficient data
- unsupported operation
- library failure
- internal error

其目的是避免：

- service 层散落大量字符串比较
- 调用方依赖第三方库异常文本

---

## 10. 错误语义与可观测性

## 10.1 错误分类建议

当前阶段建议至少统一以下类别：

- `invalid_argument`
- `failed_precondition`
- `unsupported`
- `internal`
- `cancelled`
- `deadline_exceeded`

### 对应典型场景

#### `invalid_argument`

- rating 非法
- state 非法
- 时间字段缺失或格式错误
- retrievability / optimizer 请求缺少必填字段

#### `failed_precondition`

- optimizer 输入样本不足
- card state 不足以进行当前计算
- 参数组合虽合法但不满足当前库要求

#### `unsupported`

- 当前阶段尚未支持的 reschedule mode
- 当前 worker 尚未支持的 optimizer 配置项

#### `internal`

- 第三方库抛出不可预期异常
- 内部转换逻辑错误
- response 构造失败

## 10.2 与 gRPC status 的关系

文档应明确：

- 新 typed RPC 不再主依赖 `ok/message/result_json` 风格
- 失败主路径应通过 gRPC status + 结构化 details 表达
- 如需保留 `warnings[]`，应仅用于“请求成功但有附加信息”的场景

## 10.3 日志与观测性

当前可继续复用 `python/app/grpc_logging.py:19-41` 的模式：

- 记录 `rpc started`
- 记录 `rpc finished`
- 记录 `rpc failed`
- 传递 `x-request-id`

后续文档应要求：

- server 层记录 RPC 级日志
- facade 层仅记录关键上下文与异常摘要
- 不在日志中直接打印大体积 review history 全量内容
- optimizer 请求应额外记录样本量、过滤量、耗时等摘要指标

---

## 11. 依赖管理

## 11.1 当前现状

`python/requirements.txt:3-7` 当前尚未包含：

- `py-fsrs`
- `fsrs-optimizer`

## 11.2 当前阶段建议

后续实现阶段应显式引入并固定版本，而不是使用无约束最新版本。

文档建议至少讨论：

- 依赖版本 pinning
- Python 版本兼容性
- Apple Silicon / grpcio 兼容性延续现有约束
- optimizer 是否需要延迟导入，以避免所有 worker 请求都承担其初始化成本

## 11.3 关于依赖隔离的保守建议

当前阶段不建议为了 Optimizer 单独再拆一个 Python 进程或独立服务。

更合适的做法是：

- 先在同一个 worker 进程内按模块隔离 Scheduler/Optimizer
- 若后续验证发现 Optimizer 的依赖、耗时或资源占用明显不同，再考虑进一步拆分部署

---

## 12. 与 Go contracts 的关系

本设计文档必须与 `docs/contracts/fsrs-design.md` 保持一致，而不是各自定义一套术语。

## 12.1 一致的边界

应继续保持：

- Go contracts 定义稳定 capability
- Python worker 是 capability 的一种实现来源
- Go actions / flows 不应直接依赖 proto generated details

## 12.2 一致的词汇

应保持一致的至少包括：

- `FSRSState`
- `Rating`
- scheduler setting / retrievability / review / reschedule 的输入语义
- review 只负责计算，不负责持久化
- optimizer 只负责产出参数，不负责应用参数

## 12.3 与 Go `FSRSClient` 的关系

结合 `docs/contracts/fsrs-design.md` 的最新设计，Go contracts 当前主张保留两类基础能力：

- `SetScheduler(...)`
- `Calculate(...)`

Python transport 侧可以设计得更细，但语义上应满足以下映射：

- Go `SetScheduler(...)` 对应 Python `SettingScheduler`
- Go `Calculate(...)` 可以由 Python `ReviewCard` 承载单卡 review 计算主路径
- 若 Go 后续需要 retrievability / reschedule 能力，可在 adapter 之上继续扩展 contracts，而不破坏当前基础接口

同时文档要求：

- Go `FSRSRequest` / `FSRSResult` 的字段命名应向 Python 侧 typed 结果收敛
- adapter 的职责是协议转换，不是重新设计领域模型
- Python 侧若返回 `card` / `review_log` 结构，Go 侧也应优先保持同构

---

## 13. 删除旧 `CalculateFsrs` 与迁移方案

## 13.1 当前接口的定位

当前 `CalculateFsrs(item_id, payload_json) -> ok/message/result_json` 应被明确标记为：

- 旧 placeholder
- 未启用接口
- 应直接删除的历史残留

## 13.2 推荐迁移顺序

推荐采用以下顺序：

1. 先完成本设计文档，固定 Python worker 中的 FSRS service 边界
2. 修改 `proto/kmemo/v1/kmemo.proto`，新增 typed Scheduler/Optimizer service 与 messages
3. 同时删除旧 `CalculateFsrs` 及其相关 message
4. 运行 `task proto` 生成新的 Go/Python 代码
5. 在 Python worker 中实现新的 Scheduler service
6. 在 Go 侧实现新的 adapter/facade，对接 typed RPC
7. 再实现 Optimizer service

## 13.3 为什么不能继续沿旧接口扩展

如果继续往 `payload_json` / `result_json` 中塞更多字段，将带来：

- Go/Python 两侧 schema 继续隐藏在字符串约定里
- scheduler setting 无法成为稳定独立能力
- retrievability / review / reschedule / optimize 的边界混乱
- 测试和演进都缺乏编译期约束
- adapter 很难对错误与字段做稳定折叠

因此，旧接口不再作为过渡层存在，而是应在新 typed service 落地时一起清理掉。

## 13.4 与 `KmemoProcessor` 的关系

迁移完成后：

- `KmemoProcessor` 不再承载 FSRS 的旧 `CalculateFsrs`
- 新设计主路径由专用 `FsrsSchedulerService` / `FsrsOptimizerService` 承担
- 若 `KmemoProcessor` 继续保留，也仅承载与 FSRS 无关的其他能力

---

## 14. 当前阶段的保守结论

当前阶段建议采用以下方案：

1. 在 Python 设计文档中正式区分 Scheduler 与 Optimizer 两个能力域
2. 为 Python worker 设计 typed `FsrsSchedulerService` 与 `FsrsOptimizerService`
3. 为 Scheduler 增加独立的 `SettingScheduler` 初始化接口
4. `GetCardRetrievability`、`ReviewCard`、`RescheduleCard` 先按单卡能力建模
5. `OptimizeParameters` 先作为 Optimizer v1 主 RPC
6. Go contracts、proto、Python worker 三层统一 FSRS 核心字段命名与结果结构
7. Python worker 内部通过 service + facade + converters + validation 隔离 gRPC 与第三方库
8. Python worker 不直接读数据库，也不直接决定 preset 持久化策略
9. 旧 `CalculateFsrs` 因未启用而直接删除，不做兼容保留
10. 依赖管理采用显式版本 pinning，而不是运行时临时解析

---

## 15. 后续演进方向

在不破坏当前边界的前提下，后续可以继续演进：

- 为 Scheduler 增加批量 reschedule 能力
- 为 Optimizer 增加 evaluate / compare / dry-run 能力
- 细化 optimizer metrics 与 diagnostics
- 在 proto 中引入更稳定的 enum/message 结构
- 在 Go 侧把 typed RPC 进一步收敛进更清晰的 adapter/facade
- 若 Optimizer 运行成本明显高于 Scheduler，再考虑进一步拆服务或任务化

这些演进都不应改变本文档已经固定的核心边界：

- Scheduler 与 Optimizer 分域
- scheduler 先 setting 再计算
- typed gRPC 为唯一主路径
- Python worker 封装第三方库而不是暴露第三方库
- Go 负责业务与持久化，Python 负责算法能力

---

## 16. 验证清单

本次是专项设计文档任务，因此验证以“边界是否清楚、术语是否一致、迁移是否可执行”为主。

## 16.1 架构边界检查

需要检查文档是否明确说明：

- Python worker 负责算法能力，不负责持久化
- Scheduler 与 Optimizer 是不同能力域
- server / services / facade / converters 的职责是否清楚

## 16.2 接口设计检查

需要检查文档是否完整讨论：

- `SettingScheduler`
- `GetCardRetrievability`
- `ReviewCard`
- `RescheduleCard`
- `OptimizeParameters`
- Go/Python 字段统一约束

## 16.3 迁移路径检查

需要检查文档是否清楚表达：

- 旧 `CalculateFsrs` 是待删除历史残留
- 新 typed service 是唯一主路径
- proto、Python、Go 三侧应按什么顺序演进

## 16.4 与现有代码对齐检查

需要检查文档是否与当前现实一致：

- `python/app/server.py` 是 thin dispatcher
- `python/app/services/fsrs_service.py` 仍是旧 stub
- `python/requirements.txt` 尚未引入 `py-fsrs` / `fsrs-optimizer`
- `docs/contracts/fsrs-design.md` 已固定 scheduler setting 与 typed contract 的总体方向

## 16.5 风格检查

需要检查文档是否符合仓库设计文档风格：

- 中文
- 编号分节
- 先目标、边界、原则，再接口、迁移、验证
- 不把实现细节代码堆砌成文档主体

---

## 17. 结论

当前 kmemo 在 Python 模块中最需要的，不是继续围绕旧 `CalculateFsrs(payload_json)` 做兼容设计，而是直接把 FSRS 边界收敛到新的 typed 能力模型上：

- 用 `SettingScheduler` 明确表达 py-fsrs 在计算前必须先初始化 scheduler
- 用 typed Scheduler/Optimizer gRPC 作为唯一主路径
- 用统一的 Go/Python 字段结构消除双边模型漂移
- 用清晰的 service/facade 分层隔离 transport、领域模型与第三方库
- 直接删除未启用的旧 `CalculateFsrs`，避免继续维护无效兼容层

这样后续无论是实现 Scheduler、接入 Optimizer、调整 proto，还是改造 Go adapter，整个 FSRS 能力都能沿着一条清晰、可落地、可维护的路线继续演进。
