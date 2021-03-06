# gotimer

### examples:

```go
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
```

```go
package main

import (
	"fmt"
	"time"

	"github.com/gocpp/gotimer/timeheap"
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
	th = timeheap.New(selectHandler) // time heap
	th.Start()                       // start

	// add timer for loop
	for i := uint32(1); i <= 5; i++ {
		th.AddTimer(i, 1*time.Second, true)
	}
	th.AddTimer(6, 6*time.Second, false)

	go stop(th) // stop timer
	select {}
}
```