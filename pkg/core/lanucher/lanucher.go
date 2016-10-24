package lanucher

import (
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	"sync"
	"github.com/anywhy/medis/pkg/utils/log"
	"github.com/anywhy/medis/pkg/core/task/queue"
)

type TaskLanucher interface {
	ProcessOffer(driver sched.SchedulerDriver, offers[] *mesos.Offer)
}

type Lanucher struct {
	mtu sync.Mutex
	queue *queue.TaskQueue
}

func NewLanucher() *Lanucher {

	return &Lanucher{

	}
}

func (l *Lanucher) ProcessOffer(driver sched.SchedulerDriver, offers[] *mesos.Offer) {

}

func (l *Lanucher) DeclineOffer(driver sched.SchedulerDriver, offerId *mesos.OfferID) error {
	status,err := driver.DeclineOffer(offerId, &mesos.Filters{})

	log.Warnf("declineOffer: %v", status)
	return err
}
