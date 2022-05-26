package main

import "fmt"

// 工厂接口
type FactoryInterface interface {
	CreateProduct(t string) ProductInterface
}

// 产品接口
type ProductInterface interface {
	Intro()
}

// 创建工厂结构体并实现工厂接口
type Factory1 struct {
}

func (f Factory1) CreateProduct(t string) ProductInterface {
	switch t {
	case "product11":
		return Product11{}
	default:
		return nil
	}

}

// 创建产品11并实现产品接口
type Product11 struct {
}

func (p Product11) Intro() {
	fmt.Println("this is product 1")
}

func main2() {
	// 创建工厂
	f := new(Factory1)

	p := f.CreateProduct("product11")
	p.Intro() // output:  this is product 11.

}
