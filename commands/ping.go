package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userDiscord *discordgo.User

	if i.Member != nil {
		userDiscord = i.Member.User
	} else {
		userDiscord = i.User
	}

	log.Debug().
		Str("userDiscordId", userDiscord.ID).
		Str("type", "command").
		Str("function", "ping").
		Msg("Command has been triggered")

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🏓 Pong !!!",
		},
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("userDiscordId", userDiscord.ID).
			Str("type", "command").
			Str("function", "ping").
			Msg("An error occurred while responding to the command interaction")
		return
	}
}
