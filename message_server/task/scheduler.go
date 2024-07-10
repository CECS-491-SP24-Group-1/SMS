package task

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

// Represents a simple wrapper around a `gocron.Scheduler` object.
type Scheduler struct {
	//The task objects that are to be ran by the scheduler.
	tasks []Task

	//The IDs of the hooked tasks.
	taskIDs []uuid.UUID

	//Whether the scheduler is active.
	started bool

	//The underlying scheduler object.
	handler gocron.Scheduler
}

// Registers n number of tasks with the scheduler.
func (s *Scheduler) Register(tasks ...Task) error {
	//Create a handler if it doesn't already exist
	if s.handler == nil {
		var err error
		s.handler, err = gocron.NewScheduler()
		if err != nil {
			return err
		}
	}

	//Add the periodic tasks to the handler
	for i, ptask := range tasks {
		//Add the task to the list of tasks
		s.tasks = append(s.tasks, ptask)
		s.taskIDs = append(s.taskIDs, uuid.Nil)

		//Create a new job
		job, err := s.handler.NewJob(
			gocron.DurationJob(ptask.periodicDuration()),
			gocron.NewTask(ptask.runPeriodically),
		)
		if err != nil {
			return err
		}

		//Replace the current task ID with the ID of the added job
		s.taskIDs[i] = job.ID()
	}

	//No errors so return nil
	return nil
}

// Starts the scheduler, running all startup functions in the process.
func (s *Scheduler) Start() error {
	//Run only if the scheduler is not already running
	if s.started {
		return fmt.Errorf("cannot start an already running scheduler")
	}
	s.started = true

	//Run the "on start" tasks
	for _, task := range s.tasks {
		go task.runOnStart()
	}

	//Start the handler in async mode; a goroutine is implicitly started here
	s.handler.Start()
	return nil
}

// Stops the scheduler, running all shutdown functions in the process.
func (s *Scheduler) Shutdown() error {
	// Run only if the scheduler is already running
	if !s.started {
		return fmt.Errorf("cannot shutdown an already stopped scheduler")
	}
	s.started = false

	//Run the "on stop" tasks
	for i, task := range s.tasks {
		go task.runOnStop()
		s.taskIDs[i] = uuid.Nil
	}

	//Stop the handler
	return s.handler.Shutdown()
}
