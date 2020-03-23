package animal

import (
	f "fmt"
)

func A() {
	a := "a"
	あ := Dog{} // 108
	f.Println(あ, a)
}

type Dog struct {
	Parent

	P Parent

	Name string
}
