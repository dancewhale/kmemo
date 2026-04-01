from __future__ import annotations

import time
from collections.abc import Callable

from app import request_context
from app.logging_setup import get_logger

_REQUEST_ID_HEADER = "x-request-id"


def extract_request_id(context) -> str | None:
    for item in context.invocation_metadata():
        if item.key == _REQUEST_ID_HEADER and item.value:
            return item.value
    return None


def with_rpc_logging(context, rpc_method: str, handler: Callable[[], object]):
    request_id = extract_request_id(context)
    token = request_context.set_request_id(request_id)
    logger = get_logger("kmemo.worker").bind(
        service="python-worker",
        component="grpc_server",
        grpc_method=rpc_method,
        request_id=request_id,
    )

    started_at = time.perf_counter()
    logger.info("rpc started")
    try:
        result = handler()
        duration_ms = round((time.perf_counter() - started_at) * 1000, 2)
        logger.info("rpc finished", duration_ms=duration_ms)
        return result
    except Exception as exc:
        duration_ms = round((time.perf_counter() - started_at) * 1000, 2)
        logger.exception("rpc failed", duration_ms=duration_ms, error=str(exc))
        raise
    finally:
        request_context.clear_request_id(token)
