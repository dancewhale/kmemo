from __future__ import annotations

from datetime import UTC, datetime, timedelta
from pathlib import Path
import sys

_ROOT = Path(__file__).resolve().parent.parent
_GEN_ROOT = _ROOT / "generated"
if str(_GEN_ROOT) not in sys.path:
    sys.path.insert(0, str(_GEN_ROOT))
if str(_ROOT) not in sys.path:
    sys.path.insert(0, str(_ROOT))

from google.protobuf.timestamp_pb2 import Timestamp
from kmemo.v1 import kmemo_pb2

from app.services import fsrs_service


_VALID_PARAMETERS = list(fsrs_service._DEFAULT_PARAMETERS)


def _timestamp(value: datetime) -> Timestamp:
    ts = Timestamp()
    ts.FromDatetime(value.astimezone(UTC))
    return ts


def _card_state(
    *,
    state: str = "review",
    due: datetime | None = None,
    last_review: datetime | None = None,
    stability: float = 3.2,
    difficulty: float = 5.6,
    elapsed_days: float = 3.0,
    scheduled_days: float = 3.0,
    reps: int = 3,
    lapses: int = 0,
) -> kmemo_pb2.CardState:
    card = kmemo_pb2.CardState(
        state=state,
        reps=reps,
        lapses=lapses,
        stability=stability,
        difficulty=difficulty,
        elapsed_days=elapsed_days,
        scheduled_days=scheduled_days,
        learning_steps_json="[]",
    )
    if due is not None:
        card.due.CopyFrom(_timestamp(due))
    if last_review is not None:
        card.last_review.CopyFrom(_timestamp(last_review))
    return card


def _review_request(rating: int = 3) -> kmemo_pb2.ReviewCardRequest:
    reviewed_at = datetime(2026, 4, 2, 12, 0, tzinfo=UTC)
    return kmemo_pb2.ReviewCardRequest(
        card_id="card-1",
        card=_card_state(
            due=reviewed_at,
            last_review=reviewed_at - timedelta(days=3),
        ),
        review=kmemo_pb2.ReviewInput(
            rating=rating,
            reviewed_at=_timestamp(reviewed_at),
        ),
    )


def _reschedule_request() -> kmemo_pb2.RescheduleCardRequest:
    reviewed_at = datetime(2026, 4, 2, 12, 0, tzinfo=UTC)
    return kmemo_pb2.RescheduleCardRequest(
        card_id="card-1",
        card=_card_state(
            due=reviewed_at + timedelta(days=2),
            last_review=reviewed_at - timedelta(days=2),
            scheduled_days=2.0,
            elapsed_days=2.0,
        ),
        reschedule_at=_timestamp(reviewed_at),
        review_logs=[
            kmemo_pb2.ReviewLogSnapshot(
                rating=3,
                review="good",
                elapsed_days=2.0,
                scheduled_days=2.0,
                reviewed_at=_timestamp(reviewed_at - timedelta(days=2)),
            ),
            kmemo_pb2.ReviewLogSnapshot(
                rating=3,
                review="good",
                elapsed_days=3.0,
                scheduled_days=3.0,
                reviewed_at=_timestamp(reviewed_at - timedelta(days=1)),
            ),
        ],
        mode="keep-due-close",
    )


def _optimizer_request() -> kmemo_pb2.OptimizeParametersRequest:
    base = datetime(2026, 3, 1, 9, 0, tzinfo=UTC)
    entries: list[kmemo_pb2.OptimizerReviewLogEntry] = []
    for card_index in range(1, 4):
        review_time = base
        for step, rating in enumerate((3, 3, 1, 3), start=1):
            before = _card_state(
                state="learning" if step == 1 else "review",
                due=review_time,
                last_review=review_time - timedelta(days=max(step, 1)),
                stability=2.0 + step,
                difficulty=5.0,
                elapsed_days=float(step),
                scheduled_days=float(step),
                reps=step,
                lapses=1 if rating == 1 else 0,
            )
            after = _card_state(
                state="relearning" if rating == 1 else "review",
                due=review_time + timedelta(days=step + 1),
                last_review=review_time,
                stability=2.5 + step,
                difficulty=5.5,
                elapsed_days=float(step),
                scheduled_days=float(step + 1),
                reps=step + 1,
                lapses=1 if rating == 1 else 0,
            )
            entries.append(
                kmemo_pb2.OptimizerReviewLogEntry(
                    card_id=f"card-{card_index}",
                    rating=rating,
                    reviewed_at=_timestamp(review_time),
                    card_before=before,
                    card_after=after,
                )
            )
            review_time = review_time + timedelta(days=step + 1)
    return kmemo_pb2.OptimizeParametersRequest(
        request_id="req-1",
        dataset_id="dataset-1",
        knowledge_id="knowledge-1",
        review_logs=entries,
        baseline_setting=kmemo_pb2.SchedulerSetting(
            parameters=_VALID_PARAMETERS,
            desired_retention=0.9,
            maximum_interval=365,
        ),
        config=kmemo_pb2.OptimizerConfig(
            max_iterations=1,
            enable_short_term=True,
        ),
    )


def test_review_card_works_without_setting():
    response = fsrs_service.review_card(_review_request())

    assert response.card.state in {"learning", "review", "relearning"}
    assert response.effective_setting.parameters
    assert response.effective_setting.desired_retention > 0
    assert response.effective_setting.maximum_interval > 0
    assert response.review_log.rating == 3
    assert response.review_log.review == "good"
    assert response.retrievability > 0


def test_scheduler_setting_overrides_previous_values():
    first = fsrs_service.scheduler_set_setting(
        kmemo_pb2.SchedulerSetSettingRequest(
            setting=kmemo_pb2.SchedulerSetting(
                parameters=_VALID_PARAMETERS,
                desired_retention=0.83,
                maximum_interval=100,
            )
        )
    )
    second = fsrs_service.scheduler_set_setting(
        kmemo_pb2.SchedulerSetSettingRequest(
            setting=kmemo_pb2.SchedulerSetting(
                parameters=_VALID_PARAMETERS,
                desired_retention=0.97,
                maximum_interval=200,
            )
        )
    )
    reviewed = fsrs_service.review_card(_review_request())

    assert first.effective_setting.desired_retention == 0.83
    assert second.effective_setting.desired_retention == 0.97
    assert reviewed.effective_setting.parameters == _VALID_PARAMETERS
    assert reviewed.effective_setting.desired_retention == 0.97
    assert reviewed.effective_setting.maximum_interval == 200


def test_get_card_retrievability_returns_value():
    now = datetime(2026, 4, 2, 12, 0, tzinfo=UTC)
    response = fsrs_service.get_card_retrievability(
        kmemo_pb2.GetCardRetrievabilityRequest(
            card_id="card-1",
            card=_card_state(
                due=now + timedelta(days=1),
                last_review=now - timedelta(days=3),
            ),
            now=_timestamp(now),
        )
    )

    assert 0 < response.retrievability <= 1
    assert response.evaluated_at.seconds > 0
    assert response.effective_setting.parameters


def test_reschedule_card_returns_card_and_warning():
    response = fsrs_service.reschedule_card(_reschedule_request())

    assert response.card.state in {"learning", "review", "relearning"}
    assert response.review_log.review == "reschedule"
    assert response.review_log.scheduled_days >= 0
    assert response.retrievability > 0
    assert response.warnings


def test_optimize_parameters_returns_metrics():
    response = fsrs_service.optimize_parameters(_optimizer_request())

    assert len(response.optimized_setting.parameters) == 21
    assert response.sample_count > 0
    assert response.metrics.rmse_bins >= 0
    assert response.metrics.log_loss >= 0
    assert response.dropped_sample_count >= 0
