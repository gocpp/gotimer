package timeheap

import (
	"time"

	"github.com/gocpp/gotimer/heap"
)

// Job 延时任务回调函数
type Job func(uint32)

// Task 延时任务
type Task struct {
	delay      time.Duration // 延迟时间
	expiration time.Time     // 延迟点的时间
	timeId     uint32        // 定时器唯一标识, 用于删除定时器
	index      int           // 索引
	loop       bool          // 是否是周期任务
}

// timeHeap is a heap-based priority queue
type timeHeapType []*Task

func (heap timeHeapType) Len() int {
	return len(heap)
}

func (heap timeHeapType) Less(i, j int) bool {
	return heap[i].expiration.UnixNano() < heap[j].expiration.UnixNano()
}

func (heap timeHeapType) Swap(i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
	heap[i].index = i
	heap[j].index = j
}

func (heap *timeHeapType) Push(x interface{}) {
	n := len(*heap)
	timer := x.(*Task)
	timer.index = n
	*heap = append(*heap, timer)
}

func (heap *timeHeapType) Pop() interface{} {
	old := *heap
	n := len(old)
	timer := old[n-1]
	timer.index = -1
	*heap = old[0 : n-1]
	return timer
}

func (heap *timeHeapType) Top() interface{} {
	old := *heap
	n := len(old)
	return old[n-1]
}

func (heap *timeHeapType) Clear() {
	old := *heap
	*heap = old[:0]
}

// TimeHeap 实现
type TimeHeap struct {
	timer             *time.Timer      // 定时器
	timers            timeHeapType     // 定时任务集合
	indexes           map[uint32]*Task // 定时任务的索引管理，key：定时器唯一标识，value：任务
	job               Job              // 定时器回调函数
	addTaskChannel    chan Task        // 新增任务channel
	removeTaskChannel chan uint32      // 删除任务channel
	stopChannel       chan bool        // 停止定时器channel
}

// New 创建时间轮
func New(job Job) *TimeHeap {
	if job == nil {
		return nil
	}
	tw := &TimeHeap{
		timer:             time.NewTimer(1 * time.Second),
		timers:            make(timeHeapType, 0),
		indexes:           make(map[uint32]*Task),
		job:               job,
		addTaskChannel:    make(chan Task),
		removeTaskChannel: make(chan uint32),
		stopChannel:       make(chan bool),
	}
	tw.timer.Stop() //
	return tw
}

// Start 启动
func (th *TimeHeap) Start() {
	go th.start()
}

// Stop 停止
func (th *TimeHeap) Stop() {
	th.stopChannel <- true
}

// AddTimer 添加定时器 timerId为定时器唯一标识
func (th *TimeHeap) AddTimer(timerId uint32, delay time.Duration, loop bool) {
	if delay < 0 || timerId <= 0 {
		return
	}
	th.addTaskChannel <- Task{delay: delay, expiration: time.Now().Add(delay), timeId: timerId, loop: loop}
}

// RemoveTimer 删除定时器 timerId为添加定时器时传递的定时器唯一标识
func (th *TimeHeap) RemoveTimer(timerId uint32) {
	if timerId <= 0 {
		return
	}
	th.removeTaskChannel <- timerId
}

func (th *TimeHeap) start() {
	for {
		select {
		case <-th.timer.C:
			th.timerHandler()
		case task := <-th.addTaskChannel:
			th.addTask(&task)
		case timerId := <-th.removeTaskChannel:
			th.removeTask(timerId, false)
		case <-th.stopChannel:
			th.timer.Stop()
			return
		}
	}
}

// 执行定时器
func (th *TimeHeap) timerHandler() {
	var (
		top *Task
		ok  bool
	)

	// 执行到期的定时任务
	for {
		top, ok = heap.Pop(&th.timers).(*Task)
		if !ok {
			return
		}

		// 验证
		if top.expiration.UnixNano() <= time.Now().UnixNano() {
			th.job(top.timeId) // 不新建goroutine，job回调用chan发送至工作goroutine
		} else {
			heap.Push(&th.timers, top)
			break
		}

		// 如果是周期任务则重新添加
		if top.loop {
			top.expiration = time.Now().Add(top.delay)
			heap.Push(&th.timers, top)
		}
	}
	// 重置定时器
	th.timer.Reset(top.expiration.Sub(time.Now()))
}

// 新增任务
func (th *TimeHeap) addTask(task *Task) {
	th.removeTask(task.timeId, true)

	heap.Push(&th.timers, task)
	th.indexes[task.timeId] = task

	top := heap.Top(&th.timers).(*Task)
	if task.timeId == top.timeId {
		th.timer.Reset(top.expiration.Sub(time.Now()))
	}
}

// 删除任务 isAdd是否由addTask调用
func (th *TimeHeap) removeTask(timerId uint32, isAdd bool) {
	task, ok := th.indexes[timerId]
	if !ok {
		return
	}

	delete(th.indexes, timerId)
	// 删除定任务
	heap.Remove(&th.timers, task.index)

	if isAdd {
		return
	}

	top, ok := heap.Top(&th.timers).(*Task)
	if !ok {
		th.timer.Stop()
		return
	}
	th.timer.Reset(top.expiration.Sub(time.Now()))
}
