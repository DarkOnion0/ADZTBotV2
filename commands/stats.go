package commands

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/DarkOnion0/ADZTBotV2/db"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func stats(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
				Str("function", "stats").
				Msg("An error occurred while responding to the command interaction")
			return
		}

		return
	}

	// execute this block only if the user exist in the db
	if userExists && len(i.ApplicationCommandData().Options) == 1 {
		// check which kind of stats should be done, on a user or a post ?
		statsType := i.ApplicationCommandData().Options[0].Name

		if statsType == "post" {
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
						Content: fmt.Sprintf("An error occurred while executing the GetVote function: %s", err),
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
						Content: fmt.Sprintf("An error occurred while executing the GetDiscordID function: %s", err1),
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

		} else if statsType == "user" {
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
					Msg("Something bad happen while running the GetUserInfo command (user has no post, doesnt exist...)")

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
		} else if statsType == "ranking" {
			// return the global ranking of the whole serve
			log.Debug().
				Str("type", "module").
				Str("module", "functions").
				Str("function", "updateUserRanking").
				Msg("Fetching all the user in the database")
			err2, userStatsList := db.FetchAllUsers()

			if err2 != nil {
				log.Error().
					Err(err2).
					Str("type", "module").
					Str("module", "functions").
					Str("function", "updateUserRanking").
					Msg("Something bad happen while fetching all the users")

				return
			}
			log.Debug().
				Str("type", "module").
				Str("module", "functions").
				Str("function", "updateUserRanking").
				Msg("All the user was fetched successfully")

			// Limit the scoreboard to 10 users for big guild
			if len(userStatsList) > int(i.ApplicationCommandData().Options[0].IntValue()) {
				userStatsList = userStatsList[0:11]
			}

			// Sort user
			sort.Slice(userStatsList, func(x, y int) bool {
				return userStatsList[x].GlobalScore > userStatsList[y].GlobalScore
			})

			var embedFields []*discordgo.MessageEmbedField

			// Create embed structure for every user
			for _, userStats := range userStatsList {
				emoji := "🔟"

				userStatsDiscord := discordgo.User{ID: userStats.Userid}

				// Add custom emoji for the 3 first users
				switch userStats.Rank {
				case 1:
					emoji = "🥇"
				case 2:
					emoji = "🥈"
				case 3:
					emoji = "🥉"
				case 4:
					emoji = "4️⃣"
				case 5:
					emoji = "5️⃣"
				case 6:
					emoji = "6️⃣"
				case 7:
					emoji = "7️⃣"
				case 8:
					emoji = "8️⃣"
				case 9:
					emoji = "9️⃣"
				}

				embedFields = append(embedFields, &discordgo.MessageEmbedField{Name: fmt.Sprintf("%s - %s pts", emoji, strconv.Itoa(userStats.GlobalScore)), Value: userStatsDiscord.Mention(), Inline: false})
			}

			err3 := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:  fmt.Sprintf("Server Score Ranking - Top %s", strconv.Itoa(int(i.ApplicationCommandData().Options[0].IntValue()))),
							Color:  userDiscord.AccentColor,
							Fields: embedFields,
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
					Str("statsType", "ranking").
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
}
