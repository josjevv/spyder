package db

import "github.com/rwynn/gtm"

type Fly struct {
	Id        interface{}
	Operation string
	Data      map[string]interface{}
	Timestamp int64
}

func createFly(op *gtm.Op) *Fly {
	fly := Fly{}

	fly.Id = op.Id
	fly.Operation = op.Operation
	fly.Data = op.Data
	fly.Timestamp = int64(op.Timestamp)

	return &fly
}

type FlyChans []chan *Fly
