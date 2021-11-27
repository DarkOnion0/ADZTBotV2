package db

import (
	"context"
	"log"
	"time"

	"ADZTBotV2/config"
	"go.mongodb.org/mongo-driver/bson"
)

type userRecord struct {
	Userid string
}

func CheckUser(userId string) bool {
	userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")

	var userList userRecord
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_ = userInfoCollection.FindOne(ctx, bson.D{{"userid", userId}}).Decode(&userList)

	//fmt.Println(userList, userId)

	if len(userList.Userid) == 0 {
		log.Printf("User %s doesn't exist in the database", userId)
		return false
	} else {
		log.Printf("User %s already exists in the database", userId)
		return true
	}
}

func RegisterUser(userId string) {
	if !CheckUser(userId) {
		userInfoCollection := config.Client.Database(*config.DBName).Collection("userInfo")
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		info, _ := userInfoCollection.InsertOne(ctx, userRecord{Userid: userId})
		log.Printf("A new user has been added to the database; userid=%s dbId=%s", userId, info.InsertedID)
	}
}
