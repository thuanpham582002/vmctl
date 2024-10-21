package start

import (
	"fmt"
	"github.com/spf13/cobra"
	"vmctl/common"
	"vmctl/model"
	"vmctl/util/completion"
	"vmctl/util/printcolor"
	"vmctl/util/resource"
)

type StartOptions struct {
}

func NewStartOptions() *StartOptions {
	return &StartOptions{}
}

func NewCmdStart() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "start <node_path...>",
		Short:             "Start a new virtual machine",
		ValidArgsFunction: completion.BashCompleteInstance,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			} else {
				resource.
					NewBuilder().
					SetNodePaths(args).
					Do(func(vm model.VirtualMachine) {
						printcolor.Info(fmt.Sprintf("Starting VM %s in group %s", vm.Name, vm.Group))
						if _, _, err := common.ExecShell("limactl", fmt.Sprintf("start %s", vm.Name)); err != nil {
							printcolor.Error(fmt.Sprintf("Error starting VM %s in group %s: %v", vm.Name, vm.Group, err))
							return
						}
					})
			}
		},
	}
	return cmd
}
