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

func (rqb requestBody) ChatRequest() message {
	url := "http://localhost:11434/api/chat"

	data := strings.NewReader(EncodeJson(rqb))
	response, err := http.Post(url, "", data)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return ExtractAnswer(responseBody)
}

func (rqb requestBody) MsgHistory() string {
	var chatHist []string
	for _, msg := range rqb.Messages {
		chatHist = append(chatHist, msg.Role+": "+msg.Content)
	}
	return strings.Join(chatHist, "\n\n")
}
