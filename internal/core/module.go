package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// ModuleType represents the kind of module
type ModuleType string

const (
	TypeConfig ModuleType = "config"
	TypeTool   ModuleType = "tool"
	TypePlugin ModuleType = "plugin"
)

// LinkStrategy defines how a config module is deployed
type LinkStrategy string

const (
	StrategySymlink LinkStrategy = "symlink"
	StrategyCopy    LinkStrategy = "copy"
	StrategyAppend  LinkStrategy = "append"
)

// LinkMode defines whether we link a file or directory
type LinkMode string

const (
	LinkModeFile      LinkMode = "file"
	LinkModeDirectory LinkMode = "directory"
)

// Hooks defines lifecycle hooks for a module
type Hooks struct {
	PreLink  string `yaml:"pre_link"`
	PostLink string `yaml:"post_link"`
}

// InstallInfo defines how to install a tool
type InstallInfo struct {
	Check         string `yaml:"check"`
	Apt           string `yaml:"apt"`
	PPA           string `yaml:"ppa"`
	Dnf           string `yaml:"dnf"`
	Pacman        string `yaml:"pacman"`
	Brew          string `yaml:"brew"`
	Script        string `yaml:"script"`
	GithubRelease string `yaml:"github_release"`
	AssetPattern  string `yaml:"asset_pattern"`
}

// VersionInfo defines how to detect a tool's version
type VersionInfo struct {
	Command string `yaml:"command"`
	Pattern string `yaml:"pattern"`
}

// Module is the universal module definition
type Module struct {
	Name        string       `yaml:"name"`
	Type        ModuleType   `yaml:"type"`
	Description string       `yaml:"description"`
	Source      string       `yaml:"source"`
	Target      string       `yaml:"target"`
	Strategy    LinkStrategy `yaml:"strategy"`
	LinkMode    LinkMode     `yaml:"link_mode"`
	Platforms   []string     `yaml:"platforms"`
	Group       string       `yaml:"group"`
	Optional    bool         `yaml:"optional"`
	DependsOn   []string     `yaml:"depends_on"`
	Tags        []string     `yaml:"tags"`
	Guard       string       `yaml:"guard"`
	Hooks       Hooks        `yaml:"hooks"`
	Install     InstallInfo  `yaml:"install"`
	Version     VersionInfo  `yaml:"version"`
}

// ModuleFile is the top-level YAML structure
type ModuleFile struct {
	Modules []Module `yaml:"modules"`
}

// ExpandPath expands ~ to home directory
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

// LoadModules loads modules from a YAML file
func LoadModules(path string) ([]Module, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading module file %s: %w", path, err)
	}

	var mf ModuleFile
	if err := yaml.Unmarshal(data, &mf); err != nil {
		return nil, fmt.Errorf("parsing module file %s: %w", path, err)
	}

	return mf.Modules, nil
}

// LoadAllModules loads modules from all YAML files in a directory
func LoadAllModules(dir string) ([]Module, error) {
	var all []Module

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading modules directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		modules, err := LoadModules(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		all = append(all, modules...)
	}

	return all, nil
}

// FilterByPlatform returns only modules that support the current platform
func FilterByPlatform(modules []Module) []Module {
	platform := runtime.GOOS
	// Check for WSL2
	if platform == "linux" {
		if data, err := os.ReadFile("/proc/version"); err == nil {
			if strings.Contains(strings.ToLower(string(data)), "microsoft") {
				// WSL2 is still linux but we might want to note it
				platform = "linux"
			}
		}
	}

	var filtered []Module
	for _, m := range modules {
		if len(m.Platforms) == 0 {
			filtered = append(filtered, m)
			continue
		}
		for _, p := range m.Platforms {
			if p == platform {
				filtered = append(filtered, m)
				break
			}
		}
	}
	return filtered
}

// GroupModules returns a map of group name → modules in that group
func GroupModules(modules []Module) map[string][]Module {
	groups := make(map[string][]Module)
	for _, m := range modules {
		if m.Group != "" {
			groups[m.Group] = append(groups[m.Group], m)
		}
	}
	return groups
}

// IndependentModules returns modules that don't belong to any group
func IndependentModules(modules []Module) []Module {
	var independent []Module
	for _, m := range modules {
		if m.Group == "" {
			independent = append(independent, m)
		}
	}
	return independent
}
