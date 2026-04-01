package config

// Load initializes config directories/files and returns the effective runtime configuration.
func Load() (Config, error) {
	return loadConfig()
}
