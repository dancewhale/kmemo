# FSRS 端口与适配器设计（contracts + internal/adapters/fsrs）

## 1. 目标

- 将 **FSRS 相关能力** 从「散落 `pyclient` + 上层自行拼 proto」收拢为：**明确的端口（port）+ 可替换的适配器（adapter）**。
- **上层（actions / services / app）** 只依赖 **Go 领域侧类型**（以 `internal/storage/models` 中与 FSRS 相关的结构为主），**不直接持有** `*kmemov1.*` 或 gRPC 细节。
- 迁移路径清晰：`internal/pyclient` 中与 FSRS 强相关的代码 **迁入** `internal/adapters/fsrs`（及必要的共享连接层）；`pyclient` **可整体废弃** 或 **降级为仅连接/多领域网关**（见第 8 节待你确认）。

本文档约定 **可读性**（读一个接口便知能做什么）与 **可修改性**（替换 Python 实现、替换传输层、替换持久化模型时的改动面）之间的平衡。

---

## 2. 与开源实践的对齐（摘要）

常见做法（Hexagonal / Ports & Adapters、Clean Architecture）：

| 层次 | 职责 | 依赖方向 |
|------|------|----------|
| **Port（接口）** | 描述应用需要什么能力，入参/出参用**应用语言**（此处即 DB 模型或薄 DTO） | 不依赖 gRPC/proto |
| **Adapter（实现）** | 调用外部系统；**唯一允许**依赖 `gen/kmemo/v1`、grpc 客户端 | 依赖 port + 外部类型 |
| **Mapping** | proto ↔ 模型 | 放在 adapter 侧或独立的 `wire/conv` 包，避免业务层复制粘贴 |

你提出的要求——**「gRPC 返回类型转换为 GORM 模型，方便上层使用」**——在业界通常落实为：

- **端口方法签名**保证上层拿到的是 `*models.CardSRS`、`*models.ReviewLog`、`*models.FSRSParameter` 等；
- **转换代码**多数项目放在 **adapter 包内**（或 `internal/wire/fsrs`），而不是让「接口类型」去实现转换（Go 的 interface 本身不包含实现）。


---

## 3. 推荐包结构

```text
internal/contracts/fsrs/
├── doc.go              # 包说明：本包为 FSRS 端口，禁止在上层使用 proto
├── client.go           # FSRSClient 接口定义
└── types.go            # 端口专用输入输出 DTO

internal/adapters/grpcworker/
└── conn.go             # 共享 gRPC 连接，仅负责 Dial / 生命周期管理

internal/adapters/fsrs/
├── grpc_python/
│   ├── scheduler_client.go   # 对接 FsrsSchedulerService
│   ├── optimizer_client.go   # 对接 FsrsOptimizerService（若 Go 侧启用）
│   └── conv.go               # proto ↔ contracts/fsrs DTO 映射
├── noop/
└── inmemory/
```

说明：

- `contracts/fsrs` 只放接口 + 纯 Go DTO，不在接口方法里出现 `*kmemov1.*`。
- `grpcworker` 只负责共享连接；FSRS 主路径通过专门的 `FsrsSchedulerService` / `FsrsOptimizerService` 适配器访问。
- `adapters/fsrs/grpc_python/conv.go` 负责 proto ↔ Go DTO 映射，避免上层复制 transport 细节。

---

## 4. 端口（contracts）设计

### 4.1 核心接口（建议命名）

```go
// FSRSClient 描述当前 Go 主路径使用的两类基础能力：
// 1) 先将当前 preset/setting 下发给远端 scheduler；
// 2) 再执行单卡 review 计算。
type FSRSClient interface {
    SetScheduler(ctx context.Context, in SetSchedulerInput) error
    Calculate(ctx context.Context, in ReviewInput) (*ReviewOutput, error)
}
```

说明：

- 当前 Go contracts 主路径与 `docs/python/fsrs-module-design.md` 对齐，优先保留 `SetScheduler(...)` + `Calculate(...)` 两类基础能力。
- Python transport 侧可以更细地实现 `SettingScheduler`、`GetCardRetrievability`、`ReviewCard`、`RescheduleCard`、`OptimizeParameters`，但不要求 Go contract 第一版一次性全部暴露。
- 若后续 Go 明确需要 retrievability / reschedule / optimize，再在不破坏主路径的前提下增补能力接口。

### 4.2 输入输出类型（建议）

建议使用显式 Input/Output 结构体，与 Python typed RPC 的核心语义保持同构：

- `SetSchedulerInput`：`Parameters []float64`、`DesiredRetention *float64`、`MaximumInterval *int`
- `ReviewInput`：`CardID`、`State`、`Stability`、`Difficulty`、`Due`、`LastReview`、`ElapsedDays`、`ScheduledDays`、`Reps`、`Lapses`、`Rating`、`ReviewedAt`
- `ReviewOutput`：
  - `Card`：`State`、`Stability`、`Difficulty`、`Due`
  - `ReviewLog`：`Rating`、`Review`、`ElapsedDays`、`ScheduledDays`
  - `Retrievability *float64`
  - `Warnings []string`

这样 actions 层代码可以清晰表达为：

```go
if err := client.SetScheduler(ctx, setting); err != nil { ... }
out, err := client.Calculate(ctx, reviewReq)
```

而不用在上层暴露 `SettingSchedulerRequest` / `ReviewCardResponse` 等 proto 细节。

### 4.3 与现有 models 的关系

- 优先复用 **`models.CardSRS`、`models.ReviewLog`、`models.FSRSParameter`**，与 repository 层一致，减少二次映射。
- 若未来模型与 proto 字段严重分叉，再在 `contracts/fsrs` 引入 **端口专用 DTO**，由 repository 在边界做一次转换（现阶段不必）。

---

## 5. 适配器（adapters/fsrs）职责

### 5.1 `grpc_python` 实现

1. 通过共享 `grpcworker` 连接调用 `FsrsSchedulerService` 的 `SettingScheduler`、`ReviewCard`，以及 `FsrsOptimizerService` 的 `OptimizeParameters`（若 Go 侧启用）。
2. **在 adapter 内**完成：
   - 请求：contracts/fsrs DTO → `*kmemov1.*Request`
   - 响应：`*kmemov1.*Response` → contracts/fsrs DTO
3. 错误处理：
   - gRPC `status` → 包装为带 `Unwrap` 的 `error`，必要时定义可识别的哨兵错误（如 `ErrUnavailable`、`ErrInvalidArgument`）。

补充：

- 旧 `CalculateFsrs(item_id, payload_json)` 不再作为适配目标。
- `KmemoProcessor` 不再承载 FSRS 主路径；若仓库中仍保留该 service，也仅处理与 FSRS 无关的其他能力。

### 5.2 转换代码放置（与「contracts 负责转换」的对应关系）

| 方案 | 转换代码位置 | contracts 是否 import proto | 说明 |
|------|----------------|----------------------------|------|
| **A（推荐）** | `internal/adapters/fsrs/grpc_python/conv.go` | 否 | 端口只声明 Go DTO；由 FSRS adapter 负责把 typed gRPC 映射成上层可用结构，与最新 Python FSRS 设计保持一致 |


---

## 6. `internal/pyclient` 的迁移与废弃策略

### 6.1 当前职责拆分

| 现状文件 | 内容 | 迁移去向 |
|----------|------|----------|
| `client.go` | 历史上的 raw generated client / 连接封装 | **连接**：下沉到 `internal/adapters/grpcworker` |
| `fsrs.go` | FSRS proto ↔ Go DTO 映射 | **`internal/adapters/fsrs/grpc_python/conv.go`** |
| `fsrs_test.go` | 映射单测 | 随 conv 迁移，测试表驱动覆盖边界 |

### 6.2 废弃策略

- FSRS 相关代码不再继续挂在 `internal/pyclient` 名下。
- 多领域共享连接时，仅保留一个极薄的 `grpcworker` 连接层；不要再保留语义含糊的 `pyclient` 大包。
- 若现有实现已经完成迁移，相关设计文档应以 `grpcworker + adapters/fsrs` 为准，而不是继续把 `pyclient` 当作长期结构。
---

## 7. 上层调用约定（actions / repository）

- **actions**：只依赖 `FSRSClient.SetScheduler(...)` 与 `FSRSClient.Calculate(...)`；复习写库仍用 `repository.SRSRepository.UpdateAfterReview`（或等价事务接口）。
- **典型调用顺序**：先从 `FSRSParameter` 组装 scheduler setting，下发 `SetScheduler(...)`，再调用 `Calculate(...)` 执行单卡 review 计算。
- **禁止**：actions 直接依赖 `*kmemov1.*` request/response 或调用旧 `CalculateFsrs(payload_json)`。
- **ReviewLog 主键**：若 Python 响应仅返回 review 结果字段，可由 Go 在持久化前补齐主键；adapter 不负责决定持久化主键策略。

---

## 8. 与其他 Python 能力的关系

- 若 Python worker 同时承载 HTML、source-process 等能力，建议通过 **`internal/adapters/grpcworker`** 共享连接。
- 连接复用不等于 service 合并：FSRS 仍应使用专门的 `FsrsSchedulerService` / `FsrsOptimizerService`。
- 因此，文档中的长期结构应是“共享连接层 + 按领域拆 adapter/service”，而不是“单一 `KmemoProcessor` + 单一 `pyclient`”。

---

## 9. 测试与可修改性

- **conv 单测**：不启动 Python，仅 proto ↔ model 表驱动测试（迁移 `fsrs_test.go`）。
- **集成测试**：可选 `task run:python` + 真实 gRPC（标 `integration` build tag）。
- **替换实现**：未来增加 `adapters/fsrs/native_go`，实现同一 `FSRSClient`，actions **零改**。

---

## 10. 小结

- **端口**在 `internal/contracts/fsrs`：当前 Go 主路径优先保留 `SetScheduler(...)` + `Calculate(...)` 两类基础能力。
- **Python service**：FSRS 主路径由专用 `FsrsSchedulerService` / `FsrsOptimizerService` 承担；旧 `CalculateFsrs` 与 `KmemoProcessor` 中的 FSRS 入口均应退出主设计。
- **映射职责**：proto ↔ Go DTO 的转换仍由 `internal/adapters/fsrs/grpc_python/conv.go` 承担。
- **连接层**：多领域共享连接时保留极薄的 `grpcworker`，不要继续以 `pyclient` 作为长期结构名称。

