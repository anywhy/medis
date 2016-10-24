package queue

import (
	"testing"
	util "github.com/mesos/mesos-go/mesosutil"
	"github.com/mesos/mesos-go/mesosproto"
	"fmt"
)

func Test_Add(t *testing.T) {
	a := NewTaskQueue()
	tt :=util.NewTaskInfo(
		"test",
		util.NewTaskID("test"),
		util.NewSlaveID("1231kjk"),
		[]*mesosproto.Resource{},
	)
	a.Add(tt)

	if (a.PopFront() != tt) {
		t.Errorf("error")
	}

	fmt.Print(a.Front())
}
