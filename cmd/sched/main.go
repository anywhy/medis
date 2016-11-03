package main

import (
	"fmt"
	"github.com/anywhy/medis/pkg/utils/log"
	"github.com/docopt/docopt-go"
)

var usage = `
Usage:
	medis-scheduler (--master=MASTER_ADDR) (--zk=ZK_PATH|--etcd=ETCD_ADDR) [--name=NAME] [--user=USER] [--mesos_principal=PRINCIPAL] [--check_point=CHECK_POINT] [--role=ROLE] [--failover_timeout=FAILOVER_TIMEOUT] [--secret_path=SECRET_PATH] [--auth_provider=AUTH_PROVIDER]
Options:
   --master=MASTER_ADDR			Mesos master address
   --zk=ZK_PATH				Medis framework metadata address
   --etcd=ETCD_ADDR			Medis framework metadata address
   --name=NAME				Medis framework registered with name
   --user=USER				Medis framework registered with user
   --role=ROLE				Medis framework registered with role default *
   --mesos_principal=MESOS_PRINCIPAL	Medis framework registered with mesos principal
   --check_point=CHECK_POINT		Medis framework used checkpint or not
   --failover_timeout=FAILOVER_TIMEOUT	Medis framework registered timeout
   --secret_path=SECRET_PATH		Medis framework registered with scret file path
   --auth_provider=AUTH_PROVIDER	Medis framework registered with auth provider default SASL
`

func main() {
	args, err := docopt.Parse(usage, nil, true, "", false)
	if err != nil {
		log.PanicError(err)
	}

	fmt.Println(args)
}
