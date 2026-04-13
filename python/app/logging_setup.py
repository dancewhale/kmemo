from __future__ import annotations

import logging
import os
import sys
from pathlib import Path

import structlog


def configure_logging(level: int = logging.INFO, log_file_path: str | None = None) -> None:
    timestamper = structlog.processors.TimeStamper(fmt="iso", key="ts")

    structlog.configure(
        processors=[
            structlog.contextvars.merge_contextvars,
            structlog.stdlib.add_log_level,
            timestamper,
            structlog.processors.StackInfoRenderer(),
            structlog.processors.format_exc_info,
            structlog.processors.EventRenamer("msg"),
            structlog.dev.ConsoleRenderer(),
        ],
        wrapper_class=structlog.make_filtering_bound_logger(level),
        logger_factory=structlog.stdlib.LoggerFactory(),
        cache_logger_on_first_use=True,
    )

    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setLevel(level)
    console_handler.setFormatter(logging.Formatter("%(message)s"))

    root = logging.getLogger()
    root.handlers.clear()
    root.setLevel(level)
    root.addHandler(console_handler)

    if log_file_path:
        Path(os.path.dirname(log_file_path)).mkdir(parents=True, exist_ok=True)
        file_handler = logging.FileHandler(log_file_path, encoding="utf-8")
        file_handler.setLevel(level)
        file_handler.setFormatter(logging.Formatter("%(message)s"))
        root.addHandler(file_handler)


def get_logger(name: str):
    return structlog.get_logger(name)
