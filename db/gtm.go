package db

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type tailOptions struct {
	After  TimestampGenerator
	Filter OpFilter
}

type OpChan chan *Fly

type OpLogEntry map[string]interface{}

type OpFilter func(*Fly) bool

type TimestampGenerator func(*mgo.Session) bson.MongoTimestamp

func ChainOpFilters(filters ...OpFilter) OpFilter {
	return func(op *Fly) bool {
		for _, filter := range filters {
			if filter(op) == false {
				return false
			}
		}
		return true
	}
}

func OpLogCollection(session *mgo.Session) *mgo.Collection {
	collection := session.DB("local").C("oplog.rs")
	return collection
}

func ParseTimestamp(timestamp bson.MongoTimestamp) (int32, int32) {
	ordinal := (timestamp << 32) >> 32
	ts := (timestamp >> 32)
	return int32(ts), int32(ordinal)
}

func LastOpTimestamp(session *mgo.Session) bson.MongoTimestamp {
	var opLog Fly
	collection := OpLogCollection(session)
	collection.Find(nil).Sort("-$natural").One(&opLog)
	return opLog.Timestamp
}

func GetOpLogQuery(session *mgo.Session, after bson.MongoTimestamp) *mgo.Query {
	query := bson.M{"ts": bson.M{"$gt": after}}
	collection := OpLogCollection(session)
	return collection.Find(query).LogReplay().Sort("$natural")
}

func TailOps(session *mgo.Session, channel OpChan,
	errChan chan error, timeout string, options *tailOptions) error {
	s := session.Copy()
	defer s.Close()
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		errChan <- err
		return err
	}
	if options.After == nil {
		options.After = LastOpTimestamp
	}
	currTimestamp := options.After(s)
	iter := GetOpLogQuery(s, currTimestamp).Tail(duration)
	for {
		var entry Fly
		for iter.Next(&entry) {
			err := entry.ParseEntry()

			if err != nil {
				log.Println(err)
				continue
			}

			if entry.Id != "" {
				if options.Filter == nil || options.Filter(&entry) {
					channel <- &entry
				}
			}

			currTimestamp = entry.Timestamp
		}
		if err = iter.Close(); err != nil {
			errChan <- err
			return err
		}
		if iter.Timeout() {
			continue
		}
		iter = GetOpLogQuery(s, currTimestamp).Tail(duration)
	}

	return nil
}

func tail(session *mgo.Session, options *tailOptions) (OpChan, chan error) {
	outErr := make(chan error, 20)
	outOp := make(OpChan, 20)
	go TailOps(session, outOp, outErr, "100s", options)
	return outOp, outErr
}
