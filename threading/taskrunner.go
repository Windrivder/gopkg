package threading

import (
	"github.com/windrivder/gopkg/container/typex"
	"github.com/windrivder/gopkg/rescue"
)

// A TaskRunner is used to control the concurrency of goroutines.
type TaskRunner struct {
	limitChan chan typex.PlaceholderType
}

// NewTaskRunner returns a TaskRunner.
func NewTaskRunner(concurrency int) *TaskRunner {
	return &TaskRunner{
		limitChan: make(chan typex.PlaceholderType, concurrency),
	}
}

// Schedule schedules a task to run under concurrency control.
func (rp *TaskRunner) Schedule(task func()) {
	rp.limitChan <- typex.Placeholder

	go func() {
		defer rescue.Recover(func() {
			<-rp.limitChan
		})

		task()
	}()
}
