package executor

import (
	"sync"

	"github.com/yangzhao28/jigsaw/src/common"
)

type Config struct {
	WorkerNum       int
	InputQueueSize  int
	WorkerQueueSize int
}

type job struct {
	entry common.Entry
	arg   interface{}
}

type notifier struct {
	Subject string
	Payload interface{}
}

type EventListener interface {
	Get(subject string) []common.Entry
	Add(subject string, entry common.Entry)
}

type Executor struct {
	config      Config
	queue       chan *notifier
	workerQueue chan *job
	quit        chan struct{}
	listener    EventListener
	wg          sync.WaitGroup
}

// New create new executor
func New(config Config, listener EventListener) *Executor {
	return &Executor{
		listener:    listener,
		config:      config,
		queue:       make(chan *notifier, config.InputQueueSize),
		workerQueue: make(chan *job, config.WorkerQueueSize),
		quit:        make(chan struct{}),
	}
}

func (e *Executor) serve() {
	for {
		select {
		case <-e.quit:
			return
		case job := <-e.workerQueue:
			// fmt.Printf("new job %#v %T \n", job.entry, job.entry)
			job.entry.Call(job.arg)
		}
	}
}

func (e *Executor) split() {
	for {
		select {
		case <-e.quit:
			return
		case n := <-e.queue:
			entries := e.listener.Get(n.Subject)
			for _, entry := range entries {
				e.workerQueue <- &job{
					entry: entry,
					arg:   n.Payload,
				}
			}
		}
	}
}

// Serve start serving jobs
func (e *Executor) Serve() {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.split()
	}()

	for i := 0; i < e.config.WorkerNum; i++ {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			e.serve()
		}()
	}
}

// Quit quit serve loop
func (e *Executor) Quit() {
	close(e.quit)
	e.Wait()
}

// Wait wait for workers
func (e *Executor) Wait() {
	e.wg.Wait()
}

// Enqueue put job into queue, waiting for async execution
func (e *Executor) Publish(subject string, payload interface{}) {
	e.queue <- &notifier{
		Subject: subject,
		Payload: payload,
	}
}

func (e *Executor) Register(subject string, entry common.Entry) {
	e.listener.Add(subject, entry)
}
