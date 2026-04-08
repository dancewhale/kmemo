# File 模块实现方案

## 1. 概述

本文档描述如何实现基于 file.md 的新文件存储设计。由于项目处于开发阶段，无正式用户数据，因此采用直接替换策略，不考虑向后兼容。

### 1.1 实现目标

- 采用 UUID v7 + 4层分桶的目录结构
- 实现文件系统层面的 trash 机制（与数据库软删除配合）
- 使用派生路径策略（路径由 kind+id+slug+ext 计算得出）
- 提供完整的文件对象生命周期管理接口

### 1.2 实现原则

- **直接替换**：不保留旧设计，直接实现新架构
- **清理彻底**：移除所有旧的路径存储字段
- **一次到位**：所有代码同步更新，无兼容层

---

## 2. 核心变更点

### 2.1 目录结构

#### 新目录结构
```
vault/
├── cards/              # 卡片对象（active）
│   └── 019/5/f/3/      # UUID v7 分桶
│       └── 0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee--title.html
├── assets/             # 资源对象（active）
│   └── 019/5/f/3/
│       └── 0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee.png
├── sources/            # 来源文档（active）
│   └── 019/5/f/3/
│       └── 0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee--book-name.pdf
├── trash/              # 回收站
│   ├── cards/
│   │   └── 019/5/f/3/
│   ├── assets/
│   │   └── 019/5/f/3/
│   └── sources/
│       └── 019/5/f/3/
├── tmp/                # 临时文件（staging）
└── locks/              # 对象锁
    ├── cards/
    ├── assets/
    └── sources/
```

**关键特性**：
- 使用 `vault/` 作为根目录（可配置）
- 4 层分桶：`h1h2h3/h4/h5/h6`
- `trash/` 目录用于文件级回收站
- `tmp/` 目录用于原子写入
- `locks/` 目录用于并发控制

### 2.2 路径存储策略

数据库只存储身份信息，路径由模块计算：
```go
type Asset struct {
    ID   string  // UUID v7
    Slug string  // 可选，用于可读性
    Ext  string  // 文件扩展名
}

type Card struct {
    ID   string
    Slug string
}
```

**路径计算示例**：
```go
// 输入：kind="card", id="0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee", slug="how-memory-works", ext="html"
// 输出：vault/cards/019/5/f/3/0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee--how-memory-works.html
```

### 2.3 删除策略（混合方案）
数据库软删除 + 文件 trash 机制：
```go
// 步骤 1：数据库标记删除
asset.Status = "deleted"

// 步骤 2：文件移动到 trash
fileStore.MoveFileObjectToTrash(ctx, FileObjectRef{
    Kind: "asset",
    ID:   asset.ID,
})

// 后续可恢复
fileStore.RestoreFileObjectFromTrash(ctx, ref)

// 或永久删除
fileStore.PermanentlyDeleteFileObject(ctx, ref)
```

**优势**：
- 数据库和文件系统双重保护
- 支持文件级恢复（即使数据库记录已清理）
- 可独立清理 trash（定期任务）

### 2.4 接口设计
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

## 3. 数据库模型调整

### 3.1 Asset 模型

```go
type Asset struct {
    ID          string  // UUID v7
    Slug        string  // 可选，用于可读文件名
    Ext         string  // 文件扩展名
    Checksum    string
    Size        int64
    MimeType    string
    Status      string  // active, missing, deleted
}
```

**变更**：
- 移除 `StoragePath` 字段
- 新增 `Slug` 和 `Ext` 字段

### 3.2 Card 模型

```go
type Card struct {
    ID              string  // UUID v7
    Slug            string  // 从 Title 派生
    HTMLHash        string
    AnswerHTMLHash  string
}
```

**变更**：
- 移除 `HTMLPath` 和 `AnswerHTMLPath` 字段
- 新增 `Slug` 字段

**路径计算规则**：
- 正文 HTML：`kind="card", id=card.ID, slug=card.Slug, ext="html"`
- 答案 HTML：`kind="card", id=card.ID, slug=card.Slug+"-answer", ext="html"`

### 3.3 SourceDocument 模型

```go
type SourceDocument struct {
    ID           string  // UUID v7
    Slug         string  // 从文件名派生
    Ext          string  // 原始文件扩展名
    OriginalURI  string
    OriginalHash string
}
```

**变更**：
- 移除 `FilePath` 字段
- 新增 `Slug` 和 `Ext` 字段

---

## 4. 实现步骤

### 4.1 阶段 1：实现 file 模块

#### 1.1 实现核心模块
- 在 `internal/file/` 目录实现新模块
- 实现 `FileStore` 接口及所有方法
- 实现分桶路径计算
- 实现 trash 机制
- 实现对象级锁
- 实现原子写入
- 编写单元测试

#### 1.2 配置更新
```go
// config/config.go
type Config struct {
    VaultDir string  // 文件根目录，默认 "vault"
}
```

### 4.2 阶段 2：更新 contracts 层

```go
// internal/contracts/filestore.go
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

### 4.3 阶段 3：更新数据库模型

```sql
-- Asset 表
ALTER TABLE assets ADD COLUMN slug VARCHAR(255);
ALTER TABLE assets ADD COLUMN ext VARCHAR(32);
ALTER TABLE assets DROP COLUMN storage_path;

-- Card 表
ALTER TABLE cards ADD COLUMN slug VARCHAR(255);
ALTER TABLE cards DROP COLUMN html_path;
ALTER TABLE cards DROP COLUMN answer_html_path;

-- SourceDocument 表
ALTER TABLE source_documents ADD COLUMN slug VARCHAR(255);
ALTER TABLE source_documents ADD COLUMN ext VARCHAR(32);
ALTER TABLE source_documents DROP COLUMN file_path;
```

### 4.4 阶段 4：更新业务代码

#### 4.1 更新 actions 层
```go
func (a *CreateCardAction) Execute(ctx context.Context, req CreateCardRequest) error {
    // 保存 HTML（路径自动计算）
    loc, err := a.fileStore.CreateFileObject(ctx, CreateFileObjectInput{
        Kind: "card",
        ID:   card.ID,
        Slug: card.Slug,
        Ext:  "html",
        Data: htmlContent,
    })

    // 只保存身份信息到数据库
    card.Slug = loc.Slug
    return a.cardRepo.Create(ctx, card)
}
```

#### 4.2 更新所有使用文件存储的 actions
- CreateCardAction
- UpdateCardAction
- DeleteCardAction
- CreateAssetAction
- ImportDocumentFlow
- 等等

---

## 5. 实现注意事项

### 5.1 性能考虑

**分桶目录结构**：
- 使用 UUID v7（时间有序）提高局部性
- 4 层分桶避免单目录文件过多
- 监控文件系统 inode 使用情况

**并发控制**：
- 实现对象级锁避免并发修改冲突
- 锁超时时间建议 30 秒
- 锁重试间隔建议 100 毫秒

### 5.2 错误处理

**文件缺失**：
- 使用 `Status = "missing"` 标记缺失文件
- 提供修复工具扫描并标记缺失对象
- 支持从 trash 恢复

**并发冲突**：
- 使用 `IfMatch` 条件防止 lost update
- 返回明确的 `ErrConcurrentModification` 错误

### 5.3 UUID v7 生成

确保使用标准 UUID v7 实现：
- 时间戳部分保证时间有序
- 随机部分保证唯一性
- 格式：`0195f3f0-a8b7-7c8d-b1b9-d45c8f23a1ee`

---

## 6. 测试计划

### 6.1 单元测试

- UUID v7 生成和验证
- 分桶路径计算
- 文件名构建和解析
- Trash 操作（移动、恢复、永久删除）
- 并发锁机制
- 原子写入

### 6.2 集成测试

- 完整的创建-读取-更新-删除流程
- Slug 重命名流程
- Trash 恢复流程
- 并发写入保护（IfMatch）
- 文件缺失处理

### 6.3 性能测试

- 大量文件场景（10万+ 文件）
- 并发读写性能
- 目录扫描性能
- 锁竞争场景

---

## 7. 实现时间表

| 阶段 | 任务 | 预计时间 |
|------|------|----------|
| 阶段 1 | 实现 file 模块核心功能 | 1 周 |
| 阶段 2 | 更新 contracts 层接口 | 1 天 |
| 阶段 3 | 更新数据库模型 | 1 天 |
| 阶段 4 | 更新业务代码 | 1 周 |
| 测试 | 单元测试 + 集成测试 | 3 天 |
| 验证 | 性能测试和问题修复 | 2 天 |

**总计**：约 2-3 周

---

## 8. 总结

本实现方案采用直接替换策略，特点：

1. **架构清晰**：派生路径策略，数据库只存身份信息
2. **双重保护**：数据库软删除 + 文件 trash 机制
3. **并发安全**：对象级锁 + 原子写入
4. **可扩展**：支持 Kind 扩展，支持未来新对象类型

关键实现要点：
- UUID v7 保证时间有序和唯一性
- 4 层分桶避免单目录文件过多
- 原子写入保证文件完整性
- Trash 机制支持文件恢复
