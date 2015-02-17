package main

import (
	"log"

	"github.com/aderjaan/spyder/config"
	"github.com/aderjaan/spyder/db"
)

func main() {
	log.Println("starting spyder...")

	config := config.ReadConfig()

	session := db.GetSession(config.MongoHost)
	defer session.Close()

	db.ReadOplog(session, config)

	log.Println("exiting spyder...")
}
