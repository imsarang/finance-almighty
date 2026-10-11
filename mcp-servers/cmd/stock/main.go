package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/redis/go-redis/v9"

	"github.com/imsarang/finance-almighty/mcp-servers/internal/config"
	"github.com/imsarang/finance-almighty/mcp-servers/internal/fundamentals"
	"github.com/imsarang/finance-almighty/mcp-servers/internal/queue"
	"github.com/imsarang/finance-almighty/mcp-servers/internal/stocks"
	"github.com/imsarang/finance-almighty/mcp-servers/internal/tools"
)

func main(){
	log:=slog.New(slog.NewTextHandler(os.Stderr, nil))
	if err := run(context.Background()); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	
	// load configs from env
	cfg, err:=config.Load()
	if err!=nil {
		return err
	}

	// connect redis
	rdb:= redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, DB: cfg.RedisDB})
	defer rdb.Close()

	if err:=rdb.Ping(ctx).Err(); err != nil {
		slog.Error("redis could not be connected")
		return err
	}

	// outbound adapters
	publisher := queue.NewCeleryPublisher(rdb)
	results := queue.NewRedisResultReader(rdb)

	// clients
	httpClient := &http.Client{Timeout: 5 * time.Second}
	fundamentalsProvider := fundamentals.FinnhubInit(httpClient, cfg.FinnhubApiKey)

	// services
	stockSvc := stocks.StockService(publisher, fundamentalsProvider)

	// MCP tools
	server := tools.ServerInit(cfg.ServiceName, cfg.Version,
		tools.StockToolsConstructor(stockSvc),
		tools.JobsToolsInit(results),
	)

	return server.Run(ctx, &mcp.StdioTransport{})
}