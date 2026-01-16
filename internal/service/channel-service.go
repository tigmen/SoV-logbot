package service

import (
	"context"
	"fmt"
	log "log/slog"
	"strings"

	"github.com/tigmen/SoV-logbot/internal/telegram"
)

type User struct {
	Username string
	Muted    bool
	Deaf     bool
	Activity string
}

type VoiceChannel struct {
	Name      string
	Members   []User
	messageID int
}

type Guild struct {
	groupID      string
	threadID     int
	voiceChannel map[string]*VoiceChannel
}

type Service struct {
	bot   *telegram.Bot
	guild map[string]*Guild
}

func New(bot *telegram.Bot) *Service {
	return &Service{
		bot:   bot,
		guild: make(map[string]*Guild),
	}
}

func (s *Service) SyncChannel(chname string, chatid string) {
	s.guild[chname] = &Guild{
		groupID:      chname,
		voiceChannel: make(map[string]*VoiceChannel),
	}
}

func (s *Service) UpdateChannel(ctx context.Context, guildID string, updates map[string]VoiceChannel) error {
	for key, value := range updates {
		guild, ok := s.guild[guildID]
		if ok {
			builder := strings.Builder{}
			fmt.Fprintf(&builder, "🌐 %s:\n", value.Name)
			for _, mem := range value.Members {
				if mem.Muted || mem.Deaf {
					fmt.Fprint(&builder, "🔇")
				} else {
					fmt.Fprint(&builder, "🔈")
				}
				fmt.Fprintf(&builder, " %s\n", mem.Username)
				if mem.Activity != "" {
					fmt.Fprintf(&builder, " 🕹%s\n", mem.Activity)
				}
			}

			voice, ok := guild.voiceChannel[key]
			if !ok && len(value.Members) > 0 {
				log.LogAttrs(
					ctx, log.LevelDebug,
					"New channel-message",
					log.String("key", key),
					log.String("group id", guild.groupID),
					log.String("channel name", value.Name),
					log.String("channel members", fmt.Sprint(value.Members)),
				)

				id, err := s.bot.SendMessage(
					ctx, guild.groupID,
					guild.threadID,
					builder.String(),
				)

				if err != nil {
					return err
				}

				voice = &value
				voice.messageID = id
			} else {
				if len(value.Members) > 0 {
					_, err := s.bot.EditMessage(
						ctx, guild.groupID, voice.messageID,
						builder.String(),
					)
					if err != nil {
						return err
					}
				} else {
					_, err := s.bot.DeleteMessage(
						ctx, guild.groupID, voice.messageID,
					)
					if err != nil {
						return err
					}

					delete(guild.voiceChannel, key)
				}
			}
		} else {
			return fmt.Errorf("syncerror")
		}
	}

	return nil
}

func (s *Service) Sync(guildid, chatid string, threadid int) {
	s.guild[guildid] = &Guild{
		groupID:      chatid,
		threadID:     threadid,
		voiceChannel: make(map[string]*VoiceChannel),
	}
}

func (s *Service) Register(dsus, tgus string) {
	//
}
