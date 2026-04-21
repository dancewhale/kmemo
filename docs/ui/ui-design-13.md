# 任务单 13：导入 / 入箱（Inbox Capture）第一阶段（可直接投喂 AI）

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **“导入 / 入箱（Inbox Capture）第一阶段可运行版本”**。

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
* 阅读 / 摘录 / 知识树 / 复习 / 命令 / 搜索 / 偏好 / 右键菜单 等基础联动
* `workspace.store.ts`
* `reader.store.ts`
* `tree.store.ts`
* `editor.store.ts`
* `extract.store.ts`
* `review.store.ts`
* `search.store.ts`
* `command.store.ts`
* `mock/` 数据
* 三栏工作台与统一 UI 风格

本次目标是：**实现第一阶段的 Inbox Capture 能力**，让用户能够通过若干简单入口，把新内容加入系统，并进入 inbox / reading 工作流。

当前阶段重点是：

* 入箱对象创建
* 多种轻量 capture 入口
* 基础表单与对话框
* reader/inbox 联动
* 最小可运行的数据流

不是复杂网页抓取器、PDF 解析器、浏览器剪藏系统、文件解析系统。

---

## 一、模块目标

请实现如下能力：

1. 用户能够创建新的 inbox article
2. 至少支持以下 capture 入口：

   * 新建空白文章
   * 粘贴纯文本创建文章
   * 输入 URL 创建文章占位
3. 所有新建内容默认进入 `inbox`
4. 创建成功后：

   * article 进入 reader 列表
   * 可在 inbox / reading 中看到
   * 可点击打开右侧编辑区
5. 提供一个统一的 “Capture / Add to Inbox” 入口
6. 支持从多个位置触发 capture，例如：

   * Inbox 页面工具栏
   * 左侧导航快捷按钮
   * 命令面板中的命令
   * 可选的右键菜单入口
7. 前端结构要为未来接入真实导入逻辑预留清晰边界

当前阶段 **不要求真实网络抓取与解析**，但要让 URL capture 有一个合理的“占位对象”创建流程。

---

## 二、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现网页正文抓取
3. 不要实现 PDF 导入解析
4. 不要实现文件拖拽导入解析
5. 不要实现浏览器扩展剪藏
6. 不要实现复杂队列式导入任务管理
7. 不要实现多步骤导入向导
8. 不要实现复杂元数据提取
9. 不要把所有导入逻辑塞进一个对话框组件
10. 不要让 inbox capture 和 reader store / editor store 职责混乱

---

## 三、本次建议实现/修改的文件

请至少实现和补全以下文件：

```text
src/modules/inbox-capture/
  components/
    CaptureDialog.vue
    CaptureTypeTabs.vue
    NewArticleForm.vue
    PasteTextCaptureForm.vue
    UrlCaptureForm.vue
    CaptureQuickEntry.vue
  stores/
    capture.store.ts
  types.ts
  services/
    capture.mapper.ts
    capture.validator.ts
```

并根据需要补充或修改以下文件：

```text
src/modules/reader/stores/reader.store.ts
src/modules/reader/components/ReaderToolbar.vue
src/modules/workspace/components/LeftSidebar.vue
src/modules/workspace/components/CenterPane.vue
src/modules/workspace/components/RightPane.vue
src/pages/InboxPage.vue
src/modules/command-center/services/command.registry.ts
src/modules/command-center/services/command.executor.ts
src/modules/context-actions/services/context-menu.registry.ts
src/modules/context-actions/services/context-menu.executor.ts
src/mock/articles.ts
src/mock/index.ts
```

如有必要，你可以新增少量辅助文件，例如：

```text
src/shared/constants/inbox.ts
src/shared/utils/url.ts
```

但请保持克制，不要扩展成完整 import pipeline。

---

## 四、整体设计要求

请把这次能力拆成 4 个清晰层次：

### 1. Capture 输入层

负责：

* 不同 capture 类型的表单输入
* 基础校验
* 表单切换
* UI 交互

### 2. Capture 状态层

负责：

* 对话框打开/关闭
* 当前 capture 类型
* 提交状态
* 表单草稿状态（第一阶段可轻量）

### 3. Capture 映射层

负责：

* 将表单输入映射成 article 对象
* 统一默认字段
* 生成标题、summary、content 等基础内容

### 4. Reader / Inbox 联动层

负责：

* 把新 article 加入 reader store
* 默认标记为 `inbox`
* 创建成功后打开该 article
* 切换到合理上下文

要求边界清晰，后续便于接入真实后端导入器。

---

## 五、类型设计要求

请先整理和定义清晰类型。

---

### `src/modules/inbox-capture/types.ts`

至少包含：

```ts
export type CaptureType = 'blank' | 'paste-text' | 'url'

export interface BlankArticleInput {
  title: string
  content: string
  tags?: string[]
}

export interface PasteTextCaptureInput {
  title: string
  content: string
  tags?: string[]
}

export interface UrlCaptureInput {
  url: string
  title?: string
  note?: string
  tags?: string[]
}

export interface CaptureDraftState {
  type: CaptureType
  blank: BlankArticleInput
  pasteText: PasteTextCaptureInput
  url: UrlCaptureInput
}

export interface CaptureSubmissionResult {
  articleId: string
}
```

你可以在此基础上微调，但必须满足：

1. 不同 capture 类型输入结构分开
2. 状态足够支撑多 tab 表单
3. 便于未来接真实导入 DTO
4. 类型不要过度复杂

---

### 对现有 `reader` 类型的要求

请确保 `ReaderArticle` 至少具备：

* `id`
* `title`
* `summary`
* `content`
* `sourceType`
* `sourceUrl?`
* `status`
* `updatedAt`
* `tags?`

因为 capture 最终要创建这种对象并注入 reader store。

---

## 六、Capture 表单与校验要求

请新增：

### `capture.validator.ts`

至少实现这些能力：

1. 校验 blank article：

   * title 可为空但建议自动补默认标题
   * content 至少允许为空或最小长度较低

2. 校验 paste text：

   * content 必须非空
   * title 可为空，若为空则自动生成

3. 校验 url：

   * url 必须非空
   * url 需做基础合法性校验
   * title 可为空

建议至少提供：

```ts
validateBlankInput(...)
validatePasteTextInput(...)
validateUrlInput(...)
```

要求：

* 不要把校验逻辑散落在组件里
* 当前阶段只需做基础校验
* 错误文案简洁、专业

---

## 七、Capture 映射要求

请新增：

### `capture.mapper.ts`

至少实现以下函数：

1. `buildArticleFromBlankInput(...)`
2. `buildArticleFromPasteTextInput(...)`
3. `buildArticleFromUrlInput(...)`

要求：

* 统一生成 `ReaderArticle`
* 默认 `status = 'inbox'`
* 根据 capture 类型生成合理的：

  * title
  * summary
  * content
  * sourceType
  * sourceUrl
  * updatedAt
* 不要把这些映射规则写在表单组件里

### 推荐映射规则

#### blank

* `sourceType = 'text'`
* title 为空时使用 `Untitled Inbox Note`
* summary 从 content 截断生成

#### paste-text

* `sourceType = 'text'`
* title 为空时从前 20~40 字内容生成
* summary 从 content 截断生成

#### url

* `sourceType = 'web'`
* title 优先用用户输入标题
* 未输入标题时可用：

  * `Captured from URL`
  * 或从 URL hostname/path 简单生成
* content 可使用占位结构，例如：

```text
Source URL: ...
Capture note: ...
This article is a placeholder created from a URL capture.
```

* summary 可为：

  * note 的截断
  * 或 URL capture 占位摘要

要求：

* 映射规则清晰、稳定
* 适合当前第一阶段

---

## 八、store 设计要求

---

### `capture.store.ts`

请实现 Inbox Capture 模块专用 store。

至少包含这些状态：

```ts
interface CaptureState {
  isDialogOpen: boolean
  activeType: CaptureType
  submitting: boolean
  draft: CaptureDraftState
}
```

至少包含这些 getter：

1. `activeDraft`
2. `canSubmit`
3. `dialogTitle`

至少包含这些 action：

1. `openDialog(type?: CaptureType)`
2. `closeDialog()`
3. `setActiveType(type: CaptureType)`
4. `updateBlankDraft(payload: Partial<BlankArticleInput>)`
5. `updatePasteTextDraft(payload: Partial<PasteTextCaptureInput>)`
6. `updateUrlDraft(payload: Partial<UrlCaptureInput>)`
7. `resetDraft(type?: CaptureType)`
8. `submit()`
9. `submitBlank()`
10. `submitPasteText()`
11. `submitUrl()`

### 行为要求

* 打开 dialog 时可指定默认 tab
* 提交成功后：

  * 将新 article 加入 reader store
  * 默认选中新 article
  * 切换到 `reading` 或 `inbox` 上下文（你选择一种更合理的方案并说明）
  * 打开右侧编辑区
  * 关闭 dialog
  * toast 提示成功
* store 只负责 capture 状态与提交流程，不要承担 reader 全部逻辑

### 推荐第一阶段方案

* 创建成功后切到 `reading`
* 自动打开新文章
* 原因：用户通常创建后会立即查看/编辑内容

这比留在 inbox 更顺手。

---

## 九、reader.store 联动要求

请在 `reader.store.ts` 中补充基础对象创建能力。

至少新增：

```ts
addArticle(article: ReaderArticle)
```

以及必要时补充：

```ts
updateArticleStatus(id: string, status: 'inbox' | 'reading' | 'processed')
openArticleById(id: string)
```

要求：

1. 新建 article 进入 reader 数据源
2. 新建后能被 inbox / reading 列表立即看到
3. 打开 article 时右侧能联动编辑区
4. 不要把 capture 表单逻辑写进 reader store

---

## 十、组件要求

---

### 1. `CaptureDialog.vue`

这是本次最核心的 UI 组件。

要求：

1. 使用 Element Plus `Dialog` 可以接受
2. 从 `capture.store.ts` 读取状态
3. 包含：

   * 顶部标题
   * `CaptureTypeTabs.vue`
   * 当前对应的表单组件
   * Footer 操作区
4. 支持：

   * Cancel
   * Create / Add to Inbox
5. 提交中显示 loading
6. 切换 tab 时表单状态保留
7. 风格像桌面工具的 capture dialog，不像后台表单弹窗

---

### 2. `CaptureTypeTabs.vue`

职责：

* 切换 capture 类型

至少包括：

* Blank
* Paste Text
* URL

要求：

* 紧凑
* 明确
* 切换后更新 `activeType`

---

### 3. `NewArticleForm.vue`

职责：

* 空白文章表单

至少字段：

* title
* content
* tags（可选，第一阶段可简单输入逗号分隔，或先省略）

要求：

* 支持标题为空
* 支持正文输入
* 更像轻量 capture，不像复杂文档编辑器

---

### 4. `PasteTextCaptureForm.vue`

职责：

* 粘贴文本创建文章

至少字段：

* title
* content（主输入区）
* tags（可选）

要求：

* content 为主要输入区
* 支持多行大文本输入
* 若 title 为空可自动生成

---

### 5. `UrlCaptureForm.vue`

职责：

* 输入 URL 创建占位 article

至少字段：

* url
* title（可选）
* note（可选）
* tags（可选）

要求：

* URL 输入清晰
* 能展示基础说明：

  * 当前阶段仅创建 URL capture placeholder
* 不要误导为已完成网页抓取

---

### 6. `CaptureQuickEntry.vue`

职责：

* 一个轻量 capture 入口组件
* 可用于：

  * 左侧导航区按钮
  * Inbox 页顶部按钮
  * Reader toolbar 中 “Add” 按钮

要求：

* 点击后打开 `CaptureDialog`
* 可接受默认 capture type
* 风格紧凑、适合工具条或侧边栏

---

## 十一、Inbox / Reading / Layout 集成要求

---

### 1. `InboxPage.vue`

要求：

* 继续使用 `WorkspaceShell`
* context 为 `inbox`
* 中间区域能看到 status = inbox 的文章列表
* 提供明显的 Add to Inbox / Capture 入口

---

### 2. `ReaderToolbar.vue`

要求：

* 增加一个 `Add` / `Capture` 按钮
* 打开 `CaptureDialog`
* 不要做成复杂导入器工具栏

---

### 3. `LeftSidebar.vue`

要求：

* 增加一个轻量快捷入口，例如：

  * `+ Capture`
  * 或 `Add to Inbox`
* 点击后打开 capture dialog

---

### 4. `CenterPane.vue`

要求：

* `context === 'inbox'` 时，显示 inbox 文章列表
* `context === 'reading'` 时，显示全部或 reading 相关文章列表
* 确保新创建的 inbox article 能被看到

---

### 5. `RightPane.vue`

要求：

* 创建成功并打开 article 后，右侧应显示文章编辑区
* 不要在 RightPane 里写 capture 表单逻辑
* capture dialog 作为全局弹窗或布局层弹窗存在更合适

---

## 十二、命令面板联动要求

请补充命令系统中的 capture 命令。

建议新增以下命令：

* `capture-blank-article`
* `capture-paste-text`
* `capture-url`
* `go-inbox`

### 要求

#### `capture-blank-article`

* 打开 dialog
* 默认切到 blank

#### `capture-paste-text`

* 打开 dialog
* 默认切到 paste-text

#### `capture-url`

* 打开 dialog
* 默认切到 url

#### `go-inbox`

* 切到 inbox

请在 `command.registry.ts` 与 `command.executor.ts` 中接好这部分联动。

---

## 十三、右键菜单联动要求（可选但推荐）

如果你认为当前结构合适，可补充这些 context actions：

### article 上下文菜单

* Duplicate to Inbox（可选）
* Open Article
* Mark as Reading
* Mark as Processed

### search result 为 article 时

* Open
* Open in Reading
* Copy Title

当前阶段不强求新增太多右键动作，但新增 capture 相关命令入口很有价值。

---

## 十四、搜索联动要求

新创建的 article 必须能被 Search 模块搜到。

要求：

1. 新建文章后，search 索引或搜索数据源能更新
2. 不要求复杂增量索引
3. 可采用简单方案：

   * search 每次搜索时从 store 重新构建 index
   * 或 reader 新增后触发 search store refresh

请选择一种简单清晰的第一阶段方案，并说明。

---

## 十五、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 专业、克制
3. 轻边框
4. 小圆角
5. 更像桌面知识工具的 capture 面板，不像后台导入配置页
6. CaptureDialog 要强调“快速把内容放进系统”
7. 表单输入清晰但不臃肿
8. 提交动作路径短

整体感觉应接近：

* 快速 capture
* 少思考成本
* 立即进入工作流
* 不夸张、不花哨

---

## 十六、文案与反馈要求

请统一文案语气：

* 专业
* 简洁
* 像桌面工具
* 不卖萌

建议文案风格：

* Add to Inbox
* Create Blank Article
* Paste Text
* Capture from URL
* Article added to inbox
* URL placeholder created
* Please enter some text
* Please enter a valid URL

提交成功后请接入统一 toast：

* Blank article created
* Text capture added to inbox
* URL capture placeholder created

---

## 十七、代码结构要求

请特别遵守：

1. 不要把所有 capture 表单写进一个超大组件
2. 不要把 article 创建映射规则写在表单组件里
3. `capture.mapper.ts`、`capture.validator.ts`、`capture.store.ts` 分层清晰
4. `CaptureDialog.vue` 只负责组合 UI
5. reader store 只负责文章数据，不负责 capture 表单状态
6. 所有 import 必须完整
7. 不要只写伪代码或片段
8. 代码应为“第一阶段真实可运行版本”

---

## 十八、建议的实现策略

为了降低复杂度，推荐采用以下方案：

### 推荐方案

* 全局挂一个 `CaptureDialog.vue`
* capture store 管理 dialog open/close 与 active type
* 三种 form 分开
* submit 时 mapper 生成 `ReaderArticle`
* article 加入 reader store
* 切到 reading 并打开文章
* 通过 toast 提示成功

### 为什么推荐这样做

因为这套方案已经足够支撑后续扩展：

* URL 真抓取
* 文件导入
* 剪贴板捕获
* 浏览器扩展剪藏
* OCR / PDF parsing
* 入箱批处理

但当前阶段不会过度复杂。

---

## 十九、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：说明你将如何分层完成本次 capture 系统

例如：

* input layer
* validation layer
* mapping layer
* submission layer
* integration layer

### 第四步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text
src/modules/inbox-capture/types.ts
```

然后再给完整代码。

### 第五步：最后说明：

1. 如何接入当前项目
2. 用户如何通过 blank / paste text / URL 三种方式把内容加入 inbox
3. 当前实现相对真实导入系统的简化点有哪些

---

## 二十、补充约束

```text
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入大型表单库或抓取库
- 可以使用少量 Element Plus Dialog / Input / Tabs / Button / Form 等基础控件
- 当前阶段优先保证结构清晰、可运行、风格统一
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 二十一、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。
