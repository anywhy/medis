package lanucher

import (
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	"sync"
)

type TaskLanucher interface {
	OfferResource(driver sched.SchedulerDriver, offer *mesos.Offer) error
}

type Lanucher struct {
	mtu sync.Mutex
}

func NewLanucher() *Lanucher {

	return &Lanucher{

	}
}

func (l *Lanucher) OfferResource(driver sched.SchedulerDriver, offer *mesos.Offer) error {

	return nil
}