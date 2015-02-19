package db

import (
	"log"

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

func ReadOplog(session *mgo.Session, config config.Conf) chan *gtm.Op {
	var err error
	var logChannel = make(chan *gtm.Op)

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
				if useComponent(config, "history") {
					//historyChannel <- op
				}
				if useComponent(config, "notification") {
					//notificationChannel <- op
				}
				logChannel <- op
			}
		}
	}()

	return logChannel
}

//structure of Op
// type Op struct {
//     Id        interface{}
//     Operation string
//     Namespace string
//     Data      map[string]interface{}
//     Timestamp bson.MongoTimestamp
// }
