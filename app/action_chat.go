// **********************************************************************************************100
/*
Ollama Query - A simple command-line tool to interact with the Ollama server API.
Code to issue /api/chat requests and handle streaming responses.

Created by Thomas.Cherry.gmail.com
*/

package app

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/jceaser/ollama-query/lib"
)

/*
Send a chat message with a streaming response.

curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {
      "role": "user",
      "content": "why is the sky blue?"
    }
  ]
}'
Response

A stream of JSON objects is returned:

{
  "model": "llama3.2",
  "created_at": "2023-08-04T08:52:19.385406455-07:00",
  "message": {
    "role": "assistant",
    "content": "The",
    "images": null
  },
  "done": false
}
Final response:

{
  "model": "llama3.2",
  "created_at": "2023-08-04T19:22:45.499127Z",
  "message": {
    "role": "assistant",
    "content": ""
  },
  "done": true,
  "total_duration": 4883583458,
  "load_duration": 1334875,
  "prompt_eval_count": 26,
  "prompt_eval_duration": 342546000,
  "eval_count": 282,
  "eval_duration": 4535599000
}
*/

type ChatResponse struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`

	// omit empty fields when unmarshaling
	Context            []int `json:",omitempty"`
	TotalDuration      int64 `json:",omitempty"`
	LoadDuration       int64 `json:",omitempty"`
	PromptEvalCount    int   `json:",omitempty"`
	PromptEvalDuration int64 `json:",omitempty"`
	EvalCount          int   `json:",omitempty"`
	EvalDuration       int64 `json:",omitempty"`
}

type Message struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images,omitempty"`
}

func Chat(context AppContext, args ...string) (map[string]string, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("not enough arguments provided. Usage: chat <model> <role> <message>")
	}

	modelName := args[0]
	prompt := []Message{{
		Role:    args[1],
		Content: strings.Join(args[2:], " "),
	}}

	fmt.Fprintln(context.Output, strings.Repeat("*", 80))
	fmt.Fprintf(context.Output, "Sending a chat message\n")

	requestBody := map[string]interface{}{
		"model":    modelName,
		"messages": prompt,
	}

	jsonData, err := lib.JsonFromStruct(requestBody)
	if err != nil {
		return nil, err
	}
	jsonBytes := bytes.NewBuffer(jsonData)

	resp, err := http.Post(context.HostName+"/api/chat", "application/json", jsonBytes)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		response, err := lib.StructFromJson[ChatResponse]([]byte(line))
		if err != nil {
			lib.Log.Warn.Printf("Error unmarshaling response line: %v\n", err)
			continue
		}
		fmt.Fprintf(context.Output, "%s", response.Message.Content)
		if response.Done {
			fmt.Fprintln(context.Output, "\nChat complete.")
			fmt.Fprintf(context.Output, "Stats:\n%v\n", response)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		lib.Log.Error.Printf("%v\n", err.Error())
	}
	return nil, nil
}
