package db

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/DarkOnion0/ADZTBotV2/config"
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
		Str("module", "db").
		Str("function", "checkUser").
		Str("userDiscordId", userDiscordId).
		Msg("Running the function")
	userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	var userList types.UserRecordFetch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "checkUser").
		Str("userDiscordId", userDiscordId).
		Msg("Fetching the user information")
	err1 := userInfoCollection.FindOne(ctx, bson.D{{Key: "userid", Value: userDiscordId}}).Decode(&userList)

	if err1 == mongo.ErrNoDocuments {
		log.Info().
			Str("type", "module").
			Str("module", "db").
			Str("function", "checkUser").
			Str("userDiscordId", userDiscordId).
			Msg("User doesn't exist in the database")
		return nil, false, userList.ID
	} else if err1 != nil {
		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "db").
			Str("function", "checkUser").
			Str("userDiscordId", userDiscordId).
			Msg("Something bad happen while finding the user in the database")

		return errors.New("something bad happen while finding the user in the database"), isUserExist, userList.ID
	}

	log.Info().
		Str("type", "module").
		Str("module", "db").
		Str("function", "checkUser").
		Str("userDiscordId", userDiscordId).
		Str("userDiscordId", userList.ID.Hex()).
		Msg("User exist in the database")
	return nil, true, userList.ID
}

// RegisterUser register a user if it doesn't exist in the database using his discord id
func RegisterUser(userDiscordId string) (err error) {
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "registerUser").
		Str("userDiscordId", userDiscordId).
		Msg("Running the function")

	errCheckUser, userStatus, userId := CheckUser(userDiscordId)

	if errCheckUser != nil {
		log.Error().
			Err(errCheckUser).
			Str("type", "module").
			Str("module", "db").
			Str("function", "registerUser").
			Str("userDiscordId", userDiscordId).
			Bool("userStatus", userStatus).
			Msg("Something bad happen while checking the user")

		return errors.New("something bad happen while checking the user in the database")
	}

	if !userStatus {
		userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Debug().
			Str("type", "module").
			Str("module", "db").
			Str("function", "registerUser").
			Str("userDiscordId", userDiscordId).
			Bool("userStatus", userStatus).
			Msg("Adding user to the database")
		info, err1 := userInfoCollection.InsertOne(ctx, types.UserRecordSend{Userid: userDiscordId})

		if err1 != nil {
			log.Error().
				Err(err1).
				Str("type", "module").
				Str("module", "db").
				Str("function", "registerUser").
				Str("userDiscordId", userDiscordId).
				Bool("userStatus", userStatus).
				Msg("Something bad happen while adding the user in the database")

			return errors.New("something bad happen while adding the user in the database")
		}

		log.Info().
			Str("type", "module").
			Str("module", "db").
			Str("function", "registerUser").
			Str("userDiscordId", userDiscordId).
			Str("userId", info.InsertedID.(primitive.ObjectID).Hex()).
			Bool("userStatus", userStatus).
			Msg("The user has been added successfully to the database")

		return
	}

	log.Info().
		Str("type", "module").
		Str("module", "db").
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
		Str("module", "db").
		Str("function", "getDiscordId").
		Str("userId", userId.Hex()).
		Msg("Running the function")
	userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	var userList types.UserRecordFetch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "getDiscordId").
		Str("userDiscordId", userDiscordId).
		Msg("Fetching user id from the db")
	err1 := userInfoCollection.FindOne(ctx, bson.D{{Key: "_id", Value: userId}}).Decode(&userList)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			log.Error().
				Err(err1).
				Str("type", "module").
				Str("module", "db").
				Str("function", "getDiscordId").
				Str("userDiscordId", userDiscordId).
				Msg("Something bad happen while fetching the user id from the db, suer doesn't exist")

			return errors.New("user doesn't exists in the database"), ""
		}
		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "db").
			Str("function", "getDiscordId").
			Str("userDiscordId", userDiscordId).
			Msg("Something bad happen while fetching the user id from the db")

		return errors.New("an error occurred while fetching the user id from the db"), ""
	}

	log.Info().
		Str("type", "module").
		Str("module", "db").
		Str("function", "getDiscordId").
		Str("userDiscordId", userDiscordId).
		Str("userId", userList.ID.Hex()).
		Msg("User exist in the db")

	return nil, userList.Userid
}

// GetUserInfo function get and return all the user infos according to the provided mongodb _id
func GetUserInfo(userId primitive.ObjectID) (err error, userStats types.UserInfo) {
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "getUserInfo").
		Str("userId", userId.Hex()).
		Msg("Running the function")

	// Init the mongodb collections
	postCollection := config.Client.Database(*config.DBName).Collection("post")
	userCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	// Query DB
	ctx0, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "getUserInfo").
		Str("userId", userId.Hex()).
		Msg("Fetching the user info")

	err0 := userCollection.FindOne(ctx0, bson.D{{Key: "_id", Value: userId}}).Decode(&userStats)

	if err0 != nil {
		log.Error().
			Err(err0).
			Str("type", "module").
			Str("module", "db").
			Str("function", "getUserInfo").
			Str("userId", userId.Hex()).
			Msg("Something bad happen while fetching the user info")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "db").
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
				Str("module", "db").
				Str("function", "getUserInfo").
				Str("userId", userId.Hex()).
				Msg("Something bad happen while closing the cursor")
			return
		}
	}(cursor, ctx)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			log.Error().
				Err(err1).
				Str("type", "module").
				Str("module", "db").
				Str("function", "getUserInfo").
				Str("userId", userId.Hex()).
				Msg("Something bad while fetching all the post from a user, the user has no post")
			return ErrNoDocument, types.UserInfo{}
		}

		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "db").
			Str("function", "getUserInfo").
			Str("userId", userId.Hex()).
			Msg("Something bad while fetching all the post from a user")
		return errors.New("an error occurred while fetching all the post from a user"), types.UserInfo{}
	}

	// Get all the posts in just one array
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "getUserInfo").
		Str("userId", userId.Hex()).
		Msg("Fetching all the post from the cursor")
	err2 := cursor.All(ctx, &userStats.Posts)
	if err2 != nil {
		log.Error().
			Err(err2).
			Str("type", "module").
			Str("module", "db").
			Str("function", "getUserInfo").
			Str("userId", userId.Hex()).
			Msg("Something bad while fetching all the post document")

		return ErrFetchingPost, types.UserInfo{}
	}

	// Return an error if a user as posted anything
	if len(userStats.Posts) == 0 {
		log.Error().
			Str("type", "module").
			Str("module", "db").
			Str("function", "getUserInfo").
			Str("userId", userId.Hex()).
			Msg("Something bad while fetching all the post from a user, the user has no post")
		return ErrNoPost, types.UserInfo{}
	}

	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "getUserInfo").
		Str("userId", userId.Hex()).
		Msg("Iterating over all the fetched document to count score")
	// iterate over all the fetched document
	for i := 0; i < len(userStats.Posts); i++ {
		scorePost, _ := CountScorePost(userStats.Posts[i], userId)
		log.Printf("%s", strconv.Itoa(scorePost))
		// update the score
		userStats.GlobalScore += scorePost
		log.Printf("%s", strconv.Itoa(userStats.GlobalScore))
	}

	log.Info().
		Str("type", "module").
		Str("module", "db").
		Str("function", "getUserInfo").
		Str("userId", userId.Hex()).
		Int("globalScore", userStats.GlobalScore).
		Msg("Getting user info succeed!")

	return err, userStats
}

// Fetch all the users infos of the database and return them as list
func FetchAllUsers() (err error, userStatsList []types.UserInfo) {
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "fetchAllUsers").
		Msg("Running the function")

	// init the mongodb collection
	userCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	var userList types.UserInfoList

	// Query DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "fetchAllUsers").
		Msg("Fetching all the registered users")
	cursor, err1 := userCollection.Find(ctx, bson.D{})
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Error().
				Err(err).
				Str("type", "module").
				Str("module", "db").
				Str("function", "fetchAllUsers").
				Msg("Something bad happen while closing the cursor")
			return
		}
	}(cursor, ctx)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			log.Error().
				Err(err1).
				Str("type", "module").
				Str("module", "db").
				Str("function", "fetchAllUsers").
				Msg("Something bad while fetching all registered users, there is no user in the db")
			return
		}

		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "db").
			Str("function", "fetchAllUsers").
			Msg("Something bad while fetching all the post from a user")
		return errors.New("an error occurred while fetching all the registered users"), userStatsList
	}

	// Get all the users in just one array
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "fetchAllUsers").
		Msg("Fetching all the post from the cursor")
	err2 := cursor.All(ctx, &userList)
	if err2 != nil {
		log.Error().
			Err(err2).
			Str("type", "module").
			Str("module", "db").
			Str("function", "fetchAllUsers").
			Msg("Something bad while fetching all the user document")

		return ErrFetchingPost, userStatsList
	}

	// Return an error if a user as posted anything
	if len(userList) == 0 {
		log.Error().
			Str("type", "module").
			Str("module", "db").
			Str("function", "fetchAllUsers").
			Msg("Something bad while fetching all the user, the db has no member registered")
		return ErrNoPost, userStatsList
	}

	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "fetchAllUsers").
		Msg("Iterating over all the fetched document to get all the information for each user")

	// iterate over all the fetched users and append their infos in a list
	for i := 0; i < len(userList); i++ {
		err1, userStats := GetUserInfo(userList[i].ID)

		if err1 != nil {
			switch err1 {
			case ErrNoPost:
				log.Error().
					Err(err1).
					Str("type", "module").
					Str("module", "db").
					Str("function", "fetchAllUsers").
					Msg("Something bad happen, user has no post")
			default:
				log.Error().
					Err(err1).
					Str("type", "module").
					Str("module", "db").
					Str("function", "fetchAllUsers").
					Msg("Something bad happen while fetching user info")

				return errors.New("something bad happen while fetching user info"), userStatsList
			}
		}

		userStatsList = append(userStatsList, userStats)
	}

	log.Info().
		Str("type", "module").
		Str("module", "db").
		Str("function", "fetchAllUsers").
		Msg("Getting user info succeed!")

	return
}

// Update a user info field according to the provided userId (to scope the research), the dbKey as the key to modify and the dbValue as the new value to set
func UpdateUser(userId primitive.ObjectID, dbKey string, dbValue interface{}) (err error) {
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "updateUser").
		Str("userId", userId.Hex()).
		Str("dbKey", dbKey).
		Interface("dbValue", dbValue).
		Msg("Running the function")

	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "updateUser").
		Str("userId", userId.Hex()).
		Str("dbKey", dbKey).
		Interface("dbValue", dbValue).
		Msg("Fetching the database")
	userCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	// update the user info
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "db").
		Str("function", "updateUser").
		Str("userId", userId.Hex()).
		Str("dbKey", dbKey).
		Interface("dbValue", dbValue).
		Msg("Updating the user info")
	_, err3 := userCollection.UpdateOne(ctx, bson.M{"_id": userId}, bson.D{
		{Key: "$set", Value: bson.D{{Key: dbKey, Value: dbValue}}},
	})

	if err3 != nil {
		log.Error().
			Err(err3).
			Str("type", "module").
			Str("module", "db").
			Str("function", "updateUser").
			Str("userId", userId.Hex()).
			Str("dbKey", dbKey).
			Interface("dbValue", dbValue).
			Msg("Something bad happen while updating the user info")

		return errors.New("an error append while updating the user info")
	}

	log.Info().
		Str("type", "module").
		Str("module", "db").
		Str("function", "updateUser").
		Str("userId", userId.Hex()).
		Str("dbKey", dbKey).
		Interface("dbValue", dbValue).
		Msg("The user info has been successfully updated")

	return
}
