from __future__ import annotations

from dataclasses import dataclass
from pathlib import Path

from app.source_process.errors import OutputWriteError, ValidationError


@dataclass(frozen=True)
class WorkspacePaths:
    workspace_dir: Path
    output_dir: Path
    temp_dir: Path
    article_dir: Path
    assets_dir: Path


class WorkspaceManager:
    def prepare(self, *, workspace_dir: str, output_dir: str, temp_dir: str) -> WorkspacePaths:
        if not workspace_dir.strip():
            raise ValidationError("workspace_dir is required")
        if not output_dir.strip():
            raise ValidationError("output_dir is required")
        if not temp_dir.strip():
            raise ValidationError("temp_dir is required")

        workspace_path = Path(workspace_dir)
        output_path = Path(output_dir)
        temp_path = Path(temp_dir)
        article_path = output_path / "article"
        assets_path = output_path / "assets"

        try:
            workspace_path.mkdir(parents=True, exist_ok=True)
            output_path.mkdir(parents=True, exist_ok=True)
            temp_path.mkdir(parents=True, exist_ok=True)
            article_path.mkdir(parents=True, exist_ok=True)
            assets_path.mkdir(parents=True, exist_ok=True)
        except OSError as exc:
            raise OutputWriteError(f"failed to prepare workspace: {exc}") from exc

        return WorkspacePaths(
            workspace_dir=workspace_path,
            output_dir=output_path,
            temp_dir=temp_path,
            article_dir=article_path,
            assets_dir=assets_path,
        )
