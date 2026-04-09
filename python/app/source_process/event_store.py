from __future__ import annotations

import time
from threading import Lock

from app.source_process.models import JobEventModel


class EventStore:
    def __init__(self) -> None:
        self._items: dict[str, list[JobEventModel]] = {}
        self._sequences: dict[str, int] = {}
        self._lock = Lock()

    def append(self, job_id: str, *, stage: str, message: str) -> JobEventModel:
        with self._lock:
            sequence = self._sequences.get(job_id, 0) + 1
            self._sequences[job_id] = sequence
            event = JobEventModel(
                job_id=job_id,
                sequence=sequence,
                stage=stage,
                message=message,
                created_at_unix=int(time.time()),
            )
            self._items.setdefault(job_id, []).append(event)
            return event

    def list(self, job_id: str, after_sequence: int | None = None) -> list[JobEventModel]:
        items = list(self._items.get(job_id, []))
        if after_sequence is None:
            return items
        return [event for event in items if event.sequence > after_sequence]
