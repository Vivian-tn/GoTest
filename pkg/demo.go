package pkg

import (
	"fmt"
)

type a struct {
	name *a
	i    int
}

func Demo() {
	var b a
	b.i = 1
	b.name = &b
	fmt.Println(b)
}
