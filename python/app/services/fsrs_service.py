"""FSRS scheduling service with runtime setting override support."""

from __future__ import annotations

import dataclasses
import json
import threading
from dataclasses import dataclass
from datetime import UTC, datetime, timedelta

from google.protobuf.timestamp_pb2 import Timestamp
from kmemo.v1 import kmemo_pb2

_DEFAULT_PARAMETERS = [
    0.4,
    0.6,
    2.4,
    5.8,
    4.93,
    0.94,
    0.86,
    0.01,
    1.49,
    0.14,
    0.94,
    2.18,
    0.05,
    0.34,
    1.26,
    0.29,
    2.61,
    0.57,
    0.0,
]
_DEFAULT_DESIRED_RETENTION = 0.9
_DEFAULT_MAXIMUM_INTERVAL = 36500
_RATING_TO_REVIEW = {
    1: "again",
    2: "hard",
    3: "good",
    4: "easy",
}
_RATING_TO_DAYS = {
    1: 0.0,
    2: 1.0,
    3: 3.0,
    4: 7.0,
}
_RATING_TO_STABILITY_FACTOR = {
    1: 0.7,
    2: 1.0,
    3: 1.3,
    4: 1.6,
}
_RATING_TO_DIFFICULTY_DELTA = {
    1: 0.15,
    2: 0.05,
    3: -0.05,
    4: -0.15,
}
_VALID_STATES = {"new", "learning", "review", "relearning"}


@dataclass(frozen=True)
class SchedulerSettingModel:
    parameters: tuple[float, ...]
    desired_retention: float
    maximum_interval: int


@dataclass(frozen=True)
class CardStateModel:
    state: str
    due: datetime | None
    last_review: datetime | None
    stability: float | None
    difficulty: float | None
    elapsed_days: float | None
    scheduled_days: float | None
    reps: int
    lapses: int
    learning_steps_json: str | None


@dataclass(frozen=True)
class ReviewInputModel:
    rating: int
    reviewed_at: datetime


@dataclass(frozen=True)
class ReviewLogModel:
    rating: int
    review: str
    elapsed_days: float
    scheduled_days: float
    previous_stability: float | None
    previous_difficulty: float | None
    new_stability: float | None
    new_difficulty: float | None


@dataclass(frozen=True)
class ScheduleResultModel:
    next_card: CardStateModel
    review_log: ReviewLogModel
    effective_setting: SchedulerSettingModel


class SchedulerRuntime:
    def __init__(self) -> None:
        self._lock = threading.RLock()
        self._default_setting = SchedulerSettingModel(
            parameters=tuple(_DEFAULT_PARAMETERS),
            desired_retention=_DEFAULT_DESIRED_RETENTION,
            maximum_interval=_DEFAULT_MAXIMUM_INTERVAL,
        )
        self._current_setting: SchedulerSettingModel | None = None

    def set_setting(self, setting: SchedulerSettingModel) -> SchedulerSettingModel:
        with self._lock:
            self._current_setting = setting
            return self._current_setting

    def get_effective_setting(self) -> SchedulerSettingModel:
        with self._lock:
            if self._current_setting is None:
                return dataclasses.replace(self._default_setting)
            return dataclasses.replace(self._current_setting)


_RUNTIME = SchedulerRuntime()


def scheduler_set_setting(
    request: kmemo_pb2.SchedulerSetSettingRequest,
) -> kmemo_pb2.SchedulerSetSettingResponse:
    setting = _setting_from_proto(request.setting)
    effective = _RUNTIME.set_setting(setting)
    return kmemo_pb2.SchedulerSetSettingResponse(
        ok=True,
        message="scheduler setting updated",
        effective_setting=_setting_to_proto(effective),
    )


def scheduler_schedule(
    request: kmemo_pb2.SchedulerScheduleRequest,
) -> kmemo_pb2.SchedulerScheduleResponse:
    current = _card_state_from_proto(request.current)
    review = _review_input_from_proto(request.review)
    effective_setting = _RUNTIME.get_effective_setting()
    result = _schedule(current, review, effective_setting)
    return kmemo_pb2.SchedulerScheduleResponse(
        ok=True,
        message="scheduler calculation completed",
        next=_card_state_to_proto(result.next_card),
        review_log=_review_log_to_proto(result.review_log),
        effective_setting=_setting_to_proto(result.effective_setting),
    )


def _setting_from_proto(setting: kmemo_pb2.SchedulerSetting) -> SchedulerSettingModel:
    parameters = tuple(setting.parameters) if setting.parameters else tuple(_DEFAULT_PARAMETERS)
    desired_retention = (
        setting.desired_retention if setting.HasField("desired_retention") else _DEFAULT_DESIRED_RETENTION
    )
    maximum_interval = (
        setting.maximum_interval if setting.HasField("maximum_interval") else _DEFAULT_MAXIMUM_INTERVAL
    )
    if not parameters:
        raise ValueError("parameters must not be empty")
    if desired_retention <= 0 or desired_retention >= 1:
        raise ValueError("desired_retention must be between 0 and 1")
    if maximum_interval <= 0:
        raise ValueError("maximum_interval must be positive")
    return SchedulerSettingModel(
        parameters=parameters,
        desired_retention=desired_retention,
        maximum_interval=maximum_interval,
    )


def _setting_to_proto(setting: SchedulerSettingModel) -> kmemo_pb2.SchedulerSetting:
    proto = kmemo_pb2.SchedulerSetting(parameters=list(setting.parameters))
    proto.desired_retention = setting.desired_retention
    proto.maximum_interval = setting.maximum_interval
    return proto


def _card_state_from_proto(card: kmemo_pb2.CardState) -> CardStateModel:
    state = card.state or "new"
    if state not in _VALID_STATES:
        raise ValueError(f"invalid state: {state}")
    return CardStateModel(
        state=state,
        due=_timestamp_to_datetime(card.due) if card.HasField("due") else None,
        last_review=_timestamp_to_datetime(card.last_review) if card.HasField("last_review") else None,
        stability=card.stability if card.HasField("stability") else None,
        difficulty=card.difficulty if card.HasField("difficulty") else None,
        elapsed_days=card.elapsed_days if card.HasField("elapsed_days") else None,
        scheduled_days=card.scheduled_days if card.HasField("scheduled_days") else None,
        reps=card.reps,
        lapses=card.lapses,
        learning_steps_json=card.learning_steps_json if card.HasField("learning_steps_json") else None,
    )


def _card_state_to_proto(card: CardStateModel) -> kmemo_pb2.CardState:
    proto = kmemo_pb2.CardState(
        state=card.state,
        reps=card.reps,
        lapses=card.lapses,
    )
    if card.due is not None:
        proto.due.CopyFrom(_datetime_to_timestamp(card.due))
    if card.last_review is not None:
        proto.last_review.CopyFrom(_datetime_to_timestamp(card.last_review))
    if card.stability is not None:
        proto.stability = card.stability
    if card.difficulty is not None:
        proto.difficulty = card.difficulty
    if card.elapsed_days is not None:
        proto.elapsed_days = card.elapsed_days
    if card.scheduled_days is not None:
        proto.scheduled_days = card.scheduled_days
    if card.learning_steps_json is not None:
        proto.learning_steps_json = card.learning_steps_json
    return proto


def _review_input_from_proto(review: kmemo_pb2.ReviewInput) -> ReviewInputModel:
    if review.rating not in _RATING_TO_REVIEW:
        raise ValueError(f"invalid rating: {review.rating}")
    if not review.HasField("reviewed_at"):
        raise ValueError("reviewed_at is required")
    return ReviewInputModel(
        rating=review.rating,
        reviewed_at=_timestamp_to_datetime(review.reviewed_at),
    )


def _review_log_to_proto(review_log: ReviewLogModel) -> kmemo_pb2.ReviewLogSnapshot:
    proto = kmemo_pb2.ReviewLogSnapshot(
        rating=review_log.rating,
        review=review_log.review,
        elapsed_days=review_log.elapsed_days,
        scheduled_days=review_log.scheduled_days,
    )
    if review_log.previous_stability is not None:
        proto.previous_stability = review_log.previous_stability
    if review_log.previous_difficulty is not None:
        proto.previous_difficulty = review_log.previous_difficulty
    if review_log.new_stability is not None:
        proto.new_stability = review_log.new_stability
    if review_log.new_difficulty is not None:
        proto.new_difficulty = review_log.new_difficulty
    return proto


def _schedule(
    current: CardStateModel,
    review: ReviewInputModel,
    setting: SchedulerSettingModel,
) -> ScheduleResultModel:
    previous_stability = current.stability or 1.0
    previous_difficulty = current.difficulty or 5.0
    elapsed_days = _resolve_elapsed_days(current, review.reviewed_at)
    interval_factor = _interval_factor(setting)
    scheduled_days = min(
        _RATING_TO_DAYS[review.rating] * interval_factor,
        float(setting.maximum_interval),
    )
    new_stability = max(
        0.1,
        previous_stability * _RATING_TO_STABILITY_FACTOR[review.rating] * interval_factor,
    )
    new_difficulty = min(
        10.0,
        max(1.0, previous_difficulty + _RATING_TO_DIFFICULTY_DELTA[review.rating]),
    )
    next_due = review.reviewed_at + timedelta(days=scheduled_days)
    next_state = _next_state(current.state, review.rating)
    next_card = CardStateModel(
        state=next_state,
        due=next_due,
        last_review=review.reviewed_at,
        stability=new_stability,
        difficulty=new_difficulty,
        elapsed_days=elapsed_days,
        scheduled_days=scheduled_days,
        reps=current.reps + 1,
        lapses=current.lapses + (1 if review.rating == 1 else 0),
        learning_steps_json=_normalize_learning_steps(current.learning_steps_json),
    )
    review_log = ReviewLogModel(
        rating=review.rating,
        review=_RATING_TO_REVIEW[review.rating],
        elapsed_days=elapsed_days,
        scheduled_days=scheduled_days,
        previous_stability=current.stability,
        previous_difficulty=current.difficulty,
        new_stability=new_stability,
        new_difficulty=new_difficulty,
    )
    return ScheduleResultModel(
        next_card=next_card,
        review_log=review_log,
        effective_setting=setting,
    )


def _interval_factor(setting: SchedulerSettingModel) -> float:
    parameter_head = setting.parameters[:4] if setting.parameters else tuple(_DEFAULT_PARAMETERS[:4])
    parameter_mean = sum(parameter_head) / len(parameter_head)
    retention_factor = 0.5 + setting.desired_retention
    return max(0.25, min(4.0, retention_factor * (parameter_mean / 2.0)))


def _resolve_elapsed_days(current: CardStateModel, reviewed_at: datetime) -> float:
    if current.last_review is None:
        return 0.0
    delta = reviewed_at - current.last_review
    return max(0.0, delta.total_seconds() / 86400)


def _next_state(current_state: str, rating: int) -> str:
    if rating == 1:
        return "relearning" if current_state in {"review", "relearning"} else "learning"
    if current_state == "new":
        return "learning"
    if current_state == "learning" and rating >= 3:
        return "review"
    return "review" if current_state == "review" else current_state


def _normalize_learning_steps(raw: str | None) -> str:
    if raw is None:
        return json.dumps([])
    try:
        json.loads(raw)
        return raw
    except json.JSONDecodeError:
        return json.dumps([])


def _timestamp_to_datetime(value: Timestamp) -> datetime:
    return value.ToDatetime().astimezone(UTC)


def _datetime_to_timestamp(value: datetime) -> Timestamp:
    ts = Timestamp()
    ts.FromDatetime(value.astimezone(UTC))
    return ts
