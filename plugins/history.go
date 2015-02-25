package plugins

import (
	"log"

	mgo "labix.org/v2/mgo"

	config "github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
)

func HistoryListener(settings *config.Conf, session *mgo.Session) chan *db.Fly {
	channel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on History")

		for fly := range ch {
			//ok, present := settings.Notifications[fly.Collection]
			//if !present || ok {
			//if ok {
			go historyHandler(settings, session, fly)
			//} else {
			//	log.Printf(missingSetting, fly.Operation, fly.Database, fly.Collection, present)
			//}
		}
	}(channel)

	return channel
}

func historyHandler(settings *config.Conf, session *mgo.Session, fly *db.Fly) {
	if fly.Operation != "i" {
		hist, found := db.GetHistory(session, settings.MongoDb, fly.GetId(), fly.Collection)
		if found {
			log.Print(hist)
		} else {
			log.Print("hist not found")
		}
	}
	// create from and to stuff
	//
	// add history 2 db
	log.Print("history is not implemented (yet)")
}
