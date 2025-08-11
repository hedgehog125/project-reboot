package common

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ErrTypeParseVersionedType = "parse versioned type"
)

var ErrMalformedVersionedType = NewErrorWithCategories(
	"malformed versioned type",
)
var ErrWrapperParseVersionedType = NewErrorWrapper(ErrTypeParseVersionedType)

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

func GetVersionedType(id string, version int) string {
	return fmt.Sprintf("%v_%v", id, version)
}
func ParseVersionedType(versionedType string) (string, int, *Error) {
	separatorIndex := strings.LastIndex(versionedType, "_")
	if separatorIndex == -1 {
		return "", 0, ErrWrapperParseVersionedType.Wrap(ErrMalformedVersionedType)
	}
	version, err := strconv.Atoi(versionedType[separatorIndex+1:])
	if err != nil {
		return "", 0, ErrWrapperParseVersionedType.Wrap(ErrMalformedVersionedType)
	}

	return versionedType[:separatorIndex], version, nil
}
