from __future__ import annotations


class SourceProcessError(Exception):
    """Base error for source-process runtime failures."""


class ValidationError(SourceProcessError):
    """Raised when a request is invalid at the boundary."""


class UnsupportedModeError(SourceProcessError):
    """Raised when no converter is available for the requested mode."""


class ConvertError(SourceProcessError):
    """Raised when conversion fails."""


class OutputWriteError(SourceProcessError):
    """Raised when output files cannot be written."""


class JobNotFoundError(SourceProcessError):
    """Raised when a job is not present in the in-memory store."""


class JobCancelledError(SourceProcessError):
    """Raised when a running job observes a cancel request."""
