from __future__ import annotations

import hashlib
import re
from pathlib import Path

from app.source_process.errors import ConvertError, JobCancelledError, ValidationError
from app.source_process.models import ConversionOutput, ImportResultModel, JobContext, ManifestModel, SubmitJobModel
from app.source_process.workspace import WorkspacePaths

_TITLE_RE = re.compile(r"<title>(.*?)</title>", re.IGNORECASE | re.DOTALL)


class SourceProcessPipeline:
    def run(self, ctx: JobContext, req: SubmitJobModel, paths: WorkspacePaths) -> ImportResultModel:
        html = self._load_html(req)
        self._check_cancel(ctx)

        cleaned_html = html
        raw_text = _extract_raw_text(cleaned_html)
        article_path = paths.article_dir / "1.html"
        cleaned_path = paths.output_dir / "cleaned.html"
        raw_text_path = paths.output_dir / "raw_text.txt"

        article_path.write_text(cleaned_html, encoding="utf-8")
        cleaned_path.write_text(cleaned_html, encoding="utf-8")
        raw_text_path.write_text(raw_text, encoding="utf-8")

        metadata = dict(req.metadata)
        title = _extract_title(cleaned_html)
        if title and "title" not in metadata:
            metadata["title"] = title

        content_hash = hashlib.sha256(cleaned_html.encode("utf-8")).hexdigest()
        manifest = ManifestModel(
            job_id=req.job_id,
            source_type=req.source_type,
            entry_html_path=str(article_path),
            cleaned_html_path=str(cleaned_path),
            raw_text_path=str(raw_text_path),
            assets=(),
            metadata=metadata,
            structure={"headings": [], "blocks": []},
            content_hash=content_hash,
            effective_conversion_mode=ctx.conversion_mode,
            converter_name="html_raw",
            converter_version="v1",
            cleaner_version="identity_cleaner/v1",
        )
        return ImportResultModel(
            entry_html_path=manifest.entry_html_path,
            cleaned_html_path=manifest.cleaned_html_path,
            raw_text_path=manifest.raw_text_path,
            assets=manifest.assets,
            extracted_metadata=metadata,
            manifest_path="",
            content_hash=content_hash,
            effective_conversion_mode=ctx.conversion_mode,
            converter_name="html_raw",
            converter_version="v1",
            cleaner_version="identity_cleaner/v1",
        ), manifest

    def _load_html(self, req: SubmitJobModel) -> str:
        if req.source_type.lower() != "html":
            raise ConvertError(f"unsupported source_type {req.source_type!r} in initial implementation")
        if req.raw_html:
            return req.raw_html
        if req.source_path:
            return Path(req.source_path).read_text(encoding="utf-8")
        raise ValidationError("either raw_html or source_path is required for html jobs")

    def _check_cancel(self, ctx: JobContext) -> None:
        if callable(ctx.is_cancel_requested) and ctx.is_cancel_requested():
            raise JobCancelledError(f"job {ctx.job_id!r} cancelled")


def _extract_title(html: str) -> str | None:
    match = _TITLE_RE.search(html)
    if match is None:
        return None
    title = match.group(1).strip()
    return title or None


def _extract_raw_text(html: str) -> str:
    text = re.sub(r"<[^>]+>", " ", html)
    text = re.sub(r"\s+", " ", text).strip()
    return text
