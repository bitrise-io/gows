package config

import (
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"gopkg.in/yaml.v2"
)

const (
	// ProjectConfigFilePath ...
	ProjectConfigFilePath = "./gows.yml"

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

// ProjectConfigFileAbsPath ...
func ProjectConfigFileAbsPath() (string, error) {
	return pathutil.AbsPath(ProjectConfigFilePath)
}

// ProjectConfigModel - stored in ./gows.yml
type ProjectConfigModel struct {
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

func createDefaultGOWSConfigModel() GOWSConfigModel {
	return GOWSConfigModel{
		Workspaces: map[string]WorkspaceConfigModel{},
	}
}

// WorkspaceForProjectLocation ...
func (gowsConfig GOWSConfigModel) WorkspaceForProjectLocation(projectPath string) (WorkspaceConfigModel, bool) {
	wsConfig, isFound := gowsConfig.Workspaces[projectPath]
	return wsConfig, isFound
}

// LoadGOWSConfigFromFile ...
func LoadGOWSConfigFromFile() (GOWSConfigModel, error) {
	gowsConfigFileAbsPath, err := GOWSConfigFileAbsPath()
	if err != nil {
		return GOWSConfigModel{}, fmt.Errorf("Failed to get absolute path of gows config: %s", err)
	}

	// If doesn't exist yet, return a default/empty gows config
	{
		isExists, err := pathutil.IsPathExists(gowsConfigFileAbsPath)
		if !isExists {
			log.Debugf(" (!) cows Config does not yet exists at: %s", gowsConfigFileAbsPath)
			// return an empty/default config
			return createDefaultGOWSConfigModel(), nil
		} else if err != nil {
			return GOWSConfigModel{}, err
		}
	}

	bytes, err := ioutil.ReadFile(gowsConfigFileAbsPath)
	if err != nil {
		return GOWSConfigModel{}, fmt.Errorf("Failed to read gows config file (%s), error: %s", gowsConfigFileAbsPath, err)
	}
	var gowsConfig GOWSConfigModel
	if err := yaml.Unmarshal(bytes, &gowsConfig); err != nil {
		return GOWSConfigModel{}, fmt.Errorf("Failed to parse gows config (should be valid YML, path: %s), error: %s", gowsConfigFileAbsPath, err)
	}

	return gowsConfig, nil
}

// SaveGOWSConfigToFile ...
func SaveGOWSConfigToFile(gowsConfig GOWSConfigModel) error {
	bytes, err := yaml.Marshal(gowsConfig)
	if err != nil {
		return fmt.Errorf("Failed to generate YML for gows config: %s", err)
	}

	gowsConfigFileAbsPath, err := GOWSConfigFileAbsPath()
	if err != nil {
		return fmt.Errorf("Failed to get absolute path of gows config: %s", err)
	}

	err = fileutil.WriteBytesToFile(gowsConfigFileAbsPath, bytes)
	if err != nil {
		return fmt.Errorf("Failed to write Project Config into file (%s), error: %s", gowsConfigFileAbsPath, err)
	}

	return nil
}

// LoadProjectConfigFromFile ...
func LoadProjectConfigFromFile() (ProjectConfigModel, error) {
	projectConfigFileAbsPath, err := ProjectConfigFileAbsPath()
	if err != nil {
		return ProjectConfigModel{}, fmt.Errorf("Failed to get absolute path of project config: %s", err)
	}

	bytes, err := ioutil.ReadFile(projectConfigFileAbsPath)
	if err != nil {
		return ProjectConfigModel{}, fmt.Errorf("Failed to read project config file (%s), error: %s", projectConfigFileAbsPath, err)
	}
	var projectConfig ProjectConfigModel
	if err := yaml.Unmarshal(bytes, &projectConfig); err != nil {
		return ProjectConfigModel{}, fmt.Errorf("Failed to parse project config (should be valid YML, path: %s), error: %s", projectConfigFileAbsPath, err)
	}

	return projectConfig, nil
}

// SaveProjectConfigToFile ...
func SaveProjectConfigToFile(projectConf ProjectConfigModel) error {
	bytes, err := yaml.Marshal(projectConf)
	if err != nil {
		return fmt.Errorf("Failed to parse Project Config (should be valid YML): %s", err)
	}

	err = fileutil.WriteBytesToFile(ProjectConfigFilePath, bytes)
	if err != nil {
		return fmt.Errorf("Failed to write Project Config into file (%s), error: %s", ProjectConfigFilePath, err)
	}

	return nil
}
