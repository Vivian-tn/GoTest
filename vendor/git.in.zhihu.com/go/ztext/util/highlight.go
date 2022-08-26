package util

import (
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func highlight(source, lexer string) (string, string, error) {
	// Determine lexer.
	l := lexers.Get(lexer)
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Determine formatter.
	f := html.New(html.WithClasses(true))

	// Determine style.
	s := styles.Fallback

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return "", lexer, err
	}

	b := GetBuffer()
	defer PutBuffer(b)

	err = f.Format(b, s, it)
	if err != nil {
		return "", lexer, err
	}
	str := strings.Replace(b.String(), ` class="chroma"`, "", 1)
	return `<div class="highlight">` + str + "</div>", lexer, nil
}

// Highlight use lang name to highlight given source, it returns
// hilighted source, real lang name, and error
func Highlight(source, lang string) (string, string, error) {
	highlighted, lang, err := highlight(source, lang)
	if err != nil {
		// try text highlight
		lang = "text"
		highlighted, lang, err = highlight(source, lang)
		if err != nil {
			return "", lang, err
		}
	}
	return highlighted, lang, err
}
