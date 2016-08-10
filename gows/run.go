package gows

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/gows/config"
)

// PrepareEnvironmentAndRunCommand ...
// Returns the exit code of the command and any error occured in the function
func PrepareEnvironmentAndRunCommand(cmdName string, cmdArgs ...string) (int, error) {
	projectConfig, err := config.LoadProjectConfigFromFile()
	if err != nil {
		log.Info("Run " + colorstring.Green("gows init") + " to initialize a workspace & gows config for this project")
		return 0, fmt.Errorf("Failed to read Project Config: %s", err)
	}

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
		log.Debugln("No initialized workspace dir found for this project, initializing one ...")
		if err := initWorkspaceForProjectPath(currWorkDir, false); err != nil {
			return 0, fmt.Errorf("[PrepareEnvironmentAndRunCommand] Failed to initialize Workspace for Project: %s", err)
		}
		log.Debugln("[DONE] workspace dir initialized - continue running ...")

		// reload config
		gowsConfig, err := config.LoadGOWSConfigFromFile()
		if err != nil {
			return 0, fmt.Errorf("Failed to read gows configs: %s", err)
		}
		wsConfig, isFound = gowsConfig.WorkspaceForProjectLocation(currWorkDir)
	}
	if !isFound {
		log.Info("Run " + colorstring.Green("gows init") + " to initialize a workspace & gows config for this project")
		return 0, fmt.Errorf("No Workspace configuration found for the current project / working directory: %s", currWorkDir)
	}

	userConfig, err := config.LoadUserConfigFromFile()
	if err != nil {
		log.Debug("No User Config found, using defaults")
		userConfig = config.CreateDefaultUserConfig()
	}
	log.Debugf("Using User Config: %#v", userConfig)

	origGOPATH := os.Getenv("GOPATH")
	if origGOPATH == "" {
		return 0, fmt.Errorf("You don't have a GOPATH environment - please set it; GOPATH/bin will be symlinked")
	}

	if wsConfig.WorkspaceRootPath == "" {
		return 0, fmt.Errorf("No gows Workspace root path found for the current project / working directory: %s", currWorkDir)
	}
	if projectConfig.PackageName == "" {
		return 0, errors.New("No Package Name specified - make sure you initialized the workspace (with: gows init)")
	}

	if err := pathutil.EnsureDirExist(wsConfig.WorkspaceRootPath); err != nil {
		return 0, fmt.Errorf("Failed to create workspace root directory (path: %s), error: %s", wsConfig.WorkspaceRootPath, err)
	}

	if err := createGopathBinSymlink(origGOPATH, wsConfig); err != nil {
		return 0, fmt.Errorf("Failed to create GOPATH/bin symlink, error: %s", err)
	}

	fullPackageWorkspacePath := filepath.Join(wsConfig.WorkspaceRootPath, "src", projectConfig.PackageName)

	userConfigSyncMode := userConfig.SyncMode
	if userConfigSyncMode == "" {
		userConfigSyncMode = config.DefaultSyncMode
	}
	log.Debug("[PrepareEnvironmentAndRunCommand] specified Sync Mode : ", userConfigSyncMode)

	// prepare
	switch userConfigSyncMode {
	case config.SyncModeSymlink:
		// create symlink for Project->Workspace
		log.Debugf("=> Creating Symlink: (%s) -> (%s)", currWorkDir, fullPackageWorkspacePath)
		if err := createOrUpdateSymlink(currWorkDir, fullPackageWorkspacePath); err != nil {
			return 0, fmt.Errorf("Failed to create Project->Workspace symlink, error: %s", err)
		}
		log.Debugf(" [DONE] Symlink is in place")
	case config.SyncModeCopy:
		// Sync project into workspace
		log.Debugf("=> Sync project content into workspace: (%s) -> (%s)", currWorkDir, fullPackageWorkspacePath)
		if err := syncDirWithDir(currWorkDir, fullPackageWorkspacePath); err != nil {
			return 0, fmt.Errorf("Failed to sync the project path / workdir into the Workspace, error: %s", err)
		}
		log.Debugf(" [DONE] Sync project content into workspace")
	default:
		return 0, fmt.Errorf("Unsupported Sync Mode: %s", userConfigSyncMode)
	}

	// Run the command, in the prepared Workspace
	exitCode, cmdErr := runCommand(origGOPATH, fullPackageWorkspacePath, wsConfig, cmdName, cmdArgs...)

	// cleanup / finishing
	switch userConfigSyncMode {
	case config.SyncModeSymlink:
		// nothing to do
	case config.SyncModeCopy:
		// Sync back from workspace into project
		log.Debugf("=> Sync workspace content into project: (%s) -> (%s)", fullPackageWorkspacePath, currWorkDir)
		if err := syncDirWithDir(fullPackageWorkspacePath, currWorkDir); err != nil {
			// we should return the command's exit code and error (if any)
			// maybe if the exitCode==0 and cmdErr==nil only then we could return an error here ...
			// for now we'll just print an error log, but it won't change the "output" of this function
			log.Errorf("Failed to sync back the project content from the Workspace, error: %s", err)
		} else {
			log.Debugf(" [DONE] Sync back project content from workspace")
		}
	default:
		return 0, fmt.Errorf("Unsupported Sync Mode: %s", userConfigSyncMode)
	}

	return exitCode, cmdErr
}

func syncDirWithDir(syncContentOf, syncIntoDir string) error {
	syncContentOf = filepath.Clean(syncContentOf)
	syncIntoDir = filepath.Clean(syncIntoDir)

	if err := pathutil.EnsureDirExist(syncIntoDir); err != nil {
		return fmt.Errorf("Failed to create target (at: %s), error: %s", syncIntoDir, err)
	}

	cmd := exec.Command("rsync", "-avhP", "--delete", syncContentOf+"/", syncIntoDir+"/")
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Error("[syncDirWithDir] Sync Error")
			log.Errorf("[syncDirWithDir] Output (Stdout) was: %s", out)
			log.Errorf("[syncDirWithDir] Error Output (Stderr) was: %s", exitError.Stderr)
		} else {
			log.Error("[syncDirWithDir] Failed to convert error to ExitError")
		}
		return fmt.Errorf("Failed to rsync between (%s) and (%s), error: %s", syncContentOf, syncIntoDir, err)
	}
	return nil
}

func createOrUpdateSymlink(symlinkTargetPath, symlinkLocationPath string) error {
	fileInfo, isExists, err := pathutil.PathCheckAndInfos(symlinkLocationPath)
	if err != nil {
		return fmt.Errorf("Failed to check Symlink status (at: %s), error: %s", symlinkLocationPath, err)
	}
	isSymlinkAlreadyInPlace := false
	if isExists && fileInfo.Mode()&os.ModeSymlink != 0 {
		log.Debug(" Symlink already exists")
		originPth, err := os.Readlink(symlinkLocationPath)
		if err != nil {
			return fmt.Errorf("Symlink found (at: %s), but failed to open: %s", symlinkLocationPath, err)
		}

		if originPth == symlinkTargetPath {
			isSymlinkAlreadyInPlace = true
		} else {
			// remove
			log.Warningf("Symlink already exists (at: %s), but target (%s) is not the current one (%s)", symlinkLocationPath, originPth, symlinkTargetPath)
			log.Warning("Removing and re-creating the symlink ...")
			if err := os.Remove(symlinkLocationPath); err != nil {
				return fmt.Errorf("Failed to remove Symlink (at: %s), error: %s", symlinkLocationPath, err)
			}
		}
	}

	if !isSymlinkAlreadyInPlace {
		log.Debug(" Creating symlink ...")
		// create the parent directory
		if err := os.MkdirAll(filepath.Dir(symlinkLocationPath), 0777); err != nil {
			return fmt.Errorf("Failed to create base directory for symlink into: %s", symlinkLocationPath)
		}
		// create symlink
		if err := os.Symlink(symlinkTargetPath, symlinkLocationPath); err != nil {
			return fmt.Errorf("Failed to create symlink from project directory (%s) into gows Workspace directory (%s), error: %s", symlinkTargetPath, symlinkLocationPath, err)
		}
	}

	return nil
}

func createGopathBinSymlink(origGOPATH string, wsConfig config.WorkspaceConfigModel) error {
	fullWorkspaceBinPath := filepath.Join(wsConfig.WorkspaceRootPath, "bin")
	originalGopathBinPath, err := pathutil.AbsPath(filepath.Join(origGOPATH, "bin"))
	if err != nil {
		return fmt.Errorf("Failed to get the path of 'bin' dir inside your GOPATH (%s), error: %s", origGOPATH, err)
	}

	log.Debugf("=> Creating Symlink: (%s) -> (%s)", originalGopathBinPath, fullWorkspaceBinPath)

	// create symlink for GOPATH/bin, if not yet created
	if err := createOrUpdateSymlink(originalGopathBinPath, fullWorkspaceBinPath); err != nil {
		return fmt.Errorf("Failed to create GOPATH/bin symlink, error: %s", err)
	}

	log.Debugf(" [DONE] Symlink is in place")

	return nil
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

// runCommand runs the command with it's arguments
// Returns the exit code of the command and any error occured in the function
func runCommand(originalGOPATH, cmdWorkdir string, wsConfig config.WorkspaceConfigModel, cmdName string, cmdArgs ...string) (int, error) {
	log.Debugf("[RunCommand] Command Name: %s", cmdName)
	log.Debugf("[RunCommand] Command Args: %#v", cmdArgs)
	log.Debugf("[RunCommand] Command Work Dir: %#v", cmdWorkdir)

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cmdWorkdir
	//
	cmdEnvs := os.Environ()
	cmdEnvs = filteredEnvsList(cmdEnvs, "GOPATH")
	cmdEnvs = filteredEnvsList(cmdEnvs, "PWD")
	cmdEnvs = append(cmdEnvs,
		fmt.Sprintf("GOPATH=%s", wsConfig.WorkspaceRootPath),
		fmt.Sprintf("PWD=%s", cmdWorkdir),
	)
	cmd.Env = cmdEnvs

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
