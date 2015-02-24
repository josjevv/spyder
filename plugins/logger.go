package plugins

import (
	"log"

	db "github.com/changer/spyder/db"
)

func LogListener() chan *db.Fly {
	channel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on Log")
		for fly := range ch {
			go newLog(fly)
		}
	}(channel)

	return channel
}

func newLog(fly *db.Fly) {
	log.Printf(`Caught a fly: %v in %v for Organization '%v', AppName '%v',
    Id '%v' by %v on %v, changes -> <%v>`,
		fly.Operation, fly.Collection, fly.GetOrganization(), fly.GetAppname, fly.Id,
		fly.GetUpdatedBy(), fly.GetDateUpdated(), fly.Data)

	log.Println("Event generated at: ", fly.Timestamp)
}
