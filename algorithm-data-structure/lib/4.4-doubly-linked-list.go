package lib

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ExecDoublyLinkedList() {
	s := time.Now()

	// num := 7
	lines := []string{
		"insert 5",
		"insert 2",
		"insert 3",
		"insert 1",
		"delete 3",
		"insert 6",
		"delete 5",
	}

	doublyLinkedList(lines)
	fmt.Println(time.Since(s))
}

func doublyLinkedList(lines []string) {
	sentinel := &Node_4_4{isSentinel: true}

	for _, line := range lines {
		command := strings.Split(line, " ")
		var key int
		if len(command) == 2 {
			key, _ = strconv.Atoi(command[1])
		}
		switch command[0] {
		case "insert":
			insert_4_4(sentinel, key)
		case "delete":
			deleteKey_4_4(sentinel, key)
		case "deleteFirst":
			deleteFirst_4_4(sentinel)
		case "deleteLast":
			deleteLast_4_4(sentinel)
		}
	}
	// formatAns_4_4(sentinel)
	// fmt.Println(formatAns_4_4(sentinel))
	print_4_4_1(sentinel)
	// print_4_4_2(sentinel)
}

type Node_4_4 struct {
	key        int
	next       *Node_4_4
	prev       *Node_4_4
	isSentinel bool
}

func listSearch_4_4(sentinel *Node_4_4, key int) (*Node_4_4, bool) {
	if !sentinel.isSentinel {
		return nil, false
	}

	cur := sentinel.next // 番兵の次から探す
	for cur != nil && !cur.isSentinel && cur.key != key {
		cur = cur.next
	}
	return cur, true
}

func deleteNode_4_4(t *Node_4_4) {
	if t.prev == nil || t.next == nil {
		return
	}

	t.prev.next = t.next
	t.next.prev = t.prev
}

func deleteFirst_4_4(sentinel *Node_4_4) {
	deleteNode_4_4(sentinel.next)
}

func deleteLast_4_4(sentinel *Node_4_4) {
	deleteNode_4_4(sentinel.prev)
}

func deleteKey_4_4(sentinel *Node_4_4, key int) {
	if node, ok := listSearch_4_4(sentinel, key); ok {
		deleteNode_4_4(node)
	}
}

func insert_4_4(sentinel *Node_4_4, key int) {
	if sentinel.next != nil {
		x := &Node_4_4{
			key:  key,
			prev: sentinel,
			next: sentinel.next,
		}
		sentinel.next.prev = x
		sentinel.next = x
	} else {
		x := &Node_4_4{
			key:  key,
			prev: sentinel,
			next: sentinel,
		}
		sentinel.next = x
		sentinel.prev = x

	}
}

func formatAns_4_4(sentinel *Node_4_4) string {
	if !sentinel.isSentinel {
		return ""
	}
	cur := sentinel.next
	results := make([]string, 0, 2000000)
	for cur != nil && !cur.isSentinel {
		results = append(results, strconv.Itoa(cur.key))
		cur = cur.next
	}
	return strings.Join(results, " ")
}

func print_4_4_1(sentinel *Node_4_4) {
	if !sentinel.isSentinel {
		return
	}
	cur := sentinel.next
	for cur != nil && !cur.isSentinel {
		fmt.Print(cur.key, " ")
		cur = cur.next
	}
}

func print_4_4_2(sentinel *Node_4_4) {
	if !sentinel.isSentinel {
		return
	}
	var buffer bytes.Buffer
	cur := sentinel.next
	for cur != nil && !cur.isSentinel {
		buffer.Write([]byte(fmt.Sprintf("%d ", cur.key)))
		cur = cur.next
	}
	buffer.WriteTo(os.Stdout)
}
