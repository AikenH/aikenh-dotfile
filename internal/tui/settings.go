package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SettingsField represents which field is being edited
type SettingsField int

const (
	FieldProxy SettingsField = iota
	FieldSave
	FieldBack
)

// SettingsModel handles the settings/configuration view
type SettingsModel struct {
	cursor   SettingsField
	editing  bool
	proxyBuf string
	message  string
	width    int
}

func NewSettingsModel(proxy string) SettingsModel {
	return SettingsModel{
		proxyBuf: proxy,
	}
}

func (s SettingsModel) Update(msg tea.Msg, app *App) (SettingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.editing {
			switch msg.String() {
			case "enter":
				s.editing = false
			case "esc":
				s.proxyBuf = app.state.Proxy
				s.editing = false
			case "backspace":
				if len(s.proxyBuf) > 0 {
					s.proxyBuf = s.proxyBuf[:len(s.proxyBuf)-1]
				}
			case "ctrl+u":
				s.proxyBuf = ""
			default:
				if len(msg.String()) == 1 {
					s.proxyBuf += msg.String()
				}
			}
			return s, nil
		}

		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		case "down", "j":
			if s.cursor < FieldBack {
				s.cursor++
			}
		case "enter":
			switch s.cursor {
			case FieldProxy:
				s.editing = true
			case FieldSave:
				app.state.Proxy = s.proxyBuf
				if err := app.state.Save(); err != nil {
					s.message = fmt.Sprintf("Error saving: %v", err)
				} else {
					s.message = "Settings saved ✓"
				}
			case FieldBack:
				app.NavigateTo(ViewHome)
			}
		case "q", "esc":
			app.NavigateTo(ViewHome)
		}
	}
	return s, nil
}

func (s SettingsModel) View(app App) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("⚙️  Settings"))
	b.WriteString("\n\n")

	// Proxy field
	proxyLabel := "Proxy URL"
	proxyValue := s.proxyBuf
	if proxyValue == "" {
		proxyValue = "(none)"
	}

	cursor := "  "
	if s.cursor == FieldProxy {
		cursor = "▸ "
	}

	if s.editing {
		b.WriteString(fmt.Sprintf("%s%s: %s▋\n",
			cursor,
			lipgloss.NewStyle().Foreground(colorLavender).Render(proxyLabel),
			lipgloss.NewStyle().Foreground(colorText).Render(s.proxyBuf),
		))
		b.WriteString(helpStyle.Render("    Enter: confirm • Esc: cancel • Ctrl+U: clear"))
		b.WriteString("\n")
	} else {
		b.WriteString(fmt.Sprintf("%s%s: %s\n",
			cursor,
			lipgloss.NewStyle().Foreground(colorLavender).Render(proxyLabel),
			lipgloss.NewStyle().Foreground(colorText).Render(proxyValue),
		))
	}

	b.WriteString("\n")

	// Platform info (read-only)
	b.WriteString(fmt.Sprintf("  Platform: %s\n", lipgloss.NewStyle().Foreground(colorSubtext).Render(app.platform.String())))
	b.WriteString(fmt.Sprintf("  State:    %s\n", lipgloss.NewStyle().Foreground(colorSubtext).Render(StatePath())))

	b.WriteString("\n")

	// Save button
	cursor = "  "
	if s.cursor == FieldSave {
		cursor = "▸ "
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, selectedItemStyle.Render("💾 Save")))
	} else {
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, menuItemStyle.Render("💾 Save")))
	}

	// Back button
	cursor = "  "
	if s.cursor == FieldBack {
		cursor = "▸ "
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, selectedItemStyle.Render("← Back")))
	} else {
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, menuItemStyle.Render("← Back")))
	}

	// Message
	if s.message != "" {
		b.WriteString("\n")
		if strings.Contains(s.message, "Error") {
			b.WriteString(errorStyle.Render(s.message))
		} else {
			b.WriteString(successStyle.Render(s.message))
		}
	}

	if !s.editing {
		b.WriteString("\n" + helpStyle.Render("↑/↓ navigate • Enter select/edit • q back"))
	}

	return boxStyle.Width(contentWidth(s.width)).Render(b.String())
}

// StatePath helper for display
func StatePath() string {
	return "~/.local/share/dotsetup/state.json"
}
