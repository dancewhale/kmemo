from __future__ import annotations

from pathlib import Path
import sys
import time

_ROOT = Path(__file__).resolve().parent.parent
_GEN_ROOT = _ROOT / "generated"
if str(_GEN_ROOT) not in sys.path:
    sys.path.insert(0, str(_GEN_ROOT))
if str(_ROOT) not in sys.path:
    sys.path.insert(0, str(_ROOT))

from kmemo.v1 import kmemo_pb2

from app.services import import_service


def _submit_request(tmp_path: Path, *, job_id: str = "job-1", raw_html: str = "<html><head><title>Hello</title></head><body><h1>Demo</h1><p>Test</p></body></html>") -> kmemo_pb2.SubmitImportJobRequest:
    workspace_dir = tmp_path / job_id / "workspace"
    output_dir = tmp_path / job_id / "output"
    temp_dir = tmp_path / job_id / "temp"
    return kmemo_pb2.SubmitImportJobRequest(
        job_id=job_id,
        source_type="html",
        raw_html=raw_html,
        workspace_dir=str(workspace_dir),
        output_dir=str(output_dir),
        temp_dir=str(temp_dir),
        options=kmemo_pb2.ImportOptions(conversion_mode="HTML_RAW"),
        metadata={"knowledge_id": "knowledge-1"},
    )


def _wait_for_terminal(job_id: str, timeout_s: float = 5.0):
    deadline = time.time() + timeout_s
    while time.time() < deadline:
        job = import_service.get_job(kmemo_pb2.GetJobRequest(job_id=job_id)).job
        if job.status in {"succeeded", "failed", "cancelled"}:
            return job
        time.sleep(0.05)
    raise AssertionError(f"job {job_id} did not reach terminal state in time")


def test_submit_import_job_accepts_valid_request(tmp_path: Path):
    response = import_service.submit_import_job(_submit_request(tmp_path, job_id="job-submit"))

    assert response.job_id == "job-submit"
    assert response.accepted is True
    assert response.status in {"queued", "running", "succeeded"}



def test_submit_import_job_runs_to_success_and_writes_outputs(tmp_path: Path):
    job_id = "job-success"
    import_service.submit_import_job(_submit_request(tmp_path, job_id=job_id))

    job = _wait_for_terminal(job_id)

    assert job.status == "succeeded"
    assert job.result.manifest_path.endswith("manifest.json")
    assert Path(job.result.entry_html_path).exists()
    assert Path(job.result.cleaned_html_path).exists()
    assert Path(job.result.raw_text_path).exists()
    assert Path(job.result.manifest_path).exists()
    assert job.result.extracted_metadata["title"] == "Hello"
    assert job.result.effective_conversion_mode == "HTML_RAW"



def test_get_job_returns_current_status(tmp_path: Path):
    job_id = "job-status"
    import_service.submit_import_job(_submit_request(tmp_path, job_id=job_id))

    response = import_service.get_job(kmemo_pb2.GetJobRequest(job_id=job_id))

    assert response.job.job_id == job_id
    assert response.job.status in {"queued", "running", "succeeded"}



def test_list_job_events_returns_ordered_events(tmp_path: Path):
    job_id = "job-events"
    import_service.submit_import_job(_submit_request(tmp_path, job_id=job_id))
    _wait_for_terminal(job_id)

    response = import_service.list_job_events(kmemo_pb2.ListJobEventsRequest(job_id=job_id))

    assert response.events
    assert response.events[0].sequence == 1
    assert sorted(event.sequence for event in response.events) == [event.sequence for event in response.events]



def test_cancel_job_transitions_status(tmp_path: Path):
    job_id = "job-cancel"
    import_service.submit_import_job(_submit_request(tmp_path, job_id=job_id))

    cancel = import_service.cancel_job(kmemo_pb2.CancelJobRequest(job_id=job_id))
    assert cancel.job_id == job_id
    assert cancel.status in {"cancel_requested", "cancelled", "succeeded"}



def test_get_capabilities_returns_source_process_surface():
    response = import_service.get_capabilities(kmemo_pb2.GetCapabilitiesRequest())

    assert "html" in response.source_types
    assert "HTML_RAW" in response.conversion_modes
    assert response.converters
    assert response.cleaners



def test_submit_import_job_rejects_invalid_request(tmp_path: Path):
    request = _submit_request(tmp_path, job_id="job-invalid")
    request.ClearField("raw_html")

    try:
        import_service.submit_import_job(request)
    except Exception as exc:
        assert "required" in str(exc)
    else:  # pragma: no cover
        raise AssertionError("expected invalid request to raise")
