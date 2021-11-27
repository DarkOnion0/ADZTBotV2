package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"DarkBotV2/commands"
	"DarkBotV2/config"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	GuildID  = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands.Commands globally")
	BotToken = flag.String("token", "", "Bot access token")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	/*
		MongoBD initialisation
	*/

	var errMongo error
	config.Client, errMongo = mongo.NewClient(options.Client().ApplyURI(config.DBUrl))

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	errMongo = config.Client.Connect(ctx)
	if errMongo != nil {
		log.Fatal(errMongo)
	}

	errMongo = config.Client.Ping(ctx, readpref.Primary())
	if errMongo != nil {
		log.Fatal(errMongo)
	}

	log.Println("DB is connected !")

	defer config.Client.Disconnect(ctx)

	/*
		Discordgo initialisation
	*/

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	for _, v := range commands.Commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {
			fmt.Printf("An error occured whle closing the bot: %s", err)
		}
	}(s)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down...")
}
