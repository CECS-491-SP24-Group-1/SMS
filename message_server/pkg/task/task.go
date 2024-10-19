package task

import "time"

// Defines a runnable task with start/stop actions and a periodic action.
type Task interface {
	//The time quantum for which the periodic task shall sync to.
	periodicDuration() time.Duration

	//The task to run when the scheduler starts.
	runOnStart()

	//The task to run when the scheduler stops.
	runOnStop()

	//The task to run periodically within the time quantum given by `periodicDuration()`.
	runPeriodically()
}
