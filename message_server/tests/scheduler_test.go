package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-co-op/gocron/v2"
	"wraith.me/message_server/task"
)

func TestScheduler(t *testing.T) {
	//Create tasks to run
	ftt := task.FitTeaTask{}

	//Setup the scheduler
	s := task.Scheduler{}
	s.Register(ftt)

	//Run the scheduler
	s.Start()
	fmt.Println("Started scheduler")

	//Block for 20 seconds
	time.Sleep(time.Second * 20)

	//Shutdown the scheduler
	s.Shutdown()
	fmt.Println("Stopped scheduler")
}

func TestBareScheduler(t *testing.T) {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		t.Fatal(err)
	}

	// add a job to the scheduler
	j, err := s.NewJob(
		gocron.DurationJob(
			(10 * time.Second),
		),
		gocron.NewTask(
			func(a string, b int) {
				fmt.Printf("%s %d\n", a, b)
			},
			"hello",
			1,
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	// each job has a unique id
	fmt.Println(j.ID())

	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	time.Sleep(time.Second * 21)

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		t.Fatal(err)
	}
}
