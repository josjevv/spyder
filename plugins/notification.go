package plugins

import (
	"log"

	"github.com/changer/spyder/config"
	"github.com/changer/spyder/db"
)

const missingSetting = "Ignoring: [%v] in [%v]->[%v]. Present: [%v]"
const missingParams = "Ignoring: [%v] in [%v]->[%v]. Missing Organization [%v] or AppName [%v]"

func NotificationListener(settings *config.Conf) chan *db.Fly {
	channel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on Notifications")

		for fly := range ch {
			ok, present := settings.Notifications[fly.Collection]
			//if !present || ok {
			if ok {
				go newNotification(settings, fly)
			} else {
				log.Printf(missingSetting, fly.Operation, fly.Database, fly.Collection, present)
			}
		}
	}(channel)

	return channel
}

func newNotification(settings *config.Conf, fly *db.Fly) {
	//log.Printf("%# v", pretty.Formatter(*fly))

	organization := fly.GetOrganization()
	app := fly.GetAppname()
	user := fly.GetUpdatedBy()

	if app == "" || organization == "" {
		log.Printf(missingParams, fly.Operation, fly.Database, fly.Collection, organization, app)
		return
	}

	var action string

	switch fly.Operation {
	case "i":
		action = "inserted"
	case "u":
		action = "updated"
	case "d":
		action = "deleted"
	default:
		action = "Unknown"
	}

	evtStr := "[%v] %v in [%v] - [%v]. Organization: [%v]. User [%v]"

	log.Printf(evtStr, app, action, fly.Database, fly.Collection, organization, user)
}
