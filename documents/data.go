package documents

import (
	"gopkg.in/mgo.v2"
	"log"
)

type dbConnection struct {
	Session         *mgo.Session
	DB              *mgo.Database
	UsersCollection *mgo.Collection
}

var instance *dbConnection

func Connection() *dbConnection {
	if instance == nil {
		session, err := mgo.Dial("localhost")
		if err != nil {
			log.Fatal(err)
			return nil
		}
		db := session.DB("gchat")
		instance = &dbConnection{session, db, db.C("users")}
	}
	return instance
}
