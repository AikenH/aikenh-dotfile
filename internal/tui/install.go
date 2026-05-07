package tui

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

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

// maxLogLines is how many streaming log lines we keep in the scrolling area
const maxLogLines = 6

// InstallItem represents a tool in the install list
type InstallItem struct {
	module    core.Module
	selected  bool
	status    core.InstallStatus
	version   string
	log       string        // short outcome string (e.g. "✓ installed")
	startTime time.Time     // when this module started installing
	elapsed   time.Duration // how long it took (set on completion)
}

// installDoneMsg signals that one install finished
type installDoneMsg struct {
	index  int
	result core.InstallResult
}

// logLineMsg carries a single streamed log line from the runner
type logLineMsg struct {
	line string
}

// tickMsg is sent periodically to refresh elapsed-time display while running
type tickMsg time.Time

// InstallModel handles the tool installation view
type InstallModel struct {
	items       []InstallItem
	cursor      int
	phase       InstallPhase
	progress    int      // index into the selected-items order (0-based)
	currentIdx  int      // items[] index of the module currently installing
	message     string
	streamLogs  []string // last N lines streamed in real-time
	installer   *core.Installer
	logChan     chan string // receives streamed log lines
	width       int
	height      int
}

func NewInstallModel(modules []core.Module, state *core.State) InstallModel {
	m := InstallModel{
		phase: PhaseSelect,
	}

	installer := core.NewInstaller(state.Proxy)

	for _, mod := range modules {
		if mod.Type != core.TypeTool {
			continue
		}
		item := InstallItem{
			module: mod,
			status: core.InstallStatusMissing,
		}
		if installer.CheckInstalled(mod) {
			item.status = core.InstallStatusInstalled
			item.version = installer.GetVersion(mod)
		} else if state.IsInstalled(mod.Name) {
			item.status = core.InstallStatusMissing
		}
		m.items = append(m.items, item)
	}

	return m
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
				hasSelection := false
				for _, item := range m.items {
					if item.selected {
						hasSelection = true
						break
					}
				}
				if hasSelection {
					m.resolveDeps(app)
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
				m.streamLogs = nil
				m.logChan = make(chan string, 256)
				m.installer = core.NewInstaller(app.state.Proxy)
				// Wire up log streaming via channel
				ch := m.logChan
				m.installer.SetLogFunc(func(line string) {
					select {
					case ch <- line:
					default: // drop if full (non-blocking)
					}
				})
				return m, tea.Batch(m.startNextInstall(app), m.drainLogs(), tick())
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
		if msg.index < len(m.items) {
			item := &m.items[msg.index]
			item.elapsed = time.Since(item.startTime)
			if msg.result.Error != nil {
				item.log = fmt.Sprintf("✗ %v", msg.result.Error)
			} else {
				item.status = msg.result.Status
				item.version = msg.result.Version
				item.log = "✓ installed"

				app.state.SetInstalled(item.module.Name, core.ModuleState{
					Version:  msg.result.Version,
					Strategy: "tool",
				})
				app.state.Save()
			}
		}

		m.progress++
		cmd := m.startNextInstall(app)
		if cmd == nil {
			m.phase = PhaseDone
			m.message = fmt.Sprintf("Installation complete. %d modules processed.", m.countSelected())
			return m, nil
		}
		return m, tea.Batch(cmd, m.drainLogs())

	case logLineMsg:
		// Append streamed log text, split on newlines
		for _, raw := range strings.Split(msg.line, "\n") {
			raw = strings.TrimRight(raw, "\r")
			if raw == "" {
				continue
			}
			m.streamLogs = append(m.streamLogs, raw)
			if len(m.streamLogs) > maxLogLines {
				m.streamLogs = m.streamLogs[len(m.streamLogs)-maxLogLines:]
			}
		}
		// Keep draining
		if m.phase == PhaseRun {
			return m, m.drainLogs()
		}

	case tickMsg:
		// Just re-render (for elapsed time updates)
		if m.phase == PhaseRun {
			return m, tick()
		}
	}

	return m, nil
}

// tick returns a cmd that fires a tickMsg roughly every second
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// drainLogs polls the log channel without blocking
func (m *InstallModel) drainLogs() tea.Cmd {
	ch := m.logChan
	if ch == nil {
		return nil
	}
	return func() tea.Msg {
		line, ok := <-ch
		if !ok {
			return nil
		}
		return logLineMsg{line: line}
	}
}

// resolveDeps auto-selects uninstalled dependencies
func (m *InstallModel) resolveDeps(app *App) {
	var toolMods []core.Module
	for _, item := range m.items {
		toolMods = append(toolMods, item.module)
	}
	graph := core.NewDepGraph(toolMods)

	var selectedNames []string
	for _, item := range m.items {
		if item.selected {
			selectedNames = append(selectedNames, item.module.Name)
		}
	}

	resolved, err := graph.TopoSortSubset(selectedNames, func(name string) bool {
		for _, item := range m.items {
			if item.module.Name == name && item.status == core.InstallStatusInstalled {
				return true
			}
		}
		return false
	})
	if err != nil {
		return
	}

	resolvedNames := make(map[string]bool)
	for _, mod := range resolved {
		resolvedNames[mod.Name] = true
	}
	for i, item := range m.items {
		if resolvedNames[item.module.Name] && item.status != core.InstallStatusInstalled {
			m.items[i].selected = true
		}
	}
}

func (m *InstallModel) startNextInstall(app *App) tea.Cmd {
	// Walk selected items in order; skip already-done (have log set)
	count := 0
	for i, item := range m.items {
		if !item.selected {
			continue
		}
		if count == m.progress {
			m.currentIdx = i
			m.items[i].startTime = time.Now()
			m.streamLogs = nil // clear log area for new module
			mod := item.module
			installer := m.installer
			return func() tea.Msg {
				result := installer.Install(mod)
				return installDoneMsg{index: i, result: result}
			}
		}
		count++
	}
	return nil // all done
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

// ──────────────────────────────────────────────
// View
// ──────────────────────────────────────────────

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
		m.renderRunPhase(&b, contentWidth(m.width))

	case PhaseDone:
		b.WriteString("\n")
		b.WriteString(successStyle.Render(m.message))
		b.WriteString("\n\n")
		for _, item := range m.items {
			if !item.selected {
				continue
			}
			if strings.HasPrefix(item.log, "✓") {
				ver := ""
				if item.version != "" {
					ver = " " + item.version
				}
				elapsed := ""
				if item.elapsed > 0 {
					elapsed = fmt.Sprintf("  %s", formatDuration(item.elapsed))
				}
				line := fmt.Sprintf("  ✓ %-20s%s%s", item.module.Name, ver, elapsed)
				b.WriteString(successStyle.Render(line))
			} else {
				b.WriteString(errorStyle.Render(fmt.Sprintf("  ✗ %-20s  %s", item.module.Name, item.log)))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Enter/q: back to menu"))
	}

	return boxStyle.Width(contentWidth(m.width)).Render(b.String())
}

func (m InstallModel) renderSelectList(b *strings.Builder) {
	for i, item := range m.items {
		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}

		var checkbox string
		if item.status == core.InstallStatusInstalled {
			checkbox = lipgloss.NewStyle().Foreground(colorGreen).Render("[✓]")
		} else if item.selected {
			checkbox = lipgloss.NewStyle().Foreground(colorBlue).Render("[✓]")
		} else {
			checkbox = lipgloss.NewStyle().Foreground(colorSubtext).Render("[ ]")
		}

		name := item.module.Name
		if item.status == core.InstallStatusInstalled {
			name = lipgloss.NewStyle().Foreground(colorSubtext).Render(name)
		} else if i == m.cursor {
			name = lipgloss.NewStyle().Foreground(colorLavender).Bold(true).Render(name)
		}

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

func (m InstallModel) renderRunPhase(b *strings.Builder, cw int) {
	total := m.countSelected()
	done := 0
	for _, item := range m.items {
		if item.selected && item.log != "" {
			done++
		}
	}

	// ── progress counter ──────────────────────────────────────────────
	counter := lipgloss.NewStyle().
		Foreground(colorLavender).
		Bold(true).
		Render(fmt.Sprintf("%d / %d", done, total))
	b.WriteString(fmt.Sprintf("  Installing Tools  %s\n", counter))

	// ── progress bar — use ~60% of content width, min 20 ─────────────
	barWidth := cw*6/10
	if barWidth < 20 {
		barWidth = 20
	}
	if barWidth > 60 {
		barWidth = 60
	}
	filled := 0
	if total > 0 {
		filled = barWidth * done / total
	}
	bar := lipgloss.NewStyle().Foreground(colorBlue).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(colorSurface).Render(strings.Repeat("░", barWidth-filled))
	pct := 0
	if total > 0 {
		pct = 100 * done / total
	}
	b.WriteString(fmt.Sprintf("  %s  %d%%\n\n", bar, pct))

	// ── per-module status lines ───────────────────────────────────────
	for i, item := range m.items {
		if !item.selected {
			continue
		}

		var icon, nameStr, rightStr string

		switch {
		case item.log == "" && i == m.currentIdx:
			// Currently installing
			icon = lipgloss.NewStyle().Foreground(colorYellow).Render("⏳")
			nameStr = lipgloss.NewStyle().Foreground(colorLavender).Bold(true).Render(item.module.Name)
			elapsed := time.Since(item.startTime)
			rightStr = lipgloss.NewStyle().Foreground(colorSubtext).
				Render(fmt.Sprintf("installing...  %s", formatDuration(elapsed)))
		case item.log == "":
			// Pending
			icon = "  "
			nameStr = lipgloss.NewStyle().Foreground(colorSubtext).Render(item.module.Name)
			rightStr = lipgloss.NewStyle().Foreground(colorOverlay).Render("pending")
		case strings.HasPrefix(item.log, "✓"):
			icon = lipgloss.NewStyle().Foreground(colorGreen).Render("✓")
			nameStr = lipgloss.NewStyle().Foreground(colorText).Render(item.module.Name)
			ver := item.version
			if ver == "" {
				ver = ""
			} else {
				ver = " " + ver
			}
			rightStr = lipgloss.NewStyle().Foreground(colorSubtext).
				Render(fmt.Sprintf("%s  %s", ver, formatDuration(item.elapsed)))
		default:
			icon = lipgloss.NewStyle().Foreground(colorRed).Render("✗")
			nameStr = lipgloss.NewStyle().Foreground(colorRed).Render(item.module.Name)
			rightStr = lipgloss.NewStyle().Foreground(colorRed).Render(item.log)
		}

		b.WriteString(fmt.Sprintf("  %s %-20s  %s\n", icon, nameStr, rightStr))
	}

	// ── streaming log area ───────────────────────────────────────────
	if len(m.streamLogs) > 0 {
		b.WriteString("\n")
		dividerWidth := cw - 4
		if dividerWidth < 20 {
			dividerWidth = 20
		}
		divider := lipgloss.NewStyle().Foreground(colorOverlay).Render(strings.Repeat("─", dividerWidth))
		b.WriteString("  " + divider + "\n")
		logStyle := lipgloss.NewStyle().Foreground(colorSubtext)
		maxLineWidth := cw - 6 // "  > " prefix = 4 chars + small margin
		if maxLineWidth < 20 {
			maxLineWidth = 20
		}
		for _, line := range m.streamLogs {
			if utf8.RuneCountInString(line) > maxLineWidth {
				runes := []rune(line)
				line = string(runes[:maxLineWidth-3]) + "..."
			}
			b.WriteString("  " + logStyle.Render("> "+line) + "\n")
		}
	}
}

// formatDuration formats a duration as "1.2s" or "1m23s"
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm%ds", m, s)
}
