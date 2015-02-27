package db

import (
	"errors"
	"strings"

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

func (this *Fly) GetOrganization() string {
	organization := this.updateSpec["organization"]
	if organization == nil {
		return ""
	}

	return organization.(string)
}

func (this *Fly) GetAppname() string {
	app := this.updateSpec["app_name"]
	if app == nil {
		return ""
	}

	return app.(string)
}

func (this *Fly) GetUpdatedBy() interface{} {
	user := this.updateSpec["user"]
	if user == nil {
		return ""
	}

	return user.(string)
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

func (this *Fly) ParseEntry() (err error) {
	// only parse inserts, deletes, and updates
	var _id interface{}

	if this.IsInsert() || this.IsDelete() || this.IsUpdate() {
		if this.IsUpdate() {
			_id = this.Query["_id"]
		} else {
			_id = this.Object["_id"]
		}
	} else {
		return
	}

	switch _id := _id.(type) {
	case string:
		this.Id = _id
	case bson.ObjectId:
		this.Id = _id.Hex()
	default:
		err = errors.New("Unrecognized _id")
		return
	}

	var opBson bson.M
	var ok bool

	if this.IsInsert() {
		opBson = this.Object
	} else if this.IsUpdate() {
		setOp := this.Object["$set"]
		if setOp == nil {
			err = errors.New("Cannot find $set in Update operator")
			return
		}

		opBson, ok = setOp.(bson.M)
		if !ok {
			err = errors.New("Cannot type assert $set")
			return
		}
	} else if this.IsDelete() {
		return
	}

	update_spec := opBson["update_spec"]
	if update_spec == nil {
		err = errors.New("Cannot find update_spec in OpLog")
		return
	}

	this.updateSpec = update_spec.(bson.M)
	return
}

type FlyChans []chan *Fly
