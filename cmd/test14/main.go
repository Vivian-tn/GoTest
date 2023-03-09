package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	// str := [][]string{}
	// str = append(str, []string{"q", "a"})
	// str = append(str, []string{"c", "v"})
	// str = append(str, []string{"n", "l"})
	// fmt.Println(len(str))
	// for _, strings := range str {
	// 	fmt.Println(strings)
	// }

	str1 := "hello world"
	fmt.Println(utf8.RuneCountInString(str1))

	str2 := "中國"
	fmt.Println(utf8.RuneCountInString(str2))

	str3 := ","
	fmt.Println(utf8.RuneCountInString(str3))

}
