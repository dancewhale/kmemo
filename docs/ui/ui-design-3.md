# 任务单 3：reader 模块

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **第一阶段的 reader 模块**。

这个项目是一个基于 **Wails + Go 后端** 的桌面知识工具，产品定位参考 **SuperMemo 风格的渐进阅读 / 知识工作台**。
本次只实现 **reader 模块的第一阶段可运行版本**，要求 **先用 mock 数据，不接真实后端，不实现复杂阅读器能力**，但要为后续接入 **文章预览、摘录、渐进阅读流程** 预留清晰边界。

---

## 一、模块目标

请实现一个 **reader 模块的第一阶段骨架**，能够完成：

1. 展示文章列表
2. 支持文章搜索输入
3. 支持文章状态筛选
4. 支持选中当前文章
5. 支持当前文章高亮
6. 支持右侧详情区展示当前文章摘要和元信息
7. 使用 Pinia 管理阅读列表状态
8. 使用 mock 数据驱动
9. 为后续接入真实阅读器 / Tiptap / 摘录操作预留清晰边界

当前阶段 **不要实现完整富文本阅读器**，也 **不要实现文章导入、网页解析、PDF 阅读、摘录工具条**。

---

## 二、项目上下文说明

项目前两阶段已经具备如下结构与约束：

* 有 `MainLayout`
* 有 `WorkspaceShell`
* 有 `ReadingPage.vue`
* 有 `workspace.store.ts`
* 有 `knowledge-tree` 模块
* 有 `shared/components` 基础组件
* 有 `mock/` 目录
* 页面通过 `WorkspaceShell` 根据 `context` 切换内容
* 当前 `ReadingPage.vue` 以 `context="reading"` 渲染工作台
* 当前项目强调：

  * 高信息密度
  * 桌面知识工具风格
  * 不是后台管理系统
  * 模块边界清晰
  * 组件职责清晰

---

## 三、本次要实现的文件

请至少实现和补全以下文件：

```text
src/modules/reader/
  components/
    ReaderList.vue
    ReaderListItem.vue
    ReaderToolbar.vue
  stores/
    reader.store.ts
  types.ts
```

并根据需要补充或修改以下文件，使 `reading` context 真正可运行：

```text
src/mock/articles.ts
src/mock/index.ts
src/modules/workspace/components/CenterPane.vue
src/modules/workspace/components/RightPane.vue
```

如有必要，可以新增极少量辅助文件，但请控制复杂度，不要扩展成完整子系统。

---

## 四、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现 Tiptap
3. 不要实现真实网页阅读器
4. 不要实现 PDF 阅读
5. 不要实现选区摘录
6. 不要实现复杂排序和多条件过滤器面板
7. 不要实现拖拽列表
8. 不要做成后台管理系统表格页
9. 不要把所有 reader 逻辑塞进一个组件
10. 不要提前实现未来 editor / review / import 的业务逻辑

---

## 五、类型设计要求

请先定义清晰的类型。

### `types.ts` 至少包括：

```ts
export type ReaderSourceType = 'web' | 'pdf' | 'text'
export type ReaderArticleStatus = 'inbox' | 'reading' | 'processed'

export interface ReaderArticle {
  id: string
  title: string
  summary: string
  sourceType: ReaderSourceType
  sourceUrl?: string
  status: ReaderArticleStatus
  updatedAt: string
  tags?: string[]
}

export interface ReaderFilterState {
  keyword: string
  status: 'all' | ReaderArticleStatus
}
```

你可以在此基础上微调，但请满足以下要求：

1. 类型边界清晰
2. 能支持后续对接真实后端 DTO
3. 不要过度复杂
4. 文章对象要足够支撑列表和详情区展示

---

## 六、mock 数据要求

请在 `src/mock/articles.ts` 中提供更真实的 mock 数据，至少包含：

* 8~12 条文章
* 标题不要是 `Article 1 / Article 2`
* 要有不同状态：

  * inbox
  * reading
  * processed
* 要有不同类型：

  * web
  * pdf
  * text

建议示例风格类似：

* Incremental Reading Basics
* How to Build a Personal Knowledge System
* Active Recall vs Passive Review
* Learning in Public
* Designing Reading Workflows
* Why Notes Should Be Atomic
* Progressive Summarization for Dense Material

每条文章至少要有：

* id
* title
* summary
* sourceType
* status
* updatedAt

可选：

* sourceUrl
* tags

这些 mock 数据要有“真实知识工具内容感”，不要太像测试假数据。

---

## 七、store 设计要求

请实现 `reader.store.ts`，使用 Pinia 管理阅读列表状态。

### 至少包含这些状态：

```ts
interface ReaderState {
  articles: ReaderArticle[]
  selectedArticleId: string | null
  filter: ReaderFilterState
  loading: boolean
}
```

### 至少包含这些 getter / 计算逻辑：

1. `filteredArticles`
2. `selectedArticle`
3. `inboxArticles`
4. `readingArticles`
5. `processedArticles`

### 至少包含这些 action：

1. `initialize()`
2. `setSelectedArticle(id: string)`
3. `setKeyword(keyword: string)`
4. `setStatusFilter(status: 'all' | 'inbox' | 'reading' | 'processed')`
5. `clearFilters()`

### 行为要求：

* 初始化时加载 mock 数据
* 默认选中第一篇可见文章
* 搜索按标题和 summary 的简单包含匹配即可
* 筛选和搜索应同时生效
* store 职责清晰，不要与 workspace store 混成一体

---

## 八、组件要求

---

### 1. `ReaderToolbar.vue`

这是阅读列表上方的轻量工具栏。

要求：

1. 展示模块标题，例如 `Reading Queue`
2. 提供一个搜索输入框
3. 提供一个状态筛选器
4. 提供一个清空筛选按钮
5. 从 `reader.store.ts` 读取和更新筛选状态
6. 保持紧凑，不要做复杂筛选面板

筛选项至少包括：

* All
* Inbox
* Reading
* Processed

风格要求：

* 紧凑
* 桌面知识工具感
* 不要后台检索表单风格

---

### 2. `ReaderList.vue`

这是阅读列表主容器组件。

要求：

1. 使用 `ReaderToolbar.vue`
2. 从 `reader.store.ts` 读取 `filteredArticles`
3. 首次挂载时调用 `initialize()`
4. 使用 `ReaderListItem.vue` 渲染单篇文章
5. 无数据时显示 `AppEmpty`
6. 加载态时显示 `AppLoading`

列表风格要求：

* 不用表格
* 使用紧凑列表
* 每项突出标题
* 显示摘要和状态信息
* 适合知识工作台，而不是管理后台

---

### 3. `ReaderListItem.vue`

这是单篇文章列表项组件。

props 建议至少包括：

* `article: ReaderArticle`
* `selected: boolean`

行为要求：

1. 点击列表项时选中当前文章
2. 展示标题
3. 展示 summary 摘要（简短截断）
4. 展示 sourceType
5. 展示 status
6. 展示 updatedAt
7. 当前选中项高亮
8. hover 有轻微反馈
9. 样式紧凑、专业

建议列表项视觉结构类似：

* 第一行：标题 + 状态标签
* 第二行：摘要
* 第三行：sourceType / 更新时间 / tags（如有）

注意：

* 不要用后台表格
* 不要做成卡片瀑布流
* 更像阅读工作台里的任务列表

---

## 九、Workspace 集成要求

请将 reader 模块集成到当前 `reading` context 中。

### 对 `CenterPane.vue` 的要求

当 `context === 'reading'` 时：

* 中间区域应显示 `ReaderList.vue`
* 不再只是简单占位文本

### 对 `RightPane.vue` 的要求

当 `context === 'reading'` 时：

* 根据当前选中文章显示详情信息
* 至少显示：

  * title
  * id
  * sourceType
  * status
  * updatedAt
  * summary
  * sourceUrl（如有）
* 如果未选中文章，显示空状态
* 右侧面板依然保持“详情面板 / 预览占位”的风格，而不是表单页面

另外，右侧详情区请预留未来扩展感，例如可使用这样的布局语义：

* Overview
* Source
* Summary
* Tags

但当前阶段不要做复杂编辑。

---

## 十、与 workspace.store 的协作要求

请注意 `reader.store.ts` 与 `workspace.store.ts` 的边界：

### `reader.store.ts`

负责：

* 文章列表数据
* 当前选中文章
* 搜索筛选状态

### `workspace.store.ts`

负责：

* 当前工作台上下文
* pane 尺寸
* 全局选中对象 id（如果当前结构已有）

如果项目当前已有 `selectedArticleId` 在 `workspace.store.ts` 中，可以在 reader 模块选中文章时同步更新，但请保持边界清晰：

* reader 模块内部逻辑以 `reader.store.ts` 为主
* workspace store 只做轻量同步
* 不要把 reader 列表逻辑迁移到 workspace store

---

## 十一、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 紧凑列表
3. 小圆角
4. 轻边框
5. 克制的 hover 和 active 状态
6. 更像知识工具，不像管理后台
7. 列表强调阅读任务感
8. 文本层级清楚
9. 不要花哨，不要营销站风格

建议列表项表现出以下感觉：

* 这是“阅读工作队列”
* 每项是一个正在处理的阅读对象
* 用户会频繁点击切换、预览、后续摘录
* 不是数据库记录表

---

## 十二、代码结构要求

请特别遵守：

1. 不要把筛选逻辑写在组件模板里
2. 不要把所有状态都用局部 `ref` 管
3. 不要让 `ReaderList.vue` 超级臃肿
4. 不要硬编码太多魔法字符串
5. 类型要明确
6. 代码应为“第一阶段真实可运行版本”
7. 所有 import 必须完整
8. 不要只写伪代码或片段

---

## 十三、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text
src/modules/reader/types.ts
```

然后再给完整代码。

### 第四步：最后说明如何把它接入当前项目并运行

---

## 十四、附加实现提示

为了降低未来接入真实阅读器和摘录逻辑的成本，请在结构上提前考虑：

* `ReaderList` 独立
* `ReaderToolbar` 独立
* `ReaderListItem` 独立
* `reader.store.ts` 独立管理筛选和当前选中状态
* `RightPane` 当前只做详情占位，未来可以平滑替换为阅读预览 / 编辑器容器

这会使后续接入 Tiptap、文章正文预览、选区摘录时更加自然。

---

## 十五、补充约束

```text
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 可以使用少量 Element Plus 基础控件，但不要让 reader 模块变成后台表格页
- 当前阶段优先保证结构清晰、模块边界清晰、可运行
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 十六、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。

