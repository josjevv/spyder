package plugins

import (
	"log"

	"github.com/changer/spyder/db"
)

type Handler func(fly *db.Fly)

func CreateListener(handler Handler, key string) chan *db.Fly {
	channel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on " + key)
		for fly := range ch {
			go handler(fly)
		}
	}(channel)

	return channel
}
