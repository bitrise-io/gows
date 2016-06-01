package config

import (
	"fmt"
	"io/ioutil"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"gopkg.in/yaml.v2"
)

const (
	// UserConfigFilePath ...
	UserConfigFilePath = "./gows.yml"

	gowsWorspacesRootDirPath = "$HOME/.bitrise-gows/wsdirs"
	gowsConfigFilePath       = "$HOME/.bitrise-gows/workspaces.yml"
)

// GOWSWorspacesRootDirAbsPath ...
func GOWSWorspacesRootDirAbsPath() (string, error) {
	return pathutil.AbsPath(gowsWorspacesRootDirPath)
}

// GOWSConfigFileAbsPath ...
func GOWSConfigFileAbsPath() (string, error) {
	return pathutil.AbsPath(gowsConfigFilePath)
}

// UserConfigModel - stored in ./gows-config.yml
type UserConfigModel struct {
	PackageName string `json:"package_name" yaml:"package_name"`
}

// WorkspaceConfigModel ...
type WorkspaceConfigModel struct {
	WorkspaceRootPath string `json:"workspace_root_path" yaml:"workspace_root_path"`
}

// GOWSConfigModel ...
type GOWSConfigModel struct {
	Workspaces map[string]WorkspaceConfigModel `json:"workspaces" yaml:"workspaces"`
}

// LoadGOWSConfigFromFile ...
func LoadGOWSConfigFromFile() (GOWSConfigModel, error) {
	gowsConfigFileAbsPath, err := GOWSConfigFileAbsPath()
	if err != nil {
		return GOWSConfigModel{}, err
	}

	bytes, err := ioutil.ReadFile(gowsConfigFileAbsPath)
	if err != nil {
		return GOWSConfigModel{}, err
	}
	var gowsConfig GOWSConfigModel
	if err := yaml.Unmarshal(bytes, &gowsConfig); err != nil {
		return GOWSConfigModel{}, err
	}

	return gowsConfig, nil
}

// SaveGOWSConfigToFile ...
func SaveGOWSConfigToFile(gowsConfig GOWSConfigModel) error {
	bytes, err := yaml.Marshal(gowsConfig)
	if err != nil {
		return err
	}

	gowsConfigFileAbsPath, err := GOWSConfigFileAbsPath()
	if err != nil {
		return err
	}

	err = fileutil.WriteBytesToFile(gowsConfigFileAbsPath, bytes)
	if err != nil {
		return fmt.Errorf("Failed to write User Config into file (%s), error: %s", gowsConfigFileAbsPath, err)
	}

	return nil
}
