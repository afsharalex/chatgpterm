package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"

	"github.com/afsharalex/chatgpterm/client"
)

type (
	errMsg        error
	apiReponseMsg string
)

type model struct {
	viewport viewport.Model
	textarea textarea.Model
	prompt   string
	response string
	err      error
	client   *client.Client
}

func initialModel(apiKey string) model {
	ta := textarea.New()
	ta.Placeholder = "Ask a question..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 100

	ta.SetWidth(100)
	ta.SetHeight(3)

	ta.ShowLineNumbers = false

	vp := viewport.New(100, 15)
	vp.SetContent("Welcome to ChatGPTerm! Type a question and press Enter to send.")

	ta.KeyMap.InsertNewline.SetEnabled(false)
	client := client.NewClient(apiKey)

	return model{
		textarea: ta,
		viewport: vp,
		prompt:   "",
		response: "",
		err:      nil,
		client:   client,
	}
}

func (m *model) TestMessage() tea.Msg {
	res, err := m.client.Query(m.prompt)
	if err != nil {
		log.Printf("Receieved error response: %s", err)
		return nil
	}

	// answer := wordwrap.String(res, 100)
	// m.viewport.SetContent(answer)

	// m.response = res

	return apiReponseMsg(res)
}

// TODO: Do I block here? If not, how do I let BubbleTea know that
// we've received a response and to rerender the viewport?
func (m *model) TestCallback(query string) tea.Msg {
	// m.viewport.SetContent(response)
	res, err := m.client.Query(query)
	if err != nil {
		log.Printf("Receieved error response: %s", err)
		return nil
	}

	answer := wordwrap.String(res, 100)
	m.viewport.SetContent(answer)

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
			// m.response = m.prompt
			// m.viewport.SetContent(m.response)
			// TODO: Definitely do not want to block update,
			// so what is the best way to fire off a process,
			// maybe show some loading spinner thingy and
			// let bubbletea know when we've finished?
			// My instinct says to use the tea.Cmd system.
			// Need to confirm.
			callback = m.TestMessage

			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	case apiReponseMsg:
		answer := wordwrap.String(string(msg), 100)
		m.viewport.SetContent(answer)
		// callback = nil
		m.viewport.GotoBottom()

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
	apiKey := os.Getenv("CHAT_GPT_API_KEY")
	if apiKey == "" {
		log.Fatal("CHAT_GPT_API_KEY is not set in Environment.")
	}
	p := tea.NewProgram(initialModel(apiKey))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error creating app.")
		os.Exit(1)
	}
}
