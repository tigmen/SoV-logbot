package service

import "github.com/tigmen/SoV-logbot/internal/telegram"

type VoiceChannel struct {
	Name    string
	Members []string
}

type Service struct {
	bot *telegram.Bot
}

func New(bot *telegram.Bot) *Service {
	return &Service{bot: bot}
}

func (s Service) UpdateChannel(map[string]VoiceChannel) {

}
