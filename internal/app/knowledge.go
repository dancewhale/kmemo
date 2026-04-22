package app

import (
	"strings"
	"time"

	"go.uber.org/zap"

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
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("knowledge.api")
	if strings.TrimSpace(req.Name) == "" {
		log.Info("CreateKnowledge rejected", zap.String("reason", "empty_name"))
		return "", repository.ErrInvalidInput
	}
	id, err := d.actions.Knowledge.Create(ctx, knowledge.CreateInput{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		ParentID:    req.ParentID,
	})
	if err != nil {
		log.Info("CreateKnowledge failed", zap.String("name", req.Name), zap.Error(err))
		return "", err
	}
	log.Info("CreateKnowledge ok",
		zap.String("id", id),
		zap.String("name", strings.TrimSpace(req.Name)),
		zap.String("description_excerpt", truncateRunes(strings.TrimSpace(req.Description), 120)),
		zapOptionalString("parentId", req.ParentID),
	)
	return id, nil
}

func (d *Desktop) GetKnowledge(id string) (*KnowledgeDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("knowledge.api")
	knowledgeModel, err := d.actions.Knowledge.Get(ctx, id)
	if err != nil {
		log.Info("GetKnowledge failed", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	cc, dc, err := d.actions.Knowledge.CountMapsByKnowledgeIDs(ctx, []string{id})
	if err != nil {
		log.Info("GetKnowledge failed", zap.String("id", id), zap.String("phase", "counts"), zap.Error(err))
		return nil, err
	}
	out := knowledgeDTOFromModel(knowledgeModel, cc, dc)
	log.Info("GetKnowledge ok",
		zap.String("id", out.ID),
		zap.String("name", out.Name),
		zap.Int("cardCount", out.CardCount),
		zap.Int("dueCount", out.DueCount),
		zapOptionalString("parentId", out.ParentID),
	)
	return out, nil
}

func (d *Desktop) ListKnowledge(parentID *string) ([]*KnowledgeDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("knowledge.api")
	items, err := d.actions.Knowledge.List(ctx, parentID)
	if err != nil {
		log.Info("ListKnowledge failed", zapOptionalString("parentId", parentID), zap.Error(err))
		return nil, err
	}
	ids := make([]string, len(items))
	for i, it := range items {
		ids[i] = it.ID
	}
	cc, dc, err := d.actions.Knowledge.CountMapsByKnowledgeIDs(ctx, ids)
	if err != nil {
		log.Info("ListKnowledge failed", zapOptionalString("parentId", parentID), zap.String("phase", "counts"), zap.Error(err))
		return nil, err
	}
	out := make([]*KnowledgeDTO, 0, len(items))
	for _, item := range items {
		out = append(out, knowledgeDTOFromModel(item, cc, dc))
	}
	sampleIDs := make([]string, 0, min(5, len(out)))
	for i := range out {
		if i >= 5 {
			break
		}
		sampleIDs = append(sampleIDs, out[i].ID)
	}
	log.Info("ListKnowledge ok",
		zapOptionalString("parentId", parentID),
		zap.Int("count", len(out)),
		zap.Strings("sample_ids", sampleIDs),
	)
	return out, nil
}

func (d *Desktop) GetKnowledgeTree(rootID *string) ([]*KnowledgeTreeNode, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("knowledge.api")
	if rootID != nil {
		root, err := d.actions.Knowledge.GetTree(ctx, *rootID)
		if err != nil {
			log.Info("GetKnowledgeTree failed", zap.String("rootId", *rootID), zap.Error(err))
			return nil, err
		}
		var ids []string
		appendKnowledgeSubtreeIDs(root, &ids)
		cc, dc, err := d.actions.Knowledge.CountMapsByKnowledgeIDs(ctx, ids)
		if err != nil {
			log.Info("GetKnowledgeTree failed", zap.String("rootId", *rootID), zap.String("phase", "counts"), zap.Error(err))
			return nil, err
		}
		forest := []*KnowledgeTreeNode{toKnowledgeTreeNode(root, cc, dc)}
		log.Info("GetKnowledgeTree ok",
			zap.String("mode", "subtree"),
			zap.String("rootId", *rootID),
			zap.Int("node_count", countKnowledgeTreeNodes(forest)),
		)
		return forest, nil
	}

	items, err := d.actions.Knowledge.ListAll(ctx)
	if err != nil {
		log.Info("GetKnowledgeTree failed", zap.String("mode", "forest"), zap.Error(err))
		return nil, err
	}
	ids := make([]string, len(items))
	for i, it := range items {
		ids[i] = it.ID
	}
	cc, dc, err := d.actions.Knowledge.CountMapsByKnowledgeIDs(ctx, ids)
	if err != nil {
		log.Info("GetKnowledgeTree failed", zap.String("mode", "forest"), zap.String("phase", "counts"), zap.Error(err))
		return nil, err
	}
	forest := buildKnowledgeForest(items, cc, dc)
	log.Info("GetKnowledgeTree ok",
		zap.String("mode", "forest"),
		zap.Int("root_count", len(forest)),
		zap.Int("node_count", countKnowledgeTreeNodes(forest)),
		zap.Strings("sample_root_ids", firstRootKnowledgeIDs(forest, 5)),
	)
	return forest, nil
}

func (d *Desktop) UpdateKnowledge(id string, req UpdateKnowledgeRequest) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("knowledge.api")
	if strings.TrimSpace(req.Name) == "" {
		log.Info("UpdateKnowledge rejected", zap.String("id", id), zap.String("reason", "empty_name"))
		return repository.ErrInvalidInput
	}
	err := d.actions.Knowledge.Update(ctx, id, knowledge.UpdateInput{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
	})
	if err != nil {
		log.Info("UpdateKnowledge failed", zap.String("id", id), zap.Error(err))
		return err
	}
	log.Info("UpdateKnowledge ok",
		zap.String("id", id),
		zap.String("name", strings.TrimSpace(req.Name)),
		zap.String("description_excerpt", truncateRunes(strings.TrimSpace(req.Description), 120)),
	)
	return nil
}

func (d *Desktop) DeleteKnowledge(id string) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("knowledge.api")
	err := d.actions.Knowledge.Delete(ctx, id)
	if err != nil {
		log.Info("DeleteKnowledge failed", zap.String("id", id), zap.Error(err))
		return err
	}
	log.Info("DeleteKnowledge ok", zap.String("id", id))
	return nil
}

func (d *Desktop) MoveKnowledge(id string, newParentID *string) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("knowledge.api")
	if newParentID != nil && *newParentID == id {
		log.Info("MoveKnowledge rejected", zap.String("id", id), zap.String("reason", "parent_is_self"))
		return repository.ErrInvalidInput
	}
	err := d.actions.Knowledge.Move(ctx, id, newParentID)
	if err != nil {
		log.Info("MoveKnowledge failed", zap.String("id", id), zapOptionalString("newParentId", newParentID), zap.Error(err))
		return err
	}
	log.Info("MoveKnowledge ok", zap.String("id", id), zapOptionalString("newParentId", newParentID))
	return nil
}

func (d *Desktop) ArchiveKnowledge(id string) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("knowledge.api")
	err := d.actions.Knowledge.Archive(ctx, id)
	if err != nil {
		log.Info("ArchiveKnowledge failed", zap.String("id", id), zap.Error(err))
		return err
	}
	log.Info("ArchiveKnowledge ok", zap.String("id", id))
	return nil
}

func (d *Desktop) UnarchiveKnowledge(id string) error {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("knowledge.api")
	err := d.actions.Knowledge.Unarchive(ctx, id)
	if err != nil {
		log.Info("UnarchiveKnowledge failed", zap.String("id", id), zap.Error(err))
		return err
	}
	log.Info("UnarchiveKnowledge ok", zap.String("id", id))
	return nil
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

