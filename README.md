# kmemo

类似 SuperMemo 的**渐进阅读**系统：当前仓库为**工程骨架**（Go 主程序 + Python gRPC 子服务 + Wails 桌面壳），业务逻辑以 TODO / 占位实现为主。

## 顶层目录

```text
.
├── Taskfile.yaml          # 统一任务入口（不使用 Makefile）
├── proto/kmemo/v1/        # Protobuf / gRPC 定义（唯一事实来源）
├── gen/kmemo/v1/          # Go 生成代码（由 protoc 写入；也可按需纳入版本控制）
├── cmd/kmemo/             # 无 UI 主入口（编排 / 调试）
├── desktop/               # Wails 应用（package main + frontend/）
├── internal/              # Go 内部模块（config、bootstrap、pyclient、预留 storage 等）
└── python/                # Python gRPC worker（占位服务实现）
```

## 前置依赖

- Go 1.22+（若 `go` 报 GOROOT 错误，请用 Homebrew 修复或设置正确的 `GOROOT`）
- [Task](https://taskfile.dev/)（`go install github.com/go-task/task/v3/cmd/task@latest`）
- Node.js + npm（Wails 前端）
- [Wails CLI](https://wails.io/) v2（`go install github.com/wailsapp/wails/v2/cmd/wails@latest`）
- `task bootstrap` 会创建 `.venv` 并安装 `grpcio-tools`（内含 `protoc`），以及通过 `go install` 安装 `protoc-gen-go` / `protoc-gen-go-grpc`

## 一次性初始化

```bash
task bootstrap
```

## gRPC 代码生成

修改 `proto/` 后执行：

```bash
task proto
```

等价命令（需在仓库根目录，且 `$(go env GOPATH)/bin` 在 `PATH` 中）：

```bash
mkdir -p gen python/generated/kmemo/v1
.venv/bin/python -m grpc_tools.protoc -I proto \
  --go_out=gen --go_opt=paths=source_relative \
  --go-grpc_out=gen --go-grpc_opt=paths=source_relative \
  kmemo/v1/kmemo.proto

.venv/bin/python -m grpc_tools.protoc -I proto \
  --python_out=python/generated \
  --grpc_python_out=python/generated \
  kmemo/v1/kmemo.proto
```

**不要手写** `gen/**/*.pb.go`、`gen/**/*_grpc.pb.go`、`python/generated/kmemo/v1/kmemo_pb2*.py`；它们应由上述命令生成。

## 启动

- **Python worker**（默认监听 `[::]:50051`，等价 IPv4 `127.0.0.1:50051`）：

  ```bash
  task run:python
  ```

- **无 UI Go 主程序**（默认**不**连接 Python，便于单独起服务）：

  ```bash
  task run:go
  ```

  需要连接正在运行的 Python 时：

  ```bash
  task run:go:connected
  ```

  或通过环境变量：`KMEMO_SKIP_PYTHON=0 go run ./cmd/kmemo`

- **Wails 开发模式**（在 `desktop/` 下执行；默认同样跳过 Python 连接）：

  ```bash
  cd desktop/frontend && npm install   # 首次
  task run:wails
  ```

  需要桌面端连接 Python 时自行设置 `KMEMO_SKIP_PYTHON=0`。

**推荐日常开发**：终端 A `task run:python`，终端 B `task run:wails`（或 `task run:go:connected` 做联调）。

## 测试与构建

```bash
task test        # Go + Python
task test:go
task test:python

task build       # build:go + build:wails
task build:go
task build:wails
```

## 其他任务

- `task dev`：生成 proto 并提示多终端启动顺序
- `task clean`：删除生成代码目录 `gen/kmemo`、Python 的 `kmemo_pb2*.py` 以及常见构建产物
- `task lint` / `task format`：可选工具（未安装时任务会静默跳过相关步骤）

## 架构说明（骨架阶段）

| 组件 | 职责（当前） |
|------|----------------|
| **cmd/kmemo** | 进程入口、信号处理、预留编排 |
| **internal/bootstrap** | 组装 config、pyclient、Wails 绑定对象 |
| **internal/pyclient** | gRPC 客户端封装（生成代码在 `gen/`） |
| **internal/storage / htmlproc / indexing** | 仅占位包，为 SQLite / HTML / 索引预留 |
| **internal/services** | 未来协调层；现为占位 |
| **python/app** | gRPC 服务实现，按 `fsrs` / `html` / `import` 分文件占位 |
| **desktop** | Wails `main` + 最小 Vite 前端，经 `window.go.main.App` 调用 Go |

环境变量摘要：

| 变量 | 含义 |
|------|------|
| `KMEMO_PYTHON_GRPC` | Python gRPC 地址，默认 `127.0.0.1:50051` |
| `KMEMO_SKIP_PYTHON` | `1` 时不建立 gRPC 连接（`task run:go` / `task run:wails` 默认开启） |

## 与 Wails / Go 主工程的关系

- Go module 根在仓库根目录（`go.mod`）。`desktop/` 为 **Wails 子工程目录**：`wails.json` 与 `desktop/main.go` 同属该目录，`go build` 会从父目录解析 `kmemo/...` 包。
- 前端位于 `desktop/frontend/`；`desktop/frontend/dist/` 内保留最小占位 HTML，以便在未执行 `npm run build` 时 `go:embed` 仍可编译。发布前应在 `desktop/frontend` 执行 `npm run build`。
