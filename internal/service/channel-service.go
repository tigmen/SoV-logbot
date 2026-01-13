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
	channel map[string]string
	message map[string]int
}

func New(bot *telegram.Bot) *Service {
	return &Service{
		bot:     bot,
		channel: make(map[string]string),
		message: make(map[string]int),
	}
}

func (s *Service) SyncChannel(chname string, chatid string) {
	s.channel[chname] = chatid
}

func (s *Service) UpdateChannel(ctx context.Context, updates map[string]VoiceChannel) {
	for key, value := range updates {
		_, ok := s.message[key]
		if !ok {
			name, _ok := s.channel[key]
			if _ok {
				log.LogAttrs(
					ctx, log.LevelDebug,
					"New channel-message",
					log.String("key", key),
					log.String("channel id", name),
					log.String("channel name", value.Name),
					log.String("channel members", fmt.Sprint(value.Members)),
				)
				s.bot.SendMessage(ctx, value.Name, value.Members[0])
			}
		}
	}
}
