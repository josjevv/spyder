package spyder

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func PushScheduled(session *mgo.Session, ident string, data bson.M, Expires time.Time) {
}
