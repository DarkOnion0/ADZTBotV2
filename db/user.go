package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"ADZTBotV2/config"

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

// CheckUser function check if a user exists in the database according to his discord id
func CheckUser(userId string) (bool, primitive.ObjectID) {
	userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	var userList userRecordFetch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = userInfoCollection.FindOne(ctx, bson.D{{"userid", userId}}).Decode(&userList)

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

func GetUser(userDbId primitive.ObjectID) (err error, userId string) {
	userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	var userList userRecordFetch
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err1 := userInfoCollection.FindOne(ctx, bson.D{{"_id", userDbId}}).Decode(&userList)

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
