package plugins

import (
	"log"

	"gopkg.in/mgo.v2/bson"

	config "github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
	utils "github.com/changer/spyder/utils"
)

func HistoryListener(settings *config.Conf) chan *db.Fly {
	channel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on History")

		for fly := range ch {
			log.Println("a one day fly?")
			if !isBlacklisted(settings, "blacklistcollections", fly.GetCollection()) {
				go historyHandler(settings, fly)
			}
		}
	}(channel)

	return channel
}

func isBlacklisted(settings *config.Conf, settingsKey string, key string) bool {
	blacklist, hasBlacklist := settings.History[settingsKey]
	if hasBlacklist {
		return utils.StringInSlice(blacklist, key)
	}
	return true
}

func historyHandler(settings *config.Conf, fly *db.Fly) {
	var hist = createBasicHistory(fly)
	var setMap bson.M

	if !fly.IsDelete() {
		if fly.IsInsert() {
			attachUpdateSpec(&hist, fly)
			setMap = fly.Object
			hist.Doc = bson.M{}
		} else {
			setMap = fly.Object["$set"].(bson.M)
		}

		for key, value := range setMap {
			if !isBlacklisted(settings, "blacklistfields", key) {
				if fly.IsInsert() {
					hist.Doc[key] = value
				} else {
					hist := createUpdateHistory(fly, key, value)
					hist.From = getHistoricValue(settings, fly.Id, key)
					insertHistory(settings, &hist)
				}
			}
		}
	}
	if !fly.IsUpdate() {
		insertHistory(settings, &hist)
	}
}

// TODO move all history object stuff to a separate file/class structure etc

// TODO refactor this so we also 'remember the first history record'
// or even retrieve everything up front in one query
func getHistoricValue(settings *config.Conf, id string, key string) interface{} {
	var hist history
	var lastValue interface{}

	session := db.GetSession(settings.MongoHost)
	defer session.Close()

	c := session.DB(settings.MongoDb).C("shared.history")
	err := c.Find(bson.M{"entity": bson.ObjectIdHex(id), "operation": "u", "key": key}).Sort("-timestamp").One(&hist)

	if err == nil {
		lastValue = hist.Value
	} else {
		err = c.Find(bson.M{"entity": bson.ObjectIdHex(id), "operation": "i"}).One(&hist)
		if err == nil {
			if initialValue, present := hist.Doc[key]; present {
				lastValue = initialValue
			}
		}
	}
	return lastValue
}

func insertHistory(settings *config.Conf, hist *history) {
	session := db.GetSession(settings.MongoHost)
	defer session.Close()

	c := session.DB(settings.MongoDb).C("shared.history")
	err := c.Insert(hist)
	if err != nil {
		log.Printf("Insert for history <%v> failed : %v", hist, err.Error())
	}
}

func createBasicHistory(fly *db.Fly) history {
	var h = history{}
	h.Id = bson.NewObjectId()
	h.Entity = bson.ObjectIdHex(fly.Id)
	h.Operation = fly.Operation
	return h
}

func attachUpdateSpec(h *history, fly *db.Fly) {
	h.Timestamp = fly.Timestamp
	h.Organization = bson.ObjectIdHex(fly.GetOrganization())
	h.User = bson.ObjectIdHex(fly.GetUser())
}

func createUpdateHistory(fly *db.Fly, key string, value interface{}) history {
	var h = createBasicHistory(fly)
	attachUpdateSpec(&h, fly)
	h.Key = key
	h.Value = value
	return h
}

//{"_id": 1, "operation": "u", "timestamp": 2PM , key: "x", "value": 2, from: 1}
type history struct {
	Id           bson.ObjectId       `json:"_id,omitempty" bson:"_id"`
	Entity       bson.ObjectId       `json:"entity"`
	Organization bson.ObjectId       `json:"organization" bson:"organization,omitempty"`
	User         bson.ObjectId       `json:"user" bson:"user,omitempty"`
	Timestamp    bson.MongoTimestamp `json:"timestamp"`
	Operation    string              `json:"operation"`
	Key          string              `json:"key, omitempty"`
	Value        interface{}         `json:"value, omitempty"`
	From         interface{}         `json:"from, omitempty"`
	Doc          bson.M              `json:"doc, omitempty"`
}
