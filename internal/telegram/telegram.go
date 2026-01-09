package telegram

import (
	"context"
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
		panic(err)
	}

	b.Start(ctx)

	<-ctx.Done()
	return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}
