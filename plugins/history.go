package plugins

import (
	"log"
	//"time"

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
	newHist.EntityId = bson.ObjectIdHex(fly.Id)
	newHist.User = fly.User()
	newHist.Timestamp = fly.Timestamp
	newHist.Organization = fly.GetOrganization()
	newHist.Operation = fly.Operation

	// if !fly.IsInsert() {
	// 	histories, found := getHistories(session, settings.MongoDb, fly.Id, fly.GetCollection())
	// 	if found && len(histories) > 0 {
	// 		var entry = histories[0]
	// 		log.Printf("Hist0.%v: %v", "Date", entry.Date)
	// 		log.Printf("Hist0.%v: %v", "DateCreated", entry.DateCreated)
	// 		log.Printf("Hist0.%v: %v", "DateUpdated", entry.DateUpdated)

	// 		//Entity won't get loaded
	// 		log.Printf("Hist0.%v: %v", "Entity.Id", entry.Entity.Id)
	// 		log.Printf("Hist0.%v: %v", "Entity.Ref", entry.Entity.Ref)

	// 		//id won't get loaded
	// 		log.Printf("Hist0.%v: %v", "Id", entry.Id.Hex())
	// 		log.Printf("Hist0.%v: %v", "Organization", entry.Organization.Hex())
	// 		log.Printf("Hist0.%v: %v", "User", entry.User.Hex())
	// 	}
	// }

	for key, value := range fly.Object {
		//if key != "update_spec" {
		log.Printf("fly.Object[%v] = %v", key, value)
		//}
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

//{"_id": 1, "operation": "u", "timestamp": 2PM , key: "x", "value": 2, from: 1}

type history struct {
	Id           bson.ObjectId       `json:"_id,omitempty" bson:"_id"`
	EntityId     bson.ObjectId       `json:"_id,omitempty" bson:"_id"`
	Organization bson.ObjectId       `json:"organization"`
	User         bson.ObjectId       `json:"user"`
	Timestamp    bson.MongoTimestamp `json:"timestamp"`
	Operation    string              `json:"operation"`
	Key          string              `json:"key"`
	Value        string              `json:"value"`
	From         string              `json:"from"`
}
