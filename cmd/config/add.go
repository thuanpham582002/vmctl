package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"vmctl/util/printcolor"
	"vmctl/util/validate"
)

func NewCmdConfigAddContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [PATH]",
		Short: "Add context",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			contexts := viper.GetStringSlice("contexts")
			for _, context := range args {
				if !validate.IsValidConfigFile(context) {
					printcolor.Error("Context file is invalid" + context)
					continue
				}
				if len(contexts) == 0 {
					viper.Set("current-context", context)
					printcolor.Success("Context added successfully" + context)
				}
				contexts = append(contexts, context)
				viper.Set("contexts", contexts)
			}
			err := viper.WriteConfig()
			if err != nil {
				printcolor.Error("Error adding context" + err.Error())
				return
			}
			printcolor.Success("Context added successfully")
		},
	}
	return cmd
}
