package greet

// go:generate gofake -type=Greeter
type Greeter interface {
	Greet(left, right string) string
}

type Dog struct {
	name string
}

func (d *Dog) Greet(left, right string) string {
	return left + d.name + right
}
