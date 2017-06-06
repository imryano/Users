package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

const dbUrl string = "127.0.0.1"
const dbName string = "userDB"
const userCol string = "users"

type AccessLevel int

const (
	AccessLevelUser        AccessLevel = iota
	AccessLevelContributor AccessLevel = iota
	AccessLevelAdmin       AccessLevel = iota
)

type User struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	GUID        string
	Username    string
	Password    string
	Email       string
	AccessLevel AccessLevel
	Version     int
}

func (user *User) WriteToDB() bool {
	session, err := mgo.Dial(dbUrl)
	if err == nil {
		result := &User{}
		col := session.DB(dbName).C(userCol)
		err = col.Find(bson.M{"GUID": user.GUID}).One(result)
		if err == nil {
			user.ID = result.ID
			col.UpdateId(result.ID, user)
			return true
		} else {
			col.Insert(user)
			return true
		}
	}
	return false
}
