package main

import (
	"fmt"
	"strconv"
)

type DocItems struct {
	DocId   int64
	DocType string
}
type ContentAssessScene int64

func (d ContentAssessScene) Int64() int64 {
	return int64(d)
}

func (d ContentAssessScene) String() string {
	return strconv.FormatInt(int64(d), 10)
}

// ref: 敏感词场景 https://zhihu.kdocs.cn/l/crFx8uPhab5C
const (
	DitingSceneSearchEntrance    ContentAssessScene = 100 // 搜索主入口
	DitingSceneSuggest           ContentAssessScene = 101 // 建议词 sug
	DitingSceneRecommendCard     ContentAssessScene = 102 // 搜索卡片类关联性推荐(澄清词+相关搜索+点击后推荐)
	DitingScenePositiveRecommend ContentAssessScene = 103 // 主动推荐（搜索发现+实体词）
	DitingSceneDocRelevant       ContentAssessScene = 104 // 被动搜索-doc关联性推荐（首页回退+评论页）
	DitingScenePresetWord        ContentAssessScene = 105 // 预置词
	DitingScenePrivilegedUser    ContentAssessScene = 106 // 用户白名单
	DitingSceneQueryPrefilter    ContentAssessScene = 107 // 搜索卡片类关联性推荐(澄清词+相关搜索+点击后推荐) - 前置过滤
	DitingSceneRecQuery          ContentAssessScene = 109 // 被动搜索-doc关联性推荐（边看边搜+问题页）
)

func main() {
	//fmt.Println("嗨客网(www.haicoder.net)")
	////使用 ToLower() 函数，将字符串转成小写
	//strHaiCoder := "Study Golang From HaiCoder"
	//lowerStr := strings.ToLower(strHaiCoder)
	//titleStr := strings.ToTitle(strHaiCoder)
	//fmt.Println("lowerStr =", lowerStr)
	//fmt.Println("titleStr =", titleStr)
	//
	//md5 := md5.NewMD5()
	//md5.WriteInt64(112916211)
	//md5.WriteInt64(536672553)
	//md5.WriteString("answer")
	//md5.WriteFloat64(float64(time.Now().UnixNano()) / 1000000000)
	//requestHashId := md5.HexDigest()
	//fmt.Println("0:requestHashId:", requestHashId)
	//
	//md5 := md5.NewMD5()
	//md5.WriteInt64(112916211)
	//md5.WriteInt64(190931088)
	//md5.WriteString("article")
	//md5.WriteFloat64(float64(time.Now().UnixNano()) / 1000000000)
	//requestHashId := md5.HexDigest()
	//fmt.Println("1:requestHashId:", requestHashId)
	//
	//idTypeMap := make(map[string]struct{}, 0)
	//idTypeMap["aaa"] = struct{}{}
	//idTypeMap["bbb"] = struct{}{}
	//idTypeMap["ccc"] = struct{}{}
	//str, _ := json.Marshal(idTypeMap)
	//log.Infof("===========Here RecQueryCache BatchGetRecQueryListFromDoc, idTypeMap:%v", string(str))

	//testMAp := make(map[DocItems]int64, 0)
	//testMAp[DocItems{DocId: 1, DocType: "a"}] = 1
	//testMAp[DocItems{DocId: 2, DocType: "b"}] = 2
	//tt := DocItems{DocId: 1, DocType: "a"}
	//ttt := DocItems{DocId: 1, DocType: "a"}
	//if haha, ok := testMAp[tt]; ok {
	//	fmt.Println(haha)
	//}
	//if tt == ttt {
	//	fmt.Println("=========")
	//}
	//fmt.Println(DitingSceneSearchEntrance.String())
	first := "社区"
	fmt.Println([]rune(first))
	fmt.Println([]byte(first))
}
