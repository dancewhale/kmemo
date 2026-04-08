下面是一份可直接交给代码生成模型的 **Go 版 `file` 模块设计文档**。
我尽量把接口、行为约束、并发语义、错误语义、命名约定都写清楚了，并结合了我们前面的讨论与这次你的修正：

* 语言是 **Golang**
* 函数和对象命名更强调**可读性和可理解性**
* 获取文件对象支持两种方式：

  * `kind + id`
  * `kind + id + slug + ext`
* `kind` 支持后续扩展，不写死在架构里
* 不使用 side index
* 使用 UUID v7
* 分桶规则固定为 `h1h2h3 / h4 / h5 / h6`

---

# File 模块设计文档

## 1. 目标

`file` 模块是系统的底层文件存储模块，负责：

1. 基于对象 `Kind` 和 `UUID v7` 计算规范分桶目录
2. 基于规范命名规则定位、创建、覆盖、移动、恢复和删除对象文件
3. 提供原子写入能力，避免正式文件出现半写入状态
4. 在无 side index 的前提下，通过扫描目标分桶目录解析对象文件
5. 提供对象级并发写保护，避免同一对象的并发修改冲突
6. 保持模块边界清晰，不承担 HTML 语义解析、引用图构建、复习调度等上层职责

本模块不负责：

* 解析 HTML/JSON 内容中的业务元数据
* 复习算法和间隔重复状态
* 引用关系检查
* 全文索引
* 对象是否“允许删除”的业务判断

---

## 2. 设计原则

### 2.1 身份与路径分离

对象的唯一身份由：

* `Kind`
* `UUID v7`

共同确定。

文件路径、文件名中的 slug、扩展名都不是对象身份的一部分。

### 2.2 路径是派生值

路径不作为主键保存和传递。
模块通过 `Kind + UUID v7` 推导分桶目录，再在该目录中定位对象文件。

### 2.3 不使用 side index

模块不维护 `id -> path` 的持久化或内存索引。
所有定位都直接基于文件系统进行：

1. 计算分桶目录
2. 扫描该目录
3. 解析文件名
4. 定位目标对象

### 2.4 原子写优先于性能

所有正式写入都必须通过临时文件 + `fsync` + `rename` 完成，避免正式文件处于损坏状态。

### 2.5 修改操作必须串行化

同一对象 `(kind, id)` 的所有修改型操作必须加对象级独占锁。

### 2.6 Kind 可扩展

`Kind` 不应被写死为固定 enum。
模块应允许未来增加新的对象类别，例如：

* `card`
* `source`
* `asset`
* `snapshot`
* `export`
* `note_bundle`

---

## 3. 存储布局

根目录示例：

```text
vault/
  cards/
  sources/
  assets/
  trash/
  tmp/
  locks/
```

由于 `Kind` 可扩展，目录布局不应完全写死到具体类型名。建议定义：

* Active 根目录：`<root>/<kindPlural>/...`
* Trash 根目录：`<root>/trash/<kindPlural>/...`
* Lock 根目录：`<root>/locks/<kindPlural>/...`

其中 `kindPlural` 由 Kind 的存储配置决定，例如：

* `card` -> `cards`
* `source` -> `sources`
* `asset` -> `assets`

Kind 的 plural/path segment 建议通过配置提供，而不是写死在代码里。

---

## 4. UUID v7 和分桶规则

### 4.1 UUID 格式

模块只接受合法的 **UUID v7**。

规范要求：

* 小写
* 标准连字符格式，例如：

```text
0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee
```

### 4.2 分桶规则

分桶基于 UUID 去掉连字符后的 32 位 hex 字符串：

```text
0195f3f0a8b77c8db1b9d45c8f23a1ee
```

目录层级固定为：

* 第 1 层：`h1h2h3`
* 第 2 层：`h4`
* 第 3 层：`h5`
* 第 4 层：`h6`

示例：

```text
019/5/f/3
```

### 4.3 例子

UUID：

```text
0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee
```

对应 bucket：

```text
019/5/f/3
```

card active 目录：

```text
<root>/cards/019/5/f/3/
```

asset trash 目录：

```text
<root>/trash/assets/019/5/f/3/
```

---

## 5. 文件命名规则

## 5.1 通用说明

文件名规则由 Kind 决定，但应遵守：

* 身份字段必须包含 UUID v7
* slug 仅用于可读性
* 扩展名由调用方显式指定
* 文件名必须可逆解析出至少：

  * kind
  * id
  * slug（若存在）
  * ext

## 5.2 默认推荐命名规则

### 5.2.1 Card-like / Source-like 对象

适用于需要可读 slug 的对象：

```text
<uuid>--<slug>.<ext>
```

例如：

```text
0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee--how-memory-works.html
```

### 5.2.2 Asset-like 对象

适用于不需要 slug 的对象：

```text
<uuid>.<ext>
```

例如：

```text
0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee.png
```

## 5.3 Kind 的命名策略可配置

由于 Kind 可扩展，建议定义每种 Kind 的文件命名策略：

```go
type FileNameStyle int

const (
    FileNameStyleWithSlug FileNameStyle = iota
    FileNameStyleWithoutSlug
)
```

每种 Kind 在注册时声明使用哪种命名风格。

---

## 6. Slug 规则

### 6.1 作用

slug 仅用于可读性，不参与对象身份判定。

### 6.2 规范化规则

建议统一 slugify 规则：

* 转小写
* trim 空白
* 空白转 `-`
* 去除危险字符
* 合并连续 `-`
* 最长 64 字符
* 空结果时使用 `untitled`

### 6.3 重要约束

* slug 可变
* slug 改变不影响对象身份
* 通过 `kind + id` 获取对象时，模块不能依赖 slug 必须正确

---

## 7. 核心数据结构

以下为推荐 Go 结构体设计。

## 7.1 Kind 定义

```go
type FileObjectKind string
```

示例：

```go
const (
    FileObjectKindCard   FileObjectKind = "card"
    FileObjectKindSource FileObjectKind = "source"
    FileObjectKindAsset  FileObjectKind = "asset"
)
```

但实现中不应假设 Kind 只有这几个。

---

## 7.2 Kind 配置

```go
type FileObjectKindSpec struct {
    Kind          FileObjectKind
    DirectoryName string
    FileNameStyle FileNameStyle
}
```

说明：

* `Kind`：对象类别
* `DirectoryName`：active 目录名，例如 `cards`
* `FileNameStyle`：命名风格，带 slug 或不带 slug

---

## 7.3 对象身份

```go
type FileObjectID string
```

它应当是已经规范化后的 UUID v7 字符串。

```go
type FileObjectRef struct {
    Kind FileObjectKind
    ID   FileObjectID
}
```

---

## 7.4 精确文件匹配信息

由于你要求支持通过 `kind + id + slug + ext` 直接获取，建议增加一个更具体的结构：

```go
type FileObjectNameHint struct {
    Slug string
    Ext  string
}
```

配合：

```go
type FileObjectLookup struct {
    Ref  FileObjectRef
    Name *FileObjectNameHint
}
```

语义：

* `Name == nil`：按 `kind + id` 查找
* `Name != nil`：按 `kind + id + slug + ext` 直接构造候选文件名并优先检查

---

## 7.5 写入参数

```go
type CreateFileObjectInput struct {
    Kind FileObjectKind
    ID   FileObjectID
    Slug string
    Ext  string
    Data []byte
}
```

```go
type OverwriteFileObjectInput struct {
    Kind FileObjectKind
    ID   FileObjectID
    Slug string
    Ext  string
    Data []byte
    IfMatch *FileObjectMatchCondition
}
```

---

## 7.6 并发匹配条件

```go
type FileObjectMatchCondition struct {
    ExpectedModTimeUnixNano int64
    ExpectedContentHash     string
}
```

用于解决 lost update 问题。

---

## 7.7 定位结果

```go
type FileObjectLocation struct {
    Ref      FileObjectRef
    Scope    FileObjectScope
    Path     string
    FileName string
    Slug     string
    Ext      string
    Size     int64
    ModTime  time.Time
}
```

---

## 7.8 Scope

```go
type FileObjectScope string

const (
    FileObjectScopeActive FileObjectScope = "active"
    FileObjectScopeTrash  FileObjectScope = "trash"
    FileObjectScopeAny    FileObjectScope = "any"
)
```

---

## 8. 对外接口设计

建议暴露一个主接口：

```go
type FileStore interface {
    CreateFileObject(ctx context.Context, input CreateFileObjectInput) (FileObjectLocation, error)
    OverwriteFileObject(ctx context.Context, input OverwriteFileObjectInput) (FileObjectLocation, error)

    ReadFileObject(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) ([]byte, FileObjectLocation, error)
    FindFileObject(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) (*FileObjectLocation, error)
    FileObjectExists(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) (bool, error)

    MoveFileObjectToTrash(ctx context.Context, ref FileObjectRef) (FileObjectLocation, error)
    RestoreFileObjectFromTrash(ctx context.Context, ref FileObjectRef) (FileObjectLocation, error)
    PermanentlyDeleteFileObject(ctx context.Context, ref FileObjectRef) error

    RenameFileObjectSlug(ctx context.Context, ref FileObjectRef, newSlug string) (FileObjectLocation, error)

    ScanFileObjects(ctx context.Context, options ScanFileObjectsOptions) (<-chan FileObjectLocation, <-chan error)
}
```

---

# 9. 接口语义说明

## 9.1 `CreateFileObject`

```go
CreateFileObject(ctx, input)
```

### 用途

创建一个新的对象文件。

### 输入要求

* `Kind` 必须已注册
* `ID` 必须是合法 UUID v7
* `Ext` 必须非空且合法
* `Slug` 对于需要 slug 的 Kind 可以为空，空时内部转为 `untitled`
* `Data` 必须是完整文件内容

### 行为

1. 获取 `(kind, id)` 对象锁
2. 检查 active scope 下是否已存在该对象
3. 检查 trash scope 下是否已存在该对象
4. 若任一存在，则报错
5. 计算规范文件名和 active 目录
6. 确保目录存在
7. 原子写入
8. 返回最终 `FileObjectLocation`

### 错误

* `ErrInvalidKind`
* `ErrInvalidUUIDv7`
* `ErrInvalidFileExtension`
* `ErrFileObjectAlreadyExists`
* `ErrFileObjectAlreadyExistsInTrash`
* `ErrLockTimeout`
* `ErrIO`

---

## 9.2 `OverwriteFileObject`

```go
OverwriteFileObject(ctx, input)
```

### 用途

覆盖现有 active 对象文件内容。

### 输入要求

* 对象必须存在于 active scope
* 默认不允许修改 ext
* 允许修改 slug
* 可提供 `IfMatch` 进行并发保护

### 行为

1. 获取 `(kind, id)` 对象锁
2. 在 active scope 下定位对象
3. 若不存在，报错
4. 若有多个同 id 文件，报错
5. 若 ext 与现有对象不一致，报错
6. 若设置了 `IfMatch`，校验当前对象是否匹配
7. 若 slug 变化：

   * 先将现有文件重命名为新 slug 文件名
   * 再对新文件路径执行原子内容替换
8. 若 slug 不变：

   * 直接原子替换内容
9. 返回新的 `FileObjectLocation`

### 并发语义

* 在对象锁保护下串行执行
* 默认 `last successful writer wins`
* 若设置 `IfMatch`，可防止覆盖最新修改

### 错误

* `ErrFileObjectNotFound`
* `ErrDuplicateFileObjectID`
* `ErrFileExtensionMismatch`
* `ErrConcurrentModification`
* `ErrLockTimeout`
* `ErrIO`

---

## 9.3 `FindFileObject`

```go
FindFileObject(ctx, lookup, scope)
```

### 用途

查找对象文件位置和文件信息，但不读取内容。

### 支持两种查找方式

#### 方式 A：按 `kind + id`

`lookup.Name == nil`

流程：

1. 计算目标 scope 的 bucket 目录
2. 扫描目录
3. 解析文件名
4. 匹配 `kind + id`
5. 返回唯一对象

#### 方式 B：按 `kind + id + slug + ext`

`lookup.Name != nil`

流程：

1. 计算 bucket 目录
2. 直接构造期望文件名
3. 检查该文件是否存在
4. 若存在，直接返回
5. 若不存在，可选择是否回退到目录扫描

### 建议行为

为提高可理解性，建议：

* 当 `lookup.Name != nil` 时，先尝试精确文件名定位
* 若精确文件不存在，则**回退到按 `kind + id` 目录扫描**
* 这样即便 slug 已变化，也不影响对象查找鲁棒性

### scope 规则

* `active`：只查 active
* `trash`：只查 trash
* `any`：先查 active，后查 trash

### 返回

* 找到：返回 `*FileObjectLocation`
* 找不到：返回 `nil, nil`

### 错误

* `ErrInvalidKind`
* `ErrInvalidUUIDv7`
* `ErrDuplicateFileObjectID`
* `ErrIO`

---

## 9.4 `ReadFileObject`

```go
ReadFileObject(ctx, lookup, scope)
```

### 用途

读取对象完整文件内容和文件信息。

### 行为

1. 调用 `FindFileObject`
2. 若未找到，报 `ErrFileObjectNotFound`
3. 读取完整文件内容
4. 返回内容和 location

### 并发语义

* 默认无锁
* 依赖原子写保证不会读到半截文件
* 读者可能读到旧版本或新版本，但不会读到损坏版本

---

## 9.5 `FileObjectExists`

```go
FileObjectExists(ctx, lookup, scope)
```

### 用途

判断对象是否存在。

### 行为

内部可复用 `FindFileObject`，返回布尔值。

---

## 9.6 `MoveFileObjectToTrash`

```go
MoveFileObjectToTrash(ctx, ref)
```

### 用途

将 active 对象移动到 trash。

### 行为

1. 获取对象锁
2. 在 active scope 下定位对象
3. 若不存在，报错
4. 检查 trash scope 是否已存在同 id 对象
5. 若存在，报冲突
6. 确保 trash 目标目录存在
7. 使用 `rename` 移动文件到 trash
8. 返回 trash 中的位置

### 错误

* `ErrFileObjectNotFound`
* `ErrTrashConflict`
* `ErrDuplicateFileObjectID`
* `ErrLockTimeout`
* `ErrIO`

---

## 9.7 `RestoreFileObjectFromTrash`

```go
RestoreFileObjectFromTrash(ctx, ref)
```

### 用途

从 trash 恢复对象。

### 行为

1. 获取对象锁
2. 在 trash scope 下定位对象
3. 若不存在，报错
4. 检查 active scope 是否已有同 id 对象
5. 若存在，报冲突
6. 恢复到 active 对应 bucket 目录
7. 保留 trash 中当前文件名
8. 返回 active 中的位置

### 注意

恢复的是对象，不是历史路径。
目标路径由当前文件名和 active 目录共同决定。

---

## 9.8 `PermanentlyDeleteFileObject`

```go
PermanentlyDeleteFileObject(ctx, ref)
```

### 用途

永久删除 trash 中对象。

### 行为

1. 获取对象锁
2. 在 trash scope 下定位对象
3. 若不存在，报错
4. 删除文件
5. 不允许直接删除 active 对象

### 错误

* `ErrFileObjectNotFound`
* `ErrLockTimeout`
* `ErrIO`

---

## 9.9 `RenameFileObjectSlug`

```go
RenameFileObjectSlug(ctx, ref, newSlug)
```

### 用途

仅修改文件名中的 slug，不改文件内容。

### 行为

1. 获取对象锁
2. 仅在 active scope 下定位对象
3. 若 Kind 不支持 slug，报错
4. 规范化 `newSlug`
5. 若 slug 未变化，直接返回当前位置
6. 计算新文件名
7. 使用原子 `rename` 改名
8. 返回新位置

### 约束

* 不修改 ext
* 不修改内容
* 不操作 trash 对象

---

## 9.10 `ScanFileObjects`

```go
ScanFileObjects(ctx, options)
```

### 用途

扫描 active/trash 中的对象文件。

### 建议 options

```go
type ScanFileObjectsOptions struct {
    Kind  *FileObjectKind
    Scope FileObjectScope
}
```

### 行为

* 遍历指定 Kind 和 Scope 下的目录树
* 仅输出符合命名规则的对象文件
* 不解析文件内容
* 不检查文件内容里的业务元数据

### 输出形式

Go 推荐返回 channel：

* 一个 `location` channel
* 一个 `error` channel

或者也可以直接定义 callback 风格。
如果追求实现简单，也可直接返回切片，但不适合大规模扫描。

---

# 10. 定位逻辑细节

## 10.1 按 `kind + id` 定位

给定：

* `kind`
* `id`

模块执行：

1. 计算 bucket 目录
2. 列目录
3. 对每个文件执行 `ParseFileName`
4. 筛选 `id == target id`
5. 若 0 个 -> not found
6. 若 1 个 -> success
7. 若 >1 个 -> `ErrDuplicateFileObjectID`

## 10.2 按 `kind + id + slug + ext` 定位

给定：

* `kind`
* `id`
* `slug`
* `ext`

模块执行：

1. 规范化 slug
2. 构造目标文件名
3. 直接检查 `<bucket>/<filename>` 是否存在
4. 若存在，直接返回
5. 若不存在，回退到按 `kind + id` 扫描

这样可以兼顾：

* 快速精确定位
* 对 slug 变化的容错能力

---

# 11. 文件名解析规则

建议实现统一函数：

```go
ParseFileName(kindSpec FileObjectKindSpec, fileName string) (ParsedFileName, bool)
```

```go
type ParsedFileName struct {
    ID   FileObjectID
    Slug string
    Ext  string
}
```

### 解析要求

* 必须严格匹配该 Kind 的命名规则
* UUID 必须合法且为 v7
* 扩展名必须非空
* 不允许模糊猜测

---

# 12. 原子写实现规范

## 12.1 目标

保证正式文件路径下永远只有完整版本，避免出现半截文件。

## 12.2 写入步骤

对目标文件路径 `<target>`：

1. 在 `<target>` 同目录创建临时文件
2. 向临时文件写入完整内容
3. `fsync` 临时文件
4. 关闭临时文件
5. `rename(temp, target)`
6. `fsync` 目标目录（平台允许时）

### 重要说明

临时文件必须放在目标目录所在文件系统中，避免跨文件系统导致 `rename` 退化。

## 12.3 slug 变化时的 overwrite

由于无 side index，且同目录中同 id 文件不得并存，因此 overwrite 时若 slug 变化，必须：

1. 先在锁内将旧文件 `rename` 到新 slug 文件名
2. 再对新路径执行原子内容替换

不能采用“先写新文件再删旧文件”的方式，否则会短暂出现同 id 双文件。

---

# 13. 并发模型

## 13.1 总原则

同一对象 `(kind, id)` 的所有修改型操作必须串行。

### 修改型操作包括：

* `CreateFileObject`
* `OverwriteFileObject`
* `MoveFileObjectToTrash`
* `RestoreFileObjectFromTrash`
* `PermanentlyDeleteFileObject`
* `RenameFileObjectSlug`

### 读取型操作默认无锁：

* `FindFileObject`
* `ReadFileObject`
* `FileObjectExists`
* `ScanFileObjects`

---

## 13.2 锁粒度

对象锁粒度为：

```text
(kind, id)
```

不同对象可并发操作。

---

## 13.3 锁实现建议

建议使用锁目录或锁文件，优先考虑跨进程安全。

推荐锁路径：

```text
<root>/locks/<kindPlural>/<bucket>/<uuid>.lock
```

例如：

```text
vault/locks/cards/019/5/f/3/0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee.lock
```

### 实现建议

* 优先采用目录锁：`mkdir(lockdir)` 成功即获取锁
* 锁获取失败时轮询等待，直到超时
* 超时返回 `ErrLockTimeout`

### 第一版建议

* 不做自动 stale lock 抢占
* stale lock 由运维工具或恢复工具清理

---

## 13.4 Lost update 处理

默认行为：

* 同一对象的 overwrite 在锁内串行执行
* 默认 last successful writer wins

为避免覆盖别人刚写入的版本，支持 `IfMatch`：

* `ExpectedModTimeUnixNano`
* `ExpectedContentHash`

若当前对象不匹配预期，则返回 `ErrConcurrentModification`

---

# 14. 唯一性约束

以下约束必须成立：

## 14.1 同一 scope 内唯一

对任意 `(kind, id, scope)`，最多只能有一个文件对象。

## 14.2 active 和 trash 不得同时存在

同一 `(kind, id)` 不得在 active 和 trash 中同时都存在。
否则视为异常状态。

## 14.3 出现多个匹配文件即为错误

如果扫描目录时发现同 id 的多个文件，必须返回 `ErrDuplicateFileObjectID`，不能猜测哪一个有效。

---

# 15. 错误模型

建议定义统一错误类型。

```go
var (
    ErrInvalidKind                 = errors.New("invalid file object kind")
    ErrInvalidUUIDv7               = errors.New("invalid uuid v7")
    ErrInvalidFileExtension        = errors.New("invalid file extension")
    ErrInvalidFileName             = errors.New("invalid file name")

    ErrFileObjectAlreadyExists     = errors.New("file object already exists")
    ErrFileObjectAlreadyExistsInTrash = errors.New("file object already exists in trash")
    ErrFileObjectNotFound          = errors.New("file object not found")

    ErrDuplicateFileObjectID       = errors.New("duplicate file object id")
    ErrFileExtensionMismatch       = errors.New("file extension mismatch")
    ErrConcurrentModification      = errors.New("concurrent modification")

    ErrTrashConflict               = errors.New("trash conflict")
    ErrRestoreConflict             = errors.New("restore conflict")
    ErrUnsupportedOperation        = errors.New("unsupported operation")

    ErrLockTimeout                 = errors.New("lock timeout")
    ErrIO                          = errors.New("io error")
)
```

建议同时用 `fmt.Errorf("...: %w", ErrXxx)` 包装细节。

---

# 16. Kind 扩展机制

由于 `Kind` 要支持未来扩展，建议 `FileStore` 初始化时显式注册 Kind 配置。

```go
type FileStoreConfig struct {
    RootDir   string
    Kinds     []FileObjectKindSpec
    LockWait  time.Duration
    LockRetry time.Duration
}
```

初始化时建立 `map[FileObjectKind]FileObjectKindSpec`。

模块通过 KindSpec 获取：

* active 目录名
* trash 目录名
* 文件命名风格

### 这样做的好处

* 新增 Kind 不用改核心逻辑
* 只需注册新 KindSpec
* 保持定位、扫描、读写逻辑统一

---

# 17. 推荐内部辅助函数命名

为提高可读性，建议 Go 中使用清晰直白的函数命名。

## UUID / Slug

* `NormalizeUUIDv7`
* `ValidateUUIDv7`
* `NormalizeSlug`

## Kind / Path

* `GetKindSpec`
* `BuildBucketPathFromObjectID`
* `BuildActiveObjectDirectory`
* `BuildTrashObjectDirectory`
* `BuildLockDirectoryPath`

## File name

* `BuildObjectFileName`
* `BuildExpectedObjectFilePath`
* `ParseObjectFileName`

## Locate

* `FindObjectFileByID`
* `FindObjectFileByNameHint`
* `LocateSingleObjectFile`

## Write / Move

* `WriteFileAtomically`
* `ReplaceFileAtomically`
* `MoveObjectFileToTrash`
* `RestoreObjectFileFromTrash`
* `RenameObjectFileSlug`

## Lock

* `AcquireObjectLock`
* `ReleaseObjectLock`
* `WithObjectLock`

这些名字比过短或过抽象的名字更适合后续维护和 AI 生成代码。

---

# 18. 推荐实现结构

建议目录结构：

```text
file/
  store.go
  config.go
  kinds.go
  ids.go
  slug.go
  paths.go
  filenames.go
  locate.go
  atomic_write.go
  locks.go
  errors.go
  scan.go
```

建议职责：

* `store.go`：`FileStore` 主实现
* `config.go`：配置和初始化
* `kinds.go`：KindSpec 管理
* `ids.go`：UUID v7 校验和规范化
* `slug.go`：slug 规范化
* `paths.go`：目录和 bucket 计算
* `filenames.go`：命名和解析
* `locate.go`：对象定位逻辑
* `atomic_write.go`：原子写实现
* `locks.go`：对象锁实现
* `errors.go`：错误定义
* `scan.go`：扫描逻辑

---

# 19. 非目标与边界声明

以下内容不应进入本模块：

1. 解析 HTML 中的 card metadata
2. 验证文件内容内的 `id` 是否与文件名一致
3. 构建引用图
4. 检查 asset 是否仍被引用
5. 复习历史和调度状态
6. 自动垃圾回收策略
7. side index 或数据库索引维护

这些应由上层模块完成。

---

# 20. 给代码生成模型的实现要求摘要

下面这段可以直接作为生成代码时的附加要求：

---

## 代码实现要求

请使用 **Golang** 实现一个 `file` 模块，要求如下：

1. 提供 `FileStore` 接口及其默认实现
2. 使用 `UUID v7` 作为对象 ID，必须严格校验并规范化为小写标准格式
3. 不使用 side index，不依赖数据库
4. 分桶规则固定为：`h1h2h3 / h4 / h5 / h6`
5. 支持扩展型 `Kind`，通过 Kind 配置注册
6. 文件定位支持两种方式：

   * `kind + id`
   * `kind + id + slug + ext`
7. 默认策略：

   * 先按精确文件名查找
   * 若失败再回退到按 `kind + id` 扫描目录
8. 实现对象级锁，锁粒度为 `(kind, id)`，建议基于锁目录
9. 所有修改型操作必须持锁
10. 所有正式写入必须原子写
11. `OverwriteFileObject` 支持 `IfMatch` 并发保护
12. slug 变化时不得制造同 id 双文件并存
13. 代码命名优先可读性，避免过短缩写
14. 代码中应包含必要注释和清晰错误包装
15. 不要实现业务层逻辑，只实现文件层职责

---

如果你愿意，我下一步可以继续帮你把这份文档再压缩成一版 **更适合直接喂给代码模型的“实现提示词”**，或者继续补一版 **更贴近 Go 的接口定义和伪代码骨架**。

