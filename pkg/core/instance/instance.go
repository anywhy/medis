package instance

import . "github.com/anywhy/medis/pkg/core/task"

type Instance struct {
	Id              string          `json:"id"`
	Name            string          `json:"name"`
	Type            string          `json:"type" M-S,CODIS,C`
	TaskMap         map[string]Task `json:"taskMap,omitempty"`
	PossibleNewTask []Task          `json:"possibleNewTask",omitempty`
	Lanuched        int             `json:"lanuched"`
}

func (i *Instance) IsRuning() bool {
	return true
}

func (i *Instance) IsLanuched() bool {
	return i.Lanuched == len(i.TaskMap)
}

func (i *Instance) GetTasks() []Task  {
	list := make([]Task, 0, len(i.TaskMap))

	for item := range i.TaskMap {
		list = append(list, item)
	}

	return list
}

func (i *Instance) GetPossibleTask() []Task  {
	list := make([]Task, 0, len(i.PossibleNewTask))

	for item := range i.PossibleNewTask {
		list = append(list, item)
	}

	return list
}
