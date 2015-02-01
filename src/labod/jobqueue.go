package main

import "log"

type JobQueue chan *Job

// Known job queues
var (
	IncomingJobs JobQueue
)

// Push returns true if the job is successfully added on the queue
// or false if the queue is full
func (queue JobQueue) Push(job *Job) bool {
	if len(queue) > 10 {
		return false
	}
	queue <- job
	return true
}

func (queue JobQueue) Runner() {
	for {
		job := <-queue
		if job == nil {
			break
		}
	}
	log.Printf("Ending job queue with %d unfinished jobs", len(queue))
}
