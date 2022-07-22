package main

import (
	"fmt"
	"math"
)

func main() {
	x := 3.1415
	//分别取整
	a := int(math.Ceil(x))
	b := int(math.Floor(x))
	fmt.Println("a", "b", a, b)
	fmt.Println(math.Pow10(3))
	digits := []int{1, 2, 3}
	fmt.Print(plusOne(digits))
}
func plusOne(digits []int) []int {
	sum := 0
	for k, v := range digits {
		sum += v * int(math.Pow10(k))
	}
	sum += 1

	if digits[len(digits)-1] == 9 {
		digits[len(digits)-1] = 0
		digits[len(digits)-2] = digits[len(digits)-2] + 1
	} else {
		digits[len(digits)-1] = digits[len(digits)-1] + 1
	}
	return digits
}
