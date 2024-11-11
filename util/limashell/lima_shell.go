package limashell

import "vmctl/model"

func BuildShellArgs(vm model.VirtualMachine, root bool, args ...string) []string {
	shellArgs := []string{"shell", vm.Name}
	if root {
		shellArgs = append(shellArgs, "sudo", "bash", "-c")
	} else {
		shellArgs = append(shellArgs, "bash", "-c")
	}
	return append(shellArgs, args...)
}
