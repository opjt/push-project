package tui

import (
	"push/client/internal/pkg/httpclient/auth"
	"push/client/internal/tui/state"
	"push/linker/dto"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// 인증 결과 메시지
type userValidatedMsg struct{}
type userInvalidMsg error

type LoginModel struct {
	textInput textinput.Model
	loggedIn  bool
	userInfo  *state.User
	warning   string
	loading   bool
	spinner   spinner.Model
}

func NewLoginModel(user *state.User) *LoginModel {
	ti := textinput.New()
	ti.Placeholder = "Enter username"
	ti.Focus()
	ti.CharLimit = 32
	ti.Width = 20

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	return &LoginModel{
		textInput: ti,
		spinner:   s,
		userInfo:  user,
	}
}

func (m *LoginModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.spinner.Tick,
	)
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var textInputCmd tea.Cmd
	m.textInput, textInputCmd = m.textInput.Update(msg)
	cmds = append(cmds, textInputCmd)
	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading {
			// API 호출 중 입력 무시
			return m, tea.Batch(cmds...)
		}
		switch msg.Type {
		case tea.KeyEnter:
			username := strings.TrimSpace(m.textInput.Value())
			m.warning = ""
			m.loading = true
			m.textInput.Blur()
			return m, tea.Batch(validateUserCmd(m, username), m.spinner.Tick)

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case userValidatedMsg:
		m.loggedIn = true
		m.loading = false
		m.warning = ""
		return m, nil

	case userInvalidMsg:
		m.loggedIn = false
		m.loading = false
		m.textInput.Focus()
		m.warning = "Invalid username: " + msg.Error()
		m.textInput.Reset()
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m *LoginModel) View() string {
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	content := "Login\n\n" + m.textInput.View()

	if m.loading {
		content += "\n\n" + m.spinner.View() + " Checking username..."
	}
	if m.warning != "" {
		content += "\n" + warningStyle.Render(m.warning)
	}
	content += "\n\n(Enter to login, Esc to quit)"

	return lipgloss.NewStyle().Padding(2, 4).Render(content)
}

func validateUserCmd(m *LoginModel, username string) tea.Cmd {
	return func() tea.Msg {
		req := dto.AuthLoginReq{
			UserId: username,
		}

		res, err := auth.AuthLogin(req)
		if err != nil {
			return userInvalidMsg(err)
		}

		// *m.userInfo = state.User{UserId: res.UserId}
		m.userInfo.UserId = res.UserId

		return userValidatedMsg{}
	}
}
