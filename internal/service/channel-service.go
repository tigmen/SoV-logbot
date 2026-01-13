package service

import (
	"context"
	"fmt"
	log "log/slog"

	"github.com/tigmen/SoV-logbot/internal/telegram"
)

type VoiceChannel struct {
	Name    string
	Members []string
}

type Service struct {
	bot     *telegram.Bot
	channel map[string]any
	message map[string]int
}

func New(bot *telegram.Bot) *Service {
	return &Service{
		bot:     bot,
		channel: make(map[string]any),
		message: make(map[string]int),
	}
}

func (s *Service) SyncChannel(chname string, chatid any) {
	s.channel[chname] = chatid
}

func (s *Service) UpdateChannel(ctx context.Context, updates map[string]VoiceChannel) {
	log.LogAttrs(
		ctx, log.LevelDebug,
		"Update Channel",
		log.String("channel sync", fmt.Sprint(s.channel)),
	)

	for key, value := range updates {
		_, ok := s.message[key]
		if !ok {
			s.bot.SendMessage(ctx, s.channel[key], value.Name)
		}
	}
}
