package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

const gap = 1

type viewState int

const (
	stateLogin viewState = iota
	stateChat
)

// RootModel: 전체 앱 상태를 관리, 로그인 성공 시 상태 전환 처리
type RootModel struct {
	state      viewState
	width      int
	height     int
	chatModel  *ChatModel
	loginModel *LoginModel
	TeaProgram *tea.Program
}

func NewRootModel(login *LoginModel, chat *ChatModel) *RootModel {
	return &RootModel{
		state:      stateLogin,
		loginModel: login,
		chatModel:  chat,
	}
}

func (r *RootModel) Init() tea.Cmd {
	// 로그인 모델 초기화
	r.chatModel.user = r.loginModel.userInfo

	return r.loginModel.Init()
}

func (r *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
	}
	switch r.state {
	case stateLogin:
		newModel, cmd := r.loginModel.Update(msg)
		r.loginModel = newModel.(*LoginModel)
		// 로그인 성공 시 이벤트 발생
		if r.loginModel.loggedIn {
			r.state = stateChat
			r.chatModel.Resize(r.width, r.height)
			return r, r.chatModel.Init()
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
