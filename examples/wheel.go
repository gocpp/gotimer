package main

import (
	"fmt"
	"time"

	"github.com/gocpp/gotimer/timewheel"
)

type TimerCb interface {
	Start()
	Stop()
	AddTimer(uint32, time.Duration, bool)
	RemoveTimer(uint32)
}

// timer callback
func selectHandler(id uint32) {
	fmt.Printf("timerId:%d now:%v\n", id, time.Now())
}

func stop(th TimerCb) {
	for i := uint32(1); i <= 5; i++ {
		time.Sleep(1 * time.Second)
		th.RemoveTimer(i)
	}
}

func main() {
	var th TimerCb
	th = timewheel.New(100*time.Microsecond, 60*60, selectHandler) // time wheel
	th.Start()                                                     // start timewheel

	// add timer for loop
	for i := uint32(1); i <= 5; i++ {
		th.AddTimer(i, 1*time.Second, true)
	}
	//add timer
	th.AddTimer(6, 6*time.Second, false)

	// stop timer
	go stop(th)

	// wait
	time.Sleep(10 * time.Second)
}
