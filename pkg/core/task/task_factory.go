package task

import mesos "github.com/mesos/mesos-go/mesosproto"

type TaskOpFactory struct {

}

func (t *TaskOpFactory) ApplyOffer(taskInfo *mesos.TaskInfo, offer *mesos.Offer) bool  {

	return true;
}