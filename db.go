package spyder

import (
	"log"

	"gopkg.in/mgo.v2"
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
