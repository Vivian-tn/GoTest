package main

import "fmt"

func main() {
	nums := []int{1, 3, 5, 6}
	target := 7
	res := searchInsert(nums, target)
	fmt.Print(res)
}
func searchInsert(nums []int, target int) int {
	var res int
	if target < nums[0] {
		return 0
	} else if target > nums[len(nums)-1] {
		return len(nums)
	}
	for k, num := range nums {
		if num == target {
			res = k
			break
		} else if k > 0 && nums[k-1] < target && target < nums[k] {
			res = k
		}
	}
	return res

}
