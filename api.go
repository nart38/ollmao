// ollama API
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type requestBody struct {
	Llm      string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type responseBody struct {
	Llm             string  `json:"model"`
	Created_at      string  `json:"created_at"`
	Message         message `json:"message"`
	Done            bool    `json:"done"`
	TotalDuration   uint    `json:"total_duration"`
	LoadDuration    uint    `json:"load_duration"`
	PromptEvalCount uint    `json:"prompt_eval_count"`
	PrompEvalDur    uint    `json:"prompt_eval_duration"`
	EvalCount       uint    `json:"eval_count"`
	EvalDur         uint    `json:"eval_duration"`
}

func EncodeJson(rb requestBody) string {
	finalJson, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}
	return string(finalJson)
}

func ExtractAnswer(data []byte) message {
	var answer responseBody
	json.Unmarshal(data, &answer)
	return answer.Message
}

func (rqb *requestBody) ChatRequest() (err error) {
	url := "http://localhost:11434/api/chat"

	data := strings.NewReader(EncodeJson(*rqb))
	response, err := http.Post(url, "", data)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	rqb.Messages = append(rqb.Messages, ExtractAnswer(responseBody))

	return nil
}

func (rqb requestBody) MsgHistory() string {
	var chatHist []string
	var role, content string
	for _, msg := range rqb.Messages {
		switch msg.Role {
		case "user":
			role = userStyle.Render("User:\n")
			content = userMsgStyle.Render(msg.Content)
		case "assistant":
			role = assistantStyle.Render("AI Assistant:\n")
			content = assistantMsgStyle.Render(msg.Content)
		}
		chatHist = append(chatHist, role + content)
	}
	return strings.Join(chatHist, "\n\n")
}
