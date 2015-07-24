## Spyder

Spyder reads from the oplog and creates a channel out of that.

## Install

```sh
$ go get github.com/bulletind/spyder
```

## Update requirements

Each update is expected to contain and `update_spec`. That one should contain the organization, app_name, user and timestamp.

```js
{
  "_id": ObjectId("5eeded5eeded5eeded5eeded"),
  "update_spec": {
    "timestamp": ISODate("2015-07-16T14:48:06.098Z"),
    "app_name": "adminapp",
    "organization": ObjectId("5eeded5eeded5eeded5eeded"),
    "user": ObjectId("5eeded5eeded5eeded5eeded")
  }
}
```

## Register and read channel

Services can register to spyder and filter for a collection within mongo.

```go
spyder.FlyRegistry["actionapp.actions"] = actionapp.Handler
```

Reading the channel can be done in this way:

```go
func readOplog(settings *config.Config, session *mgo.Session) {
	ops := spyder.Tail(
		session,
		progress(settings),
		&spyder.TailOptions{nil, getFilter(settings)},
	)

	// Tail returns 2 channels - one for events and one for errors
	func() {
		for {
			// loop forever receiving events
			fly := <-ops
			listener(settings, fly)
		}
	}()
}
```

## Fly

The channel is populated by flies. Fly is basically an oplog entry.

```go
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
```

## Enable replicaset in mongo

- close running mongo instance if needed
- restart mongo using right db paths etc using replSet

```sh
$ mongod --port 27017 --dbpath /data/db --replSet rs0
$ mongo         # Connect to mongo
> rs.initiate() # Initiate the replicaset
> rs.status()   # Check for status
```
