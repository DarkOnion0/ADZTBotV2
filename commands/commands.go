package commands

import (
	"fmt"
	"log"

	"ADZTBotV2/config"
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
		{
			Name:        "post",
			Description: "Publish a video or a music in a channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "kind",
					Description: "The type of post you want to publish",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "music",
							Value: "music",
						},
						{
							Name:  "video",
							Value: "video",
						},
					},
					Required: true,
				},
				{
					Name:        "url",
					Description: "The link of your publication",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "vote",
			Description: "Like or Dislike the post you want",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "postid",
					Description: "The id of the publication you want",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
				{
					Name:        "vote",
					Description: "Like or Dislike",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "like",
							Value: "+",
						},
						{
							Name:  "dislike",
							Value: "-",
						},
					},
					Required: true,
				},
			},
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
			var registerMessage string

			if i.Member != nil {
				registerMessage = fmt.Sprintf("Thanks %s, your have been successfully registered", i.Member.Mention())
				db.RegisterUser(i.Member.User.ID)
				log.Printf("Register command has been triggerd by %s", i.Member.User.ID)
			} else {
				registerMessage = fmt.Sprintf("Thanks %s, your have been successfully registered", i.User.Mention())
				db.RegisterUser(i.User.ID)
				log.Printf("Register command has been triggerd by %s", i.User.ID)
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: registerMessage,
				},
			})
			if err != nil {
				log.Fatalf("An error occured while responding to the register command interaction: %s", err)
				return
			}
		},
		"post": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			//fmt.Printf("The interaction passed, ARG1=%s ARG2=%s", i.ApplicationCommandData().Options[0].StringValue(), i.ApplicationCommandData().Options[1].StringValue())

			var userDiscord *discordgo.User

			if i.Member != nil {
				userDiscord = i.Member.User
				log.Printf("Post command has been triggerd by %s", i.Member.User.ID)
			} else {
				userDiscord = i.User
				log.Printf("Post command has been triggerd by %s", i.User.ID)
			}

			// check if user exist in the db
			userExists, userDbId := db.CheckUser(userDiscord.ID)

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
						log.Fatalf("An error occured while responding to the register post interaction: %s", err)
					}

					message := fmt.Sprintf("**Posted by:**              %s \n**PostID:**                    %s \n**Link:**                         %s", userDiscord.Mention(), postId, i.ApplicationCommandData().Options[1].StringValue())

					if i.ApplicationCommandData().Options[0].StringValue() == "music" {
						_, err2 := s.ChannelMessageSend(*config.ChannelMusic, message)

						if err2 != nil {
							log.Fatalf("An error occured while sending the music message: %s", err2)
						}
					} else {
						_, err2 := s.ChannelMessageSend(*config.ChannelVideo, message)

						if err2 != nil {
							log.Fatalf("An error occured while sending the video message: %s", err2)
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
						log.Fatalf("An error occured while responding to the register post interaction: %s", err)
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
					log.Fatalf("An error occured while responding to the register post interaction: %s", err)
				}
			}
		},
		"vote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var userDiscord *discordgo.User

			if i.Member != nil {
				userDiscord = i.Member.User
				log.Printf("Post command has been triggerd by %s", i.Member.User.ID)
			} else {
				userDiscord = i.User
				log.Printf("Post command has been triggerd by %s", i.User.ID)
			}

			// check if user exist in the db
			userExists, userDbId := db.CheckUser(userDiscord.ID)

			// execute this block only if the user exist in the db
			if userExists {
				err1, b := db.Vote(i.ApplicationCommandData().Options[0].StringValue(), i.ApplicationCommandData().Options[1].StringValue(), userDbId)

				if err1 != nil {
					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("An error occured while executing the vote function: %s", err1),
						},
					})
					if err != nil {
						log.Fatalf("An error occured while responding to the vote interaction: %s", err)
					}
					log.Fatalf("An error occured while executing the vote function: %s", err1)
				} else {

					var postStatus string

					if b {
						postStatus = "added"
					} else {
						postStatus = "updated"
					}

					err2 := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Your vote has been successfuly %s for the post %s", postStatus, i.ApplicationCommandData().Options[0].StringValue()),
						},
					})

					if err2 != nil {
						log.Fatalf("An error occured while respondinf to the register post interaction: %s", err2)
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
					log.Fatalf("An error occured while responding to the register post interaction: %s", err)
				}
			}

		},
	}
)
