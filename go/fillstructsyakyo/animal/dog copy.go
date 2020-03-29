package animal

import (
	f "fmt"
)

type Dog2 struct {
	Parent

	P Parent

	Name string
}

func A2() {
	a := Dog{} // 108
	f.Println(a)
}
