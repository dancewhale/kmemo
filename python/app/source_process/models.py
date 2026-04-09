from __future__ import annotations

from dataclasses import dataclass, field
from pathlib import Path
from typing import Any


@dataclass(frozen=True)
class ImportOptionsModel:
    conversion_mode: str
    fallback_modes: tuple[str, ...] = ()
    extract_main_content: bool = True
    sanitize_html: bool = True
    preserve_semantic_tags: bool = True
    download_remote_assets: bool = False
    inline_small_images: bool = False
    generate_toc: bool = False
    analyze_structure: bool = True
    keep_source_copy: bool = False
    enabled_cleaners: tuple[str, ...] = ()
    converter_params_json: str | None = None


@dataclass(frozen=True)
class SubmitJobModel:
    job_id: str
    source_type: str
    source_path: str | None
    source_uri: str | None
    source_url: str | None
    raw_html: str | None
    workspace_dir: str
    output_dir: str
    temp_dir: str
    metadata: dict[str, str]
    options: ImportOptionsModel
    idempotency_key: str | None = None


@dataclass(frozen=True)
class JobContext:
    job_id: str
    request_id: str | None
    workspace_dir: Path
    output_dir: Path
    temp_dir: Path
    source_type: str
    conversion_mode: str
    logger: Any
    is_cancel_requested: Any | None = None


@dataclass(frozen=True)
class ConversionOutput:
    article_html: tuple[str, ...]
    cleaned_html: str
    raw_text: str | None
    asset_paths: tuple[str, ...]
    extracted_metadata: dict[str, str]
    effective_conversion_mode: str
    converter_name: str
    converter_version: str
    cleaner_version: str | None = None


@dataclass(frozen=True)
class ImportResultModel:
    entry_html_path: str
    cleaned_html_path: str
    raw_text_path: str | None
    assets: tuple[str, ...]
    extracted_metadata: dict[str, str]
    manifest_path: str
    content_hash: str
    effective_conversion_mode: str
    converter_name: str
    converter_version: str
    cleaner_version: str | None = None


@dataclass(frozen=True)
class JobEventModel:
    job_id: str
    sequence: int
    stage: str
    message: str
    created_at_unix: int


@dataclass
class JobStatusModel:
    job_id: str
    status: str
    stage: str
    progress: float
    result_path: str | None = None
    error_code: str | None = None
    error_message: str | None = None
    result: ImportResultModel | None = None


@dataclass(frozen=True)
class ConverterCapabilityModel:
    source_type: str
    conversion_mode: str
    converter_name: str
    converter_version: str
    description: str


@dataclass(frozen=True)
class CleanerCapabilityModel:
    cleaner_name: str
    cleaner_version: str
    description: str


@dataclass(frozen=True)
class CapabilitiesModel:
    source_types: tuple[str, ...]
    conversion_modes: tuple[str, ...]
    converters: tuple[ConverterCapabilityModel, ...]
    cleaners: tuple[CleanerCapabilityModel, ...]


@dataclass(frozen=True)
class ManifestModel:
    job_id: str
    source_type: str
    entry_html_path: str
    cleaned_html_path: str
    raw_text_path: str | None
    assets: tuple[str, ...]
    metadata: dict[str, str]
    structure: dict[str, Any] = field(default_factory=dict)
    content_hash: str = ""
    effective_conversion_mode: str = ""
    converter_name: str = ""
    converter_version: str = ""
    cleaner_version: str | None = None
