package db

import (
	"time"

	"github.com/rwynn/gtm"
)

type Fly struct {
	Id        interface{}
	Operation string
	Data      map[string]interface{}
	OrgId     string
	AppName   string
	UpdatedBy string
	Timestamp *time.Time
}

func createFly(op *gtm.Op) *Fly {
	fly := Fly{}

	fly.Id = op.Id
	fly.Operation = op.Operation
	fly.Data = op.Data

	return &fly
}

type FlyChans []chan *Fly
