# 任务单 5：选区创建摘录的前端交互闭环

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **“选区创建摘录”的前端交互闭环第一阶段**。

这个项目是一个基于 **Wails + Go 后端** 的桌面知识工具，产品定位参考 **SuperMemo 风格的渐进阅读 / 知识工作台**。
当前项目已经有：

* `workspace` 模块
* `reader` 模块
* `knowledge-tree` 模块
* `editor` 模块壳
* `RightPane` 中已经能打开文章内容进行基础编辑/预览
* `tree.store.ts`
* `reader.store.ts`
* `editor.store.ts`
* `mock/` 数据
* 基础 shared 组件

本次目标是：**打通“从文章选区创建摘录”的前端闭环**。
注意：这次仍然 **不接真实后端**，只做 **前端交互流 + mock 数据变更 + UI 联动**，但要为未来真实后端接入留出清晰边界。

---

## 一、模块目标

请实现这样一条完整交互链：

1. 用户在右侧文章编辑/阅读区中选中文本
2. 出现轻量的选区操作入口
3. 用户点击“Create Extract”
4. 弹出一个摘录确认界面
5. 用户可以确认以下信息：

   * 摘录文本
   * 摘录标题（可自动生成初始值）
   * 目标父节点（挂到知识树中的哪个节点下）
   * 可选 note
6. 用户确认后，前端创建一个新的摘录对象
7. 新摘录同步进入知识树 mock 数据
8. 知识树能够显示这个新摘录节点
9. 右侧可切换展示新摘录的内容
10. 整个过程在 UI 上形成清晰闭环

当前阶段重点是 **交互流与状态流**，不是复杂业务规则。

---

## 二、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现真正的富文本高亮落点持久化
3. 不要实现复杂 annotation 系统
4. 不要实现 cloze 删除
5. 不要实现 QA 卡片自动生成
6. 不要实现复杂右键菜单
7. 不要实现复杂浮动气泡定位库
8. 不要实现真实文档 anchor 恢复算法
9. 不要实现复杂树拖拽挂载流程
10. 不要把所有逻辑塞到一个组件里

---

## 三、本次建议实现/修改的文件

请至少实现和补全以下文件：

```text id="uvyh5y"
src/modules/editor/
  components/
    ExtractToolbar.vue
    ExtractCreateDialog.vue
  stores/
    editor.store.ts
  types.ts

src/modules/knowledge-tree/
  stores/
    tree.store.ts
  types.ts

src/modules/extract/
  stores/
    extract.store.ts
  types.ts
  services/
    extract.mapper.ts
```

并根据需要补充或修改以下文件：

```text id="m2ep0n"
src/modules/editor/components/TiptapEditor.vue
src/modules/editor/components/EditorShell.vue
src/modules/workspace/components/RightPane.vue
src/mock/tree.ts
src/mock/articles.ts
src/mock/index.ts
```

如有必要，你可以新增极少量辅助文件，例如：

* `shared/composables/useTextSelection.ts`
* `shared/utils/selection.ts`

但请保持克制，不要扩展过大。

---

## 四、整体设计要求

请把这次功能拆成 4 个清晰层次：

### 1. 编辑器选区层

负责：

* 监听当前文本选区
* 提供选区信息
* 触发摘录创建入口

### 2. 摘录创建 UI 层

负责：

* 弹出确认界面
* 展示摘录文本
* 编辑标题/备注
* 选择目标父节点

### 3. 摘录状态层

负责：

* 管理当前待创建摘录
* 管理已创建摘录的 mock 数据
* 创建摘录对象

### 4. 知识树联动层

负责：

* 将新摘录挂到知识树指定父节点下
* 刷新树显示
* 可选地选中新建摘录节点

要求层次清晰，便于后续接真实后端。

---

## 五、类型设计要求

请先定义清晰类型。

---

### `src/modules/extract/types.ts`

至少包含：

```ts id="953bfx"
export interface TextSelectionRange {
  from: number
  to: number
  text: string
}

export interface PendingExtractPayload {
  sourceArticleId: string
  sourceArticleTitle: string
  quote: string
  note: string
  title: string
  parentNodeId: string | null
  selection: TextSelectionRange
}

export interface ExtractItem {
  id: string
  sourceArticleId: string
  sourceArticleTitle: string
  title: string
  quote: string
  note: string
  parentNodeId: string | null
  createdAt: string
  updatedAt: string
}
```

---

### 对 `editor/types.ts` 的补充

请补充：

```ts id="tbo4f1"
export interface EditorSelection {
  from: number
  to: number
  text: string
}

export interface EditorDocument {
  id: string
  title: string
  content: string
  contentType: 'article' | 'node' | 'extract'
  updatedAt: string
  sourceUrl?: string
  tags?: string[]
}
```

要求：

1. 选区类型与摘录类型边界清晰
2. `PendingExtractPayload` 用于 UI 确认阶段
3. `ExtractItem` 用于已创建摘录
4. 设计要便于未来接真实后端 DTO

---

## 六、store 设计要求

---

### 1. `extract.store.ts`

请新增一个独立的摘录 store，不要把摘录创建状态硬塞进 `editor.store.ts` 或 `tree.store.ts`。

至少包含以下状态：

```ts id="s6vu5p"
interface ExtractState {
  items: ExtractItem[]
  pending: PendingExtractPayload | null
  isCreateDialogOpen: boolean
  creating: boolean
}
```

至少包含这些 getter：

1. `pendingQuote`
2. `hasPendingExtract`
3. `extractCount`

至少包含这些 action：

1. `openCreateDialog(payload: PendingExtractPayload)`
2. `closeCreateDialog()`
3. `updatePendingField<K extends keyof PendingExtractPayload>(key: K, value: PendingExtractPayload[K])`
4. `createExtract()`
5. `setItems(items: ExtractItem[])`
6. `clearPending()`

行为要求：

* `openCreateDialog` 时写入 pending，并打开 dialog
* `createExtract` 时基于 pending 生成一个新的 `ExtractItem`
* 当前阶段 `createExtract` 可直接在前端创建 id / 时间戳
* 创建成功后关闭 dialog
* 创建成功后应触发知识树联动
* store 结构要便于未来把 `createExtract()` 替换成真实 API 调用

---

### 2. 对 `editor.store.ts` 的要求

请在现有 editor store 基础上，补充：

* `selection: EditorSelection | null`
* `setSelection(selection: EditorSelection | null)`
* `clearSelection()`

要求：

* 只负责保存当前编辑器选区
* 不要把摘录创建逻辑放进去
* 编辑器 store 和摘录 store 分工清晰

---

### 3. 对 `tree.store.ts` 的要求

请补充一个用于新增摘录节点的 action，例如：

```ts id="mng0is"
addExtractNode(payload: {
  parentNodeId: string | null
  title: string
  sourceArticleId: string
})
```

行为要求：

* 创建一个新的知识树节点
* 节点类型为 `extract`
* 挂到指定父节点下
* 自动刷新树可见节点
* 可选：自动展开父节点
* 可选：自动选中新节点

请不要在 tree store 里存整份 extract 对象正文，只存树节点所需字段。

---

## 七、`extract.mapper.ts` 要求

请新增 `src/modules/extract/services/extract.mapper.ts`，至少包含：

1. `buildPendingExtractFromSelection(...)`
2. `buildTreeNodeFromExtract(...)`

### 作用

* 把编辑器选区映射为待创建摘录 payload
* 把已创建摘录映射为树节点所需结构

要求：

* 不要把映射逻辑散落在组件里
* 保持函数尽量纯
* 便于未来接真实 API 和 DTO

---

## 八、编辑器选区层要求

---

### 1. `TiptapEditor.vue`

虽然当前仍然可能是简化版 textarea/contenteditable，但请实现基础选区能力。

要求：

1. 能检测用户当前选中的文本
2. 当选区非空时，向外同步 `EditorSelection`
3. 当选区为空时，清空 selection
4. 提供稳定事件给上层组件
5. 不要求实现复杂浮动定位
6. 只要能在当前编辑区域内完成“选中文本 -> 显示操作入口”即可

建议支持：

* `v-model`
* `selection-change` 事件

选区数据至少包括：

* `from`
* `to`
* `text`

当前阶段 `from` 和 `to` 可以基于 textarea / selectionStart / selectionEnd 获得。
不要求达到真正 ProseMirror 级精度。

---

### 2. `ExtractToolbar.vue`

请新增一个轻量选区工具栏组件。

要求：

1. 只在当前存在非空选区时显示
2. 至少有一个主按钮：

   * `Create Extract`
3. 可再加一两个占位按钮，例如：

   * Highlight
   * Note
4. 当前阶段只有 `Create Extract` 需要真正联动
5. UI 要轻量、紧凑、像知识工具里的浮动操作条
6. 不需要做特别复杂的绝对定位动画

建议：

* 可以将其显示在编辑器区域上方或选区下方附近
* 当前阶段定位可以适度简化，只要交互清楚即可

---

## 九、摘录创建 UI 层要求

---

### `ExtractCreateDialog.vue`

请新增一个摘录创建确认弹窗组件。

要求：

1. 使用 Element Plus `Dialog` 可以接受
2. 从 `extract.store.ts` 读取 pending 数据
3. 支持编辑以下字段：

   * `title`
   * `quote`（可只读或允许轻微修改，请说明你的选择）
   * `note`
   * `parentNodeId`
4. `parentNodeId` 选择器可先使用简单下拉框
5. 下拉选项来自知识树中的可选父节点
6. 至少提供：

   * Cancel
   * Create Extract
7. 创建时按钮显示 loading 状态
8. 创建成功后自动关闭

标题默认生成规则建议：

* 从选中文本截取前 20~40 个字符
* 或 `Extract from {article title}`
* 规则清晰即可

注意：

* 这个弹窗是“确认创建摘录”
* 不是复杂表单系统
* 风格保持紧凑、专业

---

## 十、Workspace 与右侧集成要求

请将这套交互接到当前 `reading` 场景的右侧编辑区。

### 对 `EditorShell.vue` 的要求

请让 `EditorShell.vue`：

1. 从 `editor.store.ts` 读取当前文档和选区
2. 渲染 `EditorToolbar.vue`
3. 渲染 `TiptapEditor.vue`
4. 在有选区时渲染 `ExtractToolbar.vue`
5. 点击 `Create Extract` 时，触发 `extract.store.openCreateDialog(...)`
6. 同时挂载 `ExtractCreateDialog.vue`

---

### 对 `RightPane.vue` 的要求

当 `context === 'reading'` 时：

* 右侧应继续以 `EditorShell.vue` 为主
* 不要把摘录创建对话框逻辑写死在 `RightPane.vue`
* `RightPane.vue` 只负责承载 editor 模块

---

## 十一、知识树联动要求

这是本次最重要的闭环之一。

当用户在弹窗中点击 **Create Extract** 后，必须至少完成这些动作：

1. 创建新的 `ExtractItem`
2. 将其加入 `extract.store.items`
3. 将其转换成知识树节点
4. 调用 `tree.store.addExtractNode(...)`
5. 在知识树中可见
6. 父节点自动展开
7. 新建节点最好自动选中

### 当前阶段允许的简化

* 不要求树中存储摘录全文
* 树中只保留：

  * id
  * title
  * parentId
  * type = extract

### 额外建议

* 创建成功后，可以选择：

  * 继续停留在原文章编辑区
  * 或切换到新摘录节点详情
* 请你自行选择一种更合理且更简单的第一阶段实现，并说明理由

---

## 十二、mock 数据要求

请根据需要补充 mock 数据：

### `mock/articles.ts`

确保 article 具有足够长的 `content`，便于真实选区操作。

### `mock/tree.ts`

确保存在一些可选父节点，例如：

* Reading Inbox
* Knowledge System
* Memory
* Reading Workflow

这样用户在创建摘录时能选择挂载位置。

### `mock/index.ts`

统一导出。

---

## 十三、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 紧凑、专业
3. 轻边框
4. 小圆角
5. 克制的高亮和 hover
6. 更像知识工具，不像后台管理系统
7. `ExtractToolbar` 像轻量工作条，不像大表单
8. `ExtractCreateDialog` 像确认面板，不像繁琐配置弹窗

整体感觉应接近：

* 阅读中快速摘录
* 操作路径短
* 心智负担低
* 工具味强

---

## 十四、代码结构要求

请特别遵守：

1. 不要把摘录创建逻辑全塞进 `EditorShell.vue`
2. 不要把选区逻辑全塞进 `TiptapEditor.vue`
3. 不要把树联动逻辑散落在组件各处
4. 映射逻辑放进 `extract.mapper.ts`
5. store 职责清晰：

   * editor store：文档与选区
   * extract store：摘录创建与摘录列表
   * tree store：知识树节点
6. 所有 import 必须完整
7. 不要只写伪代码或片段
8. 代码应为“第一阶段真实可运行版本”

---

## 十五、建议的实现策略

为了降低复杂度，推荐你按下面策略实现：

### 方案建议

* 当前 `TiptapEditor.vue` 内部先用 `textarea`
* 通过 `selectionStart` / `selectionEnd` 获取选区
* `ExtractToolbar` 显示在编辑区域顶部右上角或选区附近的简化位置
* 弹窗确认后：

  * 创建 `ExtractItem`
  * 创建 tree node
  * 自动展开父节点
  * 自动选中新节点

这种方案最适合作为第一阶段闭环。

---

## 十六、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text id="vhl6n1"
src/modules/extract/types.ts
```

然后再给完整代码。

### 第四步：最后说明：

1. 如何接入当前项目
2. 用户如何在界面中完成“选区创建摘录”
3. 当前实现相对真实 Tiptap/后端版的简化点有哪些

---

## 十七、补充约束

```text id="gxuc5y"
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入复杂第三方浮动定位库
- 可以使用 Element Plus Dialog / Select / Input 等基础控件
- 当前阶段优先保证结构清晰、模块边界清晰、可运行
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 十八、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。
