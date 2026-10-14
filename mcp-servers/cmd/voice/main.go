package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/gordonklaus/portaudio"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/imsarang/finance-almighty/mcp-servers/internal/config"
	"github.com/imsarang/finance-almighty/mcp-servers/internal/tools"
	"github.com/imsarang/finance-almighty/mcp-servers/internal/voice"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	if err := run(context.Background()); err != nil {
		log.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadVoice()
	if err != nil {
		return err
	}

	model, err := whisper.New(cfg.WhisperModelPath)
	if err != nil {
		return fmt.Errorf("failed to load whisper weights: %w", err)
	}
	defer model.Close()

	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf("failed to init system audio: %w", err)
	}
	defer portaudio.Terminate()

	server := tools.ServerInit(
		cfg.ServiceName,
		cfg.Version,
		tools.VoiceToolInit(voice.VoiceService(model)),
	)
	return server.Run(ctx, &mcp.StdioTransport{})
}
