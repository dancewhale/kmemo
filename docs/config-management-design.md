# 配置文件管理模块设计文档

## 1. 目标

本文档用于为 kmemo 设计一套适合当前项目阶段、同时支持后续演进的配置文件管理模块。

当前项目是一个 `Wails + Go + Python gRPC + SQLite` 的桌面应用骨架，因此配置模块不仅要解决“读取几个环境变量”的问题，还要统一解决下面这些事情：

- 跨平台管理配置目录、数据目录、缓存目录、日志目录
- 程序首次启动时自动创建必要目录和默认配置文件
- 支持 macOS / Linux / Windows 的标准用户目录约定
- 支持“默认值 + 配置文件 + 环境变量 + 命令行参数”的分层覆盖
- 让 Go Host、Wails Desktop、Python Worker 后续都能共享一致的配置来源
- 后续接入数据库路径、日志、索引、文件存储、Python worker 地址时不需要重写配置体系

---

## 2. 设计原则

## 2.1 配置管理要成为独立基础模块

配置不应继续散落在：

- `internal/config/config.go` 里手写 `os.Getenv(...)`
- `Taskfile` 的临时 env
- 启动入口里的路径拼接逻辑

建议把配置视为一个明确的基础能力，单独沉淀为：

```text
internal/config/
├── manager.go        # 配置加载与初始化总入口
├── model.go          # Config / AppConfig / 路径结构定义
├── defaults.go       # 默认值
├── loader.go         # 分层加载逻辑
├── init_fs.go        # 首次启动目录 / 文件初始化
├── path.go           # 跨平台目录解析
└── env.go            # 环境变量映射
```

这样后续：

- bootstrap 只负责调用配置模块
- app / storage / zaplog / pyclient 只消费配置结果
- 不再各处自行决定配置目录和路径规则

---

## 2.2 配置来源必须分层且可覆盖

建议采用以下优先级，从低到高：

```text
默认值
  < 配置文件
  < 环境变量
  < 命令行参数（如后续接入）
```

这是桌面应用和本地开发工具都比较稳定的做法。

这样可以同时满足：

- 普通用户依赖配置文件
- 开发环境临时用 env 覆盖
- CI / 调试脚本用命令行参数覆盖

---

## 2.3 首次启动必须自动初始化，而不是要求手工准备

程序启动后，应自动保证以下事实成立：

- 配置目录存在
- 数据目录存在
- 缓存目录存在
- 日志目录存在
- 主配置文件存在
- 必要的子目录存在

如果不存在，则自动创建。

这比“让用户先自己创建目录 / 拷贝模板”更符合桌面应用和跨平台应用的常见体验。

---

## 2.4 配置文件应人类可读、可编辑、易扩展

建议主配置文件使用 **TOML**。

原因：

- 适合桌面 / 本地应用的静态配置
- 比 JSON 更适合手工编辑
- 比 YAML 语义更收敛，坑更少
- Go 侧生态成熟
- 层级结构清晰，适合后续继续扩展

不建议当前阶段使用：

- JSON：可读性一般，不适合带注释的手工维护
- YAML：功能强但语义过于灵活，容易引入解析歧义
- INI：不适合后续较复杂的嵌套配置

---

## 2.5 配置目录与数据目录必须分离

不要把所有内容都塞进一个目录。

建议至少区分：

- 配置（config）
- 数据（data）
- 缓存（cache）
- 日志（logs）

原因：

- 配置适合备份/迁移
- 数据可能很大，不应和配置混在一起
- 缓存可清理
- 日志应单独轮转和清理

这也更符合三大桌面系统的常见目录约定。

---

## 3. 底层库选型

## 3.1 配置聚合与覆盖：`koanf`

推荐使用：

- [`knadh/koanf`](https://github.com/knadh/koanf)

### 选择原因

`koanf` 是 Go 社区里成熟、轻量、可组合的配置聚合库，适合当前项目。

它的优点：

- 支持多来源合并：文件、环境变量、flag、map 默认值
- 结构化配置模型清晰
- 不强迫项目绑定某一种来源
- 比较适合“默认值 + 文件 + env + flag”这种分层加载模型
- 比 `viper` 更轻、更可控，隐式行为更少

### 在本项目中的定位

`koanf` 负责：

- 合并默认值
- 读取 `config.toml`
- 读取 `KMEMO_*` 环境变量
- 后续接入 CLI flags
- 最终反序列化到 `config.Config`

---

## 3.2 TOML 解析：`go-toml/v2`

推荐使用：

- [`pelletier/go-toml/v2`](https://github.com/pelletier/go-toml)

### 选择原因

- 成熟、维护活跃
- 性能和兼容性较稳定
- 生态中使用广泛
- 适合配合 `koanf` 的 TOML parser/provider 使用

---

## 3.3 跨平台标准目录：Go 标准库 `os.UserConfigDir` / `os.UserCacheDir` + 必要的显式规则封装

推荐优先使用：

- `os.UserConfigDir()`
- `os.UserCacheDir()`
- `os.UserHomeDir()`

### 选择原因

这几个 API 已由 Go 标准库提供，跨平台兼容性足够稳定：

- macOS -> `~/Library/Application Support`
- Linux -> `~/.config` / XDG 约定
- Windows -> `%AppData%`

对于 data / logs 目录，建议在配置模块内部统一封装平台规则，而不是让上层业务自己拼路径。

### 为什么这里不强行再引入额外第三方目录库

这里优先选择标准库，不是因为第三方不成熟，而是因为：

- 跨平台目录解析本身已经被标准库稳定覆盖
- 这类逻辑越基础越应减少外部依赖
- 配置模块可以把平台差异集中封装起来，上层不感知

也就是说，本项目的“成熟开源选型”重点应放在配置聚合与解析上，而目录定位优先使用 Go 官方稳定 API。

---

## 4. 跨平台目录设计

建议统一定义应用名：

```text
app_name = kmemo
```

并在三个系统上使用各自推荐路径。

## 4.1 配置目录

### macOS

```text
~/Library/Application Support/kmemo/
```

### Linux

优先：

```text
$XDG_CONFIG_HOME/kmemo/
```

若未设置，则：

```text
~/.config/kmemo/
```

### Windows

```text
%AppData%\kmemo\
```

---

## 4.2 数据目录

### macOS

```text
~/Library/Application Support/kmemo/data/
```

### Linux

优先：

```text
$XDG_DATA_HOME/kmemo/data/
```

若未设置，则：

```text
~/.local/share/kmemo/data/
```

### Windows

```text
%LocalAppData%\kmemo\data\
```

---

## 4.3 缓存目录

### macOS

```text
~/Library/Caches/kmemo/
```

### Linux

优先：

```text
$XDG_CACHE_HOME/kmemo/
```

若未设置，则：

```text
~/.cache/kmemo/
```

### Windows

```text
%LocalAppData%\kmemo\cache\
```

---

## 4.4 日志目录

建议：

- 默认将日志目录置于 data 体系下的独立子目录
- 当前阶段推荐：

### macOS

```text
~/Library/Application Support/kmemo/logs/
```

### Linux

```text
$XDG_DATA_HOME/kmemo/logs/
```

先统一放在：

```text
<data_dir>/logs/
```

### Windows

```text
%LocalAppData%\kmemo\logs\
```

---

## 4.5 建议最终目录布局

建议运行时统一拿到如下目录结构：

```text
config_dir/
├── config.toml
├── config.example.toml        # 可选，便于排查与恢复
└── profiles/                  # 预留，多配置场景未来扩展

data_dir/
├── kmemo.db
├── assets/
├── imports/
├── index/
├── export/
├── cardfile/
└── logs/

cache_dir/
├── thumbnails/
├── html/
└── index/
```

---

## 5. 配置文件结构设计

建议主配置文件：

```text
config.toml
```

建议结构：

```toml
[app]
name = "kmemo"
profile = "default"

[server]
python_grpc = "127.0.0.1:50051"
skip_python = true
dial_timeout_ms = 5000

[database]
driver = "sqlite"
path = ""
slow_threshold_ms = 200
repository_debug = false

[logging]
level = "info"
format = "console"
file_enabled = false
file_name = "kmemo.log"

[paths]
data_dir = ""
cache_dir = ""
logs_dir = ""
assets_dir = ""

[feature]
enable_search = true
enable_import = true
```

说明：

- 留空字符串表示“由程序自动推导默认路径”
- 用户可以手工覆盖特定目录
- 当前字段数量控制在“够用”范围，不提前设计过多配置项

---

## 6. 配置模型设计

建议 `internal/config/model.go` 中定义一个比当前更完整的结构：

```go
type Config struct {
    App      AppConfig      `koanf:"app"`
    Server   ServerConfig   `koanf:"server"`
    Database DatabaseConfig `koanf:"database"`
    Logging  LoggingConfig  `koanf:"logging"`
    Paths    PathsConfig    `koanf:"paths"`
    Feature  FeatureConfig  `koanf:"feature"`
}

type AppConfig struct {
    Name    string `koanf:"name"`
    Profile string `koanf:"profile"`
}

type ServerConfig struct {
    PythonGRPCAddr string        `koanf:"python_grpc"`
    SkipPython     bool          `koanf:"skip_python"`
    DialTimeout    time.Duration `koanf:"dial_timeout"`
}

type DatabaseConfig struct {
    Driver          string        `koanf:"driver"`
    Path            string        `koanf:"path"`
    SlowThreshold   time.Duration `koanf:"slow_threshold"`
    RepositoryDebug bool          `koanf:"repository_debug"`
}

type LoggingConfig struct {
    Level       string `koanf:"level"`
    Format      string `koanf:"format"`
    FileEnabled bool   `koanf:"file_enabled"`
    FileName    string `koanf:"file_name"`
}

type PathsConfig struct {
    ConfigDir   string `koanf:"config_dir"`
    DataDir     string `koanf:"data_dir"`
    CacheDir    string `koanf:"cache_dir"`
    LogsDir     string `koanf:"logs_dir"`
    AssetsDir   string `koanf:"assets_dir"`
    CardFileDir string `koanf:"card_file_dir"`
}
```

注意：

- 对外暴露给业务侧的结构应是“已归一化结果”
- 配置文件中的原始字符串值，加载后应经过 `normalize()` 二次处理
- 比如相对路径、空路径、时间单位等，都应在配置模块内部整理好

---

## 7. 启动初始化流程设计

## 7.1 初始化目标

启动时，配置模块要完成两件事：

1. **Resolve**：计算本机应使用的配置/数据/缓存目录
2. **Ensure**：确保这些目录和基础文件存在

---

## 7.2 推荐流程

建议在 `bootstrap.NewHeadless(ctx)` 最前面调用：

```text
config.InitializeAndLoad(ctx)
```

内部流程如下：

```text
1. 解析平台目录
2. 生成默认路径集合
3. 创建 config/data/cache/logs 目录
4. 若 config.toml 不存在，则写入默认配置模板
5. 加载默认值
6. 合并 config.toml
7. 合并环境变量 KMEMO_*
8. 合并命令行参数（未来）
9. 归一化路径与 duration
10. 校验关键配置
11. 返回最终 Config
```

---

## 7.3 首次启动时需要创建的内容

建议最少创建：

```text
config_dir/
data_dir/
cache_dir/
logs_dir/
cardfile_dir/
config_dir/config.toml
```

其中：

- `config.toml` 仅在文件不存在时创建
- 如果文件已存在，绝不覆盖用户内容
- 目录创建使用 `os.MkdirAll`
- 文件创建使用“仅不存在时写入”的方式，避免误覆盖

---

## 7.4 初始化失败处理策略

### 可恢复失败

例如：

- 示例配置文件创建失败，但主配置文件已存在
- cache 目录创建失败，但当前路径暂时不用

这类问题可以：

- 记录 `warn`
- 尽量继续启动

### 不可恢复失败

例如：

- 无法创建主配置目录
- 无法创建 data 目录
- 主配置文件不可读且没有其他覆盖来源
- 数据库路径最终不可用

这类问题应：

- 直接返回错误
- 由 bootstrap / main 决定退出

---

## 8. 环境变量设计

建议保留并扩展现有 `KMEMO_*` 前缀。

当前已存在：

- `KMEMO_PYTHON_GRPC`
- `KMEMO_SKIP_PYTHON`
- `KMEMO_LOG_LEVEL`
- `KMEMO_REPOSITORY_DEBUG`
- `KMEMO_DB_SLOW_THRESHOLD_MS`

建议新增并规范：

- `KMEMO_CONFIG_FILE`
- `KMEMO_CONFIG_DIR`
- `KMEMO_DATA_DIR`
- `KMEMO_CACHE_DIR`
- `KMEMO_LOGS_DIR`
- `KMEMO_DB_PATH`
- `KMEMO_LOG_FILE_ENABLED`

建议映射规则：

```text
KMEMO_SERVER_PYTHON_GRPC        -> server.python_grpc
KMEMO_SERVER_SKIP_PYTHON        -> server.skip_python
KMEMO_DATABASE_PATH             -> database.path
KMEMO_LOGGING_LEVEL             -> logging.level
```

但考虑当前项目已经有存量环境变量，建议第一阶段兼容旧命名。

即：

- **内部标准字段**采用分组命名
- **外部环境变量**先兼容旧字段
- 后续再逐步收敛

---

## 9. 默认值策略

建议把默认值集中定义在一个地方，而不是散落在各处。

例如：

```go
func DefaultConfig() Config {
    return Config{
        App: AppConfig{
            Name:    "kmemo",
            Profile: "default",
        },
        Server: ServerConfig{
            PythonGRPCAddr: "127.0.0.1:50051",
            SkipPython:     true,
            DialTimeout:    5 * time.Second,
        },
        Database: DatabaseConfig{
            Driver:          "sqlite",
            SlowThreshold:   200 * time.Millisecond,
            RepositoryDebug: false,
        },
        Logging: LoggingConfig{
            Level:       "info",
            Format:      "console",
            FileEnabled: false,
            FileName:    "kmemo.log",
        },
    }
}
```

关键原则：

- 默认值应该让程序在首次启动时就能跑起来
- 默认值不应依赖仓库根目录
- 默认值应优先面向真实用户目录，而不是开发者工作区

---

## 10. 路径归一化策略

配置加载完成后，必须统一做路径归一化。

建议规则：

- 若 `database.path` 为空，则自动落到 `data_dir/kmemo.db`
- 若 `paths.logs_dir` 为空，则自动落到默认日志目录
- 若 `paths.assets_dir` 为空，则自动落到 `data_dir/assets`
- 相对路径统一转绝对路径
- `~` 开头路径统一展开到用户 home
- 保证所有关键路径最终都是可直接使用的绝对路径

这一步非常重要。

否则后续：

- bootstrap
- storage
- logging
- filestore

都会各自再做一次路径补丁，配置模块就失去意义了。

---

## 11. 建议对外暴露的接口

建议配置模块对外只暴露少量明确接口：

```go
type Manager interface {
    Initialize(ctx context.Context) (*RuntimePaths, error)
    Load(ctx context.Context) (Config, error)
    InitializeAndLoad(ctx context.Context) (Config, error)
}
```

或者更直接一些：

```go
func InitializeAndLoad(ctx context.Context) (Config, error)
func ResolvePaths() (RuntimePaths, error)
```

其中：

```go
type RuntimePaths struct {
    ConfigDir  string
    ConfigFile string
    DataDir    string
    CacheDir   string
    LogsDir    string
}
```

建议：

- 上层默认只调用 `InitializeAndLoad`
- `ResolvePaths` 和更细粒度方法主要给测试与工具命令使用

---

## 12. 与当前项目其他模块的关系

## 12.1 bootstrap

`internal/bootstrap/bootstrap.go` 不应继续直接依赖旧式 `config.Load()` 环境变量读取逻辑。

建议改为：

```text
bootstrap -> config.InitializeAndLoad -> 得到标准化 Config -> 注入 logger / pyclient / storage
```

---

## 12.2 zaplog

日志模块不应自己决定日志目录。

应由配置模块提供：

- `logging.level`
- `logging.format`
- `logging.file_enabled`
- `paths.logs_dir`

---

## 12.3 storage

数据库模块不应自己兜底数据库文件路径。

应由配置模块提供最终：

- `database.driver`
- `database.path`
- `database.slow_threshold`
- `database.repository_debug`

---

## 12.4 pyclient

Python gRPC 客户端只消费：

- `server.python_grpc`
- `server.skip_python`
- `server.dial_timeout`

不应自行从 env 再读一遍。

---

## 12.5 后续 contracts / adapters

后续文件存储、索引、导入模块都应从配置模块拿路径，而不是自行基于 cwd 拼接。

例如：

- 搜索索引 -> `cache_dir/index` 或 `data_dir/index`
- 文件存储 -> `data_dir/assets` `data_dir/cardfile` 
- 导入中间产物 -> `cache_dir/imports`

---

## 13. 示例启动时序

建议未来启动链路变成：

```text
main / wails entry
    ↓
bootstrap.NewHeadless(ctx)
    ↓
config.InitializeAndLoad(ctx)
    ↓
创建 logger
    ↓
记录配置摘要（不打印敏感信息）
    ↓
创建 storage / pyclient / desktop app
```

首次启动时：

```text
1. 发现 config_dir 不存在
2. 自动创建目录
3. 发现 config.toml 不存在
4. 写入默认模板
5. 继续按“默认值 + 文件 + env”加载
6. 启动成功
```

这就是配置模块需要提供的核心体验。

---

## 14. 验证与测试建议

配置模块实现后，建议至少覆盖以下测试：

## 14.1 路径解析测试

- macOS 路径规则
- Linux 路径规则
- Windows 路径规则
- XDG 环境变量覆盖
- 用户自定义 `KMEMO_CONFIG_DIR` 覆盖

## 14.2 初始化测试

- 首次启动自动创建目录
- 缺失 `config.toml` 时自动创建
- 已存在 `config.toml` 时不覆盖
- 部分目录已存在时仍可成功

## 14.3 加载优先级测试

- 默认值生效
- 配置文件覆盖默认值
- 环境变量覆盖配置文件
- 命令行参数覆盖环境变量（后续）

## 14.4 归一化测试

- 空数据库路径自动补全
- 相对路径转绝对路径
- `~` 路径展开
- duration / bool / int 字段解析正确

---

## 15. 推荐实施顺序

为了保持改动可控，建议按下面顺序落地：

### 第一步：先完成配置目录与默认配置文件初始化

先让程序支持：

- 找到标准配置目录
- 自动创建目录
- 自动创建 `config.toml`

此时即使仍保留部分旧 env 配置，也没问题。

### 第二步：把现有 `internal/config/config.go` 升级为分层配置加载

接入：

- 默认值
- TOML 文件
- 环境变量

并兼容当前已有 `KMEMO_*` 字段。

### 第三步：让 bootstrap 全部改为消费新 Config

把：

- logger
- storage
- pyclient

都切到新配置模型。

### 第四步：后续再加命令行 flags 与配置热重载

配置热重载不是当前阶段必须能力。

当前项目重点是：

- 目录初始化
- 配置文件初始化
- 稳定加载
- 跨平台路径一致

---

## 16. 最终建议

结合当前 kmemo 的阶段，建议采用下面这套方案：

- **配置聚合库**：`koanf`
- **TOML 解析库**：`go-toml/v2`
- **跨平台目录解析**：优先 Go 标准库 `os.UserConfigDir` / `os.UserCacheDir` / `os.UserHomeDir`，由 `internal/config/path.go` 统一封装平台规则
- **配置文件格式**：`config.toml`
- **加载优先级**：默认值 < 配置文件 < 环境变量 < 命令行参数
- **启动行为**：首次启动自动创建目录与默认配置文件
- **目录策略**：配置、数据、缓存、日志分离

这是一个足够稳定、跨平台、可理解、后续易扩展的方案。

它既不会像只靠 `os.Getenv` 那样过于原始，也不会为了“高级配置系统”引入超出当前项目阶段的复杂度。
