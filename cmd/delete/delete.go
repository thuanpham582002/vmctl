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

type DeleteOptions struct {
}

func NewDeleteOptions() *DeleteOptions {
	return &DeleteOptions{}
}

func NewCmdDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "delete <node_path...>",
		Short:             "Delete a new virtual machine",
		ValidArgsFunction: completion.BashCompleteInstance,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			} else {
				resource.
					NewBuilder().
					SetNodePaths(args).
					Do(func(vm model.VirtualMachine) {
						printcolor.Info(fmt.Sprintf("Deleting VM %s in group %s", vm.Name, vm.Group))
						if _, _, err := common.ExecShell("limactl", fmt.Sprintf("stop %s", vm.Name)); err != nil {
							printcolor.Error(fmt.Sprintf("Error deleting VM %s in group %s: %v", vm.Name, vm.Group, err))
							return
						}
						if _, _, err := common.ExecShell("limactl", fmt.Sprintf("delete %s", vm.Name)); err != nil {
							printcolor.Error(fmt.Sprintf("Error deleting VM %s in group %s: %v", vm.Name, vm.Group, err))
							return
						}
					})
			}
		},
	}
	return cmd
}
