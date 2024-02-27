package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Defining lipgloss styling.
var (
	baseStyle lipgloss.Style = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("119"))

	userStyle lipgloss.Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("118")).
			Bold(true)
	userMsgStyle lipgloss.Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("119"))

	assistantStyle lipgloss.Style = lipgloss.NewStyle().
			Foreground(lipgloss.Color("141")).
			Bold(true)
	assistantMsgStyle lipgloss.Style = lipgloss.NewStyle().
				Foreground(lipgloss.NoColor{})
)

type model struct {
	viewport    viewport.Model
	textarea    textarea.Model
	requestbody requestBody
}

func initialModel(llm string) model {
	vp := viewport.New(160, 20)

	ta := textarea.New()
	ta.Placeholder = "Enter your prompt."
	ta.CharLimit = 0
	ta.KeyMap.InsertNewline.SetEnabled(false)
	ta.ShowLineNumbers = false
	ta.Prompt = " "
	ta.Focus()

	// Following line does not work due to: https://github.com/charmbracelet/bubbletea/issues/800
	// ta.KeyMap.InsertNewline.SetKeys("shift+enter")

	rqb := requestBody{
		Llm:      llm,
		Messages: []message{},
		Stream:   false,
	}

	return model{
		viewport:    vp,
		textarea:    ta,
		requestbody: rqb,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var vpCmd tea.Cmd
	vpMsg := msg
	var taCmd tea.Cmd
	taMsg := msg

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.textarea.Focused() {
		case true:
			// We are setting vpMsg to nil to prevent viewport to process keys
			// like h, j, k, l while typing in textarea
			vpMsg = nil
			switch msg.String() {
			case "esc":
				m.textarea.Blur()
			case "ctrl+c":
				fmt.Println(m.textarea.Value())
				return m, tea.Quit
			case "enter":
				m.requestbody.Messages = append(
					m.requestbody.Messages,
					message{"user", m.textarea.Value()})
				m.viewport.SetContent(m.requestbody.MsgHistory())
				m.textarea.Reset()
				m.textarea.Placeholder = "Waiting for the model's response..."
				err := m.requestbody.ChatRequest()
				if err != nil {
					panic(err)
				}
				m.viewport.SetContent(m.requestbody.MsgHistory())
				m.viewport.GotoBottom()
				m.textarea.Placeholder = "Enter your prompt."
			}
		case false:
			taMsg = nil
			switch msg.String() {
			case "i":
				m.textarea.Focus()
			case "ctrl+c", "q":
				fmt.Println(m.textarea.Value())
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.viewport.Height = int(float64(msg.Height)*0.8) - 2
		m.viewport.Width = msg.Width - 2
		m.textarea.SetHeight(int(float64(msg.Height)*0.2) - 2)
		m.textarea.SetWidth(msg.Width - 2)
	}

	// FIXME: Passing msg directly to viewport at the following line causing
	// key strokes like h, j, k, l to also processed by viewport.
	m.viewport, vpCmd = m.viewport.Update(vpMsg)
	m.textarea, taCmd = m.textarea.Update(taMsg)
	return m, tea.Batch(vpCmd, taCmd)
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		baseStyle.Render(m.viewport.View()),
		baseStyle.Render(m.textarea.View()),
	)
}
