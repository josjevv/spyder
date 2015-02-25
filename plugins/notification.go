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
			ok, present := settings.Notifications[fly.GetCollection()]
			//if !present || ok {
			if ok {
				go newNotification(settings, fly)
			} else {
				log.Printf(missingSetting, fly.Operation, fly.GetDatabase(), fly.GetCollection(), present)
			}
		}
	}(channel)

	return channel
}

func newNotification(settings *config.Conf, fly *db.Fly) {
	organization := fly.GetOrganization()
	app := fly.GetAppname()
	user := fly.GetUpdatedBy()

	if app == "" || organization == "" {
		log.Printf(missingParams, fly.Operation, fly.GetDatabase(), fly.GetCollection(), organization, app)
		return
	}

	evtStr := "[%v] %v in [%v] - [%v]. Organization: [%v]. User [%v]"

	log.Printf(evtStr, app, fly.Operation, fly.GetDatabase(), fly.GetCollection(), organization, user)
	//log.Printf("%# v", pretty.Formatter(fly.Data))
}
