package task

import (
	"fmt"
	"time"
)

// Example task; implements `Task`.
type FitTeaTask struct{}

var _ Task = (*FitTeaTask)(nil) // Type assertion check to ensure compliance with `Task` interface.

func (ftt FitTeaTask) periodicDuration() time.Duration {
	return time.Second * 1
}

func (ftt FitTeaTask) runOnStart() {
	fmt.Printf("fit tea task; run start\n")
}

func (ftt FitTeaTask) runOnStop() {
	fmt.Printf("fit tea task; run stop\n")
}

func (ftt FitTeaTask) runPeriodically() {
	fmt.Printf("fit tea task; run periodically; time: %s\n", time.Now().Format("15:04:05.000"))
}
