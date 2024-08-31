package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type submissionMsg struct {
	Name        string
	Description string
	submitted   bool
}

type model struct {
	focusIndex int
	inputName  textinput.Model
	inputDesc  textarea.Model
	err        error
	submitted  bool
}

func initialModel(defValues *submissionMsg) model {
	ti := textinput.New()
	ti.Placeholder = "Add your todo name"
	ti.SetValue(defValues.Name)
	ti.Validate = func(value string) error {
		if len(value) == 0 {
			return fmt.Errorf("name is required")
		}
		return nil
	}
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	ta := textarea.New()
	ta.SetValue(defValues.Description)
	ta.Placeholder = "Add your todo description"
	ta.SetWidth(50)
	ta.SetHeight(3)

	return model{
		focusIndex: 0,
		inputName:  ti,
		inputDesc:  ta,
		err:        nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			if m.focusIndex == 0 {
				m.focusIndex = 1
				m.inputDesc.Focus()
				m.inputName.Blur()
			} else {
				m.focusIndex = 0
				m.inputName.Focus()
				m.inputDesc.Blur()
			}
		case tea.KeyEnter:
			if m.focusIndex == 0 {
				m.focusIndex = 1
				m.inputName.Blur()
				m.inputDesc.Focus()
			} else if m.focusIndex == 1 {
				m.submitted = true
				// Return submission message when form is completed
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd

	if m.focusIndex == 0 {
		m.inputName, cmd = m.inputName.Update(msg)
	} else {
		m.inputDesc, cmd = m.inputDesc.Update(msg)
	}

	return m, cmd
}

// View implements tea.Model.
func (m model) View() string {
	return fmt.Sprintf(
		"Name: %s\nDescription: %s\n\n%s",
		m.inputName.View(),
		m.inputDesc.View(),
		"(ctrl+c to quit, tab to switch, enter to submit)",
	)
}

func RenderCreateForm(defValues *submissionMsg) (subMsg submissionMsg, err error) {
	p := tea.NewProgram(initialModel(defValues))
	m, err := p.Run()
	if err != nil {
		return submissionMsg{}, err
	}

	finalModel, ok := m.(model)
	if !ok {
		return submissionMsg{}, fmt.Errorf("unexpected model type")
	}

	if finalModel.submitted {
		return submissionMsg{
			Name:        strings.TrimSpace(finalModel.inputName.Value()),
			Description: strings.TrimSpace(finalModel.inputDesc.Value()),
			submitted:   true,
		}, nil
	}

	return submissionMsg{}, nil
}
