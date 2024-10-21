package model

import (
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
	dir, _ := config.GetContextDir()
	return VirtualMachine{
		Template:   dir + "/" + vm.Template,
		InitScript: vm.InitScript,
		Group:      groupName,
		Name:       vmName,
	}
}
