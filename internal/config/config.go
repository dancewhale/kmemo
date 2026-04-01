package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds host-side settings. Values are intentionally small; expand as features land.
type Config struct {
	PythonGRPCAddr string
	DialTimeout    time.Duration
	// SkipPython avoids dialing the worker; useful for UI / headless smoke without Python running.
	SkipPython      bool
	LogLevel        string
	RepositoryDebug bool
	DBSlowThreshold time.Duration
}

// Load reads environment with stable defaults suitable for local development.
func Load() Config {
	addr := os.Getenv("KMEMO_PYTHON_GRPC")
	if addr == "" {
		addr = "127.0.0.1:50051"
	}

	logLevel := os.Getenv("KMEMO_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	slowThreshold := 200 * time.Millisecond
	if value := os.Getenv("KMEMO_DB_SLOW_THRESHOLD_MS"); value != "" {
		if millis, err := strconv.Atoi(value); err == nil && millis > 0 {
			slowThreshold = time.Duration(millis) * time.Millisecond
		}
	}

	return Config{
		PythonGRPCAddr:  addr,
		DialTimeout:     5 * time.Second,
		SkipPython:      os.Getenv("KMEMO_SKIP_PYTHON") == "1",
		LogLevel:        logLevel,
		RepositoryDebug: os.Getenv("KMEMO_REPOSITORY_DEBUG") == "1",
		DBSlowThreshold: slowThreshold,
	}
}
