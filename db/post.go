package db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"ADZTBotV2/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type postRecord struct {
	//ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Type     string
	Url      string
	User     primitive.ObjectID
	VoteList []postVote
}

type postVote struct {
	User primitive.ObjectID
	Vote string
}

func Post(userDbId primitive.ObjectID, postType, postUrl string) (bool, string) {
	userInfoCollection := config.Client.Database(*config.DBName).Collection("post")

	var postRecordFetch postRecord
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err1 := userInfoCollection.FindOne(ctx, bson.D{{"Url", strings.Split(postUrl, "?si=")[0]}}).Decode(&postRecordFetch)

	fmt.Println(err1)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			// add new post to the db
			postCollection := config.Client.Database(*config.DBName).Collection("post")

			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

			info, _ := postCollection.InsertOne(ctx, postRecord{Type: postType, Url: strings.Split(postUrl, "?si=")[0], User: userDbId, VoteList: []postVote{}})
			log.Printf("A new user post has been added to the database; userDbId=%s url=%s type=%s dbId=%s", userDbId, postUrl, postType, info.InsertedID)
			return false, fmt.Sprintf("%s", info.InsertedID)
		}
		log.Fatalf("Somethings bad append while fetching the post url in the Post function: %s", err1)
	} else {
		// tell back to the command that the post already exist
		return true, ""
	}

	return true, ""
}