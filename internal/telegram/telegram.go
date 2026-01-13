package telegram

import (
	"context"
	"fmt"
	log "log/slog"
	"sync"

	"github.com/go-telegram/bot"
)

type Bot struct {
	bot *bot.Bot
}

func NewBot(ctx context.Context, wg *sync.WaitGroup, token *string) (*Bot, error) {
	opts := []bot.Option{}

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

	wg.Go(func() {
		b.Start(ctx)
	})

	wg.Go(func() {
		<-ctx.Done()

		log.LogAttrs(
			ctx, log.LevelInfo,
			"Telegram session closed",
		)
	})

	return &Bot{bot: b}, nil
}

func (b Bot) SendMessage(ctx context.Context, chatid any, text string) {
	b.bot.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID: chatid,
			Text:   text,
		},
	)
}

func (b Bot) EditMessage(ctx context.Context, chatid any, messageid int, text string) {
	b.bot.EditMessageText(
		ctx,
		&bot.EditMessageTextParams{
			ChatID:    chatid,
			MessageID: messageid,
			Text:      text,
		},
	)
}
