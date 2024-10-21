package config

import (
	"github.com/bitfield/script"
	"github.com/spf13/cobra"
	"vmctl/util/config"
)

func NewCmdConfigView() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-context",
		Short: "Display the current-context",
		Run: func(cmd *cobra.Command, args []string) {
			contextPath, err := config.GetContextPath()
			if err != nil {
				panic(err)
			}
			script.NewPipe().Echo(contextPath)
		},
	}
	return cmd
}
