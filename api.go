// ollama API
package main

import (
	"encoding/json"
	"fmt"
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

func (rqb requestBody) ChatRequest(content string) string {
	url := "http://localhost:11434/api/chat"
	msg := message{"user", content}
	rqb.Messages = append(rqb.Messages, msg)

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
	rqb.Messages = append(rqb.Messages, ExtractAnswer(responseBody))
	return rqb.MsgHistory()
}

func (rqb requestBody) MsgHistory() string {
	msgString := ""
	for _, msg := range rqb.Messages {
		msgString += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}
	return msgString
}
