package config

import (
	"github.com/spf13/cobra"
)

func NewCmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Modify kubeconfig files",
		Long:  "Modify kubeconfig files using subcommands like set, view, use-context",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(NewCmdConfigGetContext())
	cmd.AddCommand(NewCmdConfigSetContext())
	cmd.AddCommand(NewCmdConfigAddContext())
	cmd.AddCommand(NewCmdConfigView())
	return cmd
}
