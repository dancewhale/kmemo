from __future__ import annotations

import json
from pathlib import Path

from app.source_process.errors import OutputWriteError
from app.source_process.models import ManifestModel


class ManifestWriter:
    def write(self, output_dir: Path, manifest: ManifestModel) -> Path:
        payload = {
            "job_id": manifest.job_id,
            "source_type": manifest.source_type,
            "entry_html_path": manifest.entry_html_path,
            "cleaned_html_path": manifest.cleaned_html_path,
            "raw_text_path": manifest.raw_text_path,
            "assets": list(manifest.assets),
            "metadata": manifest.metadata,
            "structure": manifest.structure,
            "content_hash": manifest.content_hash,
            "effective_conversion_mode": manifest.effective_conversion_mode,
            "converter_name": manifest.converter_name,
            "converter_version": manifest.converter_version,
            "cleaner_version": manifest.cleaner_version,
        }
        manifest_path = output_dir / "manifest.json"
        try:
            manifest_path.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")
        except OSError as exc:
            raise OutputWriteError(f"failed to write manifest: {exc}") from exc
        return manifest_path
