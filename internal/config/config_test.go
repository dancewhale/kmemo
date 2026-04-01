package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad_DefaultPythonAddress(t *testing.T) {
	root := t.TempDir()
	t.Setenv("KMEMO_CONFIG_DIR", filepath.Join(root, "config"))
	t.Setenv("KMEMO_DATA_DIR", filepath.Join(root, "data"))
	t.Setenv("KMEMO_CACHE_DIR", filepath.Join(root, "cache"))
	t.Setenv("KMEMO_LOGS_DIR", filepath.Join(root, "logs"))
	t.Setenv("KMEMO_PYTHON_GRPC", "")
	t.Setenv("KMEMO_LOG_LEVEL", "")
	t.Setenv("KMEMO_REPOSITORY_DEBUG", "")
	t.Setenv("KMEMO_DB_SLOW_THRESHOLD_MS", "")

	c, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if c.AppName != defaultAppName {
		t.Fatalf("expected app name %q, got %q", defaultAppName, c.AppName)
	}
	if c.Profile != defaultProfileName {
		t.Fatalf("expected profile %q, got %q", defaultProfileName, c.Profile)
	}
	if c.PythonGRPCAddr == "" {
		t.Fatalf("expected default grpc address")
	}
	if c.DialTimeout <= 0 {
		t.Fatalf("expected positive dial timeout")
	}
	if c.LogLevel != "info" {
		t.Fatalf("expected default log level info, got %q", c.LogLevel)
	}
	if c.LogFormat != defaultLogFormat {
		t.Fatalf("expected default log format %q, got %q", defaultLogFormat, c.LogFormat)
	}
	if c.RepositoryDebug {
		t.Fatalf("expected repository debug disabled by default")
	}
	if c.DBDriver != defaultDBDriver {
		t.Fatalf("expected default db driver %q, got %q", defaultDBDriver, c.DBDriver)
	}
	if c.DBSlowThreshold != 200*time.Millisecond {
		t.Fatalf("expected default slow threshold 200ms, got %v", c.DBSlowThreshold)
	}
	if c.AssetsDir != filepath.Join(c.DataDir, defaultAssetsDirName) {
		t.Fatalf("expected assets dir under data dir, got %q", c.AssetsDir)
	}
}

func TestLoad_LogOverrides(t *testing.T) {
	root := t.TempDir()
	t.Setenv("KMEMO_CONFIG_DIR", filepath.Join(root, "config"))
	t.Setenv("KMEMO_DATA_DIR", filepath.Join(root, "data"))
	t.Setenv("KMEMO_CACHE_DIR", filepath.Join(root, "cache"))
	t.Setenv("KMEMO_LOGS_DIR", filepath.Join(root, "logs"))
	t.Setenv("KMEMO_LOG_LEVEL", "debug")
	t.Setenv("KMEMO_REPOSITORY_DEBUG", "1")
	t.Setenv("KMEMO_DB_SLOW_THRESHOLD_MS", "450")

	c, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if c.LogLevel != "debug" {
		t.Fatalf("expected debug log level, got %q", c.LogLevel)
	}
	if !c.RepositoryDebug {
		t.Fatalf("expected repository debug enabled")
	}
	if c.DBSlowThreshold != 450*time.Millisecond {
		t.Fatalf("expected 450ms slow threshold, got %v", c.DBSlowThreshold)
	}
}

func TestResolveRuntimePathsForWindowsCacheLayout(t *testing.T) {
	paths, err := resolveRuntimePathsFor("windows", map[string]string{
		"APPDATA":      `C:\Users\whale\AppData\Roaming`,
		"LOCALAPPDATA": `C:\Users\whale\AppData\Local`,
	}, `C:\work\kmemo`, `C:\Users\whale`)
	if err != nil {
		t.Fatalf("resolve runtime paths: %v", err)
	}
	if paths.CacheDir != `C:\Users\whale\AppData\Local\kmemo\cache` {
		t.Fatalf("expected windows cache dir with cache suffix, got %q", paths.CacheDir)
	}
	if paths.DataDir != `C:\Users\whale\AppData\Local\kmemo\data` {
		t.Fatalf("expected windows data dir, got %q", paths.DataDir)
	}
}

func TestEnsureInitializedCreatesRequiredDirsAndConfig(t *testing.T) {
	root := t.TempDir()
	paths := runtimePaths{
		ConfigDir:  filepath.Join(root, "config"),
		ConfigFile: filepath.Join(root, "config", defaultConfigFileName),
		DataDir:    filepath.Join(root, "data"),
		CacheDir:   filepath.Join(root, "cache"),
		LogsDir:    filepath.Join(root, "logs"),
	}
	if err := ensureInitialized(paths); err != nil {
		t.Fatalf("ensure initialized: %v", err)
	}
	for _, p := range append(requiredRuntimeDirs(paths), paths.ConfigFile) {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected path %q to exist: %v", p, err)
		}
	}
}

func TestValidateConfigRejectsMissingPythonAddrWhenEnabled(t *testing.T) {
	cfg := defaultConfig()
	cfg.ConfigDir = "/tmp/config"
	cfg.ConfigFile = "/tmp/config/config.toml"
	cfg.DataDir = "/tmp/data"
	cfg.CacheDir = "/tmp/cache"
	cfg.LogsDir = "/tmp/logs"
	cfg.DBPath = "/tmp/data/kmemo.db"
	cfg.SkipPython = false
	cfg.PythonGRPCAddr = ""
	if err := validateConfig(cfg); err == nil {
		t.Fatalf("expected validation error when python grpc is empty")
	}
}
