package db

import (
	"log"

	"github.com/changer/spyder/config"
	"github.com/rwynn/gtm"
	"labix.org/v2/mgo"
)

func ReadOplog(session *mgo.Session, config config.Conf) {
	var err error

	var getFilter func(op *gtm.Op) bool
	getFilter = func(op *gtm.Op) bool {
		return op.GetDatabase() == config.MongoDb &&
			config.HasAssociation(op.GetCollection())
	}

	ops, errs := gtm.Tail(session, &gtm.Options{nil, getFilter})
	// Tail returns 2 channels - one for events and one for errors
	for {
		// loop forever receiving events
		select {
		case err = <-errs:
			// handle errors
			log.Println(err)
		case op := <-ops:
			if config.HasComponent("history") {
				//historyChannel <- op
			}
			if config.HasComponent("notification") {
				//notificationChannel <- op
			}
			//logChannel <- op

			log.Printf(`Got op <%v> for object <%v>
			   in database <%v>
			   and collection <%v>
			   and data <%v>
			   and timestamp <%v>`,
				op.Operation, op.Id, op.GetDatabase(),
				op.GetCollection(), op.Data, op.Timestamp)
		}
	}
}
