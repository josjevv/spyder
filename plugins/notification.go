package plugins

import (
	"log"

	"github.com/changer/spyder/db"
)

func NotificationListener() chan *db.Fly {
	noticeChannel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on Notifications")
		for fly := range ch {
			go newEvent(fly)
		}
	}(noticeChannel)

	return noticeChannel
}

func newEvent(fly *db.Fly) {
	organization := fly.GetOrganization()
	app_name := fly.GetAppname()
	updated_by := fly.GetUpdatedBy()

	log.Println(organization, app_name, updated_by)

	log.Println("Event generated at: ", fly.Timestamp)
}
