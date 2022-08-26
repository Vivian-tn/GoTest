package ztext

import (
	"fmt"
	"strings"
	"unicode/utf8"
	"unsafe"

	"git.in.zhihu.com/go/ztext/escape"
	"git.in.zhihu.com/go/ztext/link"
	"git.in.zhihu.com/go/ztext/util"
)

func subBlankWhite(text string) string {
	if len(text) == 0 {
		return ""
	}

	text1 := strings.Replace(text, "　", " ", -1)
	// strip 0x00A0, you can type 'ga' in normal mode in vim to see code point
	text2 := strings.Replace(text1, " ", " ", -1)
	return link.WhitespacePattern.ReplaceAllString(text2, " ")
}

// SinaLen 按照新浪的规则，判断字符长度。
// 规则:
// (1) 2个半角算1个字符
// (2) 链接不论多长,算10个字符
func SinaLen(text string) int {
	blankStrippedText := subBlankWhite(text)
	urlStrippedText := link.URLPattern.ReplaceAllString(
		blankStrippedText,
		"this_is_len20_string",
	)
	textLength := util.GB18030Len(urlStrippedText)

	return textLength/2 + textLength%2
}

// SkeletonExcerpt returns excerpt on template content given the template and
// destination length. '...' were added to the excerpted arguments.
// IMPORTANT NOTES: only supports %v
func SkeletonExcerpt(template string, length int, args ...interface{}) string {
	var lenFunc func(string) int
	if length == 140 {
		lenFunc = SinaLen
	} else {
		lenFunc = utf8.RuneCountInString
	}

	// in fact, type of escapedArgs is []string
	escapedArgs := make([]interface{}, 0, len(args))
	var sb strings.Builder
	for _, a := range args {
		newArg := ""
		if s, ok := a.(string); ok {
			newArg = escape.XhtmlEscape(s)
		} else {
			newArg = fmt.Sprintf("%v", a)
		}
		sb.WriteString(newArg) // nolint
		escapedArgs = append(escapedArgs, newArg)
	}
	argLen := lenFunc(sb.String())

	completeText := fmt.Sprintf(template, escapedArgs...)
	completeLen := lenFunc(completeText)
	completeText = subBlankWhite(completeText)
	if completeLen <= length {
		return completeText
	}

	allowedArgLen := argLen - (completeLen - length)
	newLen := 0
	newArgs := make([]interface{}, 0)

	for _, a := range escapedArgs {
		if newLen < allowedArgLen {
			aStr := a.(string)
			if newLen+lenFunc(aStr) > allowedArgLen {
				nt := excerpt(aStr, allowedArgLen-newLen)
				newArgs = append(newArgs, nt)
				newLen += lenFunc(nt)
			} else {
				newArgs = append(newArgs, aStr)
				newLen += lenFunc(aStr)
			}
		} else {
			newArgs = append(newArgs, "")
		}

	}

	return fmt.Sprintf(template, newArgs...)
}

// linkSeg segments the plain text without breaking the link text.
func linkSeg(text string) []int {
	segLen := []int{}

	matchPos := []int{} // [start1, end1, start2, end2]
	textBytes := []byte(text)

	pos := 0
	seg := link.URLPattern.FindIndex(textBytes)
	for seg != nil {
		newPos := pos + seg[1]
		matchPos = append(matchPos, pos+seg[0], newPos)
		seg = link.URLPattern.FindIndex(textBytes[newPos:])
		pos = newPos
	}

	// no match, return [1,1,1,1,1]
	if len(matchPos) == 0 {
		for range text {
			segLen = append(segLen, 1)
		}
		return segLen
	}

	pos = 0
	for i := 0; i < len(matchPos)-1; i += 2 {
		start, end := matchPos[i], matchPos[i+1]
		//before := string(textBytes[pos:start])
		//url := string(textBytes[start:end])
		beforeBytes := textBytes[pos:start]
		urlBytes := textBytes[start:end]
		before := *(*string)(unsafe.Pointer(&beforeBytes))
		url := *(*string)(unsafe.Pointer(&urlBytes))
		for range before {
			segLen = append(segLen, 1)
		}
		segLen = append(segLen, utf8.RuneCountInString(url))
		pos = end
	}

	//after := string(textBytes[pos:])
	afterBytes := textBytes[pos:]
	after := *(*string)(unsafe.Pointer(&afterBytes))
	for range after {
		segLen = append(segLen, 1)
	}

	return segLen
}

// zdom.excerpt() is deprecated, so we only make a private function here.
func excerpt(content string, length int) string {
	zdom := ZText2(content)
	if zdom == nil {
		return ""
	}
	text, _ := RenderPlainText(zdom)
	if text == "" {
		return text
	}
	if utf8.RuneCountInString(text) <= length {
		return text
	}

	runeCount := 0
	for _, segLen := range linkSeg(text) {
		tempLen := runeCount + segLen
		if tempLen+1 >= length { // 1 means len("…")
			break
		}
		runeCount = tempLen
	}

	return util.CutSentences(text, runeCount) + "…"
}
