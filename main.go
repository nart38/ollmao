package main

import (
	"fmt"
	"io"
	"net/http"
	"github.com/nart38/ollama-bubble/api"
	"strings"
)

func main() {
	api.ConstructData("Can you describe go lang to me?")
}
