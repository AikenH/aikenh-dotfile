package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/aikenhong/dotsetup/internal/core"
)

// InstallPhase tracks where we are in the install flow
type InstallPhase int

const (
	PhaseSelect  InstallPhase = iota // selecting modules
	PhaseConfirm                     // confirming install plan
	PhaseRun                         // executing installs
	PhaseDone                        // finished
)

// InstallItem represents a tool in the install list
type InstallItem struct {
	module   core.Module
	selected bool
	status   core.InstallStatus
	version  string
	log      string
}

// InstallModel handles the tool installation view
type InstallModel struct {
	items     []InstallItem
	cursor    int
	phase     InstallPhase
	progress  int    // index of currently installing module
	message   string
	logs      []string
	installer *core.Installer
	width     int
	height    int
}

func NewInstallModel(modules []core.Module, state *core.State) InstallModel {
	m := InstallModel{
		phase: PhaseSelect,
	}

	// Create installer to do live checks (not just state-based)
	installer := core.NewInstaller(state.Proxy)

	// Filter to only tool-type modules
	for _, mod := range modules {
		if mod.Type != core.TypeTool {
			continue
		}
		item := InstallItem{
			module: mod,
			status: core.InstallStatusMissing,
		}
		// Live check if already installed on this machine
		if installer.CheckInstalled(mod) {
			item.status = core.InstallStatusInstalled
			item.version = installer.GetVersion(mod)
		} else if state.IsInstalled(mod.Name) {
			// State says installed but binary check failed → might be outdated/removed
			item.status = core.InstallStatusMissing
		}
		m.items = append(m.items, item)
	}

	return m
}

// installDoneMsg signals that one install finished
type installDoneMsg struct {
	index  int
	result core.InstallResult
}

func (m InstallModel) Update(msg tea.Msg, app *App) (InstallModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.phase == PhaseSelect {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.items)-1 {
					m.cursor++
				}
			case " ", "space":
				// Only allow toggling uninstalled tools
				if m.items[m.cursor].status != core.InstallStatusInstalled {
					m.items[m.cursor].selected = !m.items[m.cursor].selected
				}
			case "a":
				for i := range m.items {
					if m.items[i].status != core.InstallStatusInstalled {
						m.items[i].selected = true
					}
				}
			case "n":
				for i := range m.items {
					m.items[i].selected = false
				}
			case "enter":
				// Check if anything is selected
				hasSelection := false
				for _, item := range m.items {
					if item.selected {
						hasSelection = true
						break
					}
				}
				if hasSelection {
					m.phase = PhaseConfirm
				}
			case "q", "esc":
				app.NavigateTo(ViewHome)
			}
		} else if m.phase == PhaseConfirm {
			switch msg.String() {
			case "y", "enter":
				m.phase = PhaseRun
				m.progress = 0
				m.installer = core.NewInstaller(app.state.Proxy)
				return m, m.startNextInstall(app)
			case "n", "esc":
				m.phase = PhaseSelect
			}
		} else if m.phase == PhaseDone {
			switch msg.String() {
			case "enter", "q", "esc":
				app.NavigateTo(ViewHome)
			}
		}

	case installDoneMsg:
		// Update item status
		if msg.index < len(m.items) {
			item := &m.items[msg.index]
			if msg.result.Error != nil {
				item.log = fmt.Sprintf("✗ %v", msg.result.Error)
				m.logs = append(m.logs, fmt.Sprintf("✗ %s: %v", item.module.Name, msg.result.Error))
			} else {
				item.status = msg.result.Status
				item.version = msg.result.Version
				item.log = "✓ installed"
				m.logs = append(m.logs, fmt.Sprintf("✓ %s %s", item.module.Name, msg.result.Version))

				// Update state
				app.state.SetInstalled(item.module.Name, core.ModuleState{
					Version:  msg.result.Version,
					Strategy: "tool",
				})
				app.state.Save()
			}
		}

		// Move to next
		m.progress++
		cmd := m.startNextInstall(app)
		if cmd == nil {
			m.phase = PhaseDone
			m.message = fmt.Sprintf("Installation complete. %d modules processed.", m.countSelected())
		}
		return m, cmd
	}

	return m, nil
}

func (m *InstallModel) startNextInstall(app *App) tea.Cmd {
	// Find next selected item
	idx := 0
	count := 0
	for i, item := range m.items {
		if item.selected {
			if count == m.progress {
				idx = i
				break
			}
			count++
		}
		if i == len(m.items)-1 {
			return nil // no more to install
		}
	}

	if count < m.progress {
		return nil
	}

	mod := m.items[idx].module
	installer := m.installer

	return func() tea.Msg {
		result := installer.Install(mod)
		return installDoneMsg{index: idx, result: result}
	}
}

func (m InstallModel) countSelected() int {
	n := 0
	for _, item := range m.items {
		if item.selected {
			n++
		}
	}
	return n
}

func (m InstallModel) View(app App) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📦 Install Tools"))
	b.WriteString("\n")

	switch m.phase {
	case PhaseSelect:
		b.WriteString(subtitleStyle.Render("Space: toggle • a: all uninstalled • Enter: proceed • q: back"))
		b.WriteString("\n\n")
		m.renderSelectList(&b)

	case PhaseConfirm:
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(colorYellow).Render("Install the following tools?"))
		b.WriteString("\n\n")
		for _, item := range m.items {
			if item.selected {
				b.WriteString(fmt.Sprintf("  • %s - %s\n", item.module.Name, item.module.Description))
			}
		}
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("y/Enter: confirm • n/Esc: cancel"))

	case PhaseRun:
		b.WriteString("\n")
		installed := 0
		total := m.countSelected()
		for _, item := range m.items {
			if !item.selected {
				continue
			}
			if item.log != "" {
				installed++
			}
		}
		b.WriteString(fmt.Sprintf("  Progress: %d / %d\n\n", installed, total))
		m.renderProgress(&b)

	case PhaseDone:
		b.WriteString("\n")
		b.WriteString(successStyle.Render(m.message))
		b.WriteString("\n\n")
		for _, log := range m.logs {
			if strings.HasPrefix(log, "✓") {
				b.WriteString(successStyle.Render("  " + log))
			} else {
				b.WriteString(errorStyle.Render("  " + log))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Enter/q: back to menu"))
	}

	return boxStyle.Render(b.String())
}

func (m InstallModel) renderSelectList(b *strings.Builder) {
	for i, item := range m.items {
		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}

		var checkbox string
		if item.status == core.InstallStatusInstalled {
			// Installed: show locked indicator
			checkbox = lipgloss.NewStyle().Foreground(colorGreen).Render("[✓]")
		} else if item.selected {
			checkbox = lipgloss.NewStyle().Foreground(colorBlue).Render("[✓]")
		} else {
			checkbox = lipgloss.NewStyle().Foreground(colorSubtext).Render("[ ]")
		}

		name := item.module.Name
		if item.status == core.InstallStatusInstalled {
			// Dim the name for installed items
			name = lipgloss.NewStyle().Foreground(colorSubtext).Render(name)
		} else if i == m.cursor {
			name = lipgloss.NewStyle().Foreground(colorLavender).Bold(true).Render(name)
		}

		// Status
		var statusStr string
		switch item.status {
		case core.InstallStatusInstalled:
			ver := item.version
			if ver == "" {
				ver = "installed"
			}
			statusStr = statusLinkedStyle.Render(fmt.Sprintf("✓ %s", ver))
		case core.InstallStatusMissing:
			statusStr = statusMissingStyle.Render("not installed")
		case core.InstallStatusOutdated:
			statusStr = statusStaleStyle.Render(fmt.Sprintf("↑ %s", item.version))
		}

		desc := lipgloss.NewStyle().Foreground(colorSubtext).Render(item.module.Description)

		b.WriteString(fmt.Sprintf("%s%s %s %s  %s\n", cursor, checkbox, name, statusStr, desc))
	}
}

func (m InstallModel) renderProgress(b *strings.Builder) {
	for _, item := range m.items {
		if !item.selected {
			continue
		}

		var icon string
		switch {
		case item.log == "":
			icon = lipgloss.NewStyle().Foreground(colorYellow).Render("⏳")
		case strings.HasPrefix(item.log, "✓"):
			icon = lipgloss.NewStyle().Foreground(colorGreen).Render("✓")
		default:
			icon = lipgloss.NewStyle().Foreground(colorRed).Render("✗")
		}

		name := lipgloss.NewStyle().Foreground(colorText).Render(item.module.Name)
		b.WriteString(fmt.Sprintf("  %s %s", icon, name))
		if item.log != "" && !strings.HasPrefix(item.log, "✓") {
			b.WriteString(fmt.Sprintf("  %s", lipgloss.NewStyle().Foreground(colorSubtext).Render(item.log)))
		}
		b.WriteString("\n")
	}
}
