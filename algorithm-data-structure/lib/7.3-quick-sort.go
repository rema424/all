package lib

import "fmt"

type card struct {
	suit  string
	value int
}

func ExecQuickSort() {
	args := []card{
		card{"D", 3},
		card{"H", 2},
		card{"D", 1},
		card{"S", 3},
		card{"D", 2},
		card{"C", 1},
	}
	fmt.Println("sort前:", args)
	fmt.Println("sort後:", quickSort(args, len(args), 0, len(args)-1))
}

func quickSort(cards []card, size, left, right int) []card {
	if left < right {
		mid := partition_2(cards, size, left, right)
		quickSort(cards, size, left, mid-1)
		quickSort(cards, size, mid+1, right)
	}
	return cards
}

func partition_2(cards []card, size, left, right int) int {
	ref := cards[right]
	smallTail := left - 1
	for curr := left; curr < right; curr++ {
		if cards[curr].value <= ref.value {
			smallTail++
			cards[smallTail], cards[curr] = cards[curr], cards[smallTail]
		}
	}
	cards[smallTail+1], cards[right] = cards[right], cards[smallTail+1]
	return smallTail + 1
}

func mergeSort(cards []card, size, left, rightX int) []card {
	if rightX-left > 1 { // つまり2個以上
		mid := (left + rightX) / 2
		mergeSort(cards, size, left, mid)
		mergeSort(cards, size, mid, rightX)
		merge(cards, size, left, mid, rightX)
	}
	return cards
}

func merge(cards []card, size, left, mid, rightX int) []card {
	leftN := mid - left
	rightN := rightX - mid
	L := make([]card, leftN+1) // 番兵を置くため1個多く
	R := make([]card, rightN+1)
	for i := 0; i < leftN; i++ {
		L[i] = cards[left+i]
	}
	for i := 0; i < rightN; i++ {
		R[i] = cards[mid+i]
	}
	return cards
}
