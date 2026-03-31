# 数据库操作接口设计文档

## 1. 设计原则

### 1.1 分层架构

```
┌─────────────────────────────────────┐
│   Business Logic / Services Layer   │  业务逻辑层
├─────────────────────────────────────┤
│      Repository Layer (本设计)       │  仓储层 - 领域对象操作
├─────────────────────────────────────┤
│         DAO Layer (GORM Gen)        │  数据访问层 - 自动生成
├─────────────────────────────────────┤
│            Storage Layer            │  存储层 - 连接管理
└─────────────────────────────────────┘
```

### 1.2 核心原则

- **单一职责**: 每个 Repository 只负责一个聚合根
- **领域驱动**: 接口设计贴近业务语义，而非数据库操作
- **可测试性**: 接口可 mock，便于单元测试
- **事务支持**: 提供显式事务管理
- **错误处理**: 统一错误类型，区分业务错误和系统错误

## 2. 目录结构

```
internal/
├── storage/
│   ├── models/          # GORM 模型定义（已存在）
│   ├── dao/             # GORM Gen 生成代码（已存在）
│   ├── repository/      # Repository 接口与实现（新增）
│   │   ├── repository.go           # 基础接口定义
│   │   ├── knowledge_repo.go       # Knowledge 仓储
│   │   ├── card_repo.go            # Card 仓储
│   │   ├── tag_repo.go             # Tag 仓储
│   │   ├── srs_repo.go             # SRS 仓储
│   │   └── errors.go               # 仓储层错误定义
│   ├── storage.go       # Storage 主入口（已存在）
│   └── transaction.go   # 事务管理（已存在）
```

## 3. 基础接口定义

### 3.1 通用 Repository 接口

```go
// Package repository 提供领域对象的持久化操作接口
package repository

import (
    "context"
    "kmemo/internal/storage/models"
)

// Repository 所有仓储的基础接口
type Repository interface {
    // WithTx 返回使用指定事务的仓储实例
    WithTx(tx Transaction) Repository
}

// Transaction 事务接口
type Transaction interface {
    Commit() error
    Rollback() error
}
```

### 3.2 错误定义

```go
// errors.go
package repository

import "errors"

var (
    // ErrNotFound 记录不存在
    ErrNotFound = errors.New("record not found")

    // ErrDuplicateKey 唯一键冲突
    ErrDuplicateKey = errors.New("duplicate key")

    // ErrInvalidInput 输入参数无效
    ErrInvalidInput = errors.New("invalid input")

    // ErrConcurrentUpdate 并发更新冲突
    ErrConcurrentUpdate = errors.New("concurrent update conflict")
)
```

## 4. Knowledge Repository 设计

### 4.1 接口定义

```go
// KnowledgeRepository Knowledge 聚合根仓储
type KnowledgeRepository interface {
    Repository

    // Create 创建知识库
    Create(ctx context.Context, k *models.Knowledge) error

    // GetByID 根据 ID 获取
    GetByID(ctx context.Context, id string) (*models.Knowledge, error)

    // Update 更新知识库
    Update(ctx context.Context, k *models.Knowledge) error

    // Delete 软删除知识库
    Delete(ctx context.Context, id string) error

    // List 列表查询
    List(ctx context.Context, opts ListKnowledgeOptions) ([]*models.Knowledge, int64, error)

    // GetTree 获取树形结构（含子节点）
    GetTree(ctx context.Context, rootID string) (*models.Knowledge, error)

    // GetChildren 获取直接子节点
    GetChildren(ctx context.Context, parentID string) ([]*models.Knowledge, error)

    // Move 移动到新父节点
    Move(ctx context.Context, id string, newParentID *string) error

    // Archive 归档
    Archive(ctx context.Context, id string) error

    // Unarchive 取消归档
    Unarchive(ctx context.Context, id string) error
}

// ListKnowledgeOptions 查询选项
type ListKnowledgeOptions struct {
    ParentID      *string  // 父节点过滤
    IncludeArchived bool   // 是否包含已归档
    Keyword       string   // 名称/描述关键词搜索
    Limit         int      // 分页大小
    Offset        int      // 分页偏移
    OrderBy       string   // 排序字段: name, created_at, updated_at
    OrderDesc     bool     // 是否降序
}
```

## 5. Card Repository 设计

### 5.1 接口定义

```go
// CardRepository Card 聚合根仓储
type CardRepository interface {
    Repository

    // Create 创建卡片
    Create(ctx context.Context, c *models.Card) error

    // GetByID 根据 ID 获取（可选预加载关联）
    GetByID(ctx context.Context, id string, preload ...string) (*models.Card, error)

    // Update 更新卡片
    Update(ctx context.Context, c *models.Card) error

    // Delete 软删除卡片
    Delete(ctx context.Context, id string) error

    // List 列表查询
    List(ctx context.Context, opts ListCardOptions) ([]*models.Card, int64, error)

    // GetByKnowledge 获取指定知识库下的卡片
    GetByKnowledge(ctx context.Context, knowledgeID string, opts ListCardOptions) ([]*models.Card, int64, error)

    // GetChildren 获取子卡片
    GetChildren(ctx context.Context, parentID string) ([]*models.Card, error)

    // UpdateStatus 批量更新状态
    UpdateStatus(ctx context.Context, ids []string, status string) error

    // AddTags 为卡片添加标签
    AddTags(ctx context.Context, cardID string, tagIDs []string) error

    // RemoveTags 移除卡片标签
    RemoveTags(ctx context.Context, cardID string, tagIDs []string) error

    // GetTags 获取卡片的所有标签
    GetTags(ctx context.Context, cardID string) ([]*models.Tag, error)
}

// ListCardOptions 卡片查询选项
type ListCardOptions struct {
    KnowledgeID      *string   // 知识库过滤
    SourceDocumentID *string   // 来源文档过滤
    ParentID         *string   // 父卡片过滤
    CardType         string    // 卡片类型过滤
    Status           string    // 状态过滤
    TagIDs           []string  // 标签过滤（AND 关系）
    Keyword          string    // 标题关键词搜索
    IsRoot           *bool     // 是否根卡片
    Limit            int
    Offset           int
    OrderBy          string    // title, created_at, updated_at, sort_order
    OrderDesc        bool
    Preload          []string  // 预加载关联: Tags, SRS, Assets, Knowledge
}
```

## 6. Tag Repository 设计

### 6.1 接口定义

```go
// TagRepository Tag 仓储
type TagRepository interface {
    Repository

    // Create 创建标签
    Create(ctx context.Context, t *models.Tag) error

    // GetByID 根据 ID 获取
    GetByID(ctx context.Context, id string) (*models.Tag, error)

    // GetBySlug 根据 slug 获取
    GetBySlug(ctx context.Context, slug string) (*models.Tag, error)

    // Update 更新标签
    Update(ctx context.Context, t *models.Tag) error

    // Delete 软删除标签
    Delete(ctx context.Context, id string) error

    // List 列表查询
    List(ctx context.Context, opts ListTagOptions) ([]*models.Tag, int64, error)

    // GetByCardID 获取卡片的所有标签
    GetByCardID(ctx context.Context, cardID string) ([]*models.Tag, error)

    // UpdateCardCount 更新标签的卡片计数缓存
    UpdateCardCount(ctx context.Context, tagID string) error
}

// ListTagOptions 标签查询选项
type ListTagOptions struct {
    Keyword   string
    Limit     int
    Offset    int
    OrderBy   string // name, sort_order, card_count, created_at
    OrderDesc bool
}
```

## 7. SRS Repository 设计

### 7.1 接口定义

```go
// SRSRepository SRS 调度仓储
type SRSRepository interface {
    Repository

    // CreateOrUpdate 创建或更新 SRS 状态
    CreateOrUpdate(ctx context.Context, srs *models.CardSRS) error

    // GetByCardID 获取卡片的 SRS 状态
    GetByCardID(ctx context.Context, cardID string) (*models.CardSRS, error)

    // GetDueCards 获取到期卡片
    GetDueCards(ctx context.Context, opts DueCardsOptions) ([]*models.CardSRS, error)

    // UpdateAfterReview 复习后更新状态
    UpdateAfterReview(ctx context.Context, cardID string, srs *models.CardSRS, log *models.ReviewLog) error

    // GetStatistics 获取学习统计
    GetStatistics(ctx context.Context, knowledgeID *string) (*SRSStatistics, error)

    // Suspend 暂停卡片
    Suspend(ctx context.Context, cardID string) error

    // Resume 恢复卡片
    Resume(ctx context.Context, cardID string) error

    // UndoLastReview 撤销最近一次复习记录
    // 将 SRS 状态恢复到上一次复习前的状态，并删除最后一条复习日志
    UndoLastReview(ctx context.Context, cardID string) error

    // GetLastReviewLog 获取卡片的最后一条复习日志
    GetLastReviewLog(ctx context.Context, cardID string) (*models.ReviewLog, error)
}

// DueCardsOptions 到期卡片查询选项
type DueCardsOptions struct {
    KnowledgeID *string    // 知识库过滤
    Limit       int        // 最多返回数量
    States      []string   // FSRS 状态过滤: new, learning, review, relearning
}

// SRSStatistics SRS 统计数据
type SRSStatistics struct {
    NewCount       int
    LearningCount  int
    ReviewCount    int
    RelearningCount int
    TotalCards     int
    DueToday       int
}
```

## 8. ReviewLog Repository 设计

### 8.1 接口定义

```go
// ReviewLogRepository 复习日志仓储
type ReviewLogRepository interface {
    Repository

    // Create 创建复习日志
    Create(ctx context.Context, log *models.ReviewLog) error

    // GetByID 根据 ID 获取
    GetByID(ctx context.Context, id string) (*models.ReviewLog, error)

    // GetByCardID 获取卡片的所有复习日志
    GetByCardID(ctx context.Context, cardID string, opts ListReviewLogOptions) ([]*models.ReviewLog, int64, error)

    // GetLastByCardID 获取卡片的最后一条复习日志
    GetLastByCardID(ctx context.Context, cardID string) (*models.ReviewLog, error)

    // Delete 删除复习日志
    Delete(ctx context.Context, id string) error

    // GetStatsByDateRange 获取指定时间范围的复习统计
    GetStatsByDateRange(ctx context.Context, startDate, endDate time.Time, knowledgeID *string) (*ReviewStatistics, error)
}

// ListReviewLogOptions 复习日志查询选项
type ListReviewLogOptions struct {
    StartDate *time.Time // 开始时间过滤
    EndDate   *time.Time // 结束时间过滤
    Rating    *int       // 评分过滤
    Limit     int
    Offset    int
    OrderDesc bool // 默认按时间降序
}

// ReviewStatistics 复习统计
type ReviewStatistics struct {
    TotalReviews int
    AvgRating    float64
    RatingCounts map[int]int // 各评分的数量
    ReviewsByDay map[string]int // 按日期统计
}
```

## 9. FSRSParameter Repository 设计

### 9.1 接口定义

```go
// FSRSParameterRepository FSRS 参数仓储
type FSRSParameterRepository interface {
    Repository

    // Create 创建参数预设
    Create(ctx context.Context, param *models.FSRSParameter) error

    // GetByID 根据 ID 获取
    GetByID(ctx context.Context, id string) (*models.FSRSParameter, error)

    // GetByName 根据名称获取
    GetByName(ctx context.Context, name string) (*models.FSRSParameter, error)

    // Update 更新参数预设
    Update(ctx context.Context, param *models.FSRSParameter) error

    // Delete 删除参数预设
    Delete(ctx context.Context, id string) error

    // List 列表查询
    List(ctx context.Context, opts ListFSRSParameterOptions) ([]*models.FSRSParameter, int64, error)

    // GetDefault 获取默认参数预设
    GetDefault(ctx context.Context) (*models.FSRSParameter, error)
}

// ListFSRSParameterOptions FSRS 参数查询选项
type ListFSRSParameterOptions struct {
    Keyword   string
    Limit     int
    Offset    int
    OrderBy   string // name, created_at, updated_at
    OrderDesc bool
}
```

## 10. Asset Repository 设计

### 8.1 接口定义

```go
// AssetRepository Asset 仓储
type AssetRepository interface {
    Repository

    // Create 创建资源
    Create(ctx context.Context, a *models.Asset) error

    // GetByID 根据 ID 获取
    GetByID(ctx context.Context, id string) (*models.Asset, error)

    // GetByCardID 获取卡片的所有资源
    GetByCardID(ctx context.Context, cardID string, kind ...string) ([]*models.Asset, error)

    // GetByChecksum 根据校验和查找（用于去重）
    GetByChecksum(ctx context.Context, checksum string) (*models.Asset, error)

    // Update 更新资源
    Update(ctx context.Context, a *models.Asset) error

    // Delete 软删除资源
    Delete(ctx context.Context, id string) error

    // MarkMissing 标记文件丢失
    MarkMissing(ctx context.Context, id string) error
}
```

## 9. SourceDocument Repository 设计

### 9.1 接口定义

```go
// SourceDocumentRepository 来源文档仓储
type SourceDocumentRepository interface {
    Repository

    // Create 创建来源文档
    Create(ctx context.Context, sd *models.SourceDocument) error

    // GetByID 根据 ID 获取
    GetByID(ctx context.Context, id string) (*models.SourceDocument, error)

    // GetByKnowledgeID 获取知识库的所有来源文档
    GetByKnowledgeID(ctx context.Context, knowledgeID string) ([]*models.SourceDocument, error)

    // Update 更新来源文档
    Update(ctx context.Context, sd *models.SourceDocument) error

    // Delete 删除来源文档
    Delete(ctx context.Context, id string) error
}
```

## 10. 事务管理

### 10.1 事务接口

```go
// TransactionManager 事务管理器
type TransactionManager interface {
    // BeginTx 开启事务
    BeginTx(ctx context.Context) (Transaction, error)

    // WithTx 在事务中执行函数
    WithTx(ctx context.Context, fn func(tx Transaction) error) error
}

// RepositoryFactory 仓储工厂（用于事务中创建仓储）
type RepositoryFactory interface {
    Knowledge() KnowledgeRepository
    Card() CardRepository
    Tag() TagRepository
    SRS() SRSRepository
    ReviewLog() ReviewLogRepository
    FSRSParameter() FSRSParameterRepository
    Asset() AssetRepository
    SourceDocument() SourceDocumentRepository
}
```

## 11. 使用示例

### 11.1 基础 CRUD 操作

```go
// 创建知识库
func CreateKnowledge(ctx context.Context, repo KnowledgeRepository) error {
    k := &models.Knowledge{
        ID:          generateID(),
        Name:        "Go 编程",
        Description: "Go 语言学习笔记",
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    return repo.Create(ctx, k)
}

// 查询卡片列表
func ListCards(ctx context.Context, repo CardRepository) ([]*models.Card, error) {
    opts := ListCardOptions{
        Status:    "active",
        Limit:     20,
        Offset:    0,
        OrderBy:   "created_at",
        OrderDesc: true,
        Preload:   []string{"Tags", "SRS"},
    }
    cards, _, err := repo.List(ctx, opts)
    return cards, err
}
```

### 11.2 事务操作示例

```go
// 创建卡片并添加标签（事务）
func CreateCardWithTags(ctx context.Context, tm TransactionManager, factory RepositoryFactory) error {
    return tm.WithTx(ctx, func(tx Transaction) error {
        cardRepo := factory.Card().WithTx(tx).(CardRepository)

        // 创建卡片
        card := &models.Card{
            ID:          generateID(),
            KnowledgeID: "knowledge-123",
            Title:       "Go 并发模式",
            CardType:    "article",
            Status:      "active",
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        }
        if err := cardRepo.Create(ctx, card); err != nil {
            return err
        }

        // 添加标签
        tagIDs := []string{"tag-1", "tag-2"}
        return cardRepo.AddTags(ctx, card.ID, tagIDs)
    })
}

// 复习卡片并记录日志（事务）
func ReviewCard(ctx context.Context, tm TransactionManager, factory RepositoryFactory, cardID string, rating int) error {
    return tm.WithTx(ctx, func(tx Transaction) error {
        srsRepo := factory.SRS().WithTx(tx).(SRSRepository)

        // 获取当前 SRS 状态
        currentSRS, err := srsRepo.GetByCardID(ctx, cardID)
        if err != nil {
            return err
        }

        // 计算新状态（调用 Python gRPC 服务）
        newSRS := calculateNewSRS(currentSRS, rating)

        // 创建复习日志
        log := &models.ReviewLog{
            ID:                 generateID(),
            CardID:             cardID,
            ReviewedAt:         time.Now(),
            Rating:             rating,
            ReviewKind:         currentSRS.FSRSState,
            PreviousStability:  currentSRS.Stability,
            PreviousDifficulty: currentSRS.Difficulty,
            NewStability:       newSRS.Stability,
            NewDifficulty:      newSRS.Difficulty,
        }

        // 更新 SRS 状态并记录日志
        return srsRepo.UpdateAfterReview(ctx, cardID, newSRS, log)
    })
}
```

### 11.3 复杂查询示例

```go
// 获取到期卡片
func GetDueCards(ctx context.Context, srsRepo SRSRepository, knowledgeID string) ([]*models.CardSRS, error) {
    opts := DueCardsOptions{
        KnowledgeID: &knowledgeID,
        Limit:       50,
        States:      []string{"learning", "review"},
    }
    return srsRepo.GetDueCards(ctx, opts)
}

// 按标签搜索卡片
func SearchCardsByTags(ctx context.Context, cardRepo CardRepository, tagIDs []string) ([]*models.Card, error) {
    opts := ListCardOptions{
        TagIDs:    tagIDs,
        Status:    "active",
        Limit:     100,
        OrderBy:   "updated_at",
        OrderDesc: true,
        Preload:   []string{"Tags"},
    }
    cards, _, err := cardRepo.List(ctx, opts)
    return cards, err
}

// 撤销最近一次复习
func UndoReview(ctx context.Context, srsRepo SRSRepository, cardID string) error {
    // 获取最后一条复习日志以确认是否可撤销
    lastLog, err := srsRepo.GetLastReviewLog(ctx, cardID)
    if err != nil {
        return err
    }
    if lastLog == nil {
        return errors.New("no review to undo")
    }

    // 撤销复习（事务中恢复 SRS 状态并删除日志）
    return srsRepo.UndoLastReview(ctx, cardID)
}
```

## 12. 实现建议

### 12.1 Repository 实现结构

```go
// knowledge_repo.go 实现示例
package repository

import (
    "context"
    "kmemo/internal/storage/dao"
    "kmemo/internal/storage/models"
    "gorm.io/gorm"
)

type knowledgeRepo struct {
    db *gorm.DB
}

func NewKnowledgeRepository(db *gorm.DB) KnowledgeRepository {
    return &knowledgeRepo{db: db}
}

func (r *knowledgeRepo) WithTx(tx Transaction) Repository {
    return &knowledgeRepo{db: tx.(*gorm.DB)}
}

func (r *knowledgeRepo) Create(ctx context.Context, k *models.Knowledge) error {
    return dao.Use(r.db).Knowledge.WithContext(ctx).Create(k)
}

func (r *knowledgeRepo) GetByID(ctx context.Context, id string) (*models.Knowledge, error) {
    k, err := dao.Use(r.db).Knowledge.WithContext(ctx).Where(
        dao.Knowledge.ID.Eq(id),
    ).First()
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, ErrNotFound
        }
        return nil, err
    }
    return k, nil
}
```

### 12.2 错误处理

```go
// 统一错误转换
func convertError(err error) error {
    if err == nil {
        return nil
    }
    if err == gorm.ErrRecordNotFound {
        return ErrNotFound
    }
    if err == gorm.ErrDuplicatedKey {
        return ErrDuplicateKey
    }
    return err
}
```

### 12.3 事务管理实现

```go
// transaction.go
package repository

import (
    "context"
    "gorm.io/gorm"
)

type transactionManager struct {
    db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
    return &transactionManager{db: db}
}

func (tm *transactionManager) BeginTx(ctx context.Context) (Transaction, error) {
    return tm.db.WithContext(ctx).Begin(), nil
}

func (tm *transactionManager) WithTx(ctx context.Context, fn func(tx Transaction) error) error {
    return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        return fn(tx)
    })
}
```

## 13. 最佳实践

### 13.1 命名规范

- **Repository 接口**: 使用领域名词 + `Repository` 后缀，如 `KnowledgeRepository`
- **方法命名**: 使用业务语义动词
  - `Create` / `Update` / `Delete`: 基础 CRUD
  - `GetByID` / `GetBySlug`: 单条查询
  - `List` / `GetChildren`: 列表查询
  - `Archive` / `Suspend`: 业务操作
- **选项结构**: 使用 `List{Entity}Options` 命名模式

### 13.2 性能优化

1. **预加载控制**: 通过 `Preload` 参数按需加载关联
2. **分页查询**: 始终使用 `Limit` 和 `Offset`
3. **索引利用**: 查询条件对应数据库索引字段
4. **批量操作**: 提供批量更新接口减少数据库往返

### 13.3 可测试性

```go
// 使用接口便于 mock
type mockKnowledgeRepo struct {
    mock.Mock
}

func (m *mockKnowledgeRepo) Create(ctx context.Context, k *models.Knowledge) error {
    args := m.Called(ctx, k)
    return args.Error(0)
}

// 测试示例
func TestCreateKnowledge(t *testing.T) {
    repo := new(mockKnowledgeRepo)
    repo.On("Create", mock.Anything, mock.Anything).Return(nil)

    err := CreateKnowledge(context.Background(), repo)
    assert.NoError(t, err)
    repo.AssertExpectations(t)
}
```

### 13.4 错误处理原则

1. **区分错误类型**: 业务错误（ErrNotFound）vs 系统错误
2. **错误传播**: Repository 层不吞噬错误，向上传播
3. **错误转换**: 将 GORM 错误转换为领域错误
4. **上下文信息**: 必要时包装错误添加上下文

## 14. 设计优势

### 14.1 分层清晰

- **DAO 层**: GORM Gen 自动生成，提供类型安全的查询构建器
- **Repository 层**: 封装业务语义，隐藏 ORM 细节
- **Service 层**: 编排多个 Repository，实现复杂业务逻辑

### 14.2 易于理解

- 接口方法名直接表达业务意图
- 参数结构化，避免过长参数列表
- 返回值统一，错误处理一致

### 14.3 可维护性强

- 接口与实现分离，便于替换实现
- 单一职责，每个 Repository 只管理一个聚合根
- 事务管理显式，避免隐式事务

### 14.4 可扩展性好

- 新增查询条件只需扩展 Options 结构
- 新增 Repository 不影响现有代码
- 支持装饰器模式添加缓存、日志等横切关注点

## 15. 实施计划

### 15.1 第一阶段：基础框架

1. 创建 `internal/storage/repository/` 目录
2. 实现基础接口定义（`repository.go`, `errors.go`）
3. 实现事务管理（`transaction.go`）
4. 创建 RepositoryFactory

### 15.2 第二阶段：核心 Repository

1. 实现 `KnowledgeRepository`
2. 实现 `CardRepository`
3. 实现 `TagRepository`
4. 编写单元测试

### 15.3 第三阶段：SRS 和资源管理

1. 实现 `SRSRepository`
2. 实现 `AssetRepository`
3. 实现 `SourceDocumentRepository`
4. 集成测试

### 15.4 第四阶段：优化和完善

1. 性能优化（查询优化、索引验证）
2. 添加缓存层（可选）
3. 完善文档和示例
4. 代码审查和重构

## 16. 总结

本设计文档提供了一套完整的数据库操作接口方案，具有以下特点：

- **分层清晰**: DAO → Repository → Service 三层架构
- **语义明确**: 接口方法名贴近业务领域
- **易于测试**: 接口可 mock，便于单元测试
- **事务支持**: 显式事务管理，保证数据一致性
- **可扩展**: 易于添加新功能和优化

该设计遵循开源项目最佳实践，参考了 DDD（领域驱动设计）和 Repository 模式，适合中大型项目使用。
