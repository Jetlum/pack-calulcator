package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration
type Config struct {
	Port      string
	PackSizes []int
}

// Load reads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		PackSizes: []int{250, 500, 1000, 2000, 5000}, // default
	}

	// Parse custom pack sizes from environment
	if packSizesEnv := os.Getenv("PACK_SIZES"); packSizesEnv != "" {
		sizes := []int{}
		for _, s := range strings.Split(packSizesEnv, ",") {
			if size, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				sizes = append(sizes, size)
			}
		}
		if len(sizes) > 0 {
			cfg.PackSizes = sizes
		}
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
