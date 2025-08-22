package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"claude_analysis/cmd/installer/internal/config"
	"claude_analysis/cmd/installer/internal/install"
	"claude_analysis/cmd/installer/internal/logger"
	"claude_analysis/cmd/installer/internal/ui"
)

func main() {
	// Ensure child processes that support NO_COLOR also disable colorized output
	os.Setenv("NO_COLOR", "1")
	// Allow self-signed certs for current process
	os.Setenv("NODE_TLS_REJECT_UNAUTHORIZED", "0")

	// Create main menu items
	items := []list.Item{
		ui.Item{
			TitleText:     "🚀 Full Installation",
			DescText:      "Node.js + Claude CLI + Configuration",
			Action:        install.RunFullInstall,
			IsFullInstall: true,
		},
		ui.Item{
			TitleText:     "🔑 Update GAISF API Key",
			DescText:      "Update GAISF token in existing configuration",
			Action:        func() error { return config.UpdateClaudeCodeSettings() },
			IsFullInstall: false,
		},
		ui.Item{
			TitleText:     "📦 Install Node.js",
			DescText:      "Install Node.js version 22+",
			Action:        install.InstallNodeJS,
			IsFullInstall: false,
		},
		ui.Item{
			TitleText:     "🤖 Install/Update Claude CLI",
			DescText:      "Install or update Claude CLI package",
			Action:        install.InstallOrUpdateClaude,
			IsFullInstall: false,
		},
		ui.Item{
			TitleText:     "❌ Exit",
			DescText:      "Quit the program",
			Action:        nil,
			IsFullInstall: false,
		},
	}

	const defaultWidth = 80
	const listHeight = 14

	l := list.New(items, ui.ItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Claude Code Installer & Configuration Tool"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = ui.TitleStyle
	l.Styles.PaginationStyle = ui.PaginationStyle
	l.Styles.HelpStyle = ui.HelpStyle

	// Create text input for forms
	ti := textinput.New()
	ti.Placeholder = "Enter value..."
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 60

	// Create spinner for processing operations
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Create GAISF configuration menu items
	gaisfItems := []list.Item{
		ui.Item{
			TitleText: "🔑 Auto-configure GAISF token",
			DescText:  "Login with username/password to get token",
			Action:    nil,
		},
		ui.Item{
			TitleText: "📝 Manual token input",
			DescText:  "Enter GAISF token manually",
			Action:    nil,
		},
		ui.Item{
			TitleText: "⏭️ Skip GAISF configuration",
			DescText:  "Continue without API authentication",
			Action:    nil,
		},
	}

	gl := list.New(gaisfItems, ui.ItemDelegate{}, defaultWidth, listHeight)
	gl.Title = "GAISF API Authentication Setup"
	gl.SetShowStatusBar(false)
	gl.SetFilteringEnabled(false)
	gl.Styles.Title = ui.TitleStyle
	gl.Styles.PaginationStyle = ui.PaginationStyle
	gl.Styles.HelpStyle = ui.HelpStyle

	m := ui.Model{
		List:        l,
		GAISFList:   gl,
		TextInput:   ti,
		Spinner:     s,
		CurrentView: ui.MainMenuView,
		GAISFConfig: ui.NewGAISFConfig(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	// Set up global logger for status updates
	messenger := ui.NewStatusMessenger(p)
	logger.GlobalLogger = messenger

	if _, err := p.Run(); err != nil {
		fmt.Printf("❌ Error running program: %v", err)
		os.Exit(1)
	}
}
