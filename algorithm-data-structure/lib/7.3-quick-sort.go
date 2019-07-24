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
	// ref := cards[right]
	return 0
}
