from __future__ import annotations

from threading import Lock


class CancellationRegistry:
    def __init__(self) -> None:
        self._cancelled: set[str] = set()
        self._lock = Lock()

    def request_cancel(self, job_id: str) -> None:
        with self._lock:
            self._cancelled.add(job_id)

    def is_cancel_requested(self, job_id: str) -> bool:
        with self._lock:
            return job_id in self._cancelled

    def clear(self, job_id: str) -> None:
        with self._lock:
            self._cancelled.discard(job_id)
