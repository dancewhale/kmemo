# 任务单 2：knowledge-tree 模块

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **第一阶段的 knowledge-tree 模块**。

这个项目是一个基于 **Wails + Go 后端** 的桌面知识工具，产品风格参考 **SuperMemo 风格的渐进阅读/知识工作台**。
本次只实现 **knowledge-tree 模块的第一阶段可运行版本**，要求 **先用 mock 数据，不接真实后端，不实现完整 he-tree 拖拽树**，但要 **为后续接入 he-tree 预留清晰边界**。

---

## 一、模块目标

请实现一个 **知识树模块的第一阶段骨架**，能够完成：

1. 展示知识节点列表
2. 支持节点层级展示
3. 支持展开 / 折叠
4. 支持选中节点
5. 支持当前节点高亮
6. 支持简单的节点图标与类型标识
7. 支持右侧详情区读取当前选中节点
8. 使用 Pinia 管理树状态
9. 使用 mock 数据驱动
10. 为后续 he-tree 接入保留结构和接口边界

当前阶段 **不要实现完整拖拽排序**，也 **不要把 he-tree 直接深度集成进去**。
你可以先不安装 he-tree，也可以仅在代码结构层面预留未来接入点。

---

## 二、项目上下文说明

项目第一阶段已经存在以下结构和约束：

* 有 `MainLayout`
* 有 `WorkspaceShell`
* 有 `KnowledgePage.vue`
* 有 `workspace.store.ts`
* 有 `shared/components` 基础组件
* 有 `mock/` 目录
* 页面通过 `WorkspaceShell` 根据 `context` 切换内容
* 当前 `KnowledgePage.vue` 会以 `context="knowledge"` 渲染工作台
* 当前项目强调：

  * 高信息密度
  * 桌面知识工具风格
  * 不是后台管理系统
  * 组件边界清晰
  * 模块职责清晰

---

## 三、本次要实现的文件

请至少实现和补全以下文件：

```text id="0vdwz2"
desktop/frontend/src/modules/knowledge-tree/
  components/
    TreePanel.vue
    TreeNodeItem.vue
  stores/
    tree.store.ts
  services/
    tree.mapper.ts
  types.ts
```

如果你认为有必要，可以补充少量辅助文件，但请控制复杂度，不要扩展成完整子系统。

另外，请根据需要补充或修改以下文件，使 knowledge 模式真的能工作：

```text id="yz63hm"
src/mock/tree.ts
src/mock/index.ts
src/modules/workspace/components/CenterPane.vue
src/modules/workspace/components/RightPane.vue
```

如果你需要新增一个轻量常量文件或工具文件，也可以，但要说明原因。

---

## 四、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现完整 he-tree 拖拽
3. 不要实现右键菜单复杂逻辑
4. 不要实现树节点重命名
5. 不要实现新增/删除节点表单
6. 不要实现复杂懒加载
7. 不要写成 Element Plus 的普通 Tree 演示页面
8. 不要写成后台管理系统左侧菜单树
9. 不要把所有树逻辑塞进一个组件
10. 不要提前实现未来 editor/review 的业务逻辑

---

## 五、类型设计要求

请先定义清晰的类型。

### `types.ts` 至少包括：

```ts id="8z3bqn"
export type KnowledgeNodeType = 'topic' | 'article' | 'extract' | 'card'

export interface KnowledgeNode {
  id: string
  parentId: string | null
  title: string
  type: KnowledgeNodeType
  expanded?: boolean
  selected?: boolean
  hasChildren?: boolean
  childCount?: number
  createdAt?: string
  updatedAt?: string
}

export interface UITreeNode extends KnowledgeNode {
  level: number
  children: UITreeNode[]
}
```

你可以在此基础上微调，但请满足以下要求：

1. 区分原始节点和 UI 节点
2. UI 节点中应包含层级信息
3. 要能支持未来接入 he-tree
4. 类型设计不要过度复杂

---

## 六、mock 数据要求

请在 `src/mock/tree.ts` 中提供一份更适合知识工具的树数据，至少包含：

* 8~15 个节点
* 有明确层级
* 节点类型混合，例如：

  * topic
  * article
  * extract
  * card

建议示例结构类似：

* Reading Inbox

  * Article: Incremental Reading Basics
  * Article: Learning in Public
* Knowledge System

  * Topic: Memory

    * Extract: Active recall note
    * Card: What is active recall?
  * Topic: Reading Workflow

    * Extract: Highlight processing pipeline

要求 mock 数据有足够真实感，不要只有 “Node 1 / Node 2”。

---

## 七、store 设计要求

请实现 `tree.store.ts`，使用 Pinia 管理知识树状态。

### 至少包含这些状态：

```ts id="b4f8zo"
interface TreeState {
  rawNodes: KnowledgeNode[]
  selectedNodeId: string | null
  expandedNodeIds: string[]
  searchKeyword: string
  loading: boolean
}
```

### 至少包含这些 getter / 计算逻辑：

1. `treeNodes`：将扁平节点转换为树结构
2. `visibleNodes`：得到当前展开状态下可见的平铺节点列表
3. `selectedNode`
4. `rootNodes`

### 至少包含这些 action：

1. `initialize()`
2. `setSelectedNode(id: string)`
3. `toggleNode(id: string)`
4. `expandNode(id: string)`
5. `collapseNode(id: string)`
6. `setSearchKeyword(keyword: string)`

### 行为要求：

* 初始化时加载 mock 数据
* 默认展开部分顶级节点
* 点击节点时更新选中状态
* 折叠父节点时隐藏子节点
* 搜索关键词先只保存在 store 中，可不实现完整过滤逻辑
* store 职责清晰，不要与 workspace store 混成一体

---

## 八、`tree.mapper.ts` 要求

请新增 `services/tree.mapper.ts`，负责将扁平节点结构映射为 UI 树结构。

至少实现：

1. `buildTree(nodes: KnowledgeNode[]): UITreeNode[]`
2. `flattenVisibleTree(nodes: UITreeNode[], expandedNodeIds: string[]): UITreeNode[]`

要求：

* 不要把映射逻辑直接写在组件里
* 逻辑清晰，便于以后接真实后端和 he-tree
* 函数尽量保持纯函数风格

---

## 九、组件要求

---

### 1. `TreePanel.vue`

这是知识树的主容器组件。

要求：

1. 展示树模块标题，例如 `Knowledge Tree`
2. 有一个轻量工具栏区域
3. 可包含一个搜索输入框占位
4. 渲染当前可见节点列表
5. 使用 `TreeNodeItem.vue` 渲染单个节点
6. 从 `tree.store.ts` 读取数据
7. 首次挂载时调用 `initialize()`
8. 无数据时显示 `AppEmpty`
9. 加载态时显示 `AppLoading`

风格要求：

* 紧凑
* 专业工具感
* 不要后台菜单样式
* 不要像通用文件浏览器
* 更像知识工作台左/中区域里的知识树面板

---

### 2. `TreeNodeItem.vue`

这是单个树节点渲染组件。

props 建议至少包括：

* `node: UITreeNode`
* `selected: boolean`

行为要求：

1. 根据 `node.level` 显示缩进
2. 如果节点有 children，则显示展开/折叠按钮
3. 点击节点主体时选中节点
4. 点击展开按钮时切换展开/折叠
5. 展示节点标题
6. 展示节点类型图标或类型标签
7. 当前选中节点高亮
8. hover 有轻微反馈
9. 样式紧凑、专业

注意：

* 不要在 `TreeNodeItem.vue` 里直接写整个树的递归渲染
* 第一阶段推荐基于 `visibleNodes` 做平铺渲染
* 组件只负责单节点展示与事件触发

---

## 十、Workspace 集成要求

请将 knowledge-tree 模块集成到当前 `knowledge` context 中。

### 对 `CenterPane.vue` 的要求

当 `context === 'knowledge'` 时：

* 中间区域应显示 `TreePanel.vue`
* 不再只是简单占位文本

### 对 `RightPane.vue` 的要求

当 `context === 'knowledge'` 时：

* 根据当前选中节点显示详情信息
* 至少显示：

  * title
  * id
  * type
  * parentId
  * childCount（如有）
* 如果未选中节点，显示空状态
* 右侧面板依然保持“详情面板”的风格，而不是表单页面

---

## 十一、与 workspace.store 的协作要求

请注意 `tree.store.ts` 与 `workspace.store.ts` 的边界：

### `tree.store.ts`

负责：

* 树数据
* 树展开状态
* 树选中状态
* 树搜索关键词

### `workspace.store.ts`

负责：

* 当前工作台上下文
* pane 尺寸
* 工作台选中对象 id（如果当前结构已有）

如果你需要做协作，请遵循：

* tree 模块内部以 `tree.store.ts` 为主
* 如果当前项目已经有 `selectedNodeId` 放在 `workspace.store.ts`，可以在选中节点时同步更新
* 但不要把所有树逻辑迁移到 workspace store

如果需要同步，请实现得简单、清晰、可读。

---

## 十二、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 小字号但清晰
3. 轻边框
4. 小圆角
5. hover 和 active 状态克制
6. 更像知识工具，不像管理后台
7. 树节点行高紧凑
8. 图标和文字间距合理
9. 不要花哨，不要过度视觉设计

知识树节点建议呈现为：

* 左侧缩进
* 展开箭头
* 类型图标
* 标题文字
* 可选的子节点数量或类型标签

---

## 十三、代码结构要求

请特别遵守：

1. 不要把扁平转树逻辑写在组件内部
2. 不要把所有状态都用局部 `ref` 管
3. 不要让 `TreePanel.vue` 超级臃肿
4. 不要硬编码太多魔法字符串
5. 尽量使用明确类型
6. 代码应为“第一阶段真实可运行版本”
7. 所有 import 必须完整
8. 不要只写伪代码或片段

---

## 十四、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text id="e8u3rz"
src/modules/knowledge-tree/types.ts
```

然后再给完整代码。

### 第四步：最后说明如何把它接入当前项目并运行

---

## 十五、附加实现提示

为了降低未来接入 he-tree 的成本，请在结构上提前考虑：

* `tree.mapper.ts` 独立
* `UITreeNode` 独立
* `TreePanel` 不直接依赖未来具体树库
* 当前阶段采用“store 提供 visibleNodes，组件平铺渲染”的方式即可

这会使后续替换为 `he-tree` 更自然。

---

## 十六、补充约束

```text id="tjdg36"
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 可以使用少量 Element Plus 基础控件，但不要让知识树变成后台风格
- 不要用递归组件一次性把所有逻辑塞满
- 当前阶段优先保证结构清晰、模块边界清晰、可运行
```

---

## 十七、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。

