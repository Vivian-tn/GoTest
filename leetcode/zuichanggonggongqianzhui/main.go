package main

import (
	"fmt"
)

func main() {
	strs := []string{"flower", "flow", "flight"}
	fmt.Println(strs[0][0])
	fmt.Println(longestCommonPrefix(strs))
}
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	for i := 0; i < len(strs[0]); i++ { //字符串选取
		for j := 1; j < len(strs); j++ { //字符串长度
			if i == len(strs[j]) || strs[j][i] != strs[0][i] {
				return strs[0][:i]
			}
		}
	}
	return strs[0]

}
func minLen(x, y int) int {
	if x < y {
		return x
	}
	return y
}
