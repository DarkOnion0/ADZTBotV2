package db

import (
	"context"
	"log"
	"time"

	"ADZTBotV2/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// This is the datastructures of every mongodb record in the userInfo collection
type userRecord struct {
	ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Userid string
}

type CheckUserT struct {
	UserExist    bool
	UserObjectId primitive.ObjectID
}

// CheckUser function check if a user exists in the database according to his discord id
func CheckUser(userId string) CheckUserT {
	userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	var userList userRecord
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_ = userInfoCollection.FindOne(ctx, bson.D{{"userid", userId}}).Decode(&userList)

	//fmt.Println(userList, userId)

	if len(userList.Userid) == 0 {
		log.Printf("User %s doesn't exist in the database", userId)
		return CheckUserT{UserObjectId: userList.ID, UserExist: false}
	} else {
		log.Printf("User %s already exists in the database", userId)
		return CheckUserT{UserObjectId: userList.ID, UserExist: true}
	}
}

// RegisterUser register a user if it doesn't exist in the database using his discord id
func RegisterUser(userId string) {
	if !CheckUser(userId).UserExist {
		userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		info, _ := userInfoCollection.InsertOne(ctx, userRecord{Userid: userId})
		log.Printf("A new user has been added to the database; userid=%s dbId=%s", userId, info.InsertedID)
	}
}
