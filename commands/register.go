package commands

import (
	"fmt"

	"github.com/DarkOnion0/ADZTBotV2/db"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func register(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var registerMessage string
	var userDiscord *discordgo.User

	if i.Member != nil {
		userDiscord = i.Member.User
	} else {
		userDiscord = i.User
	}

	log.Debug().
		Str("userDiscordId", userDiscord.ID).
		Str("type", "command").
		Str("function", "register").
		Msg("Command has been triggered")

	registerMessage = fmt.Sprintf("Thanks %s, your have been successfully registered", userDiscord.Mention())
	err1 := db.RegisterUser(userDiscord.ID)

	if err1 != nil {
		log.Error().
			Err(err1).
			Str("userDiscordId", userDiscord.ID).
			Str("type", "command").
			Str("function", "register").
			Msg("An error occurred while registering the user")

		switch err1 {
		case db.ErrUserAlreadyRegistered:
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Something bad happen while registering your user id, you are already registered",
				},
			})
			if err != nil {
				log.Error().
					Err(err).
					Str("userDiscordId", userDiscord.ID).
					Str("type", "command").
					Str("function", "register").
					Msg("An error occurred while responding to the command interaction")
				return
			}
		default:
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Something bad happen while registering your user id",
				},
			})
			if err != nil {
				log.Error().
					Err(err).
					Str("userDiscordId", userDiscord.ID).
					Str("type", "command").
					Str("function", "register").
					Msg("An error occurred while responding to the command interaction")
				return
			}
		}

		return
	}

	err2 := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: registerMessage,
		},
	})
	if err2 != nil {
		log.Error().
			Err(err2).
			Str("userDiscordId", userDiscord.ID).
			Str("type", "command").
			Str("function", "register").
			Msg("An error occurred while responding to the command interaction")
		return
	}
}
