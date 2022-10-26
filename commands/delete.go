package commands

import (
	"fmt"

	"github.com/DarkOnion0/ADZTBotV2/config"
	"github.com/DarkOnion0/ADZTBotV2/db"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func delete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		log.Error().
			Str("userDiscordId", i.User.ID).
			Str("type", "command").
			Str("function", "delete").
			Msg("An error occurred while running the command, the user has not executed the command in a server")

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Sorry but this command is only available in a discord server not here 😞",
			},
		})

		if err != nil {
			log.Error().
				Err(err).
				Str("userDiscordId", i.User.ID).
				Str("type", "command").
				Str("function", "delete").
				Msg("An error occurred while responding to the command interaction")
			return
		}

		return
	}

	// check if user exist in the db and retrieve is db id
	err0, userExists, userDbId := db.CheckUser(i.Member.User.ID)

	if err0 != nil {
		log.Error().
			Err(err0).
			Str("userDiscordId", i.Member.User.ID).
			Str("type", "command").
			Str("function", "delete").
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
				Str("userDiscordId", i.Member.User.ID).
				Str("type", "command").
				Str("function", "delete").
				Msg("An error occurred while responding to the command interaction")
			return
		}

		return
	}

	if !userExists {
		log.Error().
			Str("userDiscordId", i.Member.User.ID).
			Str("type", "command").
			Str("function", "delete").
			Msg("The user is not register in the database")

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "The requested user is not register in the database 😞",
			},
		})

		if err != nil {
			log.Error().
				Err(err).
				Str("userDiscordId", i.Member.User.ID).
				Str("type", "command").
				Str("function", "delete").
				Msg("An error occurred while responding to the command interaction")
			return
		}

		return
	}

	isBotAdmin := false

	log.Debug().
		Str("userDiscordId", i.Member.User.ID).
		Str("type", "command").
		Str("function", "delete").
		Str("userId", userDbId.Hex()).
		Bool("isBotAdmin", isBotAdmin).
		Msg("Checking if the user is a bot admin")

	if *config.BotAdminRole != "0" {
		for index := 0; index < len(i.Member.Roles); index++ {
			if i.Member.Roles[index] == *config.BotAdminRole {
				isBotAdmin = true

				log.Debug().
					Str("userDiscordId", i.Member.User.ID).
					Str("type", "command").
					Str("function", "delete").
					Bool("isBotAdmin", isBotAdmin).
					Msg("The user is a bot admin")
			}
		}
	}

	// convert the command argument from a string to a primitive.ObjectID
	postId, _ := primitive.ObjectIDFromHex(i.ApplicationCommandData().Options[0].StringValue())

	// delete the post according to the user and the post db id
	err1 := db.DeletePost(postId, userDbId, isBotAdmin)

	// check the error message from the DeletePost function and send back message the discord
	if err1 != nil {
		log.Error().
			Err(err1).
			Str("userDiscordId", i.Member.User.ID).
			Str("type", "command").
			Str("function", "delete").
			Str("userId", userDbId.Hex()).
			Bool("isBotAdmin", isBotAdmin).
			Str("postId", i.ApplicationCommandData().Options[0].StringValue()).
			Msg("Somethings bad append while deleting the post")

		switch err1 {
		case db.ErrNoDocument:
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Something bad happen deleting the post: this postId does not exist 😞",
				},
			})

			if err != nil {
				log.Error().
					Err(err).
					Str("userDiscordId", i.Member.User.ID).
					Str("type", "command").
					Str("function", "delete").
					Msg("An error occurred while responding to the command interaction")
				return
			}
			return
		case db.ErrWrongUserDbId:
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Something bad happen deleting the post: only user with the bot admin role or the post author can delete a post 😞",
				},
			})

			if err != nil {
				log.Error().
					Err(err).
					Str("userDiscordId", i.Member.User.ID).
					Str("type", "command").
					Str("function", "delete").
					Msg("An error occurred while responding to the command interaction")
				return
			}

			return
		}

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something bad happen while deleting the post 😞",
			},
		})

		if err != nil {
			log.Error().
				Err(err).
				Str("userDiscordId", i.Member.User.ID).
				Str("type", "command").
				Str("function", "delete").
				Msg("An error occurred while responding to the command interaction")
			return
		}

		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Post `%s` successfully deleted 👍", i.ApplicationCommandData().Options[0].StringValue()),
		},
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("userDiscordId", i.Member.User.ID).
			Str("type", "command").
			Str("function", "delete").
			Msg("An error occurred while responding to the command interaction")
		return
	}
}
