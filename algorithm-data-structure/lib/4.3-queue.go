package lib

import "fmt"

func ExecQueue() {
	type Process struct {
		name string
		time int
	}

	var head, tail int
	size := 100005
	Q := make([]Process, size)

	enqueue := func(p Process) {
		Q[tail] = p
		tail = (tail + 1) % size
	}

	dequeue := func() Process {
		x := Q[head]
		head = (head + 1) % size
		return x
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	args := struct {
		num       int
		quantum   int
		processes []Process
	}{
		5,
		100,
		[]Process{
			Process{"p1", 150},
			Process{"p2", 80},
			Process{"p3", 200},
			Process{"p4", 350},
			Process{"p5", 20},
		},
	}

	for i := 0; i < args.num; i++ {
		Q[i] = args.processes[i]
	}

	head, tail = 0, args.num

	var elaps int
	for head != tail {
		p := dequeue()
		t := min(args.quantum, p.time)
		p.time -= t
		elaps += t
		if p.time > 0 {
			enqueue(p)
		} else {
			fmt.Println(p.name, elaps)
		}
	}
}
