package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	spinner  spinner.Model
	loading  bool
	prompt   string
	response string
	err      error
	client   *client.Client
}

func initialModel(apiKey string) model {
	ta := textarea.New()
	ta.Placeholder = "Ask a question..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 100

	ta.SetWidth(100)
	ta.SetHeight(3)

	ta.ShowLineNumbers = false

	vp := viewport.New(100, 15)
	vp.SetContent("Welcome to ChatGPTerm! Type a question and press Enter to send.")

	ta.KeyMap.InsertNewline.SetEnabled(false)
	client := client.NewClient(apiKey)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		textarea: ta,
		viewport: vp,
		spinner:  s,
		loading:  false,
		prompt:   "",
		response: "",
		err:      nil,
		client:   client,
	}
}

func queryChatGPT(client *client.Client, query string) tea.Cmd {
	return func() tea.Msg {
		res, err := client.Query(query)
		if err != nil {
			log.Printf("Receieved error response: %s", err)
			return apiReponseMsg(err.Error())
		}

		return apiReponseMsg(res)
	}
}

func (m *model) TestMessage() tea.Msg {
	res, err := m.client.Query(m.prompt)
	if err != nil {
		log.Printf("Receieved error response: %s", err)
		return nil
	}

	return apiReponseMsg(res)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd    tea.Cmd
		vpCmd    tea.Cmd
		spCmd    tea.Cmd
		callback tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.spinner, spCmd = m.spinner.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			// m.prompt = m.textarea.Value()

			// Set the callback Cmd to our callback
			// which will return an apiReponseMsg.
			callback = queryChatGPT(m.client, m.textarea.Value())
			m.textarea.Reset()
			m.viewport.GotoBottom()
			m.loading = true
			m.viewport.SetContent("Loading...")
		}

		// Our callback returns an apiReponseMsg
		// which we check for an update the viewport accordingly.
	case apiReponseMsg:
		m.loading = false
		answer := wordwrap.String(string(msg), 100)
		m.viewport.SetContent(answer)
		m.viewport.GotoBottom()

	case errMsg:
		m.err = msg
		return m, nil
	}

	// Return the model and a batch of Cmds
	return m, tea.Batch(tiCmd, vpCmd, spCmd, callback)
}

func (m model) View() string {
	var view string

	if m.loading {

		view = fmt.Sprintf(
			"\n%s\n\n%s\n\n\n",
			m.viewport.View(),
			m.spinner.View(),
		) + "\n\n\tPress Ctrl+c or Esc to exit"
	} else {

		view = fmt.Sprintf(
			"\n%s\n\n%s",
			m.viewport.View(),
			m.textarea.View(),
		) + "\n\n\tPress Ctrl+c or Esc to exit"
	}

	return view
}

func main() {
	apiKey := os.Getenv("CHAT_GPT_API_KEY")
	if apiKey == "" {
		log.Fatal("CHAT_GPT_API_KEY is not set in Environment.")
	}

	// 	client := client.NewClient(apiKey)

	// 	res, err := client.Query("What is Rust?")
	// 	if err != nil {
	// 		log.Fatalf("Received err: %+v", err)
	// 	}

	// 	fmt.Printf("Response: %s", res)

	// TODO: Allow user to pass a flag for non-fullscreen mode.
	p := tea.NewProgram(initialModel(apiKey), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error creating app.")
		os.Exit(1)
	}
}
