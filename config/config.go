package config

import (
	"flag"

	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	flag.Parse()
}

var (
	/*
		Constant
	*/

	// The mongodb client configuration
	Client *mongo.Client
	// The software/data structure version
	Version string

	/*
		Command line flags
	*/

	// The mongodb database name
	DBName = flag.String("db", "", "The mongodb database name")
	// The mongodb connection url
	DBUrl = flag.String("url", "", "The mongodb access url")
	// The guild id, not really used or tested
	GuildID = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands.Commands globally")
	// The discord bot token
	BotToken = flag.String("token", "", "Bot access token")
	// The discord channel id where to post musics
	ChannelMusic = flag.String("chanm", "", "Discord channel id where the post of the music category will be sent to")
	// The discord channel id where to post videos
	ChannelVideo = flag.String("chanv", "", "Discord channel id where the post of the video category will be sent to")
	// The bot administrators ids
	BotAdminRole = flag.String("admin", "0", "The bot administrator discord role ID")
	// The global debug varibable, enable shiny logging...
	Debug = flag.String("debug", "false", "Sets log level to debug true/false")
	// The timer in nanoseconds for the background task
	Timer = flag.Int64("timer", 3600000000000, "Set a custom timer scheduled for all the background tasks of the bot, run every [X] nanoseconds")
)
