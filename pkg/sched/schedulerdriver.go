package sched

import (
	"github.com/anywhy/medis/pkg/models"
	"github.com/anywhy/medis/pkg/utils/log"
	"github.com/gogo/protobuf/proto"
	"github.com/mesos/mesos-go/auth"
	"github.com/mesos/mesos-go/auth/sasl"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/scheduler"
	sched "github.com/mesos/mesos-go/scheduler"
	"golang.org/x/net/context"
	"io/ioutil"
	"net"
	"os"
)

func NewSchedulerDriver(config *Config, client models.Client) (*scheduler.MesosSchedulerDriver, error) {
	fwinfo, cred := frameworkInfoBuild(config, client)

	scheduler, err := NewMedisScheduler(client)
	if err != nil {
		log.Panic("create medisScheduler error")
	}

	bindingAddress := parseIP(config.GetAddress())
	driverConfig := sched.DriverConfig{
		Scheduler:      scheduler,
		Framework:      fwinfo,
		Master:         config.GetMaster(),
		Credential:     cred,
		BindingAddress: bindingAddress,
		WithAuthContext: func(ctx context.Context) context.Context {
			ctx = auth.WithLoginProvider(ctx, config.GetAuthProvider())
			ctx = sasl.WithBindingAddress(ctx, bindingAddress)
			return ctx
		},
	}

	return sched.NewMesosSchedulerDriver(driverConfig)
}

func frameworkInfoBuild(config *Config, client models.Client) (*mesos.FrameworkInfo, *mesos.Credential) {
	// Framework info
	fwinfo := &mesos.FrameworkInfo{
		Name:       proto.String(config.GetName()),
		Checkpoint: proto.Bool(config.GetCheckpoint()),
	}

	cred := (*mesos.Credential)(nil)
	if config.GetPrincipal() != "" {
		fwinfo.Principal = proto.String(config.GetPrincipal())

		if config.GetSecret() != "" {
			cred = &mesos.Credential{
				Principal: proto.String(config.GetPrincipal()),
			}
			_, err := os.Stat(config.GetSecret())
			if err != nil {
				log.Errorf("missing secret file: ", err.Error())
			}
			secret, err := ioutil.ReadFile(config.GetSecret())
			if err != nil {
				log.Errorf("failed to read secret file: ", err.Error())
			}
			cred.Secret = proto.String(string(secret))
		}
	}

	if config.GetUser() != "" {
		fwinfo.User = proto.String(config.GetUser())
	}

	if config.GetCheckpoint() != false {
		fwinfo.Checkpoint = proto.Bool(config.GetCheckpoint())
	}

	if config.GetWebuiUrl() != "" {
		fwinfo.WebuiUrl = proto.String(config.GetWebuiUrl())
	}

	if config.GetRole() != "" {
		fwinfo.Role = proto.String(config.GetRole())
	}

	if config.GetFailoverTimeout() != 0 {
		fwinfo.FailoverTimeout = proto.Float64(config.GetFailoverTimeout())
	}

	//fwId := storage.GetFrameworkId(client)
	//if fwId != "" {
	//	fwinfo.Id = util.NewFrameworkID(fwId)
	//}

	return fwinfo, cred
}

func parseIP(address string) net.IP {
	addr, err := net.LookupIP(address)
	if err != nil {
		log.ErrorError(err)
	}
	if len(addr) < 1 {
		log.Errorf("failed to parse IP from address '%v'", address)
	}
	return addr[0]
}
