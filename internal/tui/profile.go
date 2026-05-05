package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/aikenhong/dotsetup/internal/core"
)

// ProfileModel handles profile selection
type ProfileModel struct {
	profiles []*core.Profile
	cursor   int
	active   bool
}

func NewProfileModel(repoRoot string) ProfileModel {
	profilesDir := filepath.Join(repoRoot, "profiles")
	profiles, _ := core.LoadAllProfiles(profilesDir)
	return ProfileModel{
		profiles: profiles,
		active:   len(profiles) > 0,
	}
}

func (p ProfileModel) Update(msg tea.Msg, app *App) (ProfileModel, tea.Cmd) {
	if !p.active {
		return p, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			}
		case "down", "j":
			if p.cursor < len(p.profiles)-1 {
				p.cursor++
			}
		case "enter":
			if p.cursor < len(p.profiles) {
				profile := p.profiles[p.cursor]
				selected, groupChoices := core.ResolveProfile(profile, app.allMods)

				// Apply group choices to state
				for group, choice := range groupChoices {
					if choice != "prompt" {
						app.state.SetGroupChoice(group, choice)
					}
				}

				// Pre-select modules in the modules view
				_ = selected // will be used to pre-fill selection

				app.state.Save()
				p.active = false
				app.NavigateTo(ViewHome)
			}
		case "q", "esc":
			p.active = false
			app.NavigateTo(ViewHome)
		}
	}
	return p, nil
}

func (p ProfileModel) View() string {
	if !p.active || len(p.profiles) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("🎯 Select a Profile"))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Choose a starting configuration (you can customize later)"))
	b.WriteString("\n\n")

	for i, profile := range p.profiles {
		cursor := "  "
		if i == p.cursor {
			cursor = "▸ "
		}

		name := profile.Name
		if i == p.cursor {
			name = lipgloss.NewStyle().Foreground(colorLavender).Bold(true).Render(name)
		} else {
			name = lipgloss.NewStyle().Foreground(colorText).Render(name)
		}

		desc := lipgloss.NewStyle().Foreground(colorSubtext).Render(profile.Description)
		modCount := lipgloss.NewStyle().Foreground(colorMauve).Render(
			fmt.Sprintf("(%d modules)", len(profile.Modules)))

		b.WriteString(fmt.Sprintf("%s%s %s  %s\n", cursor, name, modCount, desc))
	}

	b.WriteString("\n" + helpStyle.Render("↑/↓ navigate • Enter select • Esc skip"))

	return boxStyle.Render(b.String())
}
