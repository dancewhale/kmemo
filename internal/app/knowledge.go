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
	ctx := d.actionContext()
	knowledgeModel, err := d.actions.Knowledge.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	cc, dc, err := d.actions.Knowledge.CountMapsByKnowledgeIDs(ctx, []string{id})
	if err != nil {
		return nil, err
	}
	return knowledgeDTOFromModel(knowledgeModel, cc, dc), nil
}

func (d *Desktop) ListKnowledge(parentID *string) ([]*KnowledgeDTO, error) {
	ctx := d.actionContext()
	items, err := d.actions.Knowledge.List(ctx, parentID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(items))
	for i, it := range items {
		ids[i] = it.ID
	}
	cc, dc, err := d.actions.Knowledge.CountMapsByKnowledgeIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make([]*KnowledgeDTO, 0, len(items))
	for _, item := range items {
		out = append(out, knowledgeDTOFromModel(item, cc, dc))
	}
	return out, nil
}

func (d *Desktop) GetKnowledgeTree(rootID *string) ([]*KnowledgeTreeNode, error) {
	ctx := d.actionContext()
	if rootID != nil {
		root, err := d.actions.Knowledge.GetTree(ctx, *rootID)
		if err != nil {
			return nil, err
		}
		var ids []string
		appendKnowledgeSubtreeIDs(root, &ids)
		cc, dc, err := d.actions.Knowledge.CountMapsByKnowledgeIDs(ctx, ids)
		if err != nil {
			return nil, err
		}
		return []*KnowledgeTreeNode{toKnowledgeTreeNode(root, cc, dc)}, nil
	}

	items, err := d.actions.Knowledge.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(items))
	for i, it := range items {
		ids[i] = it.ID
	}
	cc, dc, err := d.actions.Knowledge.CountMapsByKnowledgeIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return buildKnowledgeForest(items, cc, dc), nil
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

func knowledgeDTOFromModel(m *models.Knowledge, cardCounts, dueCounts map[string]int64) *KnowledgeDTO {
	if m == nil {
		return nil
	}
	return &KnowledgeDTO{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		ParentID:    m.ParentID,
		CardCount:   int(cardCounts[m.ID]),
		DueCount:    int(dueCounts[m.ID]),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		ArchivedAt:  m.ArchivedAt,
	}
}

func appendKnowledgeSubtreeIDs(k *models.Knowledge, dst *[]string) {
	if k == nil {
		return
	}
	*dst = append(*dst, k.ID)
	for i := range k.Children {
		appendKnowledgeSubtreeIDs(&k.Children[i], dst)
	}
}

func toKnowledgeTreeNode(model *models.Knowledge, cardCounts, dueCounts map[string]int64) *KnowledgeTreeNode {
	if model == nil {
		return nil
	}
	node := &KnowledgeTreeNode{KnowledgeDTO: *knowledgeDTOFromModel(model, cardCounts, dueCounts)}
	for i := range model.Children {
		child := model.Children[i]
		node.Children = append(node.Children, toKnowledgeTreeNode(&child, cardCounts, dueCounts))
	}
	return node
}

func buildKnowledgeForest(items []*models.Knowledge, cardCounts, dueCounts map[string]int64) []*KnowledgeTreeNode {
	nodes := make(map[string]*KnowledgeTreeNode, len(items))
	roots := make([]*KnowledgeTreeNode, 0)
	for _, item := range items {
		nodes[item.ID] = &KnowledgeTreeNode{KnowledgeDTO: *knowledgeDTOFromModel(item, cardCounts, dueCounts)}
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

