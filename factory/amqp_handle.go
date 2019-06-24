package factory

import (
	"fmt"
	"log"
	"time"

	"github.com/sheirys/mine/manager/journal"
)

func (f *Factory) HandleOrder(o journal.Order) error {
	if f.inProgress {
		return fmt.Errorf("factory is busy")
	}

	log.Println("accepted new task")
	f.inProgress = true
	f.Order = o
	f.Order.Accepted = true
	f.Order.AcceptedAt = time.Now()
	f.publishState()

	if err := f.Process(); err != nil {
		log.Printf("failed to process: %s\n", err)
	}

	f.Order.Finished = true
	f.Order.FinishedAt = time.Now()
	f.publishState()

	f.inProgress = false
	log.Println("finished task")

	return nil
}
