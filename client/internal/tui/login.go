package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoginModel: 로그인 화면, 상태
type LoginModel struct {
	textInput textinput.Model
	loggedIn  bool
	warning   string
}

func NewLoginModel() *LoginModel {
	ti := textinput.New()
	ti.Placeholder = "Enter username"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 20

	return &LoginModel{
		textInput: ti,
		loggedIn:  false,
		warning:   "",
	}
}

func (m *LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			username := strings.TrimSpace(m.textInput.Value())
			if len(username) < 3 {
				m.warning = "Username must be at least 3 characters."
				m.textInput.Reset()
			} else {
				m.loggedIn = true
				m.warning = ""
				return m, nil
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// 로그인 화면에서는 크기 변경 처리 안 함
	}
	return m, cmd
}

func (m *LoginModel) View() string {
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	content := "Login\n\n" + m.textInput.View()
	if m.warning != "" {
		content += "\n" + warningStyle.Render(m.warning)
	}
	content += "\n\n(Enter to login, Esc to quit)"

	return lipgloss.NewStyle().Padding(2, 4).Render(content)
}
