package restart

import (
	"fmt"
	"github.com/spf13/cobra"
	"vmctl/cmd/start"
	"vmctl/common"
	"vmctl/model"
	"vmctl/util/completion"
	"vmctl/util/printcolor"
	"vmctl/util/resource"
)

func NewCmdRestart() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "restart <node_path...>",
		Short:             "Restart a new virtual machine",
		ValidArgsFunction: completion.BashCompleteInstance,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			} else {
				resource.NewBuilder().
					SetNodePaths(args).
					Do(restartVM)
			}
		},
	}
	return cmd
}

func restartVM(vm model.VirtualMachine) {
	printcolor.Info(fmt.Sprintf("Restarting VM %s in group %s", vm.Name, vm.Group))
	if _, _, err := common.ExecShell("limactl", "stop", vm.Name); err != nil {
		printcolor.Error(fmt.Sprintf("Error stopping VM %s in group %s: %v", vm.Name, vm.Group, err))
	}
	start.StartVM(vm)
	printcolor.Success(fmt.Sprintf("VM %s in group %s restarted successfully", vm.Name, vm.Group))
}
