package spyder

import (
	"github.com/bulletind/spyder/log"

	"gopkg.in/mgo.v2"
)

func GetSession(connString string) *mgo.Session {
	log.Debug("Attempting connecting to:", connString)
	// get a mgo session
	session, err := mgo.Dial(connString)
	if err != nil {
		panic(err)
	}

	log.Info("spyder connected to db:", connString)
	session.SetMode(mgo.Monotonic, true)
	return session
}

type Progress struct {
	Path  string
	Ident []byte
}
