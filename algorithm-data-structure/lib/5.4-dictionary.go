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

func dictionaryWithArray(num int, commands []string) {
	isPrimeNumber := func(num int) bool {
		if num <= 1 {
			return false
		} else if num == 2 {
			return true
		} else if num%2 == 0 {
			return false
		}

		for div := 3; div <= int(math.Sqrt(float64(num))); div += 2 {
			if num%div == 0 {
				return false
			}
		}
		return true
	}

	calcNextPrimeNumber := func(min int) int {
		var x int
		for i := (min + 1) / 6; ; i++ { // 素数 = 6n ± 1 > min より
			x = 6*i - 1
			if x > min && isPrimeNumber(x) {
				return x
			}

			x = 6*i + 1
			if x > min && isPrimeNumber(x) {
				return x
			}
		}
	}

	// 割り算の結果確実に余りが0にならないように大きな素数を計算に使う
	primeN := calcNextPrimeNumber(num)

	h1 := func(key int) int {
		return key % primeN
	}

	h2 := func(key int) int {
		return (key % (primeN - 1)) + 1
	}

	calcHash := func(key int, index int) int {
		return (h1(key) + index*h2(key)) % primeN
	}

	hashTable := make([]int, primeN)

	insert := func(hashTable []int, key int) int {
		for i := 0; ; i++ {
			hash := calcHash(key, i)
			if hashTable[hash] == 0 {
				hashTable[hash] = key
				return hash
			}
		}
	}

	search := func(hashTable []int, key int) int {
		for i := 0; ; i++ {
			hash := calcHash(key, i)
			if hashTable[hash] == key {
				return hash
			} else if hashTable[hash] == 0 || i >= len(hashTable) {
				return -1
			}
		}
	}

	convertCharToInt := func(char string) int {
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

	calcKey := func(str string) int {
		sum := 0
		p := 1
		for _, r := range []rune(str) {
			sum += p * convertCharToInt(string(r))
			p *= 5 // 文字の種類が4種類で、4より大きな素数
		}
		return sum
	}

	for _, line := range commands {
		command := strings.Split(line, " ")
		key := calcKey(command[1])
		switch command[0] {
		case "insert":
			insert(hashTable, key)
		case "find":
			if pos := search(hashTable, key); pos == -1 {
				fmt.Println("no")
			} else {
				fmt.Println("yes")
			}
		}
	}
}
