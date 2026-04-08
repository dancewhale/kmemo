# File 模块设计文档（kmemo 项目版）

## 1. 目标

`file` 模块是 kmemo 的底层文件存储模块，负责：

1. 基于对象 `Kind` 和 `UUID v7` 管理文件对象
2. 使用 4 层分桶目录结构（`h1h2h3/h4/h5/h6`）
3. 提供 trash 机制支持文件恢复
4. 实现原子写入和对象级并发控制
5. 与数据库软删除配合实现双重保护

本模块不负责：
- 业务元数据解析
- 复习算法和 SRS 状态
- 全文索引
- 业务层删除判断

---

## 2. 与 kmemo 架构的集成

### 2.1 在分层架构中的位置

```text
Actions 层
    ↓
Contracts 层 (FileStore 接口)
    ↓
Adapters 层 (file 模块实现)
    ↓
文件系统
```

### 2.2 与其他模块的关系

- **Repository 层**：保存对象身份（ID, Slug, Ext），不保存路径
- **Contracts 层**：定义 FileStore 接口
- **File 模块**：实现具体的文件对象管理
- **Config 模块**：提供 VaultDir 配置

---

## 3. 存储布局

### 3.1 目录结构

```text
vault/
├── cards/              # 卡片 HTML（active）
│   └── 019/5/f/3/      # UUID v7 分桶
│       ├── 0195f3f0-...-d45c8f23a1ee--how-memory-works.html
│       └── 0195f3f0-...-d45c8f23a1ee--how-memory-works-answer.html
├── assets/             # 资源文件（active）
│   └── 019/5/f/3/
│       └── 0195f3f0-...-d45c8f23a1ee.png
├── sources/            # 来源文档（active）
│   └── 019/5/f/3/
│       └── 0195f3f0-...-d45c8f23a1ee--book-name.pdf
├── trash/              # 回收站
│   ├── cards/
│   ├── assets/
│   └── sources/
├── tmp/                # 临时文件（原子写入）
└── locks/              # 对象锁
    ├── cards/
    ├── assets/
    └── sources/
```

### 3.2 Kind 配置

```go
type FileObjectKindSpec struct {
    Kind          FileObjectKind  // "card", "asset", "source"
    DirectoryName string           // "cards", "assets", "sources"
    FileNameStyle FileNameStyle    // WithSlug / WithoutSlug
}

// kmemo 的 Kind 配置
var KmemoKinds = []FileObjectKindSpec{
    {Kind: "card", DirectoryName: "cards", FileNameStyle: FileNameStyleWithSlug},
    {Kind: "asset", DirectoryName: "assets", FileNameStyle: FileNameStyleWithoutSlug},
    {Kind: "source", DirectoryName: "sources", FileNameStyle: FileNameStyleWithSlug},
}
```

---

## 4. 核心数据结构

### 4.1 对象身份

```go
type FileObjectKind string

const (
    FileObjectKindCard   FileObjectKind = "card"
    FileObjectKindAsset  FileObjectKind = "asset"
    FileObjectKindSource FileObjectKind = "source"
)

type FileObjectID string  // UUID v7

type FileObjectRef struct {
    Kind FileObjectKind
    ID   FileObjectID
}
```

### 4.2 文件名提示

```go
type FileObjectNameHint struct {
    Slug string  // 可读名称
    Ext  string  // 文件扩展名
}

type FileObjectLookup struct {
    Ref  FileObjectRef
    Name *FileObjectNameHint  // nil 表示只按 kind+id 查找
}
```

### 4.3 Scope

```go
type FileObjectScope string

const (
    FileObjectScopeActive FileObjectScope = "active"
    FileObjectScopeTrash  FileObjectScope = "trash"
    FileObjectScopeAny    FileObjectScope = "any"
)
```

### 4.4 位置信息

```go
type FileObjectLocation struct {
    Ref      FileObjectRef
    Scope    FileObjectScope
    Path     string    // 完整物理路径
    FileName string    // 文件名
    Slug     string
    Ext      string
    Size     int64
    ModTime  time.Time
}
```

---

## 5. 核心接口

```go
type FileStore interface {
    // 创建新对象
    CreateFileObject(ctx context.Context, input CreateFileObjectInput) (FileObjectLocation, error)

    // 覆盖现有对象
    OverwriteFileObject(ctx context.Context, input OverwriteFileObjectInput) (FileObjectLocation, error)

    // 读取对象
    ReadFileObject(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) ([]byte, FileObjectLocation, error)

    // 查找对象位置
    FindFileObject(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) (*FileObjectLocation, error)

    // 检查存在性
    FileObjectExists(ctx context.Context, lookup FileObjectLookup, scope FileObjectScope) (bool, error)

    // Trash 操作
    MoveFileObjectToTrash(ctx context.Context, ref FileObjectRef) (FileObjectLocation, error)
    RestoreFileObjectFromTrash(ctx context.Context, ref FileObjectRef) (FileObjectLocation, error)
    PermanentlyDeleteFileObject(ctx context.Context, ref FileObjectRef) error

    // Slug 重命名
    RenameFileObjectSlug(ctx context.Context, ref FileObjectRef, newSlug string) (FileObjectLocation, error)

    // 扫描对象
    ScanFileObjects(ctx context.Context, options ScanFileObjectsOptions) (<-chan FileObjectLocation, <-chan error)
}
```

---

## 6. 与数据库模型的配合

### 6.1 Asset 模型

```go
type Asset struct {
    ID       string  // UUID v7
    Slug     string  // 可选，用于可读文件名
    Ext      string  // 文件扩展名
    Checksum string
    Size     int64
    Status   string  // active, missing, deleted
}

// 使用示例
loc, err := fileStore.CreateFileObject(ctx, CreateFileObjectInput{
    Kind: "asset",
    ID:   asset.ID,
    Slug: asset.Slug,
    Ext:  asset.Ext,
    Data: content,
})
```

### 6.2 Card 模型

```go
type Card struct {
    ID       string  // UUID v7
    Slug     string  // 从 Title 派生
    HTMLHash string
    AnswerHTMLHash string
}

// 正文 HTML
loc, _ := fileStore.CreateFileObject(ctx, CreateFileObjectInput{
    Kind: "card",
    ID:   card.ID,
    Slug: card.Slug,
    Ext:  "html",
    Data: htmlContent,
})

// 答案 HTML
loc, _ := fileStore.CreateFileObject(ctx, CreateFileObjectInput{
    Kind: "card",
    ID:   card.ID,
    Slug: card.Slug + "-answer",
    Ext:  "html",
    Data: answerContent,
})
```

### 6.3 SourceDocument 模型

```go
type SourceDocument struct {
    ID           string  // UUID v7
    Slug         string  // 从文件名派生
    Ext          string  // 原始文件扩展名
    OriginalHash string
}
```

---

## 7. 删除策略（混合方案）

### 7.1 软删除流程

```go
// 步骤 1：数据库标记删除
asset.Status = "deleted"
repo.Asset.Update(ctx, asset)

// 步骤 2：文件移动到 trash
loc, err := fileStore.MoveFileObjectToTrash(ctx, FileObjectRef{
    Kind: "asset",
    ID:   asset.ID,
})
```

### 7.2 恢复流程

```go
// 步骤 1：恢复文件
loc, err := fileStore.RestoreFileObjectFromTrash(ctx, FileObjectRef{
    Kind: "asset",
    ID:   asset.ID,
})

// 步骤 2：数据库恢复
asset.Status = "active"
repo.Asset.Update(ctx, asset)
```

### 7.3 永久删除

```go
// 定期清理任务
func CleanupTrash(ctx context.Context, olderThan time.Duration) error {
    // 1. 查询数据库中已删除超过 N 天的记录
    assets, _ := repo.Asset.List(ctx, ListAssetOptions{
        Status: "deleted",
        DeletedBefore: time.Now().Add(-olderThan),
    })

    // 2. 永久删除文件
    for _, asset := range assets {
        fileStore.PermanentlyDeleteFileObject(ctx, FileObjectRef{
            Kind: "asset",
            ID:   asset.ID,
        })
    }

    // 3. 删除数据库记录
    repo.Asset.HardDelete(ctx, assetIDs)
}
```

---

## 8. 使用示例

### 8.1 创建卡片

```go
func (a *CreateCardAction) Execute(ctx context.Context, req CreateCardRequest) error {
    card := &models.Card{
        ID:    generateUUIDv7(),
        Slug:  slugify(req.Title),
        Title: req.Title,
    }

    // 保存 HTML
    loc, err := a.fileStore.CreateFileObject(ctx, CreateFileObjectInput{
        Kind: "card",
        ID:   card.ID,
        Slug: card.Slug,
        Ext:  "html",
        Data: req.HTMLContent,
    })
    if err != nil {
        return err
    }

    // 计算哈希
    card.HTMLHash = computeHash(req.HTMLContent)

    // 保存到数据库
    return a.cardRepo.Create(ctx, card)
}
```

### 8.2 读取资源

```go
func (a *GetAssetAction) Execute(ctx context.Context, assetID string) ([]byte, error) {
    // 从数据库获取元数据
    asset, err := a.assetRepo.GetByID(ctx, assetID)
    if err != nil {
        return nil, err
    }

    // 读取文件
    data, loc, err := a.fileStore.ReadFileObject(ctx, FileObjectLookup{
        Ref: FileObjectRef{Kind: "asset", ID: asset.ID},
        Name: &FileObjectNameHint{Slug: asset.Slug, Ext: asset.Ext},
    }, FileObjectScopeActive)

    return data, err
}
```

### 8.3 更新卡片标题（Slug 变化）

```go
func (a *UpdateCardAction) Execute(ctx context.Context, cardID string, newTitle string) error {
    card, _ := a.cardRepo.GetByID(ctx, cardID)
    newSlug := slugify(newTitle)

    if card.Slug != newSlug {
        // 重命名文件
        loc, err := a.fileStore.RenameFileObjectSlug(ctx, FileObjectRef{
            Kind: "card",
            ID:   cardID,
        }, newSlug)
        if err != nil {
            return err
        }

        // 更新数据库
        card.Slug = newSlug
    }

    card.Title = newTitle
    return a.cardRepo.Update(ctx, card)
}
```

---

## 9. 实现要点

### 9.1 分桶路径计算

```go
func BuildBucketPath(id FileObjectID) string {
    // 移除连字符
    hex := strings.ReplaceAll(string(id), "-", "")
    // 019/5/f/3
    return fmt.Sprintf("%s/%s/%s/%s",
        hex[0:3], hex[3:4], hex[4:5], hex[5:6])
}
```

### 9.2 原子写入

```go
func WriteFileAtomically(targetPath string, content []byte) error {
    dir := filepath.Dir(targetPath)
    tmpFile, _ := os.CreateTemp(dir, ".tmp-*")
    defer os.Remove(tmpFile.Name())

    tmpFile.Write(content)
    tmpFile.Sync()
    tmpFile.Close()

    return os.Rename(tmpFile.Name(), targetPath)
}
```

### 9.3 对象锁

```go
func (s *fileStore) WithObjectLock(ctx context.Context, ref FileObjectRef, fn func() error) error {
    lockPath := s.buildLockPath(ref)

    // 尝试创建锁目录
    for {
        err := os.Mkdir(lockPath, 0755)
        if err == nil {
            break  // 获取锁成功
        }
        select {
        case <-ctx.Done():
            return ErrLockTimeout
        case <-time.After(100 * time.Millisecond):
            // 重试
        }
    }
    defer os.Remove(lockPath)

    return fn()
}
```

---

## 10. 配置

```go
type FileStoreConfig struct {
    RootDir   string                    // vault 根目录
    Kinds     []FileObjectKindSpec      // Kind 配置
    LockWait  time.Duration             // 锁等待超时
    LockRetry time.Duration             // 锁重试间隔
}

// kmemo 默认配置
func DefaultConfig(vaultDir string) FileStoreConfig {
    return FileStoreConfig{
        RootDir:   vaultDir,
        Kinds:     KmemoKinds,
        LockWait:  30 * time.Second,
        LockRetry: 100 * time.Millisecond,
    }
}
```

---

## 11. 总结

本设计为 kmemo 提供了：

1. **完整的文件对象管理**：创建、读取、更新、删除、恢复
2. **双重删除保护**：数据库软删除 + 文件 trash 机制
3. **派生路径策略**：路径由 kind+id+slug+ext 计算，数据库不存路径
4. **分桶存储**：UUID v7 + 4 层分桶，支持大规模文件
5. **并发安全**：对象级锁保护并发修改

与现有架构的集成点：
- **Contracts 层**：定义 FileStore 接口
- **Repository 层**：保存 ID/Slug/Ext，不保存路径
- **Actions 层**：使用 FileStore 管理文件生命周期
