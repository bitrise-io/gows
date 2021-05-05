package goutil

import (
	"fmt"
	"github.com/whilp/git-urls"
	"strings"
)

// ParsePackageNameFromURL - returns a Go package name/id (e.g. github.com/bitrise-io/gows)
// from a git clone URL (e.g. https://github.com/bitrise-io/gows.git)
func ParsePackageNameFromURL(remoteURL string) (string, error) {
	u, err := giturls.Parse(remoteURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse remote URL (%s): %s", remoteURL, err)
	}

	packagePth := u.Path
	packagePth = strings.TrimSuffix(packagePth, ".git")
	host := u.Host
	if host == "" {
		return "", fmt.Errorf("no Host found in URL (%s)", remoteURL)
	}
	if packagePth == "" || packagePth == "/" {
		return "", fmt.Errorf("no Path found in URL (%s)", remoteURL)
	}
	if !strings.HasPrefix(packagePth, "/") {
		packagePth = "/" + packagePth
	}
	return u.Host + packagePth, nil
}
