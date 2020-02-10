package hello

import "greet"

type Hello struct {
	Msg   string
	Body  greet.Msg
	PBody *greet.Msg
}
