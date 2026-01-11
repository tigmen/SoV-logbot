package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	log "log/slog"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Start(ctx context.Context, token *string) error {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithCheckInitTimeout(time.Second),
	}

	b, err := bot.New(*token, opts...)
	if err != nil {
		return fmt.Errorf("Opening bot: %w", err)
	}

	user, err := b.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("Getting bot GetMe: %w", err)
	}

	log.LogAttrs(
		ctx, log.LevelInfo,
		"Telegram session opened",
		log.String("Name", user.Username),
	)

	b.Start(ctx)

	<-ctx.Done()

	log.LogAttrs(
		ctx, log.LevelInfo,
		"Telegram session closed",
	)

	return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	upd, err := json.Marshal(update)
	if err == nil {
		log.LogAttrs(ctx, log.LevelDebug,
			"Telegram update",
			log.String("Update", string(upd)),
		)
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
