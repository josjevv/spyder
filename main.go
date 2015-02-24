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

	chans = append(chans, plugins.LogListener())

	if useComponent(settings, "notifications") {
		channel := plugins.NotificationListener()
		chans = append(chans, channel)
	}

	if useComponent(settings, "history") {
		channel := plugins.HistoryListener(settings, session)
		chans = append(chans, channel)
	}

	db.ReadOplog(settings, session, &chans)

	log.Println("exiting spyder...")
}
