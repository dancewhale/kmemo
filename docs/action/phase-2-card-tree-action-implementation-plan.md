# Action Layer Phase 2 实现规划（Card Tree / he-tree）

## 1. 目标与范围

本阶段目标是为前端引入 he-tree（SuperMemo 风格卡片树交互）提供后端 action 层支撑，重点实现：

- `CreateCard` 在创建时自动计算并写入 `sort_order`
- 新增 `MoveCard(cardId, targetParentId, targetIndex)` 能力
- 按 `docs/ui/card-tree-he-tree-design.md` 补齐 Phase 3 所需核心函数（优先 `ReorderCardChildren`，可选 `GetCardPath`）

本规划聚焦 action/repository/app 三层 API 与事务语义，不改动前端实现细节。

## 2. 当前问题与改造目标

### 2.1 当前问题

- `CreateCard` 当前未计算同级顺序，`sort_order` 依赖默认值，导致同级排序不稳定。
- action 层不存在卡片移动能力，无法支持 he-tree 拖拽落点。
- 仓储层缺少“同级重排”批量写入接口，动作层无法保证顺序一致性。

### 2.2 改造目标

- 保证新建卡片顺序可预测：始终 append 到目标父节点（或根节点）末尾。
- 保证移动操作顺序一致：源父与目标父两侧均完成连续重排。
- 保证接口可扩展：后续可对接 `ReorderCardChildren` 与 `GetCardPath`。

## 3. 分层改造清单

## 3.1 repository 层（`internal/storage/repository`）

为 `CardRepository` 增加以下能力：

1) 查询同级最大顺序

- `GetMaxSortOrder(ctx, knowledgeID string, parentID *string) (int, error)`
- 用于 `CreateCard` 自动分配 `sort_order = max + 1`

2) 同级列表（按 `sort_order ASC`）

- `ListSiblings(ctx, knowledgeID string, parentID *string) ([]*models.Card, error)`
- 用于 move/reorder 的顺序重算基础

3) 批量更新顺序

- `BatchUpdateSortOrders(ctx, updates []CardSortOrderUpdate) error`
- `CardSortOrderUpdate{CardID string, SortOrder int, ParentID *string, IsRoot bool}`

4) 更新父子关系（单卡）

- `UpdateParent(ctx, cardID string, targetParentID *string, targetKnowledgeID string, isRoot bool) error`
- 说明：当前 MoveCard 先限定 knowledge 不变；若未来允许跨 knowledge，可复用此签名。

实现建议：

- 所有 move/reorder 写操作只在事务内调用。
- 增加 `sort_order` 与 `(knowledge_id,parent_id,is_root,deleted_at)` 相关索引评估（如后续性能瓶颈出现）。

## 3.2 action 层（`internal/actions/card/service.go`）

新增输入结构与方法：

```go
type MoveInput struct {
  CardID         string
  TargetParentID *string
  TargetIndex    int
}

type ReorderChildrenInput struct {
  ParentID        *string
  KnowledgeID     string
  OrderedChildIDs []string
}
```

### A. 改造 `Create`

新增逻辑：

1. 校验 `KnowledgeID` 存在（保留现有逻辑）。
2. 若 `ParentID != nil`，校验 parent 卡片存在且 parent.knowledge_id 与 `KnowledgeID` 一致。
3. 调 `GetMaxSortOrder(knowledgeID,parentID)` 计算 `nextOrder`。
4. 创建 `models.Card` 时写入：
   - `SortOrder: nextOrder`
   - `IsRoot: parentID == nil`
5. 保持当前文件写入 + DB 事务 + 失败补偿策略。

约束：

- `nextOrder` 以事务内读取 + 写入为准，避免并发创建冲突。
- 并发冲突场景建议在 DB 层配合唯一索引或重试策略（可选，见风险章节）。

### B. 新增 `Move`

签名建议：

```go
func (s *Service) Move(ctx context.Context, input MoveInput) error
```

核心流程（事务内）：

1. 加载 `card`；校验存在。
2. 校验 `targetParentId`：
   - 不可等于 `cardId`
   - 不可为其后代（防环）
   - 若非空，targetParent 必须存在且 knowledge 一致
3. 读取源同级列表、目标同级列表（若同父则复用一份）。
4. 从源列表移除当前卡，插入目标列表 `targetIndex`（越界则 clamp 到 `[0,len]`）。
5. 生成顺序更新集：
   - 目标列表从 0..n-1 重写 `sort_order`
   - 若跨父，还需重写源列表 `sort_order`
6. 更新当前卡 `parent_id`、`is_root`（必要时 knowledge_id 保持原值）。
7. 批量写入顺序。

行为约定：

- `targetIndex` 按“插入后索引”解释。
- 同父移动时也要稳定重排，避免空洞顺序。

### C. 新增 `ReorderChildren`（按设计文档要求）

签名建议：

```go
func (s *Service) ReorderChildren(ctx context.Context, input ReorderChildrenInput) error
```

核心逻辑：

1. 查询该父节点下当前 children。
2. 校验 `orderedChildIds` 与当前 children 集合完全一致（防漏项/脏数据）。
3. 按入参顺序写回 `sort_order`。

用途：

- he-tree 同级拖拽后可直接调用，减少 Move 的复杂度。

### D. 可选新增 `GetPath`

签名建议：

```go
func (s *Service) GetPath(ctx context.Context, cardID string) ([]*models.Card, error)
```

用于支持前端“定位并展开祖先链”，非本阶段强制实现，但建议预留。

## 3.3 app 层（`internal/app/card.go`）

新增请求结构与暴露接口：

```go
type MoveCardRequest struct {
  CardID         string  `json:"cardId"`
  TargetParentID *string `json:"targetParentId"`
  TargetIndex    int     `json:"targetIndex"`
}

type ReorderCardChildrenRequest struct {
  ParentID        *string  `json:"parentId"`
  KnowledgeID     string   `json:"knowledgeId"`
  OrderedChildIDs []string `json:"orderedChildIds"`
}
```

新增方法：

- `MoveCard(req MoveCardRequest) error`
- `ReorderCardChildren(req ReorderCardChildrenRequest) error`

校验建议：

- `CardID` 非空
- `TargetIndex >= 0`
- `OrderedChildIDs` 非空且无重复

## 4. 事务与一致性策略

## 4.1 CreateCard 一致性

- 文件创建在 DB 事务外，维持当前“事务失败 -> 删除文件补偿”机制。
- `sort_order` 计算与写入必须在同一 DB 事务内完成。

## 4.2 Move/Reorder 一致性

- 必须全程单事务：
  - 更新 parent/isRoot
  - 批量写 sort_order
- 任一步失败即回滚，保证前端刷新后顺序正确。

## 5. 详细执行顺序（建议）

1. repository 接口扩展与实现
2. action `Create` 自动排序改造
3. action `Move` 实现 + 单测
4. action `ReorderChildren` 实现 + 单测
5. app 层暴露 `MoveCard` / `ReorderCardChildren`
6. 回归测试 `ListCards` / `GetCardChildren` 排序行为

## 6. 测试计划

单测重点：

1) `CreateCard` 排序

- 在同一父节点下连续创建 3 张卡，断言 `sort_order = 0,1,2`
- 在 root 下创建与在 parent 下创建分别验证

2) `MoveCard` 同父移动

- 初始 `[A,B,C,D]`，移动 `D -> index 1`，断言 `[A,D,B,C]` 与顺序连续

3) `MoveCard` 跨父移动

- `P1:[A,B], P2:[C]`，移动 `B -> P2 index 1`，断言：
  - `P1:[A]`
  - `P2:[C,B]`
  - 两侧 `sort_order` 连续

4) 防环校验

- 目标父为自身或后代时返回 `ErrInvalidInput`

5) `ReorderChildren`

- 入参完整一致时成功重排
- 入参缺失/重复/越权时失败

## 7. 风险与缓解

- 并发创建导致顺序冲突：
  - 缓解：事务内取 max+1；必要时增加重试机制。
- 大批量 sibling 重排写放大：
  - 缓解：先实现全量重排保证正确性，后续按“区间重排”优化。
- 旧数据存在 `sort_order` 空洞：
  - 缓解：Move/Reorder 每次写回连续序号，逐步自愈。

## 8. 与 `docs/ui/card-tree-he-tree-design.md` 对齐关系

- 本文覆盖其“后端必须补充”的核心项：
  - `MoveCard`
  - `ReorderCardChildren`
- 本文补充并细化了 `CreateCard` 自动排序规则（前端 he-tree 稳定体验的必要基础）。
- `GetCardPath` 定为可选增强，保持与设计文档一致。

## 9. 交付物清单

- `internal/storage/repository/card_repo.go`：新增排序/移动相关仓储实现
- `internal/actions/card/service.go`：Create 自动排序 + Move/Reorder 业务编排
- `internal/app/card.go`：新增 Wails 暴露接口与请求结构
- 对应单元测试文件（action/repository 维度）

