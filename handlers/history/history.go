package logger

import (
	"log"

	mgo "labix.org/v2/mgo"

	config "github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
	handler "github.com/changer/spyder/handlers"
)

func GetHandler(settings config.Conf, session *mgo.Session) handler.Handler {
	return func(fly *db.Fly) {
		// get latest hist for Fly.Id

		if fly.Operation == "u" {
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
