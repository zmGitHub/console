package common

import "strings"

func GenPlaceHolders(count int) string {
	s := strings.Repeat("?,", count)
	return strings.TrimRight(s, ",")
}
