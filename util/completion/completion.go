package completion

import (
	"github.com/spf13/cobra"
	"strings"
	"vmctl/model"
	"vmctl/util/resource"
)

func BashCompleteInstance(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	vmManager, _ := resource.GetVmManager()
	nodePaths := generateNodePaths(vmManager)
	return nodePaths, cobra.ShellCompDirectiveNoFileComp
}

func BashCompleteForCommands(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var scripts []string
	resource.NewBuilder().SetNodePaths(args).Do(func(vm model.VirtualMachine) {
		for key, _ := range vm.InitScript {
			if strings.HasPrefix(key, toComplete) {
				scripts = append(scripts, key)
			}
		}
	})

	return scripts, cobra.ShellCompDirectiveDefault
}

func generateNodePaths(vmManager model.VMManager) []string {
	nodePaths := make([]string, 0)
	nodePaths = append(nodePaths, ".")
	for group, _ := range vmManager {
		nodePaths = append(nodePaths, "."+group)
		for vm, _ := range vmManager[group] {
			nodePaths = append(nodePaths, "."+group+"."+vm)
		}
	}
	return nodePaths
}
