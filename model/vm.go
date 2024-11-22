package model

import (
	"vmctl/pkg/orderedmap"
	"vmctl/util/config"
)

type VirtualMachineYaml struct {
	StaticIP   string                                `yaml:"staticIp"`
	Template   string                                `yaml:"template"`
	InitScript orderedmap.OrderedMap[string, Script] `yaml:"initScript"`
}

type VirtualMachine struct {
	StaticIP   string
	Template   string
	InitScript orderedmap.OrderedMap[string, Script]
	Group      string
	Name       string
}

func (vm *VirtualMachineYaml) ToVirtualMachine(groupName, vmName string) VirtualMachine {
	dir, _ := config.GetContextDir()
	return VirtualMachine{
		StaticIP:   vm.StaticIP,
		Template:   dir + "/" + vm.Template,
		InitScript: vm.InitScript,
		Group:      groupName,
		Name:       vmName,
	}
}
