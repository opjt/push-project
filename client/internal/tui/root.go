package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = 1

var (
	userListStyle = lipgloss.NewStyle().
			Width(20).
			Padding(1, 1).
			Border(lipgloss.NormalBorder())

	chatViewStyle = lipgloss.NewStyle().
			Padding(1, 1).
			Border(lipgloss.NormalBorder())

	senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
)

type state int

const (
	stateLogin state = iota
	stateChat
)

// RootModel: 전체 앱 상태를 관리, 로그인 성공 시 상태 전환 처리
type RootModel struct {
	state      state
	width      int
	height     int
	messages   []string
	chatModel  *ChatModel
	loginModel *LoginModel
	TeaProgram *tea.Program
}

func NewRootModel(login *LoginModel, chat *ChatModel) *RootModel {
	return &RootModel{
		state:      stateLogin,
		messages:   []string{},
		loginModel: login,
		chatModel:  chat,
	}
}

func (r *RootModel) Init() tea.Cmd {
	// 로그인 모델 초기화
	return r.loginModel.Init()
}

func (r *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
		return r, nil
	}
	switch r.state {
	case stateLogin:
		newModel, cmd := r.loginModel.Update(msg)
		r.loginModel = newModel.(*LoginModel)
		// 로그인 성공 시 이벤트 발생
		if r.loginModel.loggedIn {
			r.state = stateChat
			r.chatModel.Resize(r.width, r.height)
		}
		return r, cmd

	case stateChat:
		newModel, cmd := r.chatModel.Update(msg)
		r.chatModel = newModel.(*ChatModel)
		return r, cmd
	}
	return r, nil
}

func (r *RootModel) View() string {
	switch r.state {
	case stateLogin:
		return r.loginModel.View()
	case stateChat:

		return r.chatModel.View()
	}
	return ""
}
