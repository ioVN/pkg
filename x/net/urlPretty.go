// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package net

import (
	"fmt"
	"strings"
)

func UrlPrettyParse(url string, schemeDefault ...string) string {
	for strings.HasPrefix(url, "/") {
		url = strings.TrimPrefix(url, "/")
	}
	for strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	if strings.Contains(url, "\t") ||
		strings.Contains(url, "\r") ||
		strings.Contains(url, "\n") {
		url = strings.ReplaceAll(url, "\t", "")
		url = strings.ReplaceAll(url, "\r", "")
		url = strings.ReplaceAll(url, "\n", "")
	}
	for i := 0; i < len(url)-5 && strings.Contains(url, "//"); i++ {
		url = strings.ReplaceAll(url, "//", "/")
		url = strings.ReplaceAll(url, ":/", "://")
	}
	if !strings.HasPrefix(url, "http://") &&
		!strings.HasPrefix(url, "https://") &&
		!strings.HasPrefix(url, "ws://") &&
		!strings.HasPrefix(url, "wss://") {
		if len(schemeDefault) == 0 {
			schemeDefault = []string{"http://"}
		}
		url = fmt.Sprintf("%s%s", schemeDefault[0], url)
	} else {
		url = fmt.Sprintf("%s", url)
	}
	return strings.ReplaceAll(url, " ", "")
}
