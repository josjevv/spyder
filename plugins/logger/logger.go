package logger

import (
	"log"

	config "github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
)

type Logger struct {
}

func (logger Logger) Handle(fly *db.Fly, settings config.Conf) {
	log.Printf(`Caught a fly: %v in %v for Organization '%v', AppName '%v',
		Id '%v' by %v on %v, changes -> <%v>`,
		fly.Operation, fly.Collection, fly.Organization, fly.AppName, fly.Id,
		fly.UpdatedBy, fly.DateUpdated, fly.Data)
}
