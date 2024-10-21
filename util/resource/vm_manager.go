package resource

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"vmctl/model"
)

func GetVmManager() (model.VMManager, error) {
	bytes, err := os.ReadFile(fmt.Sprint(viper.Get("current-context")))
	if err != nil {
		return nil, err
	}
	var vmManager model.VMManager
	err = yaml.Unmarshal(bytes, &vmManager)
	if err != nil {
		return nil, err
	}
	return vmManager, nil
}
