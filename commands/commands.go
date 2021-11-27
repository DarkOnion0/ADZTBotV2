package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"ADZTBotV2/db"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name: "basic-command",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},
		{
			Name: "register",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "The first command you must type before starting using the bot !",
		},
	}
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
			if err != nil {
				log.Fatalf("An error occured while creting command handler for the basic-command: %s", err)
				return
			}
		},
		"register": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var userId string

			if i.Member != nil {
				userId = fmt.Sprintf("Thanks %s, your have been successfully registered", i.Member.Mention())
				db.RegisterUser(i.Member.User.ID)
			} else {
				userId = fmt.Sprintf("Thanks %s, your have been successfully registered", i.User.Mention())
				db.RegisterUser(i.User.ID)
				log.Printf("Register command has been triggerd by %s", i.User.ID)
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: userId,
				},
			})
			if err != nil {
				log.Fatalf("An error occured while creting command handler for the basic-command: %s", err)
				return
			}
		},
	}
)
