package main

import "labix.org/v2/mgo"
import "github.com/rwynn/gtm"
import "fmt"


func main() {
  fmt.Println("starting spyder...")

  session := getSession()
  defer session.Close()

  read(session)

  fmt.Println("exiting spyder...")
}

func getSession() *mgo.Session {
  // get a mgo session
  session, err := mgo.Dial("localhost")
  if err != nil {
    panic(err)
  }
  fmt.Println("spyder connected to localhost")
  session.SetMode(mgo.Monotonic, true)
  return session;
}

func read(session *mgo.Session) {
  var err error

  ops, errs := gtm.Tail(session, &gtm.Options{nil, getFilter})
  // Tail returns 2 channels - one for events and one for errors
  for {
    // loop forever receiving events
    select {
    case err = <-errs:
      // handle errors
      fmt.Println(err)
    case op:= <-ops:
      // op will be an insert, delete or update to mongo
      // you can check which by calling op.IsInsert(), op.IsDelete(), or op.IsUpdate()
      // op.Data will get you the full document for inserts and updates
      msg := fmt.Sprintf(`Got op <%v> for object <%v>
      in database <%v>
      and collection <%v>
      and data <%v>
      and timestamp <%v>`,
        op.Operation, op.Id, op.GetDatabase(),
        op.GetCollection(), op.Data, op.Timestamp)
      fmt.Println(msg) // or do something more interesting
    }
  }
}

func getFilter(op *gtm.Op) bool {
  return  op.Operation == "u" &&
          op.GetDatabase() == "safetyapps" &&
          op.GetCollection() == "shared.apps"
}
