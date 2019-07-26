package lib

import (
	"fmt"
	"math"
	"strings"
)

func ExecDictionary() {
	num := 11
	commands := [11]string{
		"insert AAA",
		"insert AAC",
		"find AAA",
		"find CCC",
		"insert CCC",
		"find CCC",
		"find AA",
		"find GG",
		"insert GG",
		"find AA",
		"find GG",
	}
	dictionaryWithMap(num, commands[:])
	dictionaryWithArray(num, commands[:])
	fmt.Println(searchPrimeNumber(100))
	for i := 1; i < 100; i++ {
		if primeNumber(i) {
			fmt.Println(i, searchPrimeNumber(i))
		}
	}
}

func dictionaryWithMap(num int, commands []string) {
	dict := make(map[string]bool, len(commands))
	for _, line := range commands {
		command := strings.Split(line, " ")
		switch command[0] {
		case "insert":
			dict[command[1]] = true
		case "find":
			if _, ok := dict[command[1]]; ok {
				fmt.Println("yes")
			} else {
				fmt.Println("no")
			}
		}
	}
}

const primeN = 1046527

var hashTable = make([]int, 1000000)

func dictionaryWithArray(num int, commands []string) {
	for _, line := range commands {
		command := strings.Split(line, " ")
		switch command[0] {
		case "insert":
			insert(command[1])
		case "find":
			if pos := find(command[1]); pos == -1 {
				fmt.Println("no")
			} else {
				fmt.Println("yes")
			}
		}
	}

}

func convertCharToInt(char string) int {
	switch char {
	case "A":
		return 1
	case "C":
		return 2
	case "G":
		return 3
	case "T":
		return 4
	default:
		return 0
	}
}

func getKeyByStr(str string) int {
	sum := 0
	p := 1
	// math.
	for _, r := range []rune(str) {
		sum += p * convertCharToInt(string(r))
		p *= 5
	}
	return sum
}

func h1(key int) int {
	return key % primeN
}

func h2(key int) int {
	return (key % (primeN - 1)) + 1
}

func genHash(key int, index int) int {
	return (h1(key) + index*h2(key)) % primeN
}

func find(str string) int {
	key := getKeyByStr(str)
	i := 0
	for {
		hash := genHash(key, i)
		if hashTable[hash] == key {
			return hash
		} else if hashTable[hash] == 0 || i >= primeN {
			return -1
		}
		i++
	}
}

func insert(str string) int {
	key := getKeyByStr(str)
	i := 0
	for {
		hash := genHash(key, i)
		if hashTable[hash] == 0 {
			hashTable[hash] = key
			return hash
		}
		i++
	}
}

func searchPrimeNumber(min int) int {
	i := (min + 1) / 6
	if (min+1)%6 == 0 {
		i++
	}
	i++
	for {
		x := 6*i - 1
		if primeNumber(x) {
			return x
		}
		y := 6*i + 1
		if primeNumber(y) {
			return y
		}
		i++
	}
}

// func primeNumber(x int) bool {
// 	b := true
// 	for m := 2; m < x; m++ {
// 		if x%m == 0 {
// 			b = false
// 			break
// 		}
// 	}
// 	return b
// }

func primeNumber(num int) bool {
	if num <= 1 {
		return false
	}

	if num == 2 {
		return false
	}
	if num%2 == 0 {
		return false
	}

	b := true
	root := int(math.Sqrt(float64(num)))
	div := 3
	for div <= root {
		if num%div == 0 {
			b = false
			break
		}
		div += 2
	}
	return b
}
