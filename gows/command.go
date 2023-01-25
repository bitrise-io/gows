package gows

import (
	"fmt"
	"io"
	"strings"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
)

// Opts ...
type Opts struct {
	Stdout      io.Writer
	Stderr      io.Writer
	Stdin       io.Reader
	ErrorFinder command.ErrorFinder
}

// CommandFactory ...
type CommandFactory interface {
	Create(name string, args []string, goPath, workDir string, opts *Opts) command.Command
}

type commandFactory struct {
	cmdFactory    command.Factory
	envRepository env.Repository
}

// NewCommandFactory ...
func NewCommandFactory(cmdFactory command.Factory, envRepository env.Repository) CommandFactory {
	return commandFactory{
		cmdFactory:    cmdFactory,
		envRepository: envRepository,
	}
}

// Create creates a command, prepared to run in the isolated workspace environment.
func (f commandFactory) Create(name string, args []string, goPath, workDir string, opts *Opts) command.Command {
	cmdEnvs := f.envRepository.List()
	cmdEnvs = filteredEnvsList(cmdEnvs, "GOPATH")
	cmdEnvs = filteredEnvsList(cmdEnvs, "PWD")
	cmdEnvs = append(cmdEnvs,
		fmt.Sprintf("GOPATH=%s", goPath),
		fmt.Sprintf("PWD=%s", workDir),
	)

	return f.cmdFactory.Create(name, args, &command.Opts{
		Env: cmdEnvs,
		Dir: workDir,

		Stdout:      opts.Stdout,
		Stderr:      opts.Stderr,
		Stdin:       opts.Stdin,
		ErrorFinder: opts.ErrorFinder,
	})
}

func filteredEnvsList(envsList []string, ignoreEnv string) []string {
	filteredEnvs := []string{}
	for _, envItem := range envsList {
		// an env item is a single string with the syntax: key=the value
		if !strings.HasPrefix(envItem, fmt.Sprintf("%s=", ignoreEnv)) {
			filteredEnvs = append(filteredEnvs, envItem)
		}
	}
	return filteredEnvs
}
