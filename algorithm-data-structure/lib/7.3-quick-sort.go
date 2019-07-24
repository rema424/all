package lib

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
