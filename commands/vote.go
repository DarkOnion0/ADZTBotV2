package commands

import (
	"fmt"

	"github.com/DarkOnion0/ADZTBotV2/db"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func vote(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userDiscord *discordgo.User

	if i.Member != nil {
		userDiscord = i.Member.User
	} else {
		userDiscord = i.User
	}
	log.Debug().
		Str("userDiscordId", userDiscord.ID).
		Str("type", "command").
		Str("function", "vote").
		Msg("Command has been triggered")

	// check if user exist in the db
	err0, userExists, userDbId := db.CheckUser(userDiscord.ID)

	if err0 != nil {
		log.Error().
			Err(err0).
			Str("userDiscordId", userDiscord.ID).
			Str("type", "command").
			Str("function", "vote").
			Msg("Something bad append while checking the user")

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something bad append while running the command",
			},
		})
		if err != nil {
			log.Error().
				Err(err).
				Str("userDiscordId", userDiscord.ID).
				Str("type", "command").
				Str("function", "vote").
				Msg("An error occurred while responding to the command interaction")
			return
		}

		return
	}

	// execute this block only if the user exist in the db
	if userExists {

		// Add or change a vote for a post
		err1, postAdded := db.SetVote(i.ApplicationCommandData().Options[0].StringValue(), i.ApplicationCommandData().Options[1].StringValue(), userDbId)

		if err1 != nil {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("An error occured while executing the vote function: %s", err1),
				},
			})
			if err != nil {
				log.Error().
					Err(err).
					Str("userDiscordId", userDiscord.ID).
					Str("type", "command").
					Str("function", "vote").
					Msg("An error occurred while responding to the command interaction")
				return
			}
			log.Error().
				Err(err1).
				Str("userDiscordId", userDiscord.ID).
				Str("type", "command").
				Str("function", "vote").
				Str("userVote", i.ApplicationCommandData().Options[1].StringValue()).
				Str("postId", i.ApplicationCommandData().Options[0].StringValue()).
				Msg("An error occurred while setting the vote")
			return
		}

		var postStatus string

		// choose the correct word to send back in the message
		if postAdded {
			postStatus = "added"
		} else {
			postStatus = "updated"
		}

		log.Info().
			Str("userDiscordId", userDiscord.ID).
			Str("type", "command").
			Str("function", "vote").
			Str("userVote", i.ApplicationCommandData().Options[1].StringValue()).
			Str("postId", i.ApplicationCommandData().Options[0].StringValue()).
			Str("postStatus", postStatus).
			Msg("Post has been sent to the corresponding channel")

		err2 := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Your vote has been successfuly %s for the post `%s`", postStatus, i.ApplicationCommandData().Options[0].StringValue()),
			},
		})

		if err2 != nil {
			log.Error().
				Err(err2).
				Str("userDiscordId", userDiscord.ID).
				Str("type", "command").
				Str("function", "vote").
				Msg("An error occurred while responding to the command interaction")
			return
		}

	} else {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Sorry but your a not register in the bot database 😭",
			},
		})
		if err != nil {
			log.Error().
				Err(err).
				Str("userDiscordId", userDiscord.ID).
				Str("type", "command").
				Str("function", "vote").
				Msg("An error occurred while responding to the command interaction")
			return
		}
	}
}
