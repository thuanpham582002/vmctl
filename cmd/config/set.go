package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"vmctl/util/config"
	"vmctl/util/printcolor"
)

func NewCmdConfigSetContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "set [PATH]",
		Short:             "Set the current-context",
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: ValidConfigFile,
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]

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

			err = SetConfig(path)
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
func SetConfig(config string) error {
	viper.Set("current-context", config)
	return viper.WriteConfig()
}

func ValidConfigFile(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	contextPath, err := config.GetListContext()
	if err != nil {
		printcolor.Error(err.Error())
	}
	return contextPath, cobra.ShellCompDirectiveNoFileComp
}
