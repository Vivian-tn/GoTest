package main

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"html"
	"strings"
	"unicode/utf8"
)

func main() {
	summary := `<p  ></p><img src=\"v2-da80ecd74daa88abfcb14899a7af8f8d.jpg\" data-rawwidth=\"828\" data-rawheight=\"1792\" data-size=\"normal\" data-watermark=\"watermark\" data-original-src=\"v2-da80ecd74daa88abfcb14899a7af8f8d\" data-watermark-src=\"v2-0c9b414fdde3cbb7d6dc80890608012b\" data-private-watermark-src=\"v2-48c7689e32e8fde42a24a88c7a5f58e0\"/><p><br/></p><img src=\"v2-cb6d94424361538cd742337e8858caab.jpg\" data-rawwidth=\"828\" data-rawheight=\"1792\" data-size=\"normal\" data-watermark=\"watermark\" data-original-src=\"v2-cb6d94424361538cd742337e8858caab\" data-watermark-src=\"v2-4b90e79672c6f1f66608d02b092dca5c\" data-private-watermark-src=\"v2-ce6697d64498f28dfb6787422e77d434\"/><p data-pid=\"Ceo67xmd\">之前我也觉得悬来着，甚至准备了其他路径。</p><p data-pid=\"7Nbp6y29\">不过现在都在大宣传特宣传，工作组也去了，应该不至于了。</p><p data-pid=\"-_X7MuRd\">祝大家都能享受世界杯。</p>`
	fmt.Println("===========Replace", strings.Replace(summary, " ", "", -1))
	fmt.Println("----------", StripHTML(summary))
	if utf8.RuneCountInString(summary) >= 800 {
		fmt.Println("==========")
		summary = string([]rune(summary)[:800])
		fmt.Println("==========summary", summary)
	}

	fmt.Println("----------", html.EscapeString(StripHTML(summary)))
}
func StripHTML(s string) string {
	if s == "" {
		return s
	}
	policy := bluemonday.StrictPolicy()

	return strings.Trim(policy.Sanitize(s), " ")
}
