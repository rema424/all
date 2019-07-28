package lib

import "testing"

func BenchmarkDoublyLinkedList(b *testing.B) {
	lines := []string{
		"insert 5",
		"insert 2",
		"insert 3",
		"insert 1",
		"delete 3",
		"insert 6",
		"delete 5",
	}
	for i := 0; i < b.N; i++ {
		doublyLinkedList(lines)
	}
}
