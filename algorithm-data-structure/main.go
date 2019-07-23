package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/rema424/all/algorithm-data-structure/lib"
)

func main() {
	// measure(lib.ExecMaximamProfit)
	// lib.ExecInsertionSort()
	// lib.ExecBubbleSort()
	// lib.ExecSelectionSort()
	// lib.ExecStableSort()
	// lib.ExecShellSort()
	lib.ExecStack()
	lib.ExecQueue()
}

type tm struct {
	name    string
	measure time.Duration
}

var (
	times = make([]tm, 0, 100)
	mem_s runtime.MemStats
	mem_e runtime.MemStats
)

func measure(fn func()) {
	runtime.ReadMemStats(&mem_s)
	start := time.Now()

	fn()
	times = append(times, tm{"fn", time.Since(start)})

	runtime.ReadMemStats(&mem_e)

	fmt.Printf("Memory : %f MB\n", float64(mem_e.Alloc-mem_s.Alloc)/float64(1024*1024))

	tmp := 0.0
	for _, v := range times {
		fmt.Printf("%2.8f / %2.8f sec <= %s\n", v.measure.Seconds()-tmp, v.measure.Seconds(), v.name)
		tmp = v.measure.Seconds()
	}
}
