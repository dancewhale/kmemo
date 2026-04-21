# 任务单 7：复习队列模块（Review Queue）第一阶段

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **“复习队列模块（Review Queue）第一阶段可运行版本”**。

这个项目是一个基于 **Wails + Go 后端** 的桌面知识工具，产品定位参考 **SuperMemo 风格的渐进阅读 / 知识工作台**。
当前项目已经具备：

* `workspace` 模块
* `reader` 模块
* `knowledge-tree` 模块
* `editor` 模块壳
* `extract` 模块
* 文章中选区创建摘录
* 摘录进入知识树
* 摘录详情编辑
* 阅读区与知识区联动
* `tree.store.ts`
* `reader.store.ts`
* `editor.store.ts`
* `extract.store.ts`
* `mock/` 数据
* 基础三栏工作台和右侧详情区

本次目标是：**实现复习队列第一阶段前端闭环**，让用户可以看到待复习项、打开复习卡、进行评分、切换下一项，并看到基础统计信息。

当前阶段重点是 **Review Queue 的对象流、状态流和 UI 结构**，不是复杂复习算法。

---

## 一、模块目标

请实现如下复习流程：

1. 系统中存在一组待复习项（先使用 mock 数据）
2. 用户进入 `review` 工作区
3. 中间区域显示复习队列列表
4. 右侧区域显示当前复习项详情
5. 用户可以看到：

   * 标题
   * prompt / question
   * answer / note / quote
   * 来源信息
   * due 信息
   * 当前状态
6. 用户点击评分按钮，例如：

   * Again
   * Hard
   * Good
   * Easy
7. 提交评分后：

   * 当前项记录一次 review 结果
   * 队列状态更新
   * 自动切换到下一项
8. 页面显示今日统计，例如：

   * 今日已完成数
   * 剩余数
   * Again 次数
9. 若队列为空，显示空状态

当前阶段 **不要求真实调度算法**，但要为未来后端调度接入预留清晰边界。

---

## 二、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现真正的 SM-2 / SuperMemo 算法
3. 不要实现复杂卡片模板编辑器
4. 不要实现 cloze 生成
5. 不要实现图谱式复习关系
6. 不要实现复杂筛选器面板
7. 不要实现复习历史图表系统
8. 不要实现多会话 tab
9. 不要把所有逻辑塞进 `RightPane.vue`
10. 不要让复习 store 和 extract/tree/editor store 职责混乱

---

## 三、本次建议实现/修改的文件

请至少实现和补全以下文件：

```text id="y7j4rc"
src/modules/review/
  components/
    ReviewQueue.vue
    ReviewQueueItem.vue
    ReviewCard.vue
    ReviewActions.vue
    ReviewStatsBar.vue
    ReviewMetaSection.vue
  stores/
    review.store.ts
  types.ts
  services/
    review.mapper.ts
```

并根据需要补充或修改以下文件：

```text id="z93tfe"
src/modules/workspace/components/CenterPane.vue
src/modules/workspace/components/RightPane.vue
src/pages/ReviewPage.vue
src/mock/review.ts
src/mock/index.ts
src/modules/extract/stores/extract.store.ts
src/modules/reader/stores/reader.store.ts
src/modules/knowledge-tree/stores/tree.store.ts
```

如果有必要，你可以新增少量辅助文件，例如：

```text id="n5hkh0"
src/shared/constants/review.ts
src/shared/utils/review.ts
```

但请保持克制，不要扩展过大。

---

## 四、整体设计要求

请把这次能力拆成 4 个清晰层次：

### 1. 复习对象层

负责：

* Review item 数据结构
* 当前待复习队列
* 当前选中复习项
* 评分结果记录
* 今日统计

### 2. 复习队列展示层

负责：

* 中间列表展示
* 当前项高亮
* 剩余数量
* 空状态

### 3. 复习详情与作答层

负责：

* 展示 prompt / answer / quote / note
* 展示来源信息
* 展示评分按钮
* 响应评分动作

### 4. 与其他模块的轻量联动层

负责：

* 从复习项跳到来源文章
* 从复习项跳到关联摘录/节点
* 保持与 extract/tree/reader 的轻量协作

要求边界清晰，便于未来接真实后端调度。

---

## 五、类型设计要求

请先整理和定义清晰类型。

---

### `src/modules/review/types.ts`

至少包含：

```ts id="x0o6q0"
export type ReviewItemType = 'extract' | 'card'
export type ReviewStatus = 'new' | 'learning' | 'review'
export type ReviewGrade = 'again' | 'hard' | 'good' | 'easy'

export interface ReviewItem {
  id: string
  type: ReviewItemType
  title: string
  prompt: string
  answer: string
  sourceArticleId?: string | null
  sourceArticleTitle?: string | null
  extractId?: string | null
  nodeId?: string | null
  dueAt: string
  status: ReviewStatus
  lastReviewedAt?: string | null
  reviewCount: number
}

export interface ReviewRecord {
  itemId: string
  grade: ReviewGrade
  reviewedAt: string
}

export interface ReviewStats {
  completedToday: number
  remaining: number
  againCount: number
}
```

你可以在此基础上微调，但必须满足：

1. 复习项对象独立存在
2. 能关联摘录、知识树节点和原文章
3. 能记录评分历史的最小必要信息
4. 为未来后端 DTO 留空间

---

## 六、mock 数据要求

请在 `src/mock/review.ts` 中提供一批更真实的 mock 复习数据，至少包含：

* 6~10 个复习项
* 类型至少包括：

  * extract
  * card
* 状态至少包括：

  * new
  * learning
  * review

建议内容风格与前面的知识工具场景一致，例如：

* What is incremental reading?
* Why should notes be atomic?
* What is the role of active recall?
* Extract review about progressive summarization
* Why should quote and note be separated?

每条复习项至少要有：

* id
* title
* prompt
* answer
* dueAt
* status
* reviewCount

尽量补充：

* sourceArticleId
* sourceArticleTitle
* extractId
* nodeId

这些数据应当能与现有 mock article / extract / tree 数据形成联动。

---

## 七、store 设计要求

---

### `review.store.ts`

请新增或重构 review store，使其真正承担复习队列状态管理。

至少包含这些状态：

```ts id="g4ij7d"
interface ReviewState {
  items: ReviewItem[]
  selectedReviewItemId: string | null
  records: ReviewRecord[]
  loading: boolean
  submitting: boolean
}
```

至少包含这些 getter：

1. `queueItems`
2. `selectedItem`
3. `remainingCount`
4. `completedTodayCount`
5. `againCount`
6. `stats`
7. `isQueueEmpty`

至少包含这些 action：

1. `initialize()`
2. `setSelectedItem(id: string | null)`
3. `submitGrade(grade: ReviewGrade)`
4. `openFirstAvailable()`
5. `openNext()`
6. `getItemById(id: string)`
7. `resetSession()`（可选）

### 行为要求

* 初始化时加载 mock 数据
* 默认打开第一个可用复习项
* `submitGrade(grade)` 当前阶段可采用 mock 规则：

  * 记录一条 `ReviewRecord`
  * 当前项 `reviewCount + 1`
  * `lastReviewedAt = now`
  * 将当前项从待复习队列中移除，或标记为已完成（你选一种更清晰简单的方案并说明）
  * 自动切到下一项
* `stats.remaining` 应随着评分而变化
* Again / Hard / Good / Easy 目前只需要体现在 record 中和 again 统计中，不必做真实排期算法
* store 职责清晰，不要混进 editor/tree 的复杂逻辑

### 推荐简化方案

当前阶段最推荐：

* 将“本轮已评分项”从当前 session queue 中移除
* 原始 item 可仍存在 items 中，但 queue getter 只返回未完成项
* 这样最容易实现前端闭环

---

## 八、`review.mapper.ts` 要求

请新增 `src/modules/review/services/review.mapper.ts`，至少包含：

1. `buildReviewStats(...)`
2. `normalizeReviewItem(...)`
3. `buildReviewRecord(...)`

要求：

* 不要把映射和统计逻辑散落在组件里
* 尽量保持纯函数
* 便于未来接后端返回值

---

## 九、组件要求

---

### 1. `ReviewQueue.vue`

这是中间区域的复习队列主容器。

要求：

1. 从 `review.store.ts` 读取队列
2. 首次挂载时调用 `initialize()`
3. 渲染 `ReviewStatsBar.vue`
4. 渲染复习项列表
5. 使用 `ReviewQueueItem.vue` 渲染单项
6. 无数据时显示 `AppEmpty`
7. 加载态时显示 `AppLoading`

风格要求：

* 紧凑
* 像工作队列
* 不要表格化
* 不要后台管理系统样式

---

### 2. `ReviewQueueItem.vue`

这是中间区域单个复习项列表项。

props 至少包括：

* `item: ReviewItem`
* `selected: boolean`

行为要求：

1. 点击列表项时选中当前复习项
2. 展示标题
3. 展示类型
4. 展示状态
5. 展示 due 信息
6. 当前选中项高亮
7. hover 有轻微反馈

建议视觉结构：

* 第一行：标题 + 类型/状态 tag
* 第二行：prompt 摘要
* 第三行：dueAt / reviewCount / 来源信息摘要

---

### 3. `ReviewCard.vue`

这是右侧复习详情主面板，最核心组件之一。

要求：

1. 从 `review.store.ts` 读取当前选中复习项
2. 无当前项时显示 `AppEmpty`
3. 至少展示这些区域：

   * 标题
   * 类型
   * prompt
   * answer
   * 来源信息
   * due / status / reviewCount
4. 提供“显示答案”的交互
5. 默认先展示 prompt
6. 用户点击“Show Answer”后再显示 answer
7. 下方渲染 `ReviewActions.vue`

### 交互要求

* 初始只显示问题面
* 点击 “Show Answer” 后显示答案区
* 再允许点击评分按钮
* 评分后自动切下一项

这个流程很重要，能让复习更像真正的记忆工作流。

---

### 4. `ReviewActions.vue`

这是评分按钮区。

要求：

1. 提供 4 个按钮：

   * Again
   * Hard
   * Good
   * Easy
2. 当前阶段仅在“答案已显示”后启用
3. 点击按钮后调用 `review.store.submitGrade(...)`
4. 提交中有 loading 状态
5. 风格紧凑、明确、可连续点击

可以使用轻量颜色区分，但不要过度花哨。

---

### 5. `ReviewStatsBar.vue`

这是中间区域顶部的统计条。

要求：

至少展示：

* Completed today
* Remaining
* Again count

要求：

* 紧凑
* 工作台风格
* 一眼可读

---

### 6. `ReviewMetaSection.vue`

这是右侧详情区中用于展示复习项元信息的小组件。

至少展示：

* item id
* dueAt
* status
* reviewCount
* lastReviewedAt

要求：

* 轻量
* 可复用
* 信息密度高
* 专业克制

---

## 十、Workspace 集成要求

请将 review 模块集成到当前工作台中。

---

### 对 `CenterPane.vue` 的要求

当 `context === 'review'` 时：

* 中间区域应显示 `ReviewQueue.vue`
* 不再只是简单占位文本

---

### 对 `RightPane.vue` 的要求

当 `context === 'review'` 时：

* 右侧应显示 `ReviewCard.vue`
* 不再只是简单详情占位
* 不要把 review 逻辑写死在 `RightPane.vue`
* `RightPane.vue` 只负责根据当前 context 选择渲染 review 详情组件

---

### 对 `ReviewPage.vue` 的要求

* 继续使用 `WorkspaceShell`
* context 为 `review`
* 页面本身保持轻量，不写复杂业务逻辑

---

## 十一、与其他模块的联动要求

当前阶段只需要做 **轻量联动**，不要过度复杂。

---

### 联动 1：从复习项跳到来源文章

在 `ReviewCard.vue` 中，如果当前 item 有 `sourceArticleId`，请提供一个按钮：

* `Open Source Article`

点击后：

1. `reader.store.openArticleById(sourceArticleId)`
2. `workspace.store.setContext('reading')`
3. 右侧打开文章编辑区

---

### 联动 2：从复习项跳到关联摘录

如果当前 item 有 `extractId`，请提供一个按钮：

* `Open Extract`

点击后：

1. `extract.store.setSelectedExtract(extractId)`
2. 若有关联 tree node，则同步选中树节点
3. `workspace.store.setContext('knowledge')`
4. 右侧打开摘录详情

---

### 联动 3：从复习项跳到关联知识节点

如果当前 item 有 `nodeId`，可提供轻量入口：

* `Reveal in Tree`

点击后：

1. `tree.store.setSelectedNode(nodeId)`
2. `workspace.store.setContext('knowledge')`

如果当前你认为 `Open Extract` 已能覆盖绝大多数场景，也可以将这项做成可选能力，并说明理由。

---

## 十二、队列逻辑要求

当前阶段你需要明确一个“第一阶段前端版”的队列行为。

### 推荐简化逻辑

* 初始化后，`queueItems = items 中未在本轮 session 完成的项`
* 用户评分后：

  * 写入 `records`
  * 将该项标记为本轮已完成
  * `remainingCount - 1`
  * 自动打开下一项
* 若全部完成：

  * 中间显示空状态
  * 右侧显示 “All reviews completed for now”

### Again 的特殊处理

当前阶段可简化为：

* 仍然算一次完成提交
* 但 `againCount + 1`
* 不需要重新插回队列

请在最终说明中明确指出这是为了第一阶段简化。

---

## 十三、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 紧凑、专业
3. 轻边框
4. 小圆角
5. 克制的状态色
6. 更像知识工具，不像后台管理系统
7. Review queue 要像待处理工作清单
8. Review card 要像真正的复习工作面板，不像营销页或表单页

建议视觉结构：

### 中间区

* 顶部统计条
* 下方队列列表

### 右侧区

* 标题与状态
* prompt 区
* answer 区（初始隐藏）
* 元信息区
* 操作按钮区

---

## 十四、代码结构要求

请特别遵守：

1. 不要把所有 review 逻辑放进 `ReviewCard.vue`
2. 不要把统计逻辑直接写死在模板里
3. 不要让 `RightPane.vue` 变成业务大杂烩
4. `review.store.ts` 负责状态流
5. `review.mapper.ts` 负责映射与统计
6. `ReviewActions.vue` 只负责评分按钮区
7. 所有 import 必须完整
8. 不要只写伪代码或片段
9. 代码应为“第一阶段真实可运行版本”

---

## 十五、建议的实现策略

为了降低复杂度，推荐采用以下实现方案：

### 推荐方案

* `ReviewQueue.vue` 负责中间列表和顶部统计
* `ReviewCard.vue` 负责右侧详情
* `ReviewCard.vue` 内部维护一个轻量 `showAnswer` 状态
* `submitGrade()` 后自动重置 `showAnswer = false`
* 评分后自动打开下一项
* 若无下一项，显示完成状态

### 为什么推荐这样做

因为它能以最小复杂度模拟真实复习流程：

* 先看题
* 再看答案
* 再评分
* 再切题

这已经足够有 SuperMemo 的基础感了。

---

## 十六、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text id="e7rqmb"
src/modules/review/types.ts
```

然后再给完整代码。

### 第四步：最后说明：

1. 如何接入当前项目
2. 用户如何在界面中完成一轮复习
3. 当前实现相对真实后端调度版的简化点有哪些

---

## 十七、补充约束

```text id="ak4wdv"
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入复杂图表或额外状态库
- 可以使用少量 Element Plus Button / Tag / Divider / Empty / Card 等基础控件
- 当前阶段优先保证结构清晰、模块边界清晰、可运行
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 十八、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。
