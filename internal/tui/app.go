package tui

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/aikenhong/dotsetup/internal/core"
)

// View represents which screen we're on
type View int

const (
	ViewHome View = iota
	ViewModules
	ViewInstall
	ViewSettings
	ViewStatus
	ViewProfile
)

// App is the root model
type App struct {
	view     View
	home     HomeModel
	modules  ModulesModel
	install  InstallModel
	settings SettingsModel
	profile  ProfileModel
	state    *core.State
	platform core.Platform
	allMods  []core.Module
	repoRoot string
	width    int
	height   int
	err      error
}

// NewApp initializes the application
func NewApp() (*App, error) {
	// Find repo root (where modules/ directory lives)
	repoRoot, err := findRepoRoot()
	if err != nil {
		return nil, fmt.Errorf("finding repo root: %w", err)
	}

	// Load state
	state, err := core.LoadState()
	if err != nil {
		return nil, fmt.Errorf("loading state: %w", err)
	}

	// Load modules
	modulesDir := filepath.Join(repoRoot, "modules")
	allMods, err := core.LoadAllModules(modulesDir)
	if err != nil {
		return nil, fmt.Errorf("loading modules: %w", err)
	}

	// Filter by platform
	platform := core.DetectPlatform()
	allMods = core.FilterByPlatform(allMods)

	// Determine initial view: first run → profile selection
	initialView := ViewHome
	if isFirstRun(state) {
		initialView = ViewProfile
	}

	app := &App{
		view:     initialView,
		state:    state,
		platform: platform,
		allMods:  allMods,
		repoRoot: repoRoot,
	}

	app.home = NewHomeModel()
	app.modules = NewModulesModel(allMods, state, repoRoot)
	app.install = NewInstallModel(allMods, state)
	app.settings = NewSettingsModel(state.Proxy)
	app.profile = NewProfileModel(repoRoot)

	return app, nil
}

// isFirstRun checks if this is a fresh install (no modules managed yet)
func isFirstRun(state *core.State) bool {
	return len(state.Modules) == 0 && len(state.GroupPicks) == 0
}

func (a App) Init() tea.Cmd {
	return tea.SetWindowTitle("dotsetup")
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return a, tea.Quit
		}
		// Only handle q at the top level if on home screen
		if a.view == ViewHome {
			if msg.String() == "q" {
				return a, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.modules.width = msg.Width
		a.modules.height = msg.Height
		a.install.width = msg.Width
		a.install.height = msg.Height

	case installDoneMsg:
		// Route install completion messages to install view
		newInstall, cmd := a.install.Update(msg, &a)
		a.install = newInstall
		return a, cmd
	}

	switch a.view {
	case ViewHome:
		newHome, cmd := a.home.Update(msg, &a)
		a.home = newHome
		return a, cmd
	case ViewModules:
		newMods, cmd := a.modules.Update(msg, &a)
		a.modules = newMods
		return a, cmd
	case ViewInstall:
		newInstall, cmd := a.install.Update(msg, &a)
		a.install = newInstall
		return a, cmd
	case ViewSettings:
		newSettings, cmd := a.settings.Update(msg, &a)
		a.settings = newSettings
		return a, cmd
	case ViewProfile:
		newProfile, cmd := a.profile.Update(msg, &a)
		a.profile = newProfile
		return a, cmd
	case ViewStatus:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "q", "esc", "enter":
				a.view = ViewHome
			}
		}
		return a, nil
	}

	return a, nil
}

func (a App) View() string {
	switch a.view {
	case ViewHome:
		return a.home.View(a)
	case ViewModules:
		return a.modules.View(a)
	case ViewInstall:
		return a.install.View(a)
	case ViewSettings:
		return a.settings.View(a)
	case ViewProfile:
		return a.profile.View()
	case ViewStatus:
		return a.renderStatusView()
	default:
		return "Not implemented yet"
	}
}

// NavigateTo switches to a different view
func (a *App) NavigateTo(v View) {
	a.view = v
	// Refresh sub-models when navigating
	switch v {
	case ViewModules:
		a.modules = NewModulesModel(a.allMods, a.state, a.repoRoot)
	case ViewInstall:
		a.install = NewInstallModel(a.allMods, a.state)
	case ViewSettings:
		a.settings = NewSettingsModel(a.state.Proxy)
	case ViewProfile:
		a.profile = NewProfileModel(a.repoRoot)
	}
}

func (a App) renderStatusView() string {
	var b string
	b += titleStyle.Render("📊 Status Overview") + "\n\n"

	b += fmt.Sprintf("  Platform: %s\n", a.platform.String()) + "\n"

	linked := 0
	missing := 0
	conflict := 0
	tools := 0
	toolsInstalled := 0

	for _, mod := range a.allMods {
		if mod.Type == core.TypeTool {
			tools++
			if a.state.IsInstalled(mod.Name) {
				toolsInstalled++
			}
			continue
		}
		status := core.CheckLinkStatus(mod, a.repoRoot)
		switch status {
		case core.StatusLinked:
			linked++
		case core.StatusMissing:
			missing++
		case core.StatusConflict, core.StatusStale:
			conflict++
		}
	}

	b += fmt.Sprintf("  Configs: %s linked, %s not linked, %s conflicts\n",
		statusLinkedStyle.Render(fmt.Sprintf("%d", linked)),
		statusMissingStyle.Render(fmt.Sprintf("%d", missing)),
		statusConflictStyle.Render(fmt.Sprintf("%d", conflict)),
	)
	b += fmt.Sprintf("  Tools:   %s / %s installed\n",
		statusLinkedStyle.Render(fmt.Sprintf("%d", toolsInstalled)),
		statusMissingStyle.Render(fmt.Sprintf("%d", tools)),
	)

	// Group choices
	b += "\n  Group Choices:\n"
	groups := core.GroupModules(a.allMods)
	for gName := range groups {
		choice := a.state.GetGroupChoice(gName)
		if choice == "" {
			choice = "(not set)"
		}
		b += fmt.Sprintf("    %s: %s\n", gName, choice)
	}

	b += "\n" + helpStyle.Render("q/Esc: back to menu")

	return boxStyle.Render(b)
}

func findRepoRoot() (string, error) {
	// Try current directory first
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Look for modules/ directory up the tree
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "modules")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return cwd, nil
}
