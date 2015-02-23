package main

import (
	"log"

	"github.com/changer/spyder/config"
	"github.com/changer/spyder/db"
<<<<<<< HEAD
	"github.com/changer/spyder/plugins"
=======
	logger "github.com/changer/spyder/plugins/logger"
>>>>>>> some refactoring etc
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

<<<<<<< HEADwr
<<<<<<< HEAD
	if useComponent(settings, "notifications") {
		channel := plugins.NotificationListener()
		chans = append(chans, channel)
=======
	loggerChannel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on Logger")
		for fly := range ch {
			logger.Handle(fly)
=======
	var addHandler func(handler Handler, key string)
	addHandler = func(handler Handler, key string) {
		if !useComponent(settings, key) {
			return
>>>>>>> some refactoring and finding out - go
		}

		newChannel := make(chan *db.Fly)

		go func(ch <-chan *db.Fly) {
			log.Printf("Waiting for a Fly on %v", key)
			for fly := range ch {
				handler.Handle(fly, settings)
			}
		}(newChannel)

<<<<<<< HEAD
		chans = append(chans, noticeChannel)
>>>>>>> some refactoring etc
=======
		chans = append(chans, newChannel)
>>>>>>> some refactoring and finding out - go
	}

	addHandler(logger.Logger{}, "logger")
	addHandler(logger.Logger{}, "notificationds")
	addHandler(logger.Logger{}, "history")

	log.Println(chans)

	db.ReadOplog(settings, session, &chans)

	log.Println("exiting spyder...")
}

type Handler interface {
	Handle(fly *db.Fly, settings config.Conf)
}
