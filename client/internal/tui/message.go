package tui

import (
	"context"
	"strings"

	"push/client/internal/pkg/grpc"
	"push/client/internal/tui/style"
	"push/common/lib"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	}
}

func (m *ChatModel) Init() tea.Cmd {
	ctx := context.Background()

	go func() {
		if err := m.sessionClient.Connect(ctx, m.userID, m.messageCh); err != nil {
			m.logger.Errorf("Connect failed: %v", err)
		}
	}()

	return m.listenForMessages()
}

func (m *ChatModel) listenForMessages() tea.Cmd {
	return func() tea.Msg {
		msg := <-m.messageCh
		return incomingMessage(msg)
	}
}

type incomingMessage string

func (m *ChatModel) Resize(width, height int) {
	m.width = width
	m.height = height

	userListWidth := 20
	userListVerticalPadding := userListStyle.GetPaddingTop() + userListStyle.GetPaddingBottom()
	chatViewVerticalPadding := chatViewStyle.GetPaddingTop() + chatViewStyle.GetPaddingBottom()

	m.textarea.SetWidth(m.width - userListWidth - gap*2 - 4)
	m.textarea.SetHeight(3)

	m.viewport.Width = m.width - userListWidth - gap*2 - 4
	m.viewport.Height = m.height - m.textarea.Height() - gap - chatViewVerticalPadding - 2

	m.users.SetWidth(userListWidth)
	m.users.SetHeight(m.height - userListVerticalPadding - 3)

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
					m.messages = append(m.messages, senderStyle.Render("You: ")+input)
					content := strings.Join(m.messages, "\n")
					m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
					m.textarea.Reset()
					m.viewport.GotoBottom()
				}
			}
		}

	case incomingMessage:
		m.messages = append(m.messages, senderStyle.Render("Server: ")+string(msg))
		content := strings.Join(m.messages, "\n")
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
		m.viewport.GotoBottom()
		return m, m.listenForMessages()
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
	userView := userListStyle.Render(m.users.View())
	rightView := chatViewStyle.Render(m.viewport.View() + "\n\n" + m.textarea.View())
	return lipgloss.JoinHorizontal(lipgloss.Top, userView, rightView)
}
