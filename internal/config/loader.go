package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/knadh/koanf/maps"
	"github.com/knadh/koanf/providers/confmap"
	koanf "github.com/knadh/koanf/v2"
	tomlv2 "github.com/pelletier/go-toml/v2"
)

func loadConfig() (Config, error) {
	paths, err := resolveRuntimePaths()
	if err != nil {
		return Config{}, err
	}
	if err := ensureInitialized(paths); err != nil {
		return Config{}, err
	}

	k := koanf.New(".")
	if err := k.Load(confmap.Provider(rawConfigMap(defaultRawConfig()), "."), nil); err != nil {
		return Config{}, fmt.Errorf("load config defaults: %w", err)
	}

	fileRaw, err := readRawConfigFile(paths.ConfigFile)
	if err != nil {
		return Config{}, err
	}
	if err := k.Load(confmap.Provider(rawConfigMap(fileRaw), "."), nil); err != nil {
		return Config{}, fmt.Errorf("load config file map: %w", err)
	}

	raw := defaultRawConfig()
	if err := k.Unmarshal("", &raw); err != nil {
		return Config{}, fmt.Errorf("unmarshal merged config: %w", err)
	}
	applyEnvOverrides(&raw, configEnvMap(os.Getenv))

	cfg, err := normalizeConfig(paths, raw)
	if err != nil {
		return Config{}, err
	}
	if err := validateConfig(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func readRawConfigFile(filePath string) (rawConfig, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return rawConfig{}, fmt.Errorf("read config file: %w", err)
	}
	if len(content) == 0 {
		return rawConfig{}, nil
	}
	var cfg rawConfig
	if err := tomlv2.Unmarshal(content, &cfg); err != nil {
		return rawConfig{}, fmt.Errorf("parse config file %s: %w", filePath, err)
	}
	return cfg, nil
}

func normalizeConfig(paths runtimePaths, raw rawConfig) (Config, error) {
	cfg := defaultConfig()
	cfg.AppName = firstNonEmpty(raw.App.Name, defaultAppName)
	cfg.Profile = firstNonEmpty(raw.App.Profile, defaultProfileName)
	cfg.ConfigDir = paths.ConfigDir
	cfg.ConfigFile = paths.ConfigFile
	cfg.PythonGRPCAddr = raw.Server.PythonGRPCAddr
	cfg.SkipPython = raw.Server.SkipPython
	cfg.DialTimeout = time.Duration(raw.Server.DialTimeoutMS) * time.Millisecond
	cfg.DBDriver = firstNonEmpty(raw.Database.Driver, defaultDBDriver)
	cfg.LogLevel = raw.Logging.Level
	cfg.LogFormat = firstNonEmpty(raw.Logging.Format, defaultLogFormat)
	cfg.LogFileEnabled = raw.Logging.FileEnabled
	cfg.LogFileName = firstNonEmpty(raw.Logging.FileName, defaultLogFileName)
	cfg.RepositoryDebug = raw.Database.RepositoryDebug
	cfg.DBSlowThreshold = time.Duration(raw.Database.SlowThresholdMS) * time.Millisecond

	cwd, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("get working directory: %w", err)
	}
	cfg.DataDir = normalizeOptionalPath(raw.Paths.DataDir, paths.DataDir, paths.ConfigDir)
	cfg.CacheDir = normalizeOptionalPath(raw.Paths.CacheDir, paths.CacheDir, paths.ConfigDir)
	cfg.LogsDir = normalizeOptionalPath(raw.Paths.LogsDir, paths.LogsDir, paths.ConfigDir)

	if cfg.DataDir == "" {
		cfg.DataDir = paths.DataDir
	}
	if cfg.CacheDir == "" {
		cfg.CacheDir = paths.CacheDir
	}
	if cfg.LogsDir == "" {
		cfg.LogsDir = paths.LogsDir
	}

	cfg.AssetsDir = filepath.Join(cfg.DataDir, defaultAssetsDirName)
	cfg.ImportsDir = filepath.Join(cfg.DataDir, defaultImportsDirName)
	cfg.IndexDir = filepath.Join(cfg.DataDir, defaultIndexDirName)
	cfg.ExportDir = filepath.Join(cfg.DataDir, defaultExportDirName)
	cfg.CardFileDir = filepath.Join(cfg.DataDir, defaultCardFileDirName)
	cfg.ThumbsCacheDir = filepath.Join(cfg.CacheDir, defaultThumbnailsDirName)
	cfg.HTMLCacheDir = filepath.Join(cfg.CacheDir, defaultHTMLCacheDirName)
	cfg.IndexCacheDir = filepath.Join(cfg.CacheDir, defaultIndexCacheDirName)

	if envDBPath := os.Getenv("KMEMO_DB_PATH"); envDBPath != "" {
		cfg.DBPath = absolutize(envDBPath, cwd)
	} else if raw.Database.Path != "" {
		cfg.DBPath = absolutize(raw.Database.Path, paths.ConfigDir)
	} else {
		cfg.DBPath = filepath.Join(cfg.DataDir, defaultDBFileName)
	}
	cfg.DBPath = filepath.Clean(cfg.DBPath)
	return cfg, nil
}

func normalizeOptionalPath(value, fallback, baseDir string) string {
	if value == "" {
		return filepath.Clean(fallback)
	}
	value = expandHome(value)
	if filepath.IsAbs(value) {
		return filepath.Clean(value)
	}
	return filepath.Clean(filepath.Join(baseDir, value))
}

func validateConfig(cfg Config) error {
	if cfg.ConfigDir == "" || cfg.ConfigFile == "" || cfg.DataDir == "" || cfg.CacheDir == "" || cfg.LogsDir == "" {
		return fmt.Errorf("config paths must not be empty")
	}
	if cfg.DBDriver == "" {
		return fmt.Errorf("database driver must not be empty")
	}
	if cfg.DBPath == "" {
		return fmt.Errorf("database path must not be empty")
	}
	if cfg.DialTimeout <= 0 {
		return fmt.Errorf("dial timeout must be positive")
	}
	if cfg.DBSlowThreshold <= 0 {
		return fmt.Errorf("db slow threshold must be positive")
	}
	if !cfg.SkipPython && strings.TrimSpace(cfg.PythonGRPCAddr) == "" {
		return fmt.Errorf("python grpc address must not be empty when skip_python is false")
	}
	return nil
}

func rawConfigMap(raw rawConfig) map[string]interface{} {
	return maps.Copy(map[string]interface{}{
		"app": map[string]interface{}{
			"name":    raw.App.Name,
			"profile": raw.App.Profile,
		},
		"server": map[string]interface{}{
			"python_grpc":     raw.Server.PythonGRPCAddr,
			"skip_python":     raw.Server.SkipPython,
			"dial_timeout_ms": raw.Server.DialTimeoutMS,
		},
		"database": map[string]interface{}{
			"driver":            raw.Database.Driver,
			"path":              raw.Database.Path,
			"slow_threshold_ms": raw.Database.SlowThresholdMS,
			"repository_debug":  raw.Database.RepositoryDebug,
		},
		"logging": map[string]interface{}{
			"level":        raw.Logging.Level,
			"format":       raw.Logging.Format,
			"file_enabled": raw.Logging.FileEnabled,
			"file_name":    raw.Logging.FileName,
		},
		"paths": map[string]interface{}{
			"data_dir":  raw.Paths.DataDir,
			"cache_dir": raw.Paths.CacheDir,
			"logs_dir":  raw.Paths.LogsDir,
		},
		"feature": map[string]interface{}{
			"multi_profile": raw.Feature.MultiProfile,
		},
	})
}
