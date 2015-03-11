package db

import (
	"log"

	config "github.com/changer/spyder/config"
	"gopkg.in/mgo.v2"
)

func getFilter(settings config.Conf) func(op *Fly) bool {
	return func(op *Fly) bool {
		return op.GetDatabase() == settings.MongoDb &&
			useAssociation(settings, op.GetCollection())
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

	ops, errs := tail(session, &tailOptions{nil, getFilter(settings)})
	// Tail returns 2 channels - one for events and one for errors
	func() {
		for {
			// loop forever receiving events
			select {
			case err = <-errs:
				// handle errors
				log.Println(err)

			case fly := <-ops:
				for i := range *channels {
					(*channels)[i] <- fly
				}
			}
		}
	}()
}
