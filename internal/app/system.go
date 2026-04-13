package app

import (
	"strings"

	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
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
	items, err := d.actions.FSRSParameters.List(ctx)
	if err != nil {
		return nil, err
	}
	return toFSRSParameterDTOs(items), nil
}

func (d *Desktop) GetDefaultFSRSParameter() (*FSRSParameterDTO, error) {
	ctx := d.actionContext()
	item, err := d.actions.FSRSParameters.GetDefault(ctx)
	if err != nil {
		return nil, err
	}
	return toFSRSParameterDTO(item), nil
}

func (d *Desktop) UpdateDefaultFSRSParameter(req UpdateDefaultFSRSParameterRequest) (*FSRSParameterDTO, error) {
	if strings.TrimSpace(req.ParametersJSON) == "" {
		return nil, repository.ErrInvalidInput
	}
	ctx := d.actionContext()
	item, err := d.actions.FSRSParameters.UpdateDefault(ctx, &models.FSRSParameter{
		ParametersJSON:   req.ParametersJSON,
		DesiredRetention: req.DesiredRetention,
		MaximumInterval:  req.MaximumInterval,
	})
	if err != nil {
		return nil, err
	}
	return toFSRSParameterDTO(item), nil
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
