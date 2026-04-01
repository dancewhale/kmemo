package config

import "time"

const (
	defaultAppName            = "kmemo"
	defaultProfileName        = "default"
	defaultConfigFileName     = "config.toml"
	defaultDBDriver           = "sqlite"
	defaultDBFileName         = "kmemo.db"
	defaultPythonGRPCAddr     = "127.0.0.1:50051"
	defaultDialTimeoutMS      = 5000
	defaultSlowThresholdMS    = 200
	defaultLogLevel           = "info"
	defaultLogFormat          = "console"
	defaultLogFileName        = "kmemo.log"
	defaultConfigFileTOML     = `[app]
name = "kmemo"
profile = "default"

[server]
python_grpc = "127.0.0.1:50051"
skip_python = true
dial_timeout_ms = 5000

[database]
driver = "sqlite"
path = ""
slow_threshold_ms = 200
repository_debug = false

[logging]
level = "info"
format = "console"
file_enabled = false
file_name = "kmemo.log"

[paths]
data_dir = ""
cache_dir = ""
logs_dir = ""

[feature]
multi_profile = false
`
	defaultAssetsDirName      = "assets"
	defaultImportsDirName     = "imports"
	defaultIndexDirName       = "index"
	defaultExportDirName      = "export"
	defaultCardFileDirName    = "cardfile"
	defaultThumbnailsDirName  = "thumbnails"
	defaultHTMLCacheDirName   = "html"
	defaultIndexCacheDirName  = "index"
)

func defaultRawConfig() rawConfig {
	return rawConfig{
		App: rawAppConfig{
			Name:    defaultAppName,
			Profile: defaultProfileName,
		},
		Server: rawServerConfig{
			PythonGRPCAddr: defaultPythonGRPCAddr,
			SkipPython:     true,
			DialTimeoutMS:  defaultDialTimeoutMS,
		},
		Database: rawDatabaseConfig{
			Driver:          defaultDBDriver,
			SlowThresholdMS: defaultSlowThresholdMS,
		},
		Logging: rawLoggingConfig{
			Level:       defaultLogLevel,
			Format:      defaultLogFormat,
			FileEnabled: false,
			FileName:    defaultLogFileName,
		},
		Feature: rawFeatureConfig{},
	}
}

func defaultConfig() Config {
	return Config{
		AppName:         defaultAppName,
		Profile:         defaultProfileName,
		DBDriver:        defaultDBDriver,
		PythonGRPCAddr:  defaultPythonGRPCAddr,
		DialTimeout:     time.Duration(defaultDialTimeoutMS) * time.Millisecond,
		SkipPython:      true,
		LogLevel:        defaultLogLevel,
		LogFormat:       defaultLogFormat,
		LogFileEnabled:  false,
		LogFileName:     defaultLogFileName,
		RepositoryDebug: false,
		DBSlowThreshold: time.Duration(defaultSlowThresholdMS) * time.Millisecond,
		AssetsDir:       defaultAssetsDirName,
		ImportsDir:      defaultImportsDirName,
		IndexDir:        defaultIndexDirName,
		ExportDir:       defaultExportDirName,
		CardFileDir:     defaultCardFileDirName,
		ThumbsCacheDir:  defaultThumbnailsDirName,
		HTMLCacheDir:    defaultHTMLCacheDirName,
		IndexCacheDir:   defaultIndexCacheDirName,
	}
}
