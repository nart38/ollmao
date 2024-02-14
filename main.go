package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ollama-bubble model-name")
		os.Exit(1)
	}
	llm := os.Args[1]
	p := tea.NewProgram(initialModel(llm), tea.WithAltScreen())
	p.Run()
}
