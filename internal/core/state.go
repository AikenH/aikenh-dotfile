package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ModuleState tracks the installed state of a single module
type ModuleState struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installed_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SourceHash  string    `json:"source_hash"`  // git commit hash of source
	BackupPath  string    `json:"backup_path"`  // path to backup if we replaced existing
	Strategy    string    `json:"strategy"`     // how it was deployed
	TargetPath  string    `json:"target_path"`  // resolved target path
}

// State is the persistent state for this machine
type State struct {
	Platform   string                 `json:"platform"`
	LastSync   time.Time              `json:"last_sync"`
	Proxy      string                 `json:"proxy"`
	Modules    map[string]ModuleState `json:"modules"`
	GroupPicks map[string]string      `json:"group_picks"` // group → chosen module name (or "none")
}

// StateDir returns the state directory path
func StateDir() string {
	// XDG_DATA_HOME or default
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, _ := os.UserHomeDir()
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, "dotsetup")
}

// StatePath returns the state file path
func StatePath() string {
	return filepath.Join(StateDir(), "state.json")
}

// BackupDir returns the backup directory path
func BackupDir() string {
	return filepath.Join(StateDir(), "backups")
}

// LoadState reads the state from disk, or returns a fresh state if not found
func LoadState() (*State, error) {
	path := StatePath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			p := DetectPlatform()
			return &State{
				Platform:   p.String(),
				Modules:    make(map[string]ModuleState),
				GroupPicks: make(map[string]string),
			}, nil
		}
		return nil, fmt.Errorf("reading state: %w", err)
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing state: %w", err)
	}
	if s.Modules == nil {
		s.Modules = make(map[string]ModuleState)
	}
	if s.GroupPicks == nil {
		s.GroupPicks = make(map[string]string)
	}
	return &s, nil
}

// Save writes the state to disk
func (s *State) Save() error {
	dir := StateDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating state dir: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}

	return os.WriteFile(StatePath(), data, 0644)
}

// IsInstalled checks if a module is recorded as installed
func (s *State) IsInstalled(name string) bool {
	_, ok := s.Modules[name]
	return ok
}

// SetInstalled records a module as installed
func (s *State) SetInstalled(name string, ms ModuleState) {
	ms.Name = name
	if ms.InstalledAt.IsZero() {
		ms.InstalledAt = time.Now()
	}
	ms.UpdatedAt = time.Now()
	s.Modules[name] = ms
}

// SetUninstalled removes a module from state
func (s *State) SetUninstalled(name string) {
	delete(s.Modules, name)
}

// SetGroupChoice records which module was chosen for a group
func (s *State) SetGroupChoice(group, module string) {
	s.GroupPicks[group] = module
}

// GetGroupChoice returns the chosen module for a group
func (s *State) GetGroupChoice(group string) string {
	return s.GroupPicks[group]
}
