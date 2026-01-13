package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/tigmen/SoV-logbot/internal/service"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "sync",
		Description: "Synchronize this discord server with your telegram group",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "Group ID",
				Description: "Telegram ID of your group",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "Thread ID",
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
				Name:        "Username",
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

func handleSync(ctx context.Context, svc *service.Service, s *discordgo.Session, r *discordgo.InteractionCreate, m map[string]*discordgo.ApplicationCommandInteractionDataOption) {
	v, ok := m["Thread ID"]
	threadid := 0
	if ok {
		threadid = int(v.IntValue())
	}

	svc.Sync(r.GuildID, m["Group ID"].StringValue(), threadid)
}

func handleRegister(ctx context.Context, svc *service.Service, s *discordgo.Session, r *discordgo.InteractionCreate, m map[string]*discordgo.ApplicationCommandInteractionDataOption) {
	svc.Register(r.User.Username, m["Username"].StringValue())
}
