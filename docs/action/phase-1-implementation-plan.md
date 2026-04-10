# Action Layer Phase 1 Implementation Details

## 1. Scope implemented in current code

当前仓库已建立第一阶段 action 主干，重点是先打通 `Desktop -> Actions -> Repository/Contract` 调用链，并落地核心闭环的最小实现。

已落地代码：
- 注入链路
  - `internal/bootstrap/bootstrap.go`
  - `internal/app/desktop.go`
  - `desktop/app.go`
  - `internal/actions/actions.go`
- knowledge
  - `internal/actions/knowledge/service.go`
  - `internal/app/knowledge.go`
- card
  - `internal/actions/card/service.go`
  - `internal/app/card.go`
- review
  - `internal/actions/review/service.go`
  - `internal/app/review.go`
- search + tag/search UI entry
  - `internal/actions/search/service.go`
  - `internal/app/tag_search.go`

## 2. Current dependency chain

当前对象图如下：

```text
bootstrap.NewHeadless
  -> config.Load
  -> storage.New(..., withMigrate=true)
  -> repository.NewRepositoryFactory
  -> repository.NewTransactionManager
  -> file.NewFileStore + adapters/filestore.NewAdapter
  -> fsrs scheduler (grpc_python or noop)
  -> sourceprocess processor (grpc_python or noop)
  -> actions.New(...)
  -> app.NewDesktop(..., actions)
```

这样 `internal/app` 保持薄封装，业务编排进入 `internal/actions`。

## 3. Behavior implemented by area

### 3.1 Knowledge

已实现：
- `CreateKnowledge`
- `GetKnowledge`
- `ListKnowledge`
- `GetKnowledgeTree`
- `UpdateKnowledge`
- `DeleteKnowledge`
- `MoveKnowledge`
- `ArchiveKnowledge`
- `UnarchiveKnowledge`

实现方式：
- 直接编排 `KnowledgeRepository`
- `GetKnowledgeTree(nil)` 当前通过全量 `List` 后在 app 层拼树
- 当前未补卡片数 / 到期数聚合，DTO 中先返回 0

### 3.2 Card

已实现：
- `CreateCard`
- `GetCard`
- `ListCards`
- `UpdateCard`
- `DeleteCard`
- `AddCardTags`
- `RemoveCardTags`
- `GetCardTags`
- `SuspendCard`
- `ResumeCard`

当前规则：
- `CreateCard` 会把 HTML 写入 `contracts.FileStore`，文件类型固定为 `card/<uuid>--<slug>.html`
- `article/excerpt/qa/cloze` 会创建默认 `CardSRS{fsrs_state:"new"}`
- `note` 当前不初始化 SRS
- `DeleteCard` 会先尝试把文件移入 trash，再做 DB 软删除
- `UpdateCard` 若携带 `HTMLContent`，会覆写文件并更新 `HTMLHash`

当前限制：
- `CreateCard` 还未把 DB 写入与文件写入放进统一事务/补偿链
- 搜索索引未联动
- `SourceDocument` 相关高级编排未接入

### 3.3 Review

已实现：
- `GetDueCards`
- `SubmitReview`
- `UndoLastReview`
- `GetSRSStatistics`
- `GetReviewHistory`
- `GetReviewStats`

实现方式：
- `SubmitReview`：加载 `Card + SRS` → 调 `FSRS.Review` → 调 `SRSRepository.UpdateAfterReview`
- `GetDueCards`：调用 `SRSRepository.GetDueCards(knowledgeID)` 后再补查 `Card`
- review 相关事务边界当前保持在仓储层：`SRSRepository.UpdateAfterReview` 内部统一写 `CardSRS + ReviewLog`

当前限制：
- `UndoLastReview` 直接复用仓储实现，后续若逻辑变复杂再下沉到 flow

### 3.4 Search / Tag-search

已实现：
- `SearchCardsByTags`

当前策略：
- 只走 SQLite / repository 路径
- 不引入 `SearchIndexer`
- 不实现正文全文搜索

实现说明：
- `SearchCardsByTags` 走 `CardRepository.List(TagIDs)`
- `ListCardOptions.TagIDs` 已在仓储层按多标签 AND 关系生效

当前限制：
- 目前只落地了“按标签搜索卡片”
- 标题关键词搜索仍通过 `ListCards(filters.Keyword)` 走 card list 能力，不单独拆 search action

## 4. Confirmed product rules reflected in code

以下确认已体现在当前实现或文档边界中：

- `article` 与 `excerpt/qa/cloze` 一样进入标准复习体系
- FSRS 参数第一阶段视为全局唯一，并已支持读取/更新默认参数
- `queue` 第一阶段搁置，不生成真实实现
- `importing` 第一阶段只保留未来扩展位，不按同步导入文档实现
- 搜索第一阶段只走 SQLite/DB 路径，不做索引引擎

## 5. Gaps intentionally left for next phase

### 5.1 Queue

未实现 `actions/queue`

原因：
- 当前已确认第一阶段搁置
- 持久化排序 / postpone / pin 等能力还没有落表方案

### 5.2 Importing

未实现 `actions/importing`

原因：
- `import_jobs` 真相源与状态流转尚未落模型/仓储
- `docs/wails-api-design.md` 里的同步 `ImportDocument` 设计和当前异步 job 方向不一致

### 5.3 SearchIndexer

未实现 `contracts.SearchIndexer`

原因：
- 当前 MVP 明确先走 SQLite/标题/标签能力
- `internal/indexing/indexing.go` 仍为后续阶段占位

### 5.4 Transaction boundary for multi-step card creation

当前 `CreateCard` 已升级为：
1. 校验 knowledge
2. 写 FileStore
3. 使用 `TransactionManager.WithTx` 包裹 Card / tags / SRS 写入
4. 若事务失败，对已创建文件做补偿删除

当前限制：
- 文件写入仍发生在 DB 事务之外，当前通过失败后补偿删除收口，而不是文件系统与 DB 的真正分布式一致性

## 6. Mismatch notes to keep in mind

### 6.1 With `docs/actions-layer-design.md`

当前代码已符合主方向，主要仍未进入本阶段的点是：
- `SearchIndexer` 仍仅保留为文档中的未来扩展点

### 6.2 With `docs/wails-api-design.md`

当前代码和文档仍存在以下差异：
- `ImportDocument` 尚未实现，不能按同步返回 `ImportResult` 理解为已落地
- 搜索当前实际落地的是 `SearchCardsByTags` + `ListCards(keyword)` 组合，而不是完整全文搜索能力

## 7. Next recommended implementation order

后续继续实现时，建议顺序保持为：

1. 再决定是否进入 importing / queue / SearchIndexer Phase 2

## 8. Verification done for current phase

当前已完成的定向验证：

```bash
go test ./internal/storage/repository ./internal/actions/... ./internal/app ./internal/bootstrap
```

目标：
- 确认仓储层过滤修复、action 拆分、system/fsrs_parameters 注入链和 app 暴露代码可编译
- 不扩大到无关包或全量构建
