package resource

import (
	"sort"
	"strings"
	"vmctl/model"
	"vmctl/util/printcolor"
)

var FileExtensions = []string{".sh", ".bash", ".zsh"}

type Builder struct {
	virtualMachineMap map[string]model.VirtualMachine
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) SetNodePaths(nodePaths []string) *Builder {
	vmManager, _ := GetVmManager()
	b.virtualMachineMap = b.FromNodePathsToVMs(RemoveChildNodePaths(b.Unique(nodePaths)), vmManager)
	return b
}

func (b *Builder) Unique(paths []string) []string {
	// Tạo một map để lưu trữ các giá trị duy nhất
	uniqueMap := make(map[string]bool)
	// Duyệt qua mảng và thêm vào map
	for _, path := range paths {
		uniqueMap[path] = true
	}
	// Tạo một mảng mới để lưu trữ các giá trị duy nhất
	uniquePaths := make([]string, 0)
	// Duyệt qua map và thêm vào mảng
	for path := range uniqueMap {
		uniquePaths = append(uniquePaths, path)
	}
	return uniquePaths
}

func (b *Builder) RemoveChildNodePaths(uniquePaths []string) []string {
	sort.Strings(uniquePaths)

	var result []string
	for i, path := range uniquePaths {
		isChild := false
		for j := 0; j < i; j++ {
			if strings.HasPrefix(path, uniquePaths[j]+".") {
				isChild = true
				break
			}
		}
		if !isChild {
			result = append(result, path)
		}
	}
	return result
}

func (b *Builder) FromNodePathsToVMs(nodePaths []string, vmManager model.VMManager) map[string]model.VirtualMachine {
	virtualMachineMap := make(map[string]model.VirtualMachine)
	for _, nodePath := range nodePaths {
		arr := strings.Split(strings.Trim(nodePath, "."), ".")
		if len(arr) == 1 { // Only has group
			if arr[0] == "" {
				for groupName := range vmManager {
					for vmName := range vmManager[groupName] {
						virtualMachineYaml := vmManager[groupName][vmName]
						virtualMachineMap[vmName] = virtualMachineYaml.ToVirtualMachine(groupName, vmName)
					}
				}
			} else {
				for vmName := range vmManager[arr[0]] {
					virtualMachineYaml := vmManager[arr[0]][vmName]
					virtualMachineMap[vmName] = virtualMachineYaml.ToVirtualMachine(arr[0], vmName)
				}
			}
		} else if len(arr) == 2 { // Has group and vm
			virtualMachineYaml := vmManager[arr[0]][arr[1]]
			virtualMachineMap[arr[1]] = virtualMachineYaml.ToVirtualMachine(arr[0], arr[1])
		}
	}
	return virtualMachineMap
}

func (b *Builder) Do(actions func(machine model.VirtualMachine)) {
	if len(b.virtualMachineMap) == 0 {
		printcolor.Error("No virtual machine found")
		return
	}
	for _, vm := range b.virtualMachineMap {
		actions(vm)
	}
}
