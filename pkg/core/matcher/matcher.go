package matcher

import (
	"github.com/anywhy/medis/pkg/core/instance"
	"github.com/anywhy/medis/pkg/core/task"
	mesos "github.com/mesos/mesos-go/mesosproto"
)

type OfferMatcher interface {
	MatchOffer(offer *mesos.Offer) (instance instance.Instance, tasks []*mesos.TaskInfo, offerOps []*mesos.Offer_Operation)
}
