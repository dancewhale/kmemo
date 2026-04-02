from __future__ import annotations

from datetime import UTC, datetime
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


def _timestamp(value: datetime) -> Timestamp:
    ts = Timestamp()
    ts.FromDatetime(value.astimezone(UTC))
    return ts


def _schedule_request(rating: int = 3) -> kmemo_pb2.SchedulerScheduleRequest:
    reviewed_at = datetime(2026, 4, 2, 12, 0, tzinfo=UTC)
    return kmemo_pb2.SchedulerScheduleRequest(
        card_id="card-1",
        current=kmemo_pb2.CardState(
            state="review",
            reps=3,
            lapses=0,
        ),
        review=kmemo_pb2.ReviewInput(
            rating=rating,
            reviewed_at=_timestamp(reviewed_at),
        ),
    )


def test_scheduler_schedule_works_without_setting():
    response = fsrs_service.scheduler_schedule(_schedule_request())

    assert response.ok is True
    assert response.next.state == "review"
    assert response.effective_setting.parameters
    assert response.effective_setting.desired_retention > 0
    assert response.effective_setting.maximum_interval > 0


def test_scheduler_setting_overrides_previous_values():
    first = fsrs_service.scheduler_set_setting(
        kmemo_pb2.SchedulerSetSettingRequest(
            setting=kmemo_pb2.SchedulerSetting(
                parameters=[0.5] * 19,
                desired_retention=0.83,
                maximum_interval=100,
            )
        )
    )
    second = fsrs_service.scheduler_set_setting(
        kmemo_pb2.SchedulerSetSettingRequest(
            setting=kmemo_pb2.SchedulerSetting(
                parameters=[2.0] * 19,
                desired_retention=0.97,
                maximum_interval=200,
            )
        )
    )
    scheduled = fsrs_service.scheduler_schedule(_schedule_request())

    assert first.ok is True
    assert second.ok is True
    assert scheduled.effective_setting.parameters == [2.0] * 19
    assert scheduled.effective_setting.desired_retention == 0.97
    assert scheduled.effective_setting.maximum_interval == 200
