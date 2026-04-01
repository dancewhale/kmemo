package config

import (
	"testing"
	"time"
)

func TestLoad_DefaultPythonAddress(t *testing.T) {
	t.Setenv("KMEMO_PYTHON_GRPC", "")
	t.Setenv("KMEMO_LOG_LEVEL", "")
	t.Setenv("KMEMO_REPOSITORY_DEBUG", "")
	t.Setenv("KMEMO_DB_SLOW_THRESHOLD_MS", "")

	c := Load()
	if c.PythonGRPCAddr == "" {
		t.Fatalf("expected default grpc address")
	}
	if c.DialTimeout <= 0 {
		t.Fatalf("expected positive dial timeout")
	}
	if c.LogLevel != "info" {
		t.Fatalf("expected default log level info, got %q", c.LogLevel)
	}
	if c.RepositoryDebug {
		t.Fatalf("expected repository debug disabled by default")
	}
	if c.DBSlowThreshold != 200*time.Millisecond {
		t.Fatalf("expected default slow threshold 200ms, got %v", c.DBSlowThreshold)
	}
}

func TestLoad_LogOverrides(t *testing.T) {
	t.Setenv("KMEMO_LOG_LEVEL", "debug")
	t.Setenv("KMEMO_REPOSITORY_DEBUG", "1")
	t.Setenv("KMEMO_DB_SLOW_THRESHOLD_MS", "450")

	c := Load()
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
