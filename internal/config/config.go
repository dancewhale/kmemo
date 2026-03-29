package config

import (
	"os"
	"time"
)

// Config holds host-side settings. Values are intentionally small; expand as features land.
type Config struct {
	PythonGRPCAddr string
	DialTimeout    time.Duration
	// SkipPython avoids dialing the worker; useful for UI / headless smoke without Python running.
	SkipPython bool
}

// Load reads environment with stable defaults suitable for local development.
func Load() Config {
	addr := os.Getenv("KMEMO_PYTHON_GRPC")
	if addr == "" {
		addr = "127.0.0.1:50051"
	}
	return Config{
		PythonGRPCAddr: addr,
		DialTimeout:    5 * time.Second,
		SkipPython:     os.Getenv("KMEMO_SKIP_PYTHON") == "1",
	}
}
