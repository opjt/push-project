package tui

import (
	"push/client/internal/tui/style"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChatModel: 채팅 화면, 상태
type ChatModel struct {
	users    list.Model
	viewport viewport.Model
	textarea textarea.Model

	width    int
	height   int
	messages []string
}

func NewChatModel() *ChatModel {
	return &ChatModel{}
}

func (m *ChatModel) InitChat(width, height int) {
	m.width = width
	m.height = height

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.SetWidth(width - 20 - gap*2 - 4)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(width-20-gap*2-4, height-3-gap-5)
	vp.SetContent("Welcome to the chat room!")

	users := []list.Item{
		style.UserItem("alice"),
		style.UserItem("bob"),
		style.UserItem("carol"),
		style.UserItem("dave"),
		style.UserItem("eve"),
	}
	userList := list.New(users, style.SingleLineDelegate{}, 20, height-6)
	userList.Title = "Users"
	userList.SetShowStatusBar(false)
	userList.SetFilteringEnabled(false)
	userList.DisableQuitKeybindings()

	m.textarea = ta
	m.viewport = vp
	m.users = userList
	m.messages = []string{}
}

func (m *ChatModel) Init() tea.Cmd {
	// No initialization needed for now
	return nil
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
		lCmd  tea.Cmd
		cmds  []tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.users, lCmd = m.users.Update(msg)
	cmds = append(cmds, taCmd, vpCmd, lCmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		userListWidth := 20

		userListVerticalPadding := userListStyle.GetPaddingTop() + userListStyle.GetPaddingBottom()
		chatViewVerticalPadding := chatViewStyle.GetPaddingTop() + chatViewStyle.GetPaddingBottom()

		m.textarea.SetWidth(m.width - userListWidth - gap*2 - 4)
		m.viewport.Width = m.width - userListWidth - gap*2 - 4

		m.viewport.Height = m.height - m.textarea.Height() - gap - chatViewVerticalPadding - 2

		m.users.SetWidth(userListWidth)
		m.users.SetHeight(m.height - userListVerticalPadding - 3)

		if len(m.messages) > 0 {
			content := strings.Join(m.messages, "\n")
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(content))
			m.viewport.GotoBottom()
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
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

	return m, tea.Batch(cmds...)
}

func (m *ChatModel) View() string {
	userView := userListStyle.Render(m.users.View())
	rightView := chatViewStyle.Render(m.viewport.View() + "\n\n" + m.textarea.View())
	return lipgloss.JoinHorizontal(lipgloss.Top, userView, rightView)
}
