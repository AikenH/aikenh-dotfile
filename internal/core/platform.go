package core

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Platform represents the detected system environment
type Platform struct {
	OS       string // darwin, linux
	Arch     string // amd64, arm64
	IsWSL2   bool
	HasGUI   bool
	PkgMgr   string // apt, brew, dnf, pacman
	Shell    string // zsh, bash
	HomeDir  string
}

// DetectPlatform gathers system information
func DetectPlatform() Platform {
	home, _ := os.UserHomeDir()
	p := Platform{
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		HomeDir: home,
	}

	// WSL2 detection
	if p.OS == "linux" {
		if data, err := os.ReadFile("/proc/version"); err == nil {
			if strings.Contains(strings.ToLower(string(data)), "microsoft") {
				p.IsWSL2 = true
			}
		}
	}

	// GUI detection
	p.HasGUI = detectGUI(p)

	// Package manager detection
	p.PkgMgr = detectPkgManager(p)

	// Shell detection
	p.Shell = detectShell()

	return p
}

func detectGUI(p Platform) bool {
	if p.OS == "darwin" {
		return true // macOS always has GUI
	}
	// Check DISPLAY or WAYLAND_DISPLAY
	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" {
		return true
	}
	// WSL2 with WSLg
	if p.IsWSL2 {
		if os.Getenv("DISPLAY") != "" {
			return true
		}
	}
	return false
}

func detectPkgManager(p Platform) string {
	if p.OS == "darwin" {
		if _, err := exec.LookPath("brew"); err == nil {
			return "brew"
		}
		return "none"
	}

	// Linux package managers in priority order
	managers := []struct {
		cmd  string
		name string
	}{
		{"apt-get", "apt"},
		{"dnf", "dnf"},
		{"yum", "yum"},
		{"pacman", "pacman"},
		{"zypper", "zypper"},
	}

	for _, m := range managers {
		if _, err := exec.LookPath(m.cmd); err == nil {
			return m.name
		}
	}

	// Also check if brew is available on Linux
	if _, err := exec.LookPath("brew"); err == nil {
		return "brew"
	}

	return "none"
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	if strings.Contains(shell, "bash") {
		return "bash"
	}
	return "bash"
}

// String returns a human-readable summary
func (p Platform) String() string {
	wsl := ""
	if p.IsWSL2 {
		wsl = " (WSL2)"
	}
	gui := "headless"
	if p.HasGUI {
		gui = "GUI"
	}
	return fmt.Sprintf("%s/%s%s | %s | pkg:%s | shell:%s", p.OS, p.Arch, wsl, gui, p.PkgMgr, p.Shell)
}
