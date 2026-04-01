package config

import "time"

// Config is the normalized runtime configuration consumed by the Go host.
type Config struct {
	AppName    string
	Profile    string
	ConfigDir  string
	ConfigFile string
	DataDir    string
	CacheDir   string
	LogsDir    string
	AssetsDir  string
	ImportsDir string
	IndexDir   string
	ExportDir  string
	CardFileDir string
	ThumbsCacheDir string
	HTMLCacheDir   string
	IndexCacheDir  string
	DBDriver   string
	DBPath     string
	LogFormat  string
	LogFileEnabled bool
	LogFileName string

	PythonGRPCAddr string
	DialTimeout    time.Duration
	// SkipPython avoids dialing the worker; useful for UI / headless smoke without Python running.
	SkipPython      bool
	LogLevel        string
	RepositoryDebug bool
	DBSlowThreshold time.Duration
}

type rawConfig struct {
	App      rawAppConfig      `koanf:"app" toml:"app"`
	Server   rawServerConfig   `koanf:"server" toml:"server"`
	Database rawDatabaseConfig `koanf:"database" toml:"database"`
	Logging  rawLoggingConfig  `koanf:"logging" toml:"logging"`
	Paths    rawPathsConfig    `koanf:"paths" toml:"paths"`
	Feature  rawFeatureConfig  `koanf:"feature" toml:"feature"`
}

type rawAppConfig struct {
	Name    string `koanf:"name" toml:"name"`
	Profile string `koanf:"profile" toml:"profile"`
}

type rawServerConfig struct {
	PythonGRPCAddr string `koanf:"python_grpc" toml:"python_grpc"`
	SkipPython     bool   `koanf:"skip_python" toml:"skip_python"`
	DialTimeoutMS  int    `koanf:"dial_timeout_ms" toml:"dial_timeout_ms"`
}

type rawDatabaseConfig struct {
	Driver          string `koanf:"driver" toml:"driver"`
	Path            string `koanf:"path" toml:"path"`
	SlowThresholdMS int    `koanf:"slow_threshold_ms" toml:"slow_threshold_ms"`
	RepositoryDebug bool   `koanf:"repository_debug" toml:"repository_debug"`
}

type rawLoggingConfig struct {
	Level       string `koanf:"level" toml:"level"`
	Format      string `koanf:"format" toml:"format"`
	FileEnabled bool   `koanf:"file_enabled" toml:"file_enabled"`
	FileName    string `koanf:"file_name" toml:"file_name"`
}

type rawPathsConfig struct {
	DataDir  string `koanf:"data_dir" toml:"data_dir"`
	CacheDir string `koanf:"cache_dir" toml:"cache_dir"`
	LogsDir  string `koanf:"logs_dir" toml:"logs_dir"`
}

type rawFeatureConfig struct {
	MultiProfile bool `koanf:"multi_profile" toml:"multi_profile"`
}

type runtimePaths struct {
	ConfigDir  string
	ConfigFile string
	DataDir    string
	CacheDir   string
	LogsDir    string
}
