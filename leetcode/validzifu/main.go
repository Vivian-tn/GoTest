package main

import (
	"fmt"
)

func main() {
	test := map[int]int{
		1: 1,
		2: 2,
		3: 3,
	}
	fmt.Println(test[4])
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}
	fmt.Println(pairs)
	fmt.Println(isValid("()[]{}"))
}
func isValid(s string) bool {
	n := len(s)
	if n%2 == 1 {
		return false
	}
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}
	stack := []byte{}
	for i := 0; i < n; i++ {
		if pairs[s[i]] > 0 { //栈空或者不对应
			if len(stack) == 0 || stack[len(stack)-1] != pairs[s[i]] {
				return false
			}
			stack = stack[:len(stack)-1]
		} else { //s[i]是'('/'['/'{' 入栈
			stack = append(stack, s[i])
		}
	}
	return len(stack) == 0
}
