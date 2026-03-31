# Wails UI 接口设计文档

## 1. 概述

本文档定义 kmemo 项目中 Golang 后端暴露给 Wails 前端的接口规范。

### 1.1 架构说明

```
┌─────────────────────────────────────┐
│   Frontend (Vue/React/Vanilla JS)   │
├─────────────────────────────────────┤
│   Wails Binding (window.go.main.*)  │
├─────────────────────────────────────┤
│   Desktop API (internal/app)        │  ← 本文档定义的接口层
├─────────────────────────────────────┤
│   Service Layer (业务逻辑)           │
├─────────────────────────────────────┤
│   Repository Layer (数据访问)        │
└─────────────────────────────────────┘
```

### 1.2 调用约定

- 前端通过 `window.go.main.App.MethodName()` 调用
- 所有方法返回 Promise
- 错误通过 Go error 返回，前端捕获为 rejected Promise

### 1.3 数据传输对象 (DTO)

DTO 用于前后端数据交换，避免直接暴露数据库模型。

## 2. 知识库管理接口

### 2.1 数据结构

```go
// KnowledgeDTO 知识库数据传输对象
type KnowledgeDTO struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    ParentID    *string   `json:"parentId"`
    CardCount   int       `json:"cardCount"`      // 卡片数量
    DueCount    int       `json:"dueCount"`       // 到期卡片数
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    ArchivedAt  *time.Time `json:"archivedAt"`
}

// KnowledgeTreeNode 知识库树节点
type KnowledgeTreeNode struct {
    KnowledgeDTO
    Children []*KnowledgeTreeNode `json:"children"`
}

// CreateKnowledgeRequest 创建知识库请求
type CreateKnowledgeRequest struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    ParentID    *string `json:"parentId"`
}

// UpdateKnowledgeRequest 更新知识库请求
type UpdateKnowledgeRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

### 2.2 接口方法

```go
// CreateKnowledge 创建知识库
// 返回: 新创建的知识库 ID
func (d *Desktop) CreateKnowledge(req CreateKnowledgeRequest) (string, error)

// GetKnowledge 获取知识库详情
func (d *Desktop) GetKnowledge(id string) (*KnowledgeDTO, error)

// ListKnowledge 列出知识库（平铺列表）
// parentID: 可选，过滤父节点
func (d *Desktop) ListKnowledge(parentID *string) ([]*KnowledgeDTO, error)

// GetKnowledgeTree 获取知识库树形结构
// rootID: 可选，指定根节点；为空则返回所有顶级节点
func (d *Desktop) GetKnowledgeTree(rootID *string) ([]*KnowledgeTreeNode, error)

// UpdateKnowledge 更新知识库
func (d *Desktop) UpdateKnowledge(id string, req UpdateKnowledgeRequest) error

// DeleteKnowledge 删除知识库（软删除）
func (d *Desktop) DeleteKnowledge(id string) error

// MoveKnowledge 移动知识库到新父节点
func (d *Desktop) MoveKnowledge(id string, newParentID *string) error

// ArchiveKnowledge 归档知识库
func (d *Desktop) ArchiveKnowledge(id string) error

// UnarchiveKnowledge 取消归档
func (d *Desktop) UnarchiveKnowledge(id string) error
```

## 3. 卡片管理接口

### 3.1 数据结构

```go
// CardDTO 卡片数据传输对象
type CardDTO struct {
    ID               string     `json:"id"`
    KnowledgeID      string     `json:"knowledgeId"`
    KnowledgeName    string     `json:"knowledgeName"`    // 冗余字段，便于显示
    SourceDocumentID *string    `json:"sourceDocumentId"`
    ParentID         *string    `json:"parentId"`
    Title            string     `json:"title"`
    CardType         string     `json:"cardType"`         // article/excerpt/qa/cloze/note
    HTMLPath         string     `json:"htmlPath"`
    Status           string     `json:"status"`           // active/suspended
    Tags             []*TagDTO  `json:"tags"`
    SRS              *SRSDTO    `json:"srs"`              // 可选，SRS 状态
    CreatedAt        time.Time  `json:"createdAt"`
    UpdatedAt        time.Time  `json:"updatedAt"`
}

// CardFilters 卡片查询过滤器
type CardFilters struct {
    KnowledgeID *string  `json:"knowledgeId"`
    CardType    string   `json:"cardType"`
    Status      string   `json:"status"`
    TagIDs      []string `json:"tagIds"`
    Keyword     string   `json:"keyword"`
    Limit       int      `json:"limit"`
    Offset      int      `json:"offset"`
}

// CreateCardRequest 创建卡片请求
type CreateCardRequest struct {
    KnowledgeID      string   `json:"knowledgeId"`
    SourceDocumentID *string  `json:"sourceDocumentId"`
    ParentID         *string  `json:"parentId"`
    Title            string   `json:"title"`
    CardType         string   `json:"cardType"`
    HTMLContent      string   `json:"htmlContent"`      // 前端传入 HTML 内容
    TagIDs           []string `json:"tagIds"`
}

// UpdateCardRequest 更新卡片请求
type UpdateCardRequest struct {
    Title       string `json:"title"`
    HTMLContent string `json:"htmlContent"`
    Status      string `json:"status"`
}
```

### 3.2 接口方法

```go
// CreateCard 创建卡片
func (d *Desktop) CreateCard(req CreateCardRequest) (string, error)

// GetCard 获取卡片详情
func (d *Desktop) GetCard(id string) (*CardDTO, error)

// ListCards 列出卡片
func (d *Desktop) ListCards(filters CardFilters) ([]*CardDTO, int64, error)

// UpdateCard 更新卡片
func (d *Desktop) UpdateCard(id string, req UpdateCardRequest) error

// DeleteCard 删除卡片（软删除）
func (d *Desktop) DeleteCard(id string) error

// AddCardTags 为卡片添加标签
func (d *Desktop) AddCardTags(cardID string, tagIDs []string) error

// RemoveCardTags 移除卡片标签
func (d *Desktop) RemoveCardTags(cardID string, tagIDs []string) error

// GetCardTags 获取卡片的所有标签
func (d *Desktop) GetCardTags(cardID string) ([]*TagDTO, error)

// SuspendCard 暂停卡片（不参与复习）
func (d *Desktop) SuspendCard(cardID string) error

// ResumeCard 恢复卡片
func (d *Desktop) ResumeCard(cardID string) error
```

## 4. 复习系统接口

### 4.1 数据结构

```go
// SRSDTO SRS 状态数据传输对象
type SRSDTO struct {
    CardID       string     `json:"cardId"`
    FSRSState    string     `json:"fsrsState"`        // new/learning/review/relearning
    DueAt        *time.Time `json:"dueAt"`
    LastReviewAt *time.Time `json:"lastReviewAt"`
    Stability    *float64   `json:"stability"`
    Difficulty   *float64   `json:"difficulty"`
    Reps         int        `json:"reps"`             // 复习次数
    Lapses       int        `json:"lapses"`           // 遗忘次数
}

// CardWithSRSDTO 带 SRS 状态的卡片
type CardWithSRSDTO struct {
    Card CardDTO `json:"card"`
    SRS  SRSDTO  `json:"srs"`
}

// SRSStatisticsDTO SRS 统计数据
type SRSStatisticsDTO struct {
    NewCount        int `json:"newCount"`
    LearningCount   int `json:"learningCount"`
    ReviewCount     int `json:"reviewCount"`
    RelearningCount int `json:"relearningCount"`
    TotalCards      int `json:"totalCards"`
    DueToday        int `json:"dueToday"`
}

// ReviewRequest 复习请求
type ReviewRequest struct {
    CardID string `json:"cardId"`
    Rating int    `json:"rating"`           // 1-4: Again/Hard/Good/Easy
}

// ReviewLogDTO 复习日志
type ReviewLogDTO struct {
    ID         string    `json:"id"`
    CardID     string    `json:"cardId"`
    ReviewedAt time.Time `json:"reviewedAt"`
    Rating     int       `json:"rating"`
    ReviewKind string    `json:"reviewKind"`
}
```

### 4.2 接口方法

```go
// GetDueCards 获取到期卡片
// knowledgeID: 可选，过滤知识库
// limit: 最多返回数量
func (d *Desktop) GetDueCards(knowledgeID *string, limit int) ([]*CardWithSRSDTO, error)

// SubmitReview 提交复习结果
func (d *Desktop) SubmitReview(req ReviewRequest) error

// UndoLastReview 撤销最近一次复习
func (d *Desktop) UndoLastReview(cardID string) error

// GetSRSStatistics 获取 SRS 统计
func (d *Desktop) GetSRSStatistics(knowledgeID *string) (*SRSStatisticsDTO, error)

// GetReviewHistory 获取复习历史
func (d *Desktop) GetReviewHistory(cardID string, limit int) ([]*ReviewLogDTO, error)

// GetReviewStats 获取复习统计（按日期）
func (d *Desktop) GetReviewStats(startDate, endDate time.Time) (*ReviewStatistics, error)
```

## 5. 标签管理接口

### 5.1 数据结构

```go
// TagDTO 标签数据传输对象
type TagDTO struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Slug        string    `json:"slug"`
    Color       string    `json:"color"`
    Icon        string    `json:"icon"`
    Description string    `json:"description"`
    CardCount   int       `json:"cardCount"`
    CreatedAt   time.Time `json:"createdAt"`
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
    Name        string `json:"name"`
    Slug        string `json:"slug"`
    Color       string `json:"color"`
    Icon        string `json:"icon"`
    Description string `json:"description"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
    Name        string `json:"name"`
    Color       string `json:"color"`
    Icon        string `json:"icon"`
    Description string `json:"description"`
}
```

### 5.2 接口方法

```go
// CreateTag 创建标签
func (d *Desktop) CreateTag(req CreateTagRequest) (string, error)

// GetTag 获取标签详情
func (d *Desktop) GetTag(id string) (*TagDTO, error)

// ListTags 列出所有标签
func (d *Desktop) ListTags() ([]*TagDTO, error)

// UpdateTag 更新标签
func (d *Desktop) UpdateTag(id string, req UpdateTagRequest) error

// DeleteTag 删除标签
func (d *Desktop) DeleteTag(id string) error

// SearchCardsByTags 按标签搜索卡片
func (d *Desktop) SearchCardsByTags(tagIDs []string) ([]*CardDTO, error)
```

## 6. 导入和内容处理接口

### 6.1 数据结构

```go
// ImportDocumentRequest 导入文档请求
type ImportDocumentRequest struct {
    KnowledgeID string `json:"knowledgeId"`
    FilePath    string `json:"filePath"`           // 本地文件路径
    SourceType  string `json:"sourceType"`         // pdf/epub/html/text
    Title       string `json:"title"`
    Author      string `json:"author"`
}

// ImportResult 导入结果
type ImportResult struct {
    SourceDocumentID string   `json:"sourceDocumentId"`
    CardIDs          []string `json:"cardIds"`          // 生成的卡片 ID 列表
    Message          string   `json:"message"`
}

// CleanHTMLRequest 清理 HTML 请求
type CleanHTMLRequest struct {
    HTML string `json:"html"`
}

// CleanHTMLResponse 清理 HTML 响应
type CleanHTMLResponse struct {
    CleanedHTML string `json:"cleanedHtml"`
}
```

### 6.2 接口方法

```go
// ImportDocument 导入文档
// 调用 Python gRPC 服务处理文档，生成卡片
func (d *Desktop) ImportDocument(req ImportDocumentRequest) (*ImportResult, error)

// CleanHTML 清理 HTML 内容
// 调用 Python gRPC 服务清理和标准化 HTML
func (d *Desktop) CleanHTML(req CleanHTMLRequest) (*CleanHTMLResponse, error)

// GetSourceDocuments 获取知识库的来源文档列表
func (d *Desktop) GetSourceDocuments(knowledgeID string) ([]*SourceDocumentDTO, error)
```

## 7. 系统和配置接口

### 7.1 数据结构

```go
// SystemInfo 系统信息
type SystemInfo struct {
    Version         string `json:"version"`
    PythonEndpoint  string `json:"pythonEndpoint"`
    DatabasePath    string `json:"databasePath"`
    TotalCards      int    `json:"totalCards"`
    TotalKnowledge  int    `json:"totalKnowledge"`
}

// FSRSParameterDTO FSRS 参数
type FSRSParameterDTO struct {
    ID               string   `json:"id"`
    Name             string   `json:"name"`
    ParametersJSON   string   `json:"parametersJson"`
    DesiredRetention *float64 `json:"desiredRetention"`
    MaximumInterval  *int     `json:"maximumInterval"`
}
```

### 7.2 接口方法

```go
// GetSystemInfo 获取系统信息
func (d *Desktop) GetSystemInfo() (*SystemInfo, error)

// GetVersion 获取版本号
func (d *Desktop) GetVersion() string

// PythonEndpoint 获取 Python 服务地址
func (d *Desktop) PythonEndpoint() string

// ListFSRSParameters 列出 FSRS 参数预设
func (d *Desktop) ListFSRSParameters() ([]*FSRSParameterDTO, error)

// GetDefaultFSRSParameter 获取默认 FSRS 参数
func (d *Desktop) GetDefaultFSRSParameter() (*FSRSParameterDTO, error)
```

## 8. 错误处理

### 8.1 错误码定义

```go
const (
    ErrCodeNotFound         = "NOT_FOUND"
    ErrCodeInvalidInput     = "INVALID_INPUT"
    ErrCodeDuplicateKey     = "DUPLICATE_KEY"
    ErrCodePythonService    = "PYTHON_SERVICE_ERROR"
    ErrCodeDatabaseError    = "DATABASE_ERROR"
    ErrCodeUnauthorized     = "UNAUTHORIZED"
)
```

### 8.2 错误响应

Go 的 error 会自动转换为前端的 rejected Promise，前端可以这样处理：

```javascript
try {
    const result = await window.go.main.App.CreateCard(request);
    console.log('Success:', result);
} catch (error) {
    console.error('Error:', error);
    // error 是字符串类型的错误消息
}
```

## 9. 前端调用示例

### 9.1 创建知识库

```javascript
const request = {
    name: "Go 编程",
    description: "Go 语言学习笔记",
    parentId: null
};

const knowledgeId = await window.go.main.App.CreateKnowledge(request);
console.log('Created knowledge:', knowledgeId);
```

### 9.2 获取到期卡片并复习

```javascript
// 获取到期卡片
const dueCards = await window.go.main.App.GetDueCards(null, 20);

// 复习第一张卡片
if (dueCards.length > 0) {
    const card = dueCards[0];
    await window.go.main.App.SubmitReview({
        cardId: card.card.id,
        rating: 3  // Good
    });
}
```

### 9.3 创建卡片并添加标签

```javascript
const cardRequest = {
    knowledgeId: "knowledge-123",
    title: "Go 并发模式",
    cardType: "article",
    htmlContent: "<h1>Go 并发模式</h1><p>内容...</p>",
    tagIds: ["tag-1", "tag-2"]
};

const cardId = await window.go.main.App.CreateCard(cardRequest);
```

## 10. 实施优先级

### Phase 1: 核心功能（MVP）
1. 知识库 CRUD
2. 卡片 CRUD
3. 基础 SRS（获取到期、提交复习）
4. 标签管理

### Phase 2: 增强功能
1. 导入文档
2. HTML 清理
3. 复习历史和统计
4. 撤销复习

### Phase 3: 高级功能
1. FSRS 参数管理
2. 高级搜索和过滤
3. 批量操作
4. 数据导出

## 11. 注意事项

1. **性能考虑**: 大量数据查询应使用分页
2. **并发安全**: Repository 层已处理事务，接口层无需额外处理
3. **错误处理**: 统一使用 Go error，前端统一捕获
4. **数据验证**: 在 Service 层进行业务逻辑验证
5. **HTML 存储**: HTML 内容存储为文件，接口只传递路径
