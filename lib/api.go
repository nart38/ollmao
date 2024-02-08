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
	Model    string `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool `json:"stream"`
}

func ConstructData(content string) {
	url := "http://localhost:11434/api/chat"
	msg := message{"user", content}
	rb := requestBody{"phi", []message{msg}, false}
	fmt.Println(rb)
	finalJson, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	data := strings.NewReader(string(finalJson))
	response, err := http.Post(url, "", data)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	answer, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("data: ", string(finalJson))
	fmt.Println(string(answer))

}
