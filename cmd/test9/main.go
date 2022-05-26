package main

import (
	"fmt"
	"reflect"
)

type Data struct {
}

func (d *Data) User() (string, bool) {
	return "方法调用成功", false
}

func (d *Data) Person() bool {
	return false
}

func main() {
	data := Data{}
	value := reflect.ValueOf(&data)
	num := value.NumMethod()
	fmt.Println(num)
	res := value.MethodByName("Person")
	fmt.Println(res.Call([]reflect.Value{}))
}
