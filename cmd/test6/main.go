package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	cc := make(chan int, 6)
	titleIds := []int{1, 2, 3, 4, 5, 6}
	list := make([]int, 0)
	for _, id := range titleIds {
		id := id
		wg.Add(1)
		go func() {
			defer wg.Done()
			k := add(id)
			cc <- k
		}()
	}
	wg.Wait()
	close(cc)
	for v := range cc {
		list = append(list, v)
	}
	fmt.Println(list)
}

func add(i int) int {
	return i + 1
}
