package execute

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"vmctl/common"
	"vmctl/model"
	"vmctl/util/completion"
	"vmctl/util/printcolor"
	"vmctl/util/resource"
)

type ExecuteOptions struct {
	Root     bool
	commands []string
	files    []string
}

func NewExecuteOptions() *ExecuteOptions {
	return &ExecuteOptions{
		Root: false,
	}
}

func NewCmdExecute() *cobra.Command {
	o := NewExecuteOptions()
	cmd := &cobra.Command{
		Use:               "exec <node_path...>",
		Short:             "Execute a new virtual machine",
		ValidArgsFunction: completion.BashCompleteInstance,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			} else {
				resource.
					NewBuilder().
					SetNodePaths(args).
					Do(func(vm model.VirtualMachine) {
						printcolor.Info(fmt.Sprintf("Executing VM %s in group %s", vm.Name, vm.Group))
						// Tạo chuỗi lệnh để thực thi
						shellCommand := "bash"
						if o.Root {
							shellCommand = "sudo bash"
						}

						for _, command := range o.commands {
							// init_script key containe command

							if _, ok := vm.InitScript[command]; ok {
								scriptStr, err := vm.InitScript[command].GetCommand()
								if err != nil {
									printcolor.Warning(fmt.Sprintf("Error building script for %s in %s: %v", command, vm.Name, err))
									continue
								}
								if _, _, err := common.ExecShell("limactl", "shell", vm.Name, shellCommand, "-c", "\""+scriptStr+"\""); err != nil {
									printcolor.Error(fmt.Sprintf("Error executing VM %s in group %s: %v", vm.Name, vm.Group, err))
									continue
								}
							}
							// command is a string
							if _, _, err := common.ExecShell("limactl", "shell", vm.Name, shellCommand, "-c", "\""+command+"\""); err != nil {
								printcolor.Error(fmt.Sprintf("Error executing VM %s in group %s: %v", vm.Name, vm.Group, err))
								continue
							}
						}
						for _, file := range o.files {
							if fileData, err := os.ReadFile(file); err == nil {
								fileDataStr := string(fileData)
								if _, _, err := common.ExecShell("limactl", "shell", vm.Name, shellCommand, "-c", "\""+fileDataStr+"\""); err != nil {
									printcolor.Error(fmt.Sprintf("Error executing VM %s in group %s: %v", vm.Name, vm.Group, err))
									continue
								}
							}
						}
					})
			}
		},
	}
	flagSet := cmd.Flags()
	flagSet.BoolVarP(&o.Root, "root", "r", false, "Execute all VMs in the group")
	cmd.Flags().StringSliceVarP(&o.commands, "command", "c", []string{}, "Execute a list of commands")
	cmd.Flags().StringSliceVarP(&o.files, "file", "f", []string{}, "Execute a list of files")
	err := cmd.RegisterFlagCompletionFunc("command", completion.BashCompleteForCommands)
	if err != nil {
		printcolor.Error(fmt.Sprintf("Error registering flag completion function: %v", err))
	}

	annotations := []string{}

	for _, ext := range resource.FileExtensions {
		annotations = append(annotations, strings.TrimLeft(ext, "."))
	}
	_ = cmd.Flags().SetAnnotation("file", cobra.BashCompFilenameExt, annotations)
	return cmd
}
