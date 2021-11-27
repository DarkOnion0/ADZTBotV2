package db

import (
	"context"
	"log"
	"time"

	"DarkBotV2/config"
	"go.mongodb.org/mongo-driver/bson"
)

type userRecord struct {
	Userid string
}

func RegisterUser(userId string) {
	userInfoCollection := config.Client.Database(config.DBName).Collection("userInfo")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var userList userRecord
	_ = userInfoCollection.FindOne(ctx, bson.D{{"userid", userId}}).Decode(&userList)

	//fmt.Println(userList, userId)

	if len(userList.Userid) == 0 {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		info, _ := userInfoCollection.InsertOne(ctx, userRecord{Userid: userId})
		log.Printf("A new user has been added to the database; userid=%s dbId=%s", userId, info.InsertedID)
	} else {
		log.Printf("User %s already exists in the database", userId)
	}
}
