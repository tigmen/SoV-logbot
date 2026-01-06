package main

import (
	"context"
	"flag"
	"log/slog"
	log "log/slog"
	"os"

	"github.com/tigmen/SoV-logbot/internal/discord"
	logger "github.com/tigmen/SoV-logbot/internal/logger"
)

func main() {
	ctx := context.Background()

	opts := logger.HandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := logger.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	Token := flag.String("token", "", "Bot authentication token")
	App := flag.String("app", "", "Application ID")
	Guild := flag.String("guild", "", "Guild ID")

	flag.Parse()
	if *App == "" {
		log.LogAttrs(
			ctx, log.LevelError,
			"Need appId",
		)
	}

	err := discord.Start(ctx, App, Guild, Token)
	if err != nil {
		log.LogAttrs(
			ctx, log.LevelError,
			"Failed start discord bot",
			log.String("Error", err.Error()),
		)
	}
}
