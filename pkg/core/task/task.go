package task

import mesos "github.com/mesos/mesos-go/mesosproto"

type Task struct {
	InstanceId string          `json:"instanceId"`
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	Cpus       float64         `json:"cpus"`
	Mem        float64         `json:"mem"`
	Disk       float64         `json:"disk"`
	Command    string          `json:"cmd"`
	Volumes    []*mesos.Volume `json:"volumes,omitempty"`
	Ports      []*mesos.Port   `json:"port_mappings,omitempty"`
	TaskInfo   *mesos.TaskInfo `json:"taskInfo"`
	Type       string          `json:"type"`
	Stauts     string          `json:"status"`
	Ip         string          `json:"ip"`
	Dependent  *Task           `json:"dependent"`
}
