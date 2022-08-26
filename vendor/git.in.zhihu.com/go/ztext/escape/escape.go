package escape

import "golang.org/x/net/html"

// XhtmlEscape just wrapes html.EscapString
func XhtmlEscape(value string) string {
	//performance boost for empty string. see benchmark tests in escape_test.go
	if len(value) == 0 {
		return value
	}

	return html.EscapeString(value)
}

// XhtmlUnescape just wrapes html.UnescapString
func XhtmlUnescape(value string) string {
	if len(value) == 0 {
		return value
	}

	return html.UnescapeString(value)
}
