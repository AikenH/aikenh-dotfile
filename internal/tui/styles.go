package tui

import "github.com/charmbracelet/lipgloss"

// contentWidth returns the safe inner content width for a given terminal width.
// boxStyle uses RoundedBorder (2 chars) + Padding(1,2) (4 chars) = 6 chars horizontal overhead.
func contentWidth(termWidth int) int {
	w := termWidth - 6
	if w < 60 {
		w = 60
	}
	return w
}

var (
	// Colors (Catppuccin Mocha inspired)
	colorBase     = lipgloss.Color("#1e1e2e")
	colorSurface  = lipgloss.Color("#313244")
	colorOverlay  = lipgloss.Color("#45475a")
	colorText     = lipgloss.Color("#cdd6f4")
	colorSubtext  = lipgloss.Color("#a6adc8")
	colorLavender = lipgloss.Color("#b4befe")
	colorGreen    = lipgloss.Color("#a6e3a1")
	colorRed      = lipgloss.Color("#f38ba8")
	colorYellow   = lipgloss.Color("#f9e2af")
	colorPeach    = lipgloss.Color("#fab387")
	colorBlue     = lipgloss.Color("#89b4fa")
	colorMauve    = lipgloss.Color("#cba6f7")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorLavender).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorSubtext).
			MarginBottom(1)

	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(colorLavender).
				Bold(true)

	statusLinkedStyle = lipgloss.NewStyle().
				Foreground(colorGreen)

	statusMissingStyle = lipgloss.NewStyle().
				Foreground(colorSubtext)

	statusConflictStyle = lipgloss.NewStyle().
				Foreground(colorYellow)

	statusStaleStyle = lipgloss.NewStyle().
				Foreground(colorPeach)

	groupHeaderStyle = lipgloss.NewStyle().
				Foreground(colorMauve).
				Bold(true).
				MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorOverlay).
			MarginTop(1)

	successStyle = lipgloss.NewStyle().
			Foreground(colorGreen)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorRed)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorOverlay).
			Padding(1, 2)
)
