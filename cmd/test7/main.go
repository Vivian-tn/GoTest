package main

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {

	str := "[{\"button_title\":\"报名入口\",\"button_url\":\"https://ntce.neea.edu.cn/html1/folder/16013/15-1.htm\",\"module_name\":\"考试报名\",\"button_content\":\"2022-03-23 至 2022-04-09\",\"start_time\":\"1648008000\",\"end_time\":\"1649476800\",\"is_start\":false,\"header_city\":\"北京\"},{\"button_title\":\"查询入口\",\"button_url\":\"https://ntce.neea.edu.cn/html1/folder/16013/15-1.htm\",\"module_name\":\"考试查询\",\"button_content\":\"2022-03-23 至 2022-04-09\",\"start_time\":\"1648008000\",\"end_time\":\"1649476800\",\"is_start\":false,\"header_city\":\"北京\"},{\"button_title\":\"查询入口\",\"button_url\":\"https://ntce.neea.edu.cn/html1/folder/16013/15-1.htm\",\"module_name\":\"考试查询\",\"button_content\":\"2022-03-23 至 2022-04-09\",\"start_time\":\"1648008000\",\"end_time\":\"1649476800\",\"is_start\":false,\"header_city\":\"上海\"}]"
	//str := "22222#3333333"
	fmt.Println(strings.HasPrefix(str, "[{"))
	dbButtons := strings.Split(str, ";")
	for i := range dbButtons {
		buttonParts := strings.Split(dbButtons[i], "#")
		fmt.Println(len(buttonParts))
	}
	//uuid
	randBytes := make([]byte, 15)
	rand.Read(randBytes)
	fmt.Printf("%x\n", randBytes)

	//time convert
	loc, _ := time.LoadLocation("Asia/Shanghai")
	standardTime := "2006-01-02 15:04:05"
	//standardTime1 := "2006/1/2"
	startTime, _ := time.ParseInLocation("2006/1/2 15:04:05", "2018/5/31 09:22:19", loc)
	fmt.Println(startTime)
	fmt.Println(startTime.Format(standardTime))

	timeStr := time.Now().Format(standardTime)
	fmt.Println(timeStr)
	t, _ := time.ParseInLocation(standardTime, timeStr, time.Local)
	fmt.Println(t)
	timeUnix := t.Unix()
	fmt.Println(timeUnix)

	//timeNow := time.Unix(1646648746, 0)
	startTime1, _ := strconv.ParseInt("1646648746", 10, 64)
	timeString := time.Unix(startTime1, 0).Format("2006/1/2 15:04:05")
	fmt.Println(timeString)

	start, _ := strconv.ParseInt("1648080343", 10, 64)
	TimeString := time.Unix(start, 0).Format(standardTime)
	fmt.Println(TimeString)

}
