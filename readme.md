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

### FZF / FD / Bat (`shell/zshrc.common`, `bat/config`)
- `shell/zshrc.common`：通用 shell 片段，包含 `fzf` 初始化、`fd` 作为默认文件源、`bat` 预览、以及 git helper：`gfb` / `gfc` / `gfs`
- `bat/config` → `~/.config/bat/config`
- 建议在本机 `~/.zshrc` 中显式 source：
  ```bash
  [ -f ~/Workspace/aikenh-dotfile/shell/zshrc.common ] && source ~/Workspace/aikenh-dotfile/shell/zshrc.common
  ```
- 当前行为：
  - `Ctrl-T`：模糊选择文件并将路径插入命令行
  - `Alt-C`：模糊选择目录并切换到该目录
  - `gfb`：模糊选择 branch 并切换
  - `gfc`：模糊浏览 commit，并用 `delta` 预览 `git show`
  - `gfs`：模糊浏览已修改/已暂存/未跟踪文件；预览 diff 时用 `delta`，预览新文件时用 `bat`

### Git / Delta (`git/.gitconfig`)
- `git/.gitconfig`：只放通用显示增强，不包含身份信息
- 建议在本机 `~/.gitconfig` 中 include：
  ```gitconfig
  [include]
      path = ~/Workspace/aikenh-dotfile/git/.gitconfig
  ```
- 作用：`git diff` / `git log` / `git reflog` / `git show` 默认使用 `delta`

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

---

## Roadmap

### Current strategy
- 这次新增的 `fzf` / `fd` / `bat` / `delta` 通用配置已沉淀到仓库中，供新机器复用
- 新机器优先按照本文档进行 `source` / `include` / 软链接

### Future direction
- `terminalInstall.sh` 当前更适合作为历史脚本维护，不再继续堆叠新逻辑
- 后续考虑重做一个模块化安装器：按 shell、git、bat、terminal、editor 等模块组织
- 安装器优先提供轻量 TUI 选择体验，风格参考 `fzf` / `lazygit`：先选择要启用的模块，再执行安装、链接和配置写入
- 模块设计目标：
  - 安装依赖
  - 链接或复制配置
  - 追加必要的 source/include 语句
  - 支持 macOS / Linux 的差异化处理
- 长期方向是把“通用配置”和“机器私有配置”明确分层，减少一次性脚本对当前环境的直接修改
