package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Student struct {
	Name   string   `json:"name"`   // 姓名
	Age    int      `json:"age"`    // 年龄
	Gender string   `json:"gender"` // 性别
	Score  float64  `json:"score"`  // 分数
	Course []string `json:"course"` // 课程
}

type StudentForm struct {
	Name   *string  `json:"name"`   // 姓名
	Age    *int     `json:"age"`    // 年龄
	Gender *string  `json:"gender"` // 性别
	Score  *float64 `json:"score"`  // 分数
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

	fmt.Println(stu) //{张三 20 男 78.6 [语文 数学 音乐]}
	data, err := json.Marshal(&stu)
	if err != nil {
		fmt.Println("序列化错误", err)
	} else {
		fmt.Println(data)         //[123 34 110 97 109 101 34 58 34 229 188 160 228 184 137 34 44 34 97 103 101 34 58 50 48 44 34 103 101 110 100 101 114 34 58 34 231 148 183 34 44 34 115 99 111 114 101 34 58 55 56 46 54 44 34 99 111 117 114 115 101 34 58 91 34 232 175 173 230 150 135 34 44 34 230 149 176 229 173 166 34 44 34 233 159 179 228 185 144 34 93 125]
		fmt.Println(string(data)) //{"name":"张三","age":20,"gender":"男","score":78.6,"course":["语文","数学","音乐"]}
	}
	var ss Student
	fmt.Println(ss)
	//json反序列化成结构体
	var stu1 Student
	fmt.Printf("????/%+v\n", stu1)
	str := `{"name":"张三","age":20,"gender":"男","score":78.6,"course":["语文","数学","音乐"]}`

	err1 := json.Unmarshal([]byte(str), &stu1)
	if err1 != nil {
		fmt.Println("反序列化失败")
	}
	fmt.Printf("%v", stu1) //{张三 20 男 78.6 [语文 数学 音乐]}

	loc, _ := time.LoadLocation("Asia/Shanghai")
	startTime, _ := time.ParseInLocation("2006/1/2 15:04:05", "2022/4/9 12:00:00", loc)
	fmt.Println(strconv.FormatInt(startTime.Unix(), 10))
}
