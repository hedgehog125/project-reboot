package jobscommon

import "fmt"

func JoinPaths(path1 string, path2 string) string {
	if path1 == "" {
		return path2
	}
	if path2 == "" {
		return path1
	}
	return fmt.Sprintf("%s/%s", path1, path2)
}
