package logger

import (
	"log"

	db "github.com/changer/spyder/db"
)

func Handle(fly *db.Fly) {
	log.Printf(`Caught a fly: %v in %v for Organization '%v', AppName '%v',
		Id '%v' by %v on %v, changes -> <%v>`,
		fly.Operation, fly.Collection, fly.Organization, fly.AppName, fly.Id,
		fly.UpdatedBy, fly.DateUpdated, fly.Data)
}
