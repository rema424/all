package main

import (
	"fmt"
	"sync"
)

func execSync() {
	var data int
	go func() {
		data++
	}()
	if data == 0 {
		fmt.Println("the value is 0.")
	} else {
		fmt.Printf("the value is %v.\n", data)
	}
}

func execSyncFix() {
	var memoryAccess sync.Mutex
	var value int
	go func() {
		memoryAccess.Lock()
		value++
		memoryAccess.Unlock()
	}()

	memoryAccess.Lock()
	if value == 0 {
		fmt.Printf("the value is %v.\n", value)
	} else {
		fmt.Printf("the value is %v.\n", value)
	}
	memoryAccess.Unlock()
}
