package main

import (
	"log"

	"github.com/changer/spyder/config"
	"github.com/changer/spyder/db"
	"github.com/changer/spyder/plugins"
)

func useComponent(config config.Conf, component string) bool {
	value, present := config.Components[component]
	return present && value
}

func main() {
	log.Println("starting spyder...")

	settings := config.ReadConfig()

	session := db.GetSession(settings.MongoHost)
	defer session.Close()

	chans := db.FlyChans{}

	if useComponent(settings, "notifications") {
		chans = append(chans, plugins.NotificationListener(&settings))
	}

	if useComponent(settings, "history") {
		chans = append(chans, plugins.HistoryListener(&settings))
	}

	db.ReadOplog(settings, session, &chans)

	log.Println("exiting spyder...")
}
