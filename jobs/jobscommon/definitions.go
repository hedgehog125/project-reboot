package jobscommon

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hedgehog125/project-reboot/common"
)

func GetVersionedType(id string, version int) string {
	return fmt.Sprintf("%v_%v", id, version)
}
func ParseVersionedType(versionedType string) (string, int, *common.Error) {
	separatorIndex := strings.LastIndex(versionedType, "_")
	if separatorIndex == -1 {
		return "", 0, ErrMalformedVersionedType.AddCategory(ErrTypeParseVersionedType)
	}
	version, err := strconv.Atoi(versionedType[separatorIndex+1:])
	if err != nil {
		return "", 0, ErrMalformedVersionedType.AddCategory(ErrTypeParseVersionedType)
	}

	return versionedType[:separatorIndex], version, nil
}
