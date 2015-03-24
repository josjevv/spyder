package db

import (
	"log"

	"gopkg.in/mgo.v2"
)

func GetSession(connString string) *mgo.Session {
	// get a mgo session
	session, err := mgo.Dial(connString)
	if err != nil {
		panic(err)
	}

	log.Println("spyder connected to db:" + connString)
	session.SetMode(mgo.Monotonic, true)
	return session
}
