package db

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"ADZTBotV2/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type postRecordSendT struct {
	//ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Type     string
	Url      string
	User     primitive.ObjectID
	VoteList []postVote
}

type postRecordFetchT struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Type     string
	Url      string
	User     primitive.ObjectID
	VoteList []postVote
}

type postVote struct {
	User primitive.ObjectID
	Vote string
}

// The Post function check if post exist in the database according to his link.
//
// 1. If it not exists the post will be added and the function will return (true, OBJECTID)
//
// 2. Else it will return (false, "")
func Post(userDbId primitive.ObjectID, postType, postUrl string) (bool, string) {
	postCollection := config.Client.Database(*config.DBName).Collection("post")

	var postRecordFetch postRecordFetchT
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err1 := postCollection.FindOne(ctx, bson.D{{Key: "url", Value: strings.Split(postUrl, "?si=")[0]}}).Decode(&postRecordFetch)

	//fmt.Println(err1, strings.Split(postUrl, "?si=")[0])

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			// add new post to the db
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			info, _ := postCollection.InsertOne(ctx, postRecordSendT{Type: postType, Url: strings.Split(postUrl, "?si=")[0], User: userDbId, VoteList: []postVote{{User: userDbId, Vote: "+"}}})
			log.Printf("A new user post has been added to the database; userDbId=%s url=%s type=%s dbId=%s", userDbId, postUrl, postType, info.InsertedID.(primitive.ObjectID).Hex())
			return false, info.InsertedID.(primitive.ObjectID).Hex()
		}
		log.Fatalf("Somethings bad append while fetching the post url in the Post function: %s", err1)
	} else {
		// tell back to the command that the post already exist
		return true, ""
	}

	return true, ""
}

// The SetVote function add or remove a like to post (set with the Post function)
//
// NOTE: it will return true if the post has been added and false if not
func SetVote(postId, userVote string, userId primitive.ObjectID) (error, bool) {

	postCollection := config.Client.Database(*config.DBName).Collection("post")
	postIdPrimitive, err1 := primitive.ObjectIDFromHex(postId)

	if err1 != nil {
		log.Fatalf("An error occured while convertign the post hex ObjectID to the primitve.ObjectID")
	}

	var postRecordFetch postRecordFetchT
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err2 := postCollection.FindOne(ctx, bson.D{{Key: "_id", Value: postIdPrimitive}}).Decode(&postRecordFetch)

	//fmt.Println(err2, strings.Split(postUrl, "?si=")[0])

	if err2 != nil {
		if err2 == mongo.ErrNoDocuments {
			return errors.New("the postID is not valid"), false
		}
		log.Fatalf("An error occured while fetching the post: %s", err2)
		//return errors.New("An error occurred while fetching the post"), false
	} else {
		alreadyVote := false

		for i := 0; i < len(postRecordFetch.VoteList); i++ {
			if postRecordFetch.VoteList[i].User == userId {
				alreadyVote = true
				postRecordFetch.VoteList[i].Vote = userVote
			}
		}

		if !alreadyVote {
			// add new post to the db
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, _ = postCollection.UpdateOne(ctx, bson.M{"_id": postIdPrimitive}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "votelist", Value: append(postRecordFetch.VoteList, postVote{User: userId, Vote: userVote})}}},
			})
			log.Printf("Add a vote; postId=%s userVote=%s userId=%s", postId, userVote, userId.Hex())

			return nil, true
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, _ = postCollection.UpdateOne(ctx, bson.M{"_id": postIdPrimitive}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "votelist", Value: postRecordFetch.VoteList}}},
			})
			log.Printf("Update a vote; postId=%s userVote=%s userId=%s", postId, userVote, userId.Hex())

			return nil, false
		}
	}

	return errors.New("the function shouldn't arrive there"), false
}

func GetVote(postId string, userId primitive.ObjectID) (err error, globalVote int, userVote string, postFetch postRecordFetchT) {
	userVote = "You haven't yet vote on this post 😅"
	globalVote = 0
	postCollection := config.Client.Database(*config.DBName).Collection("post")
	postIdPrimitive, err1 := primitive.ObjectIDFromHex(postId)

	if err1 != nil {
		log.Fatalf("An error occured while convertign the post hex ObjectID to the primitve.ObjectID")
	}

	var postRecordFetch postRecordFetchT
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err2 := postCollection.FindOne(ctx, bson.D{{Key: "_id", Value: postIdPrimitive}}).Decode(&postRecordFetch)

	//fmt.Println(err2, strings.Split(postUrl, "?si=")[0])

	if err2 != nil {
		if err2 == mongo.ErrNoDocuments {
			return errors.New("the postID is not valid"), globalVote, userVote, postRecordFetchT{}
		}
		log.Fatalf("An error occured while fetching the post: %s", err2)
		//return errors.New("An error occurred while fetching the post"), false
	} else {
		err = nil

		for i := 0; i < len(postRecordFetch.VoteList); i++ {
			if postRecordFetch.VoteList[i].User == userId {
				if postRecordFetch.VoteList[i].Vote == "+" {
					userVote = "Like 👍"
					globalVote += 1
				} else {
					userVote = "Dislike 👎"
					globalVote += -1
				}
			} else {
				if postRecordFetch.VoteList[i].Vote == "+" {
					globalVote += 1
				} else {
					globalVote += -1
				}
			}
		}

		return err, globalVote, userVote, postRecordFetch
	}

	return errors.New("the function shouldn't arrive there"), globalVote, userVote, postRecordFetchT{}
}
