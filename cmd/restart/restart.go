package restart

import (
	"fmt"
	"github.com/spf13/cobra"
	"vmctl/common"
	"vmctl/model"
	"vmctl/util/completion"
	"vmctl/util/printcolor"
	"vmctl/util/resource"
)

type RestartOptions struct {
}

func NewRestartOptions() *RestartOptions {
	return &RestartOptions{}
}

func NewCmdRestart() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "restart <node_path...>",
		Short:             "Restart a new virtual machine",
		ValidArgsFunction: completion.BashCompleteInstance,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			} else {
				resource.
					NewBuilder().
					SetNodePaths(args).
					Do(func(vm model.VirtualMachine) {
						printcolor.Info(fmt.Sprintf("Restarting VM %s in group %s", vm.Name, vm.Group))
						if _, _, err := common.ExecShell("limactl", fmt.Sprintf("stop %s", vm.Name)); err != nil {
							printcolor.Error(fmt.Sprintf("Error restarting VM %s in group %s: %v", vm.Name, vm.Group, err))
							return
						}
						if _, _, err := common.ExecShell("limactl", fmt.Sprintf("start %s", vm.Name)); err != nil {
							printcolor.Error(fmt.Sprintf("Error restarting VM %s in group %s: %v", vm.Name, vm.Group, err))
							return
						}
					})
			}
		},
	}
	return cmd
}
