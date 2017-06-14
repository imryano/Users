package main

import (
	"fmt"
	"github.com/imryano/Users/userPackage"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
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
	session, err := mgo.Dial(userPackage.DbUrl)
	if err == nil {
		result := &User{}
		col := session.DB(userPackage.DbName).C(userPackage.UserCol)
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
