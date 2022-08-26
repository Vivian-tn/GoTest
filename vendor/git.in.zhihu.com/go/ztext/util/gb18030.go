package util

import (
	"bytes"
	"io/ioutil"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

type rmFunc func(r rune) bool

func (rm rmFunc) Contains(r rune) bool {
	return rm(r)
}

// FromGB18030 convert text string of gb18030 encoding to utf8
func FromGB18030(text string) string {

	filter := transform.Chain(
		simplifiedchinese.GBK.NewDecoder(),
		runes.Remove((rmFunc)(func(r rune) bool {
			return r == utf8.RuneError
		})),
	)
	reader := transform.NewReader(strings.NewReader(text), filter)

	var buf bytes.Buffer
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return ""
	}
	return buf.String()
}

// ToGB18030 convert text string to gb18030 encoding.
func ToGB18030(text string) string {
	reader := transform.NewReader(strings.NewReader(text), simplifiedchinese.GB18030.NewEncoder())

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return text
	}

	return string(data)
}

// CutSentences cut string by rune count
func CutSentences(text string, n int) string {
	if n <= 0 {
		return ""
	}

	if n >= utf8.RuneCountInString(text) {
		return text
	}

	realPos := 0
	for idx := range text {
		if n == 0 {
			realPos = idx
			break
		}
		n--
	}

	//BenchmarkCutSentencesByConvertingRune      	10000000	       207 ns/op
	//BenchmarkCutSentences							50000000	        32.6 ns/op
	// return string([]rune(text)[:n])
	return text[:realPos]
}

// GB18030Len returns len of string in GB18030 encoding.
func GB18030Len(text string) int {
	return len(ToGB18030(text))
}

func addAcc(acc, idx, prev int) int {
	delta := idx - prev
	if delta <= 2 {
		// 如果字符长度 <= 2字节，累加实际长度
		acc += delta
	} else {
		// 对于中文，emoji 这类3-4字节的字符，按2计算
		acc += 2
	}
	return acc
}

// GB18030Cut cuts text by n bytes, if the ending point is in multichar character, discard it.
// 英文一个字符算一个字节
// 中文及 emoji，一个字符算两个字节
func GB18030Cut(text string, n int) string {

	if n <= 0 {
		return ""
	}

	realPos := 0
	acc := 0
	for idx := range text {
		acc = addAcc(acc, idx, realPos)
		if n < acc {
			return text[:realPos]
		}
		realPos = idx
	}

	// process the last character
	acc = addAcc(acc, len(text), realPos)
	if acc > n {
		return text[:realPos]
	}

	return text
}
