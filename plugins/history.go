package plugins

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	config "github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
	utils "github.com/changer/spyder/utils"
)

var _BLACKLIST []string = []string{"__v", "date_created", "date_updated", "update_spec"}

func HistoryListener(settings *config.Conf) chan *db.Fly {
	channel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on History")

		for fly := range ch {
			//ok, present := settings.Notifications[fly.Collection]
			//if !present || ok {
			//if ok {
			go historyHandler(settings, fly)
			//} else {
			//	log.Printf(missingSetting, fly.Operation, fly.Database, fly.Collection, present)
			//}
		}
	}(channel)

	return channel
}

func historyHandler(settings *config.Conf, fly *db.Fly) {
	session := db.GetSession(settings.MongoHost)
	defer session.Close()

	if fly.IsUpdate() {
		setMap := fly.Object["$set"].(bson.M)
		for key, value := range setMap {
			if !utils.StringInSlice(key, _BLACKLIST) {
				hist := createNewHistory(fly, key, value.(string))
				hist.From = getHistoricValue(session, settings.MongoDb, fly.Id, key)
				insertHistory(session, settings.MongoDb, &hist)
			}
		}
	}

	log.Print("history is not fully implemented (yet)...")
}

func getHistoricValue(session *mgo.Session, dbName string, id string, key string) string {
	var result history
	c := session.DB(dbName).C("shared.history")
	err := c.Find(bson.M{"entity": bson.ObjectIdHex(id), "key": key}).Sort("-timestamp").One(&result)

	if err != nil {
		err := c.Find(bson.M{"entity": bson.ObjectIdHex(id), "operation": "i"}).One(&result)
		if err != nil {
			log.Printf("History for entity %v not found : %v", id, err.Error())
		}
	}
	return result.Value
}

func insertHistory(session *mgo.Session, dbName string, hist *history) {
	c := session.DB(dbName).C("shared.history")

	err := c.Insert(hist)
	if err != nil {
		log.Printf("Insert for history <%v> failed : %v", hist, err.Error())
	}
}

func createNewHistory(fly *db.Fly, key string, value string) history {
	var h = history{}
	h.Id = bson.NewObjectId()
	h.Entity = bson.ObjectIdHex(fly.Id)
	h.Organization = bson.ObjectIdHex(fly.GetOrganization())
	h.User = bson.ObjectIdHex(fly.GetUser())
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
