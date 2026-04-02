# Contracts FSRS 设计文档

## 1. 目标

本文档用于细化 `docs/contracts-design.md` 中 `FSRSClient` 相关部分，明确 kmemo 在 contracts 层中应该如何定义 FSRS 调度能力的边界、术语、字段来源、协作方式与错误语义。

目标不是在这份文档里实现 FSRS 算法，也不是重新设计完整的复习系统，而是先把下面这些问题固定下来：

- 上层业务到底向 FSRS 能力要什么结果
- 哪些字段属于稳定的 contracts request/result
- 当前 proto 里的 `payload_json` / `result_json` 应被看作什么层级的细节
- `CardSRS`、`FSRSParameter`、`ReviewLog` 与 FSRS contract 应如何协作
- `SubmitReview`、`PreviewSchedule`、`UndoLastReview` 哪些应该调用 FSRS，哪些不应该调用
- adapter 应如何隔离 Python gRPC transport 与未来本地实现
- actions 应依赖什么错误语义，而不是依赖 gRPC 或 Python 细节

这份文档是后续设计 `internal/contracts/fsrs.go`、`internal/adapters/pyclient/fsrs_client.go` 以及 review action / flow 的专项依据。

---

## 2. 在整体分层中的位置

FSRS 能力仍然遵循 `docs/contracts-design.md` 中的 contracts 分层原则：

```text
Wails UI API
    ↓
Actions / Flows
    ↓
Contracts
    ↓
Adapters
    ↓
外部能力实现（Python gRPC / 未来本地实现）
```

同时它与 repository 的关系是并列而不是替代：

```text
SubmitReviewAction / ReviewFlow
   ├── Repository   # 读取/写入 CardSRS、ReviewLog、FSRSParameter
   └── Contracts    # 计算下一次调度结果
```

### 一句话理解

- `repository` 解决“当前卡片的调度状态与历史记录是什么”
- `contracts.FSRSClient` 解决“给定当前状态与评分，下一次调度结果应该是什么”

---

## 3. 职责边界

## 3.1 FSRS contract 应该定义什么

FSRS contract 应该定义的是**面向业务调度能力的稳定边界**，而不是某个 Python 服务或 gRPC proto 的直接映射。

它关心的问题包括：

- 当前复习评分对应的下一次调度结果
- 哪些输入是调度必须知道的状态
- 哪些输出应回写到 `CardSRS`
- 如何把 transport 错误折叠成上层可判断的 contracts 错误语义

FSRS contract 不应关心：

- gRPC client 如何调用
- proto 如何编码
- 数据库事务如何提交
- review log 如何持久化
- UI 如何展示评分按钮

## 3.2 adapters 应该负责什么

adapter 负责把 typed FSRS contract 映射到当前实际的 Python gRPC transport。

adapter 负责：

- 将 `FSRSRequest` 编码为当前 proto 所需的 `payload_json`
- 调用 `CalculateFsrs`
- 将 `result_json` 解析为 `FSRSResult`
- 把 gRPC、worker、JSON 解析错误转换成 contracts 层错误语义

adapter 不应该：

- 自己决定用哪个 `FSRSParameter`
- 自己决定是否提交 review log
- 自己修改数据库状态

## 3.3 actions / flows 应该负责什么

actions / flows 负责 review 业务编排，而不是持有 transport 细节。

例如：

- 读取 `CardSRS`
- 读取/选择 `FSRSParameter`
- 传入当前评分和 review 时间
- 调用 `contracts.FSRSClient.Calculate`
- 将结果持久化为新的 `CardSRS`
- 追加 `ReviewLog`

actions / flows 不应该：

- 直接构造 proto message
- 直接调用 `pyclient.API().CalculateFsrs(...)`
- 依赖 `payload_json` / `result_json` 的内部 schema

## 3.4 repository 应该负责什么

repository 负责持久化当前调度状态、参数预设与 review 历史，不负责调度计算本身。

当前已有的持久化对象包括：

- `CardSRS`
- `FSRSParameter`
- `ReviewLog`

repository 负责回答的问题是：

- 这张卡当前是什么 `FSRSState`
- 当前的稳定度/难度/次数是多少
- 默认参数预设是什么
- 上一次 review 的日志是什么

但 repository 不负责回答：

- 下一次 due 时间是多少
- 当前评分对应的下一状态是什么

这些属于 `FSRSClient` 的职责范围。

---

## 4. 设计原则

## 4.1 按能力定义 contract，而不是按实现来源定义

推荐：

```go
type FSRSClient interface {
    Calculate(ctx context.Context, req FSRSRequest) (*FSRSResult, error)
}
```

不推荐：

```go
type PythonFSRSGrpcClient interface {
    CalculateFsrs(ctx context.Context, req *pb.CalculateFsrsRequest) (*pb.CalculateFsrsResponse, error)
}
```

原因：后者让上层直接依赖 Python/gRPC/proto，实现边界过深。

## 4.2 contracts 层使用 typed request/result，proto JSON 只是 adapter 内部细节

当前 proto 的形状仍然是：

```text
CalculateFsrs(item_id, payload_json) -> ok, message, result_json
```

这只是当前 transport 的临时现实，不应上升为业务 contract。

contracts 层应固定的是：

- `FSRSRequest`
- `FSRSResult`

而不是：

- `payload_json`
- `result_json`

这样未来无论：

- proto 改成 typed fields
- Python 改成 Go 本地实现
- transport 不再使用 gRPC

上层 actions / flows 的 contract 都可以保持稳定。

## 4.3 调度计算与持久化必须解耦

FSRSClient 的职责是**计算**，不是**持久化**。

因此：

- `FSRSClient.Calculate` 返回下一次调度结果
- `repository.SRS.UpdateAfterReview` 负责持久化状态与日志
- undo 属于历史恢复，不属于调度重算

## 4.4 错误语义必须是业务可判断的，而不是 transport 原生细节

actions 不应该依赖：

- `status.Code(err)`
- proto 里的 `message`
- Python worker 返回的底层实现字符串

adapter 应负责把这些信息映射为 contracts 层统一错误语义。

---

## 5. 推荐接口

`docs/contracts-design.md` 中已有一个合适的起点，但结合 py-fsrs 的实际约束，专项文档需要把“调度计算”和“调度器配置”同时固定下来：

```go
type FSRSClient interface {
    SetScheduler(ctx context.Context, req FSRSSettingsRequest) error
    Calculate(ctx context.Context, req FSRSRequest) (*FSRSResult, error)
}
```

其中：

- `SetScheduler(...)` 用于把当前选中的 FSRS 参数下发给 Python 侧 scheduler
- `Calculate(...)` 用于在已配置 scheduler 的前提下执行一次调度计算
- adapter 负责保证 Python 侧在计算前已经完成 scheduler 初始化，而不是让 action 直接感知 py-fsrs 的初始化细节

### 为什么需要增加设置接口

py-fsrs 的 scheduler 在执行调度前并不是完全无状态的；它需要先基于配置完成初始化，典型参数包括：

- `parameters`
- `desired_retention`
- `maximum_interval`

因此 v1 文档不能再把这些参数仅仅视为 `Calculate(...)` 的附带输入，而应明确分成两类职责：

1. scheduler 配置：确定当前这次会话/实例使用什么 FSRS 参数
2. 调度计算：在给定 card 当前状态和 rating 的情况下，返回下一次调度结果

这样做的目的不是把上层业务绑死到 Python 实现，而是把“需要先初始化 scheduler”这个实现约束收敛在 contract + adapter 边界中，避免未来 action / flow 误以为 FSRS 是完全无状态函数调用。

### 推荐设置请求模型

```go
type FSRSSettingsRequest struct {
    Parameters        []float64
    DesiredRetention  float64
    MaximumInterval   int
}
```

设计要求：

- settings request 表达的是 scheduler 初始化参数，而不是某张卡的调度状态
- `Parameters`、`DesiredRetention`、`MaximumInterval` 的命名应与 Python 侧保持一致
- action / flow 负责从 `FSRSParameter` 读取并组装 `FSRSSettingsRequest`
- adapter 负责把该请求映射到 Python worker 的 `Scheduler.setting` 或等价初始化动作

### 当前阶段结论

当前阶段不建议一次性设计过多变体接口，但应把 scheduler 设置能力显式纳入 contract。

更合适的做法是：

- v1 保留 `SetScheduler(...) + Calculate(...)` 两个基础能力
- 由 action / flow 决定何时切换 preset，并在计算前调用 `SetScheduler(...)`
- adapter 可以在内部做缓存或幂等优化，但这种优化不应改变上层 contract
- preview / submit / undo 的差异仍然放在 action / flow 层，而不是继续在 contract 层拆更多业务接口
- 未来如果确实出现批量预览、多选项计算、离线重算等稳定需求，再扩展更细 contract

---

## 6. 词汇表与合法值

## 6.1 FSRSState

当前 v1 统一使用以下状态值：

- `new`
- `learning`
- `review`
- `relearning`

这与当前 `CardSRS.FSRSState` 注释保持一致。

文档层面的要求：

- 上层传入和持久化时都应只使用这组稳定状态词汇
- adapter / worker 不应私自扩展额外状态并泄漏给 actions

## 6.2 Rating

当前 v1 仍可使用 `int` 表达评分，但必须在语义上固定：

- `1 = Again`
- `2 = Hard`
- `3 = Good`
- `4 = Easy`

原因：

- 当前 `ReviewLog.Rating` 已是整数
- 先保留与现有存储模型一致的最小方案
- 避免当前阶段引入额外 enum 类型后又无法与现有模型直接对齐

但文档必须明确：

- 只允许上述四个值
- 非法值应视为 `ErrInvalidInput`

## 6.3 ReviewKind

`ReviewLog.ReviewKind` 用于记录一次 review 的业务语义，例如：

- `learning`
- `review`
- `relearning`

当前阶段建议把它理解为：

- 这次 review 所属的调度语境或状态类别
- 主要用于历史分析、统计和调试

它不应反过来成为 `FSRSClient` 的 transport 协议字段来源。

---

## 7. 请求与结果模型

## 7.1 推荐请求模型

继续沿用总文档里的建议形状：

```go
type FSRSRequest struct {
    CardID            string
    Rating            int
    ReviewedAtUnix    int64
    FSRSState         string
    Stability         *float64
    Difficulty        *float64
    Reps              int
    Lapses            int
    DesiredRetention  *float64
    MaximumInterval   *int
    ParametersJSON    *string
}
```

## 7.2 推荐结果模型

```go
type FSRSResult struct {
    FSRSState      string
    DueAtUnix      int64
    Stability      *float64
    Difficulty     *float64
    ElapsedDays    *float64
    ScheduledDays  *float64
}
```

### 设计要点

- v1 只暴露调度计算真正必需的字段
- 暂时不把 proto `ok/message/result_json` 暴露给上层
- 若未来需要支持“多评分预览一次返回 4 种结果”，再扩展新的结果结构，而不是污染当前单次 Calculate 的基础 contract

---

## 8. 字段来源与去向

## 8.1 来自 `CardSRS` 的字段

`FSRSRequest` 中以下字段通常来自当前卡片的持久化调度状态：

- `FSRSState`
- `Stability`
- `Difficulty`
- `Reps`
- `Lapses`

这些字段对应：

- `internal/storage/models/srs.go:8`
- `internal/storage/models/srs.go:13`
- `internal/storage/models/srs.go:14`
- `internal/storage/models/srs.go:17`
- `internal/storage/models/srs.go:18`

## 8.2 来自 `FSRSParameter` 的字段

`FSRSRequest` 中以下字段通常来自当前选中的参数预设：

- `DesiredRetention`
- `MaximumInterval`
- `ParametersJSON`

这些字段对应：

- `internal/storage/models/srs.go:37`
- `internal/storage/models/srs.go:38`
- `internal/storage/models/srs.go:39`

## 8.3 来自 action 输入 / clock 的字段

这些字段不应由 repository 或 adapter 隐式生成，而应由 action / flow 明确提供：

- `CardID`
- `Rating`
- `ReviewedAtUnix`

这样做的原因：

- 让业务调用链更清楚
- 避免 adapter 偷偷决定“当前时间”
- 便于测试 review 行为

## 8.4 `FSRSResult` 的去向

`FSRSResult` 中以下字段应被用于回写新的 `CardSRS`：

- `FSRSState`
- `DueAtUnix`
- `Stability`
- `Difficulty`
- `ElapsedDays`
- `ScheduledDays`

但以下内容仍然不属于 `FSRSClient` 自己负责：

- `LastReviewAt`
- `UpdatedAt`
- `ReviewLog` 的创建
- 事务提交

这些属于 action + repository 协作完成的部分。

---

## 9. 参数预设选择策略

## 9.1 选择策略属于 action / flow，而不是 FSRSClient

`FSRSClient` 不应内部偷偷决定“该使用哪个 preset”。

原因：

- preset 选择属于业务策略
- 后续可能存在 knowledge 级 preset、用户自定义 preset、实验性参数切换
- 如果把选择逻辑藏进 adapter，actions 将无法看清业务行为

## 9.2 当前阶段的保守规则

当前仓库里 `FSRSParameterRepository.GetDefault()` 的行为是：

- 返回最早创建的一条记录

这可以作为当前阶段的 fallback，但不能被当作最终产品语义。

因此文档建议：

- 若 action 未显式指定 preset，可暂时使用 repository default
- 但应把它视为“当前阶段可接受的缺省行为”，而不是长期固定规则

## 9.3 未来演进方向

后续可以演进为：

- 显式 default 标记
- knowledge 级 preset 绑定
- 更明确的 preset version / migration 机制

这些都不应影响当前 `FSRSClient` 的基础 contract 形状。

---

## 10. 与持久化模型的协作

## 10.1 `CardSRS`

`CardSRS` 是一张卡当前生效的调度状态快照。

FSRSClient 的输入与输出都围绕它展开：

- 输入：当前 state / difficulty / stability / reps / lapses
- 输出：下一次 state / due / elapsed / scheduled / new difficulty / new stability

## 10.2 `ReviewLog`

`ReviewLog` 是 review 的追加式历史记录，不是 FSRSClient 的直接输出对象。

建议由 action / repository 根据：

- 当前 rating
- reviewed_at
- 旧的 stability / difficulty
- 新的 stability / difficulty
- scheduled / elapsed days

组装成 log，并调用 repository 落库。

这意味着：

- `FSRSClient` 不直接创建 `ReviewLog`
- `FSRSClient` 也不应直接决定 `SnapshotJSON` 的最终 schema

## 10.3 `FSRSParameter`

`FSRSParameter` 的职责是提供：

- 算法权重（`ParametersJSON`）
- 目标保持率（`DesiredRetention`）
- 最大间隔（`MaximumInterval`）

contract 层不关心这些参数是从数据库、配置文件、内存缓存还是未来远端配置加载来的；它只关心 action 最终传入了什么参数值。

---

## 11. 典型协作链路

## 11.1 SubmitReview

推荐调用链：

```text
SubmitReviewAction
  -> repository.Card.GetByID
  -> repository.SRS.GetByCardID
  -> repository.FSRSParameter.GetDefault / GetByID
  -> contracts.FSRSClient.Calculate
  -> repository.SRS.UpdateAfterReview
```

这里的关键点：

- action 负责读取当前状态与参数
- FSRSClient 只负责计算
- repository 负责原子更新 `CardSRS` 与 `ReviewLog`

## 11.2 PreviewSchedule

推荐调用链：

```text
PreviewScheduleAction
  -> repository.SRS.GetByCardID
  -> repository.FSRSParameter.GetDefault / GetByID
  -> contracts.FSRSClient.Calculate
  -> return preview result only
```

关键点：

- preview 可以调用 `Calculate`
- 但 preview 不应持久化结果
- preview 的存在不要求先把 contract 拆成新的接口类型

## 11.3 UndoLastReview

推荐调用链：

```text
UndoLastReviewAction
  -> repository.SRS.GetLastReviewLog / UndoLastReview
  -> restore state from history
```

关键点：

- undo 不应该重新调用 `FSRSClient.Calculate`
- undo 是基于历史记录与快照做恢复
- 当前 `internal/storage/repository/srs_repo.go` 的实现方向已经体现了这一边界

---

## 12. adapter 与 proto 的关系

## 12.1 当前 proto 的定位

当前 proto 里的 FSRS RPC 还是占位形态：

```proto
message CalculateFsrsRequest {
  string item_id = 1;
  bytes payload_json = 2;
}

message CalculateFsrsResponse {
  bool ok = 1;
  string message = 2;
  bytes result_json = 3;
}
```

这说明当前 transport 仍是“JSON 套在 gRPC 里”的过渡方案。

## 12.2 adapter 的职责

adapter 应负责：

- 将 `FSRSRequest.CardID` 映射到 `item_id`
- 将剩余 typed 字段编码成 `payload_json`
- 解析 `result_json` 得到 `FSRSResult`
- 根据 `ok/message` 与底层错误判断如何映射 contracts 错误

## 12.3 为什么不能让 action 直接使用 raw gRPC client

`internal/pyclient/client.go` 当前暴露 `API()`，并明确注明：

- raw generated client 只是过渡层
- business calls 应该留在 services / adapter 中

因此文档应明确：

- action 不应该直接 `API().CalculateFsrs(...)`
- 后续应在 raw transport 之上增加更高层的 FSRS adapter/facade

### 推荐目录关系

```text
internal/contracts/
└── fsrs.go

internal/adapters/
└── pyclient/
    └── fsrs_client.go
```

未来如果需要本地实现，可扩展为：

```text
internal/adapters/
├── pyclient/fsrs_client.go
└── fsrs/local.go
```

---

## 13. 错误语义

FSRS contract 应沿用 `docs/contracts-design.md` 中统一的错误语义。

## 13.1 `ErrInvalidInput`

适用场景：

- `Rating` 不在 `1..4`
- `FSRSState` 不在允许集合中
- 必填字段缺失
- `ParametersJSON` 结构无法被 adapter/worker 接受

## 13.2 `ErrUnavailable`

适用场景：

- Python worker 不可达
- gRPC 超时
- worker 返回不可用结果
- `result_json` 无法解析且可判断为下游服务异常

## 13.3 `ErrNotFound`

在 FSRS 纯计算阶段通常不是主错误。

更常见的做法应是：

- 卡不存在 -> action / repository 先处理
- preset 不存在 -> action / repository 先处理

也就是说，`ErrNotFound` 更常出现在调用 `FSRSClient` 之前，而不是计算过程本身。

## 13.4 `ErrConflict`

在纯调度计算中一般不是主错误。

更常见的场景是：

- action 提交 review 时遇到并发状态冲突
- repository 在持久化层发现乐观冲突或历史冲突

因此文档应明确：

- `ErrConflict` 不是 FSRSClient 的主路径错误
- 它更多是 review action / repository 层的语义

---

## 14. 当前阶段的设计结论

当前阶段推荐采用以下保守方案：

1. contracts 层固定 typed `FSRSRequest` / `FSRSResult`
2. proto 的 JSON payload/result 只留在 adapter 内部
3. `FSRSClient` v1 只保留单个 `Calculate(...)`
4. preset 选择由 action / flow 完成
5. submit review 会调用 FSRS，preview 也可调用 FSRS，但 undo 不调用 FSRS
6. `CardSRS` 是当前状态快照，`ReviewLog` 是历史追加记录，二者都不由 FSRSClient 自己持久化
7. actions 只依赖 contracts 错误语义，不依赖 gRPC/status/proto 细节

---

## 15. 后续演进方向

这份文档刻意只固定当前最需要稳定下来的 contract 边界，同时为后续演进留空间。

后续可扩展方向包括：

- 将 proto 从 `payload_json/result_json` 演进为 typed fields
- 增加多评分预览接口
- 明确 `ReviewLog.SnapshotJSON` 的结构化 schema
- 引入更明确的 parameter default/version 机制
- 增加本地 Go FSRS adapter，与 Python gRPC 实现并存

这些演进不应破坏当前 `FSRSClient` 作为上层稳定 capability contract 的定位。

---

## 16. 验证清单

本次是专项设计文档任务，验证以“边界是否清楚、字段是否与当前代码一致、协作链路是否闭环”为主。

## 16.1 边界检查

需要检查文档是否清楚说明了：

- FSRSClient 负责计算，不负责持久化
- adapter 负责 transport 映射
- repository 负责 `CardSRS` / `ReviewLog` / `FSRSParameter`
- action 负责 preset 选择与 review 编排

## 16.2 字段一致性检查

需要检查文档是否与当前代码术语保持一致：

- `FSRSClient`
- `FSRSRequest`
- `FSRSResult`
- `CardSRS`
- `FSRSParameter`
- `ReviewLog`
- `FSRSState`
- `Rating`

## 16.3 协作链路检查

需要检查文档是否明确区分：

- SubmitReview
- PreviewSchedule
- UndoLastReview

并且是否清楚说明：

- 哪些调用 FSRS
- 哪些不调用 FSRS
- 哪些步骤由 repository 负责

## 16.4 transport 隔离检查

需要检查文档是否清楚表达：

- `payload_json` / `result_json` 只是 adapter 内部细节
- action 不应直接依赖 proto/gRPC client
- 当前 Python 实现只是第一阶段实现来源，而不是 contract 的一部分

---

## 17. 结论

当前 kmemo 最需要的，不是马上把 FSRS 算法完整落地到代码，而是先把 contracts 层的边界固定下来：

- 用 typed request/result 表达稳定调度能力
- 用 action + repository 组织 review 持久化流程
- 用 adapter 隔离当前 JSON-over-gRPC 的过渡 transport
- 用统一的状态、评分和错误语义减少隐含约定

这样后续无论是完善 Python worker、升级 proto，还是增加本地 Go 实现，业务层都能继续依赖一套稳定、可理解、可扩展的 FSRS contract。
