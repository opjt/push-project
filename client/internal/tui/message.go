package tui

import (
	"context"
	"fmt"
	"strings"

	"push/client/internal/pkg/grpc"
	"push/client/internal/pkg/httpclient/message"
	"push/client/internal/tui/state"
	"push/client/internal/tui/style"
	"push/common/lib/logger"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 메시지 리스트용 item 타입 정의
type messageItem struct {
	title string
	desc  string
}

func (m messageItem) FilterValue() string { return m.desc }
func (m messageItem) Title() string       { return m.title }
func (m messageItem) Description() string { return m.desc }

// ChatModel 필드 변경: viewport 대신 messagesList
type ChatModel struct {
	messagesList list.Model
	textarea     textarea.Model
	width        int
	height       int
	focusArea    string // "textarea" or "messages"
	logger       *logger.Logger

	sessionClient grpc.SessionClient
	messageCh     chan grpc.Message
	user          *state.User
	sessionActive bool
}

// NewChatModel 내 메시지 리스트 초기화
func NewChatModel(logger *logger.Logger, user *state.User, client grpc.SessionClient) *ChatModel {

	// 메시지 리스트 초기 아이템 없음
	messages := []list.Item{}

	messagesList := list.New(messages, list.NewDefaultDelegate(), 0, 0)
	messagesList.Title = "Messages"
	messagesList.SetShowStatusBar(false)
	messagesList.SetFilteringEnabled(false)
	messagesList.DisableQuitKeybindings()

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)
	ta.Blur()

	return &ChatModel{
		messagesList:  messagesList,
		textarea:      ta,
		focusArea:     "messages",
		logger:        logger,
		sessionClient: client,
		messageCh:     make(chan grpc.Message),
		user:          user,
		sessionActive: false,
	}
}

func (m *ChatModel) connectSession() tea.Cmd {
	return func() tea.Msg {
		err := m.sessionClient.Connect(context.Background(), *m.user, m.messageCh)
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

/**
func (m *ChatModel) connectThousandSessions() tea.Cmd {
	return func() tea.Msg {
		sessionCount := 2000
		errCh := make(chan error, sessionCount)

		for i := 0; i < sessionCount; i++ {
			go func(i int) {
				userCopy := *m.user
				userCopy.UserId = 1
				userCopy.SessionId = fmt.Sprintf("%s_%d", m.user.SessionId, i)

				err := m.sessionClient.Connect(context.Background(), userCopy, m.messageCh)
				if err != nil {
					errCh <- fmt.Errorf("session %d connect fail: %w", i, err)
					return
				}

				errCh <- nil
			}(i)
		}

		var failed int
		for i := 0; i < sessionCount; i++ {
			if err := <-errCh; err != nil {
				m.logger.Error(err.Error())
				failed++
			}
		}

		if failed > 0 {
			return serverErrorMsg(fmt.Sprintf("%d sessions failed to connect", failed))
		}

		m.sessionActive = true
		return nil
	}
}
**/

func (m *ChatModel) listenForMessages() tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-m.messageCh
		if !ok {
			m.sessionActive = false
			return serverErrorMsg("채널이 종료되었습니다.")

		}
		if err := message.MessageReceive(msg.MsgId); err != nil {
			return serverErrorMsg(fmt.Sprintf("Receive failed: %v", err))
		}
		return incomingMessage(msg)
	}
}

// 메시지 추가 함수 수정
func (m *ChatModel) appendMessage(title, body string) {

	newItem := messageItem{title: title, desc: body}
	items := append(m.messagesList.Items(), newItem)
	m.messagesList.SetItems(items)
	m.messagesList.Select(len(items) - 1)
}

// Resize 시 메시지 리스트 크기 조절
func (m *ChatModel) Resize(width, height int) {
	m.width = width
	m.height = height

	inputHeight := 3
	m.textarea.SetWidth(m.width - 4) // 좌우 여유 마진
	m.textarea.SetHeight(inputHeight)

	m.messagesList.SetWidth(m.width - 4)
	m.messagesList.SetHeight(m.height - inputHeight - 6) // 입력창 + 여유 공간 제외

	if m.focusArea == "textarea" {
		m.textarea.Focus()
	} else {
		m.textarea.Blur()
	}
}

func (m *ChatModel) Init() tea.Cmd {

	return tea.Batch(
		m.connectSession(),
		// m.connectThousandSessions(),
		m.listenForMessages(),
	)
}

// Update 함수 메시지 리스트 업데이트 반영
func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var taCmd, mlCmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Resize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlR:
			if m.sessionActive {
				m.appendMessage("INFO", style.InfoStyle.Render("이미 서버에 연결되어 있습니다."))
				return m, nil
			}

			m.appendMessage("INFO", style.InfoStyle.Render("세션 재연결을 시도합니다."))
			m.messageCh = make(chan grpc.Message)
			return m, tea.Batch(
				m.connectSession(),
				m.listenForMessages(),
			)

		case tea.KeyTab:
			if m.focusArea == "textarea" {
				m.focusArea = "messages"
				m.textarea.Blur()
			} else {
				m.focusArea = "textarea"
				m.textarea.Focus()
			}

		case tea.KeyEnter:
			if m.focusArea == "textarea" {
				input := m.textarea.Value()
				if strings.TrimSpace(input) != "" {
					m.appendMessage("You", input)
					m.textarea.Reset()
				}
			}
		}

	case incomingMessage:
		m.appendMessage(msg.Title, msg.Body)
		return m, m.listenForMessages()
	case serverErrorMsg:
		m.appendMessage("LOG", string(msg))
		return m, nil
	}

	if m.focusArea == "messages" {
		m.messagesList, mlCmd = m.messagesList.Update(msg)
	}

	if m.focusArea == "textarea" {
		m.textarea, taCmd = m.textarea.Update(msg)
	}
	cmds = append(cmds, taCmd, mlCmd)

	return m, tea.Batch(cmds...)
}

type incomingMessage grpc.Message
type serverErrorMsg string

// View 함수 수정: viewport 대신 messagesList.View() 사용
func (m *ChatModel) View() string {
	rightView := style.ChatViewStyle.Render(
		m.messagesList.View() + "\n" + m.textarea.View(),
	)

	status := "🔴 연결 끊김"
	if m.sessionActive {
		status = "🟢 연결됨"
	}
	statusView := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Italic(true).
		Render(fmt.Sprintf("\n\n(%s) 상태: %s", m.user.Username, status))

	return rightView + statusView
}
