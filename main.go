package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nart38/ollama-bubble/lib"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ollama-bubble model-name")
		os.Exit(1)
	}
	model := os.Args[1]

	stdin := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("User Prompt: ")
		prompt, _ := stdin.ReadString('\n')
		api.ChatRequest(prompt, model)
	}
}
