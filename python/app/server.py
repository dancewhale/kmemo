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

    def SchedulerSetSetting(self, request, context):
        return with_rpc_logging(
            context,
            "SchedulerSetSetting",
            lambda: fsrs_service.scheduler_set_setting(request),
        )

    def GetCardRetrievability(self, request, context):
        return with_rpc_logging(
            context,
            "GetCardRetrievability",
            lambda: fsrs_service.get_card_retrievability(request),
        )

    def ReviewCard(self, request, context):
        return with_rpc_logging(
            context,
            "ReviewCard",
            lambda: fsrs_service.review_card(request),
        )

    def RescheduleCard(self, request, context):
        return with_rpc_logging(
            context,
            "RescheduleCard",
            lambda: fsrs_service.reschedule_card(request),
        )

    def OptimizeParameters(self, request, context):
        return with_rpc_logging(
            context,
            "OptimizeParameters",
            lambda: fsrs_service.optimize_parameters(request),
        )

    def CleanHtml(self, request, context):
        return with_rpc_logging(
            context,
            "CleanHtml",
            lambda: html_service.clean_html(request),
        )

    def SubmitImportJob(self, request, context):
        return with_rpc_logging(
            context,
            "SubmitImportJob",
            lambda: import_service.submit_import_job(request),
        )

    def GetJob(self, request, context):
        return with_rpc_logging(
            context,
            "GetJob",
            lambda: import_service.get_job(request),
        )

    def ListJobEvents(self, request, context):
        return with_rpc_logging(
            context,
            "ListJobEvents",
            lambda: import_service.list_job_events(request),
        )

    def CancelJob(self, request, context):
        return with_rpc_logging(
            context,
            "CancelJob",
            lambda: import_service.cancel_job(request),
        )

    def GetCapabilities(self, request, context):
        return with_rpc_logging(
            context,
            "GetCapabilities",
            lambda: import_service.get_capabilities(request),
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
