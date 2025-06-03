package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
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

type userItem string

func (u userItem) Title() string       { return string(u) }
func (u userItem) Description() string { return "" }
func (u userItem) FilterValue() string { return string(u) }

type singleLineDelegate struct{}

func (d singleLineDelegate) Height() int                               { return 1 }
func (d singleLineDelegate) Spacing() int                              { return 0 }
func (d singleLineDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d singleLineDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	fmt.Fprint(w, " "+item.FilterValue())
}

type model struct {
	users    list.Model
	viewport viewport.Model
	textarea textarea.Model
	width    int
	height   int
	messages []string
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.SetWidth(50)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(50, 10)
	vp.SetContent("Welcome to the chat room!")

	users := []list.Item{
		userItem("alice"),
		userItem("bob"),
		userItem("carol"),
		userItem("dave"),
		userItem("eve"),
	}

	userList := list.New(users, singleLineDelegate{}, 20, 10)
	userList.Title = "Users"
	userList.SetShowStatusBar(false)
	userList.SetFilteringEnabled(false)
	userList.DisableQuitKeybindings()

	return model{
		users:    userList,
		textarea: ta,
		viewport: vp,
		messages: []string{},
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		// 유저 리스트의 padding + border 높이 계산
		userListVerticalPadding := userListStyle.GetPaddingTop() + userListStyle.GetPaddingBottom()

		// 채팅 뷰 패딩 + border 높이 계산
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

func (m model) View() string {
	userView := userListStyle.Render(m.users.View())
	rightView := chatViewStyle.Render(m.viewport.View() + "\n\n" + m.textarea.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, userView, rightView)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
