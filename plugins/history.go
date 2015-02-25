package plugins

import (
	"log"

	"gopkg.in/mgo.v2"

	config "github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
	bson "gopkg.in/mgo.v2/bson"
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
	var newHist = history{}
	newHist.User = fly.GetUpdatedBy()
	newHist.Date = fly.Timestamp
	newHist.Entity = historyEntity{}
	newHist.Entity.Ref = fly.GetCollection()
	newHist.Entity.Id = fly.GetId()

	if !fly.IsInsert() {
		histories, found := getHistories(session, settings.MongoDb, fly.GetId(), fly.GetCollection())
		if found && len(histories) > 0 {
			log.Println(histories[0])
		}
	}

	for key, value := range fly.Object {
		log.Printf("changes[%v] = %v", key, value)
	}

	// user: { type: $.mongoose.Schema.ObjectId, ref: 'shared.User', index: true },
	//  date: { type: Date, default: Date.now },
	//  changes: [],
	//  entity: { type: Object }

	// add history 2 db
	log.Print("history is not implemented (yet)")
}

func getHistories(session *mgo.Session, dbName string, id string, collection string) ([]history, bool) {
	var result []history
	c := session.DB(dbName).C("shared.history")

	log.Printf("Id: %v , collection: %v", id, collection)
	err := c.Find(bson.M{"entity.$id": id, "entity.$ref": collection}).Sort("-date").All(&result)

	if err != nil {
		log.Print("History for key not found : " + err.Error())
		return result, false
	}
	return result, true
}

type history struct {
	Id           bson.ObjectId       `json:"id"        bson:"_id,omitempty"`
	Organization string              `json:"organization"`
	User         string              `json:"user"`
	Date         bson.MongoTimestamp `json:"date"`
	DateCreated  string              `json:"date_created"`
	DateUpdated  string              `json:"date_updated"`
	Entity       historyEntity       `json:"entity"`
	//changes
}

type historyEntity struct {
	Ref string `json:"$ref"`
	Id  string `json:"$id"`
}
