package commands

import (
	"fmt"
	"strconv"

	"github.com/DarkOnion0/ADZTBotV2/config"
	"github.com/DarkOnion0/ADZTBotV2/db"
	"github.com/rs/zerolog"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/rs/zerolog/log"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name: "ping",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "A little command to test if the bot is really up and running",
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
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"register": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
							Content: "Something bad append while registering your user id, you are already registered",
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
							Content: "Something bad append while registering your user id",
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
		},
		"post": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"vote": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		},
		"stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var userDiscord *discordgo.User

			if i.Member != nil {
				userDiscord = i.Member.User
			} else {
				userDiscord = i.User
			}

			log.Debug().
				Str("userDiscordId", userDiscord.ID).
				Str("type", "command").
				Str("function", "stats").
				Msg("Command has been triggered")

			// check if user exist in the db
			err0, userExists, userDbId := db.CheckUser(userDiscord.ID)

			if err0 != nil {
				log.Error().
					Err(err0).
					Str("userDiscordId", userDiscord.ID).
					Str("type", "command").
					Str("function", "stats").
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
						Str("function", "stats").
						Msg("An error occurred while responding to the command interaction")
					return
				}

				return
			}

			// execute this block only if the user exist in the db
			if userExists && len(i.ApplicationCommandData().Options) != 2 {
				// check which kind of stats should be done, on a user or a post ?
				statsType := i.ApplicationCommandData().Options[0].Type.String()

				if statsType == "String" {
					log.Debug().
						Str("userDiscordId", userDiscord.ID).
						Str("type", "command").
						Str("function", "stats").
						Str("statsType", "post").
						Msg("Running the command for a specific statsType")

					// Get the votes and other information for the selected post
					err, globalVote, userVote, postFetch := db.GetVote(i.ApplicationCommandData().Options[0].StringValue(), userDbId)

					if err != nil {
						log.Error().
							Err(err).
							Str("userDiscordId", userDiscord.ID).
							Str("type", "command").
							Str("function", "vote").
							Str("statsType", "post").
							Str("userVote", userVote).
							Str("postId", i.ApplicationCommandData().Options[0].StringValue()).
							Msg("An error occurred while getting the vote")

						err1 := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fmt.Sprintf("An error occured while executing the GetVote function: %s", err),
							},
						})
						if err1 != nil {
							log.Error().
								Err(err1).
								Str("userDiscordId", userDiscord.ID).
								Str("type", "command").
								Str("function", "stats").
								Str("statsType", "post").
								Msg("An error occurred while responding to the command interaction")
							return
						}
						return
					}

					// Get the information about the user who share the post with db.GetVote function
					err1, userId := db.GetDiscordId(postFetch.User)

					if err1 != nil {
						log.Error().
							Err(err1).
							Str("userDiscordId", userDiscord.ID).
							Str("type", "command").
							Str("function", "vote").
							Str("statsType", "post").
							Str("userVote", userVote).
							Str("globalVote", strconv.Itoa(globalVote)).
							Dict("post", zerolog.Dict().
								Str("id", postFetch.ID.Hex()).
								Str("user", postFetch.User.Hex())).
							Msg("An error occurred while getting the user info of the post author")

						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: fmt.Sprintf("An error occured while executing the GetDiscordID function: %s", err1),
							},
						})
						if err != nil {
							log.Error().
								Err(err).
								Str("userDiscordId", userDiscord.ID).
								Str("type", "command").
								Str("function", "stats").
								Str("statsType", "post").
								Msg("An error occurred while responding to the command interaction")
							return
						}
						return
					}

					// init the discord class for the user who share the post to be able to mention him
					postUser := discordgo.User{ID: userId}

					log.Info().
						Str("userDiscordId", userDiscord.ID).
						Str("type", "command").
						Str("function", "vote").
						Str("statsType", "post").
						Str("userVote", userVote).
						Str("globalVote", strconv.Itoa(globalVote)).
						Dict("post", zerolog.Dict().
							Str("id", postFetch.ID.Hex()).
							Str("user", postFetch.User.Hex())).
						Msg("Post has been sent to the corresponding channel")

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
						log.Error().
							Err(err2).
							Str("userDiscordId", userDiscord.ID).
							Str("type", "command").
							Str("function", "stats").
							Str("statsType", "post").
							Msg("An error occurred while responding to the command interaction")
						return
					}

				} else if statsType == "Mentionable" {
					log.Debug().
						Str("userDiscordId", userDiscord.ID).
						Str("type", "command").
						Str("function", "stats").
						Str("statsType", "user").
						Msg("Running the command for a specific statsType")

					// query the db id of the requested user (the one chosen in the stat command)
					// not the one who has executed the discord command
					err0, userExist, userDbId := db.CheckUser(i.ApplicationCommandData().Options[0].UserValue(s).ID)

					if err0 != nil {
						log.Error().
							Err(err0).
							Str("userDiscordId", userDiscord.ID).
							Str("type", "command").
							Str("function", "stats").
							Str("statsType", "user").
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
								Str("function", "stats").
								Str("statsType", "user").
								Msg("An error occurred while responding to the command interaction")
							return
						}

						return
					}

					if !userExist {
						log.Error().
							Str("userDiscordId", userDiscord.ID).
							Str("type", "command").
							Str("function", "stats").
							Str("statsType", "user").
							Msg("The requested user is not register in the database")

						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "The requested user is not register in the database 😞",
							},
						})

						if err != nil {
							log.Error().
								Err(err).
								Str("userDiscordId", userDiscord.ID).
								Str("type", "command").
								Str("function", "stats").
								Str("statsType", "user").
								Msg("An error occurred while responding to the command interaction")
							return
						}
						return
					}

					// get the user info of the requested user
					err1, userStats := db.GetUserInfo(userDbId)

					if err1 != nil {
						log.Error().
							Err(err1).
							Str("userDiscordId", userDiscord.ID).
							Str("type", "command").
							Str("function", "stats").
							Str("statsType", "user").
							Dict("userStats", zerolog.Dict().
								Str("id", userDbId.Hex())).
							Msg("Something bad append while running the GetUserInfo command (user has no post, doesnt exist...)")

						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "An error append while running the GetUserInfo command (user has no post, doesnt exist...)",
							},
						})

						if err != nil {
							log.Error().
								Err(err).
								Str("userDiscordId", userDiscord.ID).
								Str("type", "command").
								Str("function", "stats").
								Str("statsType", "user").
								Msg("An error occurred while responding to the command interaction")
							return
						}

						return
					}

					// init the user class to be able to mention him
					postUser := discordgo.User{ID: i.ApplicationCommandData().Options[0].UserValue(s).ID}

					log.Info().
						Str("userDiscordId", userDiscord.ID).
						Str("type", "command").
						Str("function", "stats").
						Str("statsType", "user").
						Str("userId", userDbId.Hex()).
						Dict("userStats", zerolog.Dict().
							Str("id", userStats.ID.Hex()).
							Int("globalScore", userStats.GlobalScore)).
						Msg("Post has been sent to the corresponding channel")

					err3 := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title: "📊 User stats",
									Color: 15158332,
									Fields: []*discordgo.MessageEmbedField{
										{Name: "😁 User", Value: postUser.Mention(), Inline: true},
										{Name: "🆔 User database ID", Value: userStats.ID.Hex(), Inline: true},
										{Name: "🆔 User discord ID", Value: i.ApplicationCommandData().Options[0].UserValue(s).ID, Inline: true},
										{Name: "📨 Number of posts", Value: strconv.Itoa(len(userStats.Posts)), Inline: false},
										{Name: "🧮 Global Score", Value: strconv.Itoa(userStats.GlobalScore), Inline: true},
										{Name: "🏆 Ranking", Value: strconv.Itoa(userStats.Rank), Inline: false},
									},
								},
							},
						},
					})

					if err3 != nil {
						log.Error().
							Err(err3).
							Str("userDiscordId", userDiscord.ID).
							Str("type", "command").
							Str("function", "stats").
							Str("statsType", "user").
							Msg("An error occurred while responding to the command interaction")
						return
					}
				} else {
					log.Error().
						Str("userDiscordId", userDiscord.ID).
						Str("type", "command").
						Str("function", "stats").
						Str("statsType", "user").
						Msg("An error occurred while determinating the value of the passed argument")
					return
				}
			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Sorry but your a not register in the bot database or you passed more than one arguments 😭",
					},
				})
				if err != nil {
					log.Error().
						Err(err).
						Str("userDiscordId", userDiscord.ID).
						Str("type", "command").
						Str("function", "stats").
						Msg("An error occurred while responding to the command interaction")
					return
				}
			}
		},
		"delete": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
							Content: "Something bad append deleting the post: this postId does not exist 😞",
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
							Content: "Something bad append deleting the post: only user with the bot admin role or the post author can delete a post 😞",
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
						Content: "Something bad append while deleting the post 😞",
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
		},
	}
)
