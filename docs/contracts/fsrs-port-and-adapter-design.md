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
├── scheduler.go        # FSRSScheduler（或 FSRSClient）接口定义
└── types.go            # （可选）端口专用输入 DTO，避免 models 过重时再抽

internal/adapters/fsrs/
├── grpc_python/        # 默认实现：通过现有 KmemoProcessor gRPC 调 Python
│   ├── client.go       # 实现 contracts/fsrs.FSRSScheduler；持有 *grpc.ClientConn 或裸 KmemoProcessorClient
│   ├── dial.go         # 从 config 拨号（可由 pyclient.New 逻辑迁入）
│   └── conv.go         # *kmemov1.* ↔ *models.* 双向映射（由 pyclient/fsrs.go 迁入并改名整理）
├── noop/               # （可选）测试用空实现
└── inmemory/           # （可选）未来 Go 原生 FSRS 算法时的本地实现
```

说明：

- **`contracts/fsrs` 只放接口 +（可选）纯 Go DTO**，**不推荐**在接口方法里出现 `*kmemov1.*`。
- **`adapters/fsrs/grpc_python/conv.go`** 承接当前 `internal/pyclient/fsrs.go` 中的映射函数，并按「入站/出站」分块注释，便于单测。

---

## 4. 端口（contracts）设计

### 4.1 核心接口（建议命名）

```go
// FSRSScheduler 描述「调度/复习/可检索度/改期/参数优化」等对外能力。
// 所有方法入参、出参均为应用侧类型，不暴露 protobuf。
type FSRSScheduler interface {
    // Review 根据当前 SRS 快照与评分，计算新的 SRS 状态与一条待持久化的 ReviewLog（不含主键时可由上层生成 ID）。
    Review(ctx context.Context, in ReviewInput) (*ReviewOutput, error)

    GetRetrievability(ctx context.Context, in RetrievabilityInput) (*RetrievabilityOutput, error)

    Reschedule(ctx context.Context, in RescheduleInput) (*RescheduleOutput, error)

    // SetGlobalSetting 将 FSRS 参数推送到远端（若远端无状态可变为 no-op）。
    SetGlobalSetting(ctx context.Context, param *models.FSRSParameter) error

    OptimizeParameters(ctx context.Context, in OptimizeParametersInput) (*models.FSRSParameter, error)
}
```

### 4.2 输入输出类型（建议）

用 **显式 Input/Output 结构体** 代替长参数列表，便于演进（与 gRPC 字段解耦）：

- `ReviewInput`：`CardID`、`Prior *models.CardSRS`（nil 表示 new）、`Rating int`、`ReviewedAt time.Time`
- `ReviewOutput`：`Next *models.CardSRS`、`Log *models.ReviewLog`、`Retrievability float64`（若 Python 返回）、`Warnings []string`（若有）

这样 **actions 层** 代码形如：`out, err := scheduler.Review(ctx, in)`，**无需知道** `ReviewCardResponse` 的存在。

### 4.3 与现有 models 的关系

- 优先复用 **`models.CardSRS`、`models.ReviewLog`、`models.FSRSParameter`**，与 repository 层一致，减少二次映射。
- 若未来模型与 proto 字段严重分叉，再在 `contracts/fsrs` 引入 **端口专用 DTO**，由 repository 在边界做一次转换（现阶段不必）。

---

## 5. 适配器（adapters/fsrs）职责

### 5.1 `grpc_python` 实现

1. 调用 `kmemov1.KmemoProcessorClient` 的 `ReviewCard`、`GetCardRetrievability`、`RescheduleCard`、`SchedulerSetSetting`、`OptimizeParameters` 等。
2. **在 adapter 内**完成：
   - 请求：`models` / Input DTO → `*kmemov1.*Request`（逻辑从现 `pyclient/fsrs.go` 迁入）。
   - 响应：`*kmemov1.*Response` → `ReviewOutput` 等（**此处完成你要求的「gRPC → GORM 模型」**）。
3. 错误处理：
   - gRPC `status` → 包装为带 `Unwrap` 的 `error`，必要时定义 `contracts/fsrs` 或 `actions/errs` 可识别的哨兵错误（如 `ErrUnavailable`、`ErrInvalidArgument`）。

### 5.2 转换代码放置（与「contracts 负责转换」的对应关系）

| 方案 | 转换代码位置 | contracts 是否 import proto | 说明 |
|------|----------------|----------------------------|------|
| **A（推荐）** | `internal/adapters/fsrs/grpc_python/conv.go` | 否 | 端口只声明「返回 models」；**语义上**由 FSRS 适配器模块负责「把 gRPC 变成 models」，与开源习惯一致 |


---

## 6. `internal/pyclient` 的迁移与废弃策略

### 6.1 当前职责拆分

| 现状文件 | 内容 | 迁移去向 |
|----------|------|----------|
| `client.go` | gRPC 连接、`KmemoProcessorClient` 封装、FSRS 相关 RPC 透传 | **连接**：放在 `adapters/fsrs` |
| `fsrs.go` | `models` ↔ proto | **`adapters/fsrs/conv.go`**|
| `fsrs_test.go` | 映射单测 | 随 conv 迁移，测试表驱动覆盖边界（nil、缺字段） |

### 6.2 直接废弃 `pyclient` 包

可直接将 `internal/pyclient` **标记 deprecated** 并删除：
---

## 7. 上层调用约定（actions / repository）

- **actions**：只依赖 `FSRSScheduler.Review` 等；复习写库仍用 `repository.SRSRepository.UpdateAfterReview`（或等价事务接口）。
- **禁止**：actions 直接 `ReviewCard(ctx, *kmemov1.ReviewCardRequest)`。
- **ReviewLog 主键**：可在 **adapter 返回后**由 actions 统一 `uuid.NewV7()`，与当前实现一致；或在 Output 中约定「若 Log.ID 为空则上层填」。

---

## 8. 待你确认（影响目录与是否保留 pyclient）

1. **Python worker 除 FSRS 外**，是否还有 **CleanHtml、PrepareImportMaterial** 等仍要走同一 `KmemoProcessor` 连接？  
   - 是的：建议单独增加 `contracts/html`、`contracts/import` 等端口，并抽 **`internal/adapters/grpcworker`** 只负责 **Dial + 共享 Conn**，`adapters/fsrs` 注入其中的 `KmemoProcessorClient` 或子接口，**避免**为每个领域各拨一次号。  

2. **转换代码**你更希望落在 **A（adapter 内 conv）** 还是 **B（contracts/fsrs/conv 依赖 proto）**？  
   -  `A`

3. **接口命名**：更偏好 `FSRSScheduler`、`FSRSClient`.

---

## 9. 测试与可修改性

- **conv 单测**：不启动 Python，仅 proto ↔ model 表驱动测试（迁移 `fsrs_test.go`）。
- **集成测试**：可选 `task run:python` + 真实 gRPC（标 `integration` build tag）。
- **替换实现**：未来增加 `adapters/fsrs/native_go`，实现同一 `FSRSScheduler`，actions **零改**。

---

## 10. 小结

- **端口**在 `internal/contracts/fsrs`：**方法签名只使用 models / 端口 DTO**，保证上层可读、可测。
- **gRPC → GORM 模型**的转换在工程上由 **`internal/adapters/fsrs`实现**；。
- **`pyclient`**：FSRS 映射与调用迁入 **`adapters/fsrs`** 后，包级 **可废弃**；若多领域共享连接，保留 **极薄的 grpcworker**，而不是保留名为 `pyclient` 的模糊大包。

