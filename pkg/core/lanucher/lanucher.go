package lanucher

import (
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	"sync"
	"github.com/anywhy/medis/pkg/utils/log"
	"github.com/anywhy/medis/pkg/core/task/queue"
	"github.com/anywhy/medis/pkg/core/task"
)

type Lanucher interface {
	ProcessOffer(driver sched.SchedulerDriver, offers[] *mesos.Offer)
}

type TaskLanucher struct {
	mtu    sync.Mutex
	queue  *queue.TaskQueue
	taskOp task.TaskOpFactory
}

func NewLanucher() *TaskLanucher {

	return &TaskLanucher{

	}
}

func (t *TaskLanucher) ProcessOffer(driver sched.SchedulerDriver, offer *mesos.Offer) {
	t.mtu.Lock()
	defer t.mtu.Unlock()

	task := t.queue.PopFront();
	if (task != nil) {
		// apply offer to task
		if (t.taskOp.ApplyOffer(task, offer)) {
			stat, err := driver.LaunchTasks(offer.GetId(), []*mesos.TaskInfo{task}, &mesos.Filters{})
			if (err != nil) {
				log.Warnf("Lanucher task error, task: %v, driver status: %v", task, stat)
			}
		} else {
			t.DeclineOffer(driver, offer)
		}

	} else {
		t.DeclineOffer(driver, offer)
	}
}

func (t *TaskLanucher) DeclineOffer(driver sched.SchedulerDriver, offerId *mesos.OfferID) error {
	status, err := driver.DeclineOffer(offerId, &mesos.Filters{})

	log.Warnf("declineOffer: %v", status)
	return err
}
