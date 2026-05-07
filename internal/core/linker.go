package core

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// LinkStatus represents the current state of a link target
type LinkStatus string

const (
	StatusLinked   LinkStatus = "linked"    // our symlink is in place
	StatusMissing  LinkStatus = "missing"   // target doesn't exist
	StatusConflict LinkStatus = "conflict"  // target exists but isn't our symlink
	StatusStale    LinkStatus = "stale"     // symlink points to wrong source
)

// LinkResult contains the outcome of a link operation
type LinkResult struct {
	Module     string
	Status     LinkStatus
	BackupPath string
	Error      error
}

// CheckLinkStatus checks the current state of a module's target
func CheckLinkStatus(module Module, repoRoot string) LinkStatus {
	target := ExpandPath(module.Target)

	// Append-strategy modules: status is determined by guard presence, not symlink
	if module.Strategy == StrategyAppend {
		if module.Guard != "" && containsGuard(target, module.Guard) {
			return StatusLinked
		}
		return StatusMissing
	}

	source := filepath.Join(repoRoot, module.Source)

	info, err := os.Lstat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return StatusMissing
		}
		return StatusConflict
	}

	// Check if it's a symlink
	if info.Mode()&os.ModeSymlink != 0 {
		linkDest, err := os.Readlink(target)
		if err != nil {
			return StatusConflict
		}
		// Resolve to absolute path for comparison
		if !filepath.IsAbs(linkDest) {
			linkDest = filepath.Join(filepath.Dir(target), linkDest)
		}
		absSource, _ := filepath.Abs(source)
		absLink, _ := filepath.Abs(linkDest)
		if absSource == absLink {
			return StatusLinked
		}
		return StatusStale
	}

	// Target exists but is not a symlink → conflict
	return StatusConflict
}

// Link creates a symlink from source to target, with backup if needed
func Link(module Module, repoRoot string, doBackup bool) LinkResult {
	target := ExpandPath(module.Target)
	source := filepath.Join(repoRoot, module.Source)
	result := LinkResult{Module: module.Name}

	// Verify source exists
	if _, err := os.Stat(source); err != nil {
		result.Error = fmt.Errorf("source does not exist: %s", source)
		return result
	}

	// Ensure target parent directory exists
	targetDir := filepath.Dir(target)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		result.Error = fmt.Errorf("creating target directory: %w", err)
		return result
	}

	// Handle existing target
	if _, err := os.Lstat(target); err == nil {
		if doBackup {
			backupPath, err := BackupWithName(module.Name, target)
			if err != nil {
				result.Error = fmt.Errorf("backup failed: %w", err)
				return result
			}
			result.BackupPath = backupPath
		}
		// Remove existing target
		if err := os.RemoveAll(target); err != nil {
			result.Error = fmt.Errorf("removing existing target: %w", err)
			return result
		}
	}

	// Create symlink
	if err := os.Symlink(source, target); err != nil {
		result.Error = fmt.Errorf("creating symlink: %w", err)
		return result
	}

	result.Status = StatusLinked
	return result
}

// Unlink removes a symlink (does not restore backup)
func Unlink(module Module) error {
	target := ExpandPath(module.Target)

	info, err := os.Lstat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // already gone
		}
		return err
	}

	// Safety: only remove if it's a symlink
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("target %s is not a symlink, refusing to remove", target)
	}

	return os.Remove(target)
}

// Append adds content to a file with a guard check
func Append(module Module, repoRoot string) LinkResult {
	result := LinkResult{Module: module.Name}
	target := ExpandPath(module.Target)
	source := filepath.Join(repoRoot, module.Source)

	// Check guard
	if module.Guard != "" {
		if containsGuard(target, module.Guard) {
			result.Status = StatusLinked // already appended
			return result
		}
	}

	// Read source content
	content, err := os.ReadFile(source)
	if err != nil {
		result.Error = fmt.Errorf("reading source: %w", err)
		return result
	}

	// Append to target
	f, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		result.Error = fmt.Errorf("opening target for append: %w", err)
		return result
	}
	defer f.Close()

	if _, err := f.Write(content); err != nil {
		result.Error = fmt.Errorf("appending content: %w", err)
		return result
	}

	result.Status = StatusLinked
	return result
}

// Backup copies a file/dir to the backup directory with module name + timestamp
func Backup(path string) (string, error) {
	return BackupWithName(filepath.Base(path), path)
}

// BackupWithName copies a file/dir using a specific name prefix (typically module name)
func BackupWithName(name, path string) (string, error) {
	backupDir := BackupDir()
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	timestamp := time.Now().Format("20060102-150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("%s.%s", name, timestamp))

	// Use copy instead of move to preserve original until we confirm success
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		if err := copyDir(path, backupPath); err != nil {
			return "", err
		}
	} else {
		if err := copyFile(path, backupPath); err != nil {
			return "", err
		}
	}

	return backupPath, nil
}

// RunPostLink executes the post_link hook script for a module, if set
func RunPostLink(module Module, repoRoot string) error {
	script := module.Hooks.PostLink
	if script == "" {
		return nil
	}
	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("post_link hook failed: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// Restore copies a backup back to its original location
func Restore(backupPath, targetPath string) error {
	info, err := os.Stat(backupPath)
	if err != nil {
		return fmt.Errorf("backup not found: %w", err)
	}

	// Remove current target if exists
	os.RemoveAll(targetPath)

	if info.IsDir() {
		return copyDir(backupPath, targetPath)
	}
	return copyFile(backupPath, targetPath)
}

func containsGuard(filePath, guard string) bool {
	f, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), guard) {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, info.Mode())
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}
