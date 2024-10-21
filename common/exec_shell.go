package common

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
