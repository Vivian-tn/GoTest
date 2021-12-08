package pkg

import (
	"fmt"
	"regexp"
)

func ZhengZe() {
	str := "好喝的蜜桃乌龙茶，哈哈哈"
	compileRegex := regexp.MustCompile("好喝的(.*?)，哈哈哈") // 正则表达式的分组，以括号()表示，每一对括号就是我们匹配到的一个文本，可以把他们提取出来。
	matchArr := compileRegex.FindStringSubmatch(str)   // FindStringSubmatch 方法是提取出匹配的字符串，然后通过[]string返回。我们可以看到，第1个匹配到的是这个字符串本身，从第2个开始，才是我们想要的字符串。
	fmt.Println("提取字符串内容：", matchArr[len(matchArr)-1]) // 输出：蜜桃乌龙茶

	return
}
