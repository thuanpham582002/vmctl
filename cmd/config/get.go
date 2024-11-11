package config

import (
	"github.com/spf13/cobra"
	"vmctl/util/config"
	"vmctl/util/printcolor"
)

func NewCmdConfigGetContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display the current-context",
		Run: func(cmd *cobra.Command, args []string) {
			contextPath, err := config.GetContextPath()
			if err != nil {
				panic(err)
			}
			printcolor.Print(contextPath)
		},
	}
	return cmd
}
