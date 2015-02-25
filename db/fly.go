package db

import (
	"github.com/rwynn/gtm"
	"labix.org/v2/mgo/bson"
)

type Fly struct {
	Id         interface{}
	Operation  string
	Collection string
	Database   string
	Data       map[string]interface{}
	Timestamp  int64
}

func createFly(op *gtm.Op) *Fly {
	fly := Fly{}

	fly.Id = op.Id
	fly.Operation = op.Operation
	fly.Collection = op.GetCollection()
	fly.Database = op.GetDatabase()
	fly.Data = op.Data
	fly.Timestamp = int64(op.Timestamp)

	return &fly
}

func (fly *Fly) GetOrganization() string {
	organization := fly.Data["organization"]
	if organization == nil {
		return ""
	}

	organizationId := organization.(bson.ObjectId)
	return organizationId.Hex()
}

func (fly *Fly) GetAppname() string {
	return fly.Data["app_name"].(string)
}

func (fly *Fly) GetUpdatedBy() interface{} {
	user := fly.Data["updated_by"]
	if user == nil {
		return ""
	}

	userId := user.(bson.ObjectId)
	return userId.Hex()
}

type FlyChans []chan *Fly
