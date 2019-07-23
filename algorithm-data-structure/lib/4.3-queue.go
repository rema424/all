package lib

func ExecQueue() {
	type Process struct {
		name string
		time int
	}

	var head, tail, n int
	size := 100005
	Q := make([]Process, size)

	enqueue := func(p Process) {
		Q[tail] = p
		tail = (tail + 1) % size
	}

	dequeue := func() Process {

	}
}
