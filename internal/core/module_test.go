package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()
	tests := []struct {
		input    string
		expected string
	}{
		{"~/.config/ghostty", filepath.Join(home, ".config/ghostty")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		result := ExpandPath(tt.input)
		if result != tt.expected {
			t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestLoadModules(t *testing.T) {
	// Create temp YAML
	dir := t.TempDir()
	content := `modules:
  - name: test-mod
    type: config
    description: "Test module"
    source: test/file
    target: ~/.config/test
    strategy: symlink
    link_mode: file
    platforms: [darwin, linux]
    group: "test-group"
    optional: true
    tags: [test]
`
	path := filepath.Join(dir, "test.yaml")
	os.WriteFile(path, []byte(content), 0644)

	mods, err := LoadModules(path)
	if err != nil {
		t.Fatalf("LoadModules error: %v", err)
	}
	if len(mods) != 1 {
		t.Fatalf("expected 1 module, got %d", len(mods))
	}
	m := mods[0]
	if m.Name != "test-mod" {
		t.Errorf("Name = %q, want %q", m.Name, "test-mod")
	}
	if m.Type != TypeConfig {
		t.Errorf("Type = %q, want %q", m.Type, TypeConfig)
	}
	if m.Group != "test-group" {
		t.Errorf("Group = %q, want %q", m.Group, "test-group")
	}
	if m.Strategy != StrategySymlink {
		t.Errorf("Strategy = %q, want %q", m.Strategy, StrategySymlink)
	}
}

func TestFilterByPlatform(t *testing.T) {
	mods := []Module{
		{Name: "mac-only", Platforms: []string{"darwin"}},
		{Name: "linux-only", Platforms: []string{"linux"}},
		{Name: "both", Platforms: []string{"darwin", "linux"}},
		{Name: "no-platform", Platforms: nil},
	}

	filtered := FilterByPlatform(mods)
	// Should include at least "both" and "no-platform"
	if len(filtered) < 2 {
		t.Errorf("expected at least 2 modules, got %d", len(filtered))
	}

	// Check that platform filtering works
	names := make(map[string]bool)
	for _, m := range filtered {
		names[m.Name] = true
	}
	if !names["both"] {
		t.Error("expected 'both' module to be included")
	}
	if !names["no-platform"] {
		t.Error("expected 'no-platform' module to be included")
	}
}

func TestGroupModules(t *testing.T) {
	mods := []Module{
		{Name: "a", Group: "g1"},
		{Name: "b", Group: "g1"},
		{Name: "c", Group: "g2"},
		{Name: "d", Group: ""},
	}

	groups := GroupModules(mods)
	if len(groups["g1"]) != 2 {
		t.Errorf("group g1: expected 2, got %d", len(groups["g1"]))
	}
	if len(groups["g2"]) != 1 {
		t.Errorf("group g2: expected 1, got %d", len(groups["g2"]))
	}
	if _, exists := groups[""]; exists {
		t.Error("empty group should not be in groups map")
	}
}

func TestIndependentModules(t *testing.T) {
	mods := []Module{
		{Name: "a", Group: "g1"},
		{Name: "b", Group: ""},
		{Name: "c", Group: ""},
	}

	indep := IndependentModules(mods)
	if len(indep) != 2 {
		t.Errorf("expected 2 independent, got %d", len(indep))
	}
}
