package model

import (
	"os"
	"vmctl/util/config"
)

type VirtualMachineYaml struct {
	Template   string            `yaml:"template"`
	InitScript map[string]Script `yaml:"init_script"`
}

type VirtualMachine struct {
	Template   string
	InitScript map[string]Script
	Group      string
	Name       string
}

func (vm *VirtualMachineYaml) ToVirtualMachine(groupName, vmName string) VirtualMachine {
	return VirtualMachine{
		Template:   vm.Template,
		InitScript: vm.InitScript,
		Group:      groupName,
		Name:       vmName,
	}
}

func (vm *VirtualMachineYaml) Validate() error {
	// check file template is abstract path or hard path
	if _, err := os.Stat(vm.Template); err != nil {
		// do nothing
	} else {
		dir, _ := config.GetContextDir()
		// check file template is exist
		if _, err := os.Stat(dir + "/" + vm.Template); err != nil {
			vm.Template = dir + "/" + vm.Template
		}
	}
	return nil
}
