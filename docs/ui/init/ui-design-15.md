# 任务单 16：卡片详情编辑 + 卡片/摘录转复习对象工作流第一阶段

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **“卡片详情编辑 + 卡片/摘录转复习对象工作流第一阶段可运行版本”**。

这个项目是一个基于 **Wails + Go 后端** 的桌面知识工具，产品定位参考 **SuperMemo 风格的渐进阅读 / 知识工作台**。
当前项目已经具备：

* `workspace` 模块
* `reader` 模块
* `knowledge-tree` 模块
* `editor` 模块
* `extract` 模块
* `review` 模块
* `command-center` 模块
* `search` 模块
* `settings` 模块
* `context-actions` 模块
* `inbox-capture` 模块
* `object-creation` 模块
* 阅读 / 摘录 / 知识树 / 复习 / 搜索 / 命令 / 设置 / 右键菜单 / 入箱 / 对象创建 的基础联动
* `tree.store.ts`
* `extract.store.ts`
* `review.store.ts`
* `workspace.store.ts`
* `reader.store.ts`
* `editor.store.ts`
* `mock/` 数据
* 统一 UI 风格与共享组件

本次目标是：**实现第一阶段的卡片详情编辑与“卡片/摘录 -> 复习对象”工作流**，让用户能够：

1. 打开并编辑卡片详情
2. 将卡片加入复习队列
3. 将摘录转换为复习对象
4. 从复习对象跳回来源卡片/摘录/原文
5. 形成“知识对象 -> 复习对象”的最小可运行闭环

当前阶段重点是：

* card 对象详情承载
* extract -> review object 转换
* card -> review object 转换
* review item 与来源对象的关联
* 最小可运行联动

不是复杂卡片模板系统、cloze 系统、SM 算法系统、复杂富文本卡片编辑器。

---

## 一、模块目标

请实现如下能力：

1. 用户点击知识树中的 **card 节点** 时，右侧能够打开 **卡片详情编辑面板**
2. 卡片详情至少支持编辑：

   * title
   * prompt
   * answer
3. 用户可以从卡片详情中执行：

   * Add to Review
   * Open in Review（若已存在 review item）
4. 用户可以从摘录详情中执行：

   * Create Review Item from Extract
   * Open Review Item（若已存在）
5. 当卡片或摘录被加入复习后：

   * review queue 中出现对应复习对象
   * review item 能记录其来源：

     * card node / extract
     * source article（如有）
6. 用户在 review 中可跳回：

   * 来源摘录
   * 来源卡片
   * 来源文章
7. Search 和右键菜单可感知这些对象关系
8. 结构要为未来接入真实后端 review object 创建逻辑留出清晰边界

---

## 二、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现复杂 cloze 删除卡
3. 不要实现复杂卡片模板系统
4. 不要实现多面卡 / 卡片模板编辑器
5. 不要实现复杂富文本 prompt/answer 编辑
6. 不要实现真实 SM-2 / SuperMemo 排期算法
7. 不要实现复杂 backlinks 图谱
8. 不要实现复杂批量转复习对象
9. 不要把所有逻辑塞进 `RightPane.vue`
10. 不要让 tree / extract / review / card 编辑状态混成一锅

---

## 三、本次建议实现/修改的文件

请至少实现和补全以下文件：

```text id="w33y73"
src/modules/card/
  components/
    CardDetailPanel.vue
    CardMetaSection.vue
    CardReviewSection.vue
    CardEditorForm.vue
  stores/
    card.store.ts
  types.ts
  services/
    card.mapper.ts
```

并根据需要补充或修改以下文件：

```text id="vromh9"
src/modules/extract/components/ExtractDetailPanel.vue
src/modules/extract/stores/extract.store.ts
src/modules/review/stores/review.store.ts
src/modules/review/types.ts
src/modules/review/services/review.mapper.ts
src/modules/knowledge-tree/stores/tree.store.ts
src/modules/workspace/components/RightPane.vue
src/modules/workspace/components/CenterPane.vue
src/modules/context-actions/services/context-menu.registry.ts
src/modules/context-actions/services/context-menu.executor.ts
src/modules/command-center/services/command.registry.ts
src/modules/command-center/services/command.executor.ts
src/modules/search/services/search.mapper.ts
src/modules/search/services/search.index.ts
src/mock/review.ts
src/mock/tree.ts
src/mock/index.ts
```

如果有必要，你可以新增少量辅助文件，例如：

```text id="5qc31i"
src/shared/constants/card.ts
src/shared/utils/review-link.ts
```

但请保持克制，不要扩展成完整 object relation framework。

---

## 四、整体设计要求

请把这次能力拆成 4 个清晰层次：

### 1. 卡片对象层

负责：

* card 数据结构
* card detail 状态
* card 编辑与更新
* card 与 tree node 的关系

### 2. 复习对象构建层

负责：

* 从 card 构建 review item
* 从 extract 构建 review item
* 建立来源关联字段

### 3. 右侧详情承载层

负责：

* 当前选中 card 节点时显示 `CardDetailPanel`
* 当前选中 extract 时显示 `ExtractDetailPanel`
* 根据对象类型选择正确详情面板

### 4. review 联动层

负责：

* review item 与 card/extract/article 的来回跳转
* 避免 review item 成为孤立对象

要求边界清晰，便于未来接真实后端。

---

## 五、类型设计要求

请先整理和定义清晰类型。

---

### `src/modules/card/types.ts`

至少包含：

```ts id="e1yvq8"
export interface CardItem {
  id: string
  nodeId: string
  title: string
  prompt: string
  answer: string
  parentNodeId: string | null
  createdAt: string
  updatedAt: string
  reviewItemId?: string | null
}

export interface CardUpdatePayload {
  title?: string
  prompt?: string
  answer?: string
}
```

你可以在此基础上微调，但必须满足：

1. card 对象与 tree node 可关联
2. card 有可编辑内容
3. card 可关联一个 reviewItemId
4. 不要过度复杂

---

### 对 `review/types.ts` 的补充要求

请确保 `ReviewItem` 至少支持这些来源关联字段：

```ts id="gaav1c"
export interface ReviewItem {
  id: string
  type: 'extract' | 'card'
  title: string
  prompt: string
  answer: string
  sourceArticleId?: string | null
  sourceArticleTitle?: string | null
  extractId?: string | null
  nodeId?: string | null
  cardId?: string | null
  dueAt: string
  status: 'new' | 'learning' | 'review'
  lastReviewedAt?: string | null
  reviewCount: number
}
```

要求：

1. extract review item 和 card review item 都能表达
2. card review item 能关联 `cardId` 和 `nodeId`
3. extract review item 能关联 `extractId`
4. 为未来后端 DTO 留空间

---

### 对 `extract/types.ts` 的补充要求

请确保 `ExtractItem` 至少可支持：

```ts id="ca0mh5"
reviewItemId?: string | null
```

因为第一阶段里，一个 extract 可能被转成 review item。

---

## 六、映射要求

请新增：

### `card.mapper.ts`

至少实现：

1. `buildCardItemFromTreeNode(...)`
2. `buildReviewItemFromCard(card: CardItem): ReviewItem`

要求：

* 不要把 review item 构造规则写在组件里
* 从 card -> review 的映射清晰
* `prompt` 用 card.prompt
* `answer` 用 card.answer
* `title` 可用 card.title

---

### 对 `review.mapper.ts` 的补充要求

请新增或补充：

1. `buildReviewItemFromExtract(...)`
2. `buildReviewItemFromCard(...)`

如果你认为 `buildReviewItemFromCard` 放 `card.mapper.ts` 更合理，也可以，但必须保持结构清晰，并在最终说明中说明。

### 推荐规则

#### 从 extract 构建 review item

* type = `extract`
* title = extract.title
* prompt = extract.quote
* answer = extract.note 或空字符串 / 默认说明
* sourceArticleId/sourceArticleTitle 从 extract 继承
* extractId 填充

#### 从 card 构建 review item

* type = `card`
* title = card.title
* prompt = card.prompt
* answer = card.answer
* cardId = card.id
* nodeId = card.nodeId

---

## 七、store 设计要求

---

### 1. `card.store.ts`

请新增独立 card store，不要把 card 编辑状态塞进 tree store 或 review store。

至少包含这些状态：

```ts id="v9i6ee"
interface CardState {
  items: CardItem[]
  selectedCardId: string | null
  saving: boolean
}
```

至少包含这些 getter：

1. `selectedCard`
2. `getCardById(id: string)`
3. `getCardByNodeId(nodeId: string)`

至少包含这些 action：

1. `setItems(items: CardItem[])`
2. `addCard(card: CardItem)`
3. `setSelectedCard(id: string | null)`
4. `openCardByNodeId(nodeId: string)`
5. `updateCard(id: string, payload: CardUpdatePayload)`
6. `attachReviewItem(cardId: string, reviewItemId: string)`
7. `getOrCreateFromNode(nodeId: string)`（可选，但推荐）

### 行为要求

* 点击 card 节点时，可根据 nodeId 找到 card 对象
* 若 object-creation 模块创建 card 时已经写入 card store，则直接读取
* 更新 card 标题时，应同步 tree node title
* 不要让 card.store 负责 tree 展开等逻辑

---

### 2. `review.store.ts`

请补充以下能力：

1. `addReviewItem(item: ReviewItem)`
2. `getReviewItemByCardId(cardId: string)`
3. `getReviewItemByExtractId(extractId: string)`
4. 可选：

   * `openReviewItemByCardId(cardId: string)`
   * `openReviewItemByExtractId(extractId: string)`

### 行为要求

* 若 card/extract 已经存在对应 review item，则不应重复创建太多重复对象
* 第一阶段推荐做法：

  * 若已存在，就直接打开已有 review item
  * 若不存在，再创建新的

这样更接近真实工具行为。

---

### 3. `extract.store.ts`

请补充以下能力：

1. `attachReviewItem(extractId: string, reviewItemId: string)`
2. `getReviewItemIdByExtractId(extractId: string)`（可选）

### 行为要求

* extract 被转成 review item 后，能保留关联
* 从 extract detail 可直接打开 review item

---

### 4. `tree.store.ts`

请确保具备：

1. `setSelectedNode(id: string | null)`
2. `updateNodeTitle(nodeId: string, title: string)`
3. `getNodeById(id: string)` 或等价 getter

要求：

* card 标题修改时可同步 tree node 标题
* 右侧能根据当前节点 type 判断是 card / extract / topic / article

---

## 八、右侧详情承载要求

---

### `RightPane.vue`

请根据当前上下文和当前选中节点类型做更明确分流。

要求：

#### 当 `context === 'knowledge'`

* 如果当前选中节点类型是 `extract`

  * 显示 `ExtractDetailPanel.vue`
* 如果当前选中节点类型是 `card`

  * 显示 `CardDetailPanel.vue`
* 否则

  * 显示普通知识节点详情占位或已有详情面板

#### 当 `context === 'reading'`

* 保持文章 `EditorShell.vue`

#### 当 `context === 'review'`

* 保持 `ReviewCard.vue`

重点：

* `RightPane.vue` 只负责“路由式分流”
* 不要把 card / extract / review 业务逻辑写死在这里

---

## 九、卡片详情组件要求

---

### 1. `CardDetailPanel.vue`

这是本次最核心的新组件之一。

要求：

1. 用于右侧展示当前选中 card 的详情
2. 从 `card.store.ts` 读取当前 card
3. 至少展示以下区域：

   * title
   * prompt
   * answer
   * parentNodeId
   * createdAt / updatedAt
   * review 状态区
4. title / prompt / answer 可编辑
5. 保存行为可简化为：

   * 输入即更新并标记 changed
   * 或点击 Save 按钮后更新
   * 请选择一种更适合第一阶段的方案并说明

### 推荐第一阶段方案

* 使用明确的 Save 按钮
* 因为 card 是结构化对象，显式保存更容易控制状态

6. 提供动作按钮：

   * Save
   * Add to Review
   * Open Review Item（若已存在）

7. 风格要求：

   * 像知识工具的对象详情面板
   * 不是后台表单页
   * 信息层级清楚

---

### 2. `CardEditorForm.vue`

职责：

* 承担 card 的基础可编辑表单区域

至少包括：

* title
* prompt
* answer

要求：

* 组件边界清晰
* 便于后续替换成更高级 card editor
* 不要让 `CardDetailPanel.vue` 过度臃肿

---

### 3. `CardMetaSection.vue`

职责：

* 展示 card 元信息

至少包括：

* card id
* node id
* parent node id
* createdAt
* updatedAt
* review status（是否已有 review item）

要求：

* 轻量
* 高信息密度
* 与 extract meta 风格保持一致

---

### 4. `CardReviewSection.vue`

职责：

* 展示 card 与 review 的关系及操作入口

至少包括：

* 是否已加入 review
* `Add to Review`
* `Open Review Item`

要求：

* 若已有 review item，则不要重复创建
* 动作反馈清晰
* 风格紧凑

---

## 十、摘录详情增强要求

---

### `ExtractDetailPanel.vue`

请在现有基础上补充一个 review 相关区块。

要求：

1. 展示当前 extract 是否已有 review item
2. 提供按钮：

   * `Create Review Item`
   * `Open Review Item`
3. 若已有 review item，则优先显示 `Open Review Item`
4. 点击 `Create Review Item` 时：

   * 若不存在 review item，则创建
   * 创建后 attach 到 extract
   * toast 提示成功
5. 点击 `Open Review Item` 时：

   * 切到 `review`
   * 打开对应 review item

建议把这部分抽成一个轻量 section，避免 `ExtractDetailPanel.vue` 继续膨胀过大。

---

## 十一、review 联动要求

---

### `ReviewCard.vue`

请增强 ReviewCard，使其更好地支持从 review item 跳回来源对象。

要求：

#### 当 review item.type === 'card'

显示来源信息与按钮：

* Open Card
* Reveal in Tree（可选）

#### 当 review item.type === 'extract'

显示来源信息与按钮：

* Open Extract
* Open Source Article
* Reveal in Tree（若存在 nodeId）

### 行为要求

#### Open Card

1. `card.store.setSelectedCard(cardId)`
2. 若有 nodeId，则 `tree.store.setSelectedNode(nodeId)`
3. `workspace.store.setContext('knowledge')`
4. 右侧打开 `CardDetailPanel.vue`

#### Open Extract

1. `extract.store.setSelectedExtract(extractId)`
2. 若有 tree node，则同步选中树节点
3. `workspace.store.setContext('knowledge')`

#### Open Source Article

1. `reader.store.openArticleById(sourceArticleId)`
2. `workspace.store.setContext('reading')`

---

## 十二、创建 review object 的行为要求

这是本次最核心的业务动作之一。

请实现以下规则：

### 1. 从 card 创建 review item

当用户在 `CardDetailPanel.vue` 点击 `Add to Review` 时：

* 若该 card 已有关联 review item：

  * 直接打开已有 review item
  * 不重复创建
* 若没有：

  * 构建 review item
  * 加入 review store
  * attach 到 card
  * toast 提示成功

### 2. 从 extract 创建 review item

当用户在 `ExtractDetailPanel.vue` 点击 `Create Review Item` 时：

* 若该 extract 已有关联 review item：

  * 直接打开已有 review item
* 若没有：

  * 构建 review item
  * 加入 review store
  * attach 到 extract
  * toast 提示成功

### 推荐第一阶段行为

* 创建成功后不一定立即切到 review
* 更合适的第一阶段方案是：

  * toast 提示成功
  * 保持当前上下文
  * 同时提供 `Open Review Item`

理由：

* 用户有时只是把对象加入复习，并不想立刻离开当前知识编辑流

请采用这个更适合知识工具的方案，并说明。

---

## 十三、命令面板联动要求

请补充命令系统中的这些命令：

* `add-current-card-to-review`
* `open-current-card-review`
* `add-current-extract-to-review`
* `open-current-extract-review`

### 要求

#### `add-current-card-to-review`

* 若当前焦点 card 存在，则执行加入 review

#### `open-current-card-review`

* 若当前 card 已有关联 review item，则切到 review 并打开

#### `add-current-extract-to-review`

* 若当前焦点 extract 存在，则执行加入 review

#### `open-current-extract-review`

* 若当前 extract 已有关联 review item，则切到 review 并打开

请在 `command.registry.ts` 与 `command.executor.ts` 中接好这部分联动。

---

## 十四、右键菜单联动要求

建议补充以下 context actions：

### card 节点右键菜单

* Open Card
* Add to Review
* Open Review Item
* Reveal in Tree

### extract 右键菜单

* Open Extract
* Add to Review
* Open Review Item
* Back to Source Article

要求：

* 不重复创建 review item
* 动作通过 context-menu executor 统一处理
* 不要把逻辑写在业务组件里

---

## 十五、搜索联动要求

新创建或已存在的 card/review object 关系应能被 Search 感知。

要求：

1. 搜索 card 时能命中：

   * title
   * prompt
   * answer
2. 搜索 review item 时能命中：

   * title
   * prompt
   * answer
3. 打开搜索结果时：

   * card 结果 -> knowledge + CardDetailPanel
   * review 结果 -> review + ReviewCard

当前阶段可采用：

* Search 每次搜索时从 store 重新构建 index

---

## 十六、mock 数据要求

请根据需要补充 mock 数据：

### `mock/tree.ts`

确保存在一些 `card` 类型节点，或允许通过运行时创建后可立即可用。

### `mock/review.ts`

确保至少有一两条与 card/extract 关联的 mock review item，便于测试跳转。

### `mock/index.ts`

统一导出。

### 如有需要

你可以在 card store 初始化时从已有 card 节点构建部分 card 对象，但要求结构清晰，不要隐式魔法太多。

---

## 十七、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 专业、克制
3. 轻边框
4. 小圆角
5. 更像知识工具对象详情面板，不像后台配置表单
6. card 详情与 extract 详情应保持家族一致性
7. review 关系区块应清楚但不喧宾夺主
8. 操作路径短

整体感觉应接近：

* 这是一个可编辑的知识卡片对象
* 这是一个可转入复习的学习对象
* 这是一个可追溯来源的复习对象

---

## 十八、文案与反馈要求

请统一文案风格：

* 专业
* 简洁
* 像桌面工具
* 不卖萌

建议文案：

* Add to Review
* Open Review Item
* Review item created
* Card updated
* Card already linked to a review item
* Extract added to review
* Open Card
* Open Extract
* Open Source Article

提交成功后请接入统一 toast：

* Card saved
* Review item created from card
* Review item created from extract
* Opened existing review item

---

## 十九、代码结构要求

请特别遵守：

1. 不要把 card 详情逻辑写进 tree store
2. 不要把 card 编辑逻辑写进 `RightPane.vue`
3. 不要把 review 对象构造规则写在组件里
4. `card.store.ts`、`card.mapper.ts`、`review.mapper.ts`、`extract.store.ts` 职责清晰
5. `CardDetailPanel.vue` 与 `ExtractDetailPanel.vue` 保持结构对称性
6. 所有 import 必须完整
7. 不要只写伪代码或片段
8. 代码应为“第一阶段真实可运行版本”

---

## 二十、建议的实现策略

为了降低复杂度，推荐采用以下方案：

### 推荐方案

* 新增独立 `card.store.ts`
* `RightPane.vue` 根据当前 knowledge node type 分流到：

  * `CardDetailPanel.vue`
  * `ExtractDetailPanel.vue`
* card 和 extract 各自拥有 review section
* review item 通过来源字段与 card/extract/article 建立关联
* Add to Review 若已存在 review item，则直接打开已有项

### 为什么推荐这样做

因为这套方案已经足够支撑后续扩展：

* 更高级的 card schema
* card/extract 转 cloze / QA
* review template system
* 真实后端 review object creation
* 复杂多来源对象关系

但当前阶段不会过度复杂。

---

## 二十一、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：说明你将如何分层完成本次 card + review object 工作流

例如：

* card object layer
* review object mapping layer
* detail rendering layer
* integration layer

### 第四步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text id="3a1nzn"
src/modules/card/types.ts
```

然后再给完整代码。

### 第五步：最后说明：

1. 如何接入当前项目
2. 用户如何编辑 card，并将 card/extract 转为 review object
3. 当前实现相对真实后端版 / 更复杂卡片系统的简化点有哪些

---

## 二十二、补充约束

```text id="v1xw2s"
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入大型表单库或富文本库
- 可以使用少量 Element Plus Input / Button / Divider / Card / Empty / Tag / Form 等基础控件
- 当前阶段优先保证结构清晰、可运行、风格统一
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```
