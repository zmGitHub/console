package common

import (
	"net/url"
)

func GetHostFromURL(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}

	return u.Scheme + "://" + u.Host
}

func EncodeURL(base string, values map[string]string) string {
	p := url.Values{}
	for k, v := range values {
		p.Add(k, v)
	}

	if len(values) > 0 {
		return base + "?" + p.Encode()
	}

	return base
}
