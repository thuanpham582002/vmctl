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
		for key, _ := range vm.InitScript.FromOldest() {
			if strings.HasPrefix(key, toComplete) {
				scripts = append(scripts, key)
			}
		}
	})
	scripts = unique(scripts)
	return scripts, cobra.ShellCompDirectiveDefault
}

func unique(scripts []string) []string {
	uniqueScripts := make([]string, 0, len(scripts)) // Slice để chứa kết quả
	seen := make(map[string]bool)                    // Map để theo dõi các phần tử đã gặp

	for _, script := range scripts {
		if !seen[script] { // Nếu script chưa xuất hiện trong map
			uniqueScripts = append(uniqueScripts, script)
			seen[script] = true // Đánh dấu script là đã gặp
		}
	}

	return uniqueScripts
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
