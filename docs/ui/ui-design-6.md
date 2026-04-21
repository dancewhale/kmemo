# 任务单 6：摘录节点详情编辑 + 阅读/知识区联动

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **“摘录节点详情编辑 + 阅读/知识区联动”的第一阶段可运行版本**。

这个项目是一个基于 **Wails + Go 后端** 的桌面知识工具，产品定位参考 **SuperMemo 风格的渐进阅读 / 知识工作台**。
当前项目已经具备：

* `workspace` 模块
* `reader` 模块
* `knowledge-tree` 模块
* `editor` 模块壳
* `extract` 模块基础创建流程
* 已支持在阅读文章中选中文本并创建摘录
* 创建摘录后可将其挂入知识树
* `tree.store.ts`
* `reader.store.ts`
* `editor.store.ts`
* `extract.store.ts`
* `mock/` 数据
* 右侧可打开文章编辑区

本次目标是：**让摘录成为一个可以被查看、编辑、追踪来源、并与原文联动的知识对象**。

---

## 一、模块目标

请实现如下完整交互能力：

1. 用户在阅读区创建摘录后，摘录进入知识树
2. 用户点击知识树中的摘录节点
3. 右侧打开该摘录的详情/编辑界面
4. 用户可编辑摘录的：

   * 标题
   * quote（可只读或可轻度编辑，请说明选择）
   * note
5. 用户可看到摘录的来源文章信息
6. 用户可从摘录跳回原文章
7. 用户在阅读文章时，也能看到当前文章关联的摘录列表
8. 用户可点击某个摘录，重新聚焦到该摘录详情
9. 阅读区与知识树、摘录详情之间形成联动

当前阶段重点是 **知识对象联动和右侧编辑承载**，不是复杂富文本能力。

---

## 二、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现复杂富文本摘录编辑
3. 不要实现真正的 anchor 定位高亮恢复
4. 不要实现复杂 backlinks 图谱
5. 不要实现复习算法
6. 不要实现 QA / cloze 自动生成
7. 不要实现复杂多标签筛选器
8. 不要实现复杂 tab 工作区系统
9. 不要把所有逻辑塞进 `RightPane.vue`
10. 不要让 store 之间职责混乱

---

## 三、本次建议实现/修改的文件

请至少实现和补全以下文件：

```text id="udm86a"
src/modules/extract/
  components/
    ExtractDetailPanel.vue
    ExtractMetaSection.vue
    ExtractSourceSection.vue
    ArticleExtractList.vue
  stores/
    extract.store.ts
  services/
    extract.mapper.ts
  types.ts
```

并根据需要补充或修改以下文件：

```text id="uv27yb"
src/modules/editor/components/EditorShell.vue
src/modules/workspace/components/RightPane.vue
src/modules/workspace/components/CenterPane.vue
src/modules/knowledge-tree/stores/tree.store.ts
src/modules/reader/stores/reader.store.ts
src/modules/editor/stores/editor.store.ts
src/mock/articles.ts
src/mock/tree.ts
src/mock/index.ts
```

如果有必要，你可以新增少量辅助文件，例如：

```text id="e0h4z7"
src/shared/constants/entity.ts
src/shared/utils/entity.ts
```

但请保持克制，不要扩展过大。

---

## 四、整体设计要求

请把这次能力拆成 4 个清晰层次：

### 1. 摘录对象层

负责：

* 摘录数据结构
* 摘录列表
* 摘录更新
* 摘录与来源文章的关联

### 2. 知识树选择层

负责：

* 树节点选中
* 判断当前选中节点是否为摘录节点
* 为右侧提供当前聚焦对象

### 3. 右侧详情编辑层

负责：

* 根据当前焦点对象，渲染文章编辑器或摘录详情
* 摘录标题/note 编辑
* 来源信息展示
* 跳回原文入口

### 4. 阅读区联动层

负责：

* 在阅读文章详情中展示该文章关联的摘录列表
* 点击摘录后能聚焦摘录详情
* 形成阅读 -> 摘录 -> 知识树 -> 摘录详情 -> 原文 的来回路径

要求职责清晰，后续便于接真实后端。

---

## 五、类型设计要求

请先整理和补充清晰类型。

---

### `src/modules/extract/types.ts`

至少包含：

```ts id="p4i76c"
export interface TextSelectionRange {
  from: number
  to: number
  text: string
}

export interface ExtractItem {
  id: string
  sourceArticleId: string
  sourceArticleTitle: string
  title: string
  quote: string
  note: string
  parentNodeId: string | null
  treeNodeId?: string | null
  createdAt: string
  updatedAt: string
}

export interface ExtractUpdatePayload {
  title?: string
  note?: string
  quote?: string
}
```

你可以在此基础上微调，但必须满足：

1. 摘录对象可以单独存在并被编辑
2. 摘录对象能关联原文章
3. 摘录对象可关联知识树节点 id
4. 为未来接真实后端 DTO 留空间

---

### 对 `knowledge-tree/types.ts` 的补充建议

请保证树节点至少具备：

```ts id="26b2ep"
type KnowledgeNodeType = 'topic' | 'article' | 'extract' | 'card'
```

且摘录节点能通过 `type === 'extract'` 被识别。

---

### 对 `editor/types.ts` 的要求

请保留文章编辑对象与摘录详情对象的边界。
不要强行让 `EditorDocument` 承担摘录详情全部职责。

你可以允许右侧根据当前 context / focus object 类型分流：

* article -> `EditorShell`
* extract -> `ExtractDetailPanel`

这是当前阶段最合理的做法。

---

## 六、store 设计要求

---

### 1. `extract.store.ts`

请在已有基础上补充，使其能真正承担摘录对象管理。

至少包含这些状态：

```ts id="i6osvz"
interface ExtractState {
  items: ExtractItem[]
  pending: PendingExtractPayload | null
  isCreateDialogOpen: boolean
  creating: boolean
  selectedExtractId: string | null
  saving: boolean
}
```

至少包含这些 getter：

1. `selectedExtract`
2. `extractCount`
3. `getExtractsByArticleId(articleId: string)`
4. `getExtractById(id: string)`

至少包含这些 action：

1. `setSelectedExtract(id: string | null)`
2. `updateExtract(id: string, payload: ExtractUpdatePayload)`
3. `openExtract(id: string)`
4. `createExtract()`（保留原有能力）
5. `getExtractsForArticle(articleId: string)`（可返回数组）

行为要求：

* 创建摘录后，要能保存 `treeNodeId` 或能建立可追溯关系
* `updateExtract` 能更新 title / note / quote（如果你允许 quote 可编辑）
* 更新摘录标题时，如该摘录对应知识树节点存在，应同步更新树节点标题
* `selectedExtractId` 由模块内部管理，但允许与 workspace/tree 进行轻量同步

---

### 2. `tree.store.ts`

请补充或完善以下能力：

1. 能根据 `nodeId` 找到节点
2. 能判断节点类型
3. 能提供当前选中节点
4. 提供一个用于更新节点标题的 action，例如：

```ts id="2k6zuj"
updateNodeTitle(nodeId: string, title: string)
```

5. 当选中摘录节点时，能让其他模块知道“当前焦点是 extract 节点”

当前阶段可允许：

* `RightPane.vue` 中根据当前选中节点类型做判断
* 或通过少量 getter 暴露

但不要让 tree store 直接接管摘录详情状态。

---

### 3. `reader.store.ts`

请补充以下能力：

1. 当前选中文章保持可读
2. 提供根据 id 获取 article 的 getter
3. 支持从外部“打开某篇文章”，例如：

```ts id="4cf7j8"
openArticleById(id: string)
```

用于从摘录详情跳回原文。

---

### 4. `editor.store.ts`

当前阶段 editor store 继续主要服务于文章编辑器。

请确保：

* 能通过 `openDocument(...)` 打开文章
* 不要把摘录详情编辑塞进 editor store
* 文章编辑状态与摘录详情状态分离

---

## 七、组件要求

---

### 1. `ExtractDetailPanel.vue`

这是本次最核心的新组件。

要求：

1. 用于在右侧展示当前选中摘录的详情
2. 从 `extract.store.ts` 读取当前选中摘录
3. 至少展示以下区域：

   * 标题
   * quote
   * note
   * 来源文章信息
   * 创建/更新时间
4. 标题可编辑
5. note 可编辑
6. quote：

   * 可以只读展示
   * 或允许轻量编辑
   * 请你选择一种更适合第一阶段的方案并说明
7. 有保存反馈
8. 无选中摘录时显示 `AppEmpty`

风格要求：

* 它应该像一个知识对象详情面板
* 不是后台表单
* 也不是普通阅读页
* 更像“摘录对象编辑器”

---

### 2. `ExtractMetaSection.vue`

职责：

* 展示摘录基本元信息，例如：

  * id
  * createdAt
  * updatedAt
  * parentNodeId
  * sourceArticleId

要求：

* 轻量
* 可复用
* 信息密度高
* 风格专业克制

---

### 3. `ExtractSourceSection.vue`

职责：

* 展示摘录来源信息
* 至少包含：

  * 来源文章标题
  * article id
  * 一个“Open Source Article”按钮

行为要求：

* 点击按钮后：

  1. 调用 `reader.store.openArticleById(sourceArticleId)`
  2. 切换 workspace 到 `reading`
  3. 右侧打开对应文章编辑区

这一步非常重要，它实现 **摘录 -> 原文** 的回跳链路。

---

### 4. `ArticleExtractList.vue`

这是阅读文章右侧/底部显示的“当前文章关联摘录列表”。

要求：

1. 在当前打开文章时显示该文章关联的摘录
2. 从 `extract.store.ts` 获取 `getExtractsByArticleId(articleId)`
3. 每个列表项至少显示：

   * 摘录标题
   * quote 片段摘要
   * 更新时间
4. 点击某个摘录项时：

   * 设置 `extract.store.selectedExtractId`
   * 同步选中对应知识树节点（如果有）
   * workspace 切换到 `knowledge` 或保留在当前右侧直接切换展示，请你选择一种更合理且更简单的方案，并说明理由

建议当前阶段的更合理方案：

* 点击后切换到 `knowledge`
* 高亮该摘录节点
* 右侧打开摘录详情

这样知识树和摘录对象关系更清晰。

---

## 八、EditorShell 与右侧集成要求

---

### `EditorShell.vue`

请在当前文章编辑器右侧结构中补充一个“相关摘录”区域。

要求：

1. 当当前文档为 article 时，在编辑器区域下方或侧边增加 `ArticleExtractList.vue`
2. 该区域只展示与当前 article 关联的摘录
3. 不要让 `EditorShell.vue` 过度臃肿
4. 可以通过组件拆分保持清晰结构

建议布局：

* 上方：文章标题和元信息
* 中部：工具栏
* 主体：文章内容编辑区
* 下方：Related Extracts

---

### `RightPane.vue`

请根据当前焦点对象做清晰分流。

要求：

#### 当 `context === 'reading'`

* 右侧优先显示文章 `EditorShell.vue`

#### 当 `context === 'knowledge'`

* 如果当前选中知识树节点类型是 `extract`

  * 显示 `ExtractDetailPanel.vue`
* 否则

  * 继续显示节点详情占位或已有的 knowledge detail 视图

#### 当 `context === 'inbox'`

* 可继续沿用阅读式详情或轻量预览
* 请说明你的处理

#### 当 `context === 'review'`

* 不接这次的新组件

重点：

* `RightPane.vue` 只负责“根据焦点对象选择渲染哪个详情组件”
* 不要把摘录编辑逻辑写在 `RightPane.vue`

---

## 九、联动规则要求

请实现下面几条关键联动规则。

---

### 规则 1：树 -> 摘录详情

当用户点击知识树中的摘录节点时：

1. tree store 选中该节点
2. 找到对应 extract
3. extract store 设置 `selectedExtractId`
4. `RightPane.vue` 渲染 `ExtractDetailPanel.vue`

---

### 规则 2：摘录标题更新 -> 树节点标题同步

当用户在 `ExtractDetailPanel.vue` 中修改摘录标题并保存后：

1. 更新 `extract.store.items`
2. 若该摘录存在对应树节点，则同步更新 `tree.store.updateNodeTitle(...)`

---

### 规则 3：摘录 -> 原文章

当用户点击 `Open Source Article` 时：

1. `reader.store.openArticleById(sourceArticleId)`
2. `workspace.store.setContext('reading')`
3. `editor.store.openDocument(...)`
4. 右侧展示文章编辑器

---

### 规则 4：文章 -> 关联摘录

当用户打开某篇文章时：

1. `EditorShell.vue` 显示 `ArticleExtractList.vue`
2. 该列表展示所有关联摘录

---

### 规则 5：文章中的关联摘录 -> 摘录详情

当用户在 `ArticleExtractList.vue` 中点击某个摘录时：

1. `extract.store.setSelectedExtract(extract.id)`
2. 若有对应 tree node，则同步选中树节点
3. 切换到 `knowledge`
4. 右侧打开摘录详情

---

## 十、mock 数据要求

请根据需要补充 mock 数据：

### `mock/articles.ts`

确保文章内容和文章 id 完整，便于摘录回跳。

### `mock/tree.ts`

确保存在 extract 节点，或允许通过运行时动态创建后保持可访问。

### `mock/index.ts`

统一导出。

如你需要，也可以在 `mock` 中增加一些预置 extract 数据，用于让页面一打开就能测试联动。

---

## 十一、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 专业、克制
3. 轻边框
4. 小圆角
5. 编辑与详情层级清晰
6. 更像知识工具，不像后台管理系统
7. “Related Extracts” 列表应紧凑、可快速浏览
8. 摘录详情应更像知识对象面板，而不是大表单页

整体感觉应接近：

* 从阅读中提炼知识
* 知识对象可回看、可编辑、可追溯来源
* 工作流连续
* 操作切换自然

---

## 十二、代码结构要求

请特别遵守：

1. 不要把摘录详情逻辑全塞进 `RightPane.vue`
2. 不要让 `extract.store.ts` 直接承担树所有逻辑
3. 不要让 `tree.store.ts` 直接承担摘录详情编辑
4. 不要把文章编辑与摘录编辑混成一个 store
5. 映射逻辑尽量放在 `extract.mapper.ts`
6. 所有 import 必须完整
7. 不要只写伪代码或片段
8. 代码应为“第一阶段真实可运行版本”

---

## 十三、建议的实现策略

为了降低复杂度，推荐采用以下方案：

### 推荐方案

* 文章继续由 `EditorShell.vue` 承载
* 摘录详情使用独立的 `ExtractDetailPanel.vue`
* 在 `RightPane.vue` 中通过“当前 context + 当前聚焦对象类型”决定渲染哪一个
* 摘录详情中：

  * `quote` 默认只读展示
  * `title` 和 `note` 可编辑

### 为什么推荐 quote 只读

因为第一阶段里，quote 本质上是原文摘录证据。
允许随意修改 quote 会让“摘录对象”和“原文片段”边界变模糊。
因此当前阶段最合理的方案是：

* `title`：可编辑
* `note`：可编辑
* `quote`：只读显示

后续若要支持“编辑后另存为 note”再扩展。

---

## 十四、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text id="rd9w8z"
src/modules/extract/components/ExtractDetailPanel.vue
```

然后再给完整代码。

### 第四步：最后说明：

1. 如何接入当前项目
2. 用户如何在界面中完成“文章 -> 摘录 -> 树 -> 摘录详情 -> 原文”的联动
3. 当前实现相对真实后端版 / Tiptap 版的简化点有哪些

---

## 十五、补充约束

```text id="e6m8v4"
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入复杂第三方状态管理以外的新库
- 可以使用少量 Element Plus Input / Button / Card / Divider 等基础控件
- 当前阶段优先保证结构清晰、模块边界清晰、可运行
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 十六、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。
