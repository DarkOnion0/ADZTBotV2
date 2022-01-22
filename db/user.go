package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/DarkOnion0/ADZTBotV2/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// This is the datastructures of every mongodb record in the userInfo collection
type userRecordFetch struct {
	ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Userid string
}
type userRecordSend struct {
	//ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Userid string
}

type UserInfoFetch struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Posts       []PostRecordFetchT
	GlobalScore int
}

// CheckUser function check if a user exists in the database according to his discord id
func CheckUser(userId string) (bool, primitive.ObjectID) {
	userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	var userList userRecordFetch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = userInfoCollection.FindOne(ctx, bson.D{{Key: "userid", Value: userId}}).Decode(&userList)

	//fmt.Println(userList.Userid, userId)

	if len(userList.Userid) == 0 {
		log.Printf("User %s doesn't exist in the database", userId)
		return false, userList.ID
	} else {
		log.Printf("User %s already exists in the database", userId)
		return true, userList.ID
	}
}

// RegisterUser register a user if it doesn't exist in the database using his discord id
func RegisterUser(userId string) {
	userStatus, _ := CheckUser(userId)
	if !userStatus {
		userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		info, _ := userInfoCollection.InsertOne(ctx, userRecordSend{Userid: userId})
		fmt.Println(info, userId)
		log.Printf("A new user has been added to the database; userid=%s dbId=%s", userId, info.InsertedID)
	}
}

// GetDiscordId function get and return the user discord id according to the provided mongodb _id
func GetDiscordId(userDbId primitive.ObjectID) (err error, userId string) {
	userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	var userList userRecordFetch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err1 := userInfoCollection.FindOne(ctx, bson.D{{Key: "_id", Value: userDbId}}).Decode(&userList)

	//fmt.Println(userList.Userid, userId)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			return errors.New("user doesnt exists in the database"), ""
		}
		log.Fatalf("An error occured while fetching the userId: %s", err1)
		//return errors.New("An error occurred while fetching the post"), false
	} else {
		log.Printf("User %s already exists in the database", userId)
		return nil, userList.Userid
	}

	return errors.New("the function shouldn't arrive there"), ""
}

// GetUserInfo function get and return all the user infos according to the provided mongodb _id
func GetUserInfo(userDbId primitive.ObjectID) (err int, userStats UserInfoFetch) {
	// TODO update the error handling in this function

	postCollection := config.Client.Database(*config.DBName).Collection("post")

	userStats = UserInfoFetch{ID: userDbId}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err1 := postCollection.Find(ctx, bson.D{{Key: "user", Value: userDbId}})
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatalf("Something bad append %s", err)
		}
	}(cursor, ctx)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			return 1, UserInfoFetch{}
		}
		log.Fatalf("An error occured while fetching the userId: %s", err1)
		//return errors.New("An error occurred while fetching the post"), false
	}

	log.Printf("%s ", strconv.Itoa(userStats.GlobalScore))

	err2 := cursor.All(ctx, &userStats.Posts)
	if err2 != nil {
		log.Fatalf("Something bad append while fetching all the document in mongodb %s", err2)

		return 2, UserInfoFetch{}
	}

	log.Println(userStats.Posts)

	if len(userStats.Posts) == 0 {
		return 3, UserInfoFetch{}
	}

	// iterate over all the fetched document
	for i := 0; i < len(userStats.Posts); i++ {
		scorePost, _ := countScorePost(userStats.Posts[i], userDbId)
		log.Printf("%s", strconv.Itoa(scorePost))
		userStats.GlobalScore += scorePost
		log.Printf("%s", strconv.Itoa(userStats.GlobalScore))
	}

	return err, userStats
}
