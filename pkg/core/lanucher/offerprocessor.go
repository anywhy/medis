package lanucher

import (
	"github.com/anywhy/medis/pkg/core/task"
	"github.com/anywhy/medis/pkg/core/task/queue"
	"github.com/anywhy/medis/pkg/utils/log"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	"sync"
	"github.com/anywhy/medis/pkg/modules"
)

type OfferProcessor struct {
	mtu sync.Mutex
	client *modules.Client
}

func NewOfferProcessor(client *modules.Client) *OfferProcessor {

	return &OfferProcessor{
		client: client,
	}
}

func (o *OfferProcessor) ProcessOffer(driver sched.SchedulerDriver, offer *mesos.Offer) {

}

func (o *OfferProcessor) declineOffer(driver sched.SchedulerDriver, offerId *mesos.OfferID) {
	_, err := driver.DeclineOffer(offerId, &mesos.Filters{})
	if (err != nil) {
		log.Warnf("declineOffer offer error, OfferId: %s", offerId.Value)
	}
}
