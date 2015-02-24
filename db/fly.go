package db

import "github.com/rwynn/gtm"

type Fly struct {
	Id         interface{}
	Operation  string
	Collection string
	Data       map[string]interface{}
	Timestamp  int64
}

func createFly(op *gtm.Op) *Fly {
	fly := Fly{}

	fly.Id = op.Id
	fly.Operation = op.Operation
	fly.Collection = op.GetCollection()
	fly.Data = op.Data
	fly.Timestamp = int64(op.Timestamp)

	return &fly
}

func (fly Fly) GetOrganization() interface{} {
	return fly.Data["organization"]
}

func (fly Fly) GetAppname() interface{} {
	return fly.Data["app_name"]
}

func (fly Fly) GetUpdatedBy() interface{} {
	return fly.Data["updated_by"]
}

func (fly Fly) GetDateUpdated() interface{} {
	return fly.Data["date_updated"]
}

type FlyChans []chan *Fly
