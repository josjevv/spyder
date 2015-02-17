package main

import (
	"log"

	"github.com/changer/spyder/config"
	"github.com/changer/spyder/db"
)

func main() {
	log.Println("starting spyder...")

	config := config.ReadConfig()

	session := db.GetSession(config.MongoHost)
	defer session.Close()

	db.ReadOplog(session, config)

	log.Println("exiting spyder...")
}
