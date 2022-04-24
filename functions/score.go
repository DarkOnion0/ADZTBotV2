package functions

import (
	"strconv"

	"github.com/DarkOnion0/ADZTBotV2/types"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// CountScorePost function calculate the total score of a post according to the provided post (postRecord),
// it can also return the score of a specific user on this post according to the provided db id (userDbId)
func CountScorePost(postRecord types.PostRecordFetchT, userId primitive.ObjectID) (globalVote int, userVote string) {
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
