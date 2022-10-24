package commands

import (
	"github.com/bwmarrin/discordgo"
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
					Description: "Chose a post to display the stat from",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    false,
				},
				{
					Name:        "ranking",
					Description: "Display the global server ranking of the selected number of user",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Top 3",
							Value: 3,
						},
						{
							Name:  "Top 5",
							Value: 5,
						},
						{
							Name:  "Top 10",
							Value: 10,
						},
					},
					Required: false,
				},
			},
		},
	}
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ping":     ping,
		"register": register,
		"post":     post,
		"vote":     vote,
		"stats":    stats,
		"delete":   delete,
	}
)
