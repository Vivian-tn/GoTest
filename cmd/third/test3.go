package main

import (
	"fmt"
	"time"
)

//并发
func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println(s)
	}
}

type Huang struct {
	A int64  `json:"a"`
	B string `json:"b"`
}

func main() {
	//ning := make(map[int64]Huang)
	i := 8
	fmt.Printf("i is %+v", i)
}
