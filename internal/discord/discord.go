package discord

import (
	"context"
	"fmt"
	log "log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type optionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

func parseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) (om optionMap) {
	om = make(optionMap)
	for _, opt := range options {
		om[opt.Name] = opt
	}
	return
}

func interactionAuthor(i *discordgo.Interaction) *discordgo.User {
	if i.Member != nil {
		return i.Member.User
	}
	return i.User
}

func handleEcho(s *discordgo.Session, i *discordgo.InteractionCreate, opts optionMap) {
	builder := new(strings.Builder)
	if v, ok := opts["author"]; ok && v.BoolValue() {
		author := interactionAuthor(i.Interaction)
		builder.WriteString("**" + author.String() + "** says: ")
	}
	builder.WriteString(opts["message"].StringValue())

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: builder.String(),
		},
	})

	if err != nil {
	}
}

func Start(ctx context.Context, App, Guild, Token *string) error {
	_ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log.LogAttrs(
		ctx, log.LevelDebug,
		"Authentication data",
		log.String("App", *App),
		log.String("Guild", *Guild),
		log.String("Token", *Token),
	)

	session, err := discordgo.New("Bot " + *Token)
	if err != nil {
		return fmt.Errorf("Create new session: %w", err)
	}

	defer func() {
		err := session.Close()
		if err != nil {
			log.LogAttrs(
				ctx, log.LevelError,
				"Failed close discord session",
				log.String("Error", err.Error()),
			)
		} else {
			log.LogAttrs(
				ctx, log.LevelInfo,
				"Discord session closed",
			)
		}
	}()

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.LogAttrs(
			ctx, log.LevelInfo,
			"Discord session opened",
			log.String("Name", r.User.Username),
		)
	})
	session.AddHandler(func(s *discordgo.Session, r *discordgo.VoiceStateUpdate) {

		prevChId := "nil"
		if r.BeforeUpdate != nil {
			prevChId = r.BeforeUpdate.ChannelID
		}

		log.LogAttrs(
			ctx, log.LevelInfo,
			"Update user voice state",
			log.String("Username", r.Member.User.Username),
			log.String("Nick", r.Member.Nick),
			log.Bool("Muted", r.SelfMute),
			log.String("Actual channel ID", r.ChannelID),
			log.String("Previous channel ID", prevChId),
		)
	})
	err = session.Open()
	if err != nil {
		return fmt.Errorf("Open session: %w", err)
	}

	<-_ctx.Done()
	return nil
}
