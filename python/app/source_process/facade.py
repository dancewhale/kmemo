from __future__ import annotations

from concurrent.futures import ThreadPoolExecutor

from app import request_context
from app.logging_setup import get_logger
from app.source_process.cancellation import CancellationRegistry
from app.source_process.event_store import EventStore
from app.source_process.errors import JobCancelledError, SourceProcessError, ValidationError
from app.source_process.job_store import JobStore
from app.source_process.manifest import ManifestWriter
from app.source_process.models import CapabilitiesModel, JobContext, JobStatusModel, SubmitJobModel
from app.source_process.pipeline import SourceProcessPipeline
from app.source_process.registry import ConverterRegistry
from app.source_process.workspace import WorkspaceManager


class SourceProcessFacade:
    def __init__(self) -> None:
        self._logger = get_logger("kmemo.worker.source_process")
        self._registry = ConverterRegistry()
        self._workspace = WorkspaceManager()
        self._pipeline = SourceProcessPipeline()
        self._manifest_writer = ManifestWriter()
        self._jobs = JobStore()
        self._events = EventStore()
        self._cancellation = CancellationRegistry()
        self._executor = ThreadPoolExecutor(max_workers=2, thread_name_prefix="source-process")

    def submit_job(self, req: SubmitJobModel) -> JobStatusModel:
        self._validate_submit(req)
        self._jobs.create(JobStatusModel(job_id=req.job_id, status="queued", stage="accepted", progress=0.0))
        self._events.append(req.job_id, stage="accepted", message="job accepted")
        self._executor.submit(self._run_job, req)
        return self._jobs.get(req.job_id)

    def get_job(self, job_id: str) -> JobStatusModel:
        return self._jobs.get(job_id)

    def list_job_events(self, job_id: str, after_sequence: int | None = None):
        self._jobs.get(job_id)
        return self._events.list(job_id, after_sequence=after_sequence)

    def cancel_job(self, job_id: str) -> JobStatusModel:
        job = self._jobs.get(job_id)
        if job.status in {"succeeded", "failed", "cancelled"}:
            return job
        self._cancellation.request_cancel(job_id)
        self._events.append(job_id, stage=job.stage or "cancel", message="cancel requested")
        return self._jobs.update(
            job_id,
            status="cancel_requested",
            stage=job.stage or "cancel",
            progress=job.progress,
        )

    def get_capabilities(self) -> CapabilitiesModel:
        return self._registry.capabilities()

    def _run_job(self, req: SubmitJobModel) -> None:
        logger = self._logger.bind(
            service="python-worker",
            component="source_process",
            job_id=req.job_id,
            source_type=req.source_type,
            conversion_mode=req.options.conversion_mode,
            request_id=request_context.get_request_id(),
        )
        try:
            self._jobs.update(req.job_id, status="running", stage="prepare_workspace", progress=0.1)
            self._events.append(req.job_id, stage="prepare_workspace", message="job started")
            paths = self._workspace.prepare(
                workspace_dir=req.workspace_dir,
                output_dir=req.output_dir,
                temp_dir=req.temp_dir,
            )

            spec = self._registry.resolve(req.source_type, req.options.conversion_mode, req.options.fallback_modes)
            self._jobs.update(req.job_id, stage="select_converter", progress=0.25)
            self._events.append(req.job_id, stage="select_converter", message=f"selected {spec.converter_name}")

            ctx = JobContext(
                job_id=req.job_id,
                request_id=request_context.get_request_id(),
                workspace_dir=paths.workspace_dir,
                output_dir=paths.output_dir,
                temp_dir=paths.temp_dir,
                source_type=req.source_type,
                conversion_mode=spec.conversion_mode,
                logger=logger,
                is_cancel_requested=lambda: self._cancellation.is_cancel_requested(req.job_id),
            )

            self._jobs.update(req.job_id, stage="convert_source", progress=0.5)
            result, manifest = self._pipeline.run(ctx, req, paths)
            self._events.append(req.job_id, stage="convert_source", message="conversion finished")

            self._jobs.update(req.job_id, stage="write_manifest", progress=0.8)
            manifest_path = self._manifest_writer.write(paths.output_dir, manifest)
            self._events.append(req.job_id, stage="write_manifest", message="manifest written")

            completed = self._jobs.update(
                req.job_id,
                status="succeeded",
                stage="finalize",
                progress=1.0,
                result_path=str(manifest_path),
                result=result.__class__(
                    entry_html_path=result.entry_html_path,
                    cleaned_html_path=result.cleaned_html_path,
                    raw_text_path=result.raw_text_path,
                    assets=result.assets,
                    extracted_metadata=result.extracted_metadata,
                    manifest_path=str(manifest_path),
                    content_hash=result.content_hash,
                    effective_conversion_mode=result.effective_conversion_mode,
                    converter_name=result.converter_name,
                    converter_version=result.converter_version,
                    cleaner_version=result.cleaner_version,
                ),
            )
            self._events.append(req.job_id, stage="finalize", message="job succeeded")
            logger.info("job succeeded", stage=completed.stage)
        except JobCancelledError as exc:
            cancelled = self._jobs.update(
                req.job_id,
                status="cancelled",
                stage="cancelled",
                progress=1.0,
                error_code="job_cancelled",
                error_message=str(exc),
            )
            self._events.append(req.job_id, stage="cancelled", message="job cancelled")
            logger.warning("job cancelled", stage=cancelled.stage, error=str(exc))
        except SourceProcessError as exc:
            failed = self._jobs.update(
                req.job_id,
                status="failed",
                stage="failed",
                progress=1.0,
                error_code=exc.__class__.__name__.lower(),
                error_message=str(exc),
            )
            self._events.append(req.job_id, stage="failed", message=str(exc))
            logger.error("job failed", stage=failed.stage, error=str(exc))
        except Exception as exc:  # pragma: no cover - defensive catch for worker runtime
            failed = self._jobs.update(
                req.job_id,
                status="failed",
                stage="failed",
                progress=1.0,
                error_code="internal_error",
                error_message=str(exc),
            )
            self._events.append(req.job_id, stage="failed", message=str(exc))
            logger.exception("job failed", stage=failed.stage, error=str(exc))
        finally:
            self._cancellation.clear(req.job_id)

    def _validate_submit(self, req: SubmitJobModel) -> None:
        if not req.job_id.strip():
            raise ValidationError("job_id is required")
        if not req.source_type.strip():
            raise ValidationError("source_type is required")
        if req.raw_html is None and req.source_path is None:
            raise ValidationError("either raw_html or source_path is required")


_FACADE: SourceProcessFacade | None = None


def get_facade() -> SourceProcessFacade:
    global _FACADE
    if _FACADE is None:
        _FACADE = SourceProcessFacade()
    return _FACADE
