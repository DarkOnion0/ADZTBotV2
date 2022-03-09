package db

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/DarkOnion0/ADZTBotV2/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type postRecordSendT struct {
	//ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Type     string
	Url      string
	User     primitive.ObjectID
	VoteList []postVote
}

type PostRecordFetchT struct {
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

var ErrNoDocument = errors.New("the selected post doesn't exist")
var ErrWrongUserDbId = errors.New("the provided user db id is different from the wanted one")

// The Post function check if post exist in the database according to his link.
//
// 1. If it not exists the post will be added and the function will return (false, OBJECTID)
//
// 2. Else it will return (true, "")
func Post(userId primitive.ObjectID, postType, postUrl string) (postAlreadyExist bool, postId string) {
	log.Debug().
		Str("type", "module").
		Str("module", "post").
		Str("function", "post").
		Str("userId", userId.Hex()).
		Dict("post", zerolog.Dict().
			Str("type", postType).
			Str("url", postUrl)).
		Msg("Running the function")

	postCollection := config.Client.Database(*config.DBName).Collection("post")

	// check if the post already exist in the DB according to its url
	var postRecordFetch PostRecordFetchT
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Debug().
		Str("type", "module").
		Str("module", "post").
		Str("function", "post").
		Str("userId", userId.Hex()).
		Dict("post", zerolog.Dict().
			Str("type", postType).
			Str("url", postUrl)).
		Msg("Fetching the post information")

	// splitting url on ?si= to avoid adding the random identifier that spotify add at the end of the music
	// that beak the post checking mechanism
	err1 := postCollection.FindOne(ctx, bson.D{{Key: "url", Value: strings.Split(postUrl, "?si=")[0]}}).Decode(&postRecordFetch)

	if err1 != nil {
		// add new post to the db
		if err1 == mongo.ErrNoDocuments {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			info, _ := postCollection.InsertOne(ctx, postRecordSendT{Type: postType, Url: strings.Split(postUrl, "?si=")[0], User: userId, VoteList: []postVote{{User: userId, Vote: "+"}}})

			log.Info().
				Str("type", "module").
				Str("module", "post").
				Str("function", "post").
				Str("userId", userId.Hex()).
				Dict("post", zerolog.Dict().
					Str("id", info.InsertedID.(primitive.ObjectID).Hex()).
					Str("type", postType).
					Str("url", postUrl)).
				Msg("The post has been successfully added to the database")

			return false, info.InsertedID.(primitive.ObjectID).Hex()
		}

		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "post").
			Str("function", "post").
			Str("userId", userId.Hex()).
			Dict("post", zerolog.Dict().
				Str("type", postType).
				Str("url", postUrl)).
			Msg("Somethings bad append while fetching the post url")

		return true, ""
		// tell back to the command that the post already exist
	} else {
		return true, ""
	}
}

// The SetVote function add or remove a like to post (set with the Post function)
//
// NOTE: it will return true if the vote has been added and false if not
func SetVote(postId, userVote string, userId primitive.ObjectID) (error, bool) {
	log.Debug().
		Str("type", "module").
		Str("module", "post").
		Str("function", "setVote").
		Str("userId", userId.Hex()).
		Str("userVote", userVote).
		Dict("post", zerolog.Dict().
			Str("id", postId)).
		Msg("Running the function")

	postCollection := config.Client.Database(*config.DBName).Collection("post")
	postIdPrimitive, err1 := primitive.ObjectIDFromHex(postId)

	if err1 != nil {
		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "post").
			Str("function", "setVote").
			Str("userId", userId.Hex()).
			Str("userVote", userVote).
			Dict("post", zerolog.Dict().
				Str("id", postId)).
			Msg("Something bad append while converting from Hex to primitive.ObjectID")

		return errors.New("something bad append while converting from Hex to primitive.ObjectID"), false
	}

	log.Debug().
		Str("type", "module").
		Str("module", "post").
		Str("function", "setVote").
		Str("userId", userId.Hex()).
		Str("userVote", userVote).
		Dict("post", zerolog.Dict().
			Str("id", postId)).
		Msg("Fetching the post information")

	var postRecordFetch PostRecordFetchT
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err2 := postCollection.FindOne(ctx, bson.D{{Key: "_id", Value: postIdPrimitive}}).Decode(&postRecordFetch)

	//fmt.Println(err2, strings.Split(postUrl, "?si=")[0])

	if err2 != nil {
		if err2 == mongo.ErrNoDocuments {
			log.Error().
				Err(err2).
				Str("type", "module").
				Str("module", "post").
				Str("function", "setVote").
				Str("userId", userId.Hex()).
				Str("userVote", userVote).
				Dict("post", zerolog.Dict().
					Str("id", postId)).
				Msg("The postId is not valid")

			return errors.New("the postID is not valid"), false
		}

		log.Error().
			Err(err2).
			Str("type", "module").
			Str("module", "post").
			Str("function", "setVote").
			Str("userId", userId.Hex()).
			Str("userVote", userVote).
			Dict("post", zerolog.Dict().
				Str("id", postId)).
			Msg("Something bad append while fetching the post")

		return errors.New("an error occurred while fetching the post"), false
	} else {
		alreadyVote := false

		log.Debug().
			Str("type", "module").
			Str("module", "post").
			Str("function", "setVote").
			Str("userId", userId.Hex()).
			Str("userVote", userVote).
			Bool("alreadyVote", alreadyVote).
			Dict("post", zerolog.Dict().
				Str("id", postId)).
			Msg("The post is valid, starting counting the vote...")

		for i := 0; i < len(postRecordFetch.VoteList); i++ {
			if postRecordFetch.VoteList[i].User == userId {
				alreadyVote = true

				log.Debug().
					Str("type", "module").
					Str("module", "post").
					Str("function", "setVote").
					Str("userId", userId.Hex()).
					Str("userVote", userVote).
					Str("userOldVote", postRecordFetch.VoteList[i].Vote).
					Bool("alreadyVote", alreadyVote).
					Dict("post", zerolog.Dict().
						Str("id", postId)).
					Msg("User has already voted, changing is vote...")

				postRecordFetch.VoteList[i].Vote = userVote
			}
		}

		if !alreadyVote {
			log.Debug().
				Str("type", "module").
				Str("module", "post").
				Str("function", "setVote").
				Str("userId", userId.Hex()).
				Str("userVote", userVote).
				Bool("alreadyVote", alreadyVote).
				Dict("post", zerolog.Dict().
					Str("id", postId)).
				Msg("User hasn't already voted, adding is vote to the DB...")

			// add new post to the db
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err3 := postCollection.UpdateOne(ctx, bson.M{"_id": postIdPrimitive}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "votelist", Value: append(postRecordFetch.VoteList, postVote{User: userId, Vote: userVote})}}},
			})

			if err3 != nil {
				log.Error().
					Err(err3).
					Str("type", "module").
					Str("module", "post").
					Str("function", "setVote").
					Str("userId", userId.Hex()).
					Str("userVote", userVote).
					Dict("post", zerolog.Dict().
						Str("id", postId)).
					Msg("Something bad append while adding the vote")

				return errors.New("an error occurred while adding the vote"), false
			}

			log.Info().
				Str("type", "module").
				Str("module", "post").
				Str("function", "setVote").
				Str("userId", userId.Hex()).
				Str("userVote", userVote).
				Bool("alreadyVote", alreadyVote).
				Dict("post", zerolog.Dict().
					Str("id", postId)).
				Msg("The vote has been added successfully")

			return nil, true
		} else {
			log.Debug().
				Str("type", "module").
				Str("module", "post").
				Str("function", "setVote").
				Str("userId", userId.Hex()).
				Str("userVote", userVote).
				Bool("alreadyVote", alreadyVote).
				Dict("post", zerolog.Dict().
					Str("id", postId)).
				Msg("User has already voted, publishing his vote to the DB...")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err3 := postCollection.UpdateOne(ctx, bson.M{"_id": postIdPrimitive}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "votelist", Value: postRecordFetch.VoteList}}},
			})

			if err3 != nil {
				log.Error().
					Err(err3).
					Str("type", "module").
					Str("module", "post").
					Str("function", "setVote").
					Str("userId", userId.Hex()).
					Str("userVote", userVote).
					Dict("post", zerolog.Dict().
						Str("id", postId)).
					Msg("Something bad append while updating the vote")

				return errors.New("an error occurred while updating the vote"), false
			}

			log.Info().
				Str("type", "module").
				Str("module", "post").
				Str("function", "setVote").
				Str("userId", userId.Hex()).
				Str("userVote", userVote).
				Bool("alreadyVote", alreadyVote).
				Dict("post", zerolog.Dict().
					Str("id", postId)).
				Msg("The vote has been updated successfully")

			return nil, false
		}
	}
}

// GetVote function fetch and return all the information about a post according to the provided postId and userId
func GetVote(postId string, userId primitive.ObjectID) (err error, globalVote int, userVote string, postFetch PostRecordFetchT) {
	log.Debug().
		Str("type", "module").
		Str("module", "post").
		Str("function", "getVote").
		Str("userId", userId.Hex()).
		Dict("post", zerolog.Dict().
			Str("id", postId)).
		Msg("Running the function")

	postCollection := config.Client.Database(*config.DBName).Collection("post")
	postIdPrimitive, err1 := primitive.ObjectIDFromHex(postId)

	if err1 != nil {
		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "post").
			Str("function", "getVote").
			Str("userId", userId.Hex()).
			Dict("post", zerolog.Dict().
				Str("id", postId)).
			Msg("Something bad append while converting from Hex to primitive.ObjectID")

		return errors.New("something bad append while converting from Hex to primitive.ObjectID"), globalVote, userVote, postFetch
	}

	// make a mongodb request to get the post information according to the provided postId
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "post").
		Str("function", "getVote").
		Str("userId", userId.Hex()).
		Dict("post", zerolog.Dict().
			Str("id", postId)).
		Msg("Fetching the post information")
	err2 := postCollection.FindOne(ctx, bson.D{{Key: "_id", Value: postIdPrimitive}}).Decode(&postFetch)

	if err2 != nil {
		if err2 == mongo.ErrNoDocuments {
			log.Error().
				Err(err2).
				Str("type", "module").
				Str("module", "post").
				Str("function", "getVote").
				Str("userId", userId.Hex()).
				Str("userVote", userVote).
				Dict("post", zerolog.Dict().
					Str("id", postId)).
				Msg("The postId is not valid")

			return errors.New("the postID is not valid"), globalVote, userVote, postFetch
		}

		log.Error().
			Err(err2).
			Str("type", "module").
			Str("module", "post").
			Str("function", "getVote").
			Str("userId", userId.Hex()).
			Str("userVote", userVote).
			Dict("post", zerolog.Dict().
				Str("id", postId)).
			Msg("Something bad append while fetching the post")

		return errors.New("an error occurred while fetching the post"), globalVote, userVote, postFetch
		//return errors.New("An error occurred while fetching the post"), false
	} else {
		log.Info().
			Str("type", "module").
			Str("module", "post").
			Str("function", "getVote").
			Str("userId", userId.Hex()).
			Str("userVote", userVote).
			Dict("post", zerolog.Dict().
				Str("id", postId)).
			Msg("The vote has been get successfully, counting its score...")

		err = nil
		// count the vote score of the post
		globalVote, userVote = CountScorePost(postFetch, userId)

		return err, globalVote, userVote, postFetch
	}
}

// CountScorePost function calculate the total score of a post according to the provided post (postRecord),
// it can also return the score of a specific user on this post according to the provided db id (userDbId)
func CountScorePost(postRecord PostRecordFetchT, userId primitive.ObjectID) (globalVote int, userVote string) {
	userVote = "You haven't yet vote on this post 😅"
	globalVote = 0

	log.Debug().
		Str("type", "module").
		Str("module", "post").
		Str("function", "countScorePost").
		Str("userId", userId.Hex()).
		Str("userVote", userVote).
		Str("globalVote", strconv.Itoa(globalVote)).
		Dict("post", zerolog.Dict().
			Str("id", postRecord.ID.Hex())).
		Msg("Running the function")

	for i := 0; i < len(postRecord.VoteList); i++ {
		if postRecord.VoteList[i].User == userId {
			if postRecord.VoteList[i].Vote == "+" {
				userVote = "Like 👍"
				globalVote += 1
			} else {
				userVote = "Dislike 👎"
				globalVote += -1
			}
		} else {
			if postRecord.VoteList[i].Vote == "+" {
				globalVote += 1
			} else {
				globalVote += -1
			}
		}
	}

	log.Info().
		Str("type", "module").
		Str("module", "post").
		Str("function", "countScorePost").
		Str("userId", userId.Hex()).
		Str("userVote", userVote).
		Str("globalVote", strconv.Itoa(globalVote)).
		Dict("post", zerolog.Dict().
			Str("id", postRecord.ID.Hex())).
		Msg("All the score has been computed successfully")

	return
}

func DeletePost(postId, userId primitive.ObjectID, isBotAdmin bool) (err error) {
	log.Debug().
		Str("type", "module").
		Str("module", "post").
		Str("function", "deletePost").
		Str("userId", userId.Hex()).
		Bool("isBotAdmin", isBotAdmin).
		Dict("post", zerolog.Dict().
			Str("id", postId.Hex())).
		Msg("Running the function")

	postCollection := config.Client.Database(*config.DBName).Collection("post")

	var postRecordFetch PostRecordFetchT
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug().
		Str("type", "module").
		Str("module", "post").
		Str("function", "deletePost").
		Str("userId", userId.Hex()).
		Bool("isBotAdmin", isBotAdmin).
		Dict("post", zerolog.Dict().
			Str("id", postId.Hex())).
		Msg("Searching the post in the DB")
	err1 := postCollection.FindOne(ctx, bson.D{{Key: "_id", Value: postId}}).Decode(&postRecordFetch)

	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			log.Error().
				Err(err1).
				Str("type", "module").
				Str("module", "post").
				Str("function", "deletePost").
				Str("userId", userId.Hex()).
				Bool("isBotAdmin", isBotAdmin).
				Dict("post", zerolog.Dict().
					Str("id", postId.Hex())).
				Msg("The selected post doesn't exist")

			return ErrNoDocument
		}

		log.Error().
			Err(err1).
			Str("type", "module").
			Str("module", "post").
			Str("function", "deletePost").
			Str("userId", userId.Hex()).
			Bool("isBotAdmin", isBotAdmin).
			Dict("post", zerolog.Dict().
				Str("id", postId.Hex())).
			Msg("Something bad append while searching the post in the db")

		return errors.New("something bad append while searching the post in the db")
	}

	if postRecordFetch.User == userId || isBotAdmin {
		log.Debug().
			Str("type", "module").
			Str("module", "post").
			Str("function", "deletePost").
			Str("userId", userId.Hex()).
			Bool("isBotAdmin", isBotAdmin).
			Dict("post", zerolog.Dict().
				Str("id", postId.Hex())).
			Msg("The post exist, deleting it...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err2 := postCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: postId}})

		if err2 != nil {
			log.Error().
				Err(err2).
				Str("type", "module").
				Str("module", "post").
				Str("function", "deletePost").
				Str("userId", userId.Hex()).
				Bool("isBotAdmin", isBotAdmin).
				Dict("post", zerolog.Dict().
					Str("id", postId.Hex())).
				Msg("Something bad append while deleting the post")

			return errors.New("something bad append while deleting the post")
		}

		log.Info().
			Str("type", "module").
			Str("module", "post").
			Str("function", "deletePost").
			Str("userId", userId.Hex()).
			Bool("isBotAdmin", isBotAdmin).
			Dict("post", zerolog.Dict().
				Str("id", postId.Hex())).
			Msg("The post has been successfully deleted")

		return
	}

	log.Error().
		Str("type", "module").
		Str("module", "post").
		Str("function", "deletePost").
		Str("userId", userId.Hex()).
		Bool("isBotAdmin", isBotAdmin).
		Dict("post", zerolog.Dict().
			Str("id", postId.Hex())).
		Msg("The provided user doesnt match the author id and is not a bot admin")

	return ErrWrongUserDbId
}
