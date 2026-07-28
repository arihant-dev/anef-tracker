package importer

import (
	"net/url"
	"strings"
)

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func ParseCookieHeader(cookieHeaderStr string, target map[string]string) {
	parts := strings.Split(cookieHeaderStr, ";")
	for _, p := range parts {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) == 2 {
			k := strings.TrimSpace(kv[0])
			val, err := url.QueryUnescape(kv[1])
			if err != nil {
				val = kv[1]
			}
			target[k] = val
		}
	}
}
