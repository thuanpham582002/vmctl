package model

import (
	"vmctl/pkg/orderedmap"
	"vmctl/util/config"
)

type VirtualMachineYaml struct {
	Template   string                                `yaml:"template"`
	InitScript orderedmap.OrderedMap[string, Script] `yaml:"init_script"`
}

type VirtualMachine struct {
	Template   string
	InitScript orderedmap.OrderedMap[string, Script]
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
