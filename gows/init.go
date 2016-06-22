package gows

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-tools/gows/config"
)

func parsePackageNameFromURL(remoteURL string) (string, error) {
	origRemoteURL := remoteURL
	if strings.HasPrefix(remoteURL, "git@") {
		remoteURL = "ssh://" + remoteURL
	}

	u, err := url.Parse(remoteURL)
	if err != nil {
		return "", fmt.Errorf("Failed to parse remote URL (%s): %s", origRemoteURL, err)
	}

	packagePth := u.Path
	packagePth = strings.TrimSuffix(packagePth, ".git")

	// in SSH git urls like "ssh://git@github.com:bitrise-io/go-utils.git" Go parses "github.com:bitrise-io" as the host
	// fix it by splitting it and replacing ":" with "/"
	hostSplits := strings.Split(u.Host, ":")
	host := hostSplits[0]
	if len(hostSplits) > 1 {
		if len(hostSplits) > 2 {
			return "", fmt.Errorf("More than one ':' found in the Host part of the URL (%s)", origRemoteURL)
		}
		packagePth = "/" + hostSplits[1] + packagePth
	}

	if host == "" {
		return "", fmt.Errorf("No Host found in URL (%s)", origRemoteURL)
	}
	if packagePth == "" || packagePth == "/" {
		return "", fmt.Errorf("No Path found in URL (%s)", origRemoteURL)
	}

	return host + packagePth, nil
}

// AutoScanPackageName ...
func AutoScanPackageName() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Errorf("[AutoScanPackageName] (Error) Output was: %s", exitError.Stderr)
		} else {
			log.Error("[AutoScanPackageName] Failed to convert error to ExitError")
		}
		return "", fmt.Errorf("Failed to get git remote url for origin: %s", err)
	}

	outStr := string(out)
	gitRemoteStr := strings.TrimSpace(outStr)
	log.Debugf("Found Git Remote: %s", gitRemoteStr)
	packageName, err := parsePackageNameFromURL(gitRemoteStr)
	if err != nil {
		return "", fmt.Errorf("Failed to parse package name from remote URL (%s), error: %s", outStr, err)
	}
	if packageName == "" {
		return "", fmt.Errorf("Failed to parse package name from remote URL (%s), error: empty package name parsed", outStr)
	}

	return packageName, nil
}

func initGoWorkspaceAtPath(wsRootPath string) error {
	if err := os.MkdirAll(filepath.Join(wsRootPath, "src"), 0777); err != nil {
		return fmt.Errorf("Failed to create GOPATH/src directory: %s", err)
	}
	return nil
}

// initWorkspaceForProjectPath ...
// Workspaces are linked to project paths, not to package IDs!
// You can have multiple workspaces for the same package ID, but not for the
// same (project) path.
func initWorkspaceForProjectPath(projectPath string, isAllowReset bool) error {
	log.Debug("[Init] Initializing Workspace & Config ...")

	gowsWorspacesRootDirAbsPath, err := config.GOWSWorspacesRootDirAbsPath()
	if err != nil {
		return fmt.Errorf("Failed to get absolute path for gows workspaces root dir, error: %s", err)
	}

	// Create the Workspace
	gowsConfig, err := config.LoadGOWSConfigFromFile()
	if err != nil {
		return fmt.Errorf("Failed to load gows config: %s", err)
	}

	projectWorkspaceAbsPath := ""
	wsConfig, isFound := gowsConfig.WorkspaceForProjectLocation(projectPath)
	if isFound {
		if wsConfig.WorkspaceRootPath == "" {
			return fmt.Errorf("A workspace is found for this project (path: %s), but the workspace root directory path is not defined!", projectPath)
		}
		projectWorkspaceAbsPath = wsConfig.WorkspaceRootPath

		if isAllowReset {
			if err := os.RemoveAll(projectWorkspaceAbsPath); err != nil {
				return fmt.Errorf("Failed to delete previous workspace at path: %s", projectWorkspaceAbsPath)
			}
			// init a new one
			projectWorkspaceAbsPath = ""
		} else {
			log.Warning(colorstring.Yellow("A workspace already exists for this project") + " (" + projectWorkspaceAbsPath + "), will be reused.")
			log.Warning("If you want to delete the previous workspace of this project and generate a new one you should run: " + colorstring.Green("gows init -reset"))
		}
	}

	if projectWorkspaceAbsPath == "" {
		// generate one
		projectBaseWorkspaceDirName := fmt.Sprintf("%s-%d", filepath.Base(projectPath), time.Now().Unix())
		projectWorkspaceAbsPath = filepath.Join(gowsWorspacesRootDirAbsPath, projectBaseWorkspaceDirName)
	}

	log.Debugf("  projectWorkspaceAbsPath: %s", projectWorkspaceAbsPath)
	if err := initGoWorkspaceAtPath(projectWorkspaceAbsPath); err != nil {
		return fmt.Errorf("Failed to initialize workspace at path: %s", projectWorkspaceAbsPath)
	}
	log.Debugf("  Workspace successfully created")

	// Save the location into Workspace config
	{
		workspaceConf := config.WorkspaceConfigModel{
			WorkspaceRootPath: projectWorkspaceAbsPath,
		}
		gowsConfig.Workspaces[projectPath] = workspaceConf

		if err := config.SaveGOWSConfigToFile(gowsConfig); err != nil {
			return fmt.Errorf("Failed to save gows config: %s", err)
		}
	}
	log.Debug("[Init] Workspace Config saved")

	return nil
}

// Init ...
func Init(packageName string, isAllowReset bool) error {
	log.Debugf("[Init] Initializing package: %s", packageName)

	log.Debug("[Init] Initializing Project Config ...")
	{
		projectConf := config.ProjectConfigModel{
			PackageName: packageName,
		}

		if err := config.SaveProjectConfigToFile(projectConf); err != nil {
			return fmt.Errorf("Failed to write Project Config into file: %s", err)
		}
	}

	log.Debugf("[Init] Project Config saved to file: %s", config.ProjectConfigFilePath)

	// init workspace for project (path)
	currWorkDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Failed to get current working directory: %s", err)
	}

	if err := initWorkspaceForProjectPath(currWorkDir, isAllowReset); err != nil {
		return fmt.Errorf("Failed to initialize Workspace for Project: %s", err)
	}

	return nil
}
