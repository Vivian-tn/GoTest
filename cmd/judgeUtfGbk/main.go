package main

import (
	"fmt"
	"io/ioutil"

	"github.com/saintfish/chardet"
)

func isGBK(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		if data[i] <= 0x7f {
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}
func preNUm(data byte) int {
	var mask byte = 0x80
	var num int = 0
	//8bit中首个0bit前有多少个1bits
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}
func isUtf8(data []byte) bool {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			return false
		}
	}
	return true
}
func main() {
	b, err := ioutil.ReadFile(`/Users/tianning/go/testNew/_common_intervene_template (2).csv`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(isUtf8(b))
	rawBytes := []byte("/Users/tianning/go/testNew/_common_intervene_template.csv")
	detector := chardet.NewTextDetector()
	charset, err := detector.DetectBest(rawBytes)
	if err != nil {
		panic(err)
	}

	println(charset.Charset)
	println(charset.Language)
}
func validUTF8(buf []byte) bool {
	nBytes := 0
	for i := 0; i < len(buf); i++ {
		if nBytes == 0 {
			if (buf[i] & 0x80) != 0 { //与操作之后不为0，说明首位为1
				for (buf[i] & 0x80) != 0 {
					buf[i] <<= 1 //左移一位
					nBytes++     //记录字符共占几个字节
				}

				if nBytes < 2 || nBytes > 6 { //因为UTF8编码单字符最多不超过6个字节
					return false
				}

				nBytes-- //减掉首字节的一个计数
			}
		} else { //处理多字节字符
			if buf[i]&0xc0 != 0x80 { //判断多字节后面的字节是否是10开头
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}
