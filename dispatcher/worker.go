package dispatcher

import (
	"context"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type Task func(TaskResponse)

type TaskResponse struct {
	Payload       []byte
	correlationID string
	replyTo       string
	app           *Worker
}

type Worker struct {
	Address string

	conn    *amqp.Connection
	ch      *amqp.Channel
	consume <-chan amqp.Delivery

	wg     *sync.WaitGroup
	cancel context.CancelFunc
	ctx    context.Context

	tasks map[string]Task
}

func (w *Worker) Register(task string, callback Task) {
	w.tasks[task] = callback
}

func (w *Worker) Listen() {
	w.wg = &sync.WaitGroup{}
	w.ctx, w.cancel = context.WithCancel(context.Background())

	var err error
	if w.conn, err = amqp.Dial(w.Address); err != nil {
		log.Fatalf("failed to connect rabbit: %s\n", err)
		return
	}
	if w.ch, err = w.conn.Channel(); err != nil {
		log.Fatalf("failed to initiate channel: %s\n", err)
		return
	}
	if _, err = w.ch.QueueDeclare("factory-tasks", false, false, false, false, nil); err != nil {
		log.Fatalf("failed to declare queue: %s\n", err)
		return
	}
	if w.consume, err = w.ch.Consume("factory-tasks", "", true, false, false, false, nil); err != nil {
		log.Fatalf("failed to consume: %s\n", err)
		return
	}

	go func() {
		log.Println("listening routine ... ")
		w.wg.Add(1)
		for {
			select {
			case deliver, ok := <-w.consume:
				if !ok {
					log.Println("disconnected w.consume")
					w.wg.Done()
					return
				}

				log.Printf("received delivery: %s\n", string(deliver.Body))
				continue

				if deliver.Headers != nil {
					log.Printf("received task: %s\n", deliver.Headers["action"])
					resp := w.prepareResponse(deliver)
					w.tasks[deliver.Headers["action"].(string)](*resp)
				}
			case <-w.ctx.Done():
				w.wg.Done()
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.cancel()
	w.wg.Wait()
}

func (w *Worker) prepareResponse(d amqp.Delivery) *TaskResponse {
	return &TaskResponse{
		correlationID: d.CorrelationId,
		replyTo:       d.ReplyTo,
		Payload:       d.Body,
	}
}
