package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckLinkStatus_Missing(t *testing.T) {
	mod := Module{
		Name:   "test",
		Source: "test-src",
		Target: filepath.Join(t.TempDir(), "nonexistent"),
	}
	status := CheckLinkStatus(mod, t.TempDir())
	if status != StatusMissing {
		t.Errorf("expected StatusMissing, got %v", status)
	}
}

func TestCheckLinkStatus_Linked(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source-file")
	target := filepath.Join(dir, "target-link")

	// Create source
	os.WriteFile(source, []byte("hello"), 0644)

	// Create correct symlink
	os.Symlink(source, target)

	mod := Module{
		Name:   "test",
		Source: "source-file",
		Target: target,
	}
	status := CheckLinkStatus(mod, dir)
	if status != StatusLinked {
		t.Errorf("expected StatusLinked, got %v", status)
	}
}

func TestCheckLinkStatus_Conflict(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source-file")
	target := filepath.Join(dir, "target-file")

	os.WriteFile(source, []byte("hello"), 0644)
	os.WriteFile(target, []byte("different"), 0644) // regular file, not symlink

	mod := Module{
		Name:   "test",
		Source: "source-file",
		Target: target,
	}
	status := CheckLinkStatus(mod, dir)
	if status != StatusConflict {
		t.Errorf("expected StatusConflict, got %v", status)
	}
}

func TestCheckLinkStatus_Stale(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source-file")
	wrongSource := filepath.Join(dir, "wrong-source")
	target := filepath.Join(dir, "target-link")

	os.WriteFile(source, []byte("correct"), 0644)
	os.WriteFile(wrongSource, []byte("wrong"), 0644)

	// Symlink points to wrong source
	os.Symlink(wrongSource, target)

	mod := Module{
		Name:   "test",
		Source: "source-file",
		Target: target,
	}
	status := CheckLinkStatus(mod, dir)
	if status != StatusStale {
		t.Errorf("expected StatusStale, got %v", status)
	}
}

func TestLink_NewFile(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "src", "config")
	target := filepath.Join(dir, "dst", "config")

	os.MkdirAll(filepath.Dir(source), 0755)
	os.WriteFile(source, []byte("config content"), 0644)

	mod := Module{
		Name:   "test",
		Source: "src/config",
		Target: target,
	}

	result := Link(mod, dir, false)
	if result.Error != nil {
		t.Fatalf("Link error: %v", result.Error)
	}
	if result.Status != StatusLinked {
		t.Errorf("expected StatusLinked, got %v", result.Status)
	}

	// Verify symlink
	linkDest, err := os.Readlink(target)
	if err != nil {
		t.Fatalf("Readlink error: %v", err)
	}
	if linkDest != source {
		t.Errorf("symlink points to %q, want %q", linkDest, source)
	}
}

func TestLink_WithBackup(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "src", "config")
	target := filepath.Join(dir, "dst", "config")

	os.MkdirAll(filepath.Dir(source), 0755)
	os.MkdirAll(filepath.Dir(target), 0755)
	os.WriteFile(source, []byte("new content"), 0644)
	os.WriteFile(target, []byte("old content"), 0644)

	// Override backup dir for test
	origBackupDir := os.Getenv("XDG_DATA_HOME")
	os.Setenv("XDG_DATA_HOME", filepath.Join(dir, "data"))
	defer os.Setenv("XDG_DATA_HOME", origBackupDir)

	mod := Module{
		Name:   "test",
		Source: "src/config",
		Target: target,
	}

	result := Link(mod, dir, true)
	if result.Error != nil {
		t.Fatalf("Link error: %v", result.Error)
	}
	if result.BackupPath == "" {
		t.Error("expected backup path to be set")
	}

	// Verify backup exists
	if _, err := os.Stat(result.BackupPath); err != nil {
		t.Errorf("backup file not found: %v", err)
	}
}

func TestUnlink(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source")
	target := filepath.Join(dir, "link")

	os.WriteFile(source, []byte("data"), 0644)
	os.Symlink(source, target)

	mod := Module{Target: target}
	err := Unlink(mod)
	if err != nil {
		t.Fatalf("Unlink error: %v", err)
	}

	// Should be gone
	if _, err := os.Lstat(target); !os.IsNotExist(err) {
		t.Error("target should be removed after unlink")
	}
}

func TestUnlink_RefuseNonSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "regular-file")
	os.WriteFile(target, []byte("data"), 0644)

	mod := Module{Target: target}
	err := Unlink(mod)
	if err == nil {
		t.Error("should refuse to unlink a regular file")
	}
}

func TestAppend_WithGuard(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "append-content")
	target := filepath.Join(dir, "target-file")

	os.WriteFile(source, []byte("# Guard Line\nnew content\n"), 0644)
	os.WriteFile(target, []byte("existing content\n"), 0644)

	mod := Module{
		Name:     "test",
		Source:   "append-content",
		Target:   target,
		Strategy: StrategyAppend,
		Guard:    "Guard Line",
	}

	// First append
	result := Append(mod, dir)
	if result.Error != nil {
		t.Fatalf("Append error: %v", result.Error)
	}

	// Second append should be no-op (guard detected)
	result2 := Append(mod, dir)
	if result2.Error != nil {
		t.Fatalf("Second Append error: %v", result2.Error)
	}
	if result2.Status != StatusLinked {
		t.Error("second append should detect guard and return linked status")
	}

	// Verify content was only appended once
	data, _ := os.ReadFile(target)
	content := string(data)
	if count := countOccurrences(content, "new content"); count != 1 {
		t.Errorf("expected content appended once, got %d times", count)
	}
}

func countOccurrences(s, substr string) int {
	count := 0
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			count++
		}
	}
	return count
}
