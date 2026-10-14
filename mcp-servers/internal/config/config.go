package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServiceName string
	Version     string
}

func Load(defaultServiceName string) Config {
	return Config{
		ServiceName: getenv("SERVICE_NAME", defaultServiceName),
		Version:     getenv("SERVICE_VERSION", "0.1.0"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func requireEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return v, nil
}
