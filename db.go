package spyder

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetSession(connString string) *mgo.Session {
	log.Println("Attempting connecting to:", connString)
	// get a mgo session
	session, err := mgo.Dial(connString)
	if err != nil {
		panic(err)
	}

	log.Println("spyder connected to db:", connString)
	session.SetMode(mgo.Monotonic, true)
	return session
}

type Progress struct {
	Path  string
	Ident []byte
}

type Scheduled struct {
	Id      bson.ObjectId `bson:"_id"`
	Ident   string        `bson:"ident"`
	Expires time.Time     `bson:"expires"`
	Data    bson.M        `bson:"data"`
}
