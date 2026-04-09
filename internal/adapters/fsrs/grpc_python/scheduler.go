package grpcpython

import (
	"context"
	"fmt"

	kmemov1 "kmemo/gen/kmemo/v1"
	"kmemo/internal/contracts/fsrs"
	"kmemo/internal/storage/models"
)

// Scheduler implements fsrs.FSRSScheduler via KmemoProcessor gRPC (Python worker).
type Scheduler struct {
	api kmemov1.KmemoProcessorClient
}

// NewScheduler constructs a Scheduler. api must be non-nil for live calls.
func NewScheduler(api kmemov1.KmemoProcessorClient) *Scheduler {
	return &Scheduler{api: api}
}

var _ fsrs.FSRSScheduler = (*Scheduler)(nil)

// Review implements fsrs.FSRSScheduler.
func (s *Scheduler) Review(ctx context.Context, in fsrs.ReviewInput) (*fsrs.ReviewOutput, error) {
	if s == nil || s.api == nil {
		return nil, fsrs.ErrUnavailable
	}
	if in.CardID == "" {
		return nil, fmt.Errorf("%w: card_id is required", fsrs.ErrInvalidInput)
	}
	if in.ReviewedAt.IsZero() {
		return nil, fmt.Errorf("%w: reviewed_at is required", fsrs.ErrInvalidInput)
	}

	resp, err := s.api.ReviewCard(ctx, ReviewCardRequest(in.CardID, in.Prior, in.Rating, in.ReviewedAt))
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.GetCard() == nil {
		return nil, fmt.Errorf("fsrs: empty ReviewCard response")
	}

	cardID := in.CardID
	next := CardStateToModel(cardID, resp.GetCard())
	var log *models.ReviewLog
	if snap := resp.GetReviewLog(); snap != nil {
		log = ReviewLogFromProto(cardID, snap)
	} else {
		log = &models.ReviewLog{
			CardID:       cardID,
			Rating:       in.Rating,
			ReviewedAt:   in.ReviewedAt.UTC(),
			ReviewKind:   "review",
			SnapshotJSON: "{}",
		}
	}

	return &fsrs.ReviewOutput{
		Next:           next,
		Log:            log,
		Retrievability: resp.GetRetrievability(),
		Warnings:       append([]string(nil), resp.GetWarnings()...),
		Diagnostics:    resp.GetDiagnostics(),
	}, nil
}

// GetRetrievability implements fsrs.FSRSScheduler.
func (s *Scheduler) GetRetrievability(ctx context.Context, in fsrs.RetrievabilityInput) (*fsrs.RetrievabilityOutput, error) {
	if s == nil || s.api == nil {
		return nil, fsrs.ErrUnavailable
	}
	if in.CardID == "" {
		return nil, fmt.Errorf("%w: card_id is required", fsrs.ErrInvalidInput)
	}
	if in.Now.IsZero() {
		return nil, fmt.Errorf("%w: now is required", fsrs.ErrInvalidInput)
	}

	resp, err := s.api.GetCardRetrievability(ctx, RetrievabilityRequest(in.CardID, in.Prior, in.Now))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("fsrs: empty GetCardRetrievability response")
	}

	out := &fsrs.RetrievabilityOutput{
		Retrievability: resp.GetRetrievability(),
		Warnings:       append([]string(nil), resp.GetWarnings()...),
		Diagnostics:    resp.GetDiagnostics(),
	}
	if ts := resp.GetEvaluatedAt(); ts != nil {
		t := ts.AsTime().UTC()
		out.EvaluatedAt = &t
	}
	return out, nil
}

// Reschedule implements fsrs.FSRSScheduler.
func (s *Scheduler) Reschedule(ctx context.Context, in fsrs.RescheduleInput) (*fsrs.RescheduleOutput, error) {
	if s == nil || s.api == nil {
		return nil, fsrs.ErrUnavailable
	}
	if in.CardID == "" {
		return nil, fmt.Errorf("%w: card_id is required", fsrs.ErrInvalidInput)
	}
	if in.RescheduleAt.IsZero() {
		return nil, fmt.Errorf("%w: reschedule_at is required", fsrs.ErrInvalidInput)
	}

	resp, err := s.api.RescheduleCard(ctx, RescheduleCardRequest(in.CardID, in.Prior, in.ReviewLogs, in.RescheduleAt, in.Mode))
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.GetCard() == nil {
		return nil, fmt.Errorf("fsrs: empty RescheduleCard response")
	}

	cardID := in.CardID
	next := CardStateToModel(cardID, resp.GetCard())
	var log *models.ReviewLog
	if snap := resp.GetReviewLog(); snap != nil {
		log = ReviewLogFromProto(cardID, snap)
	}

	return &fsrs.RescheduleOutput{
		Next:           next,
		Log:            log,
		Retrievability: resp.GetRetrievability(),
		Warnings:       append([]string(nil), resp.GetWarnings()...),
		Diagnostics:    resp.GetDiagnostics(),
	}, nil
}

// SetGlobalSetting implements fsrs.FSRSScheduler.
func (s *Scheduler) SetGlobalSetting(ctx context.Context, param *models.FSRSParameter) error {
	if s == nil || s.api == nil {
		return fsrs.ErrUnavailable
	}
	if param == nil {
		return fmt.Errorf("%w: param is nil", fsrs.ErrInvalidInput)
	}
	setting, err := SettingFromFSRSParameter(param)
	if err != nil {
		return err
	}
	_, err = s.api.SchedulerSetSetting(ctx, &kmemov1.SchedulerSetSettingRequest{Setting: setting})
	return err
}

// OptimizeParameters implements fsrs.FSRSScheduler.
func (s *Scheduler) OptimizeParameters(ctx context.Context, in fsrs.OptimizeParametersInput) (*models.FSRSParameter, error) {
	if s == nil || s.api == nil {
		return nil, fsrs.ErrUnavailable
	}
	if in.ResultID == "" || in.ResultName == "" {
		return nil, fmt.Errorf("%w: result_id and result_name are required", fsrs.ErrInvalidInput)
	}

	req, err := OptimizeParametersRequest(in.RequestID, in.DatasetID, in.KnowledgeID, in.ReviewLogs, in.Baseline)
	if err != nil {
		return nil, err
	}

	resp, err := s.api.OptimizeParameters(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.GetOptimizedSetting() == nil {
		return nil, fmt.Errorf("fsrs: empty OptimizeParameters response")
	}

	return OptimizedSettingToModel(in.ResultID, in.ResultName, resp.GetOptimizedSetting())
}
