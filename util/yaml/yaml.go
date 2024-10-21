package yaml

import (
	"fmt"
	"strings"
)

func GetValueFromPath(data map[string]interface{}, path string) (interface{}, error) {
	keys := strings.Split(path, ".")
	value := data

	for _, key := range keys {
		if v, ok := value[key].(map[string]interface{}); ok {
			value = v
		} else if v, ok := value[key]; ok {
			return v, nil
		} else {
			return nil, fmt.Errorf("path %s not found", path)
		}
	}

	return value, nil
}
