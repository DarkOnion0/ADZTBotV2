package functions

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/DarkOnion0/ADZTBotV2/config"
	"github.com/DarkOnion0/ADZTBotV2/db"
	"github.com/DarkOnion0/ADZTBotV2/types"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/blang/semver/v4"
)

var (
	v2_0_0 semver.Version
)

func init() {
	v2_0_0 = semver.MustParse("2.0.0")
}

// This function update the database schemas to enable easy migration between breaking versions
func BotUpdater() (err error) {
	// NOTE: Fatal checking is mandatory for every error in the updater function
	// In short: log.Error()... must be replaced by log.Fatal()...

	log.Debug().
		Str("type", "module").
		Str("module", "function").
		Str("function", "botUpdater").
		Msg("Running the function")

	// Fetch the data structure version of the db
	version, err1 := db.FetchVersion()

	// This condition is required to be able to migrate bot instance < 2.0.0
	if os.Getenv("ADZTBOTV2_V1_BEFORE") == "true" && err1 == db.ErrVersionNotFound {
		log.Warn().
			Str("type", "module").
			Str("module", "function").
			Str("function", "botUpdater").
			Msg("Assuming that the bot is running a version < 2.0.0")

		version = semver.MustParse("1.0.0")
	} else if err1 == db.ErrVersionNotFound {
		log.Warn().
			Str("type", "module").
			Str("module", "function").
			Str("function", "botUpdater").
			Msg("Assuming that this is the first time that the bot connect to the DB")

		err4 := db.InitDB()

		if err4 != nil {
			log.Fatal().
				Err(err4).
				Str("type", "module").
				Str("module", "function").
				Str("function", "botUpdater").
				Msg("Something bad happen while creating the base settings of the DB")
		}

		return nil
	} else if err1 != nil {
		log.Fatal().
			Err(err1).
			Str("type", "module").
			Str("module", "function").
			Str("function", "botUpdater").
			Msg("Something bad happen while running the fetchVersion function")

		return err1
	}

	// Skip the update status if the db data structure version is the same as the bot version
	if version.EQ(config.Version) {
		log.Info().
			Str("type", "module").
			Str("module", "function").
			Str("function", "botUpdater").
			Dict("version", zerolog.Dict().
				Str("db", version.String()).
				Str("bot", config.Version.String())).
			Msg("Update finished successfully, nothing was changed")

		return
	} else if version.Major > config.Version.Major {
		log.Fatal().
			Str("type", "module").
			Str("module", "function").
			Str("function", "botUpdater").
			Dict("version", zerolog.Dict().
				Str("db", version.String()).
				Str("bot", config.Version.String())).
			Msg("Bot version is inferior to the db version")

		err = errors.New("bot version is inferior to the db version")

		return
	}

	/*
		/!\ UPDATE RULE(S) /!\
		This is were the core part of the updater starts, please be sure of what you are doing :)
	*/

	if v2_0_0.GT(version) {
		/*
			Add the version field in the database
		*/
		log.Debug().
			Str("type", "module").
			Str("module", "function").
			Str("function", "botUpdater").
			Dict("version", zerolog.Dict().
				Str("db", version.String()).
				Str("bot", config.Version.String()).
				Str("upgradingTo", v2_0_0.String())).
			Msg("Running the function")

		// init the mongodb collection
		internalCollection := config.Client.Database(*config.DBName).Collection("internal")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		log.Debug().
			Str("type", "module").
			Str("module", "function").
			Str("function", "botUpdater").
			Dict("version", zerolog.Dict().
				Str("db", version.String()).
				Str("bot", config.Version.String()).
				Str("upgradingTo", v2_0_0.String())).
			Msg("Adding version to the database")
		_, err5 := internalCollection.InsertOne(ctx, types.DBInfo{Version: v2_0_0})

		if err5 != nil {
			log.Fatal().
				Err(err5).
				Str("type", "module").
				Str("module", "function").
				Str("function", "botUpdater").
				Dict("version", zerolog.Dict().
					Str("db", version.String()).
					Str("bot", config.Version.String()).
					Str("upgradingTo", v2_0_0.String())).
				Msg("Something bad happen while adding version in the database")

			return errors.New("something bad happen while adding version in the database")
		}

		log.Info().
			Str("type", "module").
			Str("module", "function").
			Str("function", "botUpdater").
			Dict("version", zerolog.Dict().
				Str("db", version.String()).
				Str("bot", config.Version.String()).
				Str("upgradingTo", v2_0_0.String())).
			Msg("The version has been added successfully to the database")

		/*
			Add the ranking field in the user profiles
		*/

		// init the mongodb collection
		userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		log.Debug().
			Str("type", "module").
			Str("module", "function").
			Str("function", "botUpdater").
			Dict("version", zerolog.Dict().
				Str("db", version.String()).
				Str("bot", config.Version.String()).
				Str("upgradingTo", v2_0_0.String())).
			Msg("Adding version to the database")
		_, err6 := userInfoCollection.UpdateMany(ctx, bson.D{{Key: "ranking", Value: bson.D{{Key: "$exists", Value: false}}}}, bson.D{{Key: "$set", Value: bson.D{{Key: "ranking", Value: 0}}}})

		if err6 != nil {
			log.Fatal().
				Err(err6).
				Str("type", "module").
				Str("module", "function").
				Str("function", "botUpdater").
				Dict("version", zerolog.Dict().
					Str("db", version.String()).
					Str("bot", config.Version.String()).
					Str("upgradingTo", v2_0_0.String())).
				Msg("Something bad happen while adding ranking in the database")

			return errors.New("something bad happen while adding ranking in the database")
		}

		log.Info().
			Str("type", "module").
			Str("module", "function").
			Str("function", "botUpdater").
			Dict("version", zerolog.Dict().
				Str("db", version.String()).
				Str("bot", config.Version.String()).
				Str("upgradingTo", v2_0_0.String())).
			Msg("The ranking has been added successfully to the user profiles in the database")

		version = v2_0_0
	}

	// Update the data scheme version of the database
	err7 := db.UpdateVersion(config.Version)

	if err7 != nil {
		log.Fatal().
			Err(err7).
			Str("type", "module").
			Str("module", "function").
			Str("function", "botUpdater").
			Dict("version", zerolog.Dict().
				Str("db", version.String()).
				Str("bot", config.Version.String())).
			Msg("Something bad happen while updating version in the database")
	}

	log.Info().
		Err(err7).
		Str("type", "module").
		Str("module", "function").
		Str("function", "botUpdater").
		Dict("version", zerolog.Dict().
			Str("db", version.String()).
			Str("bot", config.Version.String())).
		Msg("Update finished successfully")

	return
}
