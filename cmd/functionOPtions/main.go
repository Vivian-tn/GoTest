package main

import (
	"encoding/json"
	"fmt"
)

type Skateboard struct {
	Id   int64  `json:"id"`
	Name string `json:"name"` //
	Guys string `json:"guys"` //
}
type Option func(option *Skateboard)

func WithName(name string) Option {
	return func(s *Skateboard) {
		s.Name = name
	}
}

func WithGuys(guys string) Option {
	return func(s *Skateboard) {
		s.Guys = guys
	}
}

func NewSkateboard(id int64, options ...Option) *Skateboard {
	skateboard := &Skateboard{Id: id}
	// kexuna
	for _, i := range options {
		i(skateboard)
	}
	return skateboard
}

func main() {
	res := NewSkateboard(9, WithName("Eric"), WithGuys("Vivian"))
	a, _ := json.Marshal(res)
	fmt.Println(string(a))
}
