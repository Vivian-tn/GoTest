package main

import (
	"fmt"
)

func main() {
	// arr := []int64{1, 2, 3, 4}
	// arr1 := arr[:0]
	// fmt.Println(arr1)
	//var wg sync.WaitGroup
	//for _, v := range arr {
	//	if v == 3 {
	//		break
	//	} else {
	//		fmt.Println(v)
	//	}
	//}
	//fmt.Println("aaaa")
	//token, err := strconv.ParseInt("", 10, 64)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(token)
	//var str *string
	// nums := []int{3, 2, 4}
	// fmt.Println(twoSum(nums, 6))
	var slice []int
	x := 123
	for true {
		slice = append(slice, x%10)
		x = x / 10
		if x == 0 {
			break
		}
	}
	fmt.Println(slice)

}
func twoSum(nums []int, target int) []int {
	exitNum := make(map[int]int)
	for k, v := range nums {
		exitNum[v] = k
	}
	fmt.Println(exitNum)
	for k, v := range nums {
		if _, ok := exitNum[target-v]; ok && k != exitNum[target-v] {
			fmt.Println(k, exitNum[target-v])
			return []int{k, exitNum[target-v]}
		}
	}
	return nil
}
