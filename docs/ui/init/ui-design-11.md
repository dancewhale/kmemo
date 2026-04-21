# 任务单 11：视觉与交互打磨（Polish）第一阶段

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **“视觉与交互打磨（Polish）第一阶段”**。

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
* 三栏工作台
* 搜索、命令、复习、摘录、知识树等基础联动
* 全局样式、基础 shared 组件、Pane/Splitter/StatusBar 等基础设施

本次目标是：**在不大改架构的前提下，对整个系统进行第一阶段视觉与交互统一打磨**，让产品更像一个真正可长期使用的桌面知识工具。

当前阶段重点是：

* 统一视觉语言
* 统一组件状态表现
* 统一空状态/加载态/错误态
* 提升 pane / splitter / 列表 / toolbar 的手感
* 增加必要但克制的交互反馈

不是新增复杂业务功能。

---

## 一、模块目标

请实现如下目标：

1. 统一整个应用的视觉基线：

   * 字号层级
   * 间距体系
   * 边框与圆角
   * hover / selected / active 状态
   * pane 标题和 toolbar 风格

2. 统一以下基础状态表现：

   * empty
   * loading
   * error
   * selected
   * disabled

3. 统一以下高频交互区域的手感：

   * 左侧导航
   * 中间列表
   * 树节点
   * 右侧详情/编辑区
   * 底部状态栏
   * splitter
   * 命令面板
   * 搜索结果
   * review 队列与 review 卡片

4. 提供基础全局反馈方式：

   * 操作成功提示
   * 轻量错误提示
   * 状态栏/通知统一风格

5. 在不重写模块的前提下，对 UI 进行系统性“收口”

---

## 二、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要新增复杂业务模块
3. 不要重写已有核心 store 逻辑
4. 不要引入大规模动画库
5. 不要做花哨视觉升级
6. 不要做营销站风格改版
7. 不要引入大型 design system 依赖
8. 不要把所有视觉逻辑散到每个业务组件里
9. 不要对架构做颠覆式重构
10. 不要把这个任务变成“重新写一遍 UI”

---

## 三、本次建议实现/修改的文件

请至少实现和补全以下文件：

```text id="4jlwmk"
src/shared/components/
  AppPane.vue
  AppHeader.vue
  AppEmpty.vue
  AppLoading.vue
  AppSplitter.vue
  AppStatusTag.vue
  AppActionBar.vue
  AppSectionTitle.vue
  AppFeedbackToast.vue

src/shared/composables/
  useToast.ts

src/shared/constants/
  ui.ts

src/app/styles/
  variables.scss
  theme.scss
  layout.scss
  index.scss
  reset.scss
  utilities.scss

src/modules/workspace/components/
  LeftSidebar.vue
  CenterPane.vue
  RightPane.vue
  BottomStatusBar.vue
  WorkspaceShell.vue
```

并根据需要补充或修改以下文件：

```text id="9jlwmf"
src/modules/reader/components/ReaderList.vue
src/modules/reader/components/ReaderListItem.vue
src/modules/reader/components/ReaderToolbar.vue

src/modules/knowledge-tree/components/TreePanel.vue
src/modules/knowledge-tree/components/TreeNodeItem.vue

src/modules/editor/components/EditorShell.vue
src/modules/editor/components/EditorToolbar.vue
src/modules/editor/components/TiptapEditor.vue
src/modules/editor/components/ExtractToolbar.vue

src/modules/review/components/ReviewQueue.vue
src/modules/review/components/ReviewQueueItem.vue
src/modules/review/components/ReviewCard.vue
src/modules/review/components/ReviewActions.vue
src/modules/review/components/ReviewStatsBar.vue

src/modules/search/components/GlobalSearch.vue
src/modules/search/components/SearchToolbar.vue
src/modules/search/components/SearchResultList.vue
src/modules/search/components/SearchResultItem.vue

src/modules/command-center/components/CommandPalette.vue
src/modules/command-center/components/CommandList.vue
src/modules/command-center/components/CommandListItem.vue

src/modules/settings/components/SettingsPanel.vue
src/modules/settings/components/SettingsSection.vue
```

如有必要，你可以新增少量辅助文件，例如：

```text id="5jlwmh"
src/shared/utils/ui.ts
src/shared/composables/useInteractiveState.ts
```

但请保持克制，不要扩展成完整设计系统工程。

---

## 四、整体设计要求

请把这次 polish 拆成 5 个清晰层次：

### 1. 视觉 token 层

负责：

* 间距
* 字号
* 颜色
* 边框
* 圆角
* 阴影
* 层级

### 2. 基础共享组件层

负责：

* Pane 壳
* Header 壳
* Empty / Loading / StatusTag / ActionBar 等基础视觉单元

### 3. 状态表现统一层

负责：

* hover
* active
* selected
* disabled
* loading
* empty
* error

### 4. 高密度列表与工作区层

负责：

* 树
* 列表
* 结果项
* review 队列
* 搜索结果
* 导航项

### 5. 轻量反馈层

负责：

* toast
* 状态提示
* 轻量错误提示
* splitter 手感
* toolbar 交互反馈

要求边界清晰，尽量复用 shared，而不是每个模块各做各的。

---

## 五、总体风格要求

请严格遵循以下风格目标：

1. **高信息密度**
2. **桌面知识工具风格**
3. **专业、克制**
4. **轻边框、低饱和背景**
5. **小圆角**
6. **明显但不花哨的 selected 状态**
7. **统一 hover / active 反馈**
8. **更像 SuperMemo / 专业知识工作台**
9. **不要像后台管理系统**
10. **不要像 Notion 营销站或卡片社区**

关键词可以理解为：

* dense
* compact
* quiet
* focused
* workbench
* knowledge tool
* keyboard-first
* practical

---

## 六、视觉 token 要求

请在 `app/styles/variables.scss` 和 `theme.scss` 中统一整理视觉变量。

至少整理这些 token：

### 颜色

* app background
* pane background
* panel subtle background
* border
* text primary
* text secondary
* muted text
* hover background
* selected background
* active accent
* success
* warning
* danger
* info

### 间距

* xs
* sm
* md
* lg
* xl

### 字号

* caption
* small
* body
* title
* section
* page title

### 圆角

* sm
* md
* lg

### 阴影

* 轻量 hover shadow
* overlay shadow（命令面板/dialog 可用）

### 层级

* command palette
* dialog
* floating toolbar
* base content

要求：

* 尽量用 CSS 变量或 SCSS 变量统一
* light/dark 模式兼容
* 不要把颜色直接写死在组件里

---

## 七、共享组件打磨要求

---

### 1. `AppPane.vue`

请优化为统一 pane 壳。

要求：

1. 标题区、内容区、footer 区层次清晰
2. 边框和背景风格统一
3. 支持紧凑模式
4. 支持 scrollable
5. 支持 selected / focus-like 状态（可选）
6. 用于整个应用的 pane 统一基线

---

### 2. `AppHeader.vue`

请统一作为 pane/header 工具条的基线组件。

要求：

1. 支持 title
2. 支持 subtitle
3. 支持右侧 actions slot
4. 间距和字号统一
5. 可用于 ReaderToolbar / TreePanel header / SearchToolbar header 等

---

### 3. `AppEmpty.vue`

请增强为空状态组件。

要求：

1. 至少支持：

   * title
   * description
   * icon slot 或 type
   * action slot
2. 空状态风格统一
3. 不要每个模块各写各的空状态文案和样式

---

### 4. `AppLoading.vue`

请统一为简洁加载态。

要求：

1. 可显示轻量 spinner / skeleton 风格之一
2. 支持文字说明
3. 风格克制
4. 可在列表、pane、详情区复用

---

### 5. `AppSplitter.vue`

这是本次非常重要的 polish 点之一。

请优化：

1. hover 时有明确但克制的视觉反馈
2. dragging 时状态清晰
3. 命中区域比视觉线稍宽，方便拖拽
4. 不要太花哨
5. 拖拽手感要更像桌面工具分割条

---

### 6. `AppStatusTag.vue`

新增或完善轻量状态标签组件。

要求：

* 用于 article status / review status / entity type / source type 等
* 尺寸紧凑
* 颜色克制
* 统一视觉风格

---

### 7. `AppActionBar.vue`

新增或完善轻量工具条容器。

要求：

* 统一 toolbar 区间距
* 统一按钮布局
* 适用于 editor toolbar、reader toolbar、search toolbar 等

---

### 8. `AppSectionTitle.vue`

新增或完善小节标题组件。

要求：

* 用于右侧详情中的小节标题，如：

  * Overview
  * Source
  * Summary
  * Related Extracts
* 统一层级和间距

---

### 9. `AppFeedbackToast.vue` + `useToast.ts`

请实现轻量 toast 机制。

要求：

* 支持 success / info / warning / error
* 可用于：

  * 创建摘录成功
  * 保存成功
  * 复制成功
  * 操作失败
* 当前阶段实现简单全局 toast 即可
* 不要依赖复杂第三方通知系统

---

## 八、状态表现统一要求

请统一以下状态表现，并尽量通过 shared/class/token 实现，而不是每个组件手写。

### 1. hover

用于：

* 树节点
* 列表项
* 搜索结果项
* review queue item
* 左侧导航项
* 命令列表项

要求：

* hover 背景统一
* 不同组件不应出现风格割裂

### 2. selected

用于：

* 当前文章
* 当前树节点
* 当前搜索结果
* 当前 review item
* 当前导航项

要求：

* selected 明确但不过亮
* 统一左边条/背景/边框风格之一
* 建议全项目尽量统一为一套 selected 语言

### 3. active / pressed

用于：

* toolbar button
* command item
* action buttons
* splitters dragging

要求：

* 反馈明确
* 不夸张

### 4. disabled

用于：

* 未满足条件的命令
* 未显示答案前的 review grading 按钮
* 暂不可用的 extract 操作

要求：

* 可识别
* 但不要过灰导致难看

### 5. loading / empty / error

要求：

* 三者表现统一
* 各模块不要再各写一套杂乱的状态 UI

---

## 九、模块级 polish 要求

---

### A. Workspace / Pane

#### `WorkspaceShell.vue`

要求：

1. 左中右区域的视觉边界更统一
2. pane 间距和边框统一
3. 不要看起来像临时拼接的 div

#### `LeftSidebar.vue`

要求：

1. 导航项高度、图标、文字、激活态统一
2. 更像桌面知识工具左导航
3. 支持 compact 风格

#### `BottomStatusBar.vue`

要求：

1. 更轻更稳
2. 信息密度高但不乱
3. 若显示快捷键提示，要风格统一
4. 与 pane 边界统一

---

### B. Reader

#### `ReaderList.vue` / `ReaderListItem.vue`

要求：

1. 列表项标题、摘要、meta 层级更清楚
2. selected 状态与树/搜索结果项保持一致风格
3. 摘要截断更自然
4. status/sourceType 标签统一使用共享状态标签

#### `ReaderToolbar.vue`

要求：

1. 搜索输入、筛选控件、标题区统一
2. 不要像后台列表过滤器
3. 更像阅读队列工具条

---

### C. Knowledge Tree

#### `TreePanel.vue` / `TreeNodeItem.vue`

要求：

1. 节点缩进、箭头、图标、标题间距统一
2. 当前节点高亮手感更好
3. 展开/折叠点击目标更合理
4. hover/selected 明确但不吵

---

### D. Editor / Extract

#### `EditorShell.vue`

要求：

1. 标题区、meta、toolbar、正文区、related extracts 区边界清楚
2. 正文区更像“工作区内容面板”，不是普通 textarea 页面
3. 阅读宽度和字号应用更自然

#### `EditorToolbar.vue`

要求：

1. 按钮密度统一
2. 分组合理
3. hover / active 反馈统一
4. 视觉轻量

#### `TiptapEditor.vue`

当前即便还是简化实现，也请优化：

1. 内边距
2. 可读性
3. 滚动体验
4. placeholder 风格
5. focus 状态

#### `ExtractToolbar.vue`

要求：

1. 更像轻量浮动操作条
2. 按钮视觉统一
3. 不要像突兀弹层

---

### E. Review

#### `ReviewQueue.vue` / `ReviewQueueItem.vue`

要求：

1. 队列项像“待处理复习任务”
2. 类型/status/meta 层级清楚
3. selected 与 hover 统一

#### `ReviewCard.vue`

要求：

1. 问题、答案、来源、元信息、操作区层级更清楚
2. “Show Answer” 和评分区节奏更自然
3. 不是表单页，而是复习工作面板

#### `ReviewActions.vue`

要求：

1. Again / Hard / Good / Easy 的按钮层级清楚
2. 不必很花，但要一眼能操作
3. disabled / loading 明确

#### `ReviewStatsBar.vue`

要求：

1. 统计信息更像工作台指标条
2. 与其他 toolbar/pane 统一

---

### F. Search

#### `GlobalSearch.vue` / `SearchToolbar.vue`

要求：

1. 搜索输入区域更像“统一搜索中心”
2. filter tabs 更紧凑
3. 页面不是后台筛选面板

#### `SearchResultList.vue` / `SearchResultItem.vue`

要求：

1. 结果项 title/snippet/meta 层级清楚
2. type tag 统一
3. selected/hover 与 reader/tree/review 列表项统一

---

### G. Command Palette

#### `CommandPalette.vue`

要求：

1. overlay 层级清楚
2. 面板大小、圆角、阴影、边框更统一
3. 输入框和列表之间节奏更自然
4. 更像桌面工具 command center

#### `CommandList.vue` / `CommandListItem.vue`

要求：

1. 当前 active item 清楚
2. group/shortcut/title/subtitle 层级清楚
3. disabled 命令状态可识别

---

### H. Settings

#### `SettingsPanel.vue` / `SettingsSection.vue`

要求：

1. 设置分区更统一
2. 字段标签、说明、控件间距更规整
3. 不像后台配置堆砌页面
4. 更像桌面工具 Preferences

---

## 十、错误与反馈要求

请建立统一的轻量反馈约定。

### 至少覆盖这些场景：

* 创建摘录成功
* 保存成功
* 复制成功
* 重置设置成功
* 操作失败
* 空结果提示

### 要求：

1. 优先使用轻量 toast
2. 文案短、明确
3. 不要全靠 `console.log`
4. 不要到处直接使用不同风格的 Element Plus Message 而不统一封装

### 推荐方案

* 用 `useToast.ts` 做一个轻量封装
* 内部可调用 Element Plus Message，但项目层统一出口

---

## 十一、空状态/加载态/错误态文案统一要求

请为主要区域统一一套语气风格。

### 语气要求：

* 专业
* 简洁
* 不卖萌
* 不啰嗦
* 像桌面工具

### 示例风格

* No article selected
* No extracts found for this article
* No review items available
* Start typing to search your workspace
* Nothing matched your query
* Unable to complete this action

不要出现太多风格混乱的文案。

---

## 十二、样式组织要求

请新增或完善：

### `utilities.scss`

放少量高频工具类，例如：

* text truncation
* subtle text
* compact row
* selected background
* scroll container
* interactive item base
* visually hidden（如有需要）

### `ui.ts`

可选用于统一：

* toast duration
* pane paddings
* compact heights
* z-index config

要求：

* 不要让样式变量和常量散乱
* 不要过度 utility-first 化
* 仍以组件样式 + token 为主

---

## 十三、MainLayout / 全局层 polish 要求

请检查并改善：

### `MainLayout.vue`

要求：

1. 命令面板、toast、全局 overlay 层级关系清晰
2. 根容器背景统一
3. 主题切换下 light/dark 表现自然
4. 不要出现某些 pane 背景割裂、某些区域边框风格不同的问题

---

## 十四、性能与实现边界要求

请注意，这次是 polish，不是重构。

要求：

1. 优先做样式统一和小型组件抽象
2. 不要把每个模块大改一遍
3. 不要引入重依赖
4. 尽量复用已有 shared 组件
5. 保持业务逻辑不被破坏
6. 若某些组件已有实现，只做必要改良，不要推倒重来

---

## 十五、代码结构要求

请特别遵守：

1. 不要把所有视觉样式写在一个超大 SCSS 文件里
2. 不要把所有 toast/error/loading 逻辑散到每个业务组件里
3. 尽量让 shared 组件承担视觉统一责任
4. 样式 token、共享组件、业务组件这三层边界清晰
5. 所有 import 必须完整
6. 不要只写伪代码或片段
7. 代码应为“第一阶段真实可运行 polish 版本”

---

## 十六、建议的实现策略

为了降低复杂度，推荐采用以下策略：

### 推荐方案

* 先整理 token
* 再打磨 shared 组件
* 再把业务模块接到共享视觉语言上
* 最后补统一 toast / empty / loading / splitter 反馈

### 为什么推荐这样做

因为这会最大化收益：

* 改动不大
* 收益明显
* 全局统一感会快速提升
* 后续第二阶段继续做右键菜单、批量操作、交互细节时也更稳

---

## 十七、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：说明你计划如何分层完成本次 polish

例如：

* token 层
* shared 基础组件层
* 模块接入层
* 反馈层

### 第四步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text id="j1jlwm"
src/shared/components/AppStatusTag.vue
```

然后再给完整代码。

### 第五步：最后说明：

1. 如何接入当前项目
2. 这次 polish 统一了哪些视觉与交互规则
3. 当前实现相对完整 design system / 完整产品打磨版的简化点有哪些

---

## 十八、补充约束

```text id="2jlwmn"
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入大型动画库或设计系统库
- 可以使用少量 Element Plus 基础组件，但视觉统一要由项目自身控制
- 当前阶段优先保证结构清晰、可运行、风格统一
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 十九、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。
