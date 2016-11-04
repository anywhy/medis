package instance

import (
	. "github.com/anywhy/medis/pkg/core/task"
	mesos "github.com/mesos/mesos-go/mesosproto"
)

type Instance struct {
	Id              string            `json:"id"`
	Name            string            `json:"name"`
	Type            string            `json:"type" M-S,CODIS,C`
	TaskMap         map[string]Task   `json:"taskMap,omitempty"`
	PossibleNewTask []*mesos.TaskInfo `json:"possibleNewTask",omitempty`
	Lanuched        int               `json:"lanuched"`
}

func NewInstance(id string, name string, typeIns string, taskMap map[string]Task) *Instance {
	return &Instance{
		Id:              id,
		Name:            name,
		Type:            typeIns,
		TaskMap:         taskMap,
		PossibleNewTask: make([]*mesos.TaskInfo, len(taskMap)),
		Lanuched:        0,
	}
}

func (i *Instance) IsRuning() bool {
	return true
}

func (i *Instance) IsLanuched() bool {
	return i.Lanuched == len(i.TaskMap)
}

func (i *Instance) GetTasks() []Task {
	list := make([]Task, 0, len(i.TaskMap))

	for _,item := range i.TaskMap {
		list = append(list, item)
	}

	return list
}

func (i *Instance) GetPossibleTask() []*mesos.TaskInfo {
	list := make([]*mesos.TaskInfo, 0, len(i.PossibleNewTask))

	for _, item := range i.PossibleNewTask {
		list = append(list, item)
	}

	return list
}

func (i *Instance) AddPossibleTask(tasks []*mesos.TaskInfo) {
	for _, item := range  tasks {
		i.PossibleNewTask = append(i.PossibleNewTask, item)
	}
}
