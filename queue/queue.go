package queue

import (
	"context"
	"log"
)

// The following is inspired by https://webdevstation.com/posts/simple-queue-implementation-in-golang/

// Queue holds name, list of jobs and context with cancel.
type Queue struct {
	name   string
	limit  int // The limit is how many jobs can be done every hour
	jobs   chan Job
	ctx    context.Context
	cancel context.CancelFunc
}

// NewQueue instantiates new queue.
func NewQueue(name string, limit int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())

	return &Queue{
		jobs:   make(chan Job),
		name:   name,
		limit:  limit,
		ctx:    ctx,
		cancel: cancel,
	}
}

// AddJob sends job to the channel.
func (q *Queue) AddJob(job Job) {
	go func() { q.jobs <- job }()
	log.Printf("New job %s added to %s queue", job.Name, q.name)
}
