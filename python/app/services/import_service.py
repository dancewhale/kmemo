"""Import preprocessing (skeleton only)."""

from __future__ import annotations

from kmemo.v1 import kmemo_pb2


def prepare_import(
    request: kmemo_pb2.PrepareImportMaterialRequest,
) -> kmemo_pb2.PrepareImportMaterialResponse:
    """TODO: normalize imported sources (binary → structured descriptor)."""
    return kmemo_pb2.PrepareImportMaterialResponse(
        ok=False,
        message="TODO: import preparation not implemented",
        prepared_json=b"{}",
    )
