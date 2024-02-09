// ollama API
package api

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
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

func EncodeJson(rb requestBody) string {
	finalJson, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}
	return string(finalJson)
}

func ExtractAnswer(data []byte) interface{} {
	var answer map[string]map[string]interface{}
	json.Unmarshal(data, &answer)
	return answer["message"]["content"]
}

func ChatRequest(content string, model string) {
	url := "http://localhost:11434/api/chat"
	msg := message{"user", content}
	rb := requestBody{model, []message{msg}, false}

	data := strings.NewReader(EncodeJson(rb))
	response, err := http.Post(url, "", data)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	answer := ExtractAnswer(responseBody)

	fmt.Println(answer)
}
