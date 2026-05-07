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

	// Try each install method in priority order, falling back on failure
	type method struct {
		name string
		try  func() (bool, string, error) // returns (attempted, output, err)
	}

	methods := []method{
		// 1. Native package manager
		{
			name: "pkg-manager",
			try: func() (bool, string, error) {
				cmd, err := inst.buildPkgManagerCmd(mod)
				if err != nil || cmd == "" {
					return false, "", nil // not applicable
				}
				r := inst.runner.Run(cmd)
				return true, r.Stdout + r.Stderr, r.Err
			},
		},
		// 2. GitHub release download
		{
			name: "github-release",
			try: func() (bool, string, error) {
				if mod.Install.GithubRelease == "" {
					return false, "", nil
				}
				cmd, err := inst.buildGithubReleaseCmd(mod)
				if err != nil {
					return true, "", err
				}
				r := inst.runner.Run(cmd)
				return true, r.Stdout + r.Stderr, r.Err
			},
		},
		// 3. Custom script
		{
			name: "script",
			try: func() (bool, string, error) {
				if mod.Install.Script == "" {
					return false, "", nil
				}
				r := inst.runner.RunScript(mod.Install.Script)
				return true, r.Stdout + r.Stderr, r.Err
			},
		},
	}

	var triedMethods []string
	for _, m := range methods {
		attempted, output, err := m.try()
		if !attempted {
			continue
		}
		result.Output += output
		triedMethods = append(triedMethods, m.name)
		if err != nil {
			result.Output += fmt.Sprintf("\n[%s failed: %v, trying next method...]\n", m.name, err)
			continue
		}
		if inst.CheckInstalled(mod) {
			result.Status = InstallStatusInstalled
			result.Version = inst.GetVersion(mod)
			return result
		}
		result.Output += fmt.Sprintf("\n[%s command succeeded but check still fails, trying next method...]\n", m.name)
	}

	if len(triedMethods) == 0 {
		result.Error = fmt.Errorf("no install method available for %s on %s/%s",
			mod.Name, inst.platform.OS, inst.platform.PkgMgr)
	} else {
		result.Error = fmt.Errorf("all install methods failed for %s (tried: %s)",
			mod.Name, strings.Join(triedMethods, ", "))
	}
	result.Status = InstallStatusMissing
	return result
}

// buildPkgManagerCmd returns the native package manager install command, or "" if not applicable
func (inst *Installer) buildPkgManagerCmd(mod Module) (string, error) {
	switch {
	case inst.platform.OS == "darwin" && mod.Install.Brew != "":
		return fmt.Sprintf("brew install %s", mod.Install.Brew), nil

	case inst.platform.PkgMgr == "apt" && mod.Install.Apt != "":
		inst.ensureAptUpdated()
		if mod.Install.PPA != "" {
			ppaCheck := inst.runner.Run(fmt.Sprintf("grep -r %s /etc/apt/sources.list.d/ 2>/dev/null", mod.Install.PPA))
			if ppaCheck.ExitCode != 0 {
				inst.runner.Run("sudo apt-get install -y software-properties-common")
				addResult := inst.runner.Run(fmt.Sprintf("sudo add-apt-repository -y ppa:%s", mod.Install.PPA))
				if addResult.ExitCode == 0 {
					inst.runner.Run("sudo apt-get update")
				}
			}
		}
		return fmt.Sprintf("sudo apt-get install -y %s", mod.Install.Apt), nil

	case inst.platform.PkgMgr == "dnf" && mod.Install.Dnf != "":
		return fmt.Sprintf("sudo dnf install -y %s", mod.Install.Dnf), nil

	case inst.platform.PkgMgr == "pacman" && mod.Install.Pacman != "":
		return fmt.Sprintf("sudo pacman -S --noconfirm %s", mod.Install.Pacman), nil

	case inst.platform.PkgMgr == "brew" && mod.Install.Brew != "":
		return fmt.Sprintf("brew install %s", mod.Install.Brew), nil
	}

	return "", nil
}

func (inst *Installer) buildGithubReleaseCmd(mod Module) (string, error) {
	repo := mod.Install.GithubRelease

	// Get latest version via redirect (avoids GitHub API rate limits)
	// Capture the raw tag (may or may not have "v" prefix, e.g. "v0.44.1" vs "0.19.2")
	versionCmd := fmt.Sprintf(
		`curl -sL -o /dev/null -w "%%{url_effective}" "https://github.com/%s/releases/latest" | grep -oP '/releases/tag/\K[^/]+'`,
		repo,
	)
	versionResult := inst.runner.Run(versionCmd)
	if versionResult.ExitCode != 0 || strings.TrimSpace(versionResult.Stdout) == "" {
		return "", fmt.Errorf("could not fetch latest version for %s", repo)
	}
	tag := strings.TrimSpace(versionResult.Stdout)              // raw tag, e.g. "v0.44.1" or "0.19.2"
	version := strings.TrimPrefix(tag, "v")                    // strip v for asset filenames

	// Resolve asset pattern
	osRaw := runtime.GOOS    // "linux", "darwin"
	archRaw := runtime.GOARCH // "amd64", "arm64"

	// Capitalized variants used by some projects (lazygit, delta, yazi)
	osMap := map[string]string{"darwin": "Darwin", "linux": "Linux"}
	archMap := map[string]string{"amd64": "x86_64", "arm64": "arm64"}
	osName := osRaw
	arch := archRaw
	if mapped, ok := osMap[osName]; ok {
		osName = mapped
	}
	if mapped, ok := archMap[arch]; ok {
		arch = mapped
	}

	asset := mod.Install.AssetPattern
	asset = strings.ReplaceAll(asset, "{version}", version)
	asset = strings.ReplaceAll(asset, "{os}", osName)        // "Linux" / "Darwin"
	asset = strings.ReplaceAll(asset, "{os_lower}", osRaw)   // "linux" / "darwin"
	asset = strings.ReplaceAll(asset, "{arch}", arch)        // "x86_64" / "arm64"
	asset = strings.ReplaceAll(asset, "{arch_raw}", archRaw) // "amd64" / "arm64"

	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, tag, asset)

	script := fmt.Sprintf(`
set -e
cd /tmp
curl -fsSL "%s" -o "%s"
`, downloadURL, asset)

	if strings.HasSuffix(asset, ".tar.gz") {
		binName := strings.SplitN(repo, "/", 2)[1]
		script += fmt.Sprintf(`
TMPDIR=$(mktemp -d)
tar xf "%s" -C "$TMPDIR"
BIN=$(find "$TMPDIR" -name "%s" -type f | head -1)
if [ -z "$BIN" ]; then BIN="$TMPDIR/%s"; fi
sudo install "$BIN" /usr/local/bin/%s
rm -rf "$TMPDIR" "%s"
`, asset, binName, binName, binName, asset)
	} else if strings.HasSuffix(asset, ".zip") {
		binName := strings.SplitN(repo, "/", 2)[1]
		script += fmt.Sprintf(`
TMPDIR=$(mktemp -d)
unzip -q "%s" -d "$TMPDIR"
BIN=$(find "$TMPDIR" -name "%s" -type f | head -1)
sudo install "$BIN" /usr/local/bin/%s
rm -rf "$TMPDIR" "%s"
`, asset, binName, binName, asset)
	} else if strings.HasSuffix(asset, ".deb") {
		script += fmt.Sprintf(`sudo dpkg -i "%s" && rm -f "%s"`, asset, asset)
	}

	return script, nil
}
