package main

import "fmt"

func main() {
	// pre := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	pre := []int{3, 2, 2, 3}
	res := removeElement(pre, 3)
	fmt.Print(res)
}

func removeElement(nums []int, val int) int {
	l := len(nums)
	if l == 0 {
		return 0
	}
	slow := 0
	for fast := 0; fast <= l-1; fast++ {
		if nums[fast] != val {
			nums[slow] = nums[fast]
			slow++
		}
	}
	return slow
}
