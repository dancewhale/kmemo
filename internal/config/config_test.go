package config

import "testing"

func TestLoad_DefaultPythonAddress(t *testing.T) {
	t.Setenv("KMEMO_PYTHON_GRPC", "")
	c := Load()
	if c.PythonGRPCAddr == "" {
		t.Fatalf("expected default grpc address")
	}
	if c.DialTimeout <= 0 {
		t.Fatalf("expected positive dial timeout")
	}
}
