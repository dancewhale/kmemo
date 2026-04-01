package config

import (
	"fmt"
	"os"
	pathpkg "path"
	"path/filepath"
	"runtime"
	"strings"
)

func resolveRuntimePaths() (runtimePaths, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return runtimePaths{}, fmt.Errorf("get working directory: %w", err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return runtimePaths{}, fmt.Errorf("get home directory: %w", err)
	}
	return resolveRuntimePathsFor(runtime.GOOS, pathEnvMap(os.Getenv), cwd, home)
}

func resolveRuntimePathsFor(goos string, env map[string]string, cwd, home string) (runtimePaths, error) {
	configDir := env["KMEMO_CONFIG_DIR"]
	if configDir == "" {
		base, err := defaultConfigRoot(goos, env, home)
		if err != nil {
			return runtimePaths{}, err
		}
		configDir = joinPathForOS(goos, base, defaultAppName)
	}
	configDir = absolutizeForOS(goos, configDir, cwd)

	configFile := env["KMEMO_CONFIG_FILE"]
	if configFile == "" {
		configFile = joinPathForOS(goos, configDir, defaultConfigFileName)
	}
	configFile = absolutizeForOS(goos, configFile, cwd)

	dataDir := env["KMEMO_DATA_DIR"]
	if dataDir == "" {
		base, err := defaultDataRoot(goos, env, home)
		if err != nil {
			return runtimePaths{}, err
		}
		dataDir = joinPathForOS(goos, base, defaultAppName, "data")
	}
	dataDir = absolutizeForOS(goos, dataDir, cwd)

	cacheDir := env["KMEMO_CACHE_DIR"]
	if cacheDir == "" {
		base, err := defaultCacheRoot(goos, env, home)
		if err != nil {
			return runtimePaths{}, err
		}
		if goos == "windows" {
			cacheDir = joinPathForOS(goos, base, defaultAppName, "cache")
		} else {
			cacheDir = joinPathForOS(goos, base, defaultAppName)
		}
	}
	cacheDir = absolutizeForOS(goos, cacheDir, cwd)

	logsDir := env["KMEMO_LOGS_DIR"]
	if logsDir == "" {
		logsDir = joinPathForOS(goos, dataDir, "logs")
	}
	logsDir = absolutizeForOS(goos, logsDir, cwd)

	return runtimePaths{
		ConfigDir:  configDir,
		ConfigFile: configFile,
		DataDir:    dataDir,
		CacheDir:   cacheDir,
		LogsDir:    logsDir,
	}, nil
}

func defaultConfigRoot(goos string, env map[string]string, home string) (string, error) {
	switch goos {
	case "windows":
		if v := firstNonEmpty(env["AppData"], env["APPDATA"]); v != "" {
			return v, nil
		}
		if home == "" {
			return "", fmt.Errorf("APPDATA is empty")
		}
		return joinPathForOS(goos, home, "AppData", "Roaming"), nil
	case "darwin":
		if home == "" {
			return "", fmt.Errorf("HOME is empty")
		}
		return filepath.Join(home, "Library", "Application Support"), nil
	default:
		if v := env["XDG_CONFIG_HOME"]; v != "" {
			return v, nil
		}
		if home == "" {
			return "", fmt.Errorf("HOME is empty")
		}
		return filepath.Join(home, ".config"), nil
	}
}

func defaultDataRoot(goos string, env map[string]string, home string) (string, error) {
	switch goos {
	case "windows":
		if v := firstNonEmpty(env["LocalAppData"], env["LOCALAPPDATA"]); v != "" {
			return v, nil
		}
		if home == "" {
			return "", fmt.Errorf("LOCALAPPDATA is empty")
		}
		return joinPathForOS(goos, home, "AppData", "Local"), nil
	case "darwin":
		root, err := defaultConfigRoot(goos, env, home)
		if err != nil {
			return "", err
		}
		return root, nil
	default:
		if v := env["XDG_DATA_HOME"]; v != "" {
			return v, nil
		}
		if home == "" {
			return "", fmt.Errorf("HOME is empty")
		}
		return filepath.Join(home, ".local", "share"), nil
	}
}

func defaultCacheRoot(goos string, env map[string]string, home string) (string, error) {
	switch goos {
	case "windows":
		return defaultDataRoot(goos, env, home)
	case "darwin":
		if home == "" {
			return "", fmt.Errorf("HOME is empty")
		}
		return filepath.Join(home, "Library", "Caches"), nil
	default:
		if v := env["XDG_CACHE_HOME"]; v != "" {
			return v, nil
		}
		if home == "" {
			return "", fmt.Errorf("HOME is empty")
		}
		return filepath.Join(home, ".cache"), nil
	}
}

func absolutize(path, cwd string) string {
	return absolutizeForOS(runtime.GOOS, path, cwd)
}

func absolutizeForOS(goos, path, cwd string) string {
	if path == "" {
		return ""
	}
	expanded := expandHome(path)
	if isAbsPathForOS(goos, expanded) {
		return cleanPathForOS(goos, expanded)
	}
	return cleanPathForOS(goos, joinPathForOS(goos, cwd, expanded))
}

func joinPathForOS(goos string, elems ...string) string {
	if goos == "windows" {
		parts := make([]string, 0, len(elems))
		for _, elem := range elems {
			if elem != "" {
				parts = append(parts, strings.ReplaceAll(elem, `\`, "/"))
			}
		}
		if len(parts) == 0 {
			return ""
		}
		return strings.ReplaceAll(pathpkg.Join(parts...), "/", `\`)
	}
	return filepath.Join(elems...)
}

func cleanPathForOS(goos, value string) string {
	if goos == "windows" {
		return strings.ReplaceAll(pathpkg.Clean(strings.ReplaceAll(value, `\`, "/")), "/", `\`)
	}
	return filepath.Clean(value)
}

func isAbsPathForOS(goos, value string) bool {
	if goos == "windows" {
		if len(value) >= 3 && value[1] == ':' && (value[2] == '\\' || value[2] == '/') {
			return true
		}
		return strings.HasPrefix(value, `\\`)
	}
	return filepath.IsAbs(value)
}

func expandHome(path string) string {
	if path == "" || path[0] != '~' {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return path
	}
	if path == "~" {
		return home
	}
	if len(path) > 1 && (path[1] == '/' || path[1] == '\\') {
		return filepath.Join(home, path[2:])
	}
	return path
}

func pathEnvMap(getenv func(string) string) map[string]string {
	keys := []string{
		"KMEMO_CONFIG_DIR",
		"KMEMO_CONFIG_FILE",
		"KMEMO_DATA_DIR",
		"KMEMO_CACHE_DIR",
		"KMEMO_LOGS_DIR",
		"XDG_CONFIG_HOME",
		"XDG_DATA_HOME",
		"XDG_CACHE_HOME",
		"APPDATA",
		"AppData",
		"LOCALAPPDATA",
		"LocalAppData",
	}
	m := make(map[string]string, len(keys))
	for _, key := range keys {
		m[key] = getenv(key)
	}
	return m
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
