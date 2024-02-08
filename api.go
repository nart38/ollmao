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
	role    string
	content string
}

type requestBody struct {
	model    string
	messages []message
	stream   bool
}

func ConstructData(content string) {
	url := "http://localhost:11434/api/chat"
	msg := message{"user", content}
	rb := requestBody{"phi", []message{msg}, false}
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

	fmt.Println(answer)

}
