package main

import (
	"fmt"
	"strings"
)

func main() {
	haystack := "hello"
	needle := "ll"
	fmt.Println(strStr(haystack, needle))
	tracer := "死神来了,死神bye bye"
	comma := strings.Index(tracer, ",")
	fmt.Println(comma)
	fmt.Println(tracer[comma:])
	pos := strings.Index(tracer[comma:], "死神")
	fmt.Println(pos)
}

func strStr(haystack string, needle string) int {
	index := -1
	index = strings.Index(haystack, needle)

	return index

}
