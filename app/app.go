package main

import (
	"context"
	"flag"
	"log/slog"
	log "log/slog"
	"os"
	"os/signal"
	"sync"

	"github.com/tigmen/SoV-logbot/internal/discord"
	"github.com/tigmen/SoV-logbot/internal/telegram"
	logger "github.com/tigmen/SoV-logbot/utils/logger"
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

	Ds_token := flag.String("ds-token", "", "Discord bot authentication token")
	Ds_app := flag.String("ds-app", "", "Discord application ID")
	Ds_guild := flag.String("ds-guild", "", "Discord guild ID")

	Tg_token := flag.String("tg-token", "", "Telegram bot authentication token")

	flag.Parse()
	if *Ds_app == "" {
		log.LogAttrs(
			ctx, log.LevelError,
			"Need appId",
		)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	wg := sync.WaitGroup{}
	defer wg.Wait()

	wg.Go(func() {
		err := discord.Start(ctx, Ds_app, Ds_guild, Ds_token)
		if err != nil {
			log.LogAttrs(
				ctx, log.LevelError,
				"Error discord session",
				log.String("Error", err.Error()),
			)
			cancel()
		}
	})

	wg.Go(func() {
		err := telegram.Start(ctx, Tg_token)
		if err != nil {
			log.LogAttrs(
				ctx, log.LevelError,
				"Failed telegram session",
				log.String("Error", err.Error()),
			)
			cancel()
		}
	})
}
