package main

import (
	"fmt"
	"testNew/pkg"
)

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

	//forXunHuan()
	//rangeExample()

	pkg.ForXunHuan()
	pkg.RangeKeyWords()
	pkg.ZhengZe()

	var width, height = 100, 50 //"int" is dropped

	fmt.Println("width is", width, "height is", height)
	//width is 100 height is 50

	var (
		name    = "naveen"
		age     = 29
		height1 int
	)
	fmt.Println("my name is", name, ", age is", age, "and height is", height1)
	//my name is naveen , age is 29 and height is 0

}
