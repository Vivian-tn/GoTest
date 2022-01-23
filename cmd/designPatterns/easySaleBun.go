package main

import "fmt"

type BunShop struct {
}

type maibaozi interface {
	create()
}
type Pigxiande struct {
}
type Sanxiande struct {
}

func (baozi Pigxiande) create() {
	fmt.Println("Pigxiandebaozi")
}
func (baozi Sanxiande) create() {
	fmt.Println("Sanxiandebaozi")
}

func (dianpu BunShop) generate(baozi string) maibaozi {
	switch baozi {
	case "Pigxiande":
		return Pigxiande{}
	case "Sanxiande":
		return Sanxiande{}
	default:
		return nil
	}
}
func mai4() {
	GoubuliShop := new(BunShop)
	baozi1 := GoubuliShop.generate("Pigxiande")
	baozi1.create()

	baozi2 := GoubuliShop.generate("Sanxiande")
	baozi2.create()
}
