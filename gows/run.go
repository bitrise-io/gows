package gows

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/gows/config"
)

// PrepareEnvironmentAndRunCommand ...
// Returns the exit code of the command and any error occured in the function
func PrepareEnvironmentAndRunCommand(cmdName string, cmdArgs ...string) (int, error) {
	gowsConfig, err := config.LoadGOWSConfigFromFile()
	if err != nil {
		return 0, fmt.Errorf("Failed to read gows configs: %s", err)
	}
	currWorkDir, err := os.Getwd()
	if err != nil {
		return 0, fmt.Errorf("[PrepareEnvironmentAndRunCommand] Failed to get current working directory: %s", err)
	}

	wsConfig, isFound := gowsConfig.WorkspaceForProjectLocation(currWorkDir)
	if !isFound {
		return 0, fmt.Errorf("No Workspace configuration found for the current project / working directory: %s", currWorkDir)
	}

	projectConfig, err := config.LoadProjectConfigFromFile()
	if err != nil {
		return 0, fmt.Errorf("Failed to read Project Config: %s", err)
	}

	cmdWorkdir, err := prepareGoWorkspaceEnvironment(currWorkDir, wsConfig, projectConfig)
	if err != nil {
		return 0, err
	}

	return runCommand(wsConfig.WorkspaceRootPath, cmdWorkdir, cmdName, cmdArgs...)
}

func prepareGoWorkspaceEnvironment(currWorkDir string, wsConfig config.WorkspaceConfigModel, projectConfig config.ProjectConfigModel) (string, error) {
	if wsConfig.WorkspaceRootPath == "" {
		return "", fmt.Errorf("No gows Workspace root path found for the current project / working directory: %s", currWorkDir)
	}
	if projectConfig.PackageName == "" {
		return "", errors.New("No Package Name specified - make sure you initialized the workspace (with: gows init)")
	}

	fullPackageWorkspacePath := filepath.Join(wsConfig.WorkspaceRootPath, "src", projectConfig.PackageName)

	log.Debugf("=> Creating Symlink: (%s) -> (%s)", currWorkDir, fullPackageWorkspacePath)

	// create symlink, if not yet created
	fileInfo, isExists, err := pathutil.PathCheckAndInfos(fullPackageWorkspacePath)
	if err != nil {
		return "", fmt.Errorf("Failed to check Symlink status (at: %s), error: %s", fullPackageWorkspacePath, err)
	}
	isSymlinkAlreadyInPlace := false
	if isExists && fileInfo.Mode()&os.ModeSymlink != 0 {
		log.Debug(" Symlink already exists")
		originPth, err := os.Readlink(fullPackageWorkspacePath)
		if err != nil {
			return "", fmt.Errorf("Symlink found (at: %s), but failed to open: %s", fullPackageWorkspacePath, err)
		}

		if originPth == currWorkDir {
			isSymlinkAlreadyInPlace = true
		} else {
			// remove
			log.Warning("Symlink already exists (at: %s), but target (%s) is not the current one (%s)", fullPackageWorkspacePath, originPth, currWorkDir)
			log.Warning("Removing and re-creating the symlink ...")
			if err := os.Remove(fullPackageWorkspacePath); err != nil {
				return "", fmt.Errorf("Failed to remove Symlink (at: %s), error: %s", fullPackageWorkspacePath, err)
			}
		}
	}

	if !isSymlinkAlreadyInPlace {
		log.Debug(" Creating symlink ...")
		// create the parent directory
		if err := os.MkdirAll(filepath.Dir(fullPackageWorkspacePath), 0777); err != nil {
			return "", fmt.Errorf("Failed to create base directory for symlink into: %s", fullPackageWorkspacePath)
		}
		// create symlink
		if err := os.Symlink(currWorkDir, fullPackageWorkspacePath); err != nil {
			return "", fmt.Errorf("Failed to create symlink from project directory (%s) into gows Workspace directory (%s), error: %s", currWorkDir, fullPackageWorkspacePath, err)
		}
	}

	log.Debugf(" [DONE] Symlink is in place")

	return fullPackageWorkspacePath, nil
}

// runCommand runs the command with it's arguments
// Returns the exit code of the command and any error occured in the function
func runCommand(gopath, cmdWorkdir, cmdName string, cmdArgs ...string) (int, error) {
	log.Debugf("[RunCommand] Command Name: %s", cmdName)
	log.Debugf("[RunCommand] Command Args: %#v", cmdArgs)
	log.Debugf("[RunCommand] Command Work Dir: %#v", cmdWorkdir)

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cmdWorkdir
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOPATH=%s", gopath))

	cmdExitCode := 0
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus, ok := exitError.Sys().(syscall.WaitStatus)
			if !ok {
				return 0, errors.New("Failed to cast exit status")
			}
			cmdExitCode = waitStatus.ExitStatus()
		}
		return cmdExitCode, err
	}

	return 0, nil
}
