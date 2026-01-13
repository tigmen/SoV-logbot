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

func (b Bot) SendMessage(ctx context.Context, chatid string, text string) (int, error) {
	message, err := b.bot.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID: chatid,
			Text:   text,
		},
	)
	if err != nil {
		return -1, err
	}

	return message.ID, nil
}

func (b Bot) EditMessage(ctx context.Context, chatid string, messageid int, text string) (int, error) {
	message, err := b.bot.EditMessageText(
		ctx,
		&bot.EditMessageTextParams{
			ChatID:    chatid,
			MessageID: messageid,
			Text:      text,
		},
	)
	if err != nil {
		return -1, err
	}

	return message.ID, nil
}

func (b Bot) DeleteMessage(ctx context.Context, chatid string, messageid int) (bool, error) {
	ok, err := b.bot.DeleteMessage(
		ctx,
		&bot.DeleteMessageParams{
			ChatID:    chatid,
			MessageID: messageid,
		},
	)
	if err != nil {
		return false, err
	}

	return ok, nil
}
