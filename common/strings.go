package common

import "strings"

func GetStringBetween(str string, start string, end string) string {
	startIndex := strings.Index(str, start)
	if startIndex == -1 {
		return ""
	}
	startIndex += len(start)
	endIndex := strings.Index(str[startIndex:], end)
	if endIndex == -1 {
		return ""
	}
	return str[startIndex : startIndex+endIndex]
}
