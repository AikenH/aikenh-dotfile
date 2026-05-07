package core

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/aikenhong/dotsetup/internal/executor"
)

// InstallStatus represents the current state of a tool
type InstallStatus string

const (
	InstallStatusInstalled   InstallStatus = "installed"
	InstallStatusMissing     InstallStatus = "missing"
	InstallStatusOutdated    InstallStatus = "outdated"
)

// InstallResult holds the outcome of an install operation
type InstallResult struct {
	Module   string
	Status   InstallStatus
	Version  string
	Output   string
	Error    error
}

// Installer handles tool installation
type Installer struct {
	runner     *executor.Runner
	platform   Platform
	aptUpdated bool // track whether apt-get update has been run this session
}

// NewInstaller creates an installer for the current platform
func NewInstaller(proxy string) *Installer {
	return &Installer{
		runner:   executor.NewRunner(proxy),
		platform: DetectPlatform(),
	}
}

// SetLogFunc sets the real-time log callback
func (inst *Installer) SetLogFunc(fn func(string)) {
	inst.runner.LogFunc = fn
}

// ensureAptUpdated runs apt-get update once per installer session
func (inst *Installer) ensureAptUpdated() {
	if inst.aptUpdated {
		return
	}
	// Check if apt lists exist; if not, we need to update
	checkResult := inst.runner.Run("test -d /var/lib/apt/lists/partial && ls /var/lib/apt/lists/ | grep -q '^[^.]'")
	if checkResult.ExitCode != 0 {
		inst.runner.Run("sudo apt-get update -qq")
	}
	inst.aptUpdated = true
}

// CheckInstalled returns whether a tool is already installed
func (inst *Installer) CheckInstalled(mod Module) bool {
	if mod.Install.Check == "" {
		return false
	}
	return inst.runner.Check(mod.Install.Check)
}

// GetVersion detects the current version of a tool
func (inst *Installer) GetVersion(mod Module) string {
	if mod.Version.Command == "" {
		return ""
	}
	result := inst.runner.Run(mod.Version.Command)
	if result.ExitCode != 0 {
		return ""
	}

	output := strings.TrimSpace(result.Stdout)
	if mod.Version.Pattern == "" {
		return output
	}

	re, err := regexp.Compile(mod.Version.Pattern)
	if err != nil {
		return output
	}

	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		return matches[1]
	}
	return output
}

// Install installs a tool module
func (inst *Installer) Install(mod Module) InstallResult {
	result := InstallResult{Module: mod.Name}

	// Idempotent: check if already installed
	if inst.CheckInstalled(mod) {
		result.Status = InstallStatusInstalled
		result.Version = inst.GetVersion(mod)
		return result
	}

	// Pick install method based on platform
	var installCmd string
	var err error

	switch {
	case inst.platform.OS == "darwin" && mod.Install.Brew != "":
		installCmd = fmt.Sprintf("brew install %s", mod.Install.Brew)

	case inst.platform.PkgMgr == "apt" && mod.Install.Apt != "":
		// Ensure apt cache is available
		inst.ensureAptUpdated()
		// Add PPA if specified
		if mod.Install.PPA != "" {
			ppaCheck := inst.runner.Run(fmt.Sprintf("grep -r %s /etc/apt/sources.list.d/ 2>/dev/null", mod.Install.PPA))
			if ppaCheck.ExitCode != 0 {
				inst.runner.Run("sudo apt-get install -y software-properties-common")
				addResult := inst.runner.Run(fmt.Sprintf("sudo add-apt-repository -y ppa:%s", mod.Install.PPA))
				if addResult.ExitCode != 0 {
					// PPA failed, try without it
				} else {
					inst.runner.Run("sudo apt-get update")
				}
			}
		}
		installCmd = fmt.Sprintf("sudo apt-get install -y %s", mod.Install.Apt)

	case mod.Install.GithubRelease != "":
		installCmd, err = inst.buildGithubReleaseCmd(mod)
		if err != nil {
			result.Error = err
			return result
		}

	case mod.Install.Script != "":
		runResult := inst.runner.RunScript(mod.Install.Script)
		result.Output = runResult.Stdout + runResult.Stderr
		if runResult.Err != nil {
			result.Error = fmt.Errorf("script install failed: %w\n%s", runResult.Err, runResult.Stderr)
			result.Status = InstallStatusMissing
			return result
		}
		result.Status = InstallStatusInstalled
		result.Version = inst.GetVersion(mod)
		return result

	default:
		result.Error = fmt.Errorf("no install method available for %s on %s/%s",
			mod.Name, inst.platform.OS, inst.platform.PkgMgr)
		result.Status = InstallStatusMissing
		return result
	}

	// Execute install command
	runResult := inst.runner.Run(installCmd)
	result.Output = runResult.Stdout + runResult.Stderr

	if runResult.Err != nil {
		result.Error = fmt.Errorf("install failed: %w\n%s", runResult.Err, runResult.Stderr)
		result.Status = InstallStatusMissing
		return result
	}

	// Verify installation
	if inst.CheckInstalled(mod) {
		result.Status = InstallStatusInstalled
		result.Version = inst.GetVersion(mod)
	} else {
		result.Error = fmt.Errorf("install command succeeded but check still fails")
		result.Status = InstallStatusMissing
	}

	return result
}

func (inst *Installer) buildGithubReleaseCmd(mod Module) (string, error) {
	repo := mod.Install.GithubRelease

	// Get latest version from GitHub API
	versionCmd := fmt.Sprintf(`curl -s "https://api.github.com/repos/%s/releases/latest" | grep -Po '"tag_name": "v?\K[^"]*'`, repo)
	versionResult := inst.runner.Run(versionCmd)
	if versionResult.ExitCode != 0 || strings.TrimSpace(versionResult.Stdout) == "" {
		return "", fmt.Errorf("could not fetch latest version for %s", repo)
	}
	version := strings.TrimSpace(versionResult.Stdout)

	// Resolve asset pattern
	osName := runtime.GOOS
	arch := runtime.GOARCH
	// Normalize naming conventions
	osMap := map[string]string{"darwin": "Darwin", "linux": "Linux"}
	archMap := map[string]string{"amd64": "x86_64", "arm64": "arm64"}

	if mapped, ok := osMap[osName]; ok {
		osName = mapped
	}
	if mapped, ok := archMap[arch]; ok {
		arch = mapped
	}

	asset := mod.Install.AssetPattern
	asset = strings.ReplaceAll(asset, "{version}", version)
	asset = strings.ReplaceAll(asset, "{os}", osName)
	asset = strings.ReplaceAll(asset, "{arch}", arch)

	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/v%s/%s", repo, version, asset)

	// Build download + install script
	script := fmt.Sprintf(`
set -e
cd /tmp
curl -fsSL "%s" -o "%s"
`, downloadURL, asset)

	if strings.HasSuffix(asset, ".tar.gz") {
		binName := strings.SplitN(repo, "/", 2)[1]
		script += fmt.Sprintf(`
tar xf "%s" %s 2>/dev/null || tar xf "%s"
sudo install %s /usr/local/bin/
rm -f "%s" %s
`, asset, binName, asset, binName, asset, binName)
	} else if strings.HasSuffix(asset, ".deb") {
		script += fmt.Sprintf(`sudo dpkg -i "%s" && rm -f "%s"`, asset, asset)
	}

	return script, nil
}
