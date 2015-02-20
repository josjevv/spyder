package main

import (
	"log"

	"github.com/changer/spyder/config"
	"github.com/changer/spyder/db"
)

func useComponent(config config.Conf, component string) bool {
	_, present := config.Components[component]
	return present
}

func main() {
	log.Println("starting spyder...")

	settings := config.ReadConfig()

	session := db.GetSession(settings.MongoHost)
	defer session.Close()

	chans := db.FlyChans{}

	if useComponent(settings, "notifications") {
		noticeChannel := make(chan *db.Fly)

		go func(ch <-chan *db.Fly) {
			log.Println("Waiting for a Fly on Notifications")
			for fly := range ch {
				log.Println(fly)
			}
		}(noticeChannel)

		chans = append(chans, noticeChannel)
	}

	if useComponent(settings, "history") {
		historyChannel := make(chan *db.Fly)

		go func(ch <-chan *db.Fly) {
			log.Println("Waiting for a Fly on History")
			for fly := range ch {
				log.Println(fly)
			}
		}(historyChannel)

		chans = append(chans, historyChannel)
	}

	db.ReadOplog(settings, session, &chans)

	log.Println("exiting spyder...")
}
