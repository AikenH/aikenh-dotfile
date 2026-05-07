package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// HomeModel represents the main menu
type HomeModel struct {
	choices []string
	cursor  int
	width   int
}

func NewHomeModel() HomeModel {
	return HomeModel{
		choices: []string{
			"Link Configs",
			"Install Tools",
			"Settings",
			"Status",
			"Quit",
		},
	}
}

func (h HomeModel) Update(msg tea.Msg, app *App) (HomeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if h.cursor > 0 {
				h.cursor--
			}
		case "down", "j":
			if h.cursor < len(h.choices)-1 {
				h.cursor++
			}
		case "enter":
			switch h.cursor {
			case 0: // Link Configs
				app.NavigateTo(ViewModules)
			case 1: // Install Tools
				app.NavigateTo(ViewInstall)
			case 2: // Settings
				app.NavigateTo(ViewSettings)
			case 3: // Status
				app.NavigateTo(ViewStatus)
			case 4: // Quit
				return h, tea.Quit
			}
		}
	}
	return h, nil
}

func (h HomeModel) View(app App) string {
	s := titleStyle.Render("🔧 dotsetup") + "\n"
	s += subtitleStyle.Render(fmt.Sprintf("Platform: %s", app.platform.String())) + "\n\n"

	icons := []string{"📁", "📦", "⚙️ ", "📊", "🚪"}

	for i, choice := range h.choices {
		cursor := "  "
		style := menuItemStyle
		if h.cursor == i {
			cursor = "▸ "
			style = selectedItemStyle
		}
		s += style.Render(fmt.Sprintf("%s%s %s", cursor, icons[i], choice)) + "\n"
	}

	s += "\n" + helpStyle.Render("↑/↓ navigate • enter select • q quit")

	return boxStyle.Width(contentWidth(h.width)).Render(s)
}
