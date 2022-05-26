package pkg

import "fmt"

func PanicAbnormal() {
	//1。
	m := make(map[int]interface{})
	m[1] = m
	fmt.Println(m)

}
