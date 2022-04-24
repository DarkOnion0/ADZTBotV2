package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// This is the datastructures of every mongodb record in the userInfo collection
type UserRecordFetch struct {
	ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Userid string
}
type UserRecordSend struct {
	//ID     primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Userid string
}

type UserInfoFetch struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Posts       []PostRecordFetchT
	GlobalScore int
}
