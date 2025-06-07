package style

import "github.com/charmbracelet/lipgloss"

var (
	UserListStyle = lipgloss.NewStyle().
			Width(20).
			Padding(1, 1).
			Border(lipgloss.NormalBorder())

	ChatViewStyle = lipgloss.NewStyle().
			Padding(1, 1).
			Border(lipgloss.NormalBorder())

	SenderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	ErrorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	InfoStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("36"))
)
