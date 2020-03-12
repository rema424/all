package animal

import (
	f "fmt"
)

type Dog struct {
	Parent

	P Parent

	Name string
}

func A() {
	a := Dog{} // 108
	f.Println(a)
}
