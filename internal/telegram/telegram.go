package telegram

import (
	"context"
	"fmt"
	log "log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Start(ctx context.Context, token *string) error {
	_ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithCheckInitTimeout(time.Second),
	}

	b, err := bot.New(*token, opts...)
	if err != nil {
		return fmt.Errorf("Opening bot: %w", err)
	}

	name, err := b.GetMyName(ctx, &bot.GetMyNameParams{
		LanguageCode: "ru",
	})
	if err != nil {
		return fmt.Errorf("Getting bot name: %w", err)
	}

	log.LogAttrs(
		ctx, log.LevelInfo,
		"Telegram session opened",
		log.String("Name", name.Name),
	)

	b.Start(ctx)

	<-_ctx.Done()

	ok, err := b.Close(ctx)
	if err != nil || !ok {
		return fmt.Errorf("Closing bot: %w", err)
	}

	log.LogAttrs(
		ctx, log.LevelInfo,
		"Telegram session closed",
	)

	return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
