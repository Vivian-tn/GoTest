package main

import (
	"fmt"
	"strings"
)

func main() {
	dbButton := "1#2#3##4"
	buttonParts := strings.Split(dbButton, ";")
	fmt.Println(len(buttonParts))
	//var chosedMember {}
	//fmt.Println(chosedMember)
	arr := []string{"a", "b"}
	for k, v := range arr {
		fmt.Println(k, v)
	}
	var s3 []string
	s1 := "aa"
	s2 := "aa"
	s3 = append(s3, s1, s2)
	fmt.Println(s3)
	batchPostRow1 := []string{
		"query集(多个query以半角逗号分隔)",
		"头部标题(不超过9个字符)",
		"头部跳转",
		"头部背景(填写色值)",
		"头部描述",
		"生效时间(格式为“2021/5/13 12:00:00”，需要具体到秒)",
		"过期时间(格式为“2021/5/13 12:00:00”，需要具体到秒)",
		"内容模块(格式为“标题(单分类则不填，但是冒号要保留):内容类型(包括question, answer, article, zvideo四种):对应token;标题:内容类型:对应token”)",
	}
	batchPostRow1 = append(batchPostRow1, "aaa")
	fmt.Println(batchPostRow1)

	contentTitles := make(map[string]struct{})
	contentTitles["aa"] = struct{}{}
	contentTitles["aa"] = struct{}{}
	contentTitles["aa"] = struct{}{}
	fmt.Println(contentTitles)

}
