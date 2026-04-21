# 任务单 14：新建节点 / 新建卡片 / 新建摘录对象工作流第一阶段

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **“新建节点 / 新建卡片 / 新建摘录对象工作流第一阶段可运行版本”**。

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
* 阅读 / 摘录 / 知识树 / 复习 / 搜索 / 命令 / 设置 / 右键菜单 / 入箱 的基础联动
* `workspace.store.ts`
* `reader.store.ts`
* `tree.store.ts`
* `editor.store.ts`
* `extract.store.ts`
* `review.store.ts`
* `search.store.ts`
* `command.store.ts`
* `capture.store.ts`
* `mock/` 数据
* 三栏工作台与统一 UI 风格

本次目标是：**实现第一阶段的对象创建工作流**，让用户能够主动创建三类对象，并把它们接入当前系统：

1. Knowledge Node
2. Review Card
3. Manual Extract

当前阶段重点是：

* 统一创建入口
* 创建表单与对话框
* 创建后挂入知识树/复习队列/右侧详情
* 最小可运行对象流

不是复杂模板引擎、批量创建器、富文本 schema 编辑器、真实后端实体事务系统。

---

## 一、模块目标

请实现如下能力：

1. 用户能够创建新的 **知识节点**
2. 用户能够创建新的 **卡片**
3. 用户能够创建新的 **手工摘录对象**
4. 提供统一的 “Create” / “New” 入口
5. 支持从多个位置触发创建，例如：

   * 左侧导航或工具条按钮
   * 知识树右键菜单
   * 命令面板命令
   * knowledge 页面工具栏
6. 创建时用户可以选择对象类型
7. 创建后对象应正确进入系统：

   * 新节点进入知识树
   * 新卡片进入知识树，且可加入 review 队列或作为可复习对象存在
   * 新摘录对象进入知识树，并可在右侧详情编辑
8. 创建成功后：

   * 自动选中新对象
   * 切换到合理上下文
   * 右侧打开对应详情/编辑面板
   * 给出轻量 toast 提示
9. 结构要为未来接真实后端对象创建 API 留出清晰边界

---

## 二、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现复杂模板系统
3. 不要实现卡片类型设计器
4. 不要实现 cloze 自动生成
5. 不要实现复杂富文本 schema 编辑
6. 不要实现批量创建工作流
7. 不要实现复杂标签体系
8. 不要实现复杂父子节点拖拽式创建
9. 不要把所有创建逻辑塞进一个大组件
10. 不要让 tree / extract / review / editor store 职责混乱

---

## 三、本次建议实现/修改的文件

请至少实现和补全以下文件：

```text id="6jkk6g"
src/modules/object-creation/
  components/
    CreateObjectDialog.vue
    CreateObjectTypeTabs.vue
    CreateNodeForm.vue
    CreateCardForm.vue
    CreateManualExtractForm.vue
    CreateQuickEntry.vue
  stores/
    object-creation.store.ts
  types.ts
  services/
    creation.mapper.ts
    creation.validator.ts
```

并根据需要补充或修改以下文件：

```text id="d5qjko"
src/modules/knowledge-tree/stores/tree.store.ts
src/modules/extract/stores/extract.store.ts
src/modules/review/stores/review.store.ts
src/modules/workspace/stores/workspace.store.ts
src/modules/workspace/components/LeftSidebar.vue
src/modules/workspace/components/CenterPane.vue
src/modules/workspace/components/RightPane.vue
src/modules/knowledge-tree/components/TreePanel.vue
src/modules/extract/components/ExtractDetailPanel.vue
src/modules/review/components/ReviewCard.vue
src/modules/command-center/services/command.registry.ts
src/modules/command-center/services/command.executor.ts
src/modules/context-actions/services/context-menu.registry.ts
src/modules/context-actions/services/context-menu.executor.ts
src/mock/tree.ts
src/mock/review.ts
src/mock/index.ts
```

如有必要，你可以新增少量辅助文件，例如：

```text id="xa0f6d"
src/shared/constants/object-creation.ts
src/shared/utils/id.ts
```

但请保持克制，不要扩展成完整对象工厂系统。

---

## 四、整体设计要求

请把这次能力拆成 4 个清晰层次：

### 1. 创建输入层

负责：

* 不同对象类型的表单
* 基础校验
* 表单切换
* 默认字段填充

### 2. 创建状态层

负责：

* 对话框打开/关闭
* 当前创建类型
* 草稿状态
* 提交状态

### 3. 对象映射层

负责：

* 将表单输入映射为系统对象
* 构建 tree node / extract item / review card 对象
* 统一默认字段与时间戳

### 4. 系统接入层

负责：

* 将新对象写入对应 store
* 同步知识树
* 同步右侧详情/编辑
* 必要时同步 review queue

要求边界清晰，便于未来接入真实后端对象创建 API。

---

## 五、类型设计要求

请先整理和定义清晰类型。

---

### `src/modules/object-creation/types.ts`

至少包含：

```ts id="d1kncu"
export type CreateObjectType = 'node' | 'card' | 'manual-extract'

export interface CreateNodeInput {
  title: string
  content: string
  parentNodeId: string | null
}

export interface CreateCardInput {
  title: string
  prompt: string
  answer: string
  parentNodeId: string | null
  addToReview: boolean
}

export interface CreateManualExtractInput {
  title: string
  quote: string
  note: string
  parentNodeId: string | null
  sourceArticleId?: string | null
  sourceArticleTitle?: string | null
}

export interface CreateObjectDraftState {
  type: CreateObjectType
  node: CreateNodeInput
  card: CreateCardInput
  manualExtract: CreateManualExtractInput
}

export interface CreateObjectResult {
  objectType: CreateObjectType
  id: string
  nodeId?: string | null
  extractId?: string | null
  reviewItemId?: string | null
}
```

你可以在此基础上微调，但必须满足：

1. 三类创建输入结构清晰分开
2. 能表达 parentNodeId
3. manual extract 能关联来源文章
4. card 能选择是否进入 review
5. 类型不要过度复杂

---

## 六、校验要求

请新增：

### `creation.validator.ts`

至少实现这些函数：

1. `validateNodeInput(...)`
2. `validateCardInput(...)`
3. `validateManualExtractInput(...)`

### 校验规则建议

#### node

* title 可以为空，但为空时自动生成默认标题
* content 可为空

#### card

* prompt 必填
* answer 必填
* title 可为空，若为空则根据 prompt 自动生成
* parentNodeId 可为空

#### manual extract

* quote 必填
* title 可为空，若为空则根据 quote 自动生成
* note 可为空
* parentNodeId 可为空

要求：

* 当前阶段只做基础校验
* 校验逻辑不要写在表单组件里
* 错误文案简洁、专业

---

## 七、映射要求

请新增：

### `creation.mapper.ts`

至少实现以下函数：

1. `buildTreeNodeFromCreateNodeInput(...)`
2. `buildTreeNodeFromCreateCardInput(...)`
3. `buildExtractFromManualExtractInput(...)`
4. `buildTreeNodeFromManualExtract(...)`
5. `buildReviewItemFromCreateCardInput(...)`

要求：

* 不要把对象构造规则写在表单组件里
* 映射规则清晰、稳定
* 便于未来接真实后端 DTO
* 所有对象应生成合理的：

  * id
  * title
  * content/quote/prompt/answer
  * createdAt / updatedAt
  * type
  * parentNodeId

### 推荐规则

#### 新建节点

* tree node type = `topic`
* 若 title 为空 -> `Untitled Node`
* content 可作为后续右侧编辑初始内容

#### 新建卡片

* tree node type = `card`
* title 为空时从 prompt 截断生成
* card 对象同时可映射为 review item
* 若 `addToReview = true`，则加入 review queue

#### 手工摘录对象

* extract item 单独存入 extract store
* 同时映射为 tree node type = `extract`
* title 为空时从 quote 截断生成

---

## 八、store 设计要求

---

### `object-creation.store.ts`

请实现对象创建模块专用 store。

至少包含这些状态：

```ts id="de8ah8"
interface ObjectCreationState {
  isDialogOpen: boolean
  activeType: CreateObjectType
  submitting: boolean
  draft: CreateObjectDraftState
}
```

至少包含这些 getter：

1. `activeDraft`
2. `dialogTitle`
3. `canSubmit`

至少包含这些 action：

1. `openDialog(type?: CreateObjectType)`
2. `closeDialog()`
3. `setActiveType(type: CreateObjectType)`
4. `updateNodeDraft(payload: Partial<CreateNodeInput>)`
5. `updateCardDraft(payload: Partial<CreateCardInput>)`
6. `updateManualExtractDraft(payload: Partial<CreateManualExtractInput>)`
7. `resetDraft(type?: CreateObjectType)`
8. `submit()`
9. `submitNode()`
10. `submitCard()`
11. `submitManualExtract()`

### 行为要求

* 打开 dialog 时可指定默认类型
* 提交成功后：

  * 将对象写入对应 store
  * 切到合适上下文
  * 自动选中新对象
  * 打开右侧详情/编辑区
  * 关闭 dialog
  * toast 提示成功

### 推荐第一阶段行为

#### submitNode

* 写入 tree store
* 自动展开父节点
* 选中新节点
* 切到 `knowledge`

#### submitCard

* 写入 tree store
* 若 `addToReview = true`，同时写入 review store
* 选中新卡片节点
* 切到 `knowledge`
* 右侧显示卡片详情或知识节点详情占位

#### submitManualExtract

* 写入 extract store
* 同时写入 tree store
* 同步 extractId 与 treeNodeId
* 选中新摘录
* 切到 `knowledge`
* 右侧打开摘录详情

要求：

* store 只负责创建流程，不要承担全部业务模块逻辑

---

## 九、各业务 store 联动要求

---

### 1. `tree.store.ts`

请补充或确保具备：

* `addNode(...)`
* `setSelectedNode(id: string | null)`
* `expandNode(id: string)`
* `updateNodeTitle(...)`
* 能支持新增 `topic` / `card` / `extract` 节点

要求：

* 新对象进入树后可立即显示
* 父节点自动展开
* 新节点自动选中

---

### 2. `extract.store.ts`

请补充或确保具备：

* `addExtract(...)`
* `setSelectedExtract(id: string | null)`
* `getExtractById(...)`

要求：

* manual extract 创建后能立即进入 extract store
* 与 tree node 建立可追溯关系

---

### 3. `review.store.ts`

请补充或确保具备：

* `addReviewItem(...)`
* `setSelectedItem(id: string | null)`

要求：

* 当新 card 选择 `addToReview = true` 时加入 review queue
* 后续 review 页面可以看到

---

### 4. `workspace.store.ts`

请确保可用于：

* `setContext('knowledge')`
* 必要时同步当前焦点对象 id

当前阶段不要求做复杂对象焦点总线，但要保持最小联动可用。

---

## 十、组件要求

---

### 1. `CreateObjectDialog.vue`

这是本次最核心的 UI 组件。

要求：

1. 使用 Element Plus `Dialog` 可接受
2. 从 `object-creation.store.ts` 读取状态
3. 包含：

   * 标题
   * `CreateObjectTypeTabs.vue`
   * 当前对应表单组件
   * footer 操作区
4. 支持：

   * Cancel
   * Create
5. 提交中显示 loading
6. 切换 tab 时保留各自 draft
7. 风格应像桌面工具对象创建面板，不像后台表单弹窗

---

### 2. `CreateObjectTypeTabs.vue`

职责：

* 切换创建类型

至少包括：

* Node
* Card
* Manual Extract

要求：

* 紧凑
* 清晰
* 切换后更新 `activeType`

---

### 3. `CreateNodeForm.vue`

职责：

* 新建知识节点表单

至少字段：

* title
* content
* parentNodeId

要求：

* parentNodeId 使用简单下拉框
* 下拉选项来自知识树可选节点
* 内容区可为轻量多行输入
* 保持紧凑，不做复杂文档编辑器

---

### 4. `CreateCardForm.vue`

职责：

* 新建卡片表单

至少字段：

* title
* prompt
* answer
* parentNodeId
* addToReview

要求：

* prompt / answer 是重点输入区
* `addToReview` 使用 switch/checkbox
* 保持“知识卡片创建”感，不是后台配置页

---

### 5. `CreateManualExtractForm.vue`

职责：

* 新建手工摘录对象表单

至少字段：

* title
* quote
* note
* parentNodeId
* 可选 sourceArticleId / sourceArticleTitle 显示或隐藏处理

要求：

* quote 为重点输入区
* note 为补充说明
* parentNodeId 可选
* 第一阶段可不做复杂 article selector，但要留结构空间

### 推荐第一阶段简化

* 手工摘录默认不强制绑定 source article
* 若当前 reading context 有打开文章，可自动预填 sourceArticleId/sourceArticleTitle
* 这是很合理的第一阶段增强点

---

### 6. `CreateQuickEntry.vue`

职责：

* 一个轻量创建入口组件
* 可用于：

  * 左侧导航区按钮
  * knowledge/tree toolbar 按钮
  * 未来快捷入口

要求：

* 点击后打开 `CreateObjectDialog`
* 可接受默认 type
* 紧凑、专业

---

## 十一、知识树 / 右侧详情集成要求

---

### 1. `TreePanel.vue`

要求：

* 增加一个 `New` / `Create` 按钮
* 可打开 `CreateObjectDialog`
* 默认 type 可设为 `node`

---

### 2. `RightPane.vue`

要求：

* 创建对象成功后，右侧能显示合理详情：

#### 新建 node

* 可显示 knowledge detail 占位或右侧基础详情
* 若你已有 node 详情面板，可接入

#### 新建 card

* 可先显示 node/card detail 占位
* 当前阶段不要求实现完整卡片编辑器，但至少要能打开对象详情

#### 新建 manual extract

* 应优先显示 `ExtractDetailPanel.vue`

要求：

* 不要把对象创建逻辑写在 RightPane 里
* RightPane 只负责根据当前 context 和对象类型选择渲染

---

## 十二、命令面板联动要求

请补充命令系统中的对象创建命令。

建议新增以下命令：

* `create-node`
* `create-card`
* `create-manual-extract`

### 要求

#### `create-node`

* 打开 dialog
* 默认 type = node
* 切到 knowledge 或保持当前 knowledge 相关上下文即可

#### `create-card`

* 打开 dialog
* 默认 type = card

#### `create-manual-extract`

* 打开 dialog
* 默认 type = manual-extract
* 若当前 reading 中有文章，可预填 source article 信息

请在 `command.registry.ts` 与 `command.executor.ts` 中接好联动。

---

## 十三、右键菜单联动要求

建议补充以下 context actions：

### tree-node 上下文菜单

* New Child Node
* New Card Under Node
* New Manual Extract Under Node

### extract 上下文菜单（可选）

* Create Related Card（当前阶段可先不做或 disabled）

### article 上下文菜单（可选）

* New Manual Extract
* New Card From Article（可占位）

当前阶段最重要的是把 **tree-node 上的创建动作** 接起来。

### 右键创建行为要求

若从某个 tree node 上触发创建：

* `parentNodeId` 自动预填为当前节点 id

---

## 十四、搜索联动要求

新创建的对象必须能被 Search 模块搜到。

要求：

1. 新建 node 后能被搜索到
2. 新建 card 后能被搜索到
3. 新建 manual extract 后能被搜索到

当前阶段可采用简单方案：

* Search 每次运行时重新从各 store 构建 index
* 不要求复杂增量索引

---

## 十五、Review 联动要求

当用户创建 card 并勾选 `addToReview = true` 时：

1. card 进入知识树
2. 同时创建 review item
3. review 页面可看到该项
4. 后续可从 review 跳回 card/知识树（若当前系统支持）

### 当前阶段的简化建议

* card 与 review item 可通过 `nodeId` 或独立 id 建立关联
* 不要求完整 card schema，只要 prompt / answer 能进入 review 即可

---

## 十六、默认值与预填要求

请尽量提供合理预填逻辑。

### 推荐规则

#### 从 tree node 右键创建

* 自动带上 `parentNodeId`

#### 从 reading context 创建 manual extract

* 若当前 article 已打开：

  * 自动填 `sourceArticleId`
  * 自动填 `sourceArticleTitle`

#### 从命令面板直接创建

* 若无上下文则 `parentNodeId = null`

#### 标题自动生成

* node：`Untitled Node`
* card：从 prompt 前 20~40 字生成
* manual extract：从 quote 前 20~40 字生成

---

## 十七、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 专业、克制
3. 轻边框
4. 小圆角
5. 更像桌面知识工具的对象创建面板，不像后台配置弹窗
6. 用户应感觉是在“快速创建知识对象”
7. 表单层级清楚
8. 提交路径短

整体感觉应接近：

* 快速建对象
* 低心智负担
* 立即进入工作流
* 不花哨，但很实用

---

## 十八、文案与反馈要求

请统一文案风格：

* 专业
* 简洁
* 像桌面工具
* 不卖萌

建议文案：

* Create Node
* Create Card
* Create Manual Extract
* New node created
* Card added to review
* Manual extract created
* Please enter a prompt
* Please enter a quote

提交成功后请接入统一 toast：

* Node created
* Card created
* Card added to review
* Manual extract created

---

## 十九、代码结构要求

请特别遵守：

1. 不要把三种创建表单写进一个超大组件
2. 不要把对象映射规则写在表单组件里
3. `creation.mapper.ts`、`creation.validator.ts`、`object-creation.store.ts` 分层清晰
4. `CreateObjectDialog.vue` 只负责组合 UI
5. tree / extract / review store 各自只负责自身对象状态
6. 所有 import 必须完整
7. 不要只写伪代码或片段
8. 代码应为“第一阶段真实可运行版本”

---

## 二十、建议的实现策略

为了降低复杂度，推荐采用以下方案：

### 推荐方案

* 全局挂一个 `CreateObjectDialog.vue`
* object-creation store 管理 dialog 与 active type
* 三种 form 分开
* submit 时 mapper 生成对应对象
* 根据对象类型写入不同 store
* 成功后切到 `knowledge`
* 自动选中新对象
* 通过 toast 提示成功

### 为什么推荐这样做

因为这套方案已经足够支撑后续扩展：

* 卡片模板
* 手工摘录转卡片
* 节点类型扩展
* 更复杂的 parent 选择器
* 多对象创建 wizard
* 真实后端事务创建

但当前阶段不会过度复杂。

---

## 二十一、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：说明你将如何分层完成本次对象创建系统

例如：

* form layer
* validation layer
* mapping layer
* submission layer
* integration layer

### 第四步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text id="g3c7b8"
src/modules/object-creation/types.ts
```

然后再给完整代码。

### 第五步：最后说明：

1. 如何接入当前项目
2. 用户如何通过 node / card / manual extract 三种方式创建对象
3. 当前实现相对真实对象创建系统的简化点有哪些

---

## 二十二、补充约束

```text id="mmub0e"
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入大型表单库
- 可以使用少量 Element Plus Dialog / Input / Select / Switch / Form / Button 等基础控件
- 当前阶段优先保证结构清晰、可运行、风格统一
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 二十三、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。
