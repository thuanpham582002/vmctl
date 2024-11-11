package execute

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"vmctl/common"
	"vmctl/model"
	"vmctl/util/completion"
	"vmctl/util/limashell"
	"vmctl/util/printcolor"
	"vmctl/util/resource"
)

type ExecuteOptions struct {
	Root         bool
	Commands     []string
	Files        []string
	ContinueFrom string
}

func NewExecuteOptions() *ExecuteOptions {
	return &ExecuteOptions{}
}

func NewCmdExecute() *cobra.Command {
	o := NewExecuteOptions()
	cmd := &cobra.Command{
		Use:               "execute <node_path...>",
		Short:             "Execute a new virtual machine",
		ValidArgsFunction: completion.BashCompleteInstance,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			} else {
				resource.NewBuilder().
					SetNodePaths(args).
					Do(func(vm model.VirtualMachine) {
						printcolor.Info(fmt.Sprintf("Executing VM %s in group %s", vm.Name, vm.Group))
						executeCommands(vm, o)
						executeFiles(vm, o)
						if o.ContinueFrom != "" {
							isHashFound := false
							for key, _ := range vm.InitScript.FromOldest() {
								if key == o.ContinueFrom {
									isHashFound = true
								}
								if isHashFound {
									o.Commands = append(o.Commands, key)
								}
							}
							if isHashFound {
								printcolor.Print(fmt.Sprintf("Continue from %s", o.ContinueFrom))
								executeCommands(vm, o)
							} else {
								printcolor.Warning(fmt.Sprintf("Continue from %s not found", o.ContinueFrom))
							}
						}
					})
			}
		},
		Aliases: []string{"exec"},
	}
	flagSet := cmd.Flags()
	flagSet.BoolVarP(&o.Root, "root", "r", false, "Execute all VMs in the group")
	cmd.Flags().StringSliceVarP(&o.Commands, "command", "c", []string{}, "Execute a list of commands")
	cmd.Flags().StringSliceVarP(&o.Files, "file", "f", []string{}, "Execute a list of files")
	cmd.Flags().StringVarP(&o.ContinueFrom, "continue-from", "", "", "Continue from a specific command")
	err := cmd.RegisterFlagCompletionFunc("command", completion.BashCompleteForCommands)
	if err != nil {
		printcolor.Error(fmt.Sprintf("Error registering flag completion function: %v", err))
	}
	err = cmd.RegisterFlagCompletionFunc("continue-from", completion.BashCompleteForCommands)
	if err != nil {
		printcolor.Error(fmt.Sprintf("Error registering flag completion function: %v", err))
	}
	_ = cmd.Flags().SetAnnotation("file", cobra.BashCompFilenameExt, resource.FileExtensions)
	return cmd
}

func executeCommands(vm model.VirtualMachine, o *ExecuteOptions) {
	for _, command := range o.Commands {
		scriptStr := ""
		root := o.Root
		err := error(nil)
		if script, ok := vm.InitScript.Get(command); ok {
			root = root || script.Root
			scriptStr, err = script.GetCommand()
		} else {
			scriptStr = command
		}
		if err != nil {
			printcolor.Warning(fmt.Sprintf("Error building script for %s in %s: %v", command, vm.Name, err))
			continue
		}
		shellArgs := limashell.BuildShellArgs(vm, root, scriptStr)
		if _, _, err := common.ExecShell("limactl", shellArgs...); err != nil {
			printcolor.Error(fmt.Sprintf("Error executing command %s in VM %s in group %s: %v", vm.Name, vm.Group, err))
			continue
		}
	}
}

func executeFiles(vm model.VirtualMachine, o *ExecuteOptions) {
	for _, file := range o.Files {
		shellArgs := limashell.BuildShellArgs(vm, o.Root)
		if fileData, err := os.ReadFile(file); err == nil {
			shellArgs = append(shellArgs, string(fileData))
			if _, _, err := common.ExecShell("limactl", shellArgs...); err != nil {
				printcolor.Error(fmt.Sprintf("Error executing VM %s in group %s: %v", vm.Name, vm.Group, err))
			}
		}
	}
}
