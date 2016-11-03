package lanucherqueue

import (
	"container/list"
	"sync"
)

// job queue
var jobQueue = newJobsQueue()

type JobsQueue struct {
	mut  sync.Mutex
	list *list.List
}

func newJobsQueue() *JobsQueue {
	return &JobsQueue{
		list: list.New(),
	}
}

func Add(job interface{}) {
	jobQueue.mut.Lock()
	defer jobQueue.mut.Unlock()
	jobQueue.list.PushBack(job)
}

func Size() int {
	jobQueue.mut.Lock()
	defer jobQueue.mut.Unlock()
	return jobQueue.list.Len()
}

func Front() interface{} {
	jobQueue.mut.Lock()
	defer jobQueue.mut.Unlock()

	element := jobQueue.list.Front()
	if element == nil {
		return nil
	}

	return element.Value.(interface{})
}

func ReomveFront() interface{} {
	jobQueue.mut.Lock()
	defer jobQueue.mut.Unlock()

	element := jobQueue.list.Front()
	if element == nil {
		return nil
	}

	return jobQueue.list.Remove(element)
}

func Contain(instance interface{}) bool {
	e := jobQueue.list.Front()
	for e != nil {
		if e.Value == instance {
			return true
		} else {
			e = e.Next()
		}
	}

	return false
}

func IsEmpty() bool {
	jobQueue.mut.Lock()
	defer jobQueue.mut.Unlock()
	return jobQueue.list.Len() == 0
}
