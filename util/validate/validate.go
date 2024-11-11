package validate

import "os"

func IsValidConfigFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
