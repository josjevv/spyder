package plugins

import (
	"log"

	config "github.com/changer/spyder/config"
	db "github.com/changer/spyder/db"
	history "github.com/changer/spyder/plugins/history"
)

func HistoryListener(settings *config.Conf) chan *db.Fly {
	channel := make(chan *db.Fly)

	go func(ch <-chan *db.Fly) {
		log.Println("Waiting for a Fly on History")

		for fly := range ch {
			if !history.IsBlacklisted(settings, "blacklistcollections", fly.GetCollection()) {
				handler := history.HistoryHandler{}
				go handler.Handle(settings, fly)
			}
		}
	}(channel)

	return channel
}
