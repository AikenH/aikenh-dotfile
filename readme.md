# dotsetup

Personal dotfiles + automated dev environment setup for macOS and Linux (including WSL2).

A TUI-driven tool that manages dotfile symlinking, shell init injection, and binary tool installation — all declarative, idempotent, and platform-aware.

---

## Quick Start

```bash
# Clone the repo
git clone https://github.com/aikenhong/aikenh-dotfile ~/.dotfiles
cd ~/.dotfiles

# Run the TUI (pre-built binary)
./dotsetup
```

First run lands on the **Profile Selection** screen. Pick a profile or press `Esc` to configure manually.

---

## Usage

```
./dotsetup
```

### Navigation

| Screen | How to reach | Keys |
|---|---|---|
| Home | default | `q` quit, `↑↓` navigate, `Enter` select |
| Link Configs | Home → Link Configs | `Space` toggle, `Enter` apply, `a` all, `n` none, `q` back |
| Install Tools | Home → Install Tools | `Space` toggle, `a` all uninstalled, `Enter` proceed, `q` back |
| Status | Home → Status | `q` / `Esc` back |
| Settings | Home → Settings | proxy URL edit, `Enter` save |
| Profile | first launch | `↑↓` navigate, `Enter` select, `Esc` skip |

### Link Configs view

Manages config files via **symlink** (e.g. `dotfiles/ghostty` → `~/.config/ghostty`) or **append** (e.g. shell init blocks appended to `~/.zshrc` with guard comments).

- Symlink modules: toggling applies / removes the symlink; existing files are backed up automatically to `~/.local/share/dotsetup/backups/`
- Append modules: toggling inserts / does nothing (append is one-way; manual removal from `~/.zshrc` required)
- Group modules (e.g. `starship-style`): radio-select — only one member can be active at a time

### Install Tools view

Installs binary tools via the system package manager (apt / brew / dnf / pacman) with GitHub Release fallback. Installation runs sequentially with a real-time progress bar and streaming log output.

### State file

All applied changes are tracked at `~/.local/share/dotsetup/state.json`. This file records installed modules, versions, backup paths, and group choices. The TUI reads it on start to reflect current status accurately.

---

## Build

**Prerequisites:** Go 1.21+

```bash
# Build for current platform
make build          # → ./dotsetup

# Cross-compile all targets
make build-all      # → dist/dotsetup_{os}_{arch}

# Build + run immediately
make run
```

---

## Testing

### Unit tests

```bash
go test ./internal/core/... -v
```

All 24 tests cover: module loading & platform filtering, symlink status detection, link/unlink/append operations, dependency graph resolution, and state persistence.

### Docker integration test

Test against a clean Ubuntu 22.04 environment (simulates a fresh server):

```bash
# 1. Rebuild binary first (so the image gets the latest version)
go build -o dotsetup ./cmd/dotsetup

# 2. Build Docker image
docker build -t dotsetup-test .

# 3. Launch interactive TUI test
docker run --rm -it dotsetup-test ./dotsetup

# 4. Drop into shell for manual inspection
docker run --rm -it dotsetup-test bash
```

The container runs as `tester` (non-root with sudo), with the full repo copied in. The `modules/` and `dotfiles/` directories are present exactly as in the repo.

---

## Repository Layout

```
.
├── cmd/dotsetup/        # main entry point
├── internal/
│   ├── core/            # module loading, linker, installer, state, profiles
│   └── executor/        # shell runner (proxy, timeout, log streaming)
├── modules/
│   ├── configs.yaml     # config modules (symlink / append)
│   └── tools.yaml       # tool modules (install)
├── profiles/
│   ├── essential.yaml   # headless server preset
│   └── full.yaml        # full workstation preset
├── dotfiles/            # all managed config files
│   ├── ghostty/         # Ghostty terminal config + GLSL shaders
│   ├── wezterm/         # WezTerm Lua config
│   ├── kaku/            # Kaku editor config
│   ├── zellij/          # Zellij multiplexer config + layouts
│   ├── tmux/            # Tmux config
│   ├── vim/             # Vim config (vimrc_aiken)
│   ├── nvim/            # NvChad custom config (lua/custom/)
│   ├── shell/           # Shell init snippets (starship, fzf, git-fzf, aliases)
│   ├── git/             # Git delta config
│   ├── ranger/          # Ranger file manager config
│   ├── tool-configs/    # Tool-specific shell integrations (bat, fzf, delta)
│   ├── claude/          # Claude Code settings
│   ├── starship.toml    # Starship Catppuccin powerline style
│   └── starship.toml.pil # Starship minimal style (default profile)
├── platform/
│   └── windows/         # Windows-only configs (not managed by dotsetup)
│       ├── potplayer/   # PotPlayer INI + skins
│       ├── powershell/  # PowerShell profile
│       └── vscode/      # VSCode profiles + extensions list
├── archive/             # Legacy install scripts (no longer used)
└── Dockerfile           # Ubuntu 22.04 test environment
```

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                    TUI (Bubbletea)                   │
│  HomeModel → ModulesModel / InstallModel / ...       │
│               App (root model, routes msgs)          │
└──────────────┬──────────────────────┬────────────────┘
               │                      │
        ┌──────▼──────┐        ┌──────▼──────┐
        │  core/linker│        │core/installer│
        │  Link()     │        │  Install()   │
        │  Unlink()   │        │  CheckInstalled()│
        │  Append()   │        └──────┬──────┘
        └──────┬──────┘               │
               │                ┌─────▼──────┐
        ┌──────▼──────┐         │ executor/  │
        │ core/module │         │ runner.go  │
        │  LoadAll()  │         │ Run/Script │
        │  Filter()   │         └────────────┘
        └──────┬──────┘
               │
        ┌──────▼──────┐
        │  core/state │
        │  LoadState()│
        │  Save()     │
        └─────────────┘
```

### Key design decisions

- **Declarative YAML** — all modules are defined in `modules/*.yaml`; the binary contains no hardcoded tool names
- **Idempotent** — every operation checks current state before acting; safe to run multiple times
- **Platform-aware** — each module declares `platforms: [darwin, linux]`; the loader filters at startup
- **Dependency graph** — `deps.go` implements topological sort; Install view auto-selects uninstalled deps
- **Guard-based append** — shell init snippets use a unique guard comment to prevent duplicate appends
- **Backup before overwrite** — any conflicting target is backed up with a timestamp before being replaced
- **Real-time streaming** — installer uses a buffered `chan string` bridged into Bubbletea's Update loop via `drainLogs()` + `tea.Tick`

---

## Adding a New Config Module

### 1. Add the config file

Put the config file or directory under `dotfiles/`:

```
dotfiles/
└── myapp/
    └── config.toml
```

### 2. Register it in `modules/configs.yaml`

```yaml
- name: myapp
  type: config
  description: "MyApp configuration"
  source: dotfiles/myapp/config.toml     # relative to repo root
  target: ~/.config/myapp/config.toml   # where it lands on the system
  strategy: symlink                      # symlink | copy | append
  link_mode: file                        # file | directory
  platforms: [darwin, linux]
  group: ""                              # leave empty unless mutually exclusive with another module
  optional: true
  tags: [myapp]
  hooks:
    post_link: ""                        # optional bash snippet run after linking
```

**strategy options:**

| Value | Effect |
|---|---|
| `symlink` | Creates a symlink at `target` pointing to `source` |
| `copy` | Copies the file (not tracked after copy) |
| `append` | Appends source content to `target`; use `guard:` to prevent duplicates |

**For directory configs** (e.g. an entire `~/.config/myapp/`):

```yaml
source: dotfiles/myapp          # the directory itself
target: ~/.config/myapp
strategy: symlink
link_mode: directory
```

**For shell init snippets** (appended to `~/.zshrc`):

```yaml
source: dotfiles/shell/myapp.sh
target: ~/.zshrc
strategy: append
link_mode: file
guard: "# MyApp Init Flag"     # prevents duplicate appends
```

### 3. Test it

```bash
./dotsetup
# → Link Configs → find your module → Space to toggle → Enter to apply
```

Or verify with the status view (Home → Status) to see symlink state.

---

## Adding a New Tool Module

### 1. Register it in `modules/tools.yaml`

```yaml
- name: mytool
  type: tool
  description: "My tool description"
  platforms: [darwin, linux]
  tags: [dev-tools]
  depends_on: []                   # other tool names that must be installed first
  install:
    check: "mytool --version"      # command that exits 0 if already installed
    apt: "mytool"                  # package name for apt
    brew: "mytool"                 # package name for homebrew
    dnf: "mytool"                  # package name for dnf
    pacman: "mytool"               # package name for pacman
    github_release: "owner/repo"   # GitHub release fallback (optional)
    asset_pattern: "mytool-{version}-{os_lower}-{arch_raw}.tar.gz"
    script: |                      # custom install script (last resort)
      curl -fsSL https://mytool.sh/install.sh | bash
  version:
    command: "mytool --version"
    pattern: 'mytool (\d+\.\d+\.\d+)'   # capture group 1 = version string
```

**Install method priority:** package manager → github_release → script (first that succeeds wins).

**Asset pattern variables:**

| Variable | Example |
|---|---|
| `{version}` | `0.44.1` |
| `{os}` | `Linux` / `Darwin` |
| `{os_lower}` | `linux` / `darwin` |
| `{arch}` | `x86_64` / `arm64` |
| `{arch_raw}` | `amd64` / `arm64` |

### 2. Test it

```bash
./dotsetup
# → Install Tools → find your tool → Space → Enter → y to confirm
```

Or in Docker for a clean-room test:

```bash
go build -o dotsetup ./cmd/dotsetup
docker build -t dotsetup-test .
docker run --rm -it dotsetup-test ./dotsetup
```

---

## Adding a New Profile

Profiles live in `profiles/`. Create a new YAML file:

```yaml
# profiles/server.yaml
name: server
description: "Minimal headless server setup"

modules:
  - zsh
  - starship-bin
  - fzf
  - ripgrep
  - git-delta

group_choices:
  starship-style: starship          # pick one module from the group (or "none" / "prompt")
```

The profile appears automatically in the Profile Selection screen on first run.

---

## Dotfiles Reference

| Config | Source | Target |
|---|---|---|
| Ghostty | `dotfiles/ghostty/` | `~/.config/ghostty/` |
| WezTerm | `dotfiles/wezterm/` | `~/.config/wezterm/` |
| Zellij | `dotfiles/zellij/` | `~/.config/zellij/` |
| Vim | `dotfiles/vim/vimrc_aiken` | `~/.vimrc` |
| NvChad | `dotfiles/nvim/lua/custom/` | `~/.config/nvim/lua/custom/` |
| Starship (minimal) | `dotfiles/starship.toml.pil` | `~/.config/starship.toml` |
| Starship (powerline) | `dotfiles/starship.toml` | `~/.config/starship.toml` |
| Git delta | `dotfiles/git/gitconfig-delta` | `~/.config/git/delta.gitconfig` |
| Ranger | `dotfiles/ranger/rc.conf` | `~/.config/ranger/rc.conf` |
| Claude Code | `dotfiles/claude/settings.json` | `~/.claude/settings.json` |
| Shell aliases | `dotfiles/shell/aliases.sh` | `~/.zshrc` (append) |
| FZF integration | `dotfiles/shell/fzf.sh` | `~/.zshrc` (append) |
| Git-FZF functions | `dotfiles/shell/git-fzf.sh` | `~/.zshrc` (append) |
| Starship init | `dotfiles/shell/starship.sh` | `~/.zshrc` (append) |
| Proxy functions | `dotfiles/setproxy.sh` | `~/.zshrc` (append) |

---

## Windows / Platform-specific

Windows configs are in `platform/windows/` and are **not managed by dotsetup** (manual setup):

- **VSCode** — import `.code-profile` files via VSCode Profile page; bulk install extensions: `cat platform/windows/vscode/extensions.txt | xargs -I{} code --install-extension {}`
- **PowerShell** — copy `platform/windows/powershell/Microsoft.PowerShell_profile.ps1` to `$PROFILE`
- **PotPlayer** — double-click `.reg`, copy `.ini` to `DAUM\PotPlayer\`, copy `.dsf` skins to `Skins\`
