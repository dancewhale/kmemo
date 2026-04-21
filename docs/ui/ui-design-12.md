# 任务单 12：右键菜单 + 上下文操作（Context Actions）第一阶段

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **“右键菜单 + 上下文操作（Context Actions）第一阶段可运行版本”**。

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
* `polish` 第一阶段
* 三栏工作台
* 搜索、命令、复习、摘录、知识树等基础联动
* 若干共享组件和全局样式体系

本次目标是：**实现第一阶段的通用右键菜单与对象级上下文操作能力**，让用户可以对不同对象执行对应动作，并强化“桌面知识工具”的使用感。

当前阶段重点是：

* 通用 context menu 基础设施
* 不同实体的菜单项注册
* 上下文动作执行
* 与现有模块联动
* 清晰的组件边界与可扩展结构

不是复杂批量操作系统，也不是完整命令权限框架。

---

## 一、模块目标

请实现如下能力：

1. 用户在不同对象上右键时，弹出对应上下文菜单

2. 至少支持以下对象的右键菜单：

   * 知识树节点
   * 阅读列表文章项
   * 摘录详情对象或关联摘录列表项
   * 搜索结果项
   * 复习队列项

3. 不同对象显示不同的菜单项

4. 用户点击菜单项后，能够触发对应动作

5. 菜单项支持：

   * icon（可选）
   * label
   * disabled
   * divider
   * danger（可选）

6. 菜单位置在鼠标附近弹出

7. 点击外部区域、Esc、执行后关闭菜单

8. 整个系统要有一套通用的 context action 注册与执行结构

9. 为后续扩展更多对象类型、批量菜单、二级菜单预留边界

---

## 二、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现复杂权限系统
3. 不要实现多级嵌套菜单
4. 不要实现复杂批量选择系统
5. 不要实现拖放与右键菜单联动
6. 不要实现复杂剪贴板协议
7. 不要实现系统级原生菜单
8. 不要把所有动作逻辑写死在组件里
9. 不要把所有对象共用一份 if/else 超长菜单组件
10. 不要把这个任务做成“到处复制 el-dropdown”

---

## 三、本次建议实现/修改的文件

请至少实现和补全以下文件：

```text
src/modules/context-actions/
  components/
    ContextMenu.vue
    ContextMenuList.vue
    ContextMenuItem.vue
  stores/
    context-menu.store.ts
  types.ts
  services/
    context-menu.registry.ts
    context-menu.executor.ts
    context-menu.builder.ts
```

并根据需要补充或修改以下文件：

```text
src/modules/knowledge-tree/components/TreeNodeItem.vue
src/modules/reader/components/ReaderListItem.vue
src/modules/extract/components/ArticleExtractList.vue
src/modules/extract/components/ExtractDetailPanel.vue
src/modules/search/components/SearchResultItem.vue
src/modules/review/components/ReviewQueueItem.vue
src/layouts/MainLayout.vue

src/modules/knowledge-tree/stores/tree.store.ts
src/modules/reader/stores/reader.store.ts
src/modules/extract/stores/extract.store.ts
src/modules/review/stores/review.store.ts
src/modules/workspace/stores/workspace.store.ts
src/modules/command-center/services/command.executor.ts
```

如果有必要，你可以新增少量辅助文件，例如：

```text
src/shared/constants/context-actions.ts
src/shared/utils/dom.ts
src/shared/composables/useContextMenu.ts
```

但请保持克制，不要扩展成完整 action framework。

---

## 四、整体设计要求

请把这次能力拆成 4 个清晰层次：

### 1. 菜单基础设施层

负责：

* 菜单打开/关闭
* 菜单位置
* 当前上下文对象
* 点击外部关闭
* Esc 关闭

### 2. 菜单定义层

负责：

* 针对不同实体类型定义菜单项
* 不同实体的菜单内容映射
* 动态 disabled 判断
* label/icon/danger/divider 元信息

### 3. 动作执行层

负责：

* 根据 action id 执行对应逻辑
* 跨模块联动
* 尽量复用已有 store action / executor

### 4. 组件接入层

负责：

* 在树节点、列表项、搜索项等地方触发右键菜单
* 只传上下文，不直接写业务逻辑

要求边界清晰，便于后续扩展：

* 批量操作
* 二级菜单
* 工具栏菜单
* 快捷动作栏

---

## 五、类型设计要求

请先整理和定义清晰类型。

---

### `src/modules/context-actions/types.ts`

至少包含：

```ts
export type ContextEntityType =
  | 'tree-node'
  | 'article'
  | 'extract'
  | 'search-result'
  | 'review-item'

export interface ContextMenuContext {
  entityType: ContextEntityType
  entityId: string
  nodeId?: string | null
  articleId?: string | null
  extractId?: string | null
  reviewId?: string | null
  searchResultId?: string | null
  x: number
  y: number
}

export interface ContextActionItem {
  id: string
  label: string
  icon?: string
  disabled?: boolean
  danger?: boolean
  dividerBefore?: boolean
}

export interface ContextMenuState {
  isOpen: boolean
  x: number
  y: number
  context: ContextMenuContext | null
}
```

你可以在此基础上微调，但必须满足：

1. 菜单状态和菜单项分离
2. 菜单上下文能表达不同实体类型
3. 动作项能支持基本视觉元信息
4. 为未来扩展二级菜单留空间，但当前不要实现

---

## 六、右键菜单注册要求

请新增：

### `context-menu.registry.ts`

用于集中定义不同实体类型支持的菜单项。

至少覆盖这些实体：

### 1. `tree-node`

建议菜单项：

* Open
* Rename（当前阶段可 disabled 或仅占位）
* New Child Node（可做基础版）
* Create Review Item（可占位或基础版）
* Reveal in Knowledge
* Delete Node（当前阶段可 disabled 或仅 mock）

### 2. `article`

建议菜单项：

* Open Article
* Open in Reading
* Create Extract From Current Selection（仅在可行场景展示或 disabled）
* Mark as Reading
* Mark as Processed
* Reveal Related Extracts

### 3. `extract`

建议菜单项：

* Open Extract
* Back to Source Article
* Reveal in Tree
* Add to Review
* Copy Quote
* Delete Extract（当前阶段可 mock）

### 4. `search-result`

建议菜单项：

* Open Result
* Reveal in Context
* Copy Title
* Copy Reference（可选）

### 5. `review-item`

建议菜单项：

* Open Review Item
* Open Source Article
* Open Extract
* Mark Reviewed（可 mock）
* Remove From Queue（可 mock）

要求：

1. 注册表只定义菜单项元信息，不写执行逻辑
2. 支持根据上下文动态禁用某些菜单项
3. 不要把菜单定义散落在每个组件里
4. 不同对象的菜单尽量符合该对象语义

---

## 七、菜单构建要求

请新增：

### `context-menu.builder.ts`

至少实现类似：

```ts
buildContextMenuItems(context: ContextMenuContext): ContextActionItem[]
```

要求：

1. 根据当前 context 调用 registry 生成菜单项
2. 支持动态 disabled 判断
3. 必要时可读取当前 store 状态判断：

   * 当前 extract 是否存在 source article
   * 当前 search result 是否有效
   * 当前 tree node 是否允许某动作
4. 这层只负责“生成菜单”，不负责执行动作

---

## 八、动作执行要求

请新增：

### `context-menu.executor.ts`

这是本次最关键的文件之一。

至少实现：

```ts
executeContextAction(actionId: string, context: ContextMenuContext): void | Promise<void>
```

要求：

* 真正执行右键菜单动作
* 尽量复用已有 store action 和 command executor
* 统一处理跨模块跳转
* 执行成功后关闭菜单
* 必要时调用 toast 提示

至少支持这些动作：

---

### Tree Node

#### `open-node`

* 选中该节点
* 切到 `knowledge`

#### `new-child-node`

* 在当前 node 下创建一个 mock 子节点
* 自动展开父节点
* 自动选中新节点
* 标题可先用如 `New Node` / `Untitled Node`

#### `reveal-node`

* 切到 `knowledge`
* 选中该节点

#### `rename-node`

* 当前阶段可做最简版：

  * 弹出轻量 prompt / dialog 修改标题
  * 成功后更新 tree store
* 若你认为这会让本任务过重，也可先做 disabled，并说明理由

#### `delete-node`

* 当前阶段可做 mock 删除或 disabled
* 推荐第一阶段：对非根节点支持简单删除，对根节点禁用
* 删除后更新树结构

---

### Article

#### `open-article`

* `reader.store.openArticleById(...)`
* `workspace.store.setContext('reading')`

#### `mark-article-reading`

* 更新 article 状态为 `reading`

#### `mark-article-processed`

* 更新 article 状态为 `processed`

#### `reveal-related-extracts`

* 打开该文章
* 聚焦文章相关摘录区，或至少切到 reading 并打开文章
* 若“聚焦摘录区”当前过重，可只做打开文章并说明

---

### Extract

#### `open-extract`

* `extract.store.setSelectedExtract(...)`
* 若有 tree node，则同步选中树节点
* `workspace.store.setContext('knowledge')`

#### `back-to-source-article`

* 若存在 sourceArticleId，则打开原文章

#### `reveal-extract-in-tree`

* 若存在 treeNodeId，则切到 `knowledge` 并选中树节点

#### `add-extract-to-review`

* 将该 extract 生成为一个 mock review item
* 加入 review queue
* 给出成功 toast

#### `copy-extract-quote`

* 复制 quote 到剪贴板
* 成功后 toast 提示

#### `delete-extract`

* 当前阶段可做 mock 删除
* 要同步：

  * extract.store
  * tree 节点（若有关联）
* 删除后右侧应更新为空状态或切换到合理对象

---

### Search Result

#### `open-search-result`

* 复用现有 search result 打开逻辑

#### `reveal-search-result`

* 与 open 类似，但语义是“跳到原上下文”
* 可以复用 open 逻辑

#### `copy-search-title`

* 复制 title 到剪贴板

---

### Review Item

#### `open-review-item`

* 选中 review item
* 切到 review

#### `open-review-source-article`

* 若有 sourceArticleId，则打开原文章

#### `open-review-extract`

* 若有 extractId，则打开 extract

#### `remove-review-item`

* 从当前 queue 中移除该项或标记为已移除
* toast 提示成功

---

## 九、store 设计要求

---

### `context-menu.store.ts`

请实现右键菜单专用 store。

至少包含这些状态：

```ts
interface ContextMenuStoreState {
  isOpen: boolean
  x: number
  y: number
  context: ContextMenuContext | null
}
```

至少包含这些 getter：

1. `menuItems`
2. `hasContext`

至少包含这些 action：

1. `open(context: ContextMenuContext)`
2. `close()`
3. `setPosition(x: number, y: number)`
4. `execute(actionId: string)`

### 行为要求

* `open()` 时保存 context 和位置
* 菜单项通过 builder 动态生成
* 执行动作后默认关闭
* 仅负责菜单自身状态，不承担业务实体状态

---

## 十、组件要求

---

### 1. `ContextMenu.vue`

这是全局右键菜单主组件。

要求：

1. 挂载在全局布局层，例如 `MainLayout.vue`
2. 从 `context-menu.store.ts` 读取状态
3. 在 `isOpen = true` 时显示
4. 在鼠标位置附近渲染
5. 点击外部区域关闭
6. 按 `Esc` 关闭
7. 不要依赖原生浏览器 context menu

建议实现方式：

* fixed 定位
* teleport 到 body
* z-index 高于普通 pane，低于 modal/dialog 也可接受，或略高于 dialog 视需求

---

### 2. `ContextMenuList.vue`

职责：

* 渲染菜单项列表
* 使用 `ContextMenuItem.vue`
* 支持 dividerBefore

要求：

* 列表紧凑
* 可滚动（若条目多）
* 风格像桌面工具上下文菜单

---

### 3. `ContextMenuItem.vue`

职责：

* 渲染单个菜单项

要求：

1. 展示 label
2. 可显示 icon 占位
3. disabled 状态清楚
4. danger 状态清楚但克制
5. hover / active 反馈统一
6. 点击时执行动作

风格要求：

* 像桌面工具菜单项
* 不是后台 dropdown menu 风格
* 紧凑、专业

---

## 十一、组件接入要求

请将右键菜单接入以下实体组件中。

---

### 1. `TreeNodeItem.vue`

要求：

* 右键节点时打开 tree-node context menu
* 打开前先将该节点设为当前选中节点，或至少同步当前上下文
* 不要把菜单定义写在 TreeNodeItem 内

---

### 2. `ReaderListItem.vue`

要求：

* 右键文章项时打开 article context menu
* 当前右键对象应对应 article id

---

### 3. `ArticleExtractList.vue`

要求：

* 右键摘录列表项时打开 extract context menu

---

### 4. `ExtractDetailPanel.vue`

可选要求：

* 在摘录标题区或整体 panel 提供右键支持
* 若整体 panel 右键过于泛化，可只对摘录标题/主要内容区域接入

---

### 5. `SearchResultItem.vue`

要求：

* 右键搜索结果项时打开 search-result context menu

---

### 6. `ReviewQueueItem.vue`

要求：

* 右键复习项时打开 review-item context menu

---

## 十二、MainLayout 集成要求

请将 `ContextMenu.vue` 全局挂载到 `MainLayout.vue` 或应用根层。

要求：

1. 菜单作为全局 overlay 存在
2. 与命令面板、toast、dialog 层级清晰
3. 布局层不写菜单业务逻辑
4. 仅负责挂载全局 UI

---

## 十三、与现有模块的联动要求

这次最重要的是把右键动作变成“现有系统的快捷入口”。

---

### 联动 1：树节点 -> 打开 / 新建 / 删除

* 能直接操作知识树对象

### 联动 2：文章 -> 状态流转

* 可直接从列表项右键把 article 标记为 reading / processed

### 联动 3：摘录 -> review

* 可直接从 extract 上下文菜单加入复习队列

### 联动 4：摘录 / review -> 原文回跳

* 可直接跳回 source article

### 联动 5：搜索结果 -> 打开对象

* 可直接从搜索结果进入对应上下文

### 联动 6：review -> extract / article

* 复习对象可快速跳往其来源对象

---

## 十四、菜单行为要求

请确保以下基础行为：

1. 打开菜单时阻止浏览器默认右键菜单
2. 菜单若超出视口，位置要做简单修正
3. 点击菜单外关闭
4. 执行动作后关闭
5. 切换路由/切换工作区时可自动关闭
6. 连续右键不同对象时菜单能平滑更新
7. Esc 可关闭

### 位置修正要求

至少做基础处理：

* 若右侧超出屏幕，则向左偏移
* 若底部超出屏幕，则向上偏移

不要求复杂 collision engine。

---

## 十五、反馈要求

请将部分动作接入统一 toast 反馈。

至少覆盖这些动作：

* Create child node success
* Extract added to review
* Quote copied
* Node deleted
* Review item removed
* Operation unavailable / failed

要求：

* 使用你项目已有的统一 toast 封装
* 文案简洁专业
* 不要直接散落多处 Message 调用

---

## 十六、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 桌面知识工具风格
3. 轻边框
4. 小圆角
5. hover 明确但克制
6. danger 动作清楚但不刺眼
7. disabled 项可识别
8. 菜单整体像桌面工具上下文菜单，不像后台 dropdown

视觉感觉应接近：

* 快
* 紧凑
* 专业
* 上下文明确
* 操作密度高

---

## 十七、代码结构要求

请特别遵守：

1. 不要把菜单项数组写死在每个业务组件里
2. registry / builder / executor / store 四层要清晰
3. `ContextMenu.vue` 只负责 UI 与交互
4. 执行动作放在 executor 中
5. 动态可用性判断放在 builder 或 registry 辅助逻辑中
6. 所有 import 必须完整
7. 不要只写伪代码或片段
8. 代码应为“第一阶段真实可运行版本”

---

## 十八、建议的实现策略

为了降低复杂度，推荐采用以下方案：

### 推荐方案

* 全局只挂一个 `ContextMenu.vue`
* 右键对象时传入 `{ entityType, entityId, ... }`
* builder 基于 context 生成菜单
* executor 基于 actionId + context 执行动作
* 位置修正做基础版
* 删除/重命名等敏感动作做最简 mock 版

### 为什么推荐这样做

因为它已经足够支撑后续扩展：

* 更多对象类型
* 更多动作
* 二级菜单
* 工具栏下拉动作
* 批量对象菜单
* 上下文敏感动作

但当前阶段不会过度复杂。

---

## 十九、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：说明你将如何分层完成本次 context menu 系统

例如：

* state layer
* registry layer
* builder layer
* executor layer
* integration layer

### 第四步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text
src/modules/context-actions/types.ts
```

然后再给完整代码。

### 第五步：最后说明：

1. 如何接入当前项目
2. 用户如何在不同对象上使用右键菜单
3. 当前实现相对完整桌面级上下文菜单系统的简化点有哪些

---

## 二十、补充约束

```text
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入大型菜单库
- 可以使用少量 Element Plus Icon 或基础组件，但菜单视觉和行为应由项目自身控制
- 当前阶段优先保证结构清晰、可运行、风格统一
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 二十一、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。

