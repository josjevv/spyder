package plugins

import (
	"log"
	"time"

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
	newHist.User = bson.ObjectIdHex(fly.GetUpdatedBy())
	//newHist.Date = fly.Timestamp
	newHist.Entity = historyEntity{}
	newHist.Entity.Ref = fly.GetCollection()
	newHist.Entity.Id = fly.GetId()

	if !fly.IsInsert() {
		histories, found := getHistories(session, settings.MongoDb, fly.GetId(), fly.GetCollection())
		if found && len(histories) > 0 {
			var entry = histories[0]
			log.Printf("Hist0.%v: %v", "Date", entry.Date)
			log.Printf("Hist0.%v: %v", "DateCreated", entry.DateCreated)
			log.Printf("Hist0.%v: %v", "DateUpdated", entry.DateUpdated)

			//Entity won't get loaded
			log.Printf("Hist0.%v: %v", "Entity.Id", entry.Entity.Id)
			log.Printf("Hist0.%v: %v", "Entity.Ref", entry.Entity.Ref)

			//id won't get loaded
			log.Printf("Hist0.%v: %v", "Id", entry.Id.Hex())
			log.Printf("Hist0.%v: %v", "Organization", entry.Organization.Hex())
			log.Printf("Hist0.%v: %v", "User", entry.User.Hex())
		}
	}

	for key, value := range fly.Object {
		if key != "update_spec" {
			log.Printf("changes[%v] = %v", key, value)
		}
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
	Id           bson.ObjectId `json:"_id,omitempty" bson:"_id"`
	Organization bson.ObjectId `json:"organization"`
	User         bson.ObjectId `json:"user"`
	Date         time.Time     `json:"date"`
	DateCreated  time.Time     `json:"date_created"`
	DateUpdated  time.Time     `json:"date_updated"`
	Entity       historyEntity `json:"entity"`
	//changes
}

type historyEntity struct {
	Ref string `json:"$ref"`
	Id  string `json:"$id"`
}

//use a map here?
// changes" : [ { "key" : "name", "from" : "function-group 0", "to" : "sjaak" }, { "key" : "update_spec", "to" : { "organization" : "54eddb6acfd92cb9eeedec6c", "updated_by" : "c0ffeeeeeeeeeeeeeeeeeeee", "app_name" : "adminapp", "timestamp" : 1424874745497 } } ]
