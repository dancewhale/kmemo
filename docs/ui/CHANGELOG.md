# Changelog

本目录（`docs/ui`）记录与桌面端前端界面结构、布局与交互相关的**显著变更**，便于设计与实现对照。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本号语义可参考 [Semantic Versioning](https://semver.org/lang/zh-CN/)（与发布节奏对齐时由维护者填写）。

## [Unreleased]

### Added

- `WorkspaceNavMenu`：主导航置于 `MainLayout` 的 `AppHeader` 的 `#actions` 插槽，使用 Element Plus 下拉菜单切换 Inbox / Reading / Knowledge / Review / Search。

### Changed

- Knowledge 工作区：`ReadPane` 在 `knowledge` 上下文改为与 Reading/Inbox 一致使用 `EditorShell`；通过统一的 `documentForEditor` 将树选中节点（摘录 / 卡片 / 其它节点描述）映射为 `editor.openDocument`，不再使用 `ExtractDetailPanel` / `CardDetailPanel` 与元数据块；该侧 `AppPane` 标题为「Editor」。
- 工作区外壳 `WorkspaceShell`：由「左侧栏 + 中间 + 右侧 + 底栏」改为「**TreePane**（中间列表/树）| 竖向分割条 | **ReadPane**（阅读/详情）」两栏主体 + 底栏。
- 原 `CenterPane.vue` 重命名为 `TreePane.vue`，根样式类名为 `tree-pane__*`；在 `AppPane` 的 `#header` 中增加顶栏工具条（Capture / Create / 工具栏折叠 / 设置），Inbox 列表区不再单独放置重复的 Capture 入口。
- 原 `RightPane.vue` 重命名为 `ReadPane.vue`，样式类名为 `read-pane__*`；搜索占位说明中的 “center pane” 改为 “tree pane”。
- `applyPreferences` 与持久化布局：工作区 Pinia 仍保留 `isLeftCollapsed`（用于 TreePane 顶栏紧凑模式），但不再从设置中推导左栏显示状态。

### Removed

- `LeftSidebar.vue` 及 `WorkspaceShell` 内左侧栏与首条竖向 `AppSplitter`。
- 设置中的「左栏宽度」「显示左栏」偏好项，以及 `WorkspacePreferences` / `workspace` store 中与 `showLeftPane`、`leftPaneWidth` 相关的字段与动作（如 `setShowLeftPane`、`setLeftPaneWidth`、`bumpLeftWidth`）。
- `shared/constants/layout.ts` 中仅用于已删除左栏的宽度常量。
