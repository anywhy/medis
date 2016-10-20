package mesos

import (
	resource "github.com/anywhy/medis/pkg/utils/mesos/protos"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
)

func getOfferResources(offer *mesos.Offer, name string) []*mesos.Resource {
	resources := util.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
		return res.GetName() == name
	})

	return resources
}

func getOfferScalar(offer *mesos.Offer, name string) float64 {
	resources := getOfferResources(offer, name)

	value := 0.0
	for _, res := range resources {
		value += res.GetScalar().GetValue()
	}

	return value
}

func GetOfferCpu(offer *mesos.Offer) float64  {

	return getOfferScalar(offer, resource.CPUS)
}

func GetOfferMem(offer *mesos.Offer) float64 {

	return getOfferScalar(offer, resource.MEM)
}

func GetOfferDisk(offer *mesos.Offer) float64  {

	return getOfferScalar(offer, resource.DISK)
}