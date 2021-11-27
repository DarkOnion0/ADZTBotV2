package config

import (
	"flag"

	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	flag.Parse()
}

var (
	Client   *mongo.Client
	DBName   = flag.String("db", "", "The mongodb database name")
	DBUrl    = flag.String("url", "", "The mongodb access url")
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands.Commands globally")
	BotToken = flag.String("token", "", "Bot access token")
)
