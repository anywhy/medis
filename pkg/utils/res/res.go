package mesos

import (
	"github.com/anywhy/medis/pkg/utils/log"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
)

func getOfferScalar(offer *mesos.Offer, name string) float64 {
	resources := getOfferResources(offer, name)

	value := 0.0
	for _, res := range resources {
		value += res.GetScalar().GetValue()
	}

	return value
}

func getOfferResources(offer *mesos.Offer, name string) []*mesos.Resource {
	resources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
		return res.GetName() == name
	})

	return resources
}

func GetOfferCpu(offer *mesos.Offer) float64  {

	return getOfferScalar(offer, "cpus")
}

func GetOfferMem(offer *mesos.Offer) float64 {

	return getOfferScalar(offer, "mem")
}

func GetOfferDisk(offer *mesos.Offer) float64  {

	return getOfferScalar(offer, "disk")
}

func GetOfferPorts(offer *mesos.Offer) []uint64 {
	resources := getOfferResources(offer, "ports")
	log.Infof("host: %v, url:%v", offer.GetHostname(), offer.GetUrl())
	ports := []uint64{}
	for _, res := range resources {
		for _, rans := range res.Ranges.GetRange() {
			for port := rans.GetBegin(); port < rans.GetEnd() + 1; port++ {
				ports = append(ports, port)
			}
		}
	}

	return ports;
}

func LogOffers(offers []*mesos.Offer) {
	for _, offer := range offers {
		log.Infof("Received Offer <%v> with cpus=%v mem=%v disk=%v ports=%v",
			offer.Id.GetValue(),
			GetOfferCpu(offer),
			GetOfferMem(offer),
			GetOfferDisk(offer),
			GetOfferPorts(offer))
	}
}
