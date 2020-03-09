package greeter

import "gofake/greeter/animal"

//go:generate gofake -type=Dog
type Dog struct {
	valString string
	ptrString *string
	valAnimal animal.Animal
	ptrAnimal *animal.Animal
}

// go:generate gofake -type=Greeter
type Greeter interface {
	Greet(left, right string) string
}

// func (d *Dog) Greet(left, right string) string {
// 	return left + d.name + right
// }

// go:generate gofake -type=Cat
// type Cat struct {
// 	name, address string
// }

//go:generate gofake -type=Person
type Person struct {
	ID      int64
	Profile Profile
	Address *Address
	Hobbies []Hobby
	Pets    []*Pet
	Job     *struct {
		Name string
	}
	Food struct {
		Name string
	}
	valAnimal animal.Animal
	ptrAnimal *animal.Animal
}

type Profile struct {
	Name string
}

type Address struct {
	Zip string
}

type Hobby struct {
	Name string
}

type Pet struct {
	Name string
}
