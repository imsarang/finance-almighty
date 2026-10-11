package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServiceName string
	Version string
	RedisAddr string
	RedisDB int
	FinnhubApiKey string
}


func Load() (Config, error) {
	db, err := strconv.Atoi(getenv("REDIS_DB", "0"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid REDIS_DB: %w", err)
	}
	apiKey := os.Getenv("FINNHUB_API_KEY")
	if apiKey == "" {
		return Config{}, errors.New("FINNHUB_API_KEY is required")
	}
	return Config{
		ServiceName:   getenv("SERVICE_NAME", "stock-mcp"),
		Version:       getenv("SERVICE_VERSION", "0.1.0"),
		RedisAddr:     getenv("REDIS_ADDR", "localhost:6379"),
		RedisDB:       db,
		FinnhubApiKey: apiKey,
	}, nil
}
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}