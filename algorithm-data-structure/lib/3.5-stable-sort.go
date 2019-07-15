package lib

import (
	"fmt"
	"strings"
)

func ExecStableSort() {
	cards := []Card{
		Card{"H", 4},
		Card{"C", 9},
		Card{"S", 4},
		Card{"D", 2},
		Card{"C", 3},
	}
	var bubbleCards CardList = bubble(cards)
	var selectionCards CardList = selection(cards)
	fmt.Println("buble:", bubbleCards.Show(), IsStable(bubbleCards, cards))
	fmt.Println("selection:", selectionCards.Show(), IsStable(selectionCards, cards))
}

type Card struct {
	Suit string
	Val  int
}

type CardList []Card

func (c Card) Show() string {
	return fmt.Sprintf("%s%d", c.Suit, c.Val)
}

func (cl CardList) Show() string {
	l := make([]string, len(cl))
	for i, card := range cl {
		l[i] = card.Show()
	}
	return strings.Join(l, " ")
}

func IsStable(org []Card, sorted []Card) bool {
	if len(org) != len(sorted) {
		return false
	}

	for i := 0; i < len(org); i++ {
		val := org[i].Val
		for j := i + 1; j < len(org); j++ {
			if org[j].Val != val {
				continue
			}
			for a := 0; a < len(org); a++ {
				if sorted[a].Val != val {
					continue
				}
				for b := a + 1; b < len(org); b++ {
					if org[b].Val != val {
						continue
					}
					if org[i].Show() == sorted[b].Show() && org[j].Show() == sorted[a].Show() {
						return false
					}
				}
			}
		}
	}
	return true
}

func IsStableByBubble(org []Card, sorted []Card) bool {
	bubbleSorted := bubble(org)

	if len(bubbleSorted) != len(sorted) {
		return false
	}

	for i := 0; i < len(bubbleSorted); i++ {
		if bubbleSorted[i].Show() != sorted[i].Show() {
			return false
		}
	}
	return true
}

func bubble(args []Card) []Card {
	cards := make([]Card, 0, len(args))
	cards = append(cards, args...)

	finished := false
	for i := 0; !finished; i++ {
		finished = true
		for j := len(cards) - 1; j > i; j-- {
			if cards[j-1].Val > cards[j].Val {
				cards[j-1], cards[j] = cards[j], cards[j-1]
				finished = false
			}
		}
	}
	return cards
}

func selection(args []Card) []Card {
	cards := make([]Card, 0, len(args))
	cards = append(cards, args...)

	for i := 0; i < len(cards); i++ {
		min := i
		for j := i; j < len(cards); j++ {
			if cards[j].Val < cards[min].Val {
				min = j
			}
		}
		if i != min {
			cards[i], cards[min] = cards[min], cards[i]
		}
	}
	return cards
}
