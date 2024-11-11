package delete

import (
	"fmt"
	"github.com/spf13/cobra"
	"vmctl/common"
	"vmctl/model"
	"vmctl/util/completion"
	"vmctl/util/printcolor"
	"vmctl/util/resource"
)

func NewCmdDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete <node_path...>",
		Short:             "Delete a new virtual machine",
		ValidArgsFunction: completion.BashCompleteInstance,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			} else {
				resource.NewBuilder().
					SetNodePaths(args).
					Do(DeleteVM)
			}
		},
		Aliases: []string{"del"},
	}
	return cmd
}

func DeleteVM(vm model.VirtualMachine) {
	printcolor.Info(fmt.Sprintf("Deleting VM %s in group %s", vm.Name, vm.Group))
	if _, _, err := common.ExecShell("limactl", "stop", vm.Name); err != nil {
		printcolor.Error(fmt.Sprintf("Error deleting VM %s in group %s: %v", vm.Name, vm.Group, err))
	}
	if _, _, err := common.ExecShell("limactl", "delete", vm.Name); err != nil {
		printcolor.Error(fmt.Sprintf("Error deleting VM %s in group %s: %v", vm.Name, vm.Group, err))
		return
	}
	printcolor.Success(fmt.Sprintf("VM %s in group %s deleted successfully", vm.Name, vm.Group))
}
