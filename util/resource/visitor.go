package resource

import (
	"sort"
	"strings"
)

type Info struct {
	Name  string
	Group string
}

func RemoveChildNodePaths(uniquePaths []string) []string {
	sort.Strings(uniquePaths)

	var result []string
	for i, path := range uniquePaths {
		isChild := false
		for j := 0; j < i; j++ {
			if strings.HasPrefix(path, uniquePaths[j]+".") {
				isChild = true
				break
			}
		}
		if !isChild {
			result = append(result, path)
		}
	}
	return result
}
