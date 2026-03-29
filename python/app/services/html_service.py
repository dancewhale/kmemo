"""HTML cleanup (skeleton only)."""

from __future__ import annotations

from kmemo.v1 import kmemo_pb2


def clean_html(request: kmemo_pb2.CleanHtmlRequest) -> kmemo_pb2.CleanHtmlResponse:
    """TODO: sanitize / normalize HTML via dedicated pipeline."""
    return kmemo_pb2.CleanHtmlResponse(
        ok=False,
        message="TODO: HTML cleaning not implemented",
        cleaned_html=request.html,
    )
