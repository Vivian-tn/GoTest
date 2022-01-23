package main

import "fmt"

//添加工厂接口(包子店的接口)和产品接口(包子接口)
type BunShopInterface interface {
	Generate(t string) Bun
}

type Bun interface {
	create()
}

//创建工厂结构体和产品结构体（具体包子店和具体包子）
type QSPigMeatBuns struct{}
type GDPigMeatBuns struct{}

type QSSamSunStuffingBuns struct{}
type GDSamSunStuffingBuns struct{}

// 实现产品接口
func (p QSPigMeatBuns) create() {
	fmt.Println("qs Pig")
}
func (p GDPigMeatBuns) create() {
	fmt.Println("gd Pig")
}
func (p QSSamSunStuffingBuns) create() {
	fmt.Println("qs SamSun")
}
func (p GDSamSunStuffingBuns) create() {
	fmt.Println("gd SamSun")
}

//创建对应的工厂（齐市包子店和广东包子店）
type QSBunShop struct{}

type GDBunShop struct{}

func (qs QSBunShop) Generate(t string) Bun {
	switch t {
	case "pig":
		return QSPigMeatBuns{}
	case "3s":
		return QSSamSunStuffingBuns{}
	default:
		return nil
	}
}

func (gd GDBunShop) Generate(t string) Bun {
	switch t {
	case "pig":
		return GDPigMeatBuns{}
	case "3s":
		return GDSamSunStuffingBuns{}
	default:
		return nil
	}
}

func main() {
	var b Bun
	var c Bun

	// 卖呀卖呀卖包子...
	QSFactory := new(QSBunShop)
	b = QSFactory.Generate("pig") // 传入猪肉馅的参数，会返回齐市包子铺的猪肉馅包子
	b.create()
	c = QSFactory.Generate("3s") // 传入三鲜的参数，会返回齐市包子铺的猪肉馅包子
	c.create()

	GDFactory := new(GDBunShop)
	b = GDFactory.Generate("pig") // 同样传入猪肉馅的参数，会返回广东包子铺的猪肉馅包子
	b.create()
	c = GDFactory.Generate("3s") // 传入三鲜的参数，会返回广东包子铺的猪肉馅包子
	c.create()
}
