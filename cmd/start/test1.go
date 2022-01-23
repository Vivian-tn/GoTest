package main

import (
	"fmt"
	"os"
	"strings"
	"testNew/pkg"
	"unsafe"
)

func main() {
	ss := String("Article")
	fmt.Println(*ss)
	os.Setenv("FOO", "1")
	fmt.Println("FOO:", os.Getenv("FOO"))
	fmt.Println("BAR:", os.Getenv("BAR"))

	fmt.Println()
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		fmt.Println(pair[0])
	}
	var i1 int = 1
	var i2 int8 = 2
	var i3 int16 = 3
	var i4 int32 = 4
	var i5 int64 = 5
	fmt.Println(unsafe.Sizeof(i1))
	fmt.Println(unsafe.Sizeof(i2))
	fmt.Println(unsafe.Sizeof(i3))
	fmt.Println(unsafe.Sizeof(i4))
	fmt.Println(unsafe.Sizeof(i5))
	xx := []int{10, 20, 30}
	yy := make([]int, 2)
	copy(yy, xx)
	fmt.Println(xx, yy)
	var score float32 = -1.1
	var temp float64 = 12.2e10
	fmt.Println("Score = ", score, " Temp = ", temp)
	var r byte = 0
	fmt.Println(r)
	var aa int
	aa = 10
	getType(aa)
	fmt.Printf("%q", 0x4E2D)
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

	//pkg.PanicAbnormal()
	pkg.Demo()
	pkg.ForXunHuan()
	pkg.PointerTest()
	pkg.RangeKeyWords()
	pkg.ZhengZe()
	pkg.SynMap()
	pkg.StrconvParseInt()
	fmt.Println("*&", pkg.Int(1))
}

//func main() {
//	var a int
//	var b int8
//	var c int16
//	var d int32
//	var e int64
//	slice := make([]int, 0)
//	slice = append(slice, 1)
//	fmt.Printf("int:%dnint8:%dnint16:%dnint32:%dnint64:%dn", unsafe.Sizeof(a), unsafe.Sizeof(b), unsafe.Sizeof(c), unsafe.Sizeof(d), unsafe.Sizeof(e))
//	fmt.Printf("slice:%d", unsafe.Sizeof(slice))
//}
func getType(a interface{}) {
	switch a.(type) {
	case int:
		fmt.Println("the type of a is int")
	case string:
		fmt.Println("the type of a is string")
	case float64:
		fmt.Println("the type of a is float")
	default:
		fmt.Println("unknown type")
	}
}
func String(s string) *string {
	return &s
}
