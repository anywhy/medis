package matcher

import  mesos "github.com/mesos/mesos-go/mesosproto"

type OfferMatcher interface {
	MatchOffer(offer *mesos.Offer) bool
}


