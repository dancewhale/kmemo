# zaplog 使用约定

本文档定义 `internal/zaplog` 的统一使用方式，避免业务代码中散落日志初始化、级别判断与 request_id 处理逻辑。

## 目标

- 业务代码只表达“记录什么”，不处理“怎么记录”。
- `request_id` 自动从 `context.Context` 注入。
- Debug 开关统一在底层判断，不在业务层写 `if debug`.
- SQL trace 与业务动作日志职责分离、可同时追踪。

## 统一入口

业务代码统一使用链式 API：

- `zaplog.L(ctx).Debug(...)`
- `zaplog.L(ctx).Info(...)`
- `zaplog.L(ctx).Warn(...)`
- `zaplog.L(ctx).Error(...)`
- `zaplog.L(ctx).Named("card").Debug(...)`

示例：

```go
zaplog.L(ctx).Named("card").Debug("card.create.start",
	zap.String("knowledge_id", input.KnowledgeID),
	zap.String("card_type", input.CardType),
)
```

## 上下文约定

调用链入口（如 Desktop 方法、适配器入口）应确保：

1. `ctx = zaplog.WithLogger(ctx, baseLogger)`
2. `ctx, _ = zaplog.EnsureRequestID(ctx)`

之后链路内统一传递该 `ctx`，日志自动带 `request_id`。

## 级别与开关约定

### 业务动作日志

- 使用 `zaplog.L(ctx).Debug(...)` 记录过程事件（start/success/fail）。
- Debug 是否输出由 `zaplog` 底层自动判定（`DebugEnabled(ctx)`）。
- 业务代码禁止手写 `if zaplog.DebugEnabled(ctx)`.

### GORM SQL Trace

`gorm_logger.Trace` 仅在以下条件同时满足时开启：

- `log-level == debug`
- `repository_debug == true`

判断统一走：

- `zaplog.ShouldEnableRepositoryTrace(level, repositoryDebug)`

## 字段与命名建议

- 事件名：`<domain>.<action>.<stage>`，例如：
  - `card.create.start`
  - `card.create.success`
  - `card.create.fail`
- 常用字段：
  - 目标实体：`card_id`, `knowledge_id`, `tag_id`
  - 阶段：`phase`
  - 耗时：`duration`
  - 错误：`zap.Error(err)`

## 反例（禁止）

- 业务代码中直接：
  - `zaplog.FromContext(ctx).Named(...).Debug(...)`
  - `if zaplog.DebugEnabled(ctx) { ... }`
- 在日志字段中手动重复追加 `request_id`（上下文已自动注入）。

## 兼容与扩展

- 若必须直接拿原始 `*zap.Logger`（少数基础设施场景），可用：
  - `zaplog.Logger(ctx)`
  - `zaplog.LoggerNamed(ctx, name)`
- 新增日志功能优先扩展 `zaplog`，避免在业务侧复制判定逻辑。

