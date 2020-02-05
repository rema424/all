// package comment

package greet

//go:generate gofake -type=Greeter
type Greeter interface {
	Greet(left, right string) string
}

//go:generate gofake -type=Dog

type Dog struct {
	name string
}

func (d *Dog) Greet(left, right string) string {
	return left + d.name + right
}

//go:generate gofake -type=Cat
type Cat struct {
	name, address string
}
