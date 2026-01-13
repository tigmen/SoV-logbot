package discord

import (
	"context"
	"fmt"
	log "log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/tigmen/SoV-logbot/internal/service"
)

type Bot struct {
	session *discordgo.Session
}

func NewBot(ctx context.Context, App, Guild, Token *string) (*Bot, error) {
	log.LogAttrs(
		ctx, log.LevelInfo,
		"Authentication data",
		log.String("App", *App),
		log.String("Guild", *Guild),
		log.String("Token", *Token),
	)

	session, err := discordgo.New("Bot " + *Token)
	if err != nil {
		return nil, fmt.Errorf("Create new session: %w", err)
	}

	session.Identify.Intents = discordgo.MakeIntent(
		discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates,
	)

	return &Bot{session: session}, nil
}

func (b Bot) Start(ctx context.Context, svc *service.Service) error {
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.LogAttrs(
			ctx, log.LevelInfo,
			"Discord session opened",
			log.String("Name", r.User.Username),
		)
	})

	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.VoiceStateUpdate) {
		log.LogAttrs(
			ctx, log.LevelDebug,
			"Voice State Update",
			log.String("Username", r.Member.User.Username),
		)
		guild, err := s.State.Guild(r.GuildID)
		if err != nil {
			log.LogAttrs(
				ctx, log.LevelError,
				"Error getting guild",
				log.String("Error", err.Error()),
			)
		}

		voice_channel := make(map[string]service.VoiceChannel)

		if r.BeforeUpdate != nil {
			v_Ch, err := s.State.Channel(r.BeforeUpdate.ChannelID)
			if err != nil {
				log.LogAttrs(
					ctx, log.LevelError,
					"Error getting voice channel",
					log.String("Error", err.Error()),
				)
			}

			voice_channel[r.BeforeUpdate.ChannelID] = service.VoiceChannel{
				Name:    v_Ch.Name,
				Members: make([]string, 0),
			}
		}

		for _, vs := range guild.VoiceStates {
			if vs.ChannelID != "" {
				users, ok := voice_channel[vs.ChannelID]
				if !ok {
					V_Ch, err := s.State.Channel(vs.ChannelID)
					if err != nil {
						log.LogAttrs(
							ctx, log.LevelError,
							"Error getting voice channel",
							log.String("Error", err.Error()),
						)
					}

					users = service.VoiceChannel{
						Name:    V_Ch.Name,
						Members: make([]string, 0),
					}
				}

				users.Members = append(users.Members, vs.Member.User.Username)
				voice_channel[vs.ChannelID] = users
			}

			log.LogAttrs(
				ctx, log.LevelDebug,
				"Debug guild voice states",
				log.String("GuildID", guild.ID),
				log.String("Username", vs.Member.User.Username),
				log.String("Voice channel ID", vs.ChannelID),
				log.String("Voice channel name", voice_channel[vs.ChannelID].Name),
			)
		}

		err = svc.UpdateChannel(ctx, guild.ID, voice_channel)
		if err != nil {
			log.LogAttrs(
				ctx, log.LevelError,
				"Error sending telegram message",
				log.String("Error", err.Error()),
			)
		}
	})

	err := b.session.Open()
	if err != nil {
		return fmt.Errorf("Open session: %w", err)
	}

	return nil
}

func (b Bot) Stop(ctx context.Context) error {

	err := b.session.Close()
	if err != nil {
		return err
	} else {
		log.LogAttrs(
			ctx, log.LevelInfo,
			"Discord session closed",
		)
	}

	return nil
}
