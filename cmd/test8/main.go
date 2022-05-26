package main

import (
	"fmt"
	"sort"
)

func main() {
	//b := []string{"新疆", "天津", "上海", "安徽", "北京"}
	//
	//sort.Strings(b)
	////output: [, ha 呼 哈 嚯]
	//fmt.Println("Default sort: ", b)
	//
	//sort.Sort(pkg.ByPinyin(b))
	////output: [, ha 哈 呼 嚯]
	//fmt.Println("By Pinyin sort: ", b)

	//strings := []string{"hello", "good", "students", "morning", "people", "world", "a"}
	//sort.Strings(strings)
	//fmt.Println("after sort by Strings:\t", strings) // [good hello mornig people students world]

	//var a string
	//fmt.Println("is", a)
	//
	//var mapppp map[string]string
	////mapp["a"] = "aaa"
	//if mapppp != nil {
	//	fmt.Println("is", mapppp)
	//}
	//
	//mappp := make(map[string]string)
	//fmt.Println("is", mappp)
	//
	//var sum *int
	//sum = new(int) //分配空间
	//*sum = 98
	//fmt.Println(*sum)

	//var mapp map[string]string
	//var p interface{} = mapp
	//value, ok := p.(map[string]string)
	//if !ok {
	//	fmt.Println("It's not ok for type string")
	//	return
	//}
	//if value != nil {
	//	fmt.Println("The value is ", value)
	//}

	//models := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	//ret := []int64{}
	//target := 0
	//k := 9
	//if len(models) == 10 { //删除元素4
	//	ret = append(ret, models[:target]...)
	//	ret = append(ret, models[target+1:]...)
	//} else { //4移动位置到slice下标1
	//
	//	ret = append(ret, models[:target]...)
	//	sliceList1 := models[target:k]
	//	sliceList2 := models[k+1:]
	//	ret = append(ret, models[k])
	//	ret = append(ret, sliceList1...)
	//	ret = append(ret, sliceList2...)
	//}
	//fmt.Println(ret)

	//models := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	//ret := []int64{}
	//k := 0
	//target := 8
	//if len(models) != 10 { //删除元素4
	//	ret = append(ret, models[:k]...)
	//	ret = append(ret, models[k+1:]...)
	//} else { //4移动位置到slice下标1
	//	var max, min int
	//	if k < target {
	//		max = target
	//		min = k
	//		ret = append(ret, models[:min]...)
	//		sliceList1 := models[min+1 : max+1]
	//		sliceList2 := models[max+1:]
	//		ret = append(ret, sliceList1...)
	//		ret = append(ret, models[min])
	//		ret = append(ret, sliceList2...)
	//	} else {
	//		max = k
	//		min = target
	//		ret = append(ret, models[:min]...)
	//		sliceList1 := models[min:max]
	//		sliceList2 := models[max+1:]
	//		ret = append(ret, models[max])
	//		ret = append(ret, sliceList1...)
	//		ret = append(ret, sliceList2...)
	//	}
	//
	//}
	//fmt.Println(ret)

	//demo := []int{1, 3, 5, 7, 9, 8, 2}
	//
	//swap(demo, 1, 3)
	//fmt.Println(demo)

	//s := []int{0, 1, 2, 3}
	//ret := s[:2]
	//fmt.Println(ret)
	//fmt.Println(cap(s))
	//fmt.Println(len(s))
	//fmt.Println(s[2])
	//ret = append(ret, s[3:]...)
	//fmt.Println(ret)
	//s1 := make([]int, 0)
	//fmt.Println("s1 length: ", len(s1))   // s1 length:  3
	//fmt.Println("s1 capacity: ", cap(s1)) // s1 capacity:  3
	//fmt.Printf("s1 value: %#v\n", s1)     // s1 value: []int{0, 0, 0}

	//for k := range s {
	//	fmt.Println(k)
	//	s = []int{4, 5, 6}
	//}
	//fmt.Println(s)

	m := make(map[int]string)
	m[1] = "a"
	m[2] = "c"
	m[0] = "b"

	// To store the keys in slice in sorted order
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// To perform the opertion you want
	for _, k := range keys {
		fmt.Println("Key:", k, "Value:", m[k])
	}

}

func swap(s []int, i, j int) {
	s[i], s[j] = s[j], s[i]
}
