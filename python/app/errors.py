"""Map domain failures to gRPC status codes (skeleton)."""

from __future__ import annotations

import grpc


def to_status(exc: Exception) -> grpc.StatusCode:
    """TODO: classify validation vs internal errors."""
    _ = exc
    return grpc.StatusCode.INTERNAL
