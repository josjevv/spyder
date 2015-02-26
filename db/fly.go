package db

import (
	"log"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

type Fly struct {
	Id           bson.ObjectId
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
	log.Println(organization)
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
	user := this.updateSpec["updated_by"]
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

func (this *Fly) ParseId() {
	// only parse inserts, deletes, and updates
	if this.IsInsert() || this.IsDelete() || this.IsUpdate() {
		if this.IsUpdate() {
			this.Id = this.Query["_id"].(bson.ObjectId)
		} else {
			this.Id = this.Object["_id"].(bson.ObjectId)
		}
	}

	var opBson bson.M
	var ok bool

	if this.IsInsert() {
		opBson = this.Object
	} else if this.IsUpdate() {
		setOp := this.Object["$set"]
		if setOp == nil {
			return
		}

		opBson, ok = setOp.(bson.M)
		if !ok {
			log.Println("Cannot convert $set")
			return
		}
	}

	update_spec := opBson["update_spec"]
	if update_spec == nil {
		return
	}

	this.updateSpec = update_spec.(bson.M)
}

type FlyChans []chan *Fly
