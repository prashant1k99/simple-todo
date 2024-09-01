// form.go
package list

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	ID   int
	Name string
	Desc string
}

func (i Item) Title() string       { return i.Name }
func (i Item) Description() string { return i.Desc }
func (i Item) FilterValue() string { return i.Name }

type model struct {
	list     list.Model
	selected bool
	id       int
}

type customDelegate struct{}

func (d customDelegate) Height() int                               { return 1 }
func (d customDelegate) Spacing() int                              { return 0 }
func (d customDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d customDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	it, ok := item.(Item)
	if !ok {
		return
	}

	prefix := "  "
	if index == m.Index() {
		prefix = "> "
	}

	fmt.Fprintf(w, "%s%d) %s\n", prefix, index+1, it.Title())
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if selectedItem, ok := m.list.SelectedItem().(Item); ok {
				m.selected = true
				m.id = selectedItem.ID
				return m, tea.Quit // Quit once the selection is made
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return "\n" + m.list.View()
}

func initialModel(defValues []Item) model {
	items := make([]list.Item, 0)
	for _, v := range defValues {
		items = append(items, v)
	}
	const defaultWidth = 20
	const listHeight = 10

	// Use the custom delegate
	l := list.New(items, customDelegate{}, defaultWidth, listHeight)
	l.Title = "Select an option"

	return model{list: l}
}

type SelectionResponse struct {
	Err      error
	Selected bool
	Item     Item
}

func RenderListItem(defaultValuesToRender []Item) SelectionResponse {
	p := tea.NewProgram(initialModel(defaultValuesToRender))

	m, err := p.Run()
	if err != nil {
		return SelectionResponse{
			Err: err,
		}
	}

	finalModel, ok := m.(model)
	if !ok {
		return SelectionResponse{
			Err: err,
		}
	}

	if finalModel.selected {
		return SelectionResponse{
			Selected: true,
			Item: Item{
				ID: finalModel.id,
			},
		}
	}

	return SelectionResponse{}
}
