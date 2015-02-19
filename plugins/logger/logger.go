package logger

import (
	"log"
	"sync"

	"github.com/changer/spyder/db"
	//"github.com/rwynn/gtm"
)

func Handle(oplog <-chan db.Fly) {
	var wg sync.WaitGroup

	for item := range oplog {
		wg.Add(1)
		go func(op db.Fly) {
			// log.Printf(`Got op <%v> for object <%v>
			//    in database <%v>
			//    and collection <%v>
			//    and data <%v>
			//    and timestamp <%v> <%v>`,
			// 	op.Operation, op.Id, op.GetDatabase(),
			// 	op.GetCollection(), op.Data, op.Timestamp, op.Namespace)

			log.Printf(`Got op <%v> for object <%v>
			   and data <%v>, appname <%v>`,
				op.Operation, op.Id, op.Data, op.AppName)

			wg.Done()
		}(item)
	}

	wg.Wait()
}
