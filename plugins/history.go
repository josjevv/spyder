package plugins

import (
	"log"

	mgo "labix.org/v2/mgo"

	config "github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
)

func HistoryListener(settings config.Conf, session *mgo.Session) chan *db.Fly {
	channel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on History")
		for fly := range ch {
			go newHistory(fly, settings, session)
		}
	}(channel)

	return channel
}

func newHistory(fly *db.Fly, settings config.Conf, session *mgo.Session) {
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
