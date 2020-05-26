package main

import (
	"fmt"
	"time"

	"github.com/gocpp/gotimer/timeheap"
)

// timer callback
func selectHandler(id uint32) {
	fmt.Printf("timerId:%d now:%v\n", id, time.Now())
}

func main() {
	th := timeheap.New(selectHandler) // time heap
	th.Start()                        // start

	// add timer
	th.AddTimer(1, 1*time.Second, false)

	// wait
	time.Sleep(2 * time.Second)
}
