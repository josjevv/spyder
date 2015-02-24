package plugins

import (
	"log"

	"github.com/changer/spyder/db"
)

const NotificationPlugin = "notifications"

func NotificationListener() chan *db.Fly {
	return createListener(newNotification, NotificationPlugin)
}

func newNotification(fly *db.Fly) {
	organization := fly.GetOrganization()
	app_name := fly.GetAppname()
	updated_by := fly.GetUpdatedBy()

	log.Println(organization, app_name, updated_by)

	log.Println("Event generated at: ", fly.Timestamp)
}
