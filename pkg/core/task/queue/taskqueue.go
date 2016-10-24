package queue

import (
	"sync"
	"container/list"
	mesos "github.com/mesos/mesos-go/mesosproto"
)

type TaskQueue struct {
	mut sync.Mutex
	list *list.List
}

func NewTaskQueue() *TaskQueue  {
	return &TaskQueue{
		list: list.New(),
	}
}

func (t *TaskQueue) Add(task *mesos.TaskInfo) {
	t.mut.Lock()
	defer t.mut.Unlock()
	t.list.PushBack(task)
}

func (t *TaskQueue) Size() int {
	t.mut.Lock()
	defer t.mut.Unlock()
	return t.list.Len()
}

func (t *TaskQueue) Front() *mesos.TaskInfo {
	t.mut.Lock()
	defer t.mut.Unlock()

	element := t.list.Front();
	if (element == nil) {
		return nil
	}

	return element.Value.(*mesos.TaskInfo)
}

func (t *TaskQueue) PopFront() *mesos.TaskInfo  {
	t.mut.Lock()
	defer t.mut.Unlock()

	element := t.list.Front();
	if (element == nil) {
		return nil
	}

	return t.list.Remove(element).(*mesos.TaskInfo)
}

func (t *TaskQueue) IsEmpty() bool  {
	t.mut.Lock()
	defer t.mut.Unlock()
	return t.list.Len() == 0
}
