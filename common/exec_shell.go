package common

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"vmctl/util/printcolor"
)

func ExecShell(command ...string) (string, int, error) {
	cmdStr := strings.Join(command, " ")
	cmd := exec.Command("bash", "-c", cmdStr)
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
