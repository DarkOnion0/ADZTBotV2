package config

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Client *mongo.Client
	DBName = os.Getenv("ADZTBotV2_DB_NAME")
	DBUrl  = os.Getenv("ADZTBotV2_DB_URL")
)
