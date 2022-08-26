package util

import (
	"net/url"
	"path/filepath"
	"strings"
)

// URL2TokenOld is the old version of URL2Token. See benchmarks in test file
// for performance difference.
func URL2TokenOld(rawURL string, suffix bool) string {
	result, _ := url.Parse(rawURL)
	path := strings.TrimLeft(result.Path, "/")
	ext := filepath.Ext(path) // 扩展名

	if strings.Contains(path, "/") {
		index := strings.LastIndex(path, "/")
		path = path[index+1:]
	}
	if strings.Contains(path, ".") {
		index := strings.Index(path, ".")
		path = path[:index]
	}
	if strings.Contains(path, "_") {
		index := strings.Index(path, "_")
		path = path[:index]
	}

	if suffix {
		return path + ext
	}
	return path
}

// URL2Token returns modified base of file name of a url, if suffix is true,
// returns modified base with ext part
func URL2Token(rawURL string, suffix bool) string {
	result, _ := url.Parse(rawURL)

	path := filepath.Base(result.Path)
	ext := filepath.Ext(path)

	if index := strings.Index(path, "."); index > 0 {
		path = path[:index]
	}

	if index := strings.Index(path, "_"); index > 0 {
		path = path[:index]
	}

	if suffix {
		return path + ext
	}
	return path
}
