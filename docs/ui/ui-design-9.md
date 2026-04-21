# 任务单 9：搜索模块（Search）第一阶段

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **“搜索模块（Search）第一阶段可运行版本”**。

这个项目是一个基于 **Wails + Go 后端** 的桌面知识工具，产品定位参考 **SuperMemo 风格的渐进阅读 / 知识工作台**。
当前项目已经具备：

* `workspace` 模块
* `reader` 模块
* `knowledge-tree` 模块
* `editor` 模块
* `extract` 模块
* `review` 模块
* `command-center` 模块
* 阅读 / 摘录 / 知识树 / 复习队列 / 命令面板基本联动
* `workspace.store.ts`
* `reader.store.ts`
* `tree.store.ts`
* `editor.store.ts`
* `extract.store.ts`
* `review.store.ts`
* `command.store.ts`
* `mock/` 数据
* 主工作台布局

本次目标是：**实现第一阶段 Search 模块**，让用户能够统一检索文章、摘录、知识节点、复习项，并从搜索结果直接跳转到对象详情。

当前阶段重点是 **搜索对象整合、搜索结果展示、结果跳转联动**，不是全文搜索引擎和复杂 ranking 系统。

---

## 一、模块目标

请实现如下能力：

1. 用户进入 `search` 工作区
2. 中间区域显示统一搜索面板
3. 用户输入关键词后，可以搜索：

   * articles
   * extracts
   * knowledge nodes
   * review items
4. 支持按对象类型筛选结果
5. 结果列表展示：

   * 类型
   * 标题
   * 摘要片段
   * 来源信息（如适用）
6. 用户点击某个结果后：

   * 工作台切换到对应上下文
   * 左/中/右区域同步到目标对象
7. 支持空查询时展示合理默认状态
8. 支持无结果空状态
9. Search 模块结构清晰，为未来接后端全文索引预留边界

当前阶段 **不要求真实全文搜索**，但要做到真实可用。

---

## 二、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现全文索引库
3. 不要实现复杂高亮定位算法
4. 不要实现复杂分词系统
5. 不要实现复杂排序策略配置
6. 不要实现高级筛选器面板
7. 不要实现搜索历史系统
8. 不要实现保存搜索
9. 不要把所有搜索逻辑塞进一个组件
10. 不要让 search store 和其他业务 store 职责混乱

---

## 三、本次建议实现/修改的文件

请至少实现和补全以下文件：

```text id="5z7q1m"
src/modules/search/
  components/
    GlobalSearch.vue
    SearchToolbar.vue
    SearchResultList.vue
    SearchResultItem.vue
    SearchEmptyState.vue
    SearchFilterTabs.vue
  stores/
    search.store.ts
  types.ts
  services/
    search.index.ts
    search.mapper.ts
    search.matcher.ts
```

并根据需要补充或修改以下文件：

```text id="h1p4lx"
src/modules/workspace/components/CenterPane.vue
src/modules/workspace/components/RightPane.vue
src/pages/SearchPage.vue
src/modules/command-center/services/command.executor.ts
src/modules/reader/stores/reader.store.ts
src/modules/extract/stores/extract.store.ts
src/modules/knowledge-tree/stores/tree.store.ts
src/modules/review/stores/review.store.ts
src/modules/workspace/stores/workspace.store.ts
src/mock/articles.ts
src/mock/tree.ts
src/mock/review.ts
src/mock/index.ts
```

如有必要，你可以新增少量辅助文件，例如：

```text id="xh82ui"
src/shared/constants/search.ts
src/shared/utils/search.ts
```

但请保持克制，不要扩展成完整搜索平台。

---

## 四、整体设计要求

请把这次能力拆成 4 个清晰层次：

### 1. 搜索对象标准化层

负责：

* 将 article / extract / node / review item 统一映射为可搜索对象
* 提供统一的搜索结果结构

### 2. 搜索匹配层

负责：

* 根据 query 做简单匹配
* 返回结果与基础 score
* 支持类型过滤

### 3. 搜索状态层

负责：

* query
* active type filter
* 结果列表
* 当前选中结果
* 空状态

### 4. 搜索跳转执行层

负责：

* 点击结果后跳到对应模块
* 切换 workspace context
* 同步 reader/tree/extract/review/editor 状态

要求边界清晰，便于未来替换为真实后端搜索服务。

---

## 五、类型设计要求

请先整理和定义清晰类型。

---

### `src/modules/search/types.ts`

至少包含：

```ts id="f4mjlwm"
export type SearchEntityType = 'article' | 'extract' | 'node' | 'review'
export type SearchFilterType = 'all' | SearchEntityType

export interface SearchDocument {
  id: string
  entityType: SearchEntityType
  title: string
  content: string
  subtitle?: string
  sourceId?: string | null
  sourceTitle?: string | null
  nodeId?: string | null
  articleId?: string | null
  extractId?: string | null
  reviewId?: string | null
  updatedAt?: string | null
  keywords?: string[]
}

export interface SearchResultItem {
  id: string
  entityType: SearchEntityType
  title: string
  snippet: string
  subtitle?: string
  score: number
  sourceId?: string | null
  sourceTitle?: string | null
  nodeId?: string | null
  articleId?: string | null
  extractId?: string | null
  reviewId?: string | null
  updatedAt?: string | null
}
```

你可以在此基础上微调，但必须满足：

1. 可搜索对象和搜索结果对象分离
2. 支持多实体统一搜索
3. 能支持未来替换为后端返回结果
4. 类型不要过度复杂

---

## 六、搜索对象标准化要求

请新增：

### `search.mapper.ts`

用于把现有各模块数据映射成统一搜索对象。

至少实现类似这些函数：

1. `mapArticlesToSearchDocuments(...)`
2. `mapExtractsToSearchDocuments(...)`
3. `mapNodesToSearchDocuments(...)`
4. `mapReviewsToSearchDocuments(...)`

要求：

* 不要把映射逻辑散落在组件里
* 统一将不同模块对象映射为 `SearchDocument`
* snippet 不在 mapper 中生成，snippet 由 matcher 或 result builder 生成更合适

---

### `search.index.ts`

用于整合当前项目里的所有可搜索对象。

至少实现：

```ts id="w74j0k"
buildSearchIndex(): SearchDocument[]
```

要求：

* 从各 store 或 mock 数据中读取基础数据
* 合并成统一数组
* 当前阶段可以简单实现，不需要真实倒排索引
* 结构上为未来接后端搜索服务留边界

建议：

* index builder 可以从 store 读取已有状态
* 或在 search store initialize 时调用

---

## 七、搜索匹配要求

请新增：

### `search.matcher.ts`

至少实现类似：

```ts id="e37r9c"
matchSearchDocuments(
  docs: SearchDocument[],
  query: string,
  filter: SearchFilterType
): SearchResultItem[]
```

要求：

1. 当前阶段只做简单匹配即可
2. 至少支持在以下字段中搜索：

   * title
   * content
   * subtitle
   * keywords
3. 标题命中分数高于正文命中
4. query 为空时：

   * 可以返回空列表
   * 或返回最近更新项
   * 请选择一种更适合第一阶段的方案并说明
5. 支持按类型过滤：

   * all
   * article
   * extract
   * node
   * review
6. snippet 需要从命中内容中生成简短摘要
7. 不要引入复杂 fuzzy 库

### 推荐第一阶段方案

* query 为空时返回最近更新的若干条对象
* 这样 Search 页面一打开就不空，更像工作台搜索中心

---

## 八、store 设计要求

---

### `search.store.ts`

请实现 Search 模块专用 store。

至少包含这些状态：

```ts id="kpjlwm"
interface SearchState {
  query: string
  filter: SearchFilterType
  results: SearchResultItem[]
  selectedResultId: string | null
  loading: boolean
}
```

至少包含这些 getter：

1. `hasQuery`
2. `hasResults`
3. `selectedResult`
4. `resultCount`
5. `resultsByFilter`（如果你认为有必要）

至少包含这些 action：

1. `initialize()`
2. `setQuery(query: string)`
3. `setFilter(filter: SearchFilterType)`
4. `runSearch()`
5. `setSelectedResult(id: string | null)`
6. `openResult(result: SearchResultItem)`
7. `clearSearch()`

### 行为要求

* 初始化时准备搜索索引
* query 或 filter 改变时，可自动重新搜索
* 默认选中第一条结果（若存在）
* `openResult` 是本次关键 action 之一，负责结果跳转联动
* search store 只管理搜索状态，不要把其他模块大量状态搬进来

---

## 九、搜索结果跳转要求

这是本次最关键的部分之一。

### `search.store.ts` 的 `openResult(result)` 或单独执行器，需要支持如下跳转规则：

---

### 1. 打开 article 结果

当结果类型为 `article` 时：

1. `reader.store.openArticleById(articleId 或 id)`
2. `workspace.store.setContext('reading')`
3. 右侧打开文章编辑区
4. 必要时同步 `editor.store.openDocument(...)`

---

### 2. 打开 extract 结果

当结果类型为 `extract` 时：

1. `extract.store.setSelectedExtract(extractId 或 id)`
2. 若有 treeNodeId / nodeId，则同步 `tree.store.setSelectedNode(nodeId)`
3. `workspace.store.setContext('knowledge')`
4. 右侧打开摘录详情

---

### 3. 打开 knowledge node 结果

当结果类型为 `node` 时：

1. `tree.store.setSelectedNode(nodeId 或 id)`
2. `workspace.store.setContext('knowledge')`
3. 右侧显示对应节点详情或摘录详情（如果该节点实际为 extract 节点）

---

### 4. 打开 review 结果

当结果类型为 `review` 时：

1. `review.store.setSelectedItem(reviewId 或 id)`
2. `workspace.store.setContext('review')`
3. 右侧打开 review card

---

## 十、组件要求

---

### 1. `GlobalSearch.vue`

这是 Search 页面主容器。

要求：

1. 使用 `SearchToolbar.vue`
2. 使用 `SearchFilterTabs.vue`
3. 使用 `SearchResultList.vue`
4. 首次挂载时调用 `search.store.initialize()`
5. 从 `search.store.ts` 读取 query / filter / results
6. 无结果时显示 `SearchEmptyState.vue`
7. 风格像工作台搜索中心，不是后台检索页面

建议布局：

* 顶部：搜索输入和筛选
* 主体：结果列表
* 底部或侧边可留少量搜索提示

---

### 2. `SearchToolbar.vue`

职责：

* 提供主搜索输入框
* 可带一个清空按钮
* 输入 query 后实时触发搜索

要求：

1. 输入框显眼但克制
2. 支持 placeholder，例如：

   * Search articles, extracts, nodes, reviews...
3. 支持回车时打开当前选中结果（可选）
4. 风格像知识工作台搜索，不像后台过滤表单

---

### 3. `SearchFilterTabs.vue`

职责：

* 提供类型筛选入口

至少包括：

* All
* Articles
* Extracts
* Nodes
* Reviews

要求：

1. 切换 filter 后立即更新结果
2. 视觉紧凑
3. 不要做成复杂 Tabs 页面

---

### 4. `SearchResultList.vue`

职责：

* 展示搜索结果列表
* 使用 `SearchResultItem.vue`
* 当前选中项高亮
* 支持点击打开结果

要求：

* 列表可滚动
* 无结果时展示空状态
* 不要表格化

---

### 5. `SearchResultItem.vue`

职责：

* 渲染单条搜索结果

至少展示：

* title
* entity type
* snippet
* subtitle/source
* updatedAt（如有）

行为要求：

1. 点击时打开结果
2. 当前选中项高亮
3. hover 有轻微反馈

建议视觉结构：

* 第一行：标题 + 类型 tag
* 第二行：snippet
* 第三行：source / 更新时间 / 附加说明

要更像知识结果列表，而不是数据库记录表。

---

### 6. `SearchEmptyState.vue`

职责：

* 展示空查询或无结果状态

至少支持两种状态：

1. 当前无 query，提示用户输入搜索
2. 当前有 query 但无结果，提示未找到

如果你采用“query 为空显示最近项”的方案，则空状态只用于真正无结果场景。

---

## 十一、Workspace 集成要求

请将 Search 模块集成到当前工作台中。

---

### 对 `CenterPane.vue` 的要求

当 `context === 'search'` 时：

* 中间区域显示 `GlobalSearch.vue`
* 不再只是简单占位文本

---

### 对 `RightPane.vue` 的要求

当前阶段对 `search` context，建议采用以下方案：

### 推荐方案

* `search` context 下，右侧保持轻量说明面板
* 点击结果后直接切换到目标上下文，而不是在 `search` 页右侧做复杂预览

原因：

* 第一阶段最重要的是“搜到并跳过去”
* 不必在搜索页右侧再复制一套详情渲染逻辑
* 可以降低复杂度，保持模块清晰

请按这个方向实现，并在最终说明中说明理由。

---

### 对 `SearchPage.vue` 的要求

* 继续使用 `WorkspaceShell`
* context 为 `search`
* 页面本身保持轻量，不写复杂业务逻辑

---

## 十二、与命令面板的联动要求

请补充命令系统中的搜索入口。

建议补充以下命令：

* `go-search`
* `focus-search`

### 要求

执行 `focus-search` 时：

1. `workspace.store.setContext('search')`
2. 搜索输入框自动聚焦

你可以采用一种清晰简单的实现方式，例如：

* search store 中维护一个 `focusToken`
* 或 `GlobalSearch.vue` 通过 `nextTick + ref` 处理
* 请选择一种简单可靠的方案并说明

---

## 十三、与现有模块的联动要求

请确保 Search 结果可以真正串起现有对象。

---

### 联动 1：文章 -> Search -> 文章

搜索到文章结果后：

* 跳到 reading
* 打开对应文章

---

### 联动 2：摘录 -> Search -> 摘录

搜索到摘录结果后：

* 跳到 knowledge
* 打开摘录详情
* 尽量同步树选中

---

### 联动 3：节点 -> Search -> 节点

搜索到知识节点后：

* 跳到 knowledge
* 选中该节点

---

### 联动 4：复习项 -> Search -> Review

搜索到复习项后：

* 跳到 review
* 打开该复习项

---

## 十四、mock 数据要求

请根据需要补充 mock 数据，使搜索结果看起来更真实。

### `mock/articles.ts`

确保 title / summary / content 足够有辨识度。

### `mock/tree.ts`

确保存在可搜索的 topic / article / extract / card 节点标题。

### `mock/review.ts`

确保 prompt / answer / title 有差异度。

### `extract` 数据

若当前项目已有 extract store 或 mock extract，请确保其中有若干条真实感较强的摘录对象。

目标是：

* 搜索不同关键词时确实能命中不同对象
* 搜索页面不是“假搜索”

---

## 十五、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 专业、克制
3. 轻边框
4. 小圆角
5. 更像知识工作台搜索中心，不像后台筛选页
6. 搜索结果列表要像“知识对象结果”
7. 各种实体类型要可辨识，但不要太花
8. 搜索输入区域要清晰但不过度抢眼

建议视觉感觉：

* 一个桌面知识工作台的统一搜索入口
* 结果像“对象清单”
* 点击即可进入对应工作流

---

## 十六、代码结构要求

请特别遵守：

1. 不要把搜索映射逻辑写在组件里
2. 不要把搜索执行和结果跳转写在模板中
3. `search.mapper.ts`、`search.index.ts`、`search.matcher.ts` 分层清晰
4. `search.store.ts` 负责搜索状态与打开结果
5. `GlobalSearch.vue` 只负责界面和交互
6. 所有 import 必须完整
7. 不要只写伪代码或片段
8. 代码应为“第一阶段真实可运行版本”

---

## 十七、建议的实现策略

为了降低复杂度，推荐采用以下方案：

### 推荐方案

* Search 页面中间只做统一结果列表
* 搜索结果点击后直接跳转离开 search context
* query 为空时显示最近更新项
* matcher 使用简单 title/content/keyword 包含匹配
* 类型 filter 使用前端过滤

### 为什么推荐这样做

因为第一阶段最重要的是：

* 统一搜
* 搜得到
* 点得开
* 能跳到对应对象

这已经足够构成可用搜索工作流。

---

## 十八、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text id="m2ufd9"
src/modules/search/types.ts
```

然后再给完整代码。

### 第四步：最后说明：

1. 如何接入当前项目
2. 用户如何在 Search 页面中搜索并跳转到对象
3. 当前实现相对真实全文搜索版的简化点有哪些

---

## 十九、补充约束

```text id="9u5x4k"
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入全文搜索库或重量级 fuzzy search 库
- 可以使用少量 Element Plus Input / Empty / Tag / Button / Tabs 等基础控件
- 当前阶段优先保证结构清晰、模块边界清晰、可运行
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 二十、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。
