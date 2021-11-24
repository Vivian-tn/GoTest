package main

import "fmt"

func main() {
	fmt.Println("vivian")
	fmt.Println("田宁")
	var stockcode = 123
	var enddate = "2020-12-31"
	var url = "Code=%d&endDate=%s"
	var target_url = fmt.Sprintf(url, stockcode, enddate)
	fmt.Println(target_url)

	var i, j int = 1, 2
	var b = "2131"
	fmt.Println(i, j, b)

	var a [3]int
	a[0] = 1
	a[1] = 2
	a[2] = 3
	fmt.Println(a)

	x := [5]int{10, 20, 30, 40, 50}
	array := make([]int64, 2)
	fmt.Println(x, array)

	y := [5]int{2: 10, 4: 40}
	fmt.Println(y)

	forXunHuan()
	rangeExample()
}

func forXunHuan() {
	tests := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i, test := range tests {
		if i < 3 {
			fmt.Println("default")
		} else {
			fmt.Println(test)
		}
	}
}
func rangeExample() {
	nums := []int{2, 3, 4}
	sum := 0
	for _, num := range nums {
		sum += num
	}
	fmt.Println("sum:", sum)

	for i, num := range nums {
		if num == 2 {
			fmt.Println("下标索引值是：", i)
		}
	}

	//kvs := map
}
