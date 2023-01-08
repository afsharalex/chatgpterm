package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type model struct {
	viewport viewport.Model
	textarea textarea.Model
	prompt   string
	response string
	err      error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Ask a question..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 300

	ta.SetWidth(300)
	ta.SetHeight(3)

	ta.ShowLineNumbers = false

	vp := viewport.New(300, 5)
	vp.SetContent("Welcome to ChatGPTerm! Type a question and press Enter to send.")

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea: ta,
		viewport: vp,
		prompt:   "",
		response: "",
		err:      nil,
	}
}

func (m *model) TestCallback(response string) tea.Cmd {
	m.viewport.SetContent(response)
	return nil
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd    tea.Cmd
		vpCmd    tea.Cmd
		callback tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.prompt = m.textarea.Value()
			m.response = m.prompt
			// m.viewport.SetContent(m.response)
			callback = m.TestCallback(m.response)

			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd, callback)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error creating app.")
		os.Exit(1)
	}
}
