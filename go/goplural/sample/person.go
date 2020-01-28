package sample

import "fmt"

type Greeter interface {
	Greet()
}

type (
	Greeter2 interface {
		Greet2_1(a string) string
		Greet2_2(b string, c int) (string, error)
	}

	Greeter3 interface {
		Greet3_1()
	}
)

type Person struct {
	ID, Code int64
	Name     string `db:"aaa"`
}

func (p *Person) Greeter() {
	fmt.Println("Hi! I'm", p.Name, "!")
}
