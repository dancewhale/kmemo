package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func ensureInitialized(paths runtimePaths) error {
	for _, dir := range requiredRuntimeDirs(paths) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(paths.ConfigFile), 0o755); err != nil {
		return fmt.Errorf("create config file directory: %w", err)
	}
	if _, err := os.Stat(paths.ConfigFile); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat config file: %w", err)
	}
	if err := os.WriteFile(paths.ConfigFile, []byte(defaultConfigFileTOML), 0o644); err != nil {
		return fmt.Errorf("write default config file: %w", err)
	}
	return nil
}

func requiredRuntimeDirs(paths runtimePaths) []string {
	return []string{
		paths.ConfigDir,
		paths.DataDir,
		paths.CacheDir,
		paths.LogsDir,
		filepath.Join(paths.DataDir, defaultAssetsDirName),
		filepath.Join(paths.DataDir, defaultImportsDirName),
		filepath.Join(paths.DataDir, defaultIndexDirName),
		filepath.Join(paths.DataDir, defaultExportDirName),
		filepath.Join(paths.DataDir, defaultCardFileDirName),
		filepath.Join(paths.CacheDir, defaultThumbnailsDirName),
		filepath.Join(paths.CacheDir, defaultHTMLCacheDirName),
		filepath.Join(paths.CacheDir, defaultIndexCacheDirName),
	}
}
