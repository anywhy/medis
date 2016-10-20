package main

import (
	"os"
	. "github.com/anywhy/medis/pkg/sched"

	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	"github.com/anywhy/medis/pkg/utils/log"
	"github.com/anywhy/medis/pkg/core/worker"
)

func main() {
	// Framework
	fwinfo := &mesos.FrameworkInfo{
		Name: proto.String("medis"),
		//Principal: proto.String("root"),
		User: proto.String("root"),
		Checkpoint: proto.Bool(true),
	}

	scheduler, _ := NewMedisScheduler(&worker.TaskWorker{})

	// Scheduler Driver
	config := sched.DriverConfig{
		Scheduler:      scheduler,
		Framework:      fwinfo,
		Master:         "zk://localhost:2181/mesos",
		Credential: (*mesos.Credential)(nil),
	}

	driver, err := sched.NewMesosSchedulerDriver(config)

	if err != nil {
		log.Errorf("Unable to create a MedisSchedulerDriver: %v\n", err.Error())
		os.Exit(-3)
	}

	if stat, err := driver.Run(); err != nil {
		log.Errorf("Framework stopped with status %s and error: %s\n", stat.String(), err.Error())
		os.Exit(-4)
	}


}
