package gows

import (
	"errors"
	"os"
	"os/exec"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

// RunCommand runs the command with it's arguments
// Returns the exit code of the command and any error occured in the function
func RunCommand(cmdName string, cmdArgs ...string) (int, error) {
	log.Debugf("[RunCommand] Command Name: %s", cmdName)
	log.Debugf("[RunCommand] Command Args: %#v", cmdArgs)

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmdExitCode := 0
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus, ok := exitError.Sys().(syscall.WaitStatus)
			if !ok {
				return 1, errors.New("Failed to cast exit status")
			}
			cmdExitCode = waitStatus.ExitStatus()
		}
		return cmdExitCode, err
	}

	return 0, nil
}
