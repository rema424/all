package main

import "fmt"

func main() {
	type greet struct {
		message string
	}

	a := greet{"hello"}
	b := greet{"hello"}
	c := greet{"bye"}
	d := &greet{"hello"}
	e := &greet{"hello"}

	av := fmt.Sprintf("%+v", a)
	bv := fmt.Sprintf("%#v", b)
	cv := fmt.Sprintf("%+v", c)
	dv := fmt.Sprintf("%+v", d)
	ev := fmt.Sprintf("%#v", e)

	fmt.Println(av, bv, cv, dv, ev)
	fmt.Println(av == bv)
	fmt.Println(av == cv)
	fmt.Println(av == dv)
	fmt.Println(av == ev)
	fmt.Println(dv == ev)
}
