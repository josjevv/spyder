package logger

import (
	"log"
	"sync"

	"github.com/rwynn/gtm"
)

func Handle(oplog <-chan *gtm.Op) {
	var wg sync.WaitGroup

	for item := range oplog {
		wg.Add(1)
		go func(op *gtm.Op) {
			log.Printf(`<logger.go> Got op <%v> for object <%v>
			   in database <%v>
			   and collection <%v>
			   and data <%v>
			   and timestamp <%v>`,
				op.Operation, op.Id, op.GetDatabase(),
				op.GetCollection(), op.Data, op.Timestamp)
			wg.Done()
		}(item)
	}

	wg.Wait()
}
