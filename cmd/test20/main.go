package main

import (
	"fmt"
	"strconv"
)

type Intervene struct {
	ID            string
	Query         string
	Vertical      string
	Index         int32
	OriginID      string
	InterveneType string
}

func main() {
	// ret := []*Intervene{}
	// fmt.Println("1")
	// if ret == nil {
	// 	fmt.Println("1", ret)
	// }
	// a := []int{0, 1, 2, 3, 4, 5}

	// for _, v := range a {
	// 	if v == 5 {
	// 		fmt.Println("111")
	// 		continue
	// 	}
	// 	fmt.Println("222")
	// }

	// md5 := md5.NewMD5()
	// md5.WriteString("离群索居者，不是野兽，便是神灵。——亚里士多德")
	// hashId := md5.HexDigest()
	// fmt.Println(hashId)

	// var arr []interface{}
	// // arr = append(arr, 1)
	// fmt.Println(len(arr))

	// str := make(map[int]interface{})
	// str[1] = "24546"
	// fmt.Println(str[1].(int64))

	// ctx := context.Background()
	// ctx = context.WithValue(ctx, "1", "a")
	// ctx = context.WithValue(ctx, "1", "2")
	// v := ctx.Value("1")
	// v1 := ctx.Value("1")
	// fmt.Println(v)
	// fmt.Println(v1)

	// fmt.Println(strings.Contains("的下雨天文案", "下雨天文案"))

	// arr := []int{1, 2, 3}
	// fmt.Println(arr[:3])
	// var ar []int
	// fmt.Print(ar)
	// fmt.Println(append(arr, ar...))

	s := "213sss"
	pk, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pk)
}
