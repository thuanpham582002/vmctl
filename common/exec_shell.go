package common

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"vmctl/model"
	"vmctl/util/limashell"
	"vmctl/util/printcolor"
)

func ExecShell(name string, command ...string) (string, int, error) {
	cmd := exec.Command(name, command...)
	printcolor.Info(fmt.Sprintf("Executing command: \n%s", fmt.Sprintf("%s", cmd.String())))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	// Kiểm tra exit code
	if err != nil {
		printcolor.Error(err.Error())
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			// Lấy exit code từ ExitError
			exitCode := exitError.ExitCode()
			return cmd.String(), exitCode, err
		}
	}
	return cmd.String(), 0, nil
}

func ExecShellWithOutput(name string, command ...string) (string, int, error) {
	cmd := exec.Command(name, command...)
	printcolor.Info(fmt.Sprintf("Executing command: \n%s", fmt.Sprintf("%s", cmd.String())))
	output, err := cmd.CombinedOutput()
	// Kiểm tra exit code
	if err != nil {
		printcolor.Error(err.Error())
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			// Lấy exit code từ ExitError
			exitCode := exitError.ExitCode()
			return string(output), exitCode, err
		}
	}
	return string(output), 0, nil
}

func GetNetplanConfig(vm model.VirtualMachine) (string, error) {
	shellArgs := limashell.BuildShellArgs(vm, true, "cat /etc/netplan/50-cloud-init.yaml")
	output, _, err := ExecShellWithOutput("limactl", shellArgs...)
	if err != nil {
		return "", err
	}
	return output, nil
}

func ExecuteStaticIPScript(vm model.VirtualMachine) {
	if vm.StaticIP == "" {
		return
	}
	printcolor.Info(fmt.Sprintf("Setting static IP for %s in %s", vm.Name, vm.Group))
	netplanConfig, err := GetNetplanConfig(vm)
	if err != nil {
		printcolor.Error(fmt.Sprintf("Error getting netplan config for %s in %s: %v", vm.Name, vm.Group, err))
		return
	}
	var netplanYaml map[string]interface{}
	if err := yaml.Unmarshal([]byte(netplanConfig), &netplanYaml); err != nil {
		printcolor.Error(fmt.Sprintf("Error parsing netplan config for %s in %s: %v", vm.Name, vm.Group, err))
		return
	}

	network := netplanYaml["network"].(map[string]interface{})
	ethernets := network["ethernets"].(map[string]interface{})
	lima0 := ethernets["lima0"].(map[string]interface{})

	// Remove match and set-name
	delete(lima0, "match")
	delete(lima0, "set-name")

	// Set static IP configuration
	lima0["dhcp4"] = false
	lima0["addresses"] = []string{fmt.Sprintf("%s/24", vm.StaticIP)}
	lima0["routes"] = []map[string]interface{}{
		{
			"to":  "0.0.0.0/0",
			"via": "192.168.105.1",
		},
	}

	// Set MTU to 1280 to avoid docker errors when pulling images
	lima0["mtu"] = 1280

	// Write back to yaml
	newConfig, err := yaml.Marshal(netplanYaml)
	if err != nil {
		printcolor.Error(fmt.Sprintf("Error marshaling netplan config for %s in %s: %v", vm.Name, vm.Group, err))
		return
	}

	// Write the config back to file
	writeConfigCmd := fmt.Sprintf("echo '%s' > /etc/netplan/50-cloud-init.yaml", string(newConfig))
	shellArgs := limashell.BuildShellArgs(vm, true, writeConfigCmd)
	if _, _, err := ExecShell("limactl", shellArgs...); err != nil {
		printcolor.Error(fmt.Sprintf("Error writing netplan config for %s in %s: %v", vm.Name, vm.Group, err))
		return
	}
	// Apply the netplan config
	applyCmd := "netplan apply"
	shellArgs = limashell.BuildShellArgs(vm, true, applyCmd)
	if _, _, err := ExecShell("limactl", shellArgs...); err != nil {
		printcolor.Error(fmt.Sprintf("Error applying netplan config for %s in %s: %v", vm.Name, vm.Group, err))
		return
	}
	printcolor.Success(fmt.Sprintf("Static IP set for %s in %s", vm.Name, vm.Group))
}
