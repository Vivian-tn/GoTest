package main

import (
	"fmt"
)

func main() {
	fmt.Println(romanToInt("MCMXCIV"))
}
func romanToInt(s string) (ans int) {
	var byteMap = map[byte]int{
		'I': 1,
		'V': 5,
		'X': 10,
		'L': 50,
		'C': 100,
		'D': 500,
		'M': 1000,
	}
	for v := range s {
		if v < len(s)-1 && byteMap[s[v]] < byteMap[s[v+1]] {
			ans -= byteMap[s[v]]
		} else {
			ans += byteMap[s[v]]
		}

	}
	return ans
}
