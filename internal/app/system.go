package app

import (
	"strings"

	"go.uber.org/zap"

	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
	"kmemo/internal/zaplog"
)

type FSRSParameterDTO struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	ParametersJSON   string   `json:"parametersJson"`
	DesiredRetention *float64 `json:"desiredRetention"`
	MaximumInterval  *int     `json:"maximumInterval"`
}

type UpdateDefaultFSRSParameterRequest struct {
	ParametersJSON   string   `json:"parametersJson"`
	DesiredRetention *float64 `json:"desiredRetention"`
	MaximumInterval  *int     `json:"maximumInterval"`
}

func (d *Desktop) ListFSRSParameters() ([]*FSRSParameterDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("system.api")
	items, err := d.actions.FSRSParameters.List(ctx)
	if err != nil {
		log.Info("ListFSRSParameters failed", zap.Error(err))
		return nil, err
	}
	out := toFSRSParameterDTOs(items)
	ids := make([]string, 0, min(5, len(out)))
	for i := range out {
		if i >= 5 {
			break
		}
		ids = append(ids, out[i].ID)
	}
	log.Info("ListFSRSParameters ok", zap.Int("count", len(out)), zap.Strings("sample_ids", ids))
	return out, nil
}

func (d *Desktop) GetDefaultFSRSParameter() (*FSRSParameterDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("system.api")
	item, err := d.actions.FSRSParameters.GetDefault(ctx)
	if err != nil {
		log.Info("GetDefaultFSRSParameter failed", zap.Error(err))
		return nil, err
	}
	out := toFSRSParameterDTO(item)
	log.Info("GetDefaultFSRSParameter ok",
		zap.String("id", out.ID),
		zap.String("name", out.Name),
		zap.Int("parameters_json_len", len(out.ParametersJSON)),
		zap.String("parameters_json_excerpt", truncateRunes(out.ParametersJSON, 160)),
	)
	return out, nil
}

func (d *Desktop) UpdateDefaultFSRSParameter(req UpdateDefaultFSRSParameterRequest) (*FSRSParameterDTO, error) {
	ctx := d.actionContext()
	log := zaplog.L(ctx).Named("system.api")
	if strings.TrimSpace(req.ParametersJSON) == "" {
		log.Info("UpdateDefaultFSRSParameter rejected", zap.String("reason", "empty_parameters_json"))
		return nil, repository.ErrInvalidInput
	}
	item, err := d.actions.FSRSParameters.UpdateDefault(ctx, &models.FSRSParameter{
		ParametersJSON:   req.ParametersJSON,
		DesiredRetention: req.DesiredRetention,
		MaximumInterval:  req.MaximumInterval,
	})
	if err != nil {
		log.Info("UpdateDefaultFSRSParameter failed", zap.Int("parameters_json_len", len(req.ParametersJSON)), zap.Error(err))
		return nil, err
	}
	out := toFSRSParameterDTO(item)
	log.Info("UpdateDefaultFSRSParameter ok",
		zap.String("id", out.ID),
		zap.Int("parameters_json_len", len(out.ParametersJSON)),
		zap.String("parameters_json_excerpt", truncateRunes(out.ParametersJSON, 160)),
	)
	return out, nil
}

func toFSRSParameterDTO(model *models.FSRSParameter) *FSRSParameterDTO {
	if model == nil {
		return nil
	}
	return &FSRSParameterDTO{
		ID:               model.ID,
		Name:             model.Name,
		ParametersJSON:   model.ParametersJSON,
		DesiredRetention: model.DesiredRetention,
		MaximumInterval:  model.MaximumInterval,
	}
}

func toFSRSParameterDTOs(items []*models.FSRSParameter) []*FSRSParameterDTO {
	result := make([]*FSRSParameterDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toFSRSParameterDTO(item))
	}
	return result
}
