package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Your Nickname? ")
	scanner.Scan()
	nickname := scanner.Text()

	for {
		fmt.Print("input: ")
		scanner.Scan()
		fmt.Printf("%s: %s\n", nickname, scanner.Text())
	}
}
