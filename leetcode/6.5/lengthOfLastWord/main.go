package main

import "fmt"

func main() {
	s := "Hello World"
	fmt.Println(lengthOfLastWord(s))
}
func lengthOfLastWord(s string) int {
	var ans int
	index := len(s) - 1
	for s[index] == ' ' {
		index--
	}
	for index >= 0 && s[index] != ' ' {
		ans++
		index--
	}
	return ans
}
