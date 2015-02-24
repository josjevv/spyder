package db

import (
	"log"

	config "github.com/changer/spyder/config"
	"github.com/rwynn/gtm"
	"labix.org/v2/mgo"
)

func getFilter(settings config.Conf) func(op *gtm.Op) bool {
	return func(op *gtm.Op) bool {
		return op.GetDatabase() == settings.MongoDb &&
			useAssociation(settings, op.GetCollection()) &&
			op.GetCollection() != "shared.history"
	}
}

func useAssociation(settings config.Conf, association string) bool {
	if _, present := settings.Associations["all"]; present {
		return true
	}
	_, present := settings.Associations[association]
	return present
}

func ReadOplog(settings config.Conf, session *mgo.Session, channels *FlyChans) {
	var err error

	ops, errs := gtm.Tail(session, &gtm.Options{nil, getFilter(settings)})
	// Tail returns 2 channels - one for events and one for errors
	func() {
		for {
			// loop forever receiving events
			select {
			case err = <-errs:
				// handle errors
				log.Println(err)

			case op := <-ops:
				fly := createFly(op)
				for i := range *channels {
					(*channels)[i] <- fly
				}
			}
		}
	}()
}
