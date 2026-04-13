package config

import "strconv"

func configEnvMap(getenv func(string) string) map[string]string {
	keys := []string{
		"KMEMO_PYTHON_GRPC",
		"KMEMO_SKIP_PYTHON",
		"KMEMO_LOG_LEVEL",
		"KMEMO_LOG_FILE_ENABLED",
		"KMEMO_REPOSITORY_DEBUG",
		"KMEMO_DB_SLOW_THRESHOLD_MS",
		"KMEMO_DB_PATH",
		"KMEMO_DATA_DIR",
		"KMEMO_CACHE_DIR",
		"KMEMO_LOGS_DIR",
	}
	m := make(map[string]string, len(keys))
	for _, key := range keys {
		m[key] = getenv(key)
	}
	return m
}

func applyEnvOverrides(raw *rawConfig, env map[string]string) {
	if value := env["KMEMO_PYTHON_GRPC"]; value != "" {
		raw.Server.PythonGRPCAddr = value
	}
	if value := env["KMEMO_SKIP_PYTHON"]; value != "" {
		raw.Server.SkipPython = parseEnvBool(value)
	}
	if value := env["KMEMO_LOG_LEVEL"]; value != "" {
		raw.Logging.Level = value
	}
	if value := env["KMEMO_LOG_FILE_ENABLED"]; value != "" {
		raw.Logging.FileEnabled = parseEnvBool(value)
	}
	if value := env["KMEMO_REPOSITORY_DEBUG"]; value != "" {
		raw.Database.RepositoryDebug = parseEnvBool(value)
	}
	if value := env["KMEMO_DB_SLOW_THRESHOLD_MS"]; value != "" {
		if millis, ok := parsePositiveEnvInt(value); ok {
			raw.Database.SlowThresholdMS = millis
		}
	}
	if value := env["KMEMO_DATA_DIR"]; value != "" {
		raw.Paths.DataDir = value
	}
	if value := env["KMEMO_CACHE_DIR"]; value != "" {
		raw.Paths.CacheDir = value
	}
	if value := env["KMEMO_LOGS_DIR"]; value != "" {
		raw.Paths.LogsDir = value
	}
	if value := env["KMEMO_DB_PATH"]; value != "" {
		raw.Database.Path = value
	}
}

func parseEnvBool(value string) bool {
	return value == "1" || value == "true"
}

func parsePositiveEnvInt(value string) (int, bool) {
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}
