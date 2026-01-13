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

func (s *Service) UpdateChannel(ctx context.Context, guildID string, updates map[string]VoiceChannel) error {
	for key, value := range updates {
		gid, _ok := s.channel[guildID]
		if _ok {
			_, ok := s.message[key]
			if !ok && len(value.Members) > 0 {
				log.LogAttrs(
					ctx, log.LevelDebug,
					"New channel-message",
					log.String("key", key),
					log.String("group id", gid),
					log.String("channel name", value.Name),
					log.String("channel members", fmt.Sprint(value.Members)),
				)

				id, err := s.bot.SendMessage(ctx, gid, fmt.Sprint(value.Members))
				if err != nil {
					return err
				}

				s.message[key] = id
			} else {
				if len(value.Members) > 0 {
					_, err := s.bot.EditMessage(
						ctx, gid, s.message[key],
						fmt.Sprintf("%s: %v", value.Name, value.Members),
					)
					if err != nil {
						return err
					}
				} else {
					_, err := s.bot.DeleteMessage(
						ctx, gid, s.message[key],
					)
					if err != nil {
						return err
					}

				}
			}
		} else {
			return fmt.Errorf("syncerror")
		}
	}

	return nil
}
