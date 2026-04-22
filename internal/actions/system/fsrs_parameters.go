package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

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

const defaultFSRSParameterName = "Default"

func (s *FSRSParametersService) GetDefault(ctx context.Context) (*models.FSRSParameter, error) {
	p, err := s.deps.Parameters.GetDefault(ctx)
	if err == nil {
		return p, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	builtin, err := s.deps.FSRS.GetBuiltinDefaultParameters(ctx)
	if err != nil {
		return nil, fmt.Errorf("system: fetch builtin fsrs defaults: %w", err)
	}

	payload, err := json.Marshal(builtin.Parameters)
	if err != nil {
		return nil, fmt.Errorf("system: marshal builtin fsrs parameters: %w", err)
	}

	now := time.Now().UTC()
	dr := builtin.DesiredRetention
	max := builtin.MaximumInterval
	param := &models.FSRSParameter{
		ID:               uuid.NewString(),
		Name:             defaultFSRSParameterName,
		ParametersJSON:   string(payload),
		DesiredRetention: &dr,
		MaximumInterval:  &max,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.deps.Parameters.Create(ctx, param); err != nil {
		again, ge := s.deps.Parameters.GetDefault(ctx)
		if ge == nil && again != nil {
			return again, nil
		}
		return nil, fmt.Errorf("system: persist default fsrs parameter: %w", err)
	}

	if err := s.deps.FSRS.SetGlobalSetting(ctx, param); err != nil {
		return nil, err
	}
	return param, nil
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
