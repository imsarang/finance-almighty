package config

import (
	"fmt"
	"strconv"
)

type StockConfig struct {
	Config
	RedisAddr     string
	RedisDB       int
	FinnhubApiKey string
}

func LoadStock() (StockConfig, error) {
	db, err := strconv.Atoi(getenv("REDIS_DB", "0"))
	if err != nil {
		return StockConfig{}, fmt.Errorf("invalid REDIS_DB: %w", err)
	}
	apiKey, err := requireEnv("FINNHUB_API_KEY")
	if err != nil {
		return StockConfig{}, err
	}
	return StockConfig{
		Config:        Load("stock-mcp"),
		RedisAddr:     getenv("REDIS_ADDR", "localhost:6379"),
		RedisDB:       db,
		FinnhubApiKey: apiKey,
	}, nil
}
