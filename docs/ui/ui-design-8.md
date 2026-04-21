# 任务单 8：命令面板 + 快捷键系统（Command Palette + Shortcuts）第一阶段（可直接投喂 AI）

你是一名资深 Vue 3 / TypeScript 工程师。请在一个已经存在的 Vue 3 + TypeScript + Vite + Element Plus + Vue Router + Pinia 项目中，实现 **“命令面板 + 快捷键系统（Command Palette + Shortcuts）第一阶段可运行版本”**。

这个项目是一个基于 **Wails + Go 后端** 的桌面知识工具，产品定位参考 **SuperMemo 风格的渐进阅读 / 知识工作台**。
当前项目已经具备：

* `workspace` 模块
* `reader` 模块
* `knowledge-tree` 模块
* `editor` 模块
* `extract` 模块
* `review` 模块
* 阅读 / 摘录 / 知识树 / 复习队列基本联动
* `workspace.store.ts`
* `reader.store.ts`
* `tree.store.ts`
* `editor.store.ts`
* `extract.store.ts`
* `review.store.ts`
* `mock/` 数据
* 主工作台布局
* `app/providers/shortcuts.ts` 预留文件
* `modules/command-center/` 已预留目录或可新增

本次目标是：**实现第一阶段的命令面板和快捷键系统**，让用户能够通过键盘快速完成模块跳转、聚焦、打开对象、触发核心动作。

当前阶段重点是 **命令流、快捷键分发、命令面板 UI、跨模块动作协调**，不是复杂配置系统。

---

## 一、模块目标

请实现如下能力：

1. 用户可以通过全局快捷键打开命令面板
2. 命令面板支持输入关键词搜索命令
3. 用户可以通过键盘上下选择命令并回车执行
4. 至少支持一批核心命令，例如：

   * Go to Reading
   * Go to Knowledge
   * Go to Review
   * Focus Search
   * Open First Review Item
   * Open Selected Article
   * Open Selected Extract
   * Create Extract（若当前有选区）
   * Reveal Current Extract in Tree
   * Back to Source Article
5. 支持基础全局快捷键，例如：

   * `Ctrl/Cmd + K` 打开命令面板
   * `g r` 或单快捷键切到 Reading
   * `g k` 切到 Knowledge
   * `g v` 切到 Review
   * `Esc` 关闭命令面板
6. 命令执行后，工作台上下文和对应模块状态能同步变化
7. 命令系统结构清晰，为未来扩展更多命令预留边界

当前阶段 **不要求复杂多级命令系统**，但要做到真实可用。

---

## 二、明确范围限制

这次严格不要做以下事情：

1. 不要接真实 Wails API
2. 不要实现用户自定义快捷键配置页面
3. 不要实现复杂 Vim 模式
4. 不要实现完整命令历史系统
5. 不要实现 fuzzy search 第三方重量级方案
6. 不要实现权限系统
7. 不要实现复杂上下文菜单系统
8. 不要实现快捷键冲突检测配置器
9. 不要把所有命令逻辑塞进一个组件
10. 不要把所有模块状态硬耦合在一起

---

## 三、本次建议实现/修改的文件

请至少实现和补全以下文件：

```text id="vtjtmv"
src/modules/command-center/
  components/
    CommandPalette.vue
    CommandList.vue
    CommandListItem.vue
    ShortcutHint.vue
  stores/
    command.store.ts
  types.ts
  services/
    command.registry.ts
    command.executor.ts
    command.matcher.ts
```

并根据需要补充或修改以下文件：

```text id="e3fb8r"
src/app/providers/shortcuts.ts
src/layouts/MainLayout.vue
src/modules/workspace/stores/workspace.store.ts
src/modules/review/stores/review.store.ts
src/modules/reader/stores/reader.store.ts
src/modules/extract/stores/extract.store.ts
src/modules/knowledge-tree/stores/tree.store.ts
src/modules/editor/stores/editor.store.ts
src/shared/constants/commands.ts
```

如果有必要，也可以新增少量辅助文件，例如：

```text id="yjlwm6"
src/shared/utils/shortcut.ts
src/shared/types/command.ts
src/shared/composables/useGlobalHotkeys.ts
```

但请保持克制，不要扩展成完整输入系统框架。

---

## 四、整体设计要求

请把这次能力拆成 4 个清晰层次：

### 1. 命令定义层

负责：

* 定义命令元信息
* 命令 id
* 标题
* 描述
* 所属分组
* 默认快捷键提示
* 是否可执行
* 执行条件

### 2. 命令检索与匹配层

负责：

* 根据输入关键词过滤命令
* 返回匹配结果
* 当前阶段只需要简单匹配即可

### 3. 命令执行层

负责：

* 真正调度各模块 store action
* 处理跨模块跳转
* 避免把执行逻辑写死在 UI 组件里

### 4. 快捷键监听层

负责：

* 监听全局快捷键
* 打开/关闭命令面板
* 触发部分高频命令
* 不要影响输入框和文本编辑正常使用

要求边界清晰，方便后续扩展更多命令。

---

## 五、类型设计要求

请先整理和定义清晰类型。

---

### `src/modules/command-center/types.ts`

至少包含：

```ts id="pry40h"
export type CommandGroup =
  | 'navigation'
  | 'focus'
  | 'reading'
  | 'extract'
  | 'review'
  | 'knowledge'
  | 'system'

export interface CommandItem {
  id: string
  title: string
  subtitle?: string
  group: CommandGroup
  keywords?: string[]
  shortcut?: string[]
  enabled?: boolean
}

export interface CommandMatchResult extends CommandItem {
  score: number
}

export interface CommandState {
  isOpen: boolean
  query: string
  activeIndex: number
}
```

你可以在此基础上微调，但必须满足：

1. 命令定义和命令状态分离
2. 命令对象足够支持 UI 展示和执行
3. 能支持后续扩展上下文感知命令
4. 类型不要过度复杂

---

## 六、命令注册要求

请新增一个注册层：

### `command.registry.ts`

用于集中定义系统当前支持的命令。

至少定义以下命令：

### Navigation

* `go-reading`
* `go-knowledge`
* `go-review`
* `go-inbox`
* `go-search`

### Review

* `open-first-review-item`

### Reading

* `open-selected-article`

### Extract

* `create-extract`
* `open-selected-extract`
* `back-to-source-article`

### Knowledge

* `reveal-current-extract-in-tree`

### System

* `toggle-command-palette`
* `close-command-palette`

要求：

1. 注册表统一定义命令元信息
2. 不要把命令元信息散落在组件里
3. 命令执行逻辑不要写在 registry 里
4. 每个命令应包含：

   * id
   * title
   * subtitle
   * group
   * keywords
   * shortcut（可选）

---

## 七、命令匹配要求

请新增：

### `command.matcher.ts`

至少实现一个简单命令匹配函数，例如：

```ts id="sn5or2"
matchCommands(commands: CommandItem[], query: string): CommandMatchResult[]
```

要求：

1. 当前阶段只做简单匹配即可
2. 可以按以下规则实现：

   * title 包含
   * subtitle 包含
   * keywords 包含
3. 可给 title 命中更高 score
4. query 为空时返回全部命令
5. 不需要引入复杂 fuzzy search 第三方库

这样足够第一阶段使用。

---

## 八、命令执行要求

请新增：

### `command.executor.ts`

这是这次最关键的文件之一。
它负责真正执行命令，不要把执行逻辑写在组件内部。

请至少实现：

```ts id="et21m0"
executeCommand(commandId: string): void | Promise<void>
```

要求：

* 通过各模块 store 协调行为
* 统一处理跨模块跳转
* 允许内部按需读取：

  * workspace store
  * reader store
  * tree store
  * extract store
  * review store
  * editor store

至少支持这些执行行为：

### 1. `go-reading`

* `workspace.store.setContext('reading')`

### 2. `go-knowledge`

* `workspace.store.setContext('knowledge')`

### 3. `go-review`

* `workspace.store.setContext('review')`

### 4. `go-inbox`

* `workspace.store.setContext('inbox')`

### 5. `go-search`

* `workspace.store.setContext('search')`

### 6. `open-first-review-item`

* 打开第一条可用 review item
* 切到 review

### 7. `open-selected-article`

* 若当前 reader 有选中文章，则切到 reading 并确保右侧打开它

### 8. `create-extract`

* 只有当 editor 当前存在非空选区时可生效
* 触发现有的摘录创建流程
* 若当前项目已有 `extract.store.openCreateDialog(...)` 等能力，请复用

### 9. `open-selected-extract`

* 若当前已有 `selectedExtractId`，则切到 knowledge 并打开摘录详情

### 10. `back-to-source-article`

* 若当前 selected extract 存在 sourceArticleId，则跳回文章

### 11. `reveal-current-extract-in-tree`

* 若当前 selected extract 有 treeNodeId，则切到 knowledge 并选中树节点

### 12. `toggle-command-palette`

* 打开/关闭命令面板

### 13. `close-command-palette`

* 关闭命令面板

要求：

* 执行成功后一般自动关闭命令面板
* 执行逻辑要清晰
* 不要在 executor 中夹杂 UI 渲染逻辑

---

## 九、store 设计要求

---

### `command.store.ts`

请实现命令中心 store。

至少包含这些状态：

```ts id="gk2ori"
interface CommandStoreState {
  isOpen: boolean
  query: string
  activeIndex: number
}
```

至少包含这些 getter：

1. `commands`
2. `matchedCommands`
3. `activeCommand`
4. `hasResults`

至少包含这些 action：

1. `open()`
2. `close()`
3. `toggle()`
4. `setQuery(query: string)`
5. `moveUp()`
6. `moveDown()`
7. `setActiveIndex(index: number)`
8. `executeActiveCommand()`
9. `resetState()`

### 行为要求

* 打开命令面板时，重置 query 与 activeIndex
* query 变化时自动更新匹配结果
* 结果变化后，activeIndex 应保证不越界
* 执行命令后默认关闭命令面板
* store 只管命令面板状态，不要把模块业务逻辑塞进这里

---

## 十、快捷键监听要求

请完善：

### `src/app/providers/shortcuts.ts`

要求实现全局快捷键监听，并在应用根布局中启用。

至少支持这些快捷键：

### 全局命令面板

* `Ctrl/Cmd + K`：打开或关闭命令面板

### 面板内导航

当命令面板打开时：

* `ArrowDown`：下移
* `ArrowUp`：上移
* `Enter`：执行当前命令
* `Escape`：关闭面板

### 快速导航（第一阶段可选用简单方案）

建议支持这些：

* `Alt + 1` -> Inbox
* `Alt + 2` -> Reading
* `Alt + 3` -> Knowledge
* `Alt + 4` -> Review
* `Alt + 5` -> Search

### 可选补充

* `Ctrl/Cmd + Shift + R` -> Open First Review Item

### 重要约束

快捷键系统必须避免干扰正常文本输入。

要求：

1. 当焦点在 input / textarea / contenteditable 中时

   * 不要响应普通导航类快捷键
   * 但 `Ctrl/Cmd + K` 仍可允许
2. 面板打开时，优先由面板接管方向键 / Enter / Esc
3. 不要做复杂键盘状态机
4. 第一阶段保持简单可靠

---

## 十一、组件要求

---

### 1. `CommandPalette.vue`

这是命令面板主组件，本次最核心 UI 组件。

要求：

1. 从 `command.store.ts` 读取状态
2. 当 `isOpen = true` 时显示
3. 包含：

   * 搜索输入框
   * 命令结果列表
   * 底部快捷提示区（可选）
4. 输入框自动聚焦
5. 输入 query 实时过滤命令
6. 支持键盘上下选择
7. 支持回车执行
8. 支持鼠标点击执行
9. 支持 Esc 关闭
10. 风格应像桌面工具命令面板，不像后台搜索弹窗

建议视觉结构：

* 顶部搜索输入
* 中间结果列表
* 底部一行提示：

  * Enter to run
  * ↑↓ to navigate
  * Esc to close

可以使用简洁的 overlay + center panel 方案。

---

### 2. `CommandList.vue`

职责：

* 渲染匹配命令列表
* 使用 `CommandListItem.vue`
* 无结果时显示空状态

要求：

* 列表滚动
* 当前 active item 高亮
* 不要写太多业务逻辑

---

### 3. `CommandListItem.vue`

职责：

* 渲染单条命令
* 展示：

  * title
  * subtitle
  * group
  * shortcut hint
* active 时高亮
* 点击时执行

风格要求：

* 紧凑
* 信息层级清楚
* 更像命令系统，不像普通菜单

---

### 4. `ShortcutHint.vue`

职责：

* 统一展示快捷键提示，例如：

  * `Ctrl K`
  * `Alt 3`
  * `Enter`
* 用于命令项右侧和面板底部提示区

要求：

* 简洁
* 可复用
* 不要过度装饰

---

## 十二、MainLayout 集成要求

请将命令面板挂到应用主布局。

### `MainLayout.vue`

要求：

1. 挂载 `CommandPalette.vue`
2. 在布局层初始化全局快捷键监听
3. 不要把命令执行逻辑写在布局里
4. 布局只负责：

   * 提供命令面板可见性入口
   * 启动 shortcuts provider

---

## 十三、与现有模块的联动要求

这次的价值在于把现有模块串起来。请实现下面这些联动。

---

### 联动 1：命令 -> workspace context

通过命令可切换：

* inbox
* reading
* knowledge
* review
* search

---

### 联动 2：命令 -> review

执行 `open-first-review-item` 时：

1. 打开第一条可用 review item
2. 切到 review context
3. 右侧显示 review card

---

### 联动 3：命令 -> extract 创建

执行 `create-extract` 时：

1. 若 editor 当前无选区，则该命令禁用或执行无效果并说明
2. 若有选区，则复用已有摘录创建流程

建议更好的第一阶段实现：

* 在命令列表展示时仍可显示命令
* 但 `enabled` 基于当前上下文和状态控制
* disabled 状态在 UI 中弱化显示

---

### 联动 4：命令 -> 原文回跳

执行 `back-to-source-article` 时：

1. 若当前 selected extract 有 sourceArticleId
2. reader 打开对应 article
3. workspace 切到 reading

---

### 联动 5：命令 -> 树定位

执行 `reveal-current-extract-in-tree` 时：

1. 若当前 selected extract 有 treeNodeId
2. tree store 选中该节点
3. workspace 切到 knowledge

---

## 十四、命令可用性要求

建议命令支持简单的上下文启用/禁用判断。

例如：

### 永远可用

* go-reading
* go-knowledge
* go-review
* go-inbox
* go-search
* toggle-command-palette

### 条件可用

* create-extract：需要 editor 有 selection
* open-selected-article：需要当前有 selected article
* open-selected-extract：需要当前有 selected extract
* back-to-source-article：需要 selected extract 且有 sourceArticleId
* reveal-current-extract-in-tree：需要 selected extract 且有 treeNodeId
* open-first-review-item：需要 review queue 非空

实现方式可以是：

* registry 提供静态命令定义
* executor 或 store 中动态判断 enabled
* 你选择一种清晰且简单的方案并说明

---

## 十五、UI 风格要求

请严格遵循以下风格：

1. 高信息密度
2. 专业、克制
3. 轻边框
4. 小圆角
5. 更像桌面知识工具，不像后台搜索框
6. 命令面板要有“工作台命令中心”的感觉
7. 快捷键提示要简洁清楚
8. 当前 active item 高亮明显但不过度花哨

命令面板应让人感觉像：

* 快速跳转入口
* 工作流操作中心
* 键盘驱动的知识工具核心设施

---

## 十六、代码结构要求

请特别遵守：

1. 不要把所有命令定义写在组件里
2. 不要把所有执行逻辑写在 `command.store.ts`
3. registry / matcher / executor 要分层
4. `CommandPalette.vue` 只负责 UI 与交互
5. `command.store.ts` 只负责面板状态与选中状态
6. shortcuts provider 只负责监听键盘并分发动作
7. 所有 import 必须完整
8. 不要只写伪代码或片段
9. 代码应为“第一阶段真实可运行版本”

---

## 十七、建议的实现策略

为了降低复杂度，推荐采用以下方案：

### 推荐方案

* registry：静态命令定义
* matcher：简单字符串匹配
* executor：集中调度各 store
* command store：管理 query / activeIndex / open state
* shortcuts provider：统一监听 window keydown

### 为什么推荐这样做

因为这套分层已经足够支撑后续扩展：

* 更多命令
* 上下文敏感命令
* 最近命令
* 用户自定义快捷键
* 模块级命令注入

但当前阶段不会过度复杂。

---

## 十八、输出顺序要求

请严格按以下顺序输出：

### 第一步：输出本次会修改/新增的文件列表

### 第二步：简要说明每个文件职责

### 第三步：开始逐文件输出完整代码

输出代码时请明确文件路径，例如：

```text id="yfx7ju"
src/modules/command-center/types.ts
```

然后再给完整代码。

### 第四步：最后说明：

1. 如何接入当前项目
2. 用户如何通过快捷键和命令面板操作系统
3. 当前实现相对完整命令系统的简化点有哪些

---

## 十九、补充约束

```text id="shnxta"
补充约束：
- 所有文件都要提供完整代码
- 不要省略 import
- 不要留 TODO 代替实现
- 当前阶段不要引入重量级 fuzzy search 或热键库
- 可以使用少量 Element Plus Dialog / Input / Empty / Tag 等基础控件
- 当前阶段优先保证结构清晰、模块边界清晰、可运行
- UI 风格偏 SuperMemo / 知识工作台，不偏中后台
```

---

## 二十、现在开始执行

请从“本次会修改/新增的文件列表”开始，然后按要求继续输出完整代码。

---

## 这一步完成后的最佳下一步

做完这个任务后，你的产品在“可操作性”上会提升一个档次。
接下来最值得做的是：
