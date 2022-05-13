package main

import (
	"fmt"
)

func main() {
	nums := []int{3, 2, 4}
	fmt.Println(twoSum(nums, 6))
}
func twoSum(nums []int, target int) []int {
	exitNum := make(map[int]int)
	for k, v := range nums {
		exitNum[v] = k
	}
	fmt.Println(exitNum)
	for k, v := range nums {
		if _, ok := exitNum[target-v]; ok && k != exitNum[target-v] {
			return []int{k, exitNum[target-v]}
		}
	}
	return nil
}
