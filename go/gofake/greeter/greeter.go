package greeter

import "gofake/greeter/animal"

//go:generate gofake -type=Dog
type Dog struct {
	name   string
	animal animal.Animal
	ptr    *animal.Animal
}

// go:generate gofake -type=Greeter
type Greeter interface {
	Greet(left, right string) string
}

func (d *Dog) Greet(left, right string) string {
	return left + d.name + right
}

// go:generate gofake -type=Cat
type Cat struct {
	name, address string
}
