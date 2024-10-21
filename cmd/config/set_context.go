package config

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"vmctl/util/config"
	"vmctl/util/printcolor"
)

func NewCmdConfigSetContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-context [PATH]",
		Short: "Set the current-context",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			homeDir, err := os.UserHomeDir()

			info, err := os.Stat(path)
			if err != nil {
				if os.IsNotExist(err) {
					printcolor.Print("No such file or directory")
				} else {
					printcolor.Print("Error: " + err.Error())
				}
				return
			}

			if info.IsDir() {
				printcolor.Print("This is a directory, please input file")
				return
			}

			configPath := filepath.Join(homeDir, ".vmctl", "config.yaml")
			err = SetConfig(configPath, &config.Config{CurrentContext: path})
			if err != nil {
				panic(err)
				return
			}
			printcolor.Print("Current context set to " + path)
		},
	}
	return cmd
}

// SetConfig ghi cấu hình vào file
func SetConfig(path string, config *config.Config) error {
	// Mở file để ghi, tạo mới nếu không tồn tại
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mã hóa cấu hình thành YAML
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}
