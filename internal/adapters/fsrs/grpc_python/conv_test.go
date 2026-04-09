package grpcpython

import (
	"testing"
	"time"

	kmemov1 "kmemo/gen/kmemo/v1"
	"kmemo/internal/storage/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestReviewCardRequest(t *testing.T) {
	now := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)
	stability := 3.2
	difficulty := 5.4
	elapsedDays := 2.0
	scheduledDays := 4.0
	card := &models.CardSRS{
		CardID:        "card-1",
		FSRSState:     "review",
		LastReviewAt:  ptrTime(now.Add(-48 * time.Hour)),
		DueAt:         ptrTime(now.Add(24 * time.Hour)),
		Stability:     &stability,
		Difficulty:    &difficulty,
		ElapsedDays:   &elapsedDays,
		ScheduledDays: &scheduledDays,
		Reps:          5,
		Lapses:        1,
	}

	request := ReviewCardRequest("card-1", card, 3, now)
	if request.CardId != "card-1" {
		t.Fatalf("card id = %q, want card-1", request.CardId)
	}
	if request.Review.GetRating() != 3 {
		t.Fatalf("rating = %d, want 3", request.Review.GetRating())
	}
	if request.Card.GetState() != "review" {
		t.Fatalf("state = %q, want review", request.Card.GetState())
	}
}

func TestRetrievabilityRequest(t *testing.T) {
	now := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)
	request := RetrievabilityRequest("card-1", &models.CardSRS{CardID: "card-1", FSRSState: "review"}, now)
	if request.GetCardId() != "card-1" {
		t.Fatalf("card id = %q, want card-1", request.GetCardId())
	}
	if request.GetNow() == nil {
		t.Fatal("expected now timestamp")
	}
}

func TestReviewLogFromProto(t *testing.T) {
	now := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)
	previousStability := 2.2
	responseLog := &kmemov1.ReviewLogSnapshot{
		Rating:            3,
		Review:            "good",
		ElapsedDays:       2,
		ScheduledDays:     5,
		PreviousStability: &previousStability,
		ReviewedAt:        ts(now),
	}

	model := ReviewLogFromProto("card-1", responseLog)
	if model.CardID != "card-1" {
		t.Fatalf("card id = %q, want card-1", model.CardID)
	}
	if !model.ReviewedAt.Equal(now) {
		t.Fatalf("reviewed at = %v, want %v", model.ReviewedAt, now)
	}
	if model.ReviewKind != "good" {
		t.Fatalf("review kind = %q, want good", model.ReviewKind)
	}
}

func TestOptimizedSettingToModel(t *testing.T) {
	desiredRetention := 0.9
	maximumInterval := int32(365)
	setting := &kmemov1.SchedulerSetting{
		Parameters:       []float64{0.1, 0.2, 0.3},
		DesiredRetention: &desiredRetention,
		MaximumInterval:  &maximumInterval,
	}

	model, err := OptimizedSettingToModel("preset-1", "preset", setting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if model.ID != "preset-1" {
		t.Fatalf("id = %q, want preset-1", model.ID)
	}
	if model.DesiredRetention == nil || *model.DesiredRetention != desiredRetention {
		t.Fatalf("desired retention = %v, want %v", model.DesiredRetention, desiredRetention)
	}
}

func ptrTime(v time.Time) *time.Time { return &v }

func ts(v time.Time) *timestamppb.Timestamp { return timestamppb.New(v.UTC()) }
