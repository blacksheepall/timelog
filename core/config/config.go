package config

import (
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	instance *Config
	once     sync.Once
	mu       sync.Mutex
)

func GetConfig(config_path string) *Config {
	once.Do(func() {
		instance = &Config{}
		resolvedPath := ResolveConfigPath(config_path)
		if err := cleanenv.ReadConfig(resolvedPath, instance); err != nil {
			panic(err)
		}
	})
	return instance
}

func ResolveConfigPath(defaultPath string) string {
	if path := os.Getenv("TIMELOG_CONFIG_PATH"); path != "" {
		return path
	}
	return defaultPath
}

// inject custom config for testing
func SetConfig(cfg *Config) {
	mu.Lock()
	defer mu.Unlock()
	instance = cfg
	// Ensure subsequent GetConfig calls return the injected config instead of
	// re-reading the file. This makes SetConfig/ResetConfig safe when tests
	// run from package directories where the default config.yml is not present.
	once.Do(func() {})
}

// clear the injected config for testing
func ResetConfig() {
	mu.Lock()
	defer mu.Unlock()
	instance = nil
	once = sync.Once{} // reload once to allow loading config again
}
