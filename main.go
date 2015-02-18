package main

import (
	"log"

	"github.com/changer/spyder/config"
	"github.com/changer/spyder/db"
	logger "github.com/changer/spyder/plugins/logger"
)

func main() {
	log.Println("starting spyder...")

	config := config.ReadConfig()

	session := db.GetSession(config.MongoHost)
	defer session.Close()

	logChannel := db.ReadOplog(session, config)
	logger.Handle(logChannel)

	log.Println("exiting spyder...")
}
