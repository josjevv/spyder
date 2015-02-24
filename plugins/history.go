package plugins

import (
	"log"

	mgo "labix.org/v2/mgo"

	config "github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
)

func NewHistoryHandler(settings config.Conf, session *mgo.Session) Handler {
	return func(fly *db.Fly) {
		if fly.Operation != "i" {
			hist, found := db.GetHistory(session, settings.MongoDb, "", fly.Collection) //fly.Id)
			if found {
				log.Print(hist)
			}
		}
		// create from and to stuff
		//
		// add history 2 db
		log.Print("history is not implemented (yet)")
	}
}
