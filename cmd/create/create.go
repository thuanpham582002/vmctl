package create

import (
	"fmt"
	"github.com/spf13/cobra"
	"vmctl/common"
	"vmctl/model"
	"vmctl/util/completion"
	"vmctl/util/printcolor"
	"vmctl/util/resource"
)

type CreateOptions struct {
	Root bool
}

func NewCreateOptions() *CreateOptions {
	return &CreateOptions{
		Root: false,
	}
}

func NewCmdCreate() *cobra.Command {
	o := NewCreateOptions()
	cmd := &cobra.Command{
		Use:               "create <node_path...>",
		Short:             "Create a new virtual machine",
		ValidArgsFunction: completion.BashCompleteInstance,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			} else {
				resource.
					NewBuilder().
					SetNodePaths(args).
					Do(func(vm model.VirtualMachine) {
						printcolor.Info(fmt.Sprintf("Creating VM %s in group %s", vm.Name, vm.Group))
						if _, _, err := common.ExecShell("limactl", "start", "--name="+vm.Name, vm.Template, "--tty=false"); err != nil {
							printcolor.Error(fmt.Sprintf("Error creating VM %s in group %s: %v", vm.Name, vm.Group, err))
							return
						}
						for key, script := range vm.InitScript {
							scriptCommand, err := script.GetCommand()
							if err != nil {
								printcolor.Warning(fmt.Sprintf("Error building script for %s in %s: %v", key, vm.Name, err))
								continue
							}

							// Tạo chuỗi lệnh để thực thi
							shellCommand := "bash"
							if o.Root {
								shellCommand = "sudo bash"
							}
							if _, _, err := common.ExecShell("limactl", "shell", vm.Name, shellCommand, "-c", scriptCommand); err != nil {
								printcolor.Error(fmt.Sprintf("Error executing script for %s in %s: %v", key, vm.Name, err))
								return
							}

							printcolor.Success(fmt.Sprintf("VM %s in group %s created successfully", vm.Name, vm.Group))
						}
					})
			}
		},
	}
	cmd.Flags().BoolVarP(&o.Root, "root", "r", false, "Create all VMs in the group")
	return cmd
}
