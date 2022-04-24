package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/DarkOnion0/ADZTBotV2/config"
	"github.com/DarkOnion0/ADZTBotV2/functions"
	"github.com/DarkOnion0/ADZTBotV2/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/rs/zerolog/log"
)

var ErrFetchingPost = errors.New("somethings bad append while fetching post")

var ErrNoPost = errors.New("the selected user has no post shared")

var ErrUserAlreadyRegistered = errors.New("the selected user is already registered in the database")

// CheckUser function check if a user exists in the database according to his discord id
func CheckUser(userDiscordId string) (err error, isUserExist bool, userId primitive.ObjectID) {
	log.Debug().
		Str("type", "module").
		Str("module", "user").
		Str("function", "checkUser").
		Str("userDiscordId", userDiscordId).
		Msg("Running the function")
	userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	var userList types.UserRecordFetch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "user").
		Str("function", "checkUser").
		Str("userDiscordId", userDiscordId).
		Msg("Fetching the user information")
	err1 := userInfoCollection.FindOne(ctx, bson.D{{Key: "userid", Value: userDiscordId}}).Decode(&userList)

	if err1 != nil {
		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "user").
			Str("function", "checkUser").
			Str("userDiscordId", userDiscordId).
			Msg("Something bad append while finding the user in the database")

		return errors.New("something bad append while finding the user in the database"), isUserExist, userList.ID
	}

	//fmt.Println(userList.Userid, userDiscordId)

	if len(userList.Userid) == 0 {
		log.Info().
			Str("type", "module").
			Str("module", "user").
			Str("function", "checkUser").
			Str("userDiscordId", userDiscordId).
			Msg("User doesn't exist in the database")
		return nil, false, userList.ID
	} else {
		log.Info().
			Str("type", "module").
			Str("module", "user").
			Str("function", "checkUser").
			Str("userDiscordId", userDiscordId).
			Str("userDiscordId", userList.ID.Hex()).
			Msg("User exist in the database")
		return nil, true, userList.ID
	}
}

// RegisterUser register a user if it doesn't exist in the database using his discord id
func RegisterUser(userDiscordId string) (err error) {
	log.Debug().
		Str("type", "module").
		Str("module", "user").
		Str("function", "registerUser").
		Str("userDiscordId", userDiscordId).
		Msg("Running the function")

	errCheckUser, userStatus, userId := CheckUser(userDiscordId)

	if errCheckUser != nil {
		log.Error().
			Err(errCheckUser).
			Str("type", "module").
			Str("module", "user").
			Str("function", "registerUser").
			Str("userDiscordId", userDiscordId).
			Bool("userStatus", userStatus).
			Msg("Something bad append while checking the user")

		return errors.New("something bad append while checking the user in the database")
	}

	if !userStatus {
		userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Debug().
			Str("type", "module").
			Str("module", "user").
			Str("function", "registerUser").
			Str("userDiscordId", userDiscordId).
			Bool("userStatus", userStatus).
			Msg("Adding user to the database")
		info, err1 := userInfoCollection.InsertOne(ctx, types.UserRecordSend{Userid: userDiscordId})

		if err1 != nil {
			log.Error().
				Err(err1).
				Str("type", "module").
				Str("module", "user").
				Str("function", "registerUser").
				Str("userDiscordId", userDiscordId).
				Bool("userStatus", userStatus).
				Msg("Something bad append while adding the user in the database")

			return errors.New("something bad append while adding the user in the database")
		}

		log.Info().
			Str("type", "module").
			Str("module", "user").
			Str("function", "registerUser").
			Str("userDiscordId", userDiscordId).
			Str("userId", info.InsertedID.(primitive.ObjectID).Hex()).
			Bool("userStatus", userStatus).
			Msg("The user has been added successfully to the database")

		return
	}

	log.Info().
		Str("type", "module").
		Str("module", "user").
		Str("function", "registerUser").
		Str("userDiscordId", userDiscordId).
		Str("userId", userId.Hex()).
		Bool("userStatus", userStatus).
		Msg("The user hasn't been added to the database, the selected user is already registered in the database")

	return ErrUserAlreadyRegistered

}

// GetDiscordId function get and return the user discord id according to the provided mongodb _id
func GetDiscordId(userId primitive.ObjectID) (err error, userDiscordId string) {
	log.Debug().
		Str("type", "module").
		Str("module", "user").
		Str("function", "getDiscordId").
		Str("userId", userId.Hex()).
		Msg("Running the function")
	userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	var userList types.UserRecordFetch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "user").
		Str("function", "getDiscordId").
		Str("userDiscordId", userDiscordId).
		Msg("Fetching user id from the db")
	err1 := userInfoCollection.FindOne(ctx, bson.D{{Key: "_id", Value: userId}}).Decode(&userList)

	//fmt.Println(userList.Userid, userDiscordId)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			log.Error().
				Err(err1).
				Str("type", "module").
				Str("module", "user").
				Str("function", "getDiscordId").
				Str("userDiscordId", userDiscordId).
				Msg("Something bad append while fetching the user id from the db, suer doesn't exist")

			return errors.New("user doesnt exists in the database"), ""
		}
		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "user").
			Str("function", "getDiscordId").
			Str("userDiscordId", userDiscordId).
			Msg("Something bad append while fetching the user id from the db")

		return errors.New("an error occurred while fetching the user id from the db"), ""
	}

	log.Info().
		Str("type", "module").
		Str("module", "user").
		Str("function", "getDiscordId").
		Str("userDiscordId", userDiscordId).
		Str("userId", userList.ID.Hex()).
		Msg("User exist in the db")

	return nil, userList.Userid
}

// GetUserInfo function get and return all the user infos according to the provided mongodb _id
func GetUserInfo(userId primitive.ObjectID) (err error, userStats types.UserInfoFetch) {
	log.Debug().
		Str("type", "module").
		Str("module", "user").
		Str("function", "getUserInfo").
		Str("userId", userId.Hex()).
		Msg("Running the function")

	// init the mongodb collection
	postCollection := config.Client.Database(*config.DBName).Collection("post")

	// Init the users stats var
	userStats = types.UserInfoFetch{ID: userId}

	// Query DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "user").
		Str("function", "getUserInfo").
		Str("userId", userId.Hex()).
		Msg("Fetching all the post from a user")
	cursor, err1 := postCollection.Find(ctx, bson.D{{Key: "user", Value: userId}})
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Error().
				Err(err).
				Str("type", "module").
				Str("module", "user").
				Str("function", "getUserInfo").
				Str("userId", userId.Hex()).
				Msg("Something bad append while closing the cursor")
			return
		}
	}(cursor, ctx)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			log.Error().
				Err(err1).
				Str("type", "module").
				Str("module", "user").
				Str("function", "getUserInfo").
				Str("userId", userId.Hex()).
				Msg("Something bad while fetching all the post from a user, the user has no post")
			return ErrNoDocument, types.UserInfoFetch{}
		}

		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "user").
			Str("function", "getUserInfo").
			Str("userId", userId.Hex()).
			Msg("Something bad while fetching all the post from a user")
		return errors.New("an error occurred while fetching all the post from a user"), types.UserInfoFetch{}
	}

	// Get all the posts in just one array
	log.Debug().
		Str("type", "module").
		Str("module", "user").
		Str("function", "getUserInfo").
		Str("userId", userId.Hex()).
		Msg("Fetching all the post from the cursor")
	err2 := cursor.All(ctx, &userStats.Posts)
	if err2 != nil {
		log.Error().
			Err(err2).
			Str("type", "module").
			Str("module", "user").
			Str("function", "getUserInfo").
			Str("userId", userId.Hex()).
			Msg("Something bad while fetching all the post document")

		return ErrFetchingPost, types.UserInfoFetch{}
	}

	fmt.Sprintln(userStats.Posts)

	// Return an error if a user as posted anything
	if len(userStats.Posts) == 0 {
		log.Error().
			Str("type", "module").
			Str("module", "user").
			Str("function", "getUserInfo").
			Str("userId", userId.Hex()).
			Msg("Something bad while fetching all the post from a user, the user has no post")
		return ErrNoPost, types.UserInfoFetch{}
	}

	log.Debug().
		Str("type", "module").
		Str("module", "user").
		Str("function", "getUserInfo").
		Str("userId", userId.Hex()).
		Msg("Iterating over all the fetched document to count score")
	// iterate over all the fetched document
	for i := 0; i < len(userStats.Posts); i++ {
		scorePost, _ := functions.CountScorePost(userStats.Posts[i], userId)
		log.Printf("%s", strconv.Itoa(scorePost))
		// update the score
		userStats.GlobalScore += scorePost
		log.Printf("%s", strconv.Itoa(userStats.GlobalScore))
	}

	log.Info().
		Str("type", "module").
		Str("module", "user").
		Str("function", "getUserInfo").
		Str("userId", userId.Hex()).
		Int("globalScore", userStats.GlobalScore).
		Msg("Getting user info succeed!")

	return err, userStats
}
