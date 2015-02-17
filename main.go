package main

import (
	"log"

	"github.com/aderjaan/spyder/config"
	"github.com/aderjaan/spyder/db"
	"github.com/rwynn/gtm"
	"labix.org/v2/mgo"
)

func main() {
	log.Println("starting spyder...")

	config := config.ReadConfig()

	session := db.GetSession(config.Connstring)
	defer session.Close()

	readOplog(session, config)

	log.Println("exiting spyder...")
}

func getFilter(op *gtm.Op) bool {
	//TODO work out filtering
	return op.Operation == "u" &&
		op.GetDatabase() == "safetyapps" &&
		op.GetCollection() == "shared.apps"
}

func readOplog(session *mgo.Session, config config.Config) {
	var err error

	ops, errs := gtm.Tail(session, &gtm.Options{nil, getFilter})
	// Tail returns 2 channels - one for events and one for errors
	for {
		// loop forever receiving events
		select {
		case err = <-errs:
			// handle errors
			log.Println(err)
		case op := <-ops:
			// op will be an insert, delete or update to mongo
			// you can check which by calling op.IsInsert(), op.IsDelete(), or op.IsUpdate()
			// op.Data will get you the full document for inserts and updates
			log.Printf(`Got op <%v> for object <%v>
			   in database <%v>
			   and collection <%v>
			   and data <%v>
			   and timestamp <%v>`,
				op.Operation, op.Id, op.GetDatabase(),
				op.GetCollection(), op.Data, op.Timestamp)
			//log.Println(msg) // or do something more interesting
		}
	}
}
