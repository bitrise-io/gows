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
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-tools/gows/config"
	"gopkg.in/yaml.v2"
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
		return err
	}
	if err := os.MkdirAll(filepath.Join(wsRootPath, "bin"), 0777); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(wsRootPath, "pkg"), 0777); err != nil {
		return err
	}
	return nil
}

// Init ...
func Init(packageName string) error {
	log.Debugf("[Init] Initializing package: %s", packageName)

	gowsConfig, err := config.LoadGOWSConfigFromFile()
	if err != nil {
		return err
	}

	log.Debug("[Init] Initializing User Config ...")
	{
		userConf := config.UserConfigModel{
			PackageName: packageName,
		}

		bytes, err := yaml.Marshal(userConf)
		if err != nil {
			return err
		}

		err = fileutil.WriteBytesToFile(config.UserConfigFilePath, bytes)
		if err != nil {
			return fmt.Errorf("Failed to write User Config into file (%s), error: %s", config.UserConfigFilePath, err)
		}
	}

	log.Debugf("[Init] User Config saved to file: %s", config.UserConfigFilePath)

	// Workspace Config
	log.Debug("[Init] Initializing Workspace & Config ...")
	// Create the Workspace
	gowsWorspacesRootDirAbsPath, err := config.GOWSWorspacesRootDirAbsPath()
	if err != nil {
		return fmt.Errorf("Failed to get absolute path for gows workspaces root dir, error: %s", err)
	}
	currWorkDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Failed to get current working directory: %s", err)
	}
	projectBaseWorkspaceDirName := fmt.Sprintf("%s-%d", filepath.Base(currWorkDir), time.Now().Unix())
	projectWorkspaceAbsPath := filepath.Join(gowsWorspacesRootDirAbsPath, projectBaseWorkspaceDirName)
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
		gowsConfig.Workspaces[currWorkDir] = workspaceConf

		if err := config.SaveGOWSConfigToFile(gowsConfig); err != nil {
			return err
		}
	}
	log.Debug("[Init] Workspace Config saved")

	return nil
}
