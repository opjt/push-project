package tui

import (
	"context"
	"fmt"
	"strings"

	"push/client/internal/pkg/grpc"
	"push/client/internal/tui/style"
	"push/common/lib"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatModel struct {
	users     list.Model
	viewport  viewport.Model
	textarea  textarea.Model
	width     int
	height    int
	messages  []string
	focusArea string // "textarea" or "users"
	logger    lib.Logger

	sessionClient grpc.SessionClient
	messageCh     chan string
	userID        string
	sessionActive bool
}

func NewChatModel(logger lib.Logger, client grpc.SessionClient) *ChatModel {
	users := []list.Item{
		style.UserItem("alice"),
		style.UserItem("bob"),
		style.UserItem("carol"),
		style.UserItem("dave"),
	}

	userList := list.New(users, style.SingleLineDelegate{}, 20, 0)
	userList.Title = "Users"
	userList.SetShowStatusBar(false)
	userList.SetFilteringEnabled(false)
	userList.DisableQuitKeybindings()

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)
	ta.Blur()

	vp := viewport.New(0, 0)
	vp.SetContent("Welcome to the chat room!")

	return &ChatModel{
		users:         userList,
		viewport:      vp,
		textarea:      ta,
		messages:      []string{},
		focusArea:     "textarea",
		width:         0,
		height:        0,
		logger:        logger,
		sessionClient: client,
		messageCh:     make(chan string),
		userID:        "client1", // 실제 사용자 ID로 교체 필요
		sessionActive: false,
	}
}
func (m *ChatModel) Init() tea.Cmd {

	return tea.Batch(
		m.connectSession(),
		m.listenForMessages(),
	)
}

type serverErrorMsg string

func (m *ChatModel) connectSession() tea.Cmd {
	return func() tea.Msg {
		err := m.sessionClient.Connect(context.Background(), m.userID, m.messageCh)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.Unavailable {
				return serverErrorMsg("서버에 연결할 수 없습니다.")
			}
			return serverErrorMsg(fmt.Sprintf("Connect failed: %v", err))
		}
		m.sessionActive = true
		return nil
	}
}
func (m *ChatModel) listenForMessages() tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-m.messageCh
		if !ok {
			m.sessionActive = false
			return serverErrorMsg("채널이 종료되었습니다.")

		}
		return incomingMessage(msg)
	}
}

type incomingMessage string

func (m *ChatModel) appendMessage(msg string) {
	m.messages = append(m.messages, msg)
	content := strings.Join(m.messages, "\n")
	m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
	m.viewport.GotoBottom()
}

func (m *ChatModel) Resize(width, height int) {
	// TODO: 매직넘버 줄이고 레이아웃 구성.
	m.width = width
	m.height = height

	userListWidth := 20
	userListVerticalPadding := style.UserListStyle.GetPaddingTop() + style.UserListStyle.GetPaddingBottom()
	chatViewVerticalPadding := style.ChatViewStyle.GetPaddingTop() + style.ChatViewStyle.GetPaddingBottom()

	m.textarea.SetWidth(m.width - userListWidth - gap*2 - 4)
	m.textarea.SetHeight(3)

	m.viewport.Width = m.width - userListWidth - gap*2 - 4
	m.viewport.Height = m.height - m.textarea.Height() - gap - chatViewVerticalPadding - 4

	m.users.SetWidth(userListWidth)
	m.users.SetHeight(m.height - userListVerticalPadding - 5)

	if m.focusArea == "textarea" {
		m.textarea.Focus()
	} else {
		m.textarea.Blur()
	}

	if len(m.messages) > 0 {
		content := strings.Join(m.messages, "\n")
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
		m.viewport.GotoBottom()
	}
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
		lCmd  tea.Cmd
		cmds  []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Resize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlR:
			m.appendMessage(style.InfoStyle.Render("세션 재연결을 시도합니다."))
			m.messageCh = make(chan string) // 새 채널로 갱신
			return m, tea.Batch(
				m.connectSession(),
				m.listenForMessages(),
			)

		case tea.KeyTab:
			if m.focusArea == "textarea" {
				m.focusArea = "users"
				m.textarea.Blur()
			} else {
				m.focusArea = "textarea"
				m.textarea.Focus()
			}

		case tea.KeyEnter:
			if m.focusArea == "textarea" {
				input := m.textarea.Value()
				if strings.TrimSpace(input) != "" {
					m.messages = append(m.messages, style.SenderStyle.Render("You: ")+input)
					content := strings.Join(m.messages, "\n")
					m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
					m.textarea.Reset()
					m.viewport.GotoBottom()
				}
			}
		}

	case incomingMessage:
		m.appendMessage(style.SenderStyle.Render("Server: ") + string(msg))

		return m, m.listenForMessages()
	case serverErrorMsg:
		m.appendMessage(style.ErrorStyle.Render("[LOG] ") + string(msg))
		return m, nil
	}

	if m.focusArea == "textarea" {
		m.textarea, taCmd = m.textarea.Update(msg)
	} else {
		m.users, lCmd = m.users.Update(msg)
	}

	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, taCmd, lCmd, vpCmd)

	return m, tea.Batch(cmds...)
}

func (m *ChatModel) View() string {
	userView := style.UserListStyle.Render(m.users.View())
	rightView := style.ChatViewStyle.Render(m.viewport.View() + "\n\n" + m.textarea.View())

	status := "🔴 연결 끊김"
	if m.sessionActive {
		status = "🟢 연결됨"
	}
	statusView := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Italic(true).
		Render(fmt.Sprintf("\n\n상태: %s", status))

	return lipgloss.JoinHorizontal(lipgloss.Top, userView, rightView) + statusView
}
