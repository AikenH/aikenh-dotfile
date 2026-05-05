package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/aikenhong/dotsetup/internal/core"
)

// ModuleItem represents a displayable module in the list
type ModuleItem struct {
	module   core.Module
	selected bool
	status   core.LinkStatus
	isGroup  bool   // true if this is a group header
	group    string // group name (for headers)
}

// ModulesModel handles the module selection view
type ModulesModel struct {
	items    []ModuleItem
	cursor   int
	state    *core.State
	repoRoot string
	message  string
	width    int
	height   int
}

func NewModulesModel(modules []core.Module, state *core.State, repoRoot string) ModulesModel {
	m := ModulesModel{
		state:    state,
		repoRoot: repoRoot,
	}
	m.buildItems(modules)
	return m
}

func (m *ModulesModel) buildItems(modules []core.Module) {
	m.items = nil

	// Collect groups and independent modules
	groups := core.GroupModules(modules)
	independent := core.IndependentModules(modules)

	// Add grouped modules first (only truly exclusive groups)
	groupOrder := []string{"starship-style"}
	for _, gName := range groupOrder {
		gMods, ok := groups[gName]
		if !ok {
			continue
		}
		// Add group header
		m.items = append(m.items, ModuleItem{
			isGroup: true,
			group:   gName,
		})
		// Add modules in group
		for _, mod := range gMods {
			status := core.CheckLinkStatus(mod, m.repoRoot)
			selected := status == core.StatusLinked
			// If state has a group choice, reflect it
			if choice := m.state.GetGroupChoice(gName); choice != "" {
				selected = (choice == mod.Name)
			}
			m.items = append(m.items, ModuleItem{
				module:   mod,
				selected: selected,
				status:   status,
			})
		}
	}

	// Add independent modules
	if len(independent) > 0 {
		m.items = append(m.items, ModuleItem{
			isGroup: true,
			group:   "independent",
		})
		for _, mod := range independent {
			status := core.CheckLinkStatus(mod, m.repoRoot)
			selected := status == core.StatusLinked
			m.items = append(m.items, ModuleItem{
				module:   mod,
				selected: selected,
				status:   status,
			})
		}
	}
}

func (m ModulesModel) Update(msg tea.Msg, app *App) (ModulesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.moveCursor(-1)
		case "down", "j":
			m.moveCursor(1)
		case " ", "space":
			m.toggleCurrent()
		case "a":
			m.selectAll()
		case "n":
			m.selectNone()
		case "enter":
			m.applyChanges(app)
		case "q", "esc":
			app.NavigateTo(ViewHome)
		}
	}
	return m, nil
}

func (m *ModulesModel) moveCursor(delta int) {
	newCursor := m.cursor + delta
	// Skip group headers
	for newCursor >= 0 && newCursor < len(m.items) && m.items[newCursor].isGroup {
		newCursor += delta
	}
	if newCursor >= 0 && newCursor < len(m.items) {
		m.cursor = newCursor
	}
}

func (m *ModulesModel) toggleCurrent() {
	if m.cursor >= len(m.items) || m.items[m.cursor].isGroup {
		return
	}

	item := &m.items[m.cursor]
	mod := item.module

	// If module is in a group, implement radio behavior
	if mod.Group != "" {
		if item.selected {
			// Deselect (choose none for this group)
			item.selected = false
		} else {
			// Deselect others in the same group, select this one
			for i := range m.items {
				if !m.items[i].isGroup && m.items[i].module.Group == mod.Group {
					m.items[i].selected = false
				}
			}
			item.selected = true
		}
	} else {
		// Independent module: simple toggle
		item.selected = !item.selected
	}
}

func (m *ModulesModel) selectAll() {
	// For grouped modules, we can't select all (exclusive)
	// Only select all independent modules
	for i := range m.items {
		if !m.items[i].isGroup && m.items[i].module.Group == "" {
			m.items[i].selected = true
		}
	}
}

func (m *ModulesModel) selectNone() {
	for i := range m.items {
		if !m.items[i].isGroup {
			m.items[i].selected = false
		}
	}
}

func (m *ModulesModel) applyChanges(app *App) {
	linked := 0
	unlinked := 0
	var errors []string

	for _, item := range m.items {
		if item.isGroup {
			continue
		}

		mod := item.module
		currentStatus := core.CheckLinkStatus(mod, m.repoRoot)

		if item.selected && currentStatus != core.StatusLinked {
			// Need to link
			var result core.LinkResult
			if mod.Strategy == core.StrategyAppend {
				result = core.Append(mod, m.repoRoot)
			} else {
				doBackup := currentStatus == core.StatusConflict
				result = core.Link(mod, m.repoRoot, doBackup)
			}
			if result.Error != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", mod.Name, result.Error))
			} else {
				linked++
				app.state.SetInstalled(mod.Name, core.ModuleState{
					Strategy:   string(mod.Strategy),
					TargetPath: core.ExpandPath(mod.Target),
					BackupPath: result.BackupPath,
				})
			}
			// Record group choice
			if mod.Group != "" {
				app.state.SetGroupChoice(mod.Group, mod.Name)
			}
		} else if !item.selected && currentStatus == core.StatusLinked {
			// Need to unlink
			if err := core.Unlink(mod); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", mod.Name, err))
			} else {
				unlinked++
				app.state.SetUninstalled(mod.Name)
			}
			// Clear group choice if this was the chosen one
			if mod.Group != "" && app.state.GetGroupChoice(mod.Group) == mod.Name {
				app.state.SetGroupChoice(mod.Group, "none")
			}
		}
	}

	// Save state
	if err := app.state.Save(); err != nil {
		errors = append(errors, fmt.Sprintf("saving state: %v", err))
	}

	// Build message
	parts := []string{}
	if linked > 0 {
		parts = append(parts, fmt.Sprintf("linked %d", linked))
	}
	if unlinked > 0 {
		parts = append(parts, fmt.Sprintf("unlinked %d", unlinked))
	}
	if len(errors) > 0 {
		parts = append(parts, fmt.Sprintf("%d errors", len(errors)))
	}
	if len(parts) == 0 {
		m.message = "No changes"
	} else {
		m.message = strings.Join(parts, ", ")
		if len(errors) > 0 {
			m.message += "\n" + strings.Join(errors, "\n")
		}
	}

	// Refresh status
	m.buildItems(app.allMods)
}

func (m ModulesModel) View(app App) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📁 Link Configs"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Space: toggle • Enter: apply • a: all • n: none • q: back"))
	b.WriteString("\n\n")

	for i, item := range m.items {
		if item.isGroup {
			header := formatGroupName(item.group)
			b.WriteString(groupHeaderStyle.Render(header))
			b.WriteString("\n")
			continue
		}

		// Cursor
		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}

		// Checkbox (radio for groups, checkbox for independent)
		var checkbox string
		if item.module.Group != "" {
			// Radio style
			if item.selected {
				checkbox = lipgloss.NewStyle().Foreground(colorGreen).Render("◉")
			} else {
				checkbox = lipgloss.NewStyle().Foreground(colorSubtext).Render("○")
			}
		} else {
			// Checkbox style
			if item.selected {
				checkbox = lipgloss.NewStyle().Foreground(colorGreen).Render("[✓]")
			} else {
				checkbox = lipgloss.NewStyle().Foreground(colorSubtext).Render("[ ]")
			}
		}

		// Status indicator
		statusStr := formatStatus(item.status)

		// Module name and description
		name := item.module.Name
		if i == m.cursor {
			name = lipgloss.NewStyle().Foreground(colorLavender).Bold(true).Render(name)
		} else {
			name = lipgloss.NewStyle().Foreground(colorText).Render(name)
		}

		desc := lipgloss.NewStyle().Foreground(colorSubtext).Render(item.module.Description)

		line := fmt.Sprintf("%s%s %s %s  %s", cursor, checkbox, name, statusStr, desc)
		b.WriteString(line)
		b.WriteString("\n")
	}

	// Message area
	if m.message != "" {
		b.WriteString("\n")
		if strings.Contains(m.message, "error") {
			b.WriteString(errorStyle.Render(m.message))
		} else {
			b.WriteString(successStyle.Render("✓ " + m.message))
		}
		b.WriteString("\n")
	}

	return boxStyle.Render(b.String())
}

func formatGroupName(name string) string {
	switch name {
	case "starship-style":
		return "── Starship Style (pick one) ──"
	case "independent":
		return "── Modules ──"
	default:
		return fmt.Sprintf("── %s (pick one) ──", name)
	}
}

func formatStatus(status core.LinkStatus) string {
	switch status {
	case core.StatusLinked:
		return statusLinkedStyle.Render("✓")
	case core.StatusMissing:
		return statusMissingStyle.Render("○")
	case core.StatusConflict:
		return statusConflictStyle.Render("⚠")
	case core.StatusStale:
		return statusStaleStyle.Render("↻")
	default:
		return " "
	}
}
