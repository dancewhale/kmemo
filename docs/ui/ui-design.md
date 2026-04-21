# 任务单 1：工程和主布局

你是一名资深 Vue 3 / TypeScript 工程师，请为一个基于 **Wails + Go 后端** 的桌面知识工具实现 **前端第一阶段工程骨架与主布局**。产品定位参考 **SuperMemo 风格的渐进阅读/知识工作台**，不是后台管理系统，也不是普通内容网站。

请严格按照以下要求输出和生成代码。

---

## 一、项目背景

这是一个桌面应用，后端已由 Go 实现，前端使用：

* Vue 3
* TypeScript
* Vite
* Element Plus
* Vue Router
* Pinia

当前阶段只做：

1. **工程初始化后的前端骨架**
2. **主布局**
3. **路由**
4. **基础 shared 组件**
5. **workspace 模块**
6. **mock 数据驱动的可运行界面**

暂时 **不要接真实后端**，但要为后续 Wails API 接入保留清晰边界。

---

## 二、整体目标

请实现一个 **高信息密度、桌面知识工具风格** 的主界面，具备以下结构：

* 左侧导航区
* 中间内容区
* 右侧详情/编辑区
* 底部状态栏

要求完成后项目可以运行，并且不同页面路由能展示统一工作台壳子。

---

## 三、技术与编码要求

请严格遵守：

1. 使用 **Vue 3 + TypeScript**
2. 所有 Vue 单文件组件统一使用：

```vue
<script setup lang="ts">
```

3. 使用 **Vue Router**
4. 使用 **Pinia**
5. 使用 **SCSS**
6. 使用 **Element Plus**，但仅用于基础控件，不要把整体写成后台管理系统风格
7. 使用 **createWebHashHistory**

   * 因为项目运行在 Wails 场景下，优先避免 history 路由带来的路径问题
8. 所有代码必须可维护、可扩展，避免把逻辑堆进单个组件
9. 当前阶段不接真实后端，但请预留 `src/api/wails` 目录
10. **组件层不得直接调用 Wails API**
11. 页面风格必须偏向：

* 桌面知识工具
* 专业工作台
* 紧凑
* 信息密度较高
* 克制的边框与颜色
* 非中后台风格

---

## 四、目录结构要求

请按下面的目录结构生成第一阶段代码，不需要一次性填满所有业务模块，但至少要把与“工程和主布局”有关的部分创建完整：

```text
desktop/frontend/src/
  app/
    router/
      index.ts
      routes.ts
    store/
      index.ts
    providers/
      element-plus.ts
      shortcuts.ts
    styles/
      reset.scss
      variables.scss
      theme.scss
      layout.scss
      index.scss

  api/
    wails/
      index.ts

  layouts/
    MainLayout.vue
    BlankLayout.vue

  pages/
    InboxPage.vue
    ReadingPage.vue
    KnowledgePage.vue
    ReviewPage.vue
    SearchPage.vue
    SettingsPage.vue

  modules/
    workspace/
      components/
        WorkspaceShell.vue
        LeftSidebar.vue
        CenterPane.vue
        RightPane.vue
        BottomStatusBar.vue
      stores/
        workspace.store.ts
      types.ts

  shared/
    components/
      AppPane.vue
      AppHeader.vue
      AppEmpty.vue
      AppSplitter.vue
      AppIconButton.vue
      AppLoading.vue
    composables/
      usePaneResize.ts
    constants/
      routes.ts
      layout.ts
    utils/
      storage.ts
    types/
      common.ts

  mock/
    index.ts
    tree.ts
    articles.ts
    review.ts

  App.vue
  main.ts
```

---

## 五、你这次只需要实现的范围

本次只实现以下内容：

### 1. 工程入口

* `main.ts`
* `App.vue`

### 2. 路由系统

* `app/router/routes.ts`
* `app/router/index.ts`

### 3. 应用级注册

* `app/store/index.ts`
* `app/providers/element-plus.ts`
* `app/providers/shortcuts.ts`

### 4. 全局样式

* `app/styles/reset.scss`
* `app/styles/variables.scss`
* `app/styles/theme.scss`
* `app/styles/layout.scss`
* `app/styles/index.scss`

### 5. 布局

* `layouts/MainLayout.vue`
* `layouts/BlankLayout.vue`

### 6. 页面

* `InboxPage.vue`
* `ReadingPage.vue`
* `KnowledgePage.vue`
* `ReviewPage.vue`
* `SearchPage.vue`
* `SettingsPage.vue`

### 7. workspace 模块

* `modules/workspace/types.ts`
* `modules/workspace/stores/workspace.store.ts`
* `modules/workspace/components/WorkspaceShell.vue`
* `modules/workspace/components/LeftSidebar.vue`
* `modules/workspace/components/CenterPane.vue`
* `modules/workspace/components/RightPane.vue`
* `modules/workspace/components/BottomStatusBar.vue`

### 8. shared 通用组件

* `AppPane.vue`
* `AppHeader.vue`
* `AppEmpty.vue`
* `AppSplitter.vue`
* `AppIconButton.vue`
* `AppLoading.vue`

### 9. 通用工具

* `shared/composables/usePaneResize.ts`
* `shared/constants/routes.ts`
* `shared/constants/layout.ts`
* `shared/utils/storage.ts`
* `shared/types/common.ts`

### 10. mock 数据

* `mock/tree.ts`
* `mock/articles.ts`
* `mock/review.ts`
* `mock/index.ts`

---

## 六、明确不要做的事情

这次严格不要做以下内容：

1. 不要接入真实 Wails API
2. 不要实现 he-tree
3. 不要实现 Tiptap
4. 不要实现复杂业务逻辑
5. 不要实现完整拖拽树
6. 不要实现完整复习系统
7. 不要做后台管理系统式表格页
8. 不要一次性写大量未来阶段代码
9. 不要写假大空注释，代码要真实可运行
10. 不要为了“看起来完整”而引入不必要复杂度

---

## 七、布局和交互要求

请实现一个 **主工作台布局**，具有如下视觉和交互特征：

### 总体结构

* 顶部可有一个轻量应用头部区域，也可以内嵌在主布局里
* 主区域为：

  * 左侧导航栏
  * 中间主内容区
  * 右侧详情/编辑区
* 底部为状态栏

### 左侧导航栏

至少包含这些入口：

* Inbox
* Reading
* Knowledge
* Review
* Search
* Settings

要求：

* 当前路由高亮
* 紧凑
* 专业工具感
* 支持图标 + 文本
* 不要做成后台菜单样式

### 中间内容区

根据当前页面上下文显示不同的 mock 内容，例如：

* Reading 显示文章列表
* Knowledge 显示节点列表占位
* Review 显示复习列表占位

### 右侧详情区

根据当前上下文显示详情占位，例如：

* 当前文章摘要
* 当前节点信息
* 当前复习项信息

### 底部状态栏

至少显示：

* 当前模式
* 当前选中项 id 或标题
* mock sync 状态
* 一两个快捷键提示

### 面板尺寸

* 左右面板宽度要可配置
* 最好支持一个基础版 split 组件
* 第一阶段可以先实现基础拖拽，不要求特别复杂
* 面板尺寸可以保存在 localStorage

---

## 八、页面与路由要求

请使用如下路由结构：

* `/` 重定向到 `/reading`
* `/inbox`
* `/reading`
* `/knowledge`
* `/review`
* `/search`
* `/settings`

要求：

1. `/inbox`、`/reading`、`/knowledge`、`/review`、`/search` 使用 `MainLayout`
2. `/settings` 使用 `BlankLayout`
3. 主工作区页面通过 `WorkspaceShell` 渲染
4. 页面组件本身不要写太复杂，只负责传递当前上下文

---

## 九、状态设计要求

请为 `workspace` 模块建立清晰的 Pinia store。

### `workspace.store.ts` 至少包含以下状态：

```ts
type WorkspaceContext = 'inbox' | 'reading' | 'knowledge' | 'review' | 'search'

interface WorkspaceState {
  currentContext: WorkspaceContext
  leftPaneWidth: number
  rightPaneWidth: number
  bottomPaneHeight: number
  isLeftCollapsed: boolean
  isRightCollapsed: boolean
  selectedArticleId: string | null
  selectedNodeId: string | null
  selectedReviewId: string | null
  syncStatus: 'idle' | 'syncing' | 'saved'
}
```

还需要包括：

* 设置当前 context 的 action
* 设置左右面板宽度的 action
* 设置当前选中 article/node/review 的 action
* 从 localStorage 恢复布局状态
* 持久化布局状态

请保持 store 职责清晰，不要混入未来 editor/tree/review 的复杂业务逻辑。

---

## 十、shared 组件要求

请实现以下基础组件，并确保它们足够通用。

### 1. `AppPane.vue`

作用：

* 通用面板壳
* 支持 title
* 支持 header slot / default slot / footer slot
* 支持 bordered、scrollable、padding 之类的 props

### 2. `AppHeader.vue`

作用：

* 通用标题栏
* 支持 title、subtitle
* 支持 actions slot

### 3. `AppEmpty.vue`

作用：

* 空状态展示

### 4. `AppSplitter.vue`

作用：

* 作为 pane 间的分隔条
* 支持横向/纵向
* 基础拖拽事件即可

### 5. `AppIconButton.vue`

作用：

* 轻量图标按钮封装

### 6. `AppLoading.vue`

作用：

* 通用加载状态

这些 shared 组件应偏简洁、可复用，不要耦合具体业务。

---

## 十一、mock 数据要求

请提供基础 mock 数据，保证界面能“像真实产品那样联动”。

### `mock/articles.ts`

至少提供 8 条文章数据，字段建议包括：

* id
* title
* summary
* sourceType
* status
* updatedAt

### `mock/tree.ts`

至少提供 8 条知识节点，字段建议包括：

* id
* parentId
* title
* type

### `mock/review.ts`

至少提供 5 条复习项，字段建议包括：

* id
* title
* prompt
* dueAt
* status

### `mock/index.ts`

统一导出。

---

## 十二、页面组件职责要求

### `ReadingPage.vue`

* 使用 `WorkspaceShell`
* context 为 `reading`

### `InboxPage.vue`

* 使用 `WorkspaceShell`
* context 为 `inbox`

### `KnowledgePage.vue`

* 使用 `WorkspaceShell`
* context 为 `knowledge`

### `ReviewPage.vue`

* 使用 `WorkspaceShell`
* context 为 `review`

### `SearchPage.vue`

* 使用 `WorkspaceShell`
* context 为 `search`

### `SettingsPage.vue`

* 使用 `BlankLayout` 下的简单内容页
* 不进入工作台布局

页面组件尽量薄，只做页面级上下文传递。

---

## 十三、WorkspaceShell 具体要求

`WorkspaceShell.vue` 是这次最核心的组件，请认真设计。

它需要：

1. 接收一个 `context` prop，类型为：

   * `'inbox' | 'reading' | 'knowledge' | 'review' | 'search'`

2. 负责组合：

   * `LeftSidebar`
   * `CenterPane`
   * `RightPane`
   * `BottomStatusBar`

3. 根据 store 中的 pane 宽度控制布局

4. 使用 `AppSplitter` 实现左右区域分隔

5. 页面加载时根据 `context` 设置当前 workspace context

6. 不要在这里写未来 editor/tree 的复杂逻辑

7. 当前阶段通过 mock 数据驱动中间和右侧区域展示

---

## 十四、CenterPane / RightPane 内容要求

### `CenterPane.vue`

根据 `context` 展示不同内容：

* `reading`：文章列表
* `inbox`：待处理文章列表
* `knowledge`：知识节点列表占位
* `review`：复习队列列表
* `search`：搜索占位结果

要求：

* 列表是“桌面知识工具”风格，不是后台表格
* 当前选中项高亮
* 点击列表项更新 workspace store 里的 selected 状态

### `RightPane.vue`

根据 `context` 展示不同详情占位：

* `reading`：当前文章详情 / 摘要
* `inbox`：当前待处理内容摘要
* `knowledge`：当前节点详情
* `review`：当前复习项详情
* `search`：搜索说明或结果详情占位

---

## 十五、样式要求

请为整个项目提供一套 **简洁、紧凑、专业工具风格** 的基础样式。

要求：

1. 使用 SCSS 变量管理：

   * 背景色
   * 面板色
   * 边框色
   * 主色
   * 文本色
   * 次级文本色
   * 间距
   * 圆角

2. 风格要求：

   * 高信息密度
   * 轻边框
   * 小圆角
   * 紧凑列表
   * 适合桌面知识工具
   * 不要大卡片营销站风格
   * 不要后台管理系统蓝白大表格风格

3. 支持这些基础视觉层次：

   * app 背景
   * pane 背景
   * hover
   * active
   * subtle border

---

## 十六、输出顺序要求

请严格按下面顺序输出：

### 第一步：先输出目录树

只输出本次会创建的目录树。

### 第二步：输出文件职责说明

简要说明每个文件的职责。

### 第三步：开始生成代码

按文件逐个输出代码，确保：

* 文件路径明确
* 代码完整
* 彼此可配合运行

### 第四步：最后补充启动说明

包括：

* 需要安装哪些依赖
* 如何启动项目

---

## 十七、代码质量要求

请特别注意以下几点：

1. 命名清晰
2. 不要把多个职责混进一个文件
3. 类型定义明确
4. 尽量避免魔法字符串散落
5. 组件 props / emits 合理
6. 避免过度设计
7. 但也不要为了省事把结构写乱
8. 代码必须是“第一阶段真实可用骨架”，不是伪代码

---

## 十八、依赖假设

默认项目已经安装：

* vue
* vue-router
* pinia
* element-plus
* @element-plus/icons-vue
* sass

如果你认为还需要少量轻依赖，可以说明，但尽量不要额外增加复杂依赖。

---

## 十九、你生成代码时的额外限制

1. 不要把 mock 数据直接写死在组件内部，必须放到 `src/mock`
2. 不要把路由常量写死在多个地方，尽量抽到 constants
3. 不要在 page 组件里实现复杂布局逻辑
4. 不要把布局状态存在组件内部，应该放进 workspace store
5. 不要提前实现未来 editor/tree/review 的复杂模块细节
6. 第一阶段的目标是“骨架稳定、界面能跑、结构正确”

---

## 二十、现在开始执行

请从“目录树”开始，然后按要求继续输出完整代码。

---

你也可以把下面这段作为一个更短的“附加约束”，拼到上面末尾，进一步减少 AI 跑偏：

```text
补充约束：
- 所有文件都要给出完整代码，不要只给片段
- 不要省略 import
- 不要留 TODO 代替实现
- 允许用 mock 数据完成联动
- 所有页面都必须能通过路由正常访问
- 生成代码时优先保证结构稳定和组件边界清晰
- UI 风格偏 SuperMemo/知识工作台，不偏中后台
```
