package main

import (
	"fmt"
)

func main() {
	x := 1
	fmt.Println(isPalindrome(x))
}
func isPalindrome(x int) bool {
	var slice []int
	var res bool
	if x < 0 {
		return false
	} else if x/10 == 0 {
		return true
	} else {
		for true {
			slice = append(slice, x%10)
			x = x / 10
			if x == 0 {
				break
			}
		}
		len := len(slice)
		for i := 0; i < len/2; i++ {
			if slice[i] == slice[len-1-i] {
				res = true
				continue
			} else {
				res = false
				break
			}
		}
	}
	return res
}
