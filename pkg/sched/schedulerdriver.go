package sched

import (
	"github.com/CodisLabs/codis/pkg/models"
	"github.com/anywhy/medis/pkg/storage"
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

	bindingAddress := parseIP(config.GetAddress())
	driverConfig := sched.DriverConfig{
		Scheduler:      NewMedisScheduler(client),
		Framework:      fwinfo,
		Master:         proto.String(config.GetMaster()),
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
	if config.GetPrincipal() != nil {
		fwinfo.Principal = proto.String(config.GetPrincipal())
		cred = &mesos.Credential{
			Principal: proto.String(config.GetPrincipal()),
		}

		if config.GetSecret() != nil {
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

	if config.GetUser() != nil {
		fwinfo.User = config.GetUser()
	}

	if config.GetCheckpoint() != nil {
		fwinfo.Checkpoint = config.GetCheckpoint()
	}

	if config.GetWebuiUrl() != nil {
		fwinfo.WebuiUrl = config.GetWebuiUrl()
	}

	if config.GetRole() != nil {
		fwinfo.Role = config.GetRole()
	}

	if config.GetFailoverTimeout() != nil {
		fwinfo.FailoverTimeout = config.GetFailoverTimeout()
	}

	fwId := storage.GetFrameworkId(client)
	if fwId != nil {
		fwinfo.Id = fwId
	}

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
