# 任务单 10：设置模块（Settings）+ 工作台偏好（Preferences）

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **“设置模块（Settings）+ 工作台偏好（Preferences）第一阶段可运行版本”**。

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
* 阅读 / 摘录 / 知识树 / 复习 / 命令 / 搜索基本联动
* `workspace.store.ts`
* `reader.store.ts`
* `tree.store.ts`
* `editor.store.ts`
* `extract.store.ts`
* `review.store.ts`
* `command.store.ts`
* `search.store.ts`
* `mock/` 数据
* 主工作台布局
* `SettingsPage.vue` 预留页面

本次目标是：**实现第一阶段的 Settings + Preferences 模块**，让用户能够配置工作台偏好，并把这些偏好保存在本地，影响实际界面行为。

当前阶段重点是 **偏好数据结构、设置界面、持久化、本地应用效果**，不是账号同步和复杂主题系统。

---

## 一、模块目标

请实现如下能力：

1. 用户进入 `settings` 页面
2. 可以修改一组基础工作台偏好
3. 修改后偏好会立刻作用于当前 UI
4. 偏好会保存到 `localStorage`
5. 重启页面后偏好能恢复
6. 至少支持以下偏好项：

### 布局偏好

* 默认工作区（inbox / reading / knowledge / review / search）
* 左侧面板默认宽度
* 右侧面板默认宽度
* 是否默认显示左侧面板
* 是否默认显示右侧面板

### 阅读偏好

* 阅读正文宽度（narrow / medium / wide）
* 编辑区字号（small / medium / large）
* 是否显示文章相关摘录区

### 工作台偏好

* 是否显示底部状态栏
* 是否在启动时恢复上次工作区
* 是否显示快捷键提示
* 命令面板是否显示分组标签

### 外观偏好（第一阶段只做基础版）

* 主题模式（light / dark / system 三选一中的一个基础实现）
* 或者若你认为 `system` 太复杂，至少实现 `light / dark`

7. 提供“恢复默认设置”按钮
8. 提供“导出当前偏好 JSON”能力（可选但推荐）
9. Settings 页面要像桌面工具的偏好中心，不是中后台配置表单页

---

## 二、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现账号同步设置
3. 不要实现云端配置
4. 不要实现复杂主题设计器
5. 不要实现快捷键自定义编辑器
6. 不要实现国际化系统
7. 不要实现复杂实验功能开关系统
8. 不要实现过度复杂的表单 schema 引擎
9. 不要把所有偏好逻辑塞进 SettingsPage 一个文件
10. 不要让 settings store 和其他业务 store 职责混乱

---

## 三、本次建议实现/修改的文件

请至少实现和补全以下文件：

```text id="a5ht0m"
src/modules/settings/
  components/
    SettingsPanel.vue
    SettingsSection.vue
    PreferenceField.vue
    ThemeModeSelect.vue
    WorkspacePreferenceSection.vue
    ReadingPreferenceSection.vue
    AppearancePreferenceSection.vue
    AdvancedPreferenceSection.vue
  stores/
    settings.store.ts
  types.ts
  services/
    settings.defaults.ts
    settings.storage.ts
    settings.mapper.ts
```

并根据需要补充或修改以下文件：

```text id="z4n7ks"
src/pages/SettingsPage.vue
src/layouts/MainLayout.vue
src/modules/workspace/stores/workspace.store.ts
src/modules/editor/components/EditorShell.vue
src/modules/command-center/components/CommandPalette.vue
src/modules/workspace/components/BottomStatusBar.vue
src/shared/utils/storage.ts
src/app/styles/theme.scss
src/app/styles/variables.scss
src/app/styles/layout.scss
```

如果有必要，你可以新增少量辅助文件，例如：

```text id="v7m30g"
src/shared/constants/preferences.ts
src/shared/composables/useThemeMode.ts
src/shared/types/preferences.ts
```

但请保持克制，不要扩展成完整系统设置平台。

---

## 四、整体设计要求

请把这次能力拆成 4 个清晰层次：

### 1. 偏好定义层

负责：

* 默认值
* 偏好类型
* 偏好字段结构
* 偏好版本与兼容基础

### 2. 偏好持久化层

负责：

* 从 localStorage 读取
* 写入 localStorage
* 恢复默认值
* 合并默认值与已存储值

### 3. 偏好状态层

负责：

* 当前偏好 store
* 更新单项配置
* 重置默认
* 导出配置

### 4. 偏好应用层

负责：

* 将偏好应用到工作台
* 将偏好应用到样式/UI显隐
* 尽量不要把应用逻辑写散在很多组件里

要求边界清晰，方便后续接真实后端配置同步。

---

## 五、类型设计要求

请先整理和定义清晰类型。

---

### `src/modules/settings/types.ts`

至少包含：

```ts id="6qv8f2"
export type WorkspaceDefaultContext =
  | 'inbox'
  | 'reading'
  | 'knowledge'
  | 'review'
  | 'search'

export type ReadingWidthMode = 'narrow' | 'medium' | 'wide'
export type EditorFontSizeMode = 'small' | 'medium' | 'large'
export type ThemeMode = 'light' | 'dark' | 'system'

export interface WorkspacePreferences {
  defaultContext: WorkspaceDefaultContext
  leftPaneWidth: number
  rightPaneWidth: number
  showLeftPane: boolean
  showRightPane: boolean
  showBottomStatusBar: boolean
  restoreLastContextOnLaunch: boolean
  showShortcutHints: boolean
  commandPaletteShowGroup: boolean
}

export interface ReadingPreferences {
  contentWidth: ReadingWidthMode
  editorFontSize: EditorFontSizeMode
  showRelatedExtracts: boolean
}

export interface AppearancePreferences {
  themeMode: ThemeMode
}

export interface AdvancedPreferences {
  exportDebugInfo?: boolean
}

export interface AppPreferences {
  workspace: WorkspacePreferences
  reading: ReadingPreferences
  appearance: AppearancePreferences
  advanced: AdvancedPreferences
}
```

你可以在此基础上微调，但必须满足：

1. 分组清晰
2. 偏好字段足够覆盖第一阶段需求
3. 为后续扩展预留空间
4. 类型不要过度复杂

---

## 六、默认值与存储要求

请新增：

### `settings.defaults.ts`

用于定义默认偏好配置。

至少提供：

```ts id="vvjlwm"
export const DEFAULT_PREFERENCES: AppPreferences = { ... }
```

---

### `settings.storage.ts`

用于偏好本地持久化。

至少实现：

1. `loadPreferences(): AppPreferences`
2. `savePreferences(preferences: AppPreferences): void`
3. `resetPreferences(): AppPreferences`
4. `exportPreferences(): string`

要求：

* 使用 `localStorage`
* 需要做默认值合并，避免旧结构数据缺字段时报错
* `exportPreferences()` 返回格式化 JSON 字符串
* 不要把 localStorage 操作散落在 store 和组件里

---

### `settings.mapper.ts`

如果你认为有必要，可以提供：

* 合并默认配置
* 旧配置迁移
* 偏好值规范化

当前阶段可以很轻，但建议结构上预留。

---

## 七、store 设计要求

---

### `settings.store.ts`

请实现 Settings 模块专用 store。

至少包含这些状态：

```ts id="qujlwm"
interface SettingsState {
  preferences: AppPreferences
  initialized: boolean
}
```

至少包含这些 getter：

1. `workspacePreferences`
2. `readingPreferences`
3. `appearancePreferences`
4. `advancedPreferences`
5. `themeMode`
6. `editorFontSizeClass`
7. `contentWidthClass`

至少包含这些 action：

1. `initialize()`
2. `updateWorkspacePreference<K extends keyof WorkspacePreferences>(key: K, value: WorkspacePreferences[K])`
3. `updateReadingPreference<K extends keyof ReadingPreferences>(key: K, value: ReadingPreferences[K])`
4. `updateAppearancePreference<K extends keyof AppearancePreferences>(key: K, value: AppearancePreferences[K])`
5. `updateAdvancedPreference<K extends keyof AdvancedPreferences>(key: K, value: AdvancedPreferences[K])`
6. `resetToDefaults()`
7. `exportPreferences()`
8. `applyPreferences()`

### 行为要求

* 初始化时从 localStorage 加载
* 更新任意偏好后立即保存
* 更新任意偏好后立即应用到 UI
* `applyPreferences()` 是本次关键 action 之一
* store 只负责偏好状态，不要把具体模块业务逻辑都搬进来

---

## 八、偏好应用要求

这是本次最关键的部分之一。
设置不能只是“页面上改了几个值”，必须真的作用到界面。

---

### 1. 应用于 workspace 布局

请让以下偏好生效：

* `leftPaneWidth`
* `rightPaneWidth`
* `showLeftPane`
* `showRightPane`
* `showBottomStatusBar`

要求：

1. 在 `workspace.store.ts` 中保留布局状态
2. `settings.store.applyPreferences()` 能将默认布局偏好同步给 workspace
3. 第一期不要求完全解决“用户手动拖拽后和设置页之间的双向同步冲突”，但要做一个清晰简单的方案

### 推荐第一阶段方案

* Settings 是“默认偏好”
* 当用户点击“应用偏好”或更新设置时，将对应 pane 状态同步到 workspace 当前状态
* 后续再处理更复杂的 session layout / default layout 区分

请采用这个简单可运行的方案，并说明。

---

### 2. 应用于阅读/编辑区

请让以下偏好生效：

* `contentWidth`
* `editorFontSize`
* `showRelatedExtracts`

要求：

1. `EditorShell.vue` 能根据设置改变内容宽度 class
2. `EditorShell.vue` 能根据设置改变正文区字号 class
3. `showRelatedExtracts = false` 时隐藏文章相关摘录区

---

### 3. 应用于状态栏与快捷键提示

请让以下偏好生效：

* `showBottomStatusBar`
* `showShortcutHints`
* `commandPaletteShowGroup`

要求：

1. `BottomStatusBar.vue` 能根据设置显隐
2. `CommandPalette.vue` 可根据设置决定是否显示 group label
3. 若当前项目中其他地方有快捷键 hint，可根据设置决定是否显示

---

### 4. 应用于主题模式

请至少实现基础的 light / dark 模式切换。
如果你觉得 `system` 可在第一阶段稳定实现，也可以加上。

要求：

1. 使用根节点 class 或 data-theme 方案
2. 在 `MainLayout.vue` 或 app 根层应用主题 class
3. `theme.scss` 中提供基础 light/dark 变量差异
4. 不要求重设计整套视觉，只要实现基础主题切换

### 推荐第一阶段方案

* 使用 `document.documentElement.dataset.theme = 'light' | 'dark'`
* SCSS 基于 `[data-theme='dark']` 覆盖变量
* `system` 若实现，使用 `matchMedia('(prefers-color-scheme: dark)')`

如果觉得 `system` 会让当前阶段过重，也可以先只做 light/dark，并在最终说明中说明原因。

---

## 九、组件要求

---

### 1. `SettingsPanel.vue`

这是 Settings 页面主容器。

要求：

1. 负责组合多个设置分区
2. 从 `settings.store.ts` 读取偏好
3. 提供页面级操作：

   * Reset to Defaults
   * Export Preferences（可选推荐）
4. 风格应像桌面工具偏好页，不像后台管理表单

建议布局：

* 顶部：标题 + 简要说明
* 中间：若干设置分区
* 底部：重置 / 导出操作

---

### 2. `SettingsSection.vue`

职责：

* 通用设置分区壳
* 支持：

  * title
  * description
  * default slot

要求：

* 可复用
* 边界清晰
* 视觉统一

---

### 3. `PreferenceField.vue`

职责：

* 通用单个偏好项展示壳

可支持：

* label
* description
* control slot

要求：

* 保持设置页字段一致性
* 不耦合具体业务

---

### 4. `WorkspacePreferenceSection.vue`

职责：

* 展示和编辑 workspace 相关偏好

至少包括：

* 默认工作区
* 左/右面板宽度
* 显示左面板
* 显示右面板
* 显示底部状态栏
* 启动时恢复上次工作区
* 显示快捷键提示
* 命令面板显示分组

---

### 5. `ReadingPreferenceSection.vue`

职责：

* 展示和编辑阅读区相关偏好

至少包括：

* 阅读正文宽度
* 编辑区字号
* 是否显示 Related Extracts

---

### 6. `AppearancePreferenceSection.vue`

职责：

* 展示和编辑主题偏好

至少包括：

* Theme mode

---

### 7. `AdvancedPreferenceSection.vue`

职责：

* 放少量高级项或调试项
* 第一阶段可只放：

  * 导出偏好 JSON
  * 一个 debug 开关占位（可选）

---

### 8. `ThemeModeSelect.vue`

职责：

* 单独抽象主题选择器
* 便于后续扩展
* 可使用 segmented 或 radio group 风格

---

## 十、Settings 页面要求

---

### `SettingsPage.vue`

请将其实现为独立偏好页面。

要求：

1. 使用 `BlankLayout` 或当前项目设置页布局
2. 页面主体为 `SettingsPanel.vue`
3. 页面本身保持轻量，不写复杂业务逻辑
4. 进入页面时初始化 settings store
5. 风格为桌面工具偏好中心，不是后台管理配置页

---

## 十一、与现有模块的联动要求

请确保 Settings 不只是静态页面，而是能影响现有模块。

---

### 联动 1：工作台默认上下文

请支持：

* `defaultContext`

行为要求：

* 应用于“无上次工作区可恢复时”的初始上下文
* 或当用户显式重置工作台时作为默认值

如果当前项目启动逻辑不便完整接管，也可以先在设置保存后提供一个简化行为，并说明。

---

### 联动 2：文章编辑区宽度与字号

当用户修改：

* `contentWidth`
* `editorFontSize`

`EditorShell.vue` 要立即变化。

---

### 联动 3：状态栏显隐

当用户修改：

* `showBottomStatusBar`

工作台底部状态栏立即显示/隐藏。

---

### 联动 4：命令面板显示分组

当用户修改：

* `commandPaletteShowGroup`

`CommandPalette.vue` 中命令 group label 显示/隐藏。

---

### 联动 5：快捷键提示显隐

当用户修改：

* `showShortcutHints`

命令面板/状态栏等处的快捷键信息按设置显示或隐藏。

---

### 联动 6：主题模式

当用户切换 light / dark：

* 应立即影响全局主题 class/data-theme
* 页面无需刷新

---

## 十二、导出与恢复默认要求

请实现两个明确操作：

### 1. Reset to Defaults

要求：

* 恢复全部默认设置
* 保存到 localStorage
* 立即应用到 UI
* 最好有一个轻量确认步骤或说明

### 2. Export Preferences

要求：

* 生成格式化 JSON 字符串
* 可：

  * 显示在 dialog 中
  * 或复制到剪贴板
  * 或触发下载文本文件
* 请选择一种最简单且清晰的第一阶段方案，并说明理由

### 推荐第一阶段方案

* 点击后弹出 Dialog，里面显示格式化 JSON，并提供 Copy 按钮
* 这是最适合当前阶段的方案

---

## 十三、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 专业、克制
3. 轻边框
4. 小圆角
5. 更像桌面知识工具偏好页，不像后台管理配置中心
6. 分区清晰
7. 表单控件紧凑但可读
8. 重点是“长期使用的个人工作台设置”，不是花哨视觉

建议视觉感觉：

* 类似桌面应用 Preferences / Settings
* 左右或上下分区清楚
* 重点在可配置性与可理解性

---

## 十四、代码结构要求

请特别遵守：

1. 不要把所有设置项写死在 `SettingsPage.vue`
2. 不要把 localStorage 逻辑直接散落在组件里
3. `settings.defaults.ts`、`settings.storage.ts`、`settings.store.ts` 分层清晰
4. `SettingsPanel.vue` 只负责页面组合
5. 各 section 负责各自分组
6. 主题应用逻辑不要散在多个组件中
7. 所有 import 必须完整
8. 不要只写伪代码或片段
9. 代码应为“第一阶段真实可运行版本”

---

## 十五、建议的实现策略

为了降低复杂度，推荐采用以下方案：

### 推荐方案

* settings store 初始化时加载 localStorage
* 任意偏好更新 -> 立即保存 -> 立即 apply
* 主题通过 `data-theme` 应用
* 布局偏好通过 settings -> workspace 的轻量同步应用
* 阅读区偏好通过 `EditorShell.vue` 的 class 切换实现
* 导出偏好通过 Dialog + copy 按钮实现

### 为什么推荐这样做

因为这套方案：

* 实现成本低
* 反馈即时
* 结构足够清晰
* 后续可无痛升级为：

  * 后端同步设置
  * 用户 profile
  * 多设备同步
  * 更多外观与行为偏好

---

## 十六、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text id="evjlwm"
src/modules/settings/types.ts
```

然后再给完整代码。

### 第四步：最后说明：

1. 如何接入当前项目
2. 用户如何在 Settings 页面修改并应用偏好
3. 当前实现相对完整配置系统的简化点有哪些

---

## 十七、补充约束

```text id="ctjlwm"
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入额外重量级表单库
- 可以使用少量 Element Plus Form / Switch / Radio / Slider / Select / Dialog / Button 等基础控件
- 当前阶段优先保证结构清晰、模块边界清晰、可运行
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 十八、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。

