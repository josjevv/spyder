package plugins

import (
	"log"
	//"time"

	"gopkg.in/mgo.v2"

	config "github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
	bson "gopkg.in/mgo.v2/bson"
)

var BLACKLIST []string = []string{"__v", "date_created", "date_updated", "update_spec"}

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
	log.Printf("fly.Query = %v", fly.Query)

	set := fly.Object["$set"]
	if set != nil {
		var setMap = set.(bson.M)
		for key, value := range setMap {
			if !stringInSlice(key, BLACKLIST) {
				log.Printf("fly.Object['$set'][%v] = %v", key, value)
				hist := createNewHistory(fly, key, value.(string))
				hist.From = getHistoricValue(session, settings.MongoDb, fly.Id, key)

				log.Println(hist)
				insertHistory(session, settings.MongoDb, &hist)
			}
		}
	}
	// add history 2 db
	log.Print("history is not implemented (yet)")
}

func getHistoricValue(session *mgo.Session, dbName string, id string, key string) string {
	var result history
	c := session.DB(dbName).C("shared.history")

	log.Printf("entityid %v key %v ", id, key)
	err := c.Find(bson.M{"entity": bson.ObjectIdHex(id), "key": key}).Sort("-timestamp").One(&result)

	if err != nil {
		log.Print("History for key not found : " + err.Error())
		return ""
	}
	return result.Value
}

func insertHistory(session *mgo.Session, dbName string, hist *history) {
	c := session.DB(dbName).C("shared.history")

	err := c.Insert(hist)
	if err != nil {
		log.Println("Insert for history failed : " + err.Error())
	}
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

// func getOldValue() {
// 	if !fly.IsInsert() {
// 		histories, found := getHistories(session, settings.MongoDb, fly.Id, fly.GetCollection())
// 		if found && len(histories) > 0 {
// 			var entry = histories[0]
// 			log.Printf("Hist0.%v: %v", "Date", entry.Date)
// 			log.Printf("Hist0.%v: %v", "DateCreated", entry.DateCreated)
// 			log.Printf("Hist0.%v: %v", "DateUpdated", entry.DateUpdated)

// 			//Entity won't get loaded
// 			log.Printf("Hist0.%v: %v", "Entity.Id", entry.Entity.Id)
// 			log.Printf("Hist0.%v: %v", "Entity.Ref", entry.Entity.Ref)

// 			//id won't get loaded
// 			log.Printf("Hist0.%v: %v", "Id", entry.Id.Hex())
// 			log.Printf("Hist0.%v: %v", "Organization", entry.Organization.Hex())
// 			log.Printf("Hist0.%v: %v", "User", entry.User.Hex())
// 		}
// 	}
// }

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func createNewHistory(fly *db.Fly, key string, value string) history {
	var h = history{}
	h.Id = bson.NewObjectId()
	h.Entity = bson.ObjectIdHex(fly.Id)
	h.Organization = fly.GetOrganization()
	h.User = fly.User()
	h.Timestamp = fly.Timestamp
	h.Operation = fly.Operation

	h.Key = key
	h.Value = value
	h.From = value
	return h
}

//{"_id": 1, "operation": "u", "timestamp": 2PM , key: "x", "value": 2, from: 1}
type history struct {
	Id           bson.ObjectId       `json:"_id,omitempty" bson:"_id"`
	Entity       bson.ObjectId       `json:"entity,omitempty"`
	Organization bson.ObjectId       `json:"organization"`
	User         bson.ObjectId       `json:"user"`
	Timestamp    bson.MongoTimestamp `json:"timestamp"`
	Operation    string              `json:"operation"`
	Key          string              `json:"key"`
	Value        string              `json:"value"`
	From         string              `json:"from"`
}
