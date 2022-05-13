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
	//i := 8
	//fmt.Printf("i is %+v", i)

	c := make(chan int)

	go func() {

		c <- 1
		c <- 2
		c <- 3
		close(c)
	}()

	for v := range c {
		fmt.Println(v)
	}

	queryTitles := []int{3, 2, 1}
	for _, q := range queryTitles {
		for v := range c {
			if q == v {
				fmt.Println(v)
			}
		}
	}
}
