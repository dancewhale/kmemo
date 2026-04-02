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
		CardID:            cardID,
		Rating:            int(log.Rating),
		ReviewKind:        log.Review,
		ScheduledDays:     log.ScheduledDays,
		ElapsedDays:       log.ElapsedDays,
		PreviousStability: log.PreviousStability,
		PreviousDifficulty: log.PreviousDifficulty,
		NewStability:      log.NewStability,
		NewDifficulty:     log.NewDifficulty,
		SnapshotJSON:      "{}",
	}
	return review
}
