package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	viewport    viewport.Model
	textarea    textarea.Model
	requestbody requestBody
}

func initialModel(llm string) model {
	vp := viewport.New(160, 20)

	ta := textarea.New()
	ta.Placeholder = "Prompt"
	ta.CharLimit = 0
	ta.Focus()
	ta.KeyMap.InsertNewline.SetEnabled(false)

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
	var taCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.textarea.Focused() {
		case true:
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
				m.requestbody.Messages = append(
					m.requestbody.Messages,
					m.requestbody.ChatRequest())
				m.viewport.SetContent(m.requestbody.MsgHistory())
				m.viewport.GotoBottom()
			}
		case false:
			switch msg.String() {
			case "i":
				m.textarea.Focus()
			case "ctrl+c":
				fmt.Println(m.textarea.Value())
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.viewport.Height = int(float64(msg.Height) * 0.8) - 2
		m.viewport.Width = msg.Width - 2
		m.textarea.SetHeight(int(float64(msg.Height) * 0.2) - 2)
		m.textarea.SetWidth(msg.Width - 2)
	}

	// FIXME: Passing msg directly to viewport at the following line causing
	// key strokes like h, j, k, l to also processed by viewport.
	m.viewport, vpCmd = m.viewport.Update(msg)
	m.textarea, taCmd = m.textarea.Update(msg)
	return m, tea.Batch(vpCmd, taCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"\n%s\n%s\n",
		m.viewport.View(),
		m.textarea.View(),
	)
}
