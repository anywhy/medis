package lanucher

import (
	"github.com/anywhy/medis/pkg/core/instance"
	queue "github.com/anywhy/medis/pkg/core/lanuchqueue"
	"github.com/anywhy/medis/pkg/core/matcher"
	"github.com/anywhy/medis/pkg/models"
	"github.com/anywhy/medis/pkg/utils/log"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	"sync"
)

type OfferProcessor struct {
	mtu      sync.Mutex
	client   models.Client
	matacher matcher.OfferMatcher
}

func NewOfferProcessor(client models.Client) *OfferProcessor {

	return &OfferProcessor{
		client: client,
	}
}

func (o *OfferProcessor) ProcessOffer(driver sched.SchedulerDriver, offer *mesos.Offer) {
	_, tasks, offerOps := o.matacher.MatchOffer(offer)
	if offerOps != nil {
		o.AcceptOffers(driver, offer.Id, offerOps)
	}

	if tasks != nil {
		driver.LaunchTasks([]*mesos.OfferID{offer}, tasks, &mesos.Filters{})
	}

	o.declineOffer(driver, offer.Id)
}

func (o *OfferProcessor) declineOffer(driver sched.SchedulerDriver, offerId *mesos.OfferID) {
	_, err := driver.DeclineOffer(offerId, &mesos.Filters{})
	if err != nil {
		log.Warnf("declineOffer offer error, OfferId: %s", offerId.Value)
	}
}

func (o *OfferProcessor) AcceptOffers(driver sched.SchedulerDriver, offerId *mesos.OfferID, offerOp []*mesos.Offer_Operation) {
	stat, err := driver.AcceptOffers([]*mesos.OfferID{offerId}, offerOp, &mesos.Filters{RefuseSeconds: proto.Float64(0)})
	if err != nil {
		log.WarnErrorf(err, "AcceptOffers error, offerId:%v, dirver status: %v", offerId.GetValue(), stat.Enum().String())
	}
}

func (o *OfferProcessor) RevertTasks(ins *instance.Instance, tasks []*mesos.TaskInfo) {
	log.Warnf("Instance id: %s, RevertTasks: %v", ins.Id, tasks)
	ins.AddPossibleTask(tasks)

	// add to job queue
	if !queue.Contain(ins) {
		queue.Add(ins)
	}

}
