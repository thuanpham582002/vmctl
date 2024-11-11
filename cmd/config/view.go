package config

import (
	"github.com/spf13/cobra"
	"vmctl/util/config"
	"vmctl/util/printcolor"
)

func NewCmdConfigView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "Display the current-context",
		Run: func(cmd *cobra.Command, args []string) {
			contextPath, err := config.GetListContext()
			if err != nil {
				printcolor.Error(err.Error())
			}
			for _, context := range contextPath {
				printcolor.Print(context)
			}
		},
	}
	return cmd
}
