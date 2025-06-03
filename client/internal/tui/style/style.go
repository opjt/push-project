package style

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type UserItem string

func (u UserItem) Title() string       { return string(u) }
func (u UserItem) Description() string { return "" }
func (u UserItem) FilterValue() string { return string(u) }

type SingleLineDelegate struct{}

func (d SingleLineDelegate) Height() int                               { return 1 }
func (d SingleLineDelegate) Spacing() int                              { return 0 }
func (d SingleLineDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d SingleLineDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	fmt.Fprint(w, " "+item.FilterValue())
}
