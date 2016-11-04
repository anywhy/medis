package main

import (
	"github.com/anywhy/medis/pkg/models"
	"github.com/anywhy/medis/pkg/models/zk"
	"github.com/anywhy/medis/pkg/sched"
	"github.com/anywhy/medis/pkg/utils"
	"github.com/anywhy/medis/pkg/utils/log"
	"github.com/docopt/docopt-go"
	"time"
)

var usage = `
Usage:
	medis-scheduler (--master=MASTER_ADDR) (--zk=ZK_PATH|--etcd=ETCD_ADDR) [--http_port=HTTP_PORT] [--log=LOG-FILE] [--name=NAME] [--user=USER] [--mesos_principal=PRINCIPAL] [--check_point=CHECK_POINT] [--role=ROLE] [--failover_timeout=FAILOVER_TIMEOUT] [--secret_path=SECRET_PATH] [--auth_provider=AUTH_PROVIDER] [--log-level=LOG-LEVEL]
Options:
   --master=MASTER_ADDR			Mesos master address
   --zk=ZK_PATH				Medis framework metadata address
   --etcd=ETCD_ADDR			Medis framework metadata address
   --http_port=HTTP_PORT		edis framework start with UI http port
   --name=NAME				Medis framework registered with name
   --user=USER				Medis framework registered with user
   --role=ROLE				Medis framework registered with role default *
   --mesos_principal=MESOS_PRINCIPAL	Medis framework registered with mesos principal
   --check_point=CHECK_POINT		Medis framework used checkpint or not
   --failover_timeout=FAILOVER_TIMEOUT	Medis framework registered timeout
   --secret_path=SECRET_PATH		Medis framework registered with scret file path
   --auth_provider=AUTH_PROVIDER	Medis framework registered with auth provider default SASL
   --log=LOG-FILE			Medis framework start with log file
   --log-level=LOG-LEVEL		Medis framework start with log level default INFO
`

func main() {
	args, err := docopt.Parse(usage, nil, true, "", false)
	if err != nil {
		log.PanicError(err)
	}

	if s, ok := utils.Argument(args, "--log"); ok {
		w, err := log.NewRollingFile(s, log.DailyRolling)
		if err != nil {
			log.PanicErrorf(err, "open log file %s failed", s)
		} else {
			log.StdLog = log.New(w, "")
		}
	}
	log.SetLevel(log.LevelInfo)

	if s, ok := utils.Argument(args, "--log-level"); ok {
		if !log.SetLevelString(s) {
			log.Panicf("option --log-level = %s", s)
		}
	}

	// master addr
	master := utils.ArgumentMust(args, "--master")
	config := sched.NewConfig(master)

	if name, ok := utils.Argument(args, "--name"); ok {
		config.Name = name
	}

	if user, ok := utils.Argument(args, "--user"); ok {
		config.User = user
	}

	if role, ok := utils.Argument(args, "--role"); ok {
		config.Role = role
	}

	if principal, ok := utils.Argument(args, "--mesos_principal"); ok {
		config.Principal = principal
	}

	if checkPoint, ok := utils.ArgumentBool(args, "--check_point"); ok {
		config.Checkpoint = checkPoint
	}

	if authProvider, ok := utils.Argument(args, "--auth_provider"); ok {
		config.AuthProvider = authProvider
	}

	if secretPath, ok := utils.Argument(args, "--secret_path"); ok {
		config.Secret = secretPath
	}

	if failoverTimeout, ok := utils.Argument(args, "--failover_timeout"); ok {
		config.FailoverTimeout = failoverTimeout
	}

	if s, ok := utils.Argument(args, "--zk"); ok {
		config.CoordinatorName = "zookeeper"
		config.CoordinatorAddr = s
	}

	if s, ok := utils.Argument(args, "--etcd"); ok {
		config.CoordinatorName = "etcd"
		config.CoordinatorAddr = s
	}

	var client models.Client
	switch config.CoordinatorName {

	case "zookeeper":
		client, err = zkclient.New(config.CoordinatorAddr, time.Minute)
		if err != nil {
			log.Panicf("innalid zookeeper addr='%s'", config.CoordinatorAddr)
		}

		defer client.Close()
	case "etcd":

	default:
		log.Panicf("invalid coordinator name = '%s'", config.CoordinatorName)

	}

	driver, err := sched.NewSchedulerDriver(config, client)
	if err != nil {
		log.PanicErrorf(err, "start Medis Mesos Framework error")
	}

	if stat, err := driver.Run(); err != nil {
		log.PanicErrorf("Framework stopped with status %s and error: %s\n", stat.String(), err.Error())
	}
}
