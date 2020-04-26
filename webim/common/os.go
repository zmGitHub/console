package common

import "strings"

func IsPC(os string) bool {
	return strings.HasPrefix(os, "Mac") || strings.HasPrefix(os, "Windows") ||
		strings.HasPrefix(os, "Linux")
}

func IsMobile(os string) bool {
	osCategories := []string{
		"iPhone",
		"Android",
	}

	for _, osCategory := range osCategories {
		if strings.HasPrefix(os, osCategory) {
			return true
		}
	}

	return false
}
