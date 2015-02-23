package db

import (
	"log"
	//"time"

	"github.com/rwynn/gtm"
)

type Fly struct {
	Id           interface{}
	Operation    string
	Data         map[string]interface{}
	Organization interface{}
	AppName      string
	Collection   string
	UpdatedBy    interface{}
	DateUpdated  *time.Time
}

func createFly(op *gtm.Op) *Fly {
	fly := Fly{}

	var handleKey func(key string) interface{}
	handleKey = func(key string) interface{} {
		value, present := fly.Data[key]
		if !present {
			log.Printf("Key %v not found in oplog", key)
			return ""
		}
		delete(fly.Data, key)
		return value
	}

	fly.Id = op.Id
	fly.Operation = op.Operation
	fly.Collection = op.GetCollection()
	fly.Data = op.Data
	fly.Timestamp = int64(op.Timestamp)

	delete(fly.Data, "_id")
	delete(fly.Data, "history")
	fly.Organization = handleKey("organization")
	fly.UpdatedBy = handleKey("updated_by")
	if appName, ok := handleKey("app_name").(string); ok {
		fly.AppName = appName
	}
	fly.DateUpdated = handleKey("date_updated")

	return &fly
}

type FlyChans []chan *Fly
