package db

import (
	"log"
	"time"

	config "github.com/changer/spyder/config"
	"github.com/rwynn/gtm"
	"labix.org/v2/mgo"
)

func getFilter(config config.Conf) func(op *gtm.Op) bool {
	return func(op *gtm.Op) bool {
		return op.GetDatabase() == config.MongoDb &&
			useAssociation(config, op.GetCollection())
	}
}

func useAssociation(config config.Conf, association string) bool {
	if _, present := config.Associations["all"]; present {
		return true
	}
	_, present := config.Associations[association]
	return present
}

func useComponent(config config.Conf, component string) bool {
	_, present := config.Components[component]
	return present
}

func createFly(op *gtm.Op) Fly {
	fly := Fly{}

	fly.Id = op.Id
	fly.Operation = op.Operation
	fly.Data = op.Data

	return fly
}

func ReadOplog(session *mgo.Session, config config.Conf) (chan Fly, chan Fly) {
	var err error
	var logChannel = make(chan Fly)
	var historyChannel = make(chan Fly)
	//var notificationChannel = make(chan *gtm.Op)

	ops, errs := gtm.Tail(session, &gtm.Options{nil, getFilter(config)})
	// Tail returns 2 channels - one for events and one for errors
	go func() {
		for {
			// loop forever receiving events
			select {
			case err = <-errs:
				// handle errors
				log.Println(err)
			case op := <-ops:
				fly := createFly(op)
				fly.AppName = "log"
				logChannel <- fly
				if useComponent(config, "history") {
					fly.AppName = "hist"
					historyChannel <- fly
				}
				if useComponent(config, "notification") {
					//notificationChannel <- op
				}
			}
		}
	}()

	return logChannel, historyChannel
}

type Fly struct {
	Id          interface{}
	Operation   string
	Data        map[string]interface{}
	OrgId       string
	AppName     string
	UpdatedBy   string
	DateUpdated *time.Time
}
