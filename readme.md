# Dotfiles & Dev Environment Setup

Personal development environment configuration for macOS and Linux (including WSL2).

## Principles

1. Easy to use
2. Clear and light
3. Beautify and streamline the system

---

## Terminal Emulators

### Ghostty (`ghostty/`)
- Config: `ghostty/config` → `~/.config/ghostty/config`
- Custom GLSL cursor shaders: `ghostty/shaders/` → `~/.config/ghostty/shaders/`
- Active shader: `cursor_warp.glsl`（其余 shader 可按需替换）
- Fonts: `Inconsolata LGC Nerd Font Mono` + `LXGW WenKai Mono`（CJK fallback）
- SSH integration enabled: `ssh-terminfo,ssh-env,sudo`

### WezTerm (`wezterm/`)
- Entry point: `wezterm/wezterm.lua` → `~/.config/wezterm/wezterm.lua`
- 子模块：`tabbar/`（自定义 tab 栏）、`utils/`（平台检测、GPU 适配等）
- Leader key: `Ctrl+Shift+Space`；所有默认快捷键已禁用，完整键位在 `wezterm.lua` 中定义
- 平台自动适配：Mac 用 `Cmd`，Windows 用 `Alt`；Windows 默认启动 WSL2

---

## Shell & Prompt

### Starship (`starship.toml`)
- Copy to `~/.config/starship.toml`
- Catppuccin Mocha 主题，prompt 显示：OS → 用户名 → 目录 → Git → 语言版本 → conda → 时间 → 命令耗时
- 备用配置：`starship.toml.pil`

### Proxy Functions (`setproxy.sh`)
由 `terminalInstall.sh` 自动追加到 `~/.zshrc`，提供以下命令：

```bash
proxyall [ip] [port]   # 设置 http/https/git/npm 代理（默认 192.168.31.201:7890）
proxyoff               # 取消所有代理
proxyon [port]         # WSL2 专用：自动获取宿主机 IP 并设置代理
gethostip              # 查看并 ping WSL2 宿主机 IP
```

---

## Multiplexers

### Zellij (`zellij/`)
- Config: `zellij/config.kdl` → `~/.config/zellij/config.kdl`
- Layouts: `zellij/layouts/`
- 所有默认快捷键已清除并重新定义；模式切换：`Ctrl+p`（pane）、`Ctrl+t`（tab）、`Ctrl+n`（resize）、`Ctrl+o`（session）、`Ctrl+b`（tmux 兼容模式）

### Tmux (`tmux/`)
- Config: `tmux/.tmux.conf` → `~/.tmux.conf`

---

## Editors

### Neovim / NvChad (`nvchad_custom_file/`)
- 部署顺序：先 clone NvChad 到 `~/.config/nvim`，再将 `nvchad_custom_file/lua/custom/` 复制到 `~/.config/nvim/lua/custom/`
- 包含：个人 keymaps、UI 覆盖、LSP 配置、formatter/linter（null-ls）、自定义 snippets

### Vim (`vim-dot/`)
- 最新配置：`vim-dot/vimrc_aiken` → `~/.vimrc`
- Leader key: `Space`；F5 编译运行，F6 添加文件头

### VSCode (`vscode/`)
- 使用 Profile 管理不同场景：Minimal / CPP / Python / Frontend
- 导入方式：VSCode Profile 页面 → Import Profile，或直接使用 Gist 链接（见 `vscode/readme.md`）
- 批量安装插件：
  ```bash
  cat vscode/extensions.txt | xargs -I{} code --install-extension {}
  ```

---

## File Manager

### Ranger (`ranger/`)
- Config: `ranger/rc.conf` → `~/.config/ranger/rc.conf`

---

## Windows / Cross-platform

### PotPlayer (`potplayer/`)
- 导入方式：双击 `.reg` 文件，或将 `.ini` 放到安装目录（`DAUM\PotPlayer\`）
- 皮肤 `.dsf` 文件放入安装目录下的 `Skins/` 文件夹

### PowerShell (`powershell-dot/`)
- `Microsoft.PowerShell_profile.ps1` → PowerShell `$PROFILE` 路径

---

## Install Scripts

所有脚本从仓库根目录运行，依赖 `FunctionList.sh`（不可直接执行，由其他脚本 source）。

```bash
# Ubuntu：修复默认 shell（dash → bash）
bash InitBashEnv.sh

# 安装 zsh、oh-my-zsh、插件、ranger、neofetch、btm/htop；自动追加 proxy/alias 到 ~/.zshrc
bash terminalInstall.sh

# 安装 neovim 及全部依赖（conda、node/nvm、lazygit、ripgrep、fd、ruby、nvchad）
bash nvimInstall.sh
```
