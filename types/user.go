package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// This is the data structures of every mongodb record in the userInfo collection
type UserRecordFetch struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	UserRecordSend
}

// This is the data structures of every mongodb record in the userInfo collection (without the ID field to avoid override)
type UserRecordSend struct {
	Userid string
	Rank   int
}

// This represent all the user info/data from the whole db
type UserInfo struct {
	GlobalScore int
	UserInfoFetch
}

// This represent all the user info/data from the whole db (without the ID field to avoid override)
type UserInfoFetch struct {
	Posts []PostRecordFetchT
	UserRecordFetch
}

type UserInfoList []UserInfoFetch
