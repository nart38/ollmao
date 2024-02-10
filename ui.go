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
	ta.Focus()
	// ta.Prompt = "| "
	ta.SetWidth(160)
	ta.SetHeight(5)
	ta.KeyMap.InsertNewline.SetEnabled(false)

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

	m.viewport, vpCmd = m.viewport.Update(msg)
	m.textarea, taCmd = m.textarea.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			chatHist := m.requestbody.ChatRequest(m.textarea.Value())
			fmt.Println("=====")
			fmt.Println(len(m.requestbody.Messages))
			fmt.Println("=====")
			m.viewport.SetContent(chatHist)
			m.textarea.Reset()
			// m.viewport.GotoBottom()
		}
	}
	return m, tea.Batch(vpCmd, taCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n",
		m.viewport.View(),
		m.textarea.View(),
	)
}
