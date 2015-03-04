package history

import (
	db "github.com/changer/spyder/db"
	bson "gopkg.in/mgo.v2/bson"
)

func CreateBasicHistory(fly *db.Fly) History {
	var h = History{}
	h.Id = bson.NewObjectId()
	h.Entity = bson.ObjectIdHex(fly.Id)
	h.Operation = fly.Operation
	return h
}

func (h *History) AttachUpdateSpec(fly *db.Fly) {
	h.Timestamp = fly.Timestamp
	h.Organization = bson.ObjectIdHex(fly.GetOrganization())
	h.User = bson.ObjectIdHex(fly.GetUser())
}

func CreateUpdateHistory(fly *db.Fly, key string, value interface{}) History {
	var h = CreateBasicHistory(fly)
	h.AttachUpdateSpec(fly)
	h.Key = key
	h.Value = value
	return h
}

//{"_id": 1, "operation": "u", "timestamp": 2PM , key: "x", "value": 2, from: 1, doc: {} }
type History struct {
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
