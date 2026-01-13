package main

import (
	"context"
	"flag"
	"log/slog"
	log "log/slog"
	"os"
	"os/signal"
	"sync"

	"github.com/joho/godotenv"
	"github.com/tigmen/SoV-logbot/internal/discord"
	"github.com/tigmen/SoV-logbot/internal/service"
	"github.com/tigmen/SoV-logbot/internal/telegram"
	logger "github.com/tigmen/SoV-logbot/utils/logger"
)

func main() {
	opts := logger.HandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := logger.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.LogAttrs(
			ctx, log.LevelError,
			"Error to load .env",
			log.String("Error", err.Error()),
		)
		os.Exit(1)
	}

	DS_TOKEN := flag.String("ds-token", os.Getenv("DS_TOKEN"), "Discord bot authentication token")
	DS_APP := flag.String("ds-app", os.Getenv("DS_APP"), "Discord application ID")
	DS_GUILD := flag.String("ds-guild", os.Getenv("DS_GUILD"), "Discord guild ID")

	TG_TOKEN := flag.String("tg-token", os.Getenv("TG_TOKEN"), "Telegram bot authentication token")

	flag.Parse()
	if *DS_APP == "" {
		log.LogAttrs(
			ctx, log.LevelError,
			"Need appId",
		)
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	tg, err := telegram.NewBot(ctx, wg, TG_TOKEN)
	if err != nil {
		log.LogAttrs(
			ctx, log.LevelError,
			"Error create telegram session",
			log.String("Error", err.Error()),
		)

		cancel()
	}

	bot, err := discord.NewBot(ctx, DS_APP, DS_GUILD, DS_TOKEN)
	if err != nil {
		log.LogAttrs(
			ctx, log.LevelError,
			"Error create discord session",
			log.String("Error", err.Error()),
		)

		cancel()
	}

	svc := service.New(tg)

	err = bot.Start(ctx, svc)
	if err != nil {
		log.LogAttrs(
			ctx, log.LevelError,
			"Error start discord session",
			log.String("Error", err.Error()),
		)

		cancel()
	}

	defer func() {
		err = bot.Stop(ctx)
		if err != nil {
			log.LogAttrs(
				ctx, log.LevelError,
				"Error stop discord session",
				log.String("Error", err.Error()),
			)

			cancel()
		}
	}()

	svc.SyncChannel("1457173950609100864", "-1003581757738")

	<-ctx.Done()
}
