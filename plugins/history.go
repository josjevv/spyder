package plugins

import (
	"log"

	"github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
)

const HistoryPlugin = "history"

func HistoryListener(settings *config.Conf) chan *db.Fly {
	return createListener(settings, historyHandler, HistoryPlugin)
}

func historyHandler(settings *config.Conf, fly *db.Fly) {
	if fly.Operation != "i" {
		/*
			hist, found := db.GetHistory(session, settings.MongoDb, "", fly.Collection) //fly.Id)
			if found {
				log.Print(hist)
			}
		*/
	}
	// create from and to stuff
	//
	// add history 2 db
	log.Print("history is not implemented (yet)")
}
