package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"time"
)

func main() {
	// 不存在则创建;存在则清空;读写模式;
	file, err := os.Create("./info_box_city_exam_schedule.csv")
	if err != nil {
		fmt.Println("open file is failed, err: ", err)
	}
	// 延迟关闭
	defer file.Close()

	// 写入UTF-8 BOM，防止中文乱码
	file.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(file)

	// 写入数据
	w.Write([]string{"老师编号", "老师姓名", "老师特长", "日期"})
	w.Write([]string{"t1", "老师1", "纯阳无极功", ` ` + time.Now().Format("2006-01-02 15:04:05")})
	w.Flush()

	// Map写入
	m := make(map[int][]string)
	m[0] = []string{"学生编号", "学生姓名", "学生特长"}
	m[1] = []string{"s1", "学生1", "乾坤大挪移"}
	m[2] = []string{"s2", "学生2", "乾坤大挪移"}
	m[3] = []string{"s3", "学生3", "乾坤大挪移"}
	m[4] = []string{"s4", "学生4", "乾坤大挪移"}
	m[5] = []string{"s5", "学生5", "乾坤大挪移"}
	m[6] = []string{"s6", "学生6", "乾坤大挪移"}
	m[7] = []string{"s7", "学生7", "乾坤大挪移"}
	m[8] = []string{"s8", "学生8", "乾坤大挪移"}
	m[9] = []string{"s9", "学生9", "乾坤大挪移"}
	m[10] = []string{"s10", "学生10", "乾坤大挪移"}
	//m[0] = []string{
	//	"城市（北京,天津,河北,山西,内蒙古,辽宁,吉林,黑龙江,上海,江苏,浙江,安徽,福建,江西,山东,河南,湖南,湖北,广东,广西,海南,重庆,四川,贵州,云南,西藏,陕西,甘肃,青海,宁夏,新疆,台湾,香港,澳门）",
	//	"考试日程模块名称（例：考试报名）",
	//	"配置文案",
	//	"活动开始时间（格式为“2021/5/13 12:00:00”，需要具体到秒）",
	//	"活动结束时间（格式为“2021/5/13 12:00:00”，需要具体到秒）",
	//	"按钮名称",
	//	"跳转链接",
	//}
	//m[1] = []string{
	//	"北京市",
	//	"考试报名",
	//	"2022-03-23 至 2022-04-09",
	//	"2022/3/23 00:00:00",
	//	"2022/4/9 00:00:00",
	//	"报名入口",
	//	"https://ntce.neea.edu.cn/html1/folder/16013/15-1.htm",
	//}

	// 按照key排序
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, key := range keys {
		w.Write(m[key])
		// 刷新缓冲
		w.Flush()
	}
}

var (
	BatchPutInfoButtonRow1 = []string{
		"城市（北京,天津,河北,山西,内蒙古,辽宁,吉林,黑龙江,上海,江苏,浙江,安徽,福建,江西,山东,河南,湖南,湖北,广东,广西,海南,重庆,四川,贵州,云南,西藏,陕西,甘肃,青海,宁夏,新疆,台湾,香港,澳门）",
		"考试日程模块名称（例：考试报名）",
		"配置文案",
		"活动开始时间（格式为“2021/5/13 12:00:00”，需要具体到秒）",
		"活动结束时间（格式为“2021/5/13 12:00:00”，需要具体到秒）",
		"按钮名称",
		"跳转链接",
	}
	BatchPutInfoButtonRow2 = []string{
		"北京市",
		"考试报名",
		"2022-03-23 至 2022-04-09",
		"2022/3/23 00:00:00",
		"2022/4/9 00:00:00",
		"报名入口",
		"https://ntce.neea.edu.cn/html1/folder/16013/15-1.htm",
	}
)

//func main() {
//	var filename = "./info_box_city_exam_schedule1.csv"
//	file, err := os.Create(filename)
//	if err != nil {
//		return
//	}
//	csvWriter := csv.NewWriter(file)
//	row1 := BatchPutInfoButtonRow1
//	row2 := BatchPutInfoButtonRow2
//	var newContent [][]string
//	newContent = append(newContent, row1)
//	newContent = append(newContent, row2)
//	err1 := csvWriter.WriteAll(newContent)
//	if err1 != nil {
//		return
//	}
//	csvWriter.Flush()
//	defer file.Close()
//}
