package spyder

import (
	"errors"
	"log"
	"strings"

	"github.com/kr/pretty"
	"gopkg.in/mgo.v2/bson"
)

type Fly struct {
	Id           string
	Timestamp    bson.MongoTimestamp "ts"
	HistoryID    int64               "h"
	MongoVersion int                 "v"
	Operation    string              "op"
	Namespace    string              "ns"
	Object       bson.M              "o"
	Query        bson.M              "o2"
	updateSpec   bson.M
}

func getId(item interface{}) (val string, err error) {

	switch _id := item.(type) {
	case string:
		val = _id
	case bson.ObjectId:
		val = _id.Hex()
	default:
		err = errors.New("Unrecognized _id")
	}

	return
}

func (this *Fly) GetOrganization() string {
	organization := this.updateSpec["organization"]
	if organization == nil {
		return ""
	}

	val, _ := getId(organization)
	return val
}

func (this *Fly) GetAppname() string {
	app := this.updateSpec["app_name"]
	if app == nil {
		return ""
	}

	val, _ := getId(app)
	return val
}

func (this *Fly) GetUser() string {
	user := this.updateSpec["user"]
	if user == nil {
		return ""
	}

	val, _ := getId(user)
	return val
}

func (this *Fly) IsInsert() bool {
	return this.Operation == "i"
}

func (this *Fly) IsUpdate() bool {
	return this.Operation == "u"
}

func (this *Fly) IsDelete() bool {
	return this.Operation == "d"
}

func (this *Fly) ParseNamespace() []string {
	return strings.SplitN(this.Namespace, ".", 2)
}

func (this *Fly) GetDatabase() string {
	return this.ParseNamespace()[0]
}

func (this *Fly) GetCollection() string {
	return this.ParseNamespace()[1]
}

func (this *Fly) ParseEntry(ignoreCollections bson.M) (err error, ignore bool) {
	// only parse inserts, deletes, and updates
	var _id interface{}

	if this.IsInsert() || this.IsDelete() || this.IsUpdate() {
		if this.IsUpdate() {
			_id = this.Query["_id"]
		} else {
			_id = this.Object["_id"]
		}
	} else {
		log.Print("Operation is neither of Insert, Update or Delete")
		return
	}

	this.Id, err = getId(_id)
	if err != nil {
		return
	}

	var opBson bson.M
	var ok bool

	if this.IsInsert() {
		opBson = this.Object
	} else if this.IsUpdate() {

		setOp := this.Object["$set"]
		if setOp == nil {
			logTrace(&this.Object)
			err = errors.New("Cannot find $set in Update operator")
			return
		}

		opBson, ok = setOp.(bson.M)
		if !ok {
			logTrace(&this.Object)
			err = errors.New("Cannot type assert $set")
			return
		}
	} else if this.IsDelete() {
		//Delete operation does not need any check.
		return
	}

	update_spec := opBson["update_spec"]
	if update_spec == nil {
		coll := ignoreCollections[this.GetCollection()]
		if coll == nil {
			logTrace(&this.Object)
			err = errors.New("Cannot find update_spec in OpLog " + this.GetCollection())
			return err, true
		} else {
			return err, false
		}
	}

	this.updateSpec = update_spec.(bson.M)
	return err, false
}

func (this *Fly) History(connString string, initial bool) []Fly {
	entity := bson.ObjectIdHex(this.Id)
	query := bson.M{}
	query["o2"] = bson.M{
		"_id": entity,
	}

	// query the insert
	if initial {
		delete(query, "o2")
		query["o._id"] = entity
	}

	history := []Fly{}
	collection := GetSession(connString).DB("local").C("oplog.rs")
	err := collection.Find(query).Sort("-ts").All(&history)
	if err != nil {
		panic(err)
	}

	return history
}

func logTrace(spec *bson.M) {
	log.Printf("%# v", pretty.Formatter(spec))
}

type FlyChans []chan *Fly
