package sched

import (
	"github.com/anywhy/medis/pkg/utils/log"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	. "github.com/anywhy/medis/pkg/core/lanucher"
)

type MedisScheduler struct {
	offer *OfferProcessor
}

func NewMedisScheduler() (*MedisScheduler, error) {

	return &MedisScheduler{
		offer: NewOfferProcessor(nil),
	}, nil
}

var count = 1;
/*
 * Invoked when the scheduler successfully registers with a Mesos
 * master. A unique ID (generated by the master) used for
 * distinguishing this framework from others and `MasterInfo`
 * with the ip and port of the current master are provided as arguments.
 */
func (sched *MedisScheduler) Registered(driver sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	log.Infof("MedisScheduler Registered with Master: %v, frameworkId: %v ", *masterInfo, *frameworkId)
}

/*
 * Invoked when the scheduler re-registers with a newly elected Mesos master.
 * This is only called when the scheduler has previously been registered.
 * `MasterInfo` containing the updated information about the elected master
 * is provided as an argument.
 */
func (sched *MedisScheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	log.Infof("MedisScheduler Re-Registered with Master: %v", masterInfo)
}

/*
 * Invoked when the scheduler becomes "disconnected" from the master
 * (e.g., the master fails and another is taking over).
 */
func (sched *MedisScheduler) Disconnected(driver sched.SchedulerDriver) {
	log.Warn("MedisScheduler Disconnected")
	driver.Stop(true)
}

/*
 * Invoked when resources have been offered to this framework. A
 * single offer will only contain resources from a single slave.
 * Resources associated with an offer will not be re-offered to
 * _this_ framework until either (a) this framework has rejected
 * those resources (see SchedulerDriver::launchTasks) or (b) those
 * resources have been rescinded (see Scheduler::offerRescinded).
 * Note that resources may be concurrently offered to more than one
 * framework at a time (depending on the allocator being used). In
 * that case, the first framework to launch tasks using those
 * resources will be able to use them while the other frameworks
 * will have those resources rescinded (or if a framework has
 * already launched tasks with those resources then those tasks will
 * fail with a TASK_LOST status and a message saying as much).
 */
func (sched *MedisScheduler) ResourceOffers(driver sched.SchedulerDriver, offers[] *mesos.Offer) {
	log.Infof("Resource Offers %v", offers)

	for _, offer := range offers {
		sched.offer.ProcessOffer(driver, offer)
	}

		//if (count < 3) {
		//
		//
		//	//exec := &mesos.ExecutorInfo{
		//	//	ExecutorId: util.NewExecutorID("task." + id),
		//	//	Source: proto.String("go Test"),
		//	//	Command: &mesos.CommandInfo{
		//	//		Value: proto.String("echo 'hello medis'; sleep 100"),
		//	//
		//	//	},
		//	//	Resources: []*mesos.Resource{
		//	//		util.NewScalarResource("cpus", 0.1),
		//	//		util.NewScalarResource("mem", 10),
		//	//	},
		//	//
		//	//}
		//
		//	id := uuid.New();
		//
		//	driver.LaunchTasks([]*mesos.OfferID{offer.Id}, []*mesos.TaskInfo{{
		//		Name: proto.String("test"),
		//		TaskId: &mesos.TaskID{Value: proto.String("task." + id)},
		//		SlaveId: offer.SlaveId,
		//		Command: &mesos.CommandInfo{
		//			Value: proto.String("sleep 1000"),
		//		},
		//		Resources: []*mesos.Resource{
		//			util.NewScalarResource("cpus", 0.1),
		//			util.NewScalarResource("mem", 10),
		//		},
		//	}}, &mesos.Filters{})
		//	count++
		//} else {
		//	driver.DeclineOffer(offer.Id, &mesos.Filters{})
		//}


}

/*
 * Invoked when an offer is no longer valid (e.g., the slave was
 * lost or another framework used resources in the offer). If for
 * whatever reason an offer is never rescinded (e.g., dropped
 * message, failing over framework, etc.), a framework that attempts
 * to launch tasks using an invalid offer will receive TASK_LOST
 * status updates for those tasks (see Scheduler::resourceOffers).
 */
func (sched *MedisScheduler) OfferRescinded(driver sched.SchedulerDriver, id *mesos.OfferID) {
	log.Infof("Offer '%v' rescinded.\n", *id)
}

/*
 * Invoked when the status of a task has changed (e.g., a slave is
 * lost and so the task is lost, a task finishes and an executor
 * sends a status update saying so, etc). If implicit
 * acknowledgements are being used, then returning from this
 * callback _acknowledges_ receipt of this status update! If for
 * whatever reason the scheduler aborts during this callback (or
 * the process exits) another status update will be delivered (note,
 * however, that this is currently not true if the slave sending the
 * status update is lost/fails during that time). If explicit
 * acknowledgements are in use, the scheduler must acknowledge this
 * status on the driver.
 */
func (sched *MedisScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	log.Infof("Status update: task", status.TaskId.GetValue(), " is in state ", status.State.Enum().String())

}

/*
 * Invoked when an executor sends a message. These messages are best
 * effort; do not expect a framework message to be retransmitted in
 * any reliable fashion.
 */
func (sched *MedisScheduler) FrameworkMessage(driver sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, msg string) {
	log.Infof("framework message from executor %q slave %q: %q", eid, sid, msg)
}

/*
 * Invoked when a slave has been determined unreachable (e.g.,
 * machine failure, network partition). Most frameworks will need to
 * reschedule any tasks launched on this slave on a new slave.
 *
 * NOTE: This callback is not reliably delivered. If a host or
 * network failure causes messages between the master and the
 * scheduler to be dropped, this callback may not be invoked.
 */
func (sched *MedisScheduler) SlaveLost(driver sched.SchedulerDriver, id *mesos.SlaveID) {
	log.Warnf("Slave '%v' lost", *id)
}

/*
 * Invoked when an executor has exited/terminated. Note that any
 * tasks running will have TASK_LOST status updates automagically
 * generated.
 *
 * NOTE: This callback is not reliably delivered. If a host or
 * network failure causes messages between the master and the
 * scheduler to be dropped, this callback may not be invoked.
 */
func (sched *MedisScheduler) ExecutorLost(driver sched.SchedulerDriver, eid *mesos.ExecutorID, sid *mesos.SlaveID, code int) {
	log.Errorf("Executor %q lost on slave %q code %d", eid, sid, code)

}

/*
 * Invoked when there is an unrecoverable error in the scheduler or
 * scheduler driver. The driver will be aborted BEFORE invoking this
 * callback.
 */
func (sched *MedisScheduler) Error(driver sched.SchedulerDriver, err string) {
	log.Errorf("MedisScheduler received error: %v", err)
}
