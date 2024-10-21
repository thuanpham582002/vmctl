package stop

import (
	"fmt"
	"github.com/spf13/cobra"
	"vmctl/common"
	"vmctl/model"
	"vmctl/util/completion"
	"vmctl/util/printcolor"
	"vmctl/util/resource"
)

type StopOptions struct {
}

func NewStopOptions() *StopOptions {
	return &StopOptions{}
}

func NewCmdStop() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "stop <node_path...>",
		Short:             "Stop a new virtual machine",
		ValidArgsFunction: completion.BashCompleteInstance,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			} else {
				resource.
					NewBuilder().
					SetNodePaths(args).
					Do(func(vm model.VirtualMachine) {
						printcolor.Info(fmt.Sprintf("Stopping VM %s in group %s", vm.Name, vm.Group))
						if _, _, err := common.ExecShell("limactl", fmt.Sprintf("stop %s", vm.Name)); err != nil {
							printcolor.Error(fmt.Sprintf("Error stopping VM %s in group %s: %v", vm.Name, vm.Group, err))
							return
						}
					})
			}
		},
	}
	return cmd
}
