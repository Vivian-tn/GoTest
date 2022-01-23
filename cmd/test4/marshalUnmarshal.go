package main

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name   string   `json:"name"`   // 姓名
	AgeId  int      `json:"age_id"` // 年龄
	Gender string   `json:"gender"` // 性别
	Score  float64  `json:"score"`  // 分数
	Course []string `json:"course"` // 课程
}

func main() {
	//将结构体序列化成json
	stu := Student{
		"张三",
		20,
		"男",
		78.6,
		[]string{"语文", "数学", "音乐"},
	}
	fmt.Println(stu)
	data, err := json.Marshal(&stu)
	if err != nil {
		fmt.Println("序列化错误", err)
	} else {
		fmt.Println(string(data))
	}

	//json反序列化成结构体
	var stu1 Student
	str := `{"name":"张三","age":20,"gender":"男","score":78.6,"course":["语文","数学","音乐"]}`

	err1 := json.Unmarshal([]byte(str), &stu1)
	if err1 != nil {
		fmt.Println("反序列化失败")
	}
	fmt.Println(stu1)
}
