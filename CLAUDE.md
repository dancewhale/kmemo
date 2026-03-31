# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

kmemo is a SuperMemo-style incremental reading system with a hybrid architecture:
- **Go**: Main application, desktop UI (Wails), storage layer (SQLite + GORM)
- **Python**: gRPC worker service for algorithms (FSRS scheduling, HTML processing, import handling)
- **Frontend**: Vite-based UI in `desktop/frontend/`

The project is currently a skeleton with placeholder implementations.

## Common Commands

All commands use [Task](https://taskfile.dev/). Run `task --list` to see available tasks.

### Initial Setup
```bash
task bootstrap  # Install all dependencies (Go, Python venv, npm, protoc plugins)
```

### Development Workflow
```bash
# Terminal A - Start Python gRPC worker
task run:python

# Terminal B - Start Wails desktop app (or headless Go)
task run:wails
# OR for headless without Python connection:
task run:go
# OR for headless with Python connection:
task run:go:connected
```

### Code Generation
```bash
task proto      # Regenerate gRPC code from proto/kmemo/v1/kmemo.proto
task db:gen     # Regenerate GORM Gen DAO code after model changes
task db:migrate # Run SQLite AutoMigrate to update schema
```

### Testing
```bash
task test           # Run all tests (Go + Python)
task test:go        # Go tests only
task test:python    # Python pytest only
```

### Building
```bash
task build          # Build both Go binary and Wails app
task build:go       # Build headless binary to build/kmemo
task build:wails    # Build desktop app (in desktop/build/)
```

### Database Management
```bash
task db:prepare     # Run migrate then codegen (after model edits)
task db:reset       # Delete dev DB and recreate (destructive)
```

## Architecture

### Go-Python Communication

The Go host communicates with Python via gRPC defined in `proto/kmemo/v1/kmemo.proto`:
- **CalculateFsrs**: FSRS spaced repetition scheduling
- **CleanHtml**: HTML sanitization and processing
- **PrepareImportMaterial**: Import format conversion

Go client: `internal/pyclient/client.go` wraps the generated gRPC client.
Python server: `python/app/server.py` implements the KmemoProcessor service.

Environment variables:
- `KMEMO_PYTHON_GRPC`: Python gRPC address (default: `127.0.0.1:50051`)
- `KMEMO_SKIP_PYTHON`: Set to `1` to skip Python connection (useful for Go-only development)

### Storage Layer

Located in `internal/storage/`:
- **Models**: `internal/storage/models/` defines GORM models (Knowledge, Card, Tag, Asset, SRS state, etc.)
- **DAO**: `internal/storage/dao/` contains GORM Gen generated code (DO NOT edit manually)
- **Generation**: Run `task db:gen` after editing models to regenerate DAO code
- **Migration**: `internal/storage/migrate.go` handles AutoMigrate; run via `task db:migrate`

The storage layer uses SQLite with GORM. Dev database: `kmemo.dev.db` in repo root.

### Wails Desktop Structure

- `desktop/main.go`: Wails entry point
- `desktop/app.go`: App struct with methods exposed to frontend
- `desktop/frontend/`: Vite + vanilla JS frontend
- Frontend calls Go via `window.go.main.App.*` methods

The `desktop/` directory is a Wails subproject but shares the parent Go module.

### Code Generation Workflow

1. **Proto changes**: Edit `proto/kmemo/v1/kmemo.proto` → run `task proto`
   - Generates: `gen/kmemo/v1/*.pb.go` (Go) and `python/generated/kmemo/v1/*_pb2*.py` (Python)
   - Never edit generated files manually

2. **Model changes**: Edit `internal/storage/models/*.go` → run `task db:prepare`
   - First runs `db:migrate` to update schema
   - Then runs `db:gen` to regenerate DAO code in `internal/storage/dao/`

## Important Notes

- Python 3.10-3.12 recommended (3.14+ may have grpcio compilation issues on Apple Silicon)
- If grpcio installation fails, delete `.venv` with `task clean:venv` and re-run `task bootstrap`
- Generated code directories: `gen/`, `python/generated/`, `internal/storage/dao/` - do not edit manually
- The project uses Task instead of Make; `Taskfile.yaml` is the source of truth for commands
- Default behavior skips Python connection for faster Go/Wails development; set `KMEMO_SKIP_PYTHON=0` when testing integration
