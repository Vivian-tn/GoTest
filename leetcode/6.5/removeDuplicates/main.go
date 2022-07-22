package main

import "fmt"

func main() {

	// pre := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	pre := []int{1, 1, 2}
	res := removeDuplicates(pre)
	fmt.Print(res)
}
func removeDuplicates(nums []int) int {
	l := len(nums)
	if l == 0 {
		return 0
	}
	slow := 1
	for fast := 1; fast <= l-1; fast++ {
		if nums[fast] != nums[fast-1] {
			nums[slow] = nums[fast]
			slow++
		}
	}
	return slow
}
