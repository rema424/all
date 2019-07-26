package lib

func ExecDoublyLinkedList() {
	type command struct {
		name string
		key  int
	}
	args := struct {
		num      int
		commands []command
	}{
		7,
		[]command{
			command{"insert", 5},
			command{"insert", 2},
			command{"insert", 3},
			command{"insert", 1},
			command{"delete", 3},
			command{"insert", 6},
			command{"delete", 5},
		},
	}

	type Node struct {
		key  int
		next *Node
		prev *Node
	}

	sentinel := &Node{}

	listSearch := func(key int) *Node {
		cur := sentinel.next // 番兵の次から探す
		for cur != nil && cur.key != key {
			cur = cur.next
		}
		return cur
	}

	for i := 0; i < args.num; i++ {
		switch args.commands[i].name {
		case "insert":
		case "delete":
		}
	}
}
