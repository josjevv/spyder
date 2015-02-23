package db


import (
	"github.com/rwynn/gtm"
	"log"
)

type Fly struct {
	Id        interface{}
	Operation string
	Data      map[string]interface{}
	Timestamp int64
)

type Fly struct {
	Id           interface{}
	Operation    string
	Data         map[string]interface{}
	Organization interface{}
	AppName      string
	UpdatedBy    interface{}
	DateUpdated  *time.Time
	Timestamp 	int64
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
	fly.Data = op.Data
	fly.Timestamp = int64(op.Timestamp)

	delete(fly.Data, "_id")
	delete(fly.Data, "history")
	fly.Organization = handleKey("organization")
	fly.UpdatedBy = handleKey("updated_by")
	if appName, ok := handleKey("app_name").(string); ok {
		fly.AppName = appName
	}
	if dateUpdated, ok := handleKey("date_updated").(*time.Time); ok {
		fly.DateUpdated = dateUpdated
	}

	return &fly
}

type FlyChans []chan *Fly
