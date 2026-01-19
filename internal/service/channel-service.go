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
	Members   map[string]*User
	messageID int
}

type Guild struct {
	groupID      string
	threadID     int
	voiceChannel map[string]*VoiceChannel
}

type Service struct {
	bot    *telegram.Bot
	guild  map[string]*Guild
	member map[string]string
}

func New(bot *telegram.Bot) *Service {
	return &Service{
		bot:    bot,
		guild:  make(map[string]*Guild),
		member: make(map[string]string),
	}
}

func (s *Service) SyncChannel(chname string, chatid string) {
	s.guild[chname] = &Guild{
		groupID:      chname,
		voiceChannel: make(map[string]*VoiceChannel),
	}
}

func (s *Service) Update(ctx context.Context, guildID string, channelID string, update VoiceChannel) error {
	guild, ok := s.guild[guildID]
	if ok {
		builder := strings.Builder{}
		fmt.Fprintf(&builder, "🌐 %s:\n", update.Name)
		for _, mem := range update.Members {
			if mem.Muted || mem.Deaf {
				fmt.Fprint(&builder, "🔇")
			} else {
				fmt.Fprint(&builder, "🔈")
			}
			fmt.Fprintf(&builder, " %s\n", mem.Username)
			if mem.Activity != "" {
				fmt.Fprintf(&builder, "∟🕹%s\n", mem.Activity)
			}
		}

		voice, ok := guild.voiceChannel[channelID]
		if !ok {
			if len(update.Members) <= 0 {
				return nil
			}

			log.LogAttrs(
				ctx, log.LevelDebug,
				"New channel-message",
				log.String("channelID", channelID),
				log.String("group id", guild.groupID),
				log.String("channel name", update.Name),
				log.String("channel members", fmt.Sprint(update.Members)),
			)

			id, err := s.bot.SendMessage(
				ctx, guild.groupID,
				guild.threadID,
				builder.String(),
			)

			if err != nil {
				return err
			}

			update.messageID = id
			guild.voiceChannel[channelID] = &update
		} else {
			if len(update.Members) > 0 {
				id, err := s.bot.EditMessage(
					ctx, guild.groupID, voice.messageID,
					builder.String(),
				)
				if err != nil {
					return err
				}

				update.messageID = id
				guild.voiceChannel[channelID] = &update
			} else {
				_, err := s.bot.DeleteMessage(
					ctx, guild.groupID, voice.messageID,
				)
				if err != nil {
					return err
				}

				delete(guild.voiceChannel, channelID)
			}
		}
	} else {
		return fmt.Errorf("syncerror")
	}

	return nil
}

func (s *Service) VoiceUpdate(ctx context.Context, guildID string, updates map[string]VoiceChannel) error {
	for key, value := range updates {
		for _, mem := range value.Members {
			activity, ok := s.member[mem.Username]
			if !ok {
				continue
			}

			mem.Activity = activity
		}

		s.Update(ctx, guildID, key, value)
	}

	return nil
}

func (s *Service) ActivityUpdate(ctx context.Context, guildID, channelID, username, activity string) error {
	guild, ok := s.guild[guildID]
	if !ok {
		return nil
	}

	ch, ok := guild.voiceChannel[channelID]
	if !ok {
		return nil
	}

	user, ok := ch.Members[username]
	if !ok {
		return nil
	}

	user.Activity = activity

	err := s.Update(ctx, guildID, channelID, *ch)
	if err != nil {
		return err
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
