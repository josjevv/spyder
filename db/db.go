package db

import (
	"labix.org/v2/mgo"
	"log"
)

func GetSession(connString string) *mgo.Session {
	// get a mgo session
	session, err := mgo.Dial(connString)
	if err != nil {
		panic(err)
	}
	log.Println("spyder connected to db")
	session.SetMode(mgo.Monotonic, true)
	return session
}
