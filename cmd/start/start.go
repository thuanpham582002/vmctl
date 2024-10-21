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

func NewCmdStart() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "start <node_path...>",
		Short:             "Start a new virtual machine",
		ValidArgsFunction: completion.BashCompleteInstance,
		Run:               RunStart,
	}
	return cmd
}

func RunStart(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
	} else {
		resource.NewBuilder().
			SetNodePaths(args).
			Do(StartVM)
	}
}

func StartVM(vm model.VirtualMachine) {
	printcolor.Info(fmt.Sprintf("Starting VM %s in group %s", vm.Name, vm.Group))
	if _, _, err := common.ExecShell("limactl", "start", vm.Name); err != nil {
		printcolor.Error(fmt.Sprintf("Error starting VM %s in group %s: %v", vm.Name, vm.Group, err))
		return
	}
	failedScripts := executeInitScripts(vm)
	printcolor.Success(fmt.Sprintf("VM %s in group %s started successfully", vm.Name, vm.Group))
	if len(failedScripts) > 0 {
		printcolor.Warning(fmt.Sprintf("Failed to execute script(s) %v", failedScripts))
	}
}

func executeInitScripts(vm model.VirtualMachine) []string {
	var failedScripts []string
	for key, script := range vm.InitScript.FromOldest() {
		if !script.OnBoot {
			continue
		}
		scriptCommand, err := script.GetCommand()
		if err != nil {
			printcolor.Warning(fmt.Sprintf("Error building script for %s in %s: %v", key, vm.Name, err))
			continue
		}
		shellArgs := buildShellArgs(vm, script.Root, scriptCommand)
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
