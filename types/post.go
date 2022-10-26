package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// This represent the post document data structure in the db when the post is saved/edited
// (without the ID field to avoid override)
type PostRecordSendT struct {
	// The kind of the post -> music or video
	Type string
	Url  string
	// The db id of the poster
	User     primitive.ObjectID
	VoteList []PostVote
}

// This represent the post document data structure in the db when the post is fetched
type PostRecordFetchT struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	PostRecordSendT
}

// This represent a single post vote
type PostVote struct {
	User primitive.ObjectID
	Vote string
}
