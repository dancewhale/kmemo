# 任务单 A：前后端契约统一 + Wails API 适配层

你是一名资深 **Vue 3 / TypeScript / Wails 集成工程师**。
请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 前端项目中，实现 **"前后端契约统一 + Wails API 适配层" 第一阶段**。

这个项目是一个基于 **Wails v2 + Go 后端** 的桌面知识工具，产品定位参考 **SuperMemo 风格的渐进阅读 / 知识工作台**。
当前前端已经具备：

* `workspace` 模块
* `reader` 模块
* `knowledge-tree` 模块
* `card` 模块
* `editor` 模块
* `extract` 模块
* `review` 模块
* `search` 模块
* `command-center` 模块
* `settings` 模块
* `context-actions` 模块
* `inbox-capture` 模块
* `object-creation` 模块
* 若干 shared components、stores、mock 数据与交互闭环

当前问题是：

> 前端已经有较完整的模块结构，但仍主要依赖 mock 数据。
> 现在需要建立 **统一的前后端数据契约与 Wails API 适配层**，让后续每个业务模块都能逐步替换 mock，而不破坏已有前端结构。

---

## 一、后端真实 API 全景（关键基础，必须先理解）

> 本节描述的是 Go 后端 `internal/app/*.go` 中已经实现、由 Wails v2 向前端暴露的真实方法。
> Wails v2 调用约定：`window.go.App.方法名(参数)` → 对应 `desktop.App.方法名`。
> Wails 生成的 TS 绑定位于 `wailsjs/go/main/App.ts`（构建后自动生成）。

### 1.1 Knowledge（知识空间 / 卡组）— `knowledge.go`

知识空间是顶层容器，类似 Anki 的牌组（Deck），支持嵌套树形结构。
**前端的 `knowledge-tree` 模块核心数据来自这里**，而非来自假想的 `KnowledgeNode`（含 topic/article/extract 类型的扁平节点树不在后端）。

```
CreateKnowledge(req: CreateKnowledgeRequest)       → string (id), error
GetKnowledge(id: string)                           → KnowledgeDTO, error
ListKnowledge(parentID: string | null)             → KnowledgeDTO[], error
GetKnowledgeTree(rootID: string | null)            → KnowledgeTreeNode[], error
UpdateKnowledge(id: string, req: UpdateKnowledgeRequest) → error
DeleteKnowledge(id: string)                        → error
MoveKnowledge(id: string, newParentID: string | null)    → error
ArchiveKnowledge(id: string)                       → error
UnarchiveKnowledge(id: string)                     → error
```

后端 DTO（JSON 字段名均为 camelCase，与 Go json tag 一致）：

```typescript
interface KnowledgeDTO {
  id: string
  name: string
  description: string
  parentId: string | null
  cardCount: number
  dueCount: number
  createdAt: string   // ISO 8601
  updatedAt: string
  archivedAt: string | null
}

interface KnowledgeTreeNode extends KnowledgeDTO {
  children: KnowledgeTreeNode[]
}

interface CreateKnowledgeRequest {
  name: string
  description: string
  parentId: string | null
}

interface UpdateKnowledgeRequest {
  name: string
  description: string
}
```

### 1.2 Card（卡片）— `card.go`

卡片是核心学习对象。内容存储为 HTML（`htmlContent`），属于某个知识空间（`knowledgeId`），可嵌套（`parentId`），支持标签、排序、暂停/恢复。
**前端的 `card` 模块与 `object-creation` 模块创建卡片都要对接到这里。**

```
CreateCard(req: CreateCardRequest)                  → string (id), error
GetCard(id: string)                                 → CardDTO, error
GetCardDetail(id: string)                           → CardDetailDTO, error
ListCards(filters: CardFilters)                     → ListCardsResult, error
GetCardChildren(parentID: string)                   → CardDTO[], error
UpdateCard(id: string, req: UpdateCardRequest)      → error
DeleteCard(id: string)                              → error
MoveCard(req: MoveCardRequest)                      → error
ReorderCardChildren(req: ReorderCardChildrenRequest)→ error
AddCardTags(cardID: string, tagIDs: string[])       → error
RemoveCardTags(cardID: string, tagIDs: string[])    → error
GetCardTags(cardID: string)                         → TagDTO[], error
SuspendCard(cardID: string)                         → error
ResumeCard(cardID: string)                          → error
```

后端 DTO：

```typescript
interface CardDTO {
  id: string
  knowledgeId: string
  knowledgeName: string
  sourceDocumentId: string | null
  parentId: string | null
  title: string
  cardType: string          // 'basic' | 'cloze' | 'topic' 等后端枚举值
  htmlPath: string          // 文件系统路径（前端一般不直接使用）
  htmlContent: string       // 卡片 HTML 正文
  status: string            // 'active' | 'suspended'
  tags: TagDTO[]
  srs: SRSDTO | null
  createdAt: string
  updatedAt: string
}

interface CardSummaryDTO {
  id: string
  parentId: string | null
  title: string
  cardType: string
  status: string
  createdAt: string
  updatedAt: string
}

interface CardDetailDTO extends CardDTO {
  parent: CardSummaryDTO | null
  children: CardSummaryDTO[]
}

interface ListCardsResult {
  items: CardDTO[]
  total: number
}

interface CardFilters {
  knowledgeId?: string | null
  cardType?: string
  status?: string
  tagIds?: string[]
  keyword?: string
  parentId?: string | null
  isRoot?: boolean | null
  orderBy?: string          // 'title' | 'created_at' | 'updated_at' | 'sort_order'
  orderDesc?: boolean
  limit?: number
  offset?: number
}

interface CreateCardRequest {
  knowledgeId: string
  sourceDocumentId?: string | null
  parentId?: string | null
  title: string
  cardType: string
  htmlContent: string
  tagIds?: string[]
}

interface UpdateCardRequest {
  title: string
  htmlContent: string
  status: string
}

interface MoveCardRequest {
  cardId: string
  targetParentId: string | null
  targetIndex: number
}

interface ReorderCardChildrenRequest {
  knowledgeId: string
  parentId: string | null
  orderedChildIds: string[]
}
```

### 1.3 Review / SRS — `review.go`

后端 Review 基于 **FSRS 算法**，与前端 mock 的队列式模型有本质差异：
- 没有 "创建 ReviewItem" 的接口——卡片创建后自动纳入复习调度。
- 评分是**整数 1–4**（FSRS 标准：1=Again / 2=Hard / 3=Good / 4=Easy），**不是**字符串 `'again'|'hard'|'good'|'easy'`。
- `GetDueCards` 由后端 FSRS 调度算法决定哪些卡片到期，不是前端维护队列。

```
GetDueCards(knowledgeID: string | null, limit: number) → CardWithSRSDTO[], error
SubmitReview(req: ReviewRequest)                       → error
UndoLastReview(cardID: string)                         → error
GetSRSStatistics(knowledgeID: string | null)           → SRSStatisticsDTO, error
GetReviewHistory(cardID: string, limit: number)        → ReviewLogDTO[], error
GetReviewStats(startDate: string, endDate: string)     → ReviewStatistics, error
```

后端 DTO：

```typescript
interface SRSDTO {
  cardId: string
  fsrsState: string         // 'new' | 'learning' | 'review' | 'relearning'
  dueAt: string | null
  lastReviewAt: string | null
  stability: number | null
  difficulty: number | null
  reps: number
  lapses: number
}

interface CardWithSRSDTO {
  card: CardDTO
  srs: SRSDTO
}

interface SRSStatisticsDTO {
  newCount: number
  learningCount: number
  reviewCount: number
  relearningCount: number
  totalCards: number
  dueToday: number
}

interface ReviewRequest {
  cardId: string
  rating: number            // 1=Again / 2=Hard / 3=Good / 4=Easy（FSRS 标准）
}

interface ReviewLogDTO {
  id: string
  cardId: string
  reviewedAt: string
  rating: number
  reviewKind: string
}
```

### 1.4 Tag / 搜索 — `tag_search.go`

标签 CRUD + 按标签搜索卡片。目前后端没有全文搜索接口（仅 `ListCards` 的 `keyword` 字段做简单过滤）。

```
CreateTag(req: CreateTagRequest) → string (id), error
GetTag(id: string)               → TagDTO, error
ListTags()                       → TagDTO[], error
UpdateTag(id: string, req: UpdateTagRequest) → error
DeleteTag(id: string)            → error
SearchCardsByTags(tagIDs: string[]) → CardDTO[], error
```

后端 DTO：

```typescript
interface TagDTO {
  id: string
  name: string
  slug: string
  color: string
  icon: string
  description: string
  cardCount: number
  createdAt: string
}

interface CreateTagRequest {
  name: string
  slug: string
  color: string
  icon: string
  description: string
}

interface UpdateTagRequest {
  name: string
  color: string
  icon: string
  description: string
}
```

### 1.5 System / Settings — `system.go`

FSRS 参数管理，对应前端 `settings` 模块。

```
ListFSRSParameters()                                           → FSRSParameterDTO[], error
GetDefaultFSRSParameter()                                      → FSRSParameterDTO, error
UpdateDefaultFSRSParameter(req: UpdateDefaultFSRSParameterRequest) → FSRSParameterDTO, error
```

后端 DTO：

```typescript
interface FSRSParameterDTO {
  id: string
  name: string
  parametersJson: string   // JSON 字符串，序列化后的 FSRS 参数数组
  desiredRetention: number | null
  maximumInterval: number | null
}

interface UpdateDefaultFSRSParameterRequest {
  parametersJson: string
  desiredRetention: number | null
  maximumInterval: number | null
}
```

---

## 二、前后端概念映射（重要）

> 在对接前必须清楚以下映射关系，避免用错误的假设写代码。

| 前端模块 / 概念 | 后端实体 / API | 备注 |
|---|---|---|
| `knowledge-tree` 中的树形结构 | `KnowledgeTreeNode` via `GetKnowledgeTree()` | 后端 Knowledge 是"卡组空间"，不是"话题/文章/摘录/卡片"的混合树 |
| `card` 模块（CardItem） | `CardDTO` via Card API | 卡片归属于 `knowledgeId`，内容为 HTML |
| `review` 模块（ReviewItem / 队列） | `CardWithSRSDTO` via `GetDueCards()` | FSRS 调度，无"创建 ReviewItem"操作 |
| `review` 评分 `'again'|'hard'|'good'|'easy'` | `rating: 1|2|3|4` | 前端 UI 可保留字符串标签，Mapper 层转换为整数 |
| `search` 模块（全文搜索） | `ListCards({ keyword })` + `SearchCardsByTags()` | 目前后端无全文索引，仅简单关键词过滤 |
| `settings` 模块 | FSRS Parameter API | 仅 FSRS 参数配置 |
| `reader/article`（前端 mock） | **后端暂无对应接口** | article/extract/inbox 是前端 mock 阶段概念，暂不接后端 |
| `extract`（前端 mock） | **后端暂无对应接口** | 暂保留 mock |
| `inbox-capture`（前端 mock） | **后端暂无对应接口** | 暂保留 mock |

---

## 三、本次目标

1. 统一前后端 DTO / UI Model 边界，类型完全基于后端真实字段
2. 建立 `src/api/wails` 适配层（包装真实 Wails 生成函数）
3. 建立模块级 repository 封装
4. 建立统一错误处理、返回值规范、空值兼容策略
5. 为后续逐模块接入真实后端打基础

本次重点是 **契约、类型、封装和边界**，不是一次性接完全部功能。

---

## 四、明确范围限制

1. 不要直接实现所有页面的完整后端接入
2. 不要一次性重写所有 stores
3. 不要把当前所有 mock 全部删除（reader/article/extract/inbox 模块暂保留 mock）
4. 不要让组件直接 import Wails 生成代码
5. 不要把所有类型写成 `any`
6. 不要把 DTO 和 UI Model 混用
7. 不要在这一步里实现复杂缓存系统
8. 不要把 SRS 评分从整数转换遗漏（UI 用字符串标签，传给后端前 Mapper 转整数）

---

## 五、建议实现/修改的文件

```text
src/api/
  core/
    api-result.ts        — ApiResult<T> / ApiSuccess / ApiFailure 类型
    api-error.ts         — ApiErrorCode / ApiError 接口
    api-guard.ts         — ok() / fail() / safeWailsCall() 工具函数
  types/
    dto.ts               — 所有后端 DTO 类型（严格对应 Go JSON 字段）
    requests.ts          — 所有请求类型
  wails/
    index.ts             — 统一出口，export 全部 Wails 调用函数
    knowledge.ts         — Knowledge 相关 Wails 调用
    card.ts              — Card 相关 Wails 调用
    review.ts            — Review/SRS 相关 Wails 调用
    tag.ts               — Tag 相关 Wails 调用
    system.ts            — FSRS Parameter 相关 Wails 调用

src/modules/
  knowledge-tree/
    services/
      knowledge.repository.ts   — Knowledge 数据获取，返回 ApiResult<KnowledgeSpaceItem[]>
      knowledge.mapper.ts       — KnowledgeDTO → KnowledgeSpaceItem / UITreeNode
  card/
    services/
      card.repository.ts        — Card CRUD，返回 ApiResult<CardItem>
      card.mapper.ts            — 已存在，补全 CardDTO → CardItem 映射
  review/
    services/
      review.repository.ts      — GetDueCards / SubmitReview / GetSRSStatistics
      review.mapper.ts          — 已存在，补全 CardWithSRSDTO → ReviewItem 映射（含评分转换）
  search/
    services/
      search.repository.ts      — 对接 ListCards keyword + SearchCardsByTags
  settings/
    services/
      settings.repository.ts    — FSRS Parameter API
      settings.mapper.ts        — FSRSParameterDTO → SettingsItem
```

---

## 六、Wails 调用约定说明

### 真实调用方式

Wails v2 构建后会在 `wailsjs/go/main/App.ts` 生成 TypeScript 绑定，函数名与 Go 方法名一致：

```typescript
// wailsjs/go/main/App.ts（构建后自动生成，勿手写）
export function CreateKnowledge(req: CreateKnowledgeRequest): Promise<string>;
export function GetKnowledgeTree(rootID: string | null): Promise<KnowledgeTreeNode[]>;
export function ListCards(filters: CardFilters): Promise<ListCardsResult>;
export function GetDueCards(knowledgeID: string | null, limit: number): Promise<CardWithSRSDTO[]>;
export function SubmitReview(req: ReviewRequest): Promise<void>;
// ...
```

### 适配层要求

`src/api/wails/*.ts` 必须：
1. 通过 `import { 函数名 } from '../../../wailsjs/go/main/App'` 引入（或通过别名 `@wails/go/main/App`）
2. 对每个调用包裹 `safeWailsCall`，统一返回 `ApiResult<T>`
3. 组件和 store 层不直接 import `wailsjs/` 路径

**当前 `src/api/wails/index.ts` 是占位壳，需要真正实现。**

若当前开发环境 `wailsjs/` 尚不存在（未运行 `wails build` 或 `wails dev`），可用以下两种策略：

**策略 A（推荐）**：建立类型声明文件 + 运行时检测

```typescript
// src/api/wails/_wails-bridge.ts
declare const window: Window & {
  go?: {
    App?: Record<string, (...args: unknown[]) => Promise<unknown>>
  }
}

export async function callWails<T>(method: string, ...args: unknown[]): Promise<T> {
  const fn = window.go?.App?.[method]
  if (!fn) {
    throw new Error(`Wails method not available: App.${method}`)
  }
  return fn(...args) as Promise<T>
}
```

**策略 B**：条件导入 + mock fallback（dev 环境使用）

```typescript
const isDev = import.meta.env.DEV && !window.go?.App
```

---

## 七、类型与契约设计要求

### `src/api/core/api-result.ts`

```typescript
export type ApiSuccess<T> = { ok: true; data: T }
export type ApiFailure = { ok: false; error: ApiError }
export type ApiResult<T> = ApiSuccess<T> | ApiFailure
```

### `src/api/core/api-error.ts`

```typescript
export type ApiErrorCode =
  | 'UNKNOWN'
  | 'NOT_FOUND'
  | 'VALIDATION_ERROR'
  | 'CONFLICT'
  | 'EMPTY'
  | 'BACKEND_ERROR'
  | 'WAILS_NOT_READY'

export interface ApiError {
  code: ApiErrorCode
  message: string
  details?: unknown
}
```

### `src/api/core/api-guard.ts`

提供以下工具函数：

1. `ok<T>(data: T): ApiSuccess<T>`
2. `fail(code: ApiErrorCode, message: string, details?: unknown): ApiFailure`
3. `safeWailsCall<T>(fn: () => Promise<T>): Promise<ApiResult<T>>`
4. `isWailsAvailable(): boolean`

`safeWailsCall` 需处理：
- Wails 未初始化 → `WAILS_NOT_READY`
- Go 返回错误字符串（Wails v2 将 error 序列化为字符串异常）→ `BACKEND_ERROR`
- 未知 JS 异常 → `UNKNOWN`

### `src/api/types/dto.ts`

**必须严格对应后端 JSON 字段名**。复制第一节中的 TS 接口定义：

- `KnowledgeDTO`
- `KnowledgeTreeNode`
- `CardDTO`
- `CardSummaryDTO`
- `CardDetailDTO`
- `ListCardsResult`
- `SRSDTO`
- `CardWithSRSDTO`
- `SRSStatisticsDTO`
- `ReviewLogDTO`
- `TagDTO`
- `FSRSParameterDTO`

### `src/api/types/requests.ts`

**必须严格对应后端 JSON 字段名**。复制第一节中的 Request 接口：

- `CreateKnowledgeRequest`
- `UpdateKnowledgeRequest`
- `CreateCardRequest`
- `UpdateCardRequest`
- `MoveCardRequest`
- `ReorderCardChildrenRequest`
- `CardFilters`
- `ReviewRequest`（注意 `rating: number`，不是字符串）
- `CreateTagRequest`
- `UpdateTagRequest`
- `UpdateDefaultFSRSParameterRequest`

---

## 八、Wails 薄适配层设计要求

### `src/api/wails/knowledge.ts`

至少提供：

1. `getKnowledgeTree(rootID?: string | null): Promise<ApiResult<KnowledgeTreeNode[]>>`
2. `listKnowledge(parentID?: string | null): Promise<ApiResult<KnowledgeDTO[]>>`
3. `createKnowledge(req: CreateKnowledgeRequest): Promise<ApiResult<string>>`
4. `updateKnowledge(id: string, req: UpdateKnowledgeRequest): Promise<ApiResult<void>>`
5. `deleteKnowledge(id: string): Promise<ApiResult<void>>`
6. `moveKnowledge(id: string, newParentID: string | null): Promise<ApiResult<void>>`
7. `archiveKnowledge(id: string): Promise<ApiResult<void>>`

### `src/api/wails/card.ts`

至少提供：

1. `listCards(filters: CardFilters): Promise<ApiResult<ListCardsResult>>`
2. `getCard(id: string): Promise<ApiResult<CardDTO>>`
3. `getCardDetail(id: string): Promise<ApiResult<CardDetailDTO>>`
4. `createCard(req: CreateCardRequest): Promise<ApiResult<string>>`
5. `updateCard(id: string, req: UpdateCardRequest): Promise<ApiResult<void>>`
6. `deleteCard(id: string): Promise<ApiResult<void>>`
7. `getCardChildren(parentID: string): Promise<ApiResult<CardDTO[]>>`
8. `addCardTags(cardID: string, tagIDs: string[]): Promise<ApiResult<void>>`
9. `removeCardTags(cardID: string, tagIDs: string[]): Promise<ApiResult<void>>`
10. `suspendCard(cardID: string): Promise<ApiResult<void>>`
11. `resumeCard(cardID: string): Promise<ApiResult<void>>`

### `src/api/wails/review.ts`

至少提供：

1. `getDueCards(knowledgeID: string | null, limit: number): Promise<ApiResult<CardWithSRSDTO[]>>`
2. `submitReview(req: ReviewRequest): Promise<ApiResult<void>>`（注意 `rating` 是整数）
3. `undoLastReview(cardID: string): Promise<ApiResult<void>>`
4. `getSRSStatistics(knowledgeID: string | null): Promise<ApiResult<SRSStatisticsDTO>>`
5. `getReviewHistory(cardID: string, limit: number): Promise<ApiResult<ReviewLogDTO[]>>`

### `src/api/wails/tag.ts`

至少提供：

1. `listTags(): Promise<ApiResult<TagDTO[]>>`
2. `createTag(req: CreateTagRequest): Promise<ApiResult<string>>`
3. `updateTag(id: string, req: UpdateTagRequest): Promise<ApiResult<void>>`
4. `deleteTag(id: string): Promise<ApiResult<void>>`
5. `searchCardsByTags(tagIDs: string[]): Promise<ApiResult<CardDTO[]>>`

### `src/api/wails/system.ts`

至少提供：

1. `getDefaultFSRSParameter(): Promise<ApiResult<FSRSParameterDTO>>`
2. `updateDefaultFSRSParameter(req: UpdateDefaultFSRSParameterRequest): Promise<ApiResult<FSRSParameterDTO>>`

### 要求

1. 每个调用都通过 `safeWailsCall` 包装，统一返回 `ApiResult<T>`
2. 组件和 store 层**不直接** import `wailsjs/` 路径
3. 本文件只做薄封装，不做 UI Model 映射

---

## 九、Repository 设计要求

### `knowledge.repository.ts`

至少提供：

1. `fetchKnowledgeTree(rootID?: string | null): Promise<ApiResult<KnowledgeSpaceItem[]>>`
2. `createKnowledgeSpace(name: string, description?: string, parentID?: string | null): Promise<ApiResult<string>>`
3. `updateKnowledgeSpace(id: string, name: string, description: string): Promise<ApiResult<void>>`
4. `deleteKnowledgeSpace(id: string): Promise<ApiResult<void>>`

返回值使用 UI Model（`KnowledgeSpaceItem`），而非原始 DTO。

---

### `card.repository.ts`

至少提供：

1. `fetchCardsByKnowledge(knowledgeID: string, opts?: Partial<CardFilters>): Promise<ApiResult<{ items: CardItem[]; total: number }>>`
2. `fetchCard(id: string): Promise<ApiResult<CardItem>>`
3. `fetchCardDetail(id: string): Promise<ApiResult<CardDetailItem>>`
4. `createCard(input: CreateCardInput): Promise<ApiResult<string>>`
5. `updateCard(id: string, input: UpdateCardInput): Promise<ApiResult<void>>`
6. `deleteCard(id: string): Promise<ApiResult<void>>`

---

### `review.repository.ts`

至少提供：

1. `fetchDueCards(knowledgeID: string | null, limit?: number): Promise<ApiResult<ReviewQueueItem[]>>`
2. `submitReview(cardID: string, grade: ReviewGrade): Promise<ApiResult<void>>`  
   注意：`grade` 在 UI 层是 `'again'|'hard'|'good'|'easy'`，Mapper 在此转为 `1|2|3|4`
3. `undoLastReview(cardID: string): Promise<ApiResult<void>>`
4. `fetchSRSStatistics(knowledgeID: string | null): Promise<ApiResult<SRSStats>>`

---

### `settings.repository.ts`

至少提供：

1. `fetchDefaultFSRSParameter(): Promise<ApiResult<FSRSSettings>>`
2. `updateDefaultFSRSParameter(input: UpdateFSRSSettingsInput): Promise<ApiResult<FSRSSettings>>`

---

## 十、Mapper 设计要求

每个模块必须有 mapper，不允许 DTO 直接进组件。

### `knowledge.mapper.ts`

至少实现：

1. `mapKnowledgeDtoToSpaceItem(dto: KnowledgeDTO): KnowledgeSpaceItem`
2. `mapKnowledgeTreeNodeToUINode(node: KnowledgeTreeNode, level?: number): UIKnowledgeNode`  
   将后端树形 `KnowledgeTreeNode` 映射为前端 `knowledge-tree` 使用的 `UITreeNode` 结构

---

### `card.mapper.ts`（已有，需补全）

至少实现：

1. `mapCardDtoToCardItem(dto: CardDTO): CardItem`  
   注意：`htmlContent` → `content`（前端 card 模块使用 `content` 字段）；`srs.fsrsState` → `status`
2. `mapCardDtoToReviewItem(dto: CardWithSRSDTO): ReviewQueueItem`  
   这是 review 页面展示所需的 UI Model
3. `mapCreateCardInputToDto(input: CreateCardInput): CreateCardRequest`  
   注意填充 `knowledgeId`、`cardType`，将 `prompt`/`answer` 合并为 `htmlContent`

---

### `review.mapper.ts`（已有，需修正）

至少实现：

1. `mapCardWithSRSDtoToReviewItem(dto: CardWithSRSDTO): ReviewItem`  
   主要差异：`srs.dueAt` → `dueAt`；`srs.fsrsState` → `status`；`srs.reps` → `reviewCount`
2. `mapReviewGradeToRating(grade: 'again' | 'hard' | 'good' | 'easy'): 1 | 2 | 3 | 4`  
   **关键转换：** `again→1, hard→2, good→3, easy→4`
3. `mapSRSStatsDtoToStats(dto: SRSStatisticsDTO): ReviewStats`

---

### `settings.mapper.ts`

至少实现：

1. `mapFSRSParameterDtoToSettings(dto: FSRSParameterDTO): FSRSSettings`
2. `mapUpdateSettingsInputToDto(input: UpdateFSRSSettingsInput): UpdateDefaultFSRSParameterRequest`

---

## 十一、错误处理要求

### 要求

1. repository 必须返回 `ApiResult<T>`
2. store 接入时只处理 `result.ok` / `!result.ok`
3. UI 层不直接 try/catch 原始 Wails 调用
4. 对以下情况有明确策略：
   - 后端返回空数组 → list 类返回 ok + 空数组（不算错误）
   - 单个对象不存在 → `NOT_FOUND`
   - Wails 未初始化（非 Wails 环境 / dev 模式）→ `WAILS_NOT_READY`，可 fallback mock
   - 参数校验失败 → `VALIDATION_ERROR`
   - Go error 字符串 → `BACKEND_ERROR`

### Wails v2 错误格式说明

Go 方法返回 `error` 时，Wails v2 会将其序列化为 **JS 异常**（reject），错误内容为 Go error 的字符串表示（如 `"invalid input"` 或 `"record not found"`）。

`safeWailsCall` 的 catch 块应检测常见错误字符串：

```typescript
function classifyWailsError(err: unknown): ApiError {
  const msg = String(err)
  if (msg.includes('record not found') || msg.includes('not found')) {
    return { code: 'NOT_FOUND', message: msg }
  }
  if (msg.includes('invalid input') || msg.includes('validation')) {
    return { code: 'VALIDATION_ERROR', message: msg }
  }
  return { code: 'BACKEND_ERROR', message: msg }
}
```

---

## 十二、Mock 与真实后端共存方案

当前阶段不能删除所有 mock。

### 分类处理

| 模块 | 策略 |
|---|---|
| `knowledge-tree` | 本次接真实后端，删除 `mock/tree.ts` 对应逻辑 |
| `card` | 本次接真实后端，删除 `mock/cards.ts` 对应逻辑 |
| `review` | 本次接真实后端（`GetDueCards` 替换队列 mock） |
| `reader / article` | 暂保留 mock，后端当前无此接口 |
| `extract` | 暂保留 mock，后端当前无此接口 |
| `inbox-capture` | 暂保留 mock |
| `search` | 暂保留 mock，接入 `ListCards keyword` 过滤 |
| `settings` | 本次接真实后端（FSRS 参数） |

### 切换策略

```typescript
// 在 store 的 initialize() 中：
const USE_REAL_API = isWailsAvailable()  // 判断是否在 Wails 环境内运行

async initialize() {
  if (this.initialized) return
  if (USE_REAL_API) {
    const result = await knowledgeRepository.fetchKnowledgeTree()
    if (result.ok) {
      this.nodes = result.data
    } else {
      // 降级到 mock 或显示错误
      this.nodes = mockKnowledgeNodes
    }
  } else {
    this.nodes = mockKnowledgeNodes
  }
  this.initialized = true
}
```

---

## 十三、示范性接入要求（优先模块）

这一步不要一次改完所有 store，但必须至少做 **两个示范模块**。

### 推荐示范 1：`knowledge-tree` store

从 `mock/tree.ts`（前端 mock 节点树）切换为 `GetKnowledgeTree()` API。

注意处理：
- 后端 `KnowledgeTreeNode` 是"知识空间"树，不含 topic/article/extract/card 混合类型
- 前端 `UITreeNode` 需要 `level`、`hasChildren`、`childCount` 等字段，在 Mapper 层计算

### 推荐示范 2：`review` store

从前端自维护队列模型切换为 `GetDueCards()` + `SubmitReview()` API。

注意处理：
- 移除 `addMockReviewFromExtract` / `addMockReviewFromTreeNode` 等 mock 方法
- 评分 `'again'|'hard'|'good'|'easy'` 在 repository 层统一转换为整数 `1|2|3|4`
- `SRSStatisticsDTO` 对应前端 `ReviewStats`，字段名不同，需要 Mapper

---

## 十四、与现有模块的联动要求

1. `RightPane.vue`、`CenterPane.vue`、业务组件不应直接感知 DTO
2. 业务组件仍然使用原来的 UI Model（`CardItem`、`ReviewItem` 等）
3. `object-creation` 模块的 `submitCard()` 要改为调用 `card.repository.createCard()`
4. `context-actions` 中的操作（创建复习、暂停卡片等）可通过 repository 接入，当前阶段可先保留
5. `extract`、`reader`、`inbox-capture` 相关的 context actions 暂保留 mock

---

## 十五、UI 不应大改

1. 前端各页面 UI 不改
2. 示范模块可增加 loading / error 状态展示
3. 不重做组件样式
4. `review.store` 的 `submitGrade(grade: ReviewGrade)` 接口签名保持不变，内部改为调用 repository

---

## 十六、代码结构要求

1. DTO、Request、UI Model 分层清晰
2. repository 不返回 raw DTO
3. mapper 不写在组件里
4. store 不直接调用 raw Wails adapter
5. 组件不 import `src/api/wails/*` 或 `wailsjs/*`
6. 所有 import 必须完整
7. 不要只写伪代码或片段

---

## 十七、建议的实现策略与顺序

### 推荐顺序

1. `src/api/core/`（api-result / api-error / api-guard）
2. `src/api/types/dto.ts` + `requests.ts`（严格对应后端）
3. `src/api/wails/*.ts`（薄适配层，包裹 `safeWailsCall`）
4. `knowledge.mapper.ts` + `knowledge.repository.ts`
5. `review.mapper.ts` 修正（评分转换）+ `review.repository.ts`
6. `card.mapper.ts` 补全 + `card.repository.ts`
7. 示范：`knowledge-tree` store 接入 repository
8. 示范：`review` store 接入 repository
9. `settings.repository.ts` + `settings.mapper.ts`（轻量）

---

## 十八、输出顺序要求

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：说明分层架构

* raw wails adapter layer（`src/api/wails/`）
* contract / DTO layer（`src/api/types/`）
* mapper layer（`src/modules/*/services/*.mapper.ts`）
* repository layer（`src/modules/*/services/*.repository.ts`）
* store integration（示范模块）

### 第四步：逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text
src/api/core/api-result.ts
```

然后再给完整代码。

### 第五步：最后说明

1. 当前结构如何帮助后续逐模块接入真实后端
2. 哪些位置需要替换成真实 Wails 生成函数（`wailsjs/go/main/App.ts`）
3. 当前有哪些模块暂保留 mock，以及接入条件（后端何时支持）
4. 下一步最适合优先接入哪个业务模块

---

## 十九、补充约束

```text
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入复杂请求库（无 axios / react-query）
- DTO 字段名必须与后端 Go json tag 完全一致（camelCase）
- 评分转换（grade 字符串 → rating 整数）必须在 Mapper 层完成，禁止散落在组件
- knowledge-tree 前端当前含 'topic'|'article'|'extract'|'card' 混合节点类型，
  对接后端时只映射 KnowledgeTreeNode（纯知识空间树），
  卡片列表通过 ListCards API 另行获取，不要强行把两者合并进同一棵树
- 当前阶段优先保证结构清晰、边界清晰、可逐步替换 mock
```

---

## 二十、现在开始执行

请从"本次会修改/新增的文件列表"开始，然后按要求继续输出完整代码。

---

## 完成后的建议对接顺序

完成本任务后，建议按下面顺序开始真实对接：

1. **Knowledge Tree（知识空间树）**
   - `GetKnowledgeTree()` 替换前端 mock 节点树
   - `CreateKnowledge` / `UpdateKnowledge` / `DeleteKnowledge`

2. **Card（卡片）**
   - `ListCards` 按知识空间加载卡片列表
   - `GetCardDetail` 在右侧详情面板展示
   - `CreateCard` / `UpdateCard`（含 HTML 内容）
   - `SuspendCard` / `ResumeCard`（前端卡片状态切换）

3. **Review（复习队列）**
   - `GetDueCards` 替换前端 mock 队列
   - `SubmitReview`（整数评分）
   - `GetSRSStatistics` 展示复习统计

4. **Tag（标签）**
   - `ListTags` / `CreateTag` / `DeleteTag`
   - 卡片标签联动

5. **Settings（FSRS 参数）**
   - `GetDefaultFSRSParameter` / `UpdateDefaultFSRSParameter`

6. **Search（全文搜索）**
   - 当前仅 `ListCards keyword` 过滤，真正全文搜索后端尚未实现

7. **Reader / Extract / Inbox（待后端支持）**
   - 当前后端无对应接口，暂保留 mock
