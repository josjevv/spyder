package plugins

import (
	"log"

	"github.com/changer/spyder/config"
	"github.com/changer/spyder/db"
)

type flyCallback func(fly *db.Fly)

const missingSetting = "Ignoring: [%v] in [%v]->[%v]. Disbaled in settings or no handler was registered."
const missingParams = "Ignoring: [%v] in [%v]->[%v]. Missing Organization [%v] or AppName [%v]"

var NotifyRegistry = map[string]flyCallback{}

func NotificationListener(settings *config.Conf) chan *db.Fly {
	channel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on Notifications")

		for fly := range ch {
			collection := fly.GetCollection()

			ok, present := settings.Notifications[collection]
			if !present || ok {
				callback, ok := NotifyRegistry[collection]
				if ok {
					go newNotification(settings, fly, callback)
					continue
				}
			}

			log.Printf(missingSetting, fly.Operation, fly.GetDatabase(), collection)
		}
	}(channel)

	return channel
}

func newNotification(settings *config.Conf, fly *db.Fly, callback flyCallback) {

	organization := fly.GetOrganization()
	app := fly.GetAppname()

	if app == "" || organization == "" {
		log.Printf(missingParams, fly.Operation, fly.GetDatabase(), fly.GetCollection(), organization, app)
		return
	}

	callback(fly)
}
