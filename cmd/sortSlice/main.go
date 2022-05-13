package main

import (
	"fmt"
	"sort"
	"strings"
)

//type Ordered interface {
//	type int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64,  string
//}
//
//type orderedSlice[T Ordered] []T
//
//func (s orderedSlice[T]) Len() int           { return len(s) }
//func (s orderedSlice[T]) Less(i, j int) bool { return s[i] < s[j] }
//func (s orderedSlice[T]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
//
//func OrderedSlice[T Ordered](s []T) {
//	sort.Sort(orderedSlice[T](s))
//}

func main() {
	//s1 := []int32{3, 5, 2}
	//fmt.Println(s1) // [3 5 2]
	//OrderedSlice(s1)
	//fmt.Println(s1) // [2 3 5]
	//
	//s2 := []string{"jim", "amy", "tom"}
	//fmt.Println(s2) // [jim amy tom]
	//OrderedSlice(s2)
	//fmt.Println(s2) // [amy jim tom]

	s3 := []string{"新疆", "天津", "上海", "安徽", "北京"}
	fmt.Println(s3) // [jim amy tom]
	sort.Slice(s3, func(i, j int) bool {
		return strings.Compare(s3[i], s3[j]) == -1
	})
	fmt.Println(s3) // [amy jim tom]
}
