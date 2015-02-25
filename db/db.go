package db

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetSession(connString string) *mgo.Session {
	// get a mgo session
	session, err := mgo.Dial(connString)
	if err != nil {
		panic(err)
	}
	log.Println("spyder connected to db")
	session.SetMode(mgo.Monotonic, true)
	return session
}

func GetHistories(session *mgo.Session, dbName string, id string, collection string) ([]History, bool) {
	var result []History
	c := session.DB(dbName).C("shared.history")

	log.Printf("Id: %v , collection: %v", id, collection)
	err := c.Find(bson.M{"entity.$id": id, "entity.$ref": collection}).Sort("-date").All(&result)

	if err != nil {
		log.Print("History for key not found : " + err.Error())
		return result, false
	}
	return result, true
}

type History struct {
	Id           bson.ObjectId `json:"id"        bson:"_id,omitempty"`
	Organization string        `json:"organization"`
	User         string        `json:"user"`
	Date         string        `json:"date"`
	DateCreated  string        `json:"date_created"`
	DateUpdated  string        `json:"date_updated"`
	Entity       HistoryEntity `json:"entity"`
	//changes
}

type HistoryEntity struct {
	Ref string `json:"$ref"`
	Id  string `json:"$id"`
}
