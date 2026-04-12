# Vue 渐进阅读前端代码框架设计文档

## 1. 背景与目标

当前 `desktop/frontend` 仍是原生 JS + Vite 的占位骨架，只用于验证 Wails 与 Go Desktop API 调用链是否打通，尚未形成可持续演进的前端工程结构。用户已明确决定放弃这套前端实现，直接切换到 Vue，建设一个类似 SuperMemo 的渐进阅读桌面应用前端。

本轮目标不是立即实现全部页面，而是先产出一份可实践、可分阶段落地的前端代码框架设计文档，指导后续在 `desktop/frontend` 下重建 Vue 工程。

本设计文档遵循以下约束：

- 保留 Wails 宿主：`desktop/main.go`、`desktop/app.go`、`desktop/wails.json`
- 保留 Go 暴露的 Desktop API 作为前端唯一后端入口：`internal/app/*.go`
- 前端工程目录仍保持在 `desktop/frontend`
- 构建产物仍输出到 `desktop/frontend/dist`，以满足 Wails embed 约束
- 设计以“渐进阅读工作台”为核心，而不是普通后台 CRUD 页面

## 2. 设计目标

本轮 Vue 前端框架应满足以下目标：

1. 支撑桌面端主工作台
   - 知识树导航
   - 卡片列表 / 卡片树浏览
   - 卡片详情查看
   - 卡片创建与编辑
   - 复习入口

2. 支撑后续渐进阅读扩展
   - 阅读区 / 摘录流
   - 复习工作台
   - 系统设置 / FSRS 参数管理
   - importing / queue / 搜索增强等后续模块

3. 遵循社区最佳实践
   - Vue 3 + Vite + TypeScript
   - Pinia 做状态管理
   - Vue Router 做工作区路由切分
   - 统一 bridge/service 层封装 Wails API
   - 明确页面、功能、领域实体、共享模块边界

4. 可分阶段落地
   - 先完成框架和主工作台闭环
   - 再逐步补 review、system、reading 等模块

## 3. 技术选型

## 3.1 推荐技术栈

建议采用：

- Vue 3
- Vite
- TypeScript
- Pinia
- Vue Router
- Sass 或 scoped CSS + 全局 design tokens
- Vitest（后续需要时）
- Vue Test Utils（后续需要时）

## 3.2 选型理由

### Vue 3
适合构建信息密度较高的桌面工作台，组合式 API 更利于拆分知识树、卡片区、详情区、复习区等复杂交互逻辑。

### TypeScript
当前系统的 Go Desktop API 已具备明确 DTO 结构，使用 TypeScript 更适合在前端同步建立 `CardDTO / KnowledgeTreeNode / TagDTO / SRSDTO / FSRSParameterDTO` 等类型，减少跨层误用。

### Pinia
比直接事件总线或手写共享状态更适合中型桌面前端。能够清晰管理：
- 当前工作区
- 选中的 knowledge/card
- 卡片列表过滤条件
- review session 状态
- system settings 状态

### Vue Router
当前第一阶段虽然以主工作台为主，但代码结构应提前支持：
- workspace
- review
- settings
等独立工作区，不建议把所有内容都堆在单页状态切换里。

### Wails Bridge 封装
前端不能在组件里散落 `window.go.main.App.*` 调用。应统一通过 `src/bridge/` 封装成 service，保证：
- 可测试
- 可替换
- 可集中处理错误与数据转换

## 4. 保留与替换边界

## 4.1 保留部分

以下边界继续保留，不作为本轮前端重构对象：

- `desktop/main.go`
- `desktop/app.go`
- `desktop/wails.json`
- `internal/app/card.go`
- `internal/app/knowledge.go`
- `internal/app/tag_search.go`
- `internal/app/review.go`
- `internal/app/system.go`
- `docs/wails-api-design.md`

原因：
- 这些文件已经定义了 Wails 宿主和 Go 暴露接口边界
- 当前前端重构的重点是“前端工程与信息架构”，不是立即重写 Go 宿主层

## 4.2 替换部分

以下前端骨架应视为直接替换对象：

- `desktop/frontend/src/main.js`
- `desktop/frontend/index.html`（保留壳，但入口改为 Vue）
- `desktop/frontend/package.json`
- `desktop/frontend/vite.config.js`
- `desktop/frontend/src/` 下的当前占位实现

原因：
- 当前前端没有成熟可继承的页面结构、组件体系、样式体系与状态设计
- 继续在原生 JS 壳上叠加复杂界面会形成后续维护负担

## 5. 前端工程结构设计

建议在 `desktop/frontend/src` 下采用以下结构：

```text
src/
  main.ts
  App.vue
  router/
    index.ts
  stores/
    app.ts
    workspace.ts
    card.ts
    review.ts
    settings.ts
  layouts/
    AppShell.vue
    WorkspaceLayout.vue
    ReviewLayout.vue
    SettingsLayout.vue
  pages/
    workspace/
      WorkspacePage.vue
    review/
      ReviewPage.vue
    settings/
      SettingsPage.vue
  features/
    knowledge-tree/
    card-list/
    card-tree/
    card-detail/
    card-editor/
    review-session/
    fsrs-settings/
    tag-filter/
  entities/
    knowledge/
    card/
    tag/
    review/
    fsrs/
  bridge/
    appBridge.ts
    knowledgeService.ts
    cardService.ts
    reviewService.ts
    tagService.ts
    systemService.ts
  shared/
    components/
    composables/
    utils/
    constants/
  styles/
    tokens.scss
    reset.scss
    theme.scss
    layout.scss
  types/
    dto.ts
    workspace.ts
```

## 5.1 目录职责

### `pages/`
负责路由级场景，不承担复杂数据调用细节。

### `layouts/`
负责工作区布局骨架，例如三栏主工作台、review 工作区、settings 工作区。

### `features/`
负责面向用户动作流的功能模块，例如：
- 浏览知识树
- 选择卡片
- 新建卡片
- 提交复习
- 更新 FSRS 参数

### `entities/`
负责领域实体展示与复用组件，例如：
- CardListItem
- KnowledgeTreeNodeView
- TagBadge
- ReviewLogList
- FSRSParameterForm

### `stores/`
负责全局或跨模块状态。

### `bridge/`
负责封装所有 Wails API 调用，禁止业务组件直接读取 `window.go.main.App`。

### `shared/`
放基础组件、工具函数、常量与 composables。

## 6. 路由与工作区设计

## 6.1 路由结构

建议第一版采用轻量路由：

```text
/
  -> /workspace
/workspace
/review
/settings
```

说明：
- `/workspace` 是主入口，承载知识树、卡片列表、详情与创建流
- `/review` 承载 Due cards 与复习交互
- `/settings` 承载 FSRS 参数与后续系统设置

第一阶段虽然主力是 `/workspace`，但路由结构应一开始就预留，避免后续从单页强拆。

## 6.2 主工作台布局

主工作台建议采用三栏布局：

1. 左栏：Knowledge Sidebar
2. 中栏：Card Collection Area
3. 右栏：Card Detail / Editor Panel

顶部可增加轻量工具栏：
- 新建卡片
- 标签筛选
- 搜索输入
- 到期卡片入口
- 工作区切换入口

## 7. 状态管理设计

## 7.1 Store 划分建议

### `appStore`
保存应用级状态：
- 当前路由工作区
- 全局 loading
- 全局错误提示
- Wails 环境状态

### `workspaceStore`
保存主工作台状态：
- `selectedKnowledgeId`
- `selectedCardId`
- `workspaceMode`（browse/create/edit）
- pane 折叠状态

### `cardStore`
保存卡片相关状态：
- 卡片列表
- 总数
- 当前过滤条件
- 当前卡片详情
- 当前编辑中的卡片草稿

### `reviewStore`
保存复习状态：
- due cards
- 当前 review card
- review history
- review statistics

### `settingsStore`
保存系统设置状态：
- FSRS 参数列表
- 默认参数详情
- 更新中的表单状态

## 7.2 数据流原则

采用以下原则：

1. 页面负责触发初始化
2. store 负责状态与动作编排
3. bridge/service 负责 Desktop API 调用
4. entities/features 只消费 store 或显式 props/emits

避免：
- 组件直接调用 `window.go.main.App`
- 页面内重复拼接 DTO
- 跨组件事件层层透传

## 8. Wails Bridge 设计

## 8.1 统一入口

建议建立统一 bridge 封装：

- `appBridge.ts` 负责读取 `window.go.main.App`
- 各业务 service 再基于该 bridge 封装具体方法

例如：
- `knowledgeService.ts`
- `cardService.ts`
- `reviewService.ts`
- `tagService.ts`
- `systemService.ts`

## 8.2 Bridge 层职责

1. 统一调用 Desktop API
2. 做运行环境检查
3. 做 DTO 类型映射
4. 集中处理错误包装
5. 屏蔽 Vue 组件对底层绑定对象的直接依赖

## 8.3 现有可直接复用的 Desktop API

### Knowledge
- `GetKnowledgeTree`
- `ListKnowledge`
- `CreateKnowledge`
- `UpdateKnowledge`
- `MoveKnowledge`
- `DeleteKnowledge`

### Card
- `CreateCard`
- `GetCard`
- `ListCards`
- `UpdateCard`
- `DeleteCard`
- `GetCardTags`
- `AddCardTags`
- `RemoveCardTags`
- `SuspendCard`
- `ResumeCard`

### Tag / Search
- `ListTags`
- `SearchCardsByTags`

### Review
- `GetDueCards`
- `SubmitReview`
- `UndoLastReview`
- `GetReviewHistory`
- `GetSRSStatistics`
- `GetReviewStats`

### System
- `ListFSRSParameters`
- `GetDefaultFSRSParameter`
- `UpdateDefaultFSRSParameter`

## 9. 信息架构设计

本项目的前端不应被设计成普通表格管理后台，而应按“渐进阅读软件”组织。

## 9.1 第一阶段信息架构

### Workspace
面向知识浏览与卡片操作：
- 知识树导航
- 卡片列表
- 卡片树
- 卡片详情
- 创建/编辑面板

### Review Workspace
面向复习：
- due cards 列表
- 当前复习卡片
- review controls
- review history
- SRS statistics

### Settings Workspace
面向系统与算法参数：
- FSRS 默认参数查看/更新
- 后续系统设置扩展位

## 9.2 第二阶段预留架构

预留但不在本轮展开实现：
- 阅读器 / source material workspace
- 摘录与卡片生成流
- queue workspace
- importing workspace
- 增强搜索 workspace

## 10. 已知后端接口缺口

虽然当前 Desktop API 已足以支撑第一阶段主工作台，但仍存在几个会影响 Vue 前端实现质量的关键缺口。

## 10.1 渐进阅读完整流仍缺 source/import 边界

如果未来要真正接近 SuperMemo 式渐进阅读，前端最终还需要：
- source/document 视图
- 阅读摘录流程
- importing job 视图

但这些不进入当前 Vue 框架第一阶段实现，只在文档中保留扩展位。

## 11. 分阶段落地方案

## Phase 0：设计文档阶段

目标：
- 完成 `docs/ui` 下的 Vue 前端代码框架设计文档
- 固化工程结构、路由、store、bridge 与信息架构
- 明确必须先补的 Desktop API 缺口

交付物：
- 本文档

## Phase 1：Vue 工程替换前端骨架

目标：
- 将 `desktop/frontend` 从原生 JS 改为 Vue 3 + Vite + TS
- 保持 `desktop/wails.json` 中：
  - `frontend:install`
  - `frontend:build`
  - `frontend:dev:watcher`
  仍然成立
- 保持输出仍进入 `desktop/frontend/dist`

核心内容：
- 新建 `main.ts`
- 新建 `App.vue`
- 接入 Router
- 接入 Pinia
- 建立基础 layout
- 建立 bridge 层
- 建立全局样式与主题 token

## Phase 2：Workspace 主闭环

目标：
- 打通知识树 -> 卡片列表 -> 卡片详情 -> 新建卡片 闭环

调用链：
- `GetKnowledgeTree(nil)`
- `ListCards(filters)`
- `GetCard(id)`
- `ListTags()`
- `CreateCard(req)`

交付结果：
- 左栏可浏览 knowledge tree
- 中栏可展示卡片列表
- 右栏可显示卡片元数据
- 右栏可打开创建面板并完成创建

## Phase 3：Review 闭环

目标：
- 建立独立 Review Workspace
- 打通到期卡片加载、提交复习、查看历史与统计

调用链：
- `GetDueCards`
- `SubmitReview`
- `UndoLastReview`
- `GetReviewHistory`
- `GetSRSStatistics`

交付结果：
- review 页面可完成最小复习流

## Phase 4：System / FSRS 闭环

目标：
- 建立 settings 页面
- 打通 FSRS 参数查看与更新

调用链：
- `ListFSRSParameters`
- `GetDefaultFSRSParameter`
- `UpdateDefaultFSRSParameter`

交付结果：
- settings 页面可更新默认 FSRS 参数

## Phase 5：增强阶段

目标：
- 补 card tree
- 补卡片 HTML 详情/编辑器
- 补阅读/摘录/source workspace
- 再决定是否进入 importing / queue / 搜索增强

## 12. 验证方案

## 12.1 文档阶段

验证项：
- 文档位于 `docs/ui`
- 文档明确技术栈、目录结构、路由、store、bridge、信息架构、接口映射与缺口
- 文档步骤可直接指导实现

## 12.2 框架阶段

验证项：
- `npm install` 成功
- `npm run build` 成功
- Wails 仍能加载 `frontend/dist`
- `task run:wails` 下前端能正常启动

## 12.3 功能阶段

最小闭环验证：

1. 进入 `/workspace`
2. 左侧知识树成功加载
3. 点击 knowledge 后中栏显示卡片列表
4. 点击卡片后右栏展示详情
5. 创建卡片后列表自动刷新
6. `/review` 可加载 due cards 并提交复习
7. `/settings` 可查看与更新默认 FSRS 参数

## 13. 本轮建议结论

本轮最合适的实施顺序是：

1. 在 `docs/ui` 先产出 Vue 前端代码框架设计文档
2. 保留 Wails 宿主与现有 Desktop API 边界
3. 用 Vue 3 + Vite + TypeScript + Pinia + Vue Router 重建 `desktop/frontend`
4. 第一阶段优先打通 Workspace 主闭环
5. 第二阶段再补 Review 与 System
6. card tree、HTML 正文读取、reading/importing 等能力作为后续阶段扩展

这样既能遵循 Vue 社区项目的最佳实践，也能保持与当前 Wails/Desktop 后端接口的低摩擦衔接。
