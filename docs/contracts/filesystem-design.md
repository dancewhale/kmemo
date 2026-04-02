# Contracts 文件系统设计文档

## 1. 目标

本文档用于细化 `docs/contracts-design.md` 中 `FileStore` 相关部分，明确 kmemo 在 contracts 层中应该如何定义“文件系统能力”的边界、术语、路径策略和生命周期约束。

目标不是提前把所有文件操作都实现成复杂抽象，而是先回答下面这些后续实现一定会遇到的问题：

- actions / flows 向文件系统能力真正要的是什么
- repository 与 contracts 在“文件引用”和“物理文件”之间如何分工
- 哪些文件属于 durable data，哪些只是 cache
- 哪些路径可以进入数据库，哪些路径不应被业务层持久化
- adapter 应负责哪些路径分配、原子写入、删除与恢复逻辑
- 当前已有的 `Asset.StoragePath`、`Card.HTMLPath`、`SourceDocument.FilePath` 等字段，应遵循什么统一规则

这份文档是后续设计 `internal/contracts/filestore.go`、对应 adapter、以及相关 actions / flows 的专项依据。

---

## 2. 与整体分层的关系

文件系统能力仍然遵循 `docs/contracts-design.md` 中的总体分层：

```text
Wails UI API
    ↓
Actions
    ↓
Rules / Flows
    ↓
Contracts
    ↓
Adapters
    ↓
外部能力实现（本地文件系统 / 未来可替换存储）
```

同时它与 repository 是并列关系，而不是替代关系：

```text
Actions / Flows
   ├── Repository   # 保存业务对象、元数据、相对路径引用
   └── Contracts    # 保存、读取、移动、删除物理 artifact
```

### 一句话理解

- `repository` 解决“业务对象记住了哪些文件引用”
- `contracts` 解决“系统如何安全地管理这些文件本体”

---

## 3. 职责边界

## 3.1 contracts 应该定义什么

contracts 层应该定义的是**面向业务能力的文件系统约定**，而不是 `os` 包调用细节。

这里的重点包括：

- 文件对象的语义分类
- durable / cache 边界
- 路径是否允许进入数据库
- 路径以什么形式保存
- 生命周期的关键不变量
- 对上层暴露什么错误语义

contracts 关心的问题应该是：

- 保存一份卡片 HTML 后，什么时候它才算“可被业务引用”
- 一个 asset 丢失时，上层应该收到什么语义错误
- 某类 artifact 应保存到 durable data 还是 cache
- repository 中保存的路径是否可迁移、可备份、可恢复

## 3.2 adapters 应该负责什么

adapters 负责把 contracts 的抽象能力落实到具体物理文件系统。

adapter 负责：

- 基于 config 解析根目录
- 决定 artifact 的物理落盘路径
- 创建目录
- 进行临时写入与原子 promote
- 执行删除、移动、存在性检查
- 把底层 `os.PathError`、权限错误、磁盘异常转换为 contracts 层错误语义

adapter 不应该把底层路径格式、平台差异、目录拼接规则泄漏给上层。

## 3.3 actions / flows 应该负责什么

actions / flows 负责业务编排，而不是手写物理路径。

例如：

- 导入文档后，决定先保存原始来源，再保存转换后的主 HTML，再更新数据库
- 创建卡片时，决定先清理 HTML，再持久化 HTML artifact，再写入 `card.html_path`
- 删除卡片时，决定先做逻辑删除、后做物理清理，还是立即清理

但 actions / flows 不应该：

- 自己拼接 `data/cardfile/...`
- 直接调用 `os.WriteFile`
- 在业务代码里处理平台路径差异

## 3.4 repository 应该负责什么

repository 只负责保存业务元数据和 artifact 引用，不负责物理文件本身。

例如：

- `Asset.StoragePath`
- `Card.HTMLPath`
- `Card.AnswerHTMLPath`
- `SourceDocument.FilePath`

这些字段属于 repository / model 的职责范围；但如何把这些逻辑引用映射到真实磁盘路径，属于 adapter 负责的事情。

## 3.5 config 应该负责什么

config 模块负责提供稳定的根目录和标准子目录，不负责业务对象生命周期。

当前已经存在的目录语义包括：

- `DataDir`
- `CacheDir`
- `LogsDir`
- `AssetsDir`
- `ImportsDir`
- `IndexDir`
- `ExportDir`
- `CardFileDir`
- `ThumbsCacheDir`
- `HTMLCacheDir`
- `IndexCacheDir`

config 负责“这些目录在哪里”；contracts + adapters 负责“业务对象如何使用这些目录”。

---

## 4. 核心术语

为了避免 `file` 这个词过于宽泛，本文档使用 **artifact** 作为核心术语。

## 4.1 Artifact

artifact 指由 kmemo 管理、并具有明确业务语义的文件对象，例如：

- 卡片正文 HTML
- 卡片答案 HTML
- 资源附件
- 导入来源文件
- 转换后的主 HTML
- 搜索索引文件
- 缩略图缓存

artifact 不是任意磁盘文件，而是**被系统明确管理和约束的文件实体**。

## 4.2 Durable artifact

durable artifact 指会被业务对象长期引用、应支持迁移、备份、恢复的文件对象。

典型特征：

- 可以成为业务对象的一部分
- 应随 `DataDir` 一起迁移
- 不应仅靠 cache 重建来恢复
- 可以进入数据库引用字段

## 4.3 Cache artifact

cache artifact 指可重新生成、不可作为业务真相来源的派生文件。

典型特征：

- 删除后不影响业务主数据真实性
- 可通过源数据重新构建
- 不应被持久化为 canonical path
- 适合放在 `CacheDir` 相关目录中

## 4.4 Relative path

relative path 指**相对于 kmemo 管理的数据根目录 `DataDir`** 的路径。

contracts 层中的路径约束，默认都围绕 `DataDir` 展开，而不是围绕某个机器本地绝对路径展开。

---

## 5. Artifact 分类（Artifact Taxonomy）

当前代码和目录语义已经足够支持一套清晰的 artifact 分类。

## 5.1 Card content artifacts

用于卡片渲染的持久内容。

对应字段：

- `Card.HTMLPath`
- `Card.AnswerHTMLPath`

特点：

- 属于 durable artifact
- 应保存在 `DataDir` 体系内
- 应允许被 UI、渲染器、搜索流程稳定引用
- 应有内容哈希（`HTMLHash` / `AnswerHTMLHash`）辅助一致性判断

推荐目录族：

- `CardFileDir`

## 5.2 Asset artifacts

用于卡片内容渲染或展示的媒体与附件。

对应字段：

- `Asset.StoragePath`

当前模型已覆盖：

- `image`
- `audio`
- `video`
- `file`
- `cover`
- `thumbnail`

设计建议：

- `image / audio / video / file / cover` 默认视为 durable artifact
- `thumbnail` 默认视为 cache artifact，除非未来产品要求把缩略图当作正式导出资产

推荐目录族：

- durable -> `AssetsDir`
- cache thumbnail -> `ThumbsCacheDir`

## 5.3 Source / import artifacts

用于保存导入来源、转换结果和后续可追溯来源内容。

对应字段：

- `SourceDocument.FilePath`
- `SourceDocument.OriginalURI`
- `SourceDocument.OriginalHash`

建议拆分的语义包括：

- 原始导入文件
- 转换后的主 HTML
- 未来可扩展的附属资源（图片、样式、附件）

设计建议：

- `OriginalURI` 是来源标识，不等价于本地持久化路径
- `FilePath` 应明确为 durable artifact 的相对路径字段
- 未来即使同时保留原始文件和转换结果，也仍应统一落在导入/来源 artifact 体系下

推荐目录族：

- `ImportsDir`

## 5.4 Index artifacts

用于全文索引或搜索子系统持久化与缓存的文件对象。

相关现状：

- `SearchIndexState` 已经表达索引同步状态
- 但具体磁盘索引文件的生命周期与布局尚未定义

设计建议：

- `IndexDir` 用于 durable index data
- `IndexCacheDir` 用于可重建的索引缓存、临时段、构建中间产物
- DB 中不应把 cache index path 当作业务真相来源保存

## 5.5 分类汇总表

| Artifact 类别 | 当前字段/对象 | Durable | 可重建 | 是否应进入 DB 引用 | 推荐目录族 | 主要负责方 |
| --- | --- | --- | --- | --- | --- | --- |
| Card content | `Card.HTMLPath`, `Card.AnswerHTMLPath` | 是 | 否 | 是 | `CardFileDir` | actions + file adapter + repository |
| Asset | `Asset.StoragePath` | 是 | 视类型而定 | 是 | `AssetsDir` | actions + file adapter + repository |
| Thumbnail cache | `Asset.Kind=thumbnail` 等派生物 | 否 | 是 | 否 | `ThumbsCacheDir` | adapter / async flow |
| Source/import | `SourceDocument.FilePath` | 是 | 否 | 是 | `ImportsDir` | import flow + file adapter + repository |
| Persistent index | 搜索索引持久文件 | 视实现建议为是 | 可重建但不应依赖临时态 | 通常不直接进业务表 | `IndexDir` | indexing flow + adapter |
| Index cache | 索引构建缓存 | 否 | 是 | 否 | `IndexCacheDir` | indexing flow + adapter |
| HTML cache | 派生 HTML 缓存 | 否 | 是 | 否 | `HTMLCacheDir` | html flow + adapter |

---

## 6. Durable 与 Cache 的判定规则

## 6.1 Durable 的判定标准

一个 artifact 满足以下条件时，应优先视为 durable：

- 会被业务对象长期引用
- 删除后无法仅凭 cache 恢复业务完整性
- 需要随 `DataDir` 一起迁移/备份
- 未来导出、同步、恢复时需要保留

当前建议归入 durable 的目录族：

- `AssetsDir`
- `ImportsDir`
- `CardFileDir`
- `IndexDir`（若索引实现选择保留持久化索引数据）

## 6.2 Cache 的判定标准

一个 artifact 满足以下条件时，应优先视为 cache：

- 可以根据 durable data 再生成
- 不应成为数据库中业务对象的唯一文件引用
- 适合被清理、重建、回收
- 对启动性能、预览性能、搜索性能有帮助，但不是业务真相来源

当前建议归入 cache 的目录族：

- `ThumbsCacheDir`
- `HTMLCacheDir`
- `IndexCacheDir`

## 6.3 关于 `IndexDir` 与 `IndexCacheDir`

索引场景比普通文件更容易混淆 durable 与 cache，因此这里需要单独说明：

- `IndexDir` 表示可被长期保留的索引主数据
- `IndexCacheDir` 表示构建过程、查询加速或临时段等可重建产物

当前阶段不强制规定最终搜索实现必须采用哪种物理索引方案，但 contracts 层应先固定原则：

- 若某索引文件是系统可持续使用的主索引结果，应归入 durable
- 若某索引文件只是构建临时产物或可随时清空的加速缓存，应归入 cache

---

## 7. 相对路径策略

## 7.1 总体规则

凡是由 kmemo 管理、且会进入业务持久化模型的 artifact 路径，统一保存为：

**相对于 `DataDir` 的路径**。

这是一条强约束，而不是“推荐但可选”。

## 7.2 适用对象

当前至少应明确覆盖以下字段：

- `Asset.StoragePath`
- `Card.HTMLPath`
- `Card.AnswerHTMLPath`
- `SourceDocument.FilePath`

其中：

- `Asset.StoragePath` 当前已经有清晰注释
- `Card.HTMLPath` / `AnswerHTMLPath` 当前已与 asset 采取一致策略
- `SourceDocument.FilePath` 本文档明确补齐为同样策略

## 7.3 为什么统一相对 `DataDir`

这样设计的主要原因：

1. **可迁移**
   - 换机、云同步、整体迁移用户目录时，不会因为绝对路径变化而失效
2. **跨平台**
   - 避免把 Windows / macOS / Linux 的路径格式差异泄漏到业务模型
3. **可备份**
   - 数据与路径引用都围绕同一个数据根组织，备份更简单
4. **减少耦合**
   - repository 不必知道 config 最终把根目录解析到了哪里

## 7.4 上层不应该做的事

业务层和 repository 不应该：

- 保存绝对路径
- 直接拼接 `filepath.Join(DataDir, ...)` 作为持久化字段
- 假定 artifact 一定落在某个固定子目录算法下

这些解析都应该由 adapter 在运行时完成。

## 7.5 路径安全要求

contracts 层虽然不负责实现路径清洗逻辑，但应定义路径安全约束：

- 相对路径必须是归一化后的安全路径
- 不允许通过 `..` 越出受管数据根目录
- 不允许把绝对路径暴露为稳定 contracts 接口的一部分
- adapter 在解析 relative path 时必须校验结果仍在目标根目录之下

---

## 8. 生命周期设计

如果 contracts 层只停留在 `Save / Read / Delete / Move / Exists`，那么很多关键问题会被埋进 actions 或 adapters 中，难以统一约束。

因此文件系统能力至少应在概念上覆盖以下生命周期阶段。

## 8.1 Staging

staging 表示 artifact 还处于临时写入或处理中，尚不能被业务对象长期引用。

适用场景：

- 导入文件刚写入磁盘但尚未校验哈希
- 卡片 HTML 生成完但还未成功关联 DB 记录
- 索引构建过程中的临时文件

约束：

- staging path 不能写入业务持久化表
- staging artifact 可以被覆盖或清理
- staging 目录可位于专门的临时空间或目标目录的临时子路径中

## 8.2 Promotion

promotion 表示将 staging artifact 提升为最终 durable artifact。

推荐约束：

- 需要长期持久化的 artifact，应先写 staging，再 promote
- promote 成功后，artifact 才可被 repository 记录为正式路径
- adapter 应优先使用原子 rename / move 能力完成 promote（平台允许时）

业务上的关键不变量：

- **DB 不能引用 staging path**
- **只有 promoted artifact 才允许进入业务引用字段**

## 8.3 Deletion

删除需要区分两个层次：

1. **逻辑删除**：业务对象不再引用该 artifact，或状态转为 deleted
2. **物理删除**：磁盘文件真正被移除

当前模型中已有类似语义：

- `Asset.Status = active / missing / deleted`

建议：

- contracts 层应允许“逻辑删除先于物理删除”
- 对于可能仍被引用、或需要做一致性检查的 artifact，不要求每次删除都立即物理清空

## 8.4 Missing / Reconciliation

missing 表示 DB 仍保留 artifact 元数据，但磁盘文件已不存在。

这类情况不应直接等同于“系统不可恢复损坏”，而应视为：

- 一种可识别的一致性问题
- 可由修复流程、重建流程、重新导入流程、补偿逻辑处理

这与现有模型是吻合的：

- `Asset.Status` 已经预留 `missing`
- `SearchIndexState` 已经预留 `pending / stale / failed` 这类异步修复语义

## 8.5 生命周期结论

文件系统 contract 的重点不是把所有状态机都编码进接口，而是先固定几个重要不变量：

- 持久化引用只能指向 promoted artifact
- staging 与正式路径必须区分
- cache artifact 可以直接清理，但 durable artifact 的删除需要考虑业务引用关系
- missing 应有统一语义，不应让上层直接消费底层 syscall 细节

---

## 9. FileStore 合同的推荐演进方式

## 9.1 当前阶段仍可保留基础 `FileStore`

`docs/contracts-design.md` 中目前建议的通用接口如下：

```go
type FileStore interface {
    Save(ctx context.Context, path string, content []byte) error
    Read(ctx context.Context, path string) ([]byte, error)
    Delete(ctx context.Context, path string) error
    Move(ctx context.Context, fromPath, toPath string) error
    Exists(ctx context.Context, path string) (bool, error)
}
```

对于第一阶段，这是可接受的最小能力面。

## 9.2 但不建议让业务层长期直接传物理 path

随着 card content、asset、source import、index 等语义变强，如果仍让上层手工传 `path string`，会带来几个问题：

- actions 开始自己拼接路径
- 目录结构与业务层耦合
- 更换布局时改动面过大
- 很难统一 staging / promotion 规则

因此更推荐的演进方向是：

- 基础层提供通用读写能力
- 语义层为不同 artifact 提供更稳定的 request/result 结构

## 9.3 推荐的两层思路

### A. 基础能力接口

用于封装底层文件操作与生命周期基础能力，例如：

- write temporary artifact
- promote artifact
- read artifact
- delete artifact
- check existence

### B. 语义化能力接口

在需要时逐步引入，例如：

- `SaveCardHTML(...)`
- `SaveSourceDocument(...)`
- `SaveAsset(...)`
- `DeleteArtifact(...)`

这样做的好处是：

- 上层表达的是业务意图
- adapter 仍保有布局与路径分配主导权
- contracts 可以在不暴露实现细节的前提下逐步进化

## 9.4 当前阶段的保守建议

当前不必一次性引入很多复杂接口。

更合适的做法是：

- 保留 `FileStore` 作为总称
- 在设计上明确：**业务层不应手拼最终 path**
- 后续需要时，再把“路径分配”“staging/promote”“按 artifact kind 保存”逐步语义化

---

## 10. 错误语义

文件系统 contract 应沿用 `docs/contracts-design.md` 中统一的错误语义，而不是单独再发明一套文件系统错误类型。

## 10.1 推荐沿用的错误语义

- `ErrUnavailable`
- `ErrInvalidInput`
- `ErrNotFound`
- `ErrConflict`

## 10.2 在文件系统场景中的含义

### `ErrNotFound`

适用场景：

- 请求读取的 artifact 不存在
- 请求删除的目标路径不存在
- repository 中存在引用，但 adapter 解析后发现物理文件缺失

### `ErrInvalidInput`

适用场景：

- 非法 artifact kind
- 非法相对路径
- 试图传入绝对路径作为 contracts 层稳定输入
- 路径越界到受管根目录之外

### `ErrConflict`

适用场景：

- promote 目标已存在，且当前策略不允许覆盖
- 一个逻辑对象尝试绑定多个互斥 artifact
- 并发写入导致目标 artifact 状态冲突

### `ErrUnavailable`

适用场景：

- 底层磁盘不可用
- 目录不可写
- 权限异常导致无法完成写入/移动
- 平台级文件系统异常阻止正常完成操作

必要时 adapter 可以保留原始错误用于日志或包装，但不应把平台 syscall 细节直接暴露给 actions。

---

## 11. 推荐目录布局

下面的结构是**推荐实现布局**，目的是帮助后续 adapter 和 flow 设计达成一致；它不是让业务层把这些路径硬编码到数据库中。

```text
data/
├── assets/
│   └── ...
├── imports/
│   └── ...
├── cardfile/
│   └── ...
├── index/
│   └── ...
└── export/
    └── ...

cache/
├── thumbnails/
│   └── ...
├── html/
│   └── ...
└── index/
    └── ...
```

### 设计说明

- `data/` 下存 durable artifact
- `cache/` 下存可重建 artifact
- 目录名应与当前 config 模块命名保持一致
- 后续即使具体文件命名算法变化，也不应影响 contracts 层的 durable/cache 与 relative-path 规则

---

## 12. 与当前代码的对齐建议

当前代码已经给出了一些非常好的基础约束，后续实现应尽量沿用，而不是重新设计另一套语义。

## 12.1 已经明确且应保留的约束

- `Asset.StoragePath` 相对 `DataDir`
- `Card.HTMLPath` / `Card.AnswerHTMLPath` 相对 `DataDir`
- config 已经统一提供 `AssetsDir`、`ImportsDir`、`CardFileDir`、`IndexDir`、`ThumbsCacheDir`、`HTMLCacheDir`、`IndexCacheDir`

## 12.2 本文档明确补齐的约束

- `SourceDocument.FilePath` 也应统一视为相对 `DataDir`
- cache path 不应作为业务模型的 canonical path 保存
- 缺失文件应被视为一致性问题，并通过统一错误语义或状态记录处理

## 12.3 后续实现应避免的方向

- 在 repository 中引入绝对路径
- 在 actions 中拼接物理路径
- 让 UI 或上层直接感知 `AssetsDir` 的真实系统绝对位置
- 让 adapter 把 `os.PathError` 原样返回给业务动作

---

## 13. 验证清单

本次任务是设计文档任务，因此验证不以运行代码为主，而以设计一致性检查为主。

## 13.1 边界检查

需要检查文档是否清楚说明了：

- contracts 与 adapters 的分工
- contracts 与 repository 的并列关系
- config 只负责根目录，不负责 artifact 生命周期
- actions / flows 负责业务编排，而不是直接做物理文件读写

## 13.2 路径策略检查

需要检查文档是否明确说明了：

- durable artifact 的业务引用路径统一相对 `DataDir`
- cache artifact 不应成为业务模型唯一真相来源
- `SourceDocument.FilePath` 的相对路径语义已补齐
- 路径越界与绝对路径泄漏被视为非法输入

## 13.3 生命周期检查

需要检查文档是否完整覆盖了：

- staging
- promotion
- deletion
- missing / reconciliation

并且是否明确了：

- 只有 promoted artifact 才能进入 DB 引用字段
- missing 是可恢复的一致性问题

## 13.4 与当前代码术语一致性检查

需要检查文档是否与现有命名一致：

- `FileStore`
- `DataDir` / `CacheDir` / `AssetsDir` / `ImportsDir` / `IndexDir` / `CardFileDir`
- `Asset.StoragePath`
- `Card.HTMLPath` / `Card.AnswerHTMLPath`
- `SourceDocument.FilePath`

---

## 14. 结论

当前 kmemo 最需要的，不是一个立即覆盖所有未来场景的复杂文件系统抽象，而是一套足够稳定、足够清晰的 contracts 设计约束：

- 用 artifact 而不是泛化的 file 来表达业务对象
- 用 durable / cache 区分真相数据与可重建产物
- 用相对 `DataDir` 的路径保证可迁移、可备份、可跨平台
- 用 staging / promotion / deletion / reconciliation 约束文件生命周期
- 让 adapter 负责物理布局与平台细节，让 actions / repository 只面对稳定语义

这样后续无论实现继续基于本地文件系统，还是扩展到更复杂的存储方式，contracts 层都能保持可读、可扩展、可理解。
