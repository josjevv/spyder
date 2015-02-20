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
