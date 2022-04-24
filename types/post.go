package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostRecordSendT struct {
	//ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Type     string
	Url      string
	User     primitive.ObjectID
	VoteList []PostVote
}

type PostRecordFetchT struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Type     string
	Url      string
	User     primitive.ObjectID
	VoteList []PostVote
}

type PostVote struct {
	User primitive.ObjectID
	Vote string
}
