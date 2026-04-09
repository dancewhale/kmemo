from __future__ import annotations

from dataclasses import replace
from threading import Lock

from app.source_process.errors import JobNotFoundError
from app.source_process.models import JobStatusModel


class JobStore:
    def __init__(self) -> None:
        self._items: dict[str, JobStatusModel] = {}
        self._lock = Lock()

    def create(self, job: JobStatusModel) -> JobStatusModel:
        with self._lock:
            self._items[job.job_id] = replace(job)
            return replace(job)

    def get(self, job_id: str) -> JobStatusModel:
        with self._lock:
            job = self._items.get(job_id)
            if job is None:
                raise JobNotFoundError(f"job {job_id!r} not found")
            return replace(job)

    def update(self, job_id: str, **changes) -> JobStatusModel:
        with self._lock:
            current = self._items.get(job_id)
            if current is None:
                raise JobNotFoundError(f"job {job_id!r} not found")
            updated = replace(current, **changes)
            self._items[job_id] = updated
            return replace(updated)
