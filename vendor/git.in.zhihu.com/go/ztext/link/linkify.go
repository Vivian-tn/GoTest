package link

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"git.in.zhihu.com/go/ztext/node"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	// WhitespacePattern is pattern for white space
	WhitespacePattern *regexp.Regexp = regexp.MustCompile("\\s+")
	// URLPattern is pattern for url
	URLPattern *regexp.Regexp = makeURLPattern()
	//UrlPattern        *regexp.Regexp = makeURLPattern()
	tagPattern = regexp.MustCompile(`</?[^>]+>`)

	avoidElements = map[atom.Atom]struct{}{
		atom.Textarea: {},
		atom.Pre:      {},
		atom.Code:     {},
		atom.Head:     {},
		atom.Select:   {},
		atom.A:        {},
	}
	avoidHostPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)^localhost`),
		regexp.MustCompile(`(?i)\bexample\.(?:com|org|net)$`),
		regexp.MustCompile(`(?i)^127\.0\.0\.1$`),
	}
	avoidClasses = map[string]struct{}{
		"nolink": {},
	}
)

func makeURLPattern() *regexp.Regexp {
	tmpl := fmt.Sprintf(
		`(?im)(?P<body>(?:(?:(?:(%s)://)([a-z0-9\.\_]{1,20}@)?(?:[a-z0-9\-\_]{1,63}\.){1,6}%s)|(?:([a-z0-9\.\_]{1,20}@)?(?:[a-z0-9\-\_]{1,63}\.){1,6}%s))(?::\d{1,6})?(?:(?:\/[%s]{1,1000}(?:[%s]{0,1000})(?:[%s]{0,1000})(?:[%s]{0,1000}))?\/?))([^%s]|$)`,
		`http|https|ftp`, // protocols
		fmt.Sprintf(`(?:%s)`, strings.Join(TLDs, "|")),        // tlds
		fmt.Sprintf(`(?:%s)`, strings.Join(STRICT_TLDs, "|")), // tlds
		`\w\.\?\-%,=&;\+“”!~:'@\/\|#`,                         // any
		`\w\.\?\-%,=&;\+“”!~:'@\/\|#`,                         // any
		`\w\.\?\-%,=&;\+“”!~:'@\/\|#`,                         // any
		`\w\.\?\-%,=&;\+“”!~:'@\/\|#`,                         // any
		`\w\.\?\-%,=&;\+“”!~:'@\/\|#`,                         // any
	)

	return regexp.MustCompile(tmpl)
}

// UnifyLink makes link unified. eg. www.zhihu.com -> http://www.zhihu.com
func UnifyLink(text string) string {
	converter := func(s string) string {
		submatches := URLPattern.FindAllStringSubmatch(s, -1)
		url, proto, login, tail := submatches[0][1], submatches[0][2], submatches[0][4], submatches[0][5]
		href := url
		protoLen := 0
		if len(proto) == 0 {
			href = "http://" + href
		} else {
			protoLen = len(proto) + 3 // +3 for "://"
		}

		parts := strings.Split(url[protoLen:], "/")
		hostPart := parts[0]

		if len(login) > 0 && hostPart == url {
			return url + tail
		}
		return href + tail
	}

	match := tagPattern.FindAllStringIndex(text, -1)
	if match == nil {
		// not "< >" tag in text
		return URLPattern.ReplaceAllStringFunc(text, converter)
	}

	start := 0
	var sb strings.Builder
	for _, m := range match {
		s, e := m[0], m[1]
		sb.WriteString(URLPattern.ReplaceAllStringFunc(text[start:s], converter))
		sb.WriteString(text[s:e])
		start = e
	}
	sb.WriteString(text[start:])
	return sb.String()
}

// AutoLink turns any URLs into links. Reimplementation of
// lxml.html.clean.autolink function
func AutoLink(root *node.Node, linkRegexes []*regexp.Regexp) {
	// pass avoidElements
	if _, ok := avoidElements[root.DataAtom]; ok {
		return
	}

	// pass avoidClasses
	if className, ok := root.GetAttr("class"); ok {
		names := strings.Split(className, " ")
		for _, name := range names {
			if _, ok = avoidClasses[name]; ok {
				return
			}
		}
	}

	for child := (*node.Node)(root.FirstChild); child != nil; child = (*node.Node)(child.NextSibling) {
		AutoLink(child, linkRegexes)
		if tail := child.Tail(); tail != "" {
			text, tailChildren := linkText(tail, linkRegexes, avoidHostPatterns)
			if len(tailChildren) > 0 {
				_ = child.SetTail(text)
				// we should insert tail child after child's tail node.
				lastInserted := (*node.Node)(child.NextSibling)
				for _, tail := range tailChildren {
					_ = root.InsertAfter(tail, lastInserted)
					lastInserted = tail
				}
			}
		}
	}
	if rootText := root.Text(); rootText != "" {
		text, preChildren := linkText(rootText, linkRegexes, avoidHostPatterns)
		if len(preChildren) > 0 {
			if root.Type == html.TextNode {
				parentNode := (*node.Node)(root.Parent)
				parentNode.SetText(text)
				lastInserted := root
				for _, child := range preChildren {
					_ = parentNode.InsertAfter(child, lastInserted)
					lastInserted = child
				}
			} else {
				root.SetText(text)
				originalFirstChild := (*node.Node)(root.FirstChild)
				for _, child := range preChildren {
					root.InsertBefore(child, originalFirstChild)
				}
			}
		}
	}
}

// Reimplement for lxml/html/clean.py:_link_text
func linkText(text string, linkRegexes []*regexp.Regexp, avoidHostPatterns []*regexp.Regexp) (string, []*node.Node) {
	leadingText := ""
	links := []*node.Node{}

	for {
		var bestMatchedIndex []int
		var bestMatchedPattern *regexp.Regexp
	LinkPatternLoop:
		// TODO(taozle) support linkPattern contains multi items.
		for _, linkRegex := range linkRegexes {
			matchIndexes := linkRegex.FindAllStringSubmatchIndex(text, -1)
			if matchIndexes == nil {
				continue
			}

			for _, matchIndex := range matchIndexes {
				host := ""
				for idx, groupName := range linkRegex.SubexpNames() {
					if groupName == "host" {
						startPos := matchIndex[idx*2]
						endPos := matchIndex[idx*2+1]
						host = text[startPos:endPos]
					}
				}
				if host == "" {
					continue
				}

				hostMatched := true
				for _, hostPattern := range avoidHostPatterns {
					if hostPattern.FindString(host) != "" {
						hostMatched = false
						break
					}
				}

				if hostMatched {
					bestMatchedIndex = matchIndex
					bestMatchedPattern = linkRegex
					break LinkPatternLoop
				}
			}
		}

		if bestMatchedIndex == nil {
			if len(links) > 0 {
				// assert not links[-1].tail

				// no parent associated, so we can't use SetTail method,
				// here we manually add a textnode to the end.
				// original code: links[len(links)-1].SetTail(text)
				links = append(links, node.NewTextNode(text))
			} else {
				// assert not leading_text
				leadingText = text
			}
			break
		}

		bestMatchedGroups := make(map[string]string)
		for idx, groupName := range bestMatchedPattern.SubexpNames() {
			if idx != 0 && groupName != "" {
				startPos := bestMatchedIndex[idx*2]
				endPos := bestMatchedIndex[idx*2+1]
				bestMatchedGroups[groupName] = text[startPos:endPos]
			}
		}

		link := text[bestMatchedIndex[0]:bestMatchedIndex[1]]
		end := bestMatchedIndex[1]
		if strings.HasSuffix(link, ".") || strings.HasSuffix(link, ",") {
			end--
			link = link[:len(link)-1]
		}

		prevText := text[:bestMatchedIndex[0]]
		if len(links) > 0 {
			//assert not links[-1].tail

			// no parent associated, so we can't use SetTail method,
			// here we manually add a textnode to the end.
			// original code: links[len(links)-1].SetTail(text)
			links = append(links, node.NewTextNode(prevText))
		} else {
			// assert not leading_text
			leadingText = prevText
		}

		// FIXME(taozle) passing a node factory?
		anchor := node.NewNode(atom.A, "href", link)
		body, ok := bestMatchedGroups["body"]
		if !ok {
			body = link
		}
		if strings.HasSuffix(body, ".") || strings.HasSuffix(body, ",") {
			body = body[:len(body)-1]
		}
		anchor.SetText(body)
		links = append(links, anchor)
		text = text[end:]
	}

	return leadingText, links
}

// IsLinkOfDomain detects if link is url of given domain
func IsLinkOfDomain(link, domain string) bool {
	linkInfo, err := url.Parse(link)
	if err != nil {
		return false
	}

	return linkInfo.Host == domain || strings.HasSuffix(linkInfo.Host, "."+domain)
}

// IsLinkUseSpecialScheme detects if link is url of given scheme
func IsLinkUseSpecialScheme(link, scheme string) bool {
	linkInfo, err := url.Parse(link)
	if err != nil {
		return false
	}

	return linkInfo.Scheme == scheme
}

// IsOutLinkDomain detects if link is url of given domains
func IsOutLinkDomain(link string, domains []string) bool {
	linkInfo, err := url.Parse(link)
	if err != nil {
		return false
	}

	for _, domain := range domains {
		if linkInfo.Host == domain {
			return true
		}
	}

	return false
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// ShortenLink ...
// If url is illegal, then return at most first 23 characters of url.
func ShortenLink(link string) (string, error) {
	parsed, err := url.Parse(link)
	if err != nil {
		return link[:min(23, len(link))], nil
	}

	netloc := parsed.Host
	path := parsed.Path
	if !strings.Contains(link, path) {
		// url.Parse will unescape the provided string, so here we should escape it to original text.
		// FIXME(taozle) find a better way to handle this.
		path = (&url.URL{Path: path}).EscapedPath()
	}
	netloc = strings.TrimPrefix(netloc, "www.")

	maxLen := 23 - len(netloc)
	if maxLen < 0 {
		maxLen = 0
	}

	shortended := false
	if len(path) > maxLen {
		shortended = true
		path = path[:maxLen]
	}

	//if parsed.params:
	//if parsed.RawQuery {
	//	return netloc + "/...", nil
	//}

	if !shortended && len(parsed.RawQuery) > 0 {
		shortended = true
		path += "?"
	} else if !shortended && len(parsed.Fragment) > 0 {
		shortended = true
		path += "#"
	}

	if shortended {
		path += "..."
	}

	return netloc + path, nil
}
