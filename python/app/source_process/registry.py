from __future__ import annotations

from dataclasses import dataclass

from app.source_process.errors import UnsupportedModeError
from app.source_process.models import CapabilitiesModel, CleanerCapabilityModel, ConverterCapabilityModel


@dataclass(frozen=True)
class ConverterSpec:
    source_type: str
    conversion_mode: str
    converter_name: str
    converter_version: str
    description: str


class ConverterRegistry:
    def __init__(self) -> None:
        self._specs: dict[tuple[str, str], ConverterSpec] = {
            ("html", "AUTO"): ConverterSpec(
                source_type="html",
                conversion_mode="AUTO",
                converter_name="html_raw",
                converter_version="v1",
                description="Resolve HTML AUTO requests to the built-in raw HTML converter.",
            ),
            ("html", "HTML_RAW"): ConverterSpec(
                source_type="html",
                conversion_mode="HTML_RAW",
                converter_name="html_raw",
                converter_version="v1",
                description="Treat HTML input as the source of truth and write normalized output files.",
            ),
        }

    def resolve(self, source_type: str, conversion_mode: str, fallback_modes: tuple[str, ...]) -> ConverterSpec:
        normalized_source_type = source_type.lower().strip()
        modes = (conversion_mode, *fallback_modes)
        for candidate in modes:
            normalized_mode = candidate.strip() or "AUTO"
            spec = self._specs.get((normalized_source_type, normalized_mode))
            if spec is not None:
                return spec
        raise UnsupportedModeError(
            f"no converter available for source_type={source_type!r} conversion_mode={conversion_mode!r}"
        )

    def capabilities(self) -> CapabilitiesModel:
        converters = tuple(
            ConverterCapabilityModel(
                source_type=spec.source_type,
                conversion_mode=spec.conversion_mode,
                converter_name=spec.converter_name,
                converter_version=spec.converter_version,
                description=spec.description,
            )
            for spec in self._specs.values()
        )
        cleaners = (
            CleanerCapabilityModel(
                cleaner_name="identity_cleaner",
                cleaner_version="v1",
                description="Leave HTML content unchanged for the initial source-process implementation.",
            ),
        )
        return CapabilitiesModel(
            source_types=("html",),
            conversion_modes=tuple(sorted({spec.conversion_mode for spec in self._specs.values()})),
            converters=converters,
            cleaners=cleaners,
        )
