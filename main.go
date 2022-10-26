package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/DarkOnion0/ADZTBotV2/commands"
	"github.com/DarkOnion0/ADZTBotV2/config"
	"github.com/DarkOnion0/ADZTBotV2/functions"
	"github.com/rs/zerolog"

	"github.com/bwmarrin/discordgo"

	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var s *discordgo.Session

func init() {
	// enable or not the debug level (default is Info)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *config.Debug == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger()
	// activate the pretty logger for dev purpose only if the debug mode is enabled
	if *config.Debug == "true" {
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().
		Str("type", "main").
		Str("function", "init").
		Msg("Logger is configured!")

	log.Debug().
		Str("type", "main").
		Str("function", "init").
		Msg("Debug mode is enabled!")
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + *config.BotToken)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("type", "main").
			Str("function", "init").
			Msg("An error occurred, invalid bot parameters")
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
	config.Client, errMongo = mongo.NewClient(options.Client().ApplyURI(*config.DBUrl))
	if errMongo != nil {
		log.Fatal().
			Err(errMongo).
			Str("type", "main").
			Str("function", "main").
			Msg("Something bad happen while creating a the MongoDB client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errMongo = config.Client.Connect(ctx)
	if errMongo != nil {
		log.Fatal().
			Err(errMongo).
			Str("type", "main").
			Str("function", "main").
			Msg("Something bad happen while connecting to MongoDB")
	}

	errMongo = config.Client.Ping(ctx, readpref.Primary())
	if errMongo != nil {
		log.Fatal().
			Err(errMongo).
			Str("type", "main").
			Str("function", "main").
			Msg("Something bad happen while pinging the database, is it online ?!")
	}

	log.Info().
		Str("type", "main").
		Str("function", "main").
		Msg("DB is connected !")

	defer func(c *mongo.Client) {
		err := c.Disconnect(ctx)
		if err != nil {
			log.Fatal().
				Err(err).
				Str("type", "main").
				Str("function", "main").
				Msg("An error occurred while closing the bot")
		}
	}(config.Client)

	// Init and update db
	log.Debug().
		Str("type", "main").
		Str("function", "main").
		Msg("Updating the bot")

	err0 := functions.BotUpdater()

	if err0 != nil {
		log.Error().
			Err(err0).
			Str("type", "main").
			Str("function", "main").
			Msg("Something bad happen while updating the bot")
	}

	log.Info().
		Str("type", "main").
		Str("function", "main").
		Msg("Updating the bot finished successfully")

	/*
		Timer job(s)
	*/
	go func() {
		functions.UpdateUserRankingCron()

		for {
			time.Sleep(time.Duration(*config.Timer))
			functions.UpdateUserRankingCron()
		}
	}()

	/*
		Discordgo initialization
	*/

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Info().
			Str("type", "main").
			Str("function", "main").
			Msg("Bot is up!")
	})
	err := s.Open()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("type", "main").
			Str("function", "main").
			Msg("Cannot open the discord session")
	}

	for _, v := range commands.Commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, *config.GuildID, v)
		if err != nil {
			log.Panic().
				Err(err).
				Str("type", "main").
				Str("function", "main").
				Str("command", v.Name).
				Msg("Cannot create command")
		}
		log.Debug().
			Str("type", "main").
			Str("function", "main").
			Str("command", v.Name).
			Msg("Command has been successfully created")
	}

	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {
			log.Fatal().
				Err(err).
				Str("type", "main").
				Str("function", "main").
				Msg("An error occurred while closing the bot")
		}
	}(s)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Info().
		Str("type", "main").
		Str("function", "main").
		Msg("Gracefully shutting down the bot...")
}
