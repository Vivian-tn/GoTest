package pier

import (
	gfmt "fmt"
	gurl "net/url"
	"path/filepath"
	"strings"
)

var (
	domainScheduler DomainScheduler = &ModulusScheduler{}
)

func getURL(token, size, fmt, quality string) string {
	token, fmtVal := parseToken(token)

	domain := domainScheduler.GetDomain(GetCDNDomainsByLRU(), token)

	if fmt == "" {
		fmt = fmtVal
	}

	if size != "" {
		token = gfmt.Sprintf("%s_%s.%s", token, size, fmt)
	} else {
		token = gfmt.Sprintf("%s.%s", token, fmt)
	}

	if quality == "" {
		return strings.Join([]string{domain, token}, "/")
	} else {
		return strings.Join([]string{domain, quality, token}, "/")
	}
}

func GetFullURL(token, size, fmt, quality string, secure bool) string {
	if token == "" {
		return ""
	}

	if isURL(token) {
		return token
	}

	prefix := "https://"
	if !secure {
		prefix = "http://"
	}

	return prefix + getURL(token, size, fmt, quality)
}

func GetURL(token, size, fmt, quality string) string {
	if token == "" {
		return ""
	}

	if isURL(token) {
		s := strings.Replace(token, "http://", "", 1)
		s = strings.Replace(s, "https://", "", 1)
		s = strings.Replace(s, "//", "", 1)
		return s
	}

	return getURL(token, size, fmt, quality)
}

func URL2Token(url string) (token string, suffix string, err error) {
	parsedUrl, err := gurl.Parse(url)
	if err != nil {
		return
	}
	path := strings.TrimLeft(parsedUrl.Path, "/")
	suffix = strings.TrimLeft(filepath.Ext(path), ".")
	if strings.Contains(path, "/") {
		lastIndex := strings.LastIndex(path, "/")
		path = path[lastIndex+1:]
	}
	if strings.Contains(path, ".") {
		index := strings.Index(path, ".")
		path = path[:index]
	}
	if strings.Contains(path, "_") {
		index := strings.Index(path, "_")
		path = path[:index]
	}
	token = path
	return
}
