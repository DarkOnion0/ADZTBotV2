package db

import (
	"context"
	"errors"
	"time"

	"github.com/DarkOnion0/ADZTBotV2/config"
	"github.com/blang/semver/v4"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
)

// The document containing the version was not found in the database
var ErrVersionNotFound error = errors.New("the version was not found in the database")

// This function update the DB version
func UpdateVersion(version semver.Version) (err error) {
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "updateVersion").
		Msg("Running the function")

	// init the mongodb collection
	internalCollection := config.Client.Database(*config.DBName).Collection("internal")

	// update the version number in the db
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
			Msg("Something bad happen while updating version in the database")

		return errors.New("something bad happen while updating version in the database")
	}

	log.Info().
		Str("type", "module").
		Str("module", "db").
		Str("function", "updateVersion").
		Msg("The DB version has been updated successfully")

	return
}
