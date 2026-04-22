# Changelog

本目录（`docs/ui`）记录与桌面端前端界面结构、布局与交互相关的**显著变更**，便于设计与实现对照。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本号语义可参考 [Semantic Versioning](https://semver.org/lang/zh-CN/)（与发布节奏对齐时由维护者填写）。

## [Unreleased]

### Added

- `WorkspaceNavMenu`：主导航置于 `MainLayout` 的 `AppHeader` 的 `#actions` 插槽，使用 Element Plus 下拉菜单切换 **Reading / Knowledge**（旧书签路径 `inbox` / `review` / `search` 重定向到 Reading）。
- **Knowledge 树（TreePane / TreePanel）**：`knowledge-tree` store 新增 getter `treeListEmptyMessage`，用于空列表文案区分「尚无知识库节点」与「筛选无匹配」。

### Changed

- **Knowledge 树数据源**：`useTreeStore` 的 `initialize()` 仅从 Wails `GetKnowledgeTree`（经 `knowledge.repository.fetchKnowledgeTree`）加载扁平节点；不再在请求失败或 `!isWailsAvailable()` 时回退到 `mock/tree` 演示数据。失败或非桌面环境时 `rawNodes` / `expandedNodeIds` 置空，并用 toast 提示（失败沿用接口错误信息；无 Wails 时提示需在桌面应用中加载）。成功且无根节点时展开列表清空，避免引用已失效 id。
- **TreePanel**：空状态 `AppEmpty` 的 `message` 改为绑定 `tree.treeListEmptyMessage`，替代原先单一的 “No nodes match the current filter”。
- 工作区上下文 `WorkspaceContext` 仅保留 `reading` 与 `knowledge`；从设置或持久化读出的旧值（inbox / review / search）在启动与写入时归一为上述二者之一。
- 打开复习项（`openReviewItemInWorkspace`）不再进入已移除的 Review 路由，改为按关联跳转到 Reading（文章来源）或 Knowledge（节点 / 摘录）。
- 快捷键：`Alt+1` / `Alt+2` 与 `g r` / `g k` 对应 Reading / Knowledge；移除原 Inbox、Review、Search 与全局搜索相关的快捷键绑定。
- Knowledge 工作区：`ReadPane` 在 `knowledge` 上下文改为与 Reading/Inbox 一致使用 `EditorShell`；通过统一的 `documentForEditor` 将树选中节点（摘录 / 卡片 / 其它节点描述）映射为 `editor.openDocument`，不再使用 `ExtractDetailPanel` / `CardDetailPanel` 与元数据块；该侧 `AppPane` 标题为「Editor」。
- 工作区外壳 `WorkspaceShell`：由「左侧栏 + 中间 + 右侧 + 底栏」改为「**TreePane**（中间列表/树）| 竖向分割条 | **ReadPane**（阅读/详情）」两栏主体 + 底栏。
- 原 `CenterPane.vue` 重命名为 `TreePane.vue`，根样式类名为 `tree-pane__*`；在 `AppPane` 的 `#header` 中增加顶栏工具条（Capture / Create / 工具栏折叠 / 设置），Inbox 列表区不再单独放置重复的 Capture 入口。
- 原 `RightPane.vue` 重命名为 `ReadPane.vue`，样式类名为 `read-pane__*`；搜索占位说明中的 “center pane” 改为 “tree pane”。
- `applyPreferences` 与持久化布局：工作区 Pinia 仍保留 `isLeftCollapsed`（用于 TreePane 顶栏紧凑模式），但不再从设置中推导左栏显示状态。

### Removed

- 前端演示用知识树 mock：`desktop/frontend/src/mock/tree.ts` 已删除；`src/mock/index.ts` 中不再导出 `./tree`。
- **Inbox / Review / Search** 三个独立路由页面与主导航入口；`TreePane` / `ReadPane` 中仅服务这三类工作区的 UI 分支；`modules/search` 整模块（索引、全局搜索 UI、与捕获/对象创建中的索引刷新联动）。
- Review 页专用队列 UI：`ReviewQueue`、`ReviewQueueItem`、`ReviewCard`、`ReviewMetaSection`（复习数据仍由 `review` store 等保留，卡片/摘录侧栏复习流程不受影响）。
- `LeftSidebar.vue` 及 `WorkspaceShell` 内左侧栏与首条竖向 `AppSplitter`。
- 设置中的「左栏宽度」「显示左栏」偏好项，以及 `WorkspacePreferences` / `workspace` store 中与 `showLeftPane`、`leftPaneWidth` 相关的字段与动作（如 `setShowLeftPane`、`setLeftPaneWidth`、`bumpLeftWidth`）。
- `shared/constants/layout.ts` 中仅用于已删除左栏的宽度常量。
