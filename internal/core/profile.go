package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Profile represents a predefined module selection
type Profile struct {
	Name         string            `yaml:"name"`
	Description  string            `yaml:"description"`
	Modules      []string          `yaml:"modules"`
	GroupChoices map[string]string  `yaml:"group_choices"`
}

// LoadProfile reads a profile from a YAML file
func LoadProfile(path string) (*Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading profile %s: %w", path, err)
	}

	var p Profile
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parsing profile %s: %w", path, err)
	}
	return &p, nil
}

// LoadAllProfiles loads all profiles from the profiles directory
func LoadAllProfiles(dir string) ([]*Profile, error) {
	var profiles []*Profile

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading profiles directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		p, err := LoadProfile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// ResolveProfile returns the list of modules to install/link based on a profile
// It resolves group choices and filters by available modules
func ResolveProfile(profile *Profile, allModules []Module) (selected []Module, groupChoices map[string]string) {
	groupChoices = make(map[string]string)
	modMap := make(map[string]Module)
	for _, m := range allModules {
		modMap[m.Name] = m
	}

	// Resolve explicit modules
	for _, name := range profile.Modules {
		if mod, ok := modMap[name]; ok {
			selected = append(selected, mod)
		}
	}

	// Resolve group choices
	for group, choice := range profile.GroupChoices {
		if choice == "none" || choice == "" {
			groupChoices[group] = "none"
			continue
		}
		if choice == "prompt" {
			groupChoices[group] = "prompt"
			continue
		}
		// Specific module chosen
		if mod, ok := modMap[choice]; ok {
			selected = append(selected, mod)
			groupChoices[group] = choice
		}
	}

	return selected, groupChoices
}
