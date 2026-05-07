# dotsetup · 开发现状文档

> 最后更新: 2026-05-05 · 分支: `tui` · 基础版本: v0.2.0

## 项目概述

dotsetup 是一个 Go + Bubble Tea 构建的 TUI 环境管理器，用于管理个人 dotfile 的 symlink、工具安装和配置同步。

**仓库路径**: `aikenh-dotfile/` (同一个 dotfile 仓库内)  
**二进制**: `dotsetup` (~4.9MB, 单文件)  
**代码量**: ~3900 行 Go (含 24 个单元测试)

---

## 核心架构

```
cmd/dotsetup/main.go         ← CLI 入口 (status/link/rollback/sync-back + TUI)
internal/tui/                 ← Bubble Tea TUI 层
internal/core/                ← 业务逻辑 (module/linker/deps/state/installer/profile)
internal/executor/            ← Shell 执行器 (proxy inject, timeout)
modules/configs.yaml          ← Config 模块注册表
modules/tools.yaml            ← Tool 模块注册表
profiles/essential|full.yaml  ← 预设 Profile
```

### 数据流

```
modules/*.yaml → Core 解析 → State 读写 (~/.local/share/dotsetup/) → Executor → 系统
```

### 状态存储

- **State 文件**: `~/.local/share/dotsetup/state.json`
- **备份目录**: `~/.local/share/dotsetup/backups/`
- 备份命名: `{模块名}.{YYYYMMDD-HHMMSS}`

---

## 模块系统

### 三种类型

| 类型 | 部署方式 | 示例 |
|------|---------|------|
| config | symlink / append | ghostty, vim, starship, zellij |
| tool | apt / brew / script / github_release | neovim, fzf, bat, volta |
| plugin | post_hook 驱动 | vim-commentary |

### 互斥组 (Group)

仅当多个模块指向 **同一个 target 路径** 时才需要 group。当前只有一个：

- `starship-style`: starship (minimal) | starship-powerline (catppuccin)

ghostty/wezterm/kaku/zellij/tmux 等虽然是同类工具，但 target 路径不同，**不设 group**，TUI 中显示为普通 checkbox。

### 依赖关系

在 YAML 中用 `depends_on` 声明。Install TUI 在确认安装时调用 `TopoSortSubset` 自动解析并加入未安装的依赖。

---

## CLI 命令

```bash
dotsetup              # 启动 TUI
dotsetup status       # 展示所有模块状态 (实时检测)
dotsetup link <mod>   # 非交互式 link 指定模块 (支持 post_hook)
dotsetup rollback <mod>  # 反向: unlink + 恢复 backup
dotsetup sync-back    # 把机器当前配置拷回仓库 (文件+目录)
dotsetup help         # 帮助
dotsetup version      # 版本
```

所有命令支持 `--dry-run` / `-n` 标志。

---

## TUI 视图

| View | 功能 | 入口 |
|------|------|------|
| Profile | 首次运行自动弹出，选 essential/full | state.json 不存在时 |
| Home | 主菜单 5 项 | 默认 |
| Modules | Config link 管理 (checkbox + radio for groups) | Home → Link Configs |
| Install | Tool 安装 (已装灰色不可选, deps 自动解析) | Home → Install Tools |
| Settings | 代理地址配置 | Home → Settings |
| Status | 概览 (linked/missing/conflict 统计) | Home → Status |

---

## 六大场景完成状态

| 场景 | 流程 | 状态 |
|------|------|------|
| A. 初始化新环境 | bootstrap.sh → TUI Profile → Install → Link | ✅ |
| B. 增量安装 | TUI 选择 → deps resolve → install | ✅ |
| C. 更新配置 | git pull → symlink 自动生效 | ✅ (设计正确) |
| D. 冲突处理 | detect conflict → backup → link | ✅ |
| E. 反向同步 | sync-back (文件+目录递归) → git diff → commit | ✅ |
| F. 删除/回滚 | rollback → unlink → restore backup | ✅ |

---

## 当前未 commit 的改动

```
M  internal/core/linker.go    ← 备份命名改用模块名前缀
M  internal/tui/app.go        ← Profile 首次运行集成
M  internal/tui/install.go    ← Install 接入依赖图
M  internal/tui/modules.go    ← 移除假互斥组
M  modules/configs.yaml       ← 删除 terminal-emulator/multiplexer group
?? Dockerfile                 ← Docker 测试环境
?? docs/architecture.html     ← Kami 风格架构说明图
```

建议 commit message: `fix: remove false exclusive groups, improve backup naming, add Dockerfile`

---

## 后续可做 (非阻塞, 按需)

| 方向 | 说明 | 复杂度 |
|------|------|--------|
| Docker 冒烟测试 | `docker build && docker run` 验证全新 Ubuntu 流程 | 低 |
| zshrc 模块化 | 把 fzf 配置块、alias、git-fzf 函数拆成独立 append 模块 | 中 |
| yazi 配置入库 | ~/.config/yazi/ 目前不在仓库中 | 低 |
| `dotsetup update` | 对 append 类型模块做增量检测 | 中 |
| Tool uninstall | `brew uninstall` / rm binary | 中 |
| Profile 中的 group_choices=prompt 实现 | TUI 中首次弹出 radio 选择 starship style | 低 |

---

## 构建 & 测试

```bash
# 开发
make build              # 编译当前平台
make run                # 编译 + 启动 TUI
go test ./... -v        # 24 个单元测试
go vet ./...            # 静态分析

# 交叉编译
make build-all          # dist/ 下生成 4 个平台二进制

# Docker 测试 (新机器模拟)
GOOS=linux GOARCH=amd64 go build -o dotsetup ./cmd/dotsetup/
docker build -t dotsetup-test .
docker run -it dotsetup-test

# 发布
git tag v0.3.0
git push --tags         # 触发 .github/workflows/release.yaml
```

---

## 关键设计决策记录

1. **Symlink 而非 Copy** — 仓库即 source of truth, git pull 即更新
2. **YAML 注册表** — 新增模块只改配置, 无需改代码
3. **互斥组仅用于同 target** — ghostty/wezterm 等不互斥
4. **备份按模块名** — `~/.local/share/dotsetup/backups/{name}.{timestamp}`
5. **Tool 检测用 live binary check** — 不依赖 state.json, 真实反映机器状态
6. **单二进制分发** — bootstrap.sh 检测 Go → 编译 / 下载预编译
7. **升级检测已砍** — symlink 随 git pull 更新, tool 版本交给系统包管理
