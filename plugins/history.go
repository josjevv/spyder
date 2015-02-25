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
	var newHist = db.Hist{}
	newHist.User = fly.GetUpdatedBy()
	newHist.Date = fly.Timestamp
	newHist.Entity = db.HistoryEntity{}
	newHist.Entity.Ref = fly.Collection
	newHist.Entity.Id = fly.GetId()

	if !fly.IsInsert() {
		histories, found := db.GetHistories(session, settings.MongoDb, fly.GetId(), fly.Collection)
	}

	for key, value := range fly.Data {
		changes[key] = value
	}

	// user: { type: $.mongoose.Schema.ObjectId, ref: 'shared.User', index: true },
	//  date: { type: Date, default: Date.now },
	//  changes: [],
	//  entity: { type: Object }

	// add history 2 db
	log.Print("history is not implemented (yet)")
}
