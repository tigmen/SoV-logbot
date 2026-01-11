package telegram

import (
	"context"
	"fmt"
	log "log/slog"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	chatID = 1417379260
)

func Start(ctx context.Context, wg *sync.WaitGroup, token *string) (func(string), error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithCheckInitTimeout(time.Second),
	}

	b, err := bot.New(*token, opts...)
	if err != nil {
		return nil, fmt.Errorf("Opening bot: %w", err)
	}

	user, err := b.GetMe(ctx)
	if err != nil {
		return nil, fmt.Errorf("Getting bot GetMe: %w", err)
	}

	log.LogAttrs(
		ctx, log.LevelInfo,
		"Telegram session opened",
		log.String("Name", user.Username),
	)

	b.Start(ctx)

	wg.Go(func() {
		<-ctx.Done()

		log.LogAttrs(
			ctx, log.LevelInfo,
			"Telegram session closed",
		)
	})

	send := func(text string) {
		b.SendMessage(
			ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   text,
			},
		)
	}

	return send, nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
}
