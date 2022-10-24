package config

import (
	"flag"

	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	flag.Parse()
}

var (
	Client       *mongo.Client
	DBName       = flag.String("db", "", "The mongodb database name")
	DBUrl        = flag.String("url", "", "The mongodb access url")
	GuildID      = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands.Commands globally")
	BotToken     = flag.String("token", "", "Bot access token")
	ChannelMusic = flag.String("chanm", "", "Discord channel id where the post of the music category will be sent to")
	ChannelVideo = flag.String("chanv", "", "Discord channel id where the post of the video category will be sent to")
	BotAdminRole = flag.String("admin", "0", "The bot administrator discord role ID")
	Debug        = flag.String("debug", "false", "Sets log level to debug true/(false)")
	Timer        = flag.Int64("timer", 3600000000000, "Set a custom timer scheduled for all the background tasks of the bot, run every [X] nanoseconds")
)
