package functions

import (
	"errors"
	"sort"

	"github.com/rs/zerolog/log"

	"github.com/DarkOnion0/ADZTBotV2/db"
)

// Fetch all the user and their stats in the database in order to class them increasingly according to their global score
func UpdateUserRanking() (err error) {
	log.Debug().
		Str("type", "module").
		Str("module", "functions").
		Str("function", "updateUserRanking").
		Msg("Running the function")

	log.Debug().
		Str("type", "module").
		Str("module", "functions").
		Str("function", "updateUserRanking").
		Msg("Fetching all the user in the database")
	err2, userStatsList := db.FetchAllUsers()

	if err2 != nil {
		log.Error().
			Err(err2).
			Str("type", "module").
			Str("module", "functions").
			Str("function", "updateUserRanking").
			Msg("Something bad happen while fetching all the users")

		return errors.New("something bad happen while fetching all the users")
	}
	log.Debug().
		Str("type", "module").
		Str("module", "functions").
		Str("function", "updateUserRanking").
		Msg("All the user was fetched successfully")

	// Order list in decreasing order
	sort.Slice(userStatsList, func(x, y int) bool {
		return userStatsList[x].GlobalScore > userStatsList[y].GlobalScore
	})

	// Add user rank to the DB
	for i := 0; i < len(userStatsList); i++ {
		err1 := db.UpdateUser(userStatsList[i].ID, "rank", i+1)

		if err1 != nil {
			log.Error().
				Err(err1).
				Str("type", "module").
				Str("module", "functions").
				Str("function", "updateUserRanking").
				Str("userId", userStatsList[i].ID.Hex()).
				Int("rank", i+1).
				Msg("Something bad happen while adding the user rank to the db")
		} else {
			log.Debug().
				Str("type", "module").
				Str("module", "functions").
				Str("function", "updateUserRanking").
				Str("userId", userStatsList[i].ID.Hex()).
				Int("rank", i+1).
				Msg("User rank has been successfully updated")
		}
	}

	log.Info().
		Str("type", "module").
		Str("module", "functions").
		Str("function", "updateUserRanking").
		Msg("Finished updating the users rank successfully")

	// TODO document new codes since 113fed51d0715ee5fd650b10b30baa97760a43c5

	return
}

// This function is used to warp the UpdateUserRanking function into a timer job
func UpdateUserRankingCron() {
	err := UpdateUserRanking()

	if err != nil {
		log.Error().
			Err(err).
			Str("type", "timer").
			Str("function", "updateUserRanking").
			Msg("Something bad happen while updating the user ranking")
	}
}
