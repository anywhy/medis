package main

import (
	"os"
	"github.com/gogo/protobuf/proto"

	. "github.com/anywhy/medis/pkg/sched"

	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
	"github.com/anywhy/medis/pkg/utils/log"
)

func main() {

	exec := &mesos.ExecutorInfo{
		ExecutorId: util.NewExecutorID("medis"),
		Name: proto.String("Test Executor (Go)"),
		Source: proto.String("go Test"),
		Command: &mesos.CommandInfo{
			Value: proto.String("while [ true ] ; do echo 'Hello Medis' ; sleep 5 ; done"),

		},

	}

	// Framework
	fwinfo := &mesos.FrameworkInfo{
		Name: proto.String("medis"),
		Principal: proto.String("yangdx"),
	}

	scheduler, _ := NewMedisScheduler(exec)

	// Scheduler Driver
	config := sched.DriverConfig{
		Scheduler:      scheduler,
		Framework:      fwinfo,
		Master:         "zk://192.168.30.1:2181/mesos",
		Credential: (*mesos.Credential)(nil),
	}

	driver, err := sched.NewMesosSchedulerDriver(config)

	if err != nil {
		log.Errorf("Unable to create a SchedulerDriver: %v\n", err.Error())
		os.Exit(-3)
	}

	if stat, err := driver.Run(); err != nil {
		log.Errorf("Framework stopped with status %s and error: %s\n", stat.String(), err.Error())
		os.Exit(-4)
	}


}
