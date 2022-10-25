package db

import (
	"context"
	"errors"
	"time"

	"github.com/DarkOnion0/ADZTBotV2/config"
	"github.com/DarkOnion0/ADZTBotV2/types"
	"github.com/blang/semver/v4"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// The document containing the version was not found in the database
var ErrVersionNotFound error = errors.New("the version was not found in the database")

// This function fetch the version in the database and return it
func FetchVersion() (version semver.Version, err error) {
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "fetchVersion").
		Msg("Running the function")

	internalCollection := config.Client.Database(*config.DBName).Collection("internal")

	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "fetchVersion").
		Msg("Fetching the version information")

	var dbInfo types.DBInfo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err1 := internalCollection.FindOne(ctx, bson.D{{Key: "version", Value: bson.D{{Key: "$exists", Value: true}}}}).Decode(&dbInfo)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			log.Warn().
				Err(err1).
				Str("type", "module").
				Str("module", "db").
				Str("function", "fetchVersion").
				Msg("The version was not found in the database")

			err = ErrVersionNotFound

			return
		}

		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "db").
			Str("function", "fetchVersion").
			Msg("Something bad append while fetching the version")

		err = errors.New("an error occurred while fetching the version")

		return
	}

	version = dbInfo.Version

	log.Info().
		Str("type", "module").
		Str("module", "db").
		Str("function", "fetchVersion").
		Str("version", version.String()).
		Msg("The function finished successfully")

	return
}

// This functions update the DB version
func UpdateVersion(version semver.Version) (err error) {
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "updateVersion").
		Msg("Running the function")

	// init the mongodb collection
	internalCollection := config.Client.Database(*config.DBName).Collection("internal")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "updateVersion").
		Msg("Updating version in the database")
	_, err5 := internalCollection.UpdateOne(ctx, bson.D{{Key: "version", Value: bson.D{{Key: "$exists", Value: true}}}}, bson.D{{Key: "$set", Value: bson.D{{Key: "version", Value: version}}}})

	if err5 != nil {
		log.Error().
			Str("type", "module").
			Str("module", "db").
			Str("function", "updateVersion").
			Msg("Something bad append while updating version in the database")

		return errors.New("something bad append while updating version in the database")
	}

	log.Info().
		Str("type", "module").
		Str("module", "db").
		Str("function", "updateVersion").
		Msg("The DB version has been updated successfully")

	return
}

// Create some mandatory settings of the database
func InitDB() (err error) {
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "initDB").
		Msg("Running the function")

	// init the mongodb collection
	internalCollection := config.Client.Database(*config.DBName).Collection("internal")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "initDB").
		Msg("Adding default info to the database")
	_, err5 := internalCollection.InsertOne(ctx, types.DBInfo{Version: config.Version})

	if err5 != nil {
		log.Error().
			Str("type", "module").
			Str("module", "db").
			Str("function", "initDB").
			Msg("Something bad append while adding default info in the database")

		return errors.New("something bad append while adding default version in the database")
	}

	log.Info().
		Str("type", "module").
		Str("module", "db").
		Str("function", "initDB").
		Msg("The DB has been init successfully")

	return
}
