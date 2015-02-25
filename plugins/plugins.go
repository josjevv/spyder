package plugins

import (
	"log"

	"github.com/changer/spyder/config"
	"github.com/changer/spyder/db"
)

type Handler func(settings *config.Conf, fly *db.Fly)

func createListener(settings *config.Conf, handler Handler, key string) chan *db.Fly {
	channel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on " + key)
		for fly := range ch {
			go handler(settings, fly)
		}
	}(channel)

	return channel
}
