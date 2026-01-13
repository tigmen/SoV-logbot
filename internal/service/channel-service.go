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

func (s *Service) UpdateChannel(ctx context.Context, guildID string, updates map[string]VoiceChannel) {
	for key, value := range updates {
		name, _ok := s.channel[guildID]
		if _ok {
			_, ok := s.message[key]
			if !ok {
				log.LogAttrs(
					ctx, log.LevelDebug,
					"New channel-message",
					log.String("key", key),
					log.String("group id", name),
					log.String("channel name", value.Name),
					log.String("channel members", fmt.Sprint(value.Members)),
				)
				s.bot.SendMessage(ctx, name, value.Members[0])
			}
		} else {
			//
		}
	}
}
