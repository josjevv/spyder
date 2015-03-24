package db

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TailOptions struct {
	After  TimestampGenerator
	Filter *bson.M
}

type OpChan chan *Fly

type OpLogEntry map[string]interface{}

type OpFilter mgo.Query

type TimestampGenerator func(*mgo.Session) bson.MongoTimestamp

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

func GetOpLogQuery(session *mgo.Session, after bson.MongoTimestamp, filter *bson.M) *mgo.Query {
	query := bson.M{"ts": bson.M{"$gt": after}}

	if filter != nil {
		for k, v := range *filter {
			query[k] = v
		}
	}

	collection := OpLogCollection(session)
	return collection.Find(query).LogReplay().Sort("$natural")
}

func tailOps(session *mgo.Session, channel OpChan,
	errChan chan error, timeout string, options *TailOptions) error {

	s := session.Copy()
	defer s.Close()

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		errChan <- err
		return err
	}

	//TODO: Add Logic to save & find last executed timestamp

	if options.After == nil {
		options.After = LastOpTimestamp
	}

	currTimestamp := options.After(s)
	iter := GetOpLogQuery(s, currTimestamp, options.Filter).Tail(duration)
	for {
		var entry Fly
		for iter.Next(&entry) {
			err := entry.ParseEntry()

			if err != nil {
				errChan <- err
				log.Println(err)
				continue
			}

			channel <- &entry

			currTimestamp = entry.Timestamp
		}

		if err = iter.Close(); err != nil {
			errChan <- err
			return err
		}

		if iter.Timeout() {
			continue
		}

		iter = GetOpLogQuery(s, currTimestamp, options.Filter).Tail(duration)
	}

	return nil
}

func Tail(session *mgo.Session, options *TailOptions) (OpChan, chan error) {
	outErr := make(chan error, 20)
	outOp := make(OpChan, 20)
	go tailOps(session, outOp, outErr, "100s", options)
	return outOp, outErr
}
