"""gRPC server entrypoint for the Python worker."""

from __future__ import annotations

import logging
import sys
from concurrent import futures
from pathlib import Path

# Generated protos live in ../generated (created by `task proto`).
_GEN_ROOT = Path(__file__).resolve().parent.parent / "generated"
if _GEN_ROOT.is_dir():
    sys.path.insert(0, str(_GEN_ROOT))

import grpc

from kmemo.v1 import kmemo_pb2_grpc

from app.grpc_logging import with_rpc_logging
from app.logging_setup import configure_logging, get_logger
from app.services import fsrs_service, html_service, import_service

_LOG = get_logger("kmemo.worker")


class KmemoProcessor(kmemo_pb2_grpc.KmemoProcessorServicer):
    """Thin dispatcher; logic stays in app.services.* (placeholders for now)."""

    def CalculateFsrs(self, request, context):
        return with_rpc_logging(
            context,
            "CalculateFsrs",
            lambda: fsrs_service.calculate_fsrs(request),
        )

    def CleanHtml(self, request, context):
        return with_rpc_logging(
            context,
            "CleanHtml",
            lambda: html_service.clean_html(request),
        )

    def PrepareImportMaterial(self, request, context):
        return with_rpc_logging(
            context,
            "PrepareImportMaterial",
            lambda: import_service.prepare_import(request),
        )


def serve(address: str = "[::]:50051") -> None:
    configure_logging(logging.INFO)
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    kmemo_pb2_grpc.add_KmemoProcessorServicer_to_server(KmemoProcessor(), server)
    server.add_insecure_port(address)
    server.start()
    _LOG.info(
        "worker started",
        service="python-worker",
        component="grpc_server",
        target=address,
    )
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
