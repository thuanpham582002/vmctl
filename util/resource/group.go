package resource

import (
	"fmt"
	"vmctl/model"
)

func RunGroupAction(action func(group model.Group) error, groupName string) error {
	vmManager, err := GetVmManager()
	if err != nil {
		return err
	}
	group, ok := vmManager[groupName]
	if !ok {
		return fmt.Errorf("group %s not found", groupName) // Thay đổi để trả về lỗi cụ thể
	}
	return action(group)
}

func RunVMAction(action func(vm model.VirtualMachineYaml) error, groupName, vmName string) error {
	vmManager, err := GetVmManager()
	if err != nil {
		return err
	}
	group, ok := vmManager[groupName]
	if !ok {
		return fmt.Errorf("group %s not found", groupName) // Thay đổi để trả về lỗi cụ thể
	}
	vm, ok := group[vmName]
	if !ok {
		return fmt.Errorf("virtual machine %s not found in group %s", vmName, groupName) // Thay đổi để trả về lỗi cụ thể
	}
	_ = vm.Validate()
	return action(vm)
}
