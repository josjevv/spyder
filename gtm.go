package spyder

import (
	"github.com/bulletind/spyder/log"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const BASE = 10

var LAST_TIMESTAMP = []byte("last_timestamp")

func MakeNameSpace() []string {
	return []string{}
}

type TailOptions struct {
	After  TimestampGenerator
	Filter *bson.M
}

func (self *Progress) tx(db *bolt.DB) *bolt.Tx {
	tx, err := db.Begin(true)
	if err != nil {
		log.Error(err)
		return nil
	}

	return tx
}

func (self *Progress) commit(tx *bolt.Tx) {
	// Commit the transaction and check for error.
	log.Debug("Committing the changes")

	if err := tx.Commit(); err != nil {
		log.Error("Rolling back transaction", err)
		tx.Rollback()
	}
}

type OpChan chan Fly

type OpLogEntry map[string]interface{}

type OpFilter mgo.Query

type TimestampGenerator func(*mgo.Session) bson.MongoTimestamp

func OpLogCollection(session *mgo.Session) *mgo.Collection {
	collection := session.DB("local").C("oplog.rs")
	return collection
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

func tailOps(session *mgo.Session, progress *Progress, channel OpChan,
	timeout string, options *TailOptions,
) error {

	stateDB, err := bolt.Open(progress.Path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer stateDB.Close()

	s := session.Copy()
	defer s.Close()

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		return err
	}

	if options.After == nil {
		options.After = func(session *mgo.Session) bson.MongoTimestamp {

			tx := progress.tx(stateDB)
			bucket := tx.Bucket(LAST_TIMESTAMP)
			defer progress.commit(tx)

			v := bucket.Get(progress.Ident)

			timestamp, err := strconv.ParseInt(string(v), BASE, 64)

			if err != nil || timestamp == 0 {
				log.Error("Cannot convert saved timestamp", err)
				return LastOpTimestamp(s)
			}

			log.Info("Restored Last Timestamp", timestamp)
			return bson.MongoTimestamp(timestamp)
		}
	}

	currTimestamp := options.After(s)

	iter := GetOpLogQuery(s, currTimestamp, options.Filter).Tail(duration)
	for {

		entry := Fly{}

		for iter.Next(&entry) {
			currTimestamp = entry.Timestamp

			log.Debug("Saving Current Timestamp", currTimestamp)
			tx := progress.tx(stateDB)
			bucket := tx.Bucket(LAST_TIMESTAMP)
			bucket.Put(progress.Ident, []byte(strconv.FormatInt(int64(currTimestamp), BASE)))
			progress.commit(tx)

			err := entry.ParseEntry()

			if err != nil {
				log.Warn(err)
				continue
			}

			channel <- entry

		}

		if err = iter.Close(); err != nil {
			panic(err)
		}

		if iter.Timeout() {
			continue
		}

		iter = GetOpLogQuery(s, currTimestamp, options.Filter).Tail(duration)
	}

	iter.Close()
	return nil
}

func Tail(session *mgo.Session, progress *Progress, options *TailOptions) OpChan {
	outOp := make(OpChan, 20)

	stateDB, err := bolt.Open(progress.Path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer stateDB.Close()

	tx := progress.tx(stateDB)

	// Use the transaction...
	_, err = tx.CreateBucketIfNotExists(LAST_TIMESTAMP)
	if err != nil {
		panic(err)
	}

	progress.commit(tx)

	go tailOps(session, progress, outOp, "100s", options)
	return outOp
}
