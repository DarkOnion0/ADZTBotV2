package commands

import (
	"fmt"
	"log"
	"strconv"

	"github.com/DarkOnion0/ADZTBotV2/config"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/DarkOnion0/ADZTBotV2/db"
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
			Name:        "delete",
			Description: "Delete a post",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "postid",
					Description: "The id of the publication you want",
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
		{
			Name:        "stats",
			Description: "display some infos and stats about a user or a post",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Description: "Chose a user to display the stats from",
					Type:        discordgo.ApplicationCommandOptionMentionable,
					//Autocomplete: true,
					Required: false,
				},
				{
					Name:        "post",
					Description: "Chose a post to display the stat from ",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
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

					// check where the post should be sent to
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
						log.Fatalf("An error occured while responding to the vote interaction: %s", err)
					}
					log.Fatalf("An error occured while executing the vote function: %s", err1)
				} else {

					var postStatus string

					// choose the correct word to send back in the message
					if postAdded {
						postStatus = "added"
					} else {
						postStatus = "updated"
					}

					err2 := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Your vote has been successfuly %s for the post `%s`", postStatus, i.ApplicationCommandData().Options[0].StringValue()),
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
		"stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			fmt.Println(userExists, len(i.ApplicationCommandData().Options))

			// execute this block only if the user exist in the db
			if userExists && len(i.ApplicationCommandData().Options) != 2 {
				// check which kind of stats should be done, on a user or a post ?
				statsType := i.ApplicationCommandData().Options[0].Type.String()

				if statsType == "String" {
					log.Println("Running the stats command for a post")
					// Get the votes and other information for a post
					err, globalVote, userVote, postFetch := db.GetVote(i.ApplicationCommandData().Options[0].StringValue(), userDbId)

					if err != nil {
						err1 := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fmt.Sprintf("An error occured while executing the vote function: %s", err),
							},
						})
						if err1 != nil {
							log.Fatalf("An error occured while responding to the vote interaction: %s", err1)
						}
						log.Fatalf("An error occured while executing the vote function: %s", err)
					}

					// Get the information about the user who share the post with db.GetVote function
					err1, userId := db.GetDiscordId(postFetch.User)

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
					}

					// init the discord class for the user who share the post to be able to mention him
					postUser := discordgo.User{ID: userId}

					err2 := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title: "📈 Post Stats",
									Color: 16705372,
									Fields: []*discordgo.MessageEmbedField{
										{Name: "✍️ Author", Value: postUser.Mention(), Inline: false},
										{Name: "📨 Post", Value: postFetch.Url, Inline: false},
										{Name: "🏆 Vote", Value: strconv.Itoa(globalVote), Inline: false},
										{Name: "🎓 User Vote", Value: userVote, Inline: false},
									},
								},
							},
						},
					})

					if err2 != nil {
						log.Fatalf("An error occured while respondinf to the register post interaction: %s", err2)
					}

				} else if statsType == "Mentionable" {
					log.Println("Running the stats command for a user")
					// query the db id of the requested user (the one chosen in the stat command)
					userExist, userDbId := db.CheckUser(i.ApplicationCommandData().Options[0].UserValue(s).ID)

					if !userExist {
						log.Printf("The requested user is not register in the database")

						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "The requested user is not register in the database 😞",
							},
						})

						if err != nil {
							log.Fatalf("An error occured while sending back the checkuser error message to discord: %s", err)
						}
						return
					}

					// get the user info of the requested user
					err1, userStats := db.GetUserInfo(userDbId)

					if err1 != 0 {
						log.Printf("An error append while running the GetUserInfo command (user has no post, doesnt exist...): %s", strconv.Itoa(err1))

						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "An error append while running the GetUserInfo command (user has no post, doesnt exist...)",
							},
						})

						if err != nil {
							log.Fatalf("An error occured while sending back the GetUserInfo error message to discord: %s", err)
						}
					}

					// init the user class to be able to mention him
					postUser := discordgo.User{ID: i.ApplicationCommandData().Options[0].UserValue(s).ID}

					err3 := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title: "📊 User stats",
									Color: 15158332,
									Fields: []*discordgo.MessageEmbedField{
										{Name: "😁 User", Value: postUser.Mention(), Inline: false},
										{Name: "🆔 User database ID", Value: userStats.ID.Hex(), Inline: true},
										{Name: "🆔 User discord ID", Value: i.ApplicationCommandData().Options[0].UserValue(s).ID, Inline: true},
										{Name: "📨 Number of posts", Value: strconv.Itoa(len(userStats.Posts)), Inline: false},
										{Name: "🏆 Global Score", Value: strconv.Itoa(userStats.GlobalScore), Inline: false},
									},
								},
							},
						},
					})

					if err3 != nil {
						log.Fatalf("An error occured while sending back the user stats: %s", err3)
					}
				} else {
					log.Fatalf("Something bad append while determinating the value of the passed argument")
				}
			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry but your a not register in the bot database or you passed more than one arguments 😭",
					},
				})
				if err != nil {
					log.Fatalf("An error occured while responding to the register post interaction: %s", err)
				}
			}
		},
		"delete": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var userDiscord *discordgo.User

			if i.Member != nil {
				userDiscord = i.Member.User
				log.Printf("Post command has been triggerd by %s", i.Member.User.ID)
			} else {
				userDiscord = i.User
				log.Printf("Post command has been triggerd by %s", i.User.ID)
			}

			// check if user exist in the db and retrieve is db id
			userExists, userDbId := db.CheckUser(userDiscord.ID)

			if !userExists {
				log.Printf("The requested user %s is not register in the database", userDiscord.ID)

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "The requested user is not register in the database 😞",
					},
				})

				if err != nil {
					log.Fatalf("An error occured while sending back the checkuser error message to discord: %s", err)
				}
				return
			}

			// convert the command argument from a string to a primitive.ObjectID
			postId, _ := primitive.ObjectIDFromHex(i.ApplicationCommandData().Options[0].StringValue())

			// delete the post according to the user and the post db id
			err1 := db.DeletePost(postId, userDbId)

			// check the error message from the DeletePost function and send back message the discord
			if err1 != nil {
				log.Printf("Somethings bad append while deleting a post: %s", err1)

				switch err1 {
				case db.ErrNoDocument:
					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Something bad append deleting the post: this postId does not exist 😞",
						},
					})

					if err != nil {
						log.Fatalf("An error occured while sending back the error message to discord: %s", err)
					}
					return
				case db.ErrWrongUserDbId:
					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Something bad append deleting the post: only user with the bot admin role or the post author can delete a post 😞",
						},
					})

					if err != nil {
						log.Fatalf("An error occured while sending back the error message to discord: %s", err)
					}
					return
				}

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Something bad append while deleting the post 😞",
					},
				})

				if err != nil {
					log.Fatalf("An error occured while sending back the checkuser error message to discord: %s", err)
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
				log.Fatalf("An error occured while responding to the vote interaction: %s", err)
			}
		},
	}
)
