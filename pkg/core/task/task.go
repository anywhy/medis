package task

import mesos "github.com/mesos/mesos-go/mesosproto"

type Task struct {
	InstanceId string         `json:"instanceId"`
	Id         string         `json:"id"`
	Type       string         `json:"type"`
	Stauts     string         `json:"status"`
	Ip         string         `json:"ip"`
	Port       int            `json:"port"`
	TaskInfo   mesos.TaskInfo `json:"taskInfo"`
	Dependent  Task           `json:"dependent"`
}
