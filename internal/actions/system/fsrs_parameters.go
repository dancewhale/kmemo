package system

import (
	"context"

	fsrscontract "kmemo/internal/contracts/fsrs"
	"kmemo/internal/storage/models"
	"kmemo/internal/storage/repository"
)

type FSRSParametersDependencies struct {
	Parameters repository.FSRSParameterRepository
	FSRS       fsrscontract.FSRSScheduler
}

type FSRSParametersService struct {
	deps FSRSParametersDependencies
}

func NewFSRSParametersService(deps FSRSParametersDependencies) *FSRSParametersService {
	return &FSRSParametersService{deps: deps}
}

func (s *FSRSParametersService) List(ctx context.Context) ([]*models.FSRSParameter, error) {
	items, _, err := s.deps.Parameters.List(ctx, repository.ListFSRSParameterOptions{})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *FSRSParametersService) GetDefault(ctx context.Context) (*models.FSRSParameter, error) {
	return s.deps.Parameters.GetDefault(ctx)
}

func (s *FSRSParametersService) UpdateDefault(ctx context.Context, input *models.FSRSParameter) (*models.FSRSParameter, error) {
	current, err := s.deps.Parameters.GetDefault(ctx)
	if err != nil {
		return nil, err
	}
	current.ParametersJSON = input.ParametersJSON
	current.DesiredRetention = input.DesiredRetention
	current.MaximumInterval = input.MaximumInterval
	if err := s.deps.Parameters.Update(ctx, current); err != nil {
		return nil, err
	}
	if s.deps.FSRS != nil {
		if err := s.deps.FSRS.SetGlobalSetting(ctx, current); err != nil {
			return nil, err
		}
	}
	return current, nil
}
