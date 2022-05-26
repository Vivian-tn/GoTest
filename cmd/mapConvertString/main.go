package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	//cityList := make(map[string]struct{})
	cityList := map[string]struct{}{
		"会计师":   {},
		"审计师":   {},
		"金融分析师": {},
		"教授":    {},
		"教师":    {},
		"医生":    {},
		"医师":    {},
		"律师":    {},
		"法官":    {},
		"检察官":   {},
		"警察":    {},
		"交警":    {},
		"记者":    {},
		"新闻":    {},
		"日报":    {},
		"早报":    {},
		"晚报":    {},
		"周刊":    {},
		"智库":    {},
	}
	fmt.Println(MapToJson(cityList))
}

func MapToJson(param map[string]struct{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}
