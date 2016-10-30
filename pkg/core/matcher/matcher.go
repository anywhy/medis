package matcher

import  mesos "github.com/mesos/mesos-go/mesosproto"

type Matcher interface {
	MatchOffer(offer *mesos.Offer) []*mesos.Offer_Operation
}


