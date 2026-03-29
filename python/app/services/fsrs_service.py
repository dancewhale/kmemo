"""FSRS scheduling (skeleton only)."""

from __future__ import annotations

from kmemo.v1 import kmemo_pb2


def calculate_fsrs(request: kmemo_pb2.CalculateFsrsRequest) -> kmemo_pb2.CalculateFsrsResponse:
    """TODO: integrate FSRS once payloads are defined."""
    return kmemo_pb2.CalculateFsrsResponse(
        ok=False,
        message="TODO: FSRS not implemented",
        result_json=b"{}",
    )
