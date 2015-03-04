package history

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"

	config "github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
	utils "github.com/changer/spyder/utils"
)

type HistoryHandler struct {
	Settings       *config.Conf
	Fly            *db.Fly
	RetrievedFirst bool
	FirstHist      *History
	Histories      *mgo.Collection
}

func (handler HistoryHandler) Handle(settings *config.Conf, fly *db.Fly) {
	handler.Settings = settings
	handler.Fly = fly
	var hist = CreateBasicHistory(handler.Fly)
	var setMap bson.M

	//session and collection for handler
	var session = db.GetSession(handler.Settings.MongoHost)
	defer session.Close()
	handler.Histories = session.DB(handler.Settings.MongoDb).C("shared.history")

	if !handler.Fly.IsDelete() {
		if handler.Fly.IsInsert() {
			hist.AttachUpdateSpec(handler.Fly)
			setMap = handler.Fly.Object
			hist.Doc = bson.M{}
		} else {
			setMap = handler.Fly.Object["$set"].(bson.M)
		}

		for key, value := range setMap {
			if !IsBlacklisted(handler.Settings, "blacklistfields", key) {
				if handler.Fly.IsInsert() {
					hist.Doc[key] = value
				} else {
					hist := CreateUpdateHistory(handler.Fly, key, value)
					hist.From = handler.getHistoricValue(key)
					handler.insert(&hist)
				}
			}
		}
	}
	if !handler.Fly.IsUpdate() {
		handler.insert(&hist)
	}
}

func (handler *HistoryHandler) insert(hist *History) {
	if handler.Fly.GetCollection() == "shared.history" {
		log.Println("History showing up in fly handler. How....?")
		return
	}

	err := handler.Histories.Insert(hist)
	if err != nil {
		log.Printf("Insert for history <%v> failed : %v", hist, err.Error())
	}
}

func (handler *HistoryHandler) getHistoricValue(key string) interface{} {
	var hist History
	var lastValue interface{}

	err := handler.Histories.Find(bson.M{"entity": bson.ObjectIdHex(handler.Fly.Id), "operation": "u", "key": key}).Sort("-timestamp").One(&hist)

	if err == nil {
		lastValue = hist.Value
	} else {
		if handler.loadFirstHistory() {
			if initialValue, present := handler.FirstHist.Doc[key]; present {
				lastValue = initialValue
			}
		}
	}
	return lastValue
}

func (handler *HistoryHandler) loadFirstHistory() bool {
	if !handler.RetrievedFirst {
		handler.Histories.Find(bson.M{"entity": bson.ObjectIdHex(handler.Fly.Id), "operation": "i"}).One(&handler.FirstHist)
		handler.RetrievedFirst = true
	}
	return handler.FirstHist != nil
}

func IsBlacklisted(settings *config.Conf, settingsKey string, key string) bool {
	blacklist, hasBlacklist := settings.History[settingsKey]
	if hasBlacklist {
		return utils.StringInSlice(blacklist, key)
	}
	return false
}
