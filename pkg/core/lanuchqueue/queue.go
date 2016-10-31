package lanucherqueue

import (
	"container/list"
	"github.com/anywhy/medis/pkg/core/instance"
	"sync"
)

type JobsQueue struct {
	mut  sync.Mutex
	list *list.List
}

func NewJobsQueue() *JobsQueue {
	return &JobsQueue{
		list: list.New(),
	}
}

func (q *JobsQueue) Add(job *instance.Instance) {
	q.mut.Lock()
	defer q.mut.Unlock()
	q.list.PushBack(job)
}

func (q *JobsQueue) Size() int {
	q.mut.Lock()
	defer q.mut.Unlock()
	return q.list.Len()
}

func (q *JobsQueue) Front() *instance.Instance {
	q.mut.Lock()
	defer q.mut.Unlock()

	element := q.list.Front()
	if element == nil {
		return nil
	}

	return element.Value.(*instance.Instance)
}

func (q *JobsQueue) ReomveFront() *instance.Instance {
	q.mut.Lock()
	defer q.mut.Unlock()

	element := q.list.Front()
	if element == nil {
		return nil
	}

	return q.list.Remove(element)
}

func (q *JobsQueue) IsEmpty() bool {
	q.mut.Lock()
	defer q.mut.Unlock()
	return q.list.Len() == 0
}
