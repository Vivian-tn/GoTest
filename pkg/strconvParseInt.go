package pkg

import (
	"strconv"
)

/*
参数1 数字的字符串形式

参数2 数字字符串的进制 比如二进制 八进制 十进制 十六进制

参数3 返回结果的bit大小 也就是int8 int16 int32 int64
*/
func StrconvParseInt() {
	i, err := strconv.ParseInt("123", 10, 32)
	if err != nil {
		panic(err)
	}
	println(i)
}
func Int(i int) *int { return &i }
