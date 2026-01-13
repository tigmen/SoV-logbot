package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tigmen/SoV-logbot/internal/service"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "sync",
		Description: "Synchronize this discord server with your telegram group",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "group-id",
				Description: "Telegram ID of your group",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "thread-id",
				Description: "Telegram ID for thread in your group",
				Type:        discordgo.ApplicationCommandOptionInteger,
			},
		},
	},
	{
		Name:        "register",
		Description: "Synchronize your telegram username with your discord username",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "username",
				Description: "Your telegram username without '@': \"@example\" - \"example\"",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
}

func parseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	m := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, opt := range options {
		m[opt.Name] = opt
	}
	return m
}

func handleSync(svc *service.Service, s *discordgo.Session, r *discordgo.InteractionCreate, m map[string]*discordgo.ApplicationCommandInteractionDataOption) {
	v, ok := m["thread-id"]
	threadid := 0
	if ok {
		threadid = int(v.IntValue())
	}

	svc.Sync(r.GuildID, m["group-id"].StringValue(), threadid)

	s.InteractionRespond(
		r.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This Guild syncronizated with your Group",
			},
		},
	)
}

func handleRegister(svc *service.Service, s *discordgo.Session, r *discordgo.InteractionCreate, m map[string]*discordgo.ApplicationCommandInteractionDataOption) {
	svc.Register(r.User.Username, m["username"].StringValue())
	s.InteractionRespond(
		r.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Register complete",
			},
		},
	)
}
