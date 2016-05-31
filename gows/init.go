package gows

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-tools/gows/config"
	"gopkg.in/yaml.v2"
)

const (
	userConfigFilePath = "./gows.yml"
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
	log.Debugf("Found Git Remote: %s", outStr)
	packageName, err := parsePackageNameFromURL(strings.TrimSpace(outStr))
	if err != nil {
		return "", fmt.Errorf("Failed to parse package name from remote URL (%s), error: %s", outStr, err)
	}
	if packageName == "" {
		return "", fmt.Errorf("Failed to parse package name from remote URL (%s), error: empty package name parsed", outStr)
	}

	return packageName, nil
}

// Init ...
func Init(packageName string) error {
	log.Debugf("[Init] Initializing package: %s", packageName)
	userConf := config.UserConfigModel{
		PackageName: packageName,
	}

	bytes, err := yaml.Marshal(userConf)
	if err != nil {
		return err
	}

	err = fileutil.WriteBytesToFile(userConfigFilePath, bytes)
	if err != nil {
		return fmt.Errorf("Failed to write User Config into file (%s), error: %s", userConfigFilePath, err)
	}

	log.Debugf("[Init] User Config saved to file: %s", userConfigFilePath)

	return nil
}
