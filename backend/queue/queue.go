package queue

import (
	"sync"

	"ai-clipping-backend/models"
)

type JobQueue struct {
	Jobs  chan *models.Job
	Store map[string]*models.Job
	Mutex sync.RWMutex
}

func NewJobQueue(buffer int) *JobQueue {
	return &JobQueue{
		Jobs:  make(chan *models.Job, buffer),
		Store: make(map[string]*models.Job),
	}
}

func (q *JobQueue) Add(job *models.Job) {
	q.Mutex.Lock()
	q.Store[job.ID] = job
	q.Mutex.Unlock()

	q.Jobs <- job
}

func (q *JobQueue) Get(id string) (*models.Job, bool) {
	q.Mutex.RLock()
	defer q.Mutex.RUnlock()

	job, ok := q.Store[id]
	return job, ok
}
