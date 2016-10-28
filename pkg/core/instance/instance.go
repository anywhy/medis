package instance

import . "github.com/anywhy/medis/pkg/core/task"

type Instance struct {
	Id              string          `json:"id"`
	Name            string          `json:"name"`
	Type            string          `json:"type" M-S,CODIS,C`
	TaskMap         map[string]Task `json:"taskMap,omitempty"`
	PossibleNewTask []Task		`json:"possibleNewTask",omitempty`
}

func (i *Instance) IsRuning() bool {
	return true
}
