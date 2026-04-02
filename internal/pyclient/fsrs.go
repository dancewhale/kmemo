package pyclient

import (
	"encoding/json"
	"fmt"
	"time"

	kmemov1 "kmemo/gen/kmemo/v1"
	"kmemo/internal/storage/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func SettingFromFSRSParameter(param *models.FSRSParameter) (*kmemov1.SchedulerSetting, error) {
	if param == nil {
		return nil, fmt.Errorf("fsrs parameter is nil")
	}

	var parameters []float64
	if err := json.Unmarshal([]byte(param.ParametersJSON), &parameters); err != nil {
		return nil, fmt.Errorf("parse fsrs parameters json: %w", err)
	}

	setting := &kmemov1.SchedulerSetting{
		Parameters: parameters,
	}
	if param.DesiredRetention != nil {
		setting.DesiredRetention = param.DesiredRetention
	}
	if param.MaximumInterval != nil {
		v := int32(*param.MaximumInterval)
		setting.MaximumInterval = &v
	}
	return setting, nil
}

func CardStateFromModel(card *models.CardSRS) *kmemov1.CardState {
	if card == nil {
		return &kmemov1.CardState{State: "new"}
	}
	state := &kmemov1.CardState{
		State:  card.FSRSState,
		Reps:   int32(card.Reps),
		Lapses: int32(card.Lapses),
	}
	if card.DueAt != nil {
		state.Due = timestamppb.New(card.DueAt.UTC())
	}
	if card.LastReviewAt != nil {
		state.LastReview = timestamppb.New(card.LastReviewAt.UTC())
	}
	if card.Stability != nil {
		state.Stability = card.Stability
	}
	if card.Difficulty != nil {
		state.Difficulty = card.Difficulty
	}
	if card.ElapsedDays != nil {
		state.ElapsedDays = card.ElapsedDays
	}
	if card.ScheduledDays != nil {
		state.ScheduledDays = card.ScheduledDays
	}
	if card.LearningStepsJSON != nil {
		state.LearningStepsJson = card.LearningStepsJSON
	}
	return state
}

func ReviewInput(rating int, reviewedAt time.Time) *kmemov1.ReviewInput {
	return &kmemov1.ReviewInput{
		Rating:     int32(rating),
		ReviewedAt: timestamppb.New(reviewedAt.UTC()),
	}
}

func RetrievabilityRequestFromModel(card *models.CardSRS, now time.Time) *kmemov1.GetCardRetrievabilityRequest {
	cardID := ""
	if card != nil {
		cardID = card.CardID
	}
	return &kmemov1.GetCardRetrievabilityRequest{
		CardId: cardID,
		Card:   CardStateFromModel(card),
		Now:    timestamppb.New(now.UTC()),
	}
}

func ReviewCardRequestFromModel(card *models.CardSRS, rating int, reviewedAt time.Time) *kmemov1.ReviewCardRequest {
	cardID := ""
	if card != nil {
		cardID = card.CardID
	}
	return &kmemov1.ReviewCardRequest{
		CardId: cardID,
		Card:   CardStateFromModel(card),
		Review: ReviewInput(rating, reviewedAt),
	}
}

func RescheduleCardRequestFromModel(card *models.CardSRS, reviewLogs []*models.ReviewLog, rescheduleAt time.Time, mode string) *kmemov1.RescheduleCardRequest {
	cardID := ""
	if card != nil {
		cardID = card.CardID
	}
	request := &kmemov1.RescheduleCardRequest{
		CardId:       cardID,
		Card:         CardStateFromModel(card),
		RescheduleAt: timestamppb.New(rescheduleAt.UTC()),
		ReviewLogs:   ReviewLogsToProto(reviewLogs),
	}
	if mode != "" {
		request.Mode = &mode
	}
	return request
}

func OptimizeParametersRequestFromModels(requestID, datasetID, knowledgeID string, reviewLogs []*models.ReviewLog, baseline *models.FSRSParameter) (*kmemov1.OptimizeParametersRequest, error) {
	request := &kmemov1.OptimizeParametersRequest{
		RequestId:  requestID,
		DatasetId:  datasetID,
		KnowledgeId: knowledgeID,
		ReviewLogs: ReviewLogsToOptimizerEntries(reviewLogs),
	}
	if baseline != nil {
		setting, err := SettingFromFSRSParameter(baseline)
		if err != nil {
			return nil, err
		}
		request.BaselineSetting = setting
	}
	return request, nil
}

func CardStateToModel(cardID string, state *kmemov1.CardState) *models.CardSRS {
	if state == nil {
		return &models.CardSRS{CardID: cardID}
	}
	model := &models.CardSRS{
		CardID:    cardID,
		FSRSState: state.State,
		Reps:      int(state.Reps),
		Lapses:    int(state.Lapses),
	}
	if state.Due != nil {
		t := state.Due.AsTime().UTC()
		model.DueAt = &t
	}
	if state.LastReview != nil {
		t := state.LastReview.AsTime().UTC()
		model.LastReviewAt = &t
	}
	if state.Stability != nil {
		v := state.GetStability()
		model.Stability = &v
	}
	if state.Difficulty != nil {
		v := state.GetDifficulty()
		model.Difficulty = &v
	}
	if state.ElapsedDays != nil {
		v := state.GetElapsedDays()
		model.ElapsedDays = &v
	}
	if state.ScheduledDays != nil {
		v := state.GetScheduledDays()
		model.ScheduledDays = &v
	}
	if state.LearningStepsJson != nil {
		v := state.GetLearningStepsJson()
		model.LearningStepsJSON = &v
	}
	return model
}

func ReviewLogFromProto(cardID string, log *kmemov1.ReviewLogSnapshot) *models.ReviewLog {
	if log == nil {
		return &models.ReviewLog{CardID: cardID}
	}
	review := &models.ReviewLog{
		CardID:             cardID,
		Rating:             int(log.Rating),
		ReviewKind:         log.Review,
		ScheduledDays:      log.ScheduledDays,
		ElapsedDays:        log.ElapsedDays,
		PreviousStability:  log.PreviousStability,
		PreviousDifficulty: log.PreviousDifficulty,
		NewStability:       log.NewStability,
		NewDifficulty:      log.NewDifficulty,
		SnapshotJSON:       "{}",
	}
	if log.ReviewedAt != nil {
		review.ReviewedAt = log.ReviewedAt.AsTime().UTC()
	}
	return review
}

func ReviewLogsToProto(logs []*models.ReviewLog) []*kmemov1.ReviewLogSnapshot {
	result := make([]*kmemov1.ReviewLogSnapshot, 0, len(logs))
	for _, log := range logs {
		if log == nil {
			continue
		}
		item := &kmemov1.ReviewLogSnapshot{
			Rating:        int32(log.Rating),
			Review:        log.ReviewKind,
			ScheduledDays: log.ScheduledDays,
			ElapsedDays:   log.ElapsedDays,
			ReviewedAt:    timestamppb.New(log.ReviewedAt.UTC()),
		}
		if log.PreviousStability != nil {
			item.PreviousStability = log.PreviousStability
		}
		if log.PreviousDifficulty != nil {
			item.PreviousDifficulty = log.PreviousDifficulty
		}
		if log.NewStability != nil {
			item.NewStability = log.NewStability
		}
		if log.NewDifficulty != nil {
			item.NewDifficulty = log.NewDifficulty
		}
		result = append(result, item)
	}
	return result
}

func ReviewLogsToOptimizerEntries(logs []*models.ReviewLog) []*kmemov1.OptimizerReviewLogEntry {
	result := make([]*kmemov1.OptimizerReviewLogEntry, 0, len(logs))
	for _, log := range logs {
		if log == nil {
			continue
		}
		result = append(result, &kmemov1.OptimizerReviewLogEntry{
			CardId:     log.CardID,
			Rating:     int32(log.Rating),
			ReviewedAt: timestamppb.New(log.ReviewedAt.UTC()),
			CardBefore: &kmemov1.CardState{State: "review"},
			CardAfter: &kmemov1.CardState{
				State:         optimizerState(log.ReviewKind),
				ScheduledDays: &log.ScheduledDays,
				ElapsedDays:   &log.ElapsedDays,
			},
		})
	}
	return result
}

func OptimizedSettingToModel(id, name string, setting *kmemov1.SchedulerSetting) (*models.FSRSParameter, error) {
	if setting == nil {
		return nil, fmt.Errorf("optimized setting is nil")
	}
	payload, err := json.Marshal(setting.Parameters)
	if err != nil {
		return nil, fmt.Errorf("marshal optimized parameters: %w", err)
	}
	result := &models.FSRSParameter{
		ID:             id,
		Name:           name,
		ParametersJSON: string(payload),
	}
	if setting.DesiredRetention != nil {
		v := setting.GetDesiredRetention()
		result.DesiredRetention = &v
	}
	if setting.MaximumInterval != nil {
		v := int(setting.GetMaximumInterval())
		result.MaximumInterval = &v
	}
	return result, nil
}

func optimizerState(reviewKind string) string {
	switch reviewKind {
	case "again", "relearning", "reschedule":
		return "relearning"
	case "hard", "good", "easy", "review":
		return "review"
	default:
		return "review"
	}
}
