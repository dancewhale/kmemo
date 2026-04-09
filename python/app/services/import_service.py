"""Job-based source-process service facade."""

from __future__ import annotations

from kmemo.v1 import kmemo_pb2

from app.source_process import get_facade
from app.source_process.models import ImportOptionsModel, SubmitJobModel


def submit_import_job(
    request: kmemo_pb2.SubmitImportJobRequest,
) -> kmemo_pb2.SubmitImportJobResponse:
    job = get_facade().submit_job(_submit_job_from_proto(request))
    return kmemo_pb2.SubmitImportJobResponse(
        job_id=job.job_id,
        status=job.status,
        accepted=job.status in {"queued", "running", "succeeded"},
    )


def get_job(request: kmemo_pb2.GetJobRequest) -> kmemo_pb2.GetJobResponse:
    job = get_facade().get_job(request.job_id)
    return kmemo_pb2.GetJobResponse(job=_job_to_proto(job))


def list_job_events(
    request: kmemo_pb2.ListJobEventsRequest,
) -> kmemo_pb2.ListJobEventsResponse:
    after_sequence = request.after_sequence if request.HasField("after_sequence") else None
    events = get_facade().list_job_events(request.job_id, after_sequence=after_sequence)
    return kmemo_pb2.ListJobEventsResponse(events=[_job_event_to_proto(event) for event in events])


def cancel_job(request: kmemo_pb2.CancelJobRequest) -> kmemo_pb2.CancelJobResponse:
    job = get_facade().cancel_job(request.job_id)
    return kmemo_pb2.CancelJobResponse(job_id=job.job_id, status=job.status)


def get_capabilities(
    request: kmemo_pb2.GetCapabilitiesRequest,
) -> kmemo_pb2.GetCapabilitiesResponse:
    _ = request
    caps = get_facade().get_capabilities()
    return kmemo_pb2.GetCapabilitiesResponse(
        source_types=list(caps.source_types),
        conversion_modes=list(caps.conversion_modes),
        converters=[
            kmemo_pb2.ConverterCapability(
                source_type=item.source_type,
                conversion_mode=item.conversion_mode,
                converter_name=item.converter_name,
                converter_version=item.converter_version,
                description=item.description,
            )
            for item in caps.converters
        ],
        cleaners=[
            kmemo_pb2.CleanerCapability(
                cleaner_name=item.cleaner_name,
                cleaner_version=item.cleaner_version,
                description=item.description,
            )
            for item in caps.cleaners
        ],
    )


def _submit_job_from_proto(request: kmemo_pb2.SubmitImportJobRequest) -> SubmitJobModel:
    options = request.options if request.HasField("options") else kmemo_pb2.ImportOptions()
    return SubmitJobModel(
        job_id=request.job_id,
        source_type=request.source_type,
        source_path=request.source_path if request.HasField("source_path") else None,
        source_uri=request.source_uri if request.HasField("source_uri") else None,
        source_url=request.source_url if request.HasField("source_url") else None,
        raw_html=request.raw_html if request.HasField("raw_html") else None,
        workspace_dir=request.workspace_dir,
        output_dir=request.output_dir,
        temp_dir=request.temp_dir,
        metadata=dict(request.metadata),
        options=ImportOptionsModel(
            conversion_mode=options.conversion_mode or "AUTO",
            fallback_modes=tuple(options.fallback_modes),
            extract_main_content=options.extract_main_content,
            sanitize_html=options.sanitize_html,
            preserve_semantic_tags=options.preserve_semantic_tags,
            download_remote_assets=options.download_remote_assets,
            inline_small_images=options.inline_small_images,
            generate_toc=options.generate_toc,
            analyze_structure=options.analyze_structure,
            keep_source_copy=options.keep_source_copy,
            enabled_cleaners=tuple(options.enabled_cleaners),
            converter_params_json=options.converter_params_json if options.HasField("converter_params_json") else None,
        ),
        idempotency_key=request.idempotency_key if request.HasField("idempotency_key") else None,
    )


def _job_to_proto(job) -> kmemo_pb2.SourceProcessJob:
    payload = kmemo_pb2.SourceProcessJob(
        job_id=job.job_id,
        status=job.status,
        stage=job.stage,
        progress=job.progress,
    )
    if job.result_path is not None:
        payload.result_path = job.result_path
    if job.error_code is not None:
        payload.error_code = job.error_code
    if job.error_message is not None:
        payload.error_message = job.error_message
    if job.result is not None:
        payload.result.CopyFrom(
            kmemo_pb2.ImportResult(
                entry_html_path=job.result.entry_html_path,
                cleaned_html_path=job.result.cleaned_html_path,
                raw_text_path=job.result.raw_text_path or "",
                assets=list(job.result.assets),
                extracted_metadata=job.result.extracted_metadata,
                manifest_path=job.result.manifest_path,
                content_hash=job.result.content_hash,
                effective_conversion_mode=job.result.effective_conversion_mode,
                converter_name=job.result.converter_name,
                converter_version=job.result.converter_version,
                cleaner_version=job.result.cleaner_version or "",
            )
        )
    return payload


def _job_event_to_proto(event) -> kmemo_pb2.SourceProcessJobEvent:
    return kmemo_pb2.SourceProcessJobEvent(
        job_id=event.job_id,
        sequence=event.sequence,
        stage=event.stage,
        message=event.message,
        created_at_unix=event.created_at_unix,
    )
