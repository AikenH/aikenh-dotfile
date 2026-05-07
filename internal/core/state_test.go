package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestState_LoadSave(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("XDG_DATA_HOME", dir)
	defer os.Setenv("XDG_DATA_HOME", "")

	// Fresh state
	state, err := LoadState()
	if err != nil {
		t.Fatalf("LoadState error: %v", err)
	}
	if state.Modules == nil {
		t.Fatal("Modules should not be nil")
	}

	// Set some values
	state.Proxy = "http://127.0.0.1:7890"
	state.SetInstalled("vim", ModuleState{
		Version:  "9.0",
		Strategy: "symlink",
	})
	state.SetGroupChoice("terminal-emulator", "ghostty")

	// Save
	if err := state.Save(); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	// Verify file exists
	statePath := filepath.Join(dir, "dotsetup", "state.json")
	if _, err := os.Stat(statePath); err != nil {
		t.Fatalf("state file not found: %v", err)
	}

	// Reload
	state2, err := LoadState()
	if err != nil {
		t.Fatalf("Reload error: %v", err)
	}
	if state2.Proxy != "http://127.0.0.1:7890" {
		t.Errorf("Proxy = %q, want %q", state2.Proxy, "http://127.0.0.1:7890")
	}
	if !state2.IsInstalled("vim") {
		t.Error("vim should be installed")
	}
	if state2.GetGroupChoice("terminal-emulator") != "ghostty" {
		t.Errorf("group choice = %q, want ghostty", state2.GetGroupChoice("terminal-emulator"))
	}
}

func TestState_SetUninstalled(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("XDG_DATA_HOME", dir)
	defer os.Setenv("XDG_DATA_HOME", "")

	state, _ := LoadState()
	state.SetInstalled("test", ModuleState{Version: "1.0"})

	if !state.IsInstalled("test") {
		t.Error("should be installed")
	}

	state.SetUninstalled("test")
	if state.IsInstalled("test") {
		t.Error("should not be installed after removal")
	}
}

func TestState_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	os.Setenv("XDG_DATA_HOME", dir)
	defer os.Setenv("XDG_DATA_HOME", "")

	state, _ := LoadState()
	state.SetInstalled("vim", ModuleState{Version: "9.0"})
	state.Save()

	// Read raw JSON and verify it's well-formed
	statePath := filepath.Join(dir, "dotsetup", "state.json")
	data, _ := os.ReadFile(statePath)

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("state.json is not valid JSON: %v", err)
	}

	// Should have expected keys
	if _, ok := parsed["modules"]; !ok {
		t.Error("state.json missing 'modules' key")
	}
	if _, ok := parsed["platform"]; !ok {
		t.Error("state.json missing 'platform' key")
	}
}
