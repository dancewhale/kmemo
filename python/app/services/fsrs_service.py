"""FSRS scheduling service backed by the real fsrs/fsrs-optimizer libraries."""

from __future__ import annotations

import dataclasses
import hashlib
import json
import threading
from dataclasses import dataclass
from datetime import UTC, datetime

import pandas as pd
from fsrs import Card, Rating, ReviewLog, Scheduler, State
from fsrs_optimizer import Optimizer
from google.protobuf.timestamp_pb2 import Timestamp
from kmemo.v1 import kmemo_pb2

_DEFAULT_SCHEDULER = Scheduler()
_DEFAULT_PARAMETERS = tuple(float(value) for value in _DEFAULT_SCHEDULER.parameters)
_DEFAULT_DESIRED_RETENTION = float(_DEFAULT_SCHEDULER.desired_retention)
_DEFAULT_MAXIMUM_INTERVAL = int(_DEFAULT_SCHEDULER.maximum_interval)
_DEFAULT_LEARNING_STEPS_JSON = json.dumps([])
_VALID_STATES = {"new", "learning", "review", "relearning"}
_RATING_TO_REVIEW = {
    1: "again",
    2: "hard",
    3: "good",
    4: "easy",
}
_PROTO_STATE_TO_FSRS_STATE = {
    "new": State.Learning,
    "learning": State.Learning,
    "review": State.Review,
    "relearning": State.Relearning,
}
_FSRS_STATE_TO_PROTO_STATE = {
    State.Learning: "learning",
    State.Review: "review",
    State.Relearning: "relearning",
}
_RATING_TO_FSRS = {
    1: Rating.Again,
    2: Rating.Hard,
    3: Rating.Good,
    4: Rating.Easy,
}
_PROTO_STATE_TO_OPTIMIZER_STATE = {
    "new": 0,
    "learning": 1,
    "review": 2,
    "relearning": 3,
}


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
    reviewed_at: datetime | None


@dataclass(frozen=True)
class ReviewResultModel:
    card: CardStateModel
    review_log: ReviewLogModel
    retrievability: float | None
    effective_setting: SchedulerSettingModel
    warnings: tuple[str, ...] = ()
    diagnostics: str | None = None


@dataclass(frozen=True)
class RetrievabilityResultModel:
    retrievability: float
    evaluated_at: datetime
    effective_setting: SchedulerSettingModel
    warnings: tuple[str, ...] = ()
    diagnostics: str | None = None


@dataclass(frozen=True)
class OptimizeResultModel:
    optimized_setting: SchedulerSettingModel
    metrics: dict[str, float]
    sample_count: int
    dropped_sample_count: int
    warnings: tuple[str, ...] = ()
    diagnostics: str | None = None


class SchedulerRuntime:
    def __init__(self) -> None:
        self._lock = threading.RLock()
        self._default_setting = SchedulerSettingModel(
            parameters=_DEFAULT_PARAMETERS,
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
        effective_setting=_setting_to_proto(effective),
    )


def get_card_retrievability(
    request: kmemo_pb2.GetCardRetrievabilityRequest,
) -> kmemo_pb2.GetCardRetrievabilityResponse:
    card = _card_state_from_proto(request.card)
    if card.state == "new":
        raise ValueError("retrievability requires a scheduled card state")
    if card.stability is None or card.difficulty is None or card.last_review is None:
        raise ValueError("retrievability requires stability, difficulty, and last_review")
    if not request.HasField("now"):
        raise ValueError("now is required")

    evaluated_at = _timestamp_to_datetime(request.now)
    effective_setting = _RUNTIME.get_effective_setting()
    scheduler = _scheduler_from_setting(effective_setting)
    retrievability = scheduler.get_card_retrievability(
        _fsrs_card_from_model(card, request.card_id),
        current_datetime=evaluated_at,
    )
    return kmemo_pb2.GetCardRetrievabilityResponse(
        retrievability=float(retrievability),
        evaluated_at=_datetime_to_timestamp(evaluated_at),
        effective_setting=_setting_to_proto(effective_setting),
    )


def review_card(
    request: kmemo_pb2.ReviewCardRequest,
) -> kmemo_pb2.ReviewCardResponse:
    current = _card_state_from_proto(request.card)
    review = _review_input_from_proto(request.review)
    effective_setting = _RUNTIME.get_effective_setting()
    scheduler = _scheduler_from_setting(effective_setting)
    next_card, review_log = scheduler.review_card(
        _fsrs_card_from_model(current, request.card_id),
        _RATING_TO_FSRS[review.rating],
        review_datetime=review.reviewed_at,
    )
    result = ReviewResultModel(
        card=_card_state_from_fsrs(next_card, current, review.reviewed_at, review.rating),
        review_log=_review_log_from_fsrs(review_log, current, next_card),
        retrievability=float(
            scheduler.get_card_retrievability(next_card, current_datetime=review.reviewed_at)
        ),
        effective_setting=effective_setting,
    )
    return _review_result_to_proto(result)


def reschedule_card(
    request: kmemo_pb2.RescheduleCardRequest,
) -> kmemo_pb2.RescheduleCardResponse:
    current = _card_state_from_proto(request.card)
    if not request.review_logs:
        raise ValueError("review_logs must not be empty")
    if not request.HasField("reschedule_at"):
        raise ValueError("reschedule_at is required")

    reschedule_at = _timestamp_to_datetime(request.reschedule_at)
    effective_setting = _RUNTIME.get_effective_setting()
    scheduler = _scheduler_from_setting(effective_setting)
    fsrs_card = _fsrs_card_from_model(current, request.card_id)
    review_logs = [_fsrs_review_log_from_proto(item, request.card_id) for item in request.review_logs]
    rescheduled_card = scheduler.reschedule_card(fsrs_card, review_logs)

    warnings: list[str] = []
    if request.HasField("mode") and request.mode:
        warnings.append(f"reschedule mode '{request.mode}' is not interpreted by fsrs and was ignored")

    result_card = _card_state_from_fsrs(
        rescheduled_card,
        current,
        reschedule_at,
        None,
        reps=current.reps,
        lapses=current.lapses,
    )
    result = ReviewResultModel(
        card=result_card,
        review_log=ReviewLogModel(
            rating=0,
            review="reschedule",
            elapsed_days=_resolve_elapsed_days(current, reschedule_at),
            scheduled_days=_resolve_scheduled_days(result_card, reschedule_at),
            previous_stability=current.stability,
            previous_difficulty=current.difficulty,
            new_stability=result_card.stability,
            new_difficulty=result_card.difficulty,
            reviewed_at=reschedule_at,
        ),
        retrievability=float(
            scheduler.get_card_retrievability(rescheduled_card, current_datetime=reschedule_at)
        ),
        effective_setting=effective_setting,
        warnings=tuple(warnings),
    )
    return kmemo_pb2.RescheduleCardResponse(
        card=_card_state_to_proto(result.card),
        review_log=_review_log_to_proto(result.review_log),
        retrievability=result.retrievability,
        effective_setting=_setting_to_proto(result.effective_setting),
        warnings=list(result.warnings),
    )


def optimize_parameters(
    request: kmemo_pb2.OptimizeParametersRequest,
) -> kmemo_pb2.OptimizeParametersResponse:
    dataset = _optimizer_dataset_from_request(request)
    if dataset.empty:
        raise ValueError("review_logs did not produce any optimizer samples")

    optimizer = Optimizer(
        float_delta_t=request.config.float_delta_t if request.HasField("config") and request.config.HasField("float_delta_t") else False,
        enable_short_term=request.config.enable_short_term if request.HasField("config") and request.config.HasField("enable_short_term") else True,
    )
    optimizer.define_model()
    optimizer.initialize_parameters(dataset=dataset.copy(), verbose=False)
    optimizer.train(
        n_epoch=request.config.max_iterations
        if request.HasField("config") and request.config.HasField("max_iterations")
        else 5,
        verbose=False,
    )
    loss_before, loss_after = optimizer.evaluate(save_to_file=False)
    calibration_metrics, figures = optimizer.calibration_graph(verbose=False)
    for figure in figures:
        figure.clf()

    try:
        optimizer.find_optimal_retention(verbose=False)
        optimal_retention = float(optimizer.optimal_retention)
    except Exception:
        optimal_retention = 0.0

    optimized_setting = SchedulerSettingModel(
        parameters=tuple(float(value) for value in optimizer.w),
        desired_retention=(
            request.baseline_setting.desired_retention
            if request.HasField("baseline_setting") and request.baseline_setting.HasField("desired_retention")
            else _DEFAULT_DESIRED_RETENTION
        ),
        maximum_interval=(
            request.baseline_setting.maximum_interval
            if request.HasField("baseline_setting") and request.baseline_setting.HasField("maximum_interval")
            else _DEFAULT_MAXIMUM_INTERVAL
        ),
    )
    metrics_all = calibration_metrics.get("all", {})
    result = OptimizeResultModel(
        optimized_setting=optimized_setting,
        metrics={
            "rmse_bins": float(metrics_all.get("RMSE(bins)", 0.0)),
            "log_loss": float(metrics_all.get("LogLoss", loss_after)),
            "calibration": float(loss_before - loss_after),
            "optimal_retention": optimal_retention,
        },
        sample_count=int(len(dataset)),
        dropped_sample_count=max(0, len(request.review_logs) - len(dataset)),
        diagnostics=json.dumps(
            {
                "request_id": request.request_id,
                "dataset_id": request.dataset_id,
                "knowledge_id": request.knowledge_id,
                "loss_before": float(loss_before),
                "loss_after": float(loss_after),
            },
            ensure_ascii=False,
        ),
    )
    return kmemo_pb2.OptimizeParametersResponse(
        optimized_setting=_setting_to_proto(result.optimized_setting),
        metrics=_optimizer_metrics_to_proto(result.metrics),
        sample_count=result.sample_count,
        dropped_sample_count=result.dropped_sample_count,
        warnings=list(result.warnings),
        diagnostics=result.diagnostics,
    )


def _setting_from_proto(setting: kmemo_pb2.SchedulerSetting) -> SchedulerSettingModel:
    parameters = tuple(setting.parameters) if setting.parameters else _DEFAULT_PARAMETERS
    desired_retention = (
        setting.desired_retention if setting.HasField("desired_retention") else _DEFAULT_DESIRED_RETENTION
    )
    maximum_interval = (
        setting.maximum_interval if setting.HasField("maximum_interval") else _DEFAULT_MAXIMUM_INTERVAL
    )
    try:
        Scheduler(
            parameters=parameters,
            desired_retention=desired_retention,
            maximum_interval=maximum_interval,
        )
    except Exception as exc:
        raise ValueError(str(exc)) from exc
    return SchedulerSettingModel(
        parameters=tuple(float(value) for value in parameters),
        desired_retention=float(desired_retention),
        maximum_interval=int(maximum_interval),
    )


def _setting_to_proto(setting: SchedulerSettingModel) -> kmemo_pb2.SchedulerSetting:
    proto = kmemo_pb2.SchedulerSetting(parameters=list(setting.parameters))
    proto.desired_retention = setting.desired_retention
    proto.maximum_interval = setting.maximum_interval
    return proto


def _scheduler_from_setting(setting: SchedulerSettingModel) -> Scheduler:
    return Scheduler(
        parameters=setting.parameters,
        desired_retention=setting.desired_retention,
        maximum_interval=setting.maximum_interval,
    )


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
    if review_log.reviewed_at is not None:
        proto.reviewed_at.CopyFrom(_datetime_to_timestamp(review_log.reviewed_at))
    return proto


def _review_result_to_proto(result: ReviewResultModel) -> kmemo_pb2.ReviewCardResponse:
    response = kmemo_pb2.ReviewCardResponse(
        card=_card_state_to_proto(result.card),
        review_log=_review_log_to_proto(result.review_log),
        effective_setting=_setting_to_proto(result.effective_setting),
        warnings=list(result.warnings),
    )
    if result.retrievability is not None:
        response.retrievability = result.retrievability
    if result.diagnostics is not None:
        response.diagnostics = result.diagnostics
    return response


def _optimizer_metrics_to_proto(metrics: dict[str, float]) -> kmemo_pb2.OptimizerMetrics:
    proto = kmemo_pb2.OptimizerMetrics()
    if "rmse_bins" in metrics:
        proto.rmse_bins = metrics["rmse_bins"]
    if "log_loss" in metrics:
        proto.log_loss = metrics["log_loss"]
    if "calibration" in metrics:
        proto.calibration = metrics["calibration"]
    if "optimal_retention" in metrics and metrics["optimal_retention"] > 0:
        proto.optimal_retention = metrics["optimal_retention"]
    return proto


def _fsrs_card_from_model(card: CardStateModel, card_id: str) -> Card:
    state = _PROTO_STATE_TO_FSRS_STATE[card.state]
    step = None
    if state in {State.Learning, State.Relearning}:
        step = 0
    return Card(
        card_id=_card_identifier(card_id),
        state=state,
        step=step,
        stability=card.stability,
        difficulty=card.difficulty,
        due=card.due,
        last_review=card.last_review,
    )


def _card_state_from_fsrs(
    card: Card,
    previous: CardStateModel,
    reviewed_at: datetime,
    rating: int | None,
    *,
    reps: int | None = None,
    lapses: int | None = None,
) -> CardStateModel:
    return CardStateModel(
        state=_FSRS_STATE_TO_PROTO_STATE.get(card.state, previous.state),
        due=_ensure_utc(card.due),
        last_review=_ensure_utc(card.last_review) or reviewed_at,
        stability=float(card.stability) if card.stability is not None else None,
        difficulty=float(card.difficulty) if card.difficulty is not None else None,
        elapsed_days=_resolve_elapsed_days(previous, reviewed_at),
        scheduled_days=_resolve_scheduled_days_from_due(card.due, reviewed_at),
        reps=previous.reps + 1 if reps is None else reps,
        lapses=previous.lapses + (1 if rating == 1 else 0) if lapses is None else lapses,
        learning_steps_json=_normalize_learning_steps(previous.learning_steps_json),
    )


def _review_log_from_fsrs(
    review_log: ReviewLog,
    previous_card: CardStateModel,
    next_card: Card,
) -> ReviewLogModel:
    reviewed_at = _ensure_utc(review_log.review_datetime)
    scheduled_days = _resolve_scheduled_days_from_due(next_card.due, reviewed_at)
    return ReviewLogModel(
        rating=int(review_log.rating.value),
        review=_RATING_TO_REVIEW[int(review_log.rating.value)],
        elapsed_days=_resolve_elapsed_days(previous_card, reviewed_at),
        scheduled_days=scheduled_days,
        previous_stability=previous_card.stability,
        previous_difficulty=previous_card.difficulty,
        new_stability=float(next_card.stability) if next_card.stability is not None else None,
        new_difficulty=float(next_card.difficulty) if next_card.difficulty is not None else None,
        reviewed_at=reviewed_at,
    )


def _fsrs_review_log_from_proto(log: kmemo_pb2.ReviewLogSnapshot, card_id: str) -> ReviewLog:
    if log.rating not in _RATING_TO_FSRS:
        raise ValueError(f"invalid review_log rating: {log.rating}")
    if not log.HasField("reviewed_at"):
        raise ValueError("review_log.reviewed_at is required")
    return ReviewLog(
        card_id=_card_identifier(card_id),
        rating=_RATING_TO_FSRS[log.rating],
        review_datetime=_timestamp_to_datetime(log.reviewed_at),
        review_duration=None,
    )


def _optimizer_dataset_from_request(request: kmemo_pb2.OptimizeParametersRequest) -> pd.DataFrame:
    grouped: dict[str, list[kmemo_pb2.OptimizerReviewLogEntry]] = {}
    for entry in request.review_logs:
        if not entry.card_id:
            raise ValueError("optimizer review_logs.card_id is required")
        if not entry.HasField("reviewed_at"):
            raise ValueError("optimizer review_logs.reviewed_at is required")
        if entry.rating not in _RATING_TO_REVIEW:
            raise ValueError(f"invalid optimizer rating: {entry.rating}")
        grouped.setdefault(entry.card_id, []).append(entry)

    rows: list[dict[str, object]] = []
    for card_id, entries in grouped.items():
        ordered = sorted(entries, key=lambda item: _timestamp_to_datetime(item.reviewed_at))
        rating_history: list[str] = []
        delta_history: list[str] = []
        last_reviewed_at: datetime | None = None
        for index, entry in enumerate(ordered, start=1):
            reviewed_at = _timestamp_to_datetime(entry.reviewed_at)
            if entry.HasField("card_before") and entry.card_before.HasField("elapsed_days"):
                delta_t_value = float(entry.card_before.elapsed_days)
            elif last_reviewed_at is not None:
                delta_t_value = max(0.0, (reviewed_at - last_reviewed_at).total_seconds() / 86400)
            else:
                delta_t_value = 0.0
            delta_t = int(round(delta_t_value))
            rows.append(
                {
                    "review_time": int(reviewed_at.timestamp() * 1000),
                    "card_id": _card_identifier(card_id),
                    "review_rating": int(entry.rating),
                    "review_state": _optimizer_review_state(entry),
                    "delta_t": delta_t,
                    "i": index,
                    "t_history": ",".join(delta_history),
                    "r_history": ",".join(rating_history),
                    "last_rating": int(rating_history[-1]) if rating_history else int(entry.rating),
                    "y": 0 if entry.rating == 1 else 1,
                }
            )
            rating_history.append(str(entry.rating))
            delta_history.append(str(delta_t))
            last_reviewed_at = reviewed_at

    if not rows:
        return pd.DataFrame(columns=["review_time", "card_id", "review_rating", "review_state", "delta_t", "i", "t_history", "r_history", "last_rating", "y"])

    dataset = pd.DataFrame(rows)
    dataset = dataset[(dataset["i"] > 1) & (dataset["delta_t"] > 0)].copy()
    if dataset.empty:
        return dataset
    dataset["review_time"] = dataset["review_time"].astype(int)
    dataset["card_id"] = dataset["card_id"].astype(int)
    dataset["review_rating"] = dataset["review_rating"].astype(int)
    dataset["review_state"] = dataset["review_state"].astype(int)
    dataset["delta_t"] = dataset["delta_t"].astype(int)
    dataset["i"] = dataset["i"].astype(int)
    dataset["t_history"] = dataset["t_history"].astype(str)
    dataset["r_history"] = dataset["r_history"].astype(str)
    dataset["last_rating"] = dataset["last_rating"].astype(int)
    dataset["y"] = dataset["y"].astype(int)
    return dataset


def _optimizer_review_state(entry: kmemo_pb2.OptimizerReviewLogEntry) -> int:
    if entry.HasField("card_before") and entry.card_before.state:
        state = entry.card_before.state
    elif entry.HasField("card_after") and entry.card_after.state:
        state = entry.card_after.state
    else:
        state = "review"
    if state not in _PROTO_STATE_TO_OPTIMIZER_STATE:
        raise ValueError(f"invalid optimizer state: {state}")
    return _PROTO_STATE_TO_OPTIMIZER_STATE[state]


def _card_identifier(card_id: str) -> int:
    if not card_id:
        return 0
    if card_id.isdigit():
        return int(card_id)
    digest = hashlib.sha256(card_id.encode("utf-8")).hexdigest()
    return int(digest[:15], 16)


def _resolve_elapsed_days(current: CardStateModel, reviewed_at: datetime) -> float:
    if current.elapsed_days is not None:
        return float(current.elapsed_days)
    if current.last_review is None:
        return 0.0
    delta = reviewed_at - current.last_review
    return max(0.0, delta.total_seconds() / 86400)


def _resolve_scheduled_days(card: CardStateModel, reviewed_at: datetime) -> float:
    if card.scheduled_days is not None:
        return float(card.scheduled_days)
    return _resolve_scheduled_days_from_due(card.due, reviewed_at)


def _resolve_scheduled_days_from_due(due: datetime | None, reviewed_at: datetime) -> float:
    if due is None:
        return 0.0
    return max(0.0, (_ensure_utc(due) - reviewed_at).total_seconds() / 86400)


def _normalize_learning_steps(raw: str | None) -> str:
    if raw is None:
        return _DEFAULT_LEARNING_STEPS_JSON
    try:
        json.loads(raw)
        return raw
    except json.JSONDecodeError:
        return _DEFAULT_LEARNING_STEPS_JSON


def _timestamp_to_datetime(value: Timestamp) -> datetime:
    return value.ToDatetime().astimezone(UTC)


def _datetime_to_timestamp(value: datetime) -> Timestamp:
    ts = Timestamp()
    ts.FromDatetime(value.astimezone(UTC))
    return ts


def _ensure_utc(value: datetime | None) -> datetime | None:
    if value is None:
        return None
    return value.astimezone(UTC)
