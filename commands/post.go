package commands

import (
	"fmt"

	"github.com/DarkOnion0/ADZTBotV2/config"
	"github.com/DarkOnion0/ADZTBotV2/db"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func post(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//fmt.Printf("The interaction passed, ARG1=%s ARG2=%s", i.ApplicationCommandData().Options[0].StringValue(), i.ApplicationCommandData().Options[1].StringValue())

	var userDiscord *discordgo.User

	if i.Member != nil {
		userDiscord = i.Member.User
	} else {
		userDiscord = i.User
	}
	log.Debug().
		Str("userDiscordId", userDiscord.ID).
		Str("type", "command").
		Str("function", "post").
		Msg("Command has been triggered")

	// check if user exist in the db
	err0, userExists, userDbId := db.CheckUser(userDiscord.ID)

	if err0 != nil {
		log.Error().
			Err(err0).
			Str("userDiscordId", userDiscord.ID).
			Str("type", "command").
			Str("function", "post").
			Msg("Something bad happen while checking the user")

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something bad happen while running the command",
			},
		})
		if err != nil {
			log.Error().
				Err(err).
				Str("userDiscordId", userDiscord.ID).
				Str("type", "command").
				Str("function", "post").
				Msg("An error occurred while responding to the command interaction")
			return
		}

		return
	}

	// execute this block only if the user exist in the db
	if userExists {
		// execute this block if the post does not exist in the database
		postExist, postId := db.Post(userDbId, i.ApplicationCommandData().Options[0].StringValue(), i.ApplicationCommandData().Options[1].StringValue())
		if !postExist {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Your post has been successfully registered",
				},
			})
			if err != nil {
				log.Error().
					Err(err).
					Str("userDiscordId", userDiscord.ID).
					Str("type", "command").
					Str("function", "post").
					Msg("An error occurred while responding to the command interaction")
				return
			}

			message := fmt.Sprintf("**Posted by:**              %s \n**PostID:**                    %s \n**Link:**                         %s", userDiscord.Mention(), postId, i.ApplicationCommandData().Options[1].StringValue())

			// check where the post should be sent to
			if i.ApplicationCommandData().Options[0].StringValue() == "music" {
				_, err2 := s.ChannelMessageSend(*config.ChannelMusic, message)
				log.Info().
					Str("userDiscordId", userDiscord.ID).
					Str("type", "command").
					Str("function", "post").
					Str("postType", "music").
					Str("channel", *config.ChannelMusic).
					Msg("Post has been sent to the corresponding channel")

				if err2 != nil {
					log.Error().
						Err(err2).
						Str("userDiscordId", userDiscord.ID).
						Str("type", "command").
						Str("function", "post").
						Str("postType", "music").
						Msg("An error occurred while sending the post")
					return
				}
			} else {
				_, err2 := s.ChannelMessageSend(*config.ChannelVideo, message)
				log.Info().
					Str("userDiscordId", userDiscord.ID).
					Str("type", "command").
					Str("function", "post").
					Str("postType", "video").
					Str("channel", *config.ChannelVideo).
					Msg("Post has been sent to the corresponding channel")

				if err2 != nil {
					log.Error().
						Err(err2).
						Str("userDiscordId", userDiscord.ID).
						Str("type", "command").
						Str("function", "post").
						Str("postType", "video").
						Msg("An error occurred while sending the post")
					return
				}
			}

		} else {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Sorry but this link already exist in the database 🤷",
				},
			})
			if err != nil {
				log.Error().
					Err(err).
					Str("userDiscordId", userDiscord.ID).
					Str("type", "command").
					Str("function", "post").
					Msg("An error occurred while responding to the command interaction")
				return
			}
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
				Str("function", "post").
				Msg("An error occurred while responding to the command interaction")
		}
	}
}
