package grpcpython

import (
	"encoding/json"
	"fmt"
	"time"

	kmemov1 "kmemo/gen/kmemo/v1"
	"kmemo/internal/storage/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// --- Outbound: models / app types → protobuf (requests) ---

// SettingFromFSRSParameter maps a stored FSRS preset into a protobuf SchedulerSetting.
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

// CardStateFromModel maps CardSRS into protobuf CardState for RPC requests.
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

func reviewInputProto(rating int, reviewedAt time.Time) *kmemov1.ReviewInput {
	return &kmemov1.ReviewInput{
		Rating:     int32(rating),
		ReviewedAt: timestamppb.New(reviewedAt.UTC()),
	}
}

// RetrievabilityRequest builds GetCardRetrievabilityRequest. cardID may be empty if prior carries CardID.
func RetrievabilityRequest(cardID string, prior *models.CardSRS, now time.Time) *kmemov1.GetCardRetrievabilityRequest {
	cardID = effectiveCardID(cardID, prior)
	return &kmemov1.GetCardRetrievabilityRequest{
		CardId: cardID,
		Card:   CardStateFromModel(prior),
		Now:    timestamppb.New(now.UTC()),
	}
}

// ReviewCardRequest builds ReviewCardRequest. cardID may be empty if prior carries CardID.
func ReviewCardRequest(cardID string, prior *models.CardSRS, rating int, reviewedAt time.Time) *kmemov1.ReviewCardRequest {
	cardID = effectiveCardID(cardID, prior)
	return &kmemov1.ReviewCardRequest{
		CardId: cardID,
		Card:   CardStateFromModel(prior),
		Review: reviewInputProto(rating, reviewedAt),
	}
}

// RescheduleCardRequest builds RescheduleCardRequest.
func RescheduleCardRequest(cardID string, prior *models.CardSRS, reviewLogs []*models.ReviewLog, rescheduleAt time.Time, mode string) *kmemov1.RescheduleCardRequest {
	cardID = effectiveCardID(cardID, prior)
	request := &kmemov1.RescheduleCardRequest{
		CardId:       cardID,
		Card:         CardStateFromModel(prior),
		RescheduleAt: timestamppb.New(rescheduleAt.UTC()),
		ReviewLogs:   reviewLogsToProto(reviewLogs),
	}
	if mode != "" {
		request.Mode = &mode
	}
	return request
}

// OptimizeParametersRequest builds OptimizeParametersRequest.
func OptimizeParametersRequest(requestID, datasetID, knowledgeID string, reviewLogs []*models.ReviewLog, baseline *models.FSRSParameter) (*kmemov1.OptimizeParametersRequest, error) {
	request := &kmemov1.OptimizeParametersRequest{
		RequestId:   requestID,
		DatasetId:   datasetID,
		KnowledgeId: knowledgeID,
		ReviewLogs:  reviewLogsToOptimizerEntries(reviewLogs),
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

func effectiveCardID(explicit string, prior *models.CardSRS) string {
	if explicit != "" {
		return explicit
	}
	if prior != nil {
		return prior.CardID
	}
	return ""
}

// --- Inbound: protobuf → models (responses) ---

// CardStateToModel maps protobuf CardState into CardSRS for persistence.
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

// ReviewLogFromProto maps ReviewLogSnapshot into ReviewLog (SnapshotJSON defaults to "{}" if unset upstream).
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

func reviewLogsToProto(logs []*models.ReviewLog) []*kmemov1.ReviewLogSnapshot {
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

func reviewLogsToOptimizerEntries(logs []*models.ReviewLog) []*kmemov1.OptimizerReviewLogEntry {
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

// OptimizedSettingToModel maps an optimized SchedulerSetting into FSRSParameter for storage.
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
