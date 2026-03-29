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

from app.services import fsrs_service, html_service, import_service

_LOG = logging.getLogger("kmemo.worker")


class KmemoProcessor(kmemo_pb2_grpc.KmemoProcessorServicer):
    """Thin dispatcher; logic stays in app.services.* (placeholders for now)."""

    def CalculateFsrs(self, request, context):
        _ = context
        return fsrs_service.calculate_fsrs(request)

    def CleanHtml(self, request, context):
        _ = context
        return html_service.clean_html(request)

    def PrepareImportMaterial(self, request, context):
        _ = context
        return import_service.prepare_import(request)


def serve(address: str = "[::]:50051") -> None:
    logging.basicConfig(level=logging.INFO)
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    kmemo_pb2_grpc.add_KmemoProcessorServicer_to_server(KmemoProcessor(), server)
    server.add_insecure_port(address)
    server.start()
    _LOG.info("kmemo python worker listening on %s", address)
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
