package config

// UserConfigModel - stored in ./gows-config.yml
type UserConfigModel struct {
	PackageName string `json:"package_name" yaml:"package_name"`
}

// WorkspaceConfigModel - stored in ./.gows
type WorkspaceConfigModel struct {
	WorkspaceRootPath string `json:"workspace_root_path" yaml:"workspace_root_path"`
}
