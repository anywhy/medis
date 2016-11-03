package lanucher

import (
	"github.com/anywhy/medis/pkg/core/instance"
	queue "github.com/anywhy/medis/pkg/core/lanuchqueue"
	"github.com/anywhy/medis/pkg/core/matcher"
	"github.com/anywhy/medis/pkg/modules"
	"github.com/anywhy/medis/pkg/utils/log"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	"sync"
)

type OfferProcessor struct {
	mtu      sync.Mutex
	client   modules.Client
	matacher matcher.OfferMatcher
}

func NewOfferProcessor(client modules.Client) *OfferProcessor {

	return &OfferProcessor{
		client: client,
	}
}

func (o *OfferProcessor) ProcessOffer(driver sched.SchedulerDriver, offer *mesos.Offer) {
	if ins, taks, offerOps := o.matacher.MatchOffer(offer); taks != nil {
		if o.AcceptOffers(driver, offer.Id, offerOps) {
			driver.LaunchTasks([]*mesos.OfferID{offer.Id}, taks, &mesos.Filters{})
			return
		}

		o.RevertTasks(ins, taks)
	}

	o.declineOffer(driver, offer)
}

func (o *OfferProcessor) declineOffer(driver sched.SchedulerDriver, offerId *mesos.OfferID) {
	_, err := driver.DeclineOffer(offerId, &mesos.Filters{})
	if err != nil {
		log.Warnf("declineOffer offer error, OfferId: %s", offerId.Value)
	}
}

func (o *OfferProcessor) AcceptOffers(driver sched.SchedulerDriver, offerId *mesos.OfferID, offerOp []*mesos.Offer_Operation) bool {
	_, err := driver.AcceptOffers([]*mesos.OfferID{offerId}, offerOp, &mesos.Filters{RefuseSeconds: 0})

	return nil == err
}

func (o *OfferProcessor) RevertTasks(ins *instance.Instance, tasks []*mesos.TaskInfo) {
	log.Warnf("Instance id: %s, RevertTasks: %v", ins.Id, tasks)
	ins.AddPossibleTask(tasks)

	// add to job queue
	if !queue.Contain(ins) {
		queue.Add(ins)
	}

}
