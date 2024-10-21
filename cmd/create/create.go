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
	return &CreateOptions{}
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
				resource.NewBuilder().
					SetNodePaths(args).
					Do(createVM(o))
			}
		},
	}
	cmd.Flags().BoolVarP(&o.Root, "root", "r", false, "Create all VMs in the group")
	return cmd
}

func createVM(o *CreateOptions) func(vm model.VirtualMachine) {
	return func(vm model.VirtualMachine) {
		printcolor.Info(fmt.Sprintf("Creating VM %s in group %s", vm.Name, vm.Group))
		if _, _, err := common.ExecShell("limactl", "start", "--name="+vm.Name, vm.Template, "--tty=false"); err != nil {
			printcolor.Error(fmt.Sprintf("Error creating VM %s in group %s: %v", vm.Name, vm.Group, err))
			return
		}
		failedScripts := executeInitScripts(vm, o.Root)
		printcolor.Success(fmt.Sprintf("VM %s in group %s created successfully", vm.Name, vm.Group))
		if len(failedScripts) > 0 {
			printcolor.Warning(fmt.Sprintf("Failed to execute script(s) %v", failedScripts))
		}
	}
}

func executeInitScripts(vm model.VirtualMachine, root bool) []string {
	var failedScripts []string
	for key, script := range vm.InitScript {
		scriptStr, err := script.GetCommand()
		if err != nil {
			printcolor.Error(fmt.Sprintf("Error getting script for %s in %s: %v", key, vm.Name, err))
			failedScripts = append(failedScripts, key)
			continue
		}
		shellArgs := buildShellArgs(vm, root || script.Root, "bash", "-c", scriptStr)
		if _, _, err := common.ExecShell("limactl", shellArgs...); err != nil {
			printcolor.Error(fmt.Sprintf("Error executing script for %s in %s: %v", key, vm.Name, err))
			failedScripts = append(failedScripts, key)
		}
	}
	return failedScripts
}

func buildShellArgs(vm model.VirtualMachine, root bool, args ...string) []string {
	shellArgs := []string{"shell", vm.Name}
	if root {
		shellArgs = append(shellArgs, "sudo")
	}
	return append(shellArgs, args...)
}
