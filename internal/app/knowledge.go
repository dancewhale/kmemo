package app

import (
	"context"
	"strings"
	"time"

	"kmemo/internal/actions/knowledge"
	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
	"kmemo/internal/zaplog"
)

type KnowledgeDTO struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ParentID    *string    `json:"parentId"`
	CardCount   int        `json:"cardCount"`
	DueCount    int        `json:"dueCount"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
}

type KnowledgeTreeNode struct {
	KnowledgeDTO
	Children []*KnowledgeTreeNode `json:"children"`
}

type CreateKnowledgeRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ParentID    *string `json:"parentId"`
}

type UpdateKnowledgeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (d *Desktop) CreateKnowledge(req CreateKnowledgeRequest) (string, error) {
	if strings.TrimSpace(req.Name) == "" {
		return "", repository.ErrInvalidInput
	}
	return d.actions.Knowledge.Create(d.actionContext(), knowledge.CreateInput{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		ParentID:    req.ParentID,
	})
}

func (d *Desktop) GetKnowledge(id string) (*KnowledgeDTO, error) {
	knowledgeModel, err := d.actions.Knowledge.Get(d.actionContext(), id)
	if err != nil {
		return nil, err
	}
	return toKnowledgeDTO(knowledgeModel), nil
}

func (d *Desktop) ListKnowledge(parentID *string) ([]*KnowledgeDTO, error) {
	items, err := d.actions.Knowledge.List(d.actionContext(), parentID)
	if err != nil {
		return nil, err
	}
	return toKnowledgeDTOs(items), nil
}

func (d *Desktop) GetKnowledgeTree(rootID *string) ([]*KnowledgeTreeNode, error) {
	ctx := d.actionContext()
	if rootID != nil {
		root, err := d.actions.Knowledge.GetTree(ctx, *rootID)
		if err != nil {
			return nil, err
		}
		return []*KnowledgeTreeNode{toKnowledgeTreeNode(root)}, nil
	}

	items, err := d.actions.Knowledge.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildKnowledgeForest(items), nil
}

func (d *Desktop) UpdateKnowledge(id string, req UpdateKnowledgeRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return repository.ErrInvalidInput
	}
	return d.actions.Knowledge.Update(d.actionContext(), id, knowledge.UpdateInput{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
	})
}

func (d *Desktop) DeleteKnowledge(id string) error {
	return d.actions.Knowledge.Delete(d.actionContext(), id)
}

func (d *Desktop) MoveKnowledge(id string, newParentID *string) error {
	if newParentID != nil && *newParentID == id {
		return repository.ErrInvalidInput
	}
	return d.actions.Knowledge.Move(d.actionContext(), id, newParentID)
}

func (d *Desktop) ArchiveKnowledge(id string) error {
	return d.actions.Knowledge.Archive(d.actionContext(), id)
}

func (d *Desktop) UnarchiveKnowledge(id string) error {
	return d.actions.Knowledge.Unarchive(d.actionContext(), id)
}

func (d *Desktop) actionContext() context.Context {
	ctx := zaplog.WithLogger(context.Background(), d.logger)
	ctx, _ = zaplog.EnsureRequestID(ctx)
	return ctx
}

func toKnowledgeDTO(model *models.Knowledge) *KnowledgeDTO {
	if model == nil {
		return nil
	}
	return &KnowledgeDTO{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		ParentID:    model.ParentID,
		CardCount:   0,
		DueCount:    0,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		ArchivedAt:  model.ArchivedAt,
	}
}

func toKnowledgeDTOs(items []*models.Knowledge) []*KnowledgeDTO {
	result := make([]*KnowledgeDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toKnowledgeDTO(item))
	}
	return result
}

func toKnowledgeTreeNode(model *models.Knowledge) *KnowledgeTreeNode {
	if model == nil {
		return nil
	}
	node := &KnowledgeTreeNode{KnowledgeDTO: *toKnowledgeDTO(model)}
	for i := range model.Children {
		child := model.Children[i]
		node.Children = append(node.Children, toKnowledgeTreeNode(&child))
	}
	return node
}

func buildKnowledgeForest(items []*models.Knowledge) []*KnowledgeTreeNode {
	nodes := make(map[string]*KnowledgeTreeNode, len(items))
	roots := make([]*KnowledgeTreeNode, 0)
	for _, item := range items {
		nodes[item.ID] = &KnowledgeTreeNode{KnowledgeDTO: *toKnowledgeDTO(item)}
	}
	for _, item := range items {
		node := nodes[item.ID]
		if item.ParentID == nil {
			roots = append(roots, node)
			continue
		}
		parent := nodes[*item.ParentID]
		if parent == nil {
			roots = append(roots, node)
			continue
		}
		parent.Children = append(parent.Children, node)
	}
	return roots
}

