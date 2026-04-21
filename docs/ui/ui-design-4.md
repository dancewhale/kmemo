# 任务单 4：editor 模块壳

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **第一阶段的 editor 模块壳**。

这个项目是一个基于 **Wails + Go 后端** 的桌面知识工具，产品风格参考 **SuperMemo 风格的渐进阅读 / 知识工作台**。
本次只实现 **editor 模块的第一阶段可运行版本**，要求 **先用 mock 数据，不接真实后端，不完整接入 Tiptap**，但要为后续接入 **富文本编辑、选区摘录、标注、高亮、卡片生成** 预留清晰边界。

---

## 一、模块目标

请实现一个 **editor 模块壳**，能够完成：

1. 展示一个右侧编辑区壳子
2. 展示文档标题、更新时间、内容类型等基础信息
3. 提供一个基础工具栏壳
4. 提供一个可编辑区域占位
5. 通过 Pinia 管理当前编辑文档状态
6. 支持 dirty 状态
7. 支持简单“保存状态”展示
8. 使用 mock 数据驱动
9. 为后续接入 Tiptap 保留清晰边界

当前阶段 **不要完整实现富文本编辑器**，也 **不要实现复杂命令系统、选区工具栏、摘录生成**。
但组件拆分必须合理，不能把未来所有逻辑都塞进一个文件。

---

## 二、项目上下文说明

当前项目已经具备：

* `MainLayout`
* `WorkspaceShell`
* `RightPane.vue`
* `reader` 模块
* `knowledge-tree` 模块
* `workspace.store.ts`
* `shared/components`
* `mock/` 目录

当前工作台结构中：

* 中间区域展示列表 / 树 / 队列
* 右侧区域承担详情 / 预览 / 编辑职责

这次 editor 模块的目标，是把右侧区域从“简单详情面板”升级为“可承载未来编辑器的稳定容器”。

项目强调：

* 高信息密度
* 桌面知识工具风格
* 不是后台管理系统
* 模块边界清晰
* 组件职责清晰
* 后续要便于接入 Tiptap

---

## 三、本次要实现的文件

请至少实现和补全以下文件：

```text
src/modules/editor/
  components/
    EditorShell.vue
    TiptapEditor.vue
    EditorToolbar.vue
  stores/
    editor.store.ts
  types.ts
```

并根据需要补充或修改以下文件，使 editor 模块能在当前工作台中真正运行：

```text
src/mock/articles.ts
src/mock/index.ts
src/modules/workspace/components/RightPane.vue
```

如有必要，你也可以新增极少量辅助文件，但请控制复杂度，不要扩展成完整子系统。

---

## 四、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要完整接入 Tiptap 扩展体系
3. 不要实现选区摘录
4. 不要实现高亮注释
5. 不要实现卡片生成
6. 不要实现复杂快捷键系统
7. 不要实现协同编辑
8. 不要实现复杂自动保存策略
9. 不要把 editor 模块写成一个大组件
10. 不要提前实现未来 review / command-center 的复杂逻辑

---

## 五、类型设计要求

请先定义清晰的类型。

### `types.ts` 至少包括：

```ts
export type EditorDocumentType = 'article' | 'node' | 'extract'

export interface EditorDocument {
  id: string
  title: string
  content: string
  contentType: EditorDocumentType
  updatedAt: string
  sourceUrl?: string
  tags?: string[]
}

export interface EditorStateSnapshot {
  originalContent: string
  currentContent: string
  dirty: boolean
  saving: boolean
}
```

你可以在此基础上微调，但必须满足：

1. 明确区分“文档元信息”和“编辑状态”
2. 能支持未来对接 Tiptap JSON / HTML / Markdown 等扩展
3. 不要过度复杂
4. 当前阶段足够支撑右侧编辑区展示与基本编辑

---

## 六、mock 数据要求

请沿用或补充 `src/mock/articles.ts` 中的文章数据，使其具备足够的正文内容。

要求：

1. 每篇文章至少有：

   * `id`
   * `title`
   * `summary`
   * `content`
   * `sourceType`
   * `status`
   * `updatedAt`

2. `content` 必须是真实感较强的长文本，而不是占位一句话

3. 内容风格应符合知识工具场景，例如：

   * Incremental Reading
   * Active Recall
   * Progressive Summarization
   * Knowledge Workflow
   * Note Design
   * Reading Strategy

这样 editor 模块才能展示出真正像“编辑 / 阅读工作区”的感觉。

---

## 七、store 设计要求

请实现 `editor.store.ts`，使用 Pinia 管理当前编辑文档状态。

### 至少包含这些状态：

```ts
interface EditorState {
  currentDocument: EditorDocument | null
  originalContent: string
  currentContent: string
  dirty: boolean
  saving: boolean
  lastSavedAt: string | null
}
```

### 至少包含这些 getter / 计算逻辑：

1. `hasDocument`
2. `documentTitle`
3. `documentType`
4. `saveStatusText`

### 至少包含这些 action：

1. `openDocument(doc: EditorDocument)`
2. `setContent(content: string)`
3. `markSaved()`
4. `resetContent()`
5. `clearDocument()`

### 行为要求：

* 打开文档时，初始化 `originalContent` 和 `currentContent`
* 当内容变化时自动更新 `dirty`
* `markSaved()` 可以是 mock 行为，不需要真实 API
* `saveStatusText` 需要能体现：

  * idle
  * edited / unsaved
  * saving
  * saved

注意：

* store 只管编辑状态
* 不要把 reader 的列表逻辑挪进 editor store
* 不要把 workspace 的 pane 状态挪进 editor store

---

## 八、组件要求

---

### 1. `EditorShell.vue`

这是 editor 模块的主容器组件。

要求：

1. 展示文档头部区域
2. 展示 `EditorToolbar.vue`
3. 展示 `TiptapEditor.vue`
4. 当没有文档时显示 `AppEmpty`
5. 从 `editor.store.ts` 读取当前文档和编辑状态
6. 风格像“右侧知识编辑面板”，不是后台表单页

头部区域至少显示：

* 文档标题
* 文档类型
* 更新时间
* 保存状态

可以使用较轻量的布局，例如上方 metadata 条 + 中部编辑器区域。

---

### 2. `EditorToolbar.vue`

这是编辑器上方的工具栏壳。

当前阶段不要做完整富文本能力，但要预留结构。

要求：

1. 展示基础按钮占位，例如：

   * Bold
   * Italic
   * Highlight
   * Quote
   * Extract
2. 当前阶段这些按钮可以：

   * 仅做 UI 展示
   * 或给出 mock 提示
3. 工具栏不要写复杂业务逻辑
4. 工具栏风格紧凑、专业、可扩展
5. 工具栏与具体编辑器实现解耦

目标是：后续换成真实 Tiptap 命令时，不需要重做组件边界。

---

### 3. `TiptapEditor.vue`

虽然名字叫 `TiptapEditor.vue`，但当前阶段不要完整接入 Tiptap。

要求：

1. 提供一个可编辑区域壳
2. 当前阶段可以用：

   * `textarea`
   * 或 `contenteditable`
   * 或一个极简版输入区域
3. 但组件边界必须保留为未来的 Tiptap 容器
4. 接收 `modelValue`
5. 发出 `update:modelValue`
6. 有良好的占位文本
7. 样式要像编辑工作区，而不是普通表单输入框
8. 支持基本滚动
9. 支持较长文本编辑

不要把它写成普通表单 `el-input type="textarea"` 的后台配置页感觉。
它应更像一个“知识内容编辑面板”。

---

## 九、Workspace 集成要求

请将 editor 模块集成到当前工作台右侧区域。

### 对 `RightPane.vue` 的要求

当前右侧区域在不同 context 下展示不同内容。

请调整逻辑如下：

### 当 `context === 'reading'`

* 右侧优先展示 `EditorShell.vue`
* 当前选中的 article 应映射为 `EditorDocument`
* 打开阅读项时，右侧不是只显示纯详情，而是显示“可编辑/可预览的文档工作区”

### 当 `context === 'knowledge'`

* 当前阶段可以继续保留知识节点详情展示
* 先不要强行把 knowledge 节点也接到 editor
* 但代码结构上请避免未来无法扩展

### 当 `context === 'inbox'`

* 可以与 `reading` 类似，展示选中文章的 editor 壳
* 或保留轻量详情，但要说明你的选择

### 当 `context === 'review'`

* 不接 editor，继续保留复习详情面板

重点：

* editor 当前主要承接 `reading` 场景
* 为未来“阅读 + 编辑 + 摘录”的闭环打基础

---

## 十、Reader 与 Editor 的协作要求

请注意 `reader.store.ts` 与 `editor.store.ts` 的边界：

### `reader.store.ts`

负责：

* 列表数据
* 当前选中的 article id
* 搜索筛选

### `editor.store.ts`

负责：

* 当前打开的文档
* 内容变化
* dirty 状态
* 保存状态

推荐做法：

* 在 `RightPane.vue` 中监听当前选中的 article
* 将其映射为 `EditorDocument`
* 然后调用 `editorStore.openDocument(...)`

请不要：

* 直接把 reader store 的整份数据塞进 editor store
* 在 editor store 内直接实现文章列表逻辑

---

## 十一、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 编辑区应有“专业知识工作台”感
3. 小圆角
4. 轻边框
5. 克制的 hover 和 active 状态
6. 更像知识工具，不像后台表单
7. 工具栏要简洁
8. 内容区要舒适可读
9. 不要营销站风格，不要花哨

编辑区建议视觉结构：

* 上方：标题和元信息
* 中间：工具栏
* 下方：正文编辑区

正文区域应当让人感觉未来能承载：

* 阅读
* 标注
* 摘录
* 问答卡生成

---

## 十二、代码结构要求

请特别遵守：

1. 不要把所有 editor 逻辑写在 `RightPane.vue`
2. 不要把所有状态都用局部 `ref` 管
3. 不要让 `EditorShell.vue` 过度臃肿
4. 不要硬编码太多魔法字符串
5. 类型要明确
6. 代码应为“第一阶段真实可运行版本”
7. 所有 import 必须完整
8. 不要只写伪代码或片段

---

## 十三、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text
src/modules/editor/types.ts
```

然后再给完整代码。

### 第四步：最后说明如何把它接入当前项目并运行

---

## 十四、附加实现提示

为了降低未来接入真实 Tiptap 的成本，请在结构上提前考虑：

* `EditorShell` 独立
* `EditorToolbar` 独立
* `TiptapEditor` 独立
* `editor.store.ts` 独立管理打开文档与 dirty 状态
* 当前 `TiptapEditor.vue` 虽然先不用真实 Tiptap，但接口风格要尽量接近未来富文本容器的职责

建议让 `TiptapEditor.vue` 至少具备：

* `modelValue`
* `placeholder`
* `readonly`
* `update:modelValue`

这样后续替换为真实 Tiptap 时更自然。

---

## 十五、补充约束

```text
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入过多额外依赖
- 可以使用少量 Element Plus 基础控件，但不要把 editor 模块做成后台表单页
- 当前阶段优先保证结构清晰、模块边界清晰、可运行
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 十六、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。

