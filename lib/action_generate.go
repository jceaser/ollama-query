// **********************************************************************************************100
/*
Ollama Query - A simple command-line tool to interact with the Ollama server API.
Code to issue /api/show requests and display model details.

Created by Thomas.Cherry.gmail.com
*/

package lib

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type ResponseFromJson struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`

	// omit empty fields when unmarshaling
	Context            []int `json:",omitempty"` //this should be passed along with the generate request and will be returned in the response, but it is not always present in the response so we need to make it omitempty
	TotalDuration      int64 `json:",omitempty"`
	LoadDuration       int64 `json:",omitempty"`
	PromptEvalCount    int   `json:",omitempty"`
	PromptEvalDuration int64 `json:",omitempty"`
	EvalCount          int   `json:",omitempty"`
	EvalDuration       int64 `json:",omitempty"`
}

func (r ResponseFromJson) String2() string {
	s := fmt.Sprintf(
		"Model: %s\nCreated At: %s\nResponse: %s\nDone: %v\nContext: %v\nTotal Duration: %d\nLoad Duration: %d\nPrompt Eval Count: %d\nPrompt Eval Duration: %d\nEval Count: %d\nEval Duration: %d",
		r.Model, r.CreatedAt, r.Response, r.Done, r.Context, r.TotalDuration, r.LoadDuration, r.PromptEvalCount, r.PromptEvalDuration, r.EvalCount, r.EvalDuration)
	return s
}

func (r ResponseFromJson) String() string {
	var sb strings.Builder
	val := reflect.ValueOf(r)
	typ := reflect.TypeOf(r)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		sb.WriteString(fmt.Sprintf("%20s: %v\n", field.Name, value.Interface()))
	}

	return sb.String()
}

/*

curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt": "Why is the sky blue?"
}'

response:
{
  "model": "llama3.2",
  "created_at": "2023-08-04T08:52:19.385406455-07:00",
  "response": "The",
  "done": false
}

Last response looks like this:
{
	"model": "llama3.2",
	"created_at": "2023-08-04T19:22:45.499127Z",
	"response": "",
	"done": true,
	"context": [1, 2, 3],
	"total_duration": 10706818083,
	"load_duration": 6338219291,
	"prompt_eval_count": 26,
	"prompt_eval_duration": 130079000,
	"eval_count": 259,
	"eval_duration": 4232710000
}

*/

func GenerateText(context AppContext, args ...string) (map[string]string, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments provided. Usage: generate <model_name> <prompt>")
	}
	modelName := args[0]
	prompt := strings.Join(args[1:], " ")

	// Create the request body
	requestBody := map[string]interface{}{
		"model":  modelName,
		"prompt": prompt,
	}
	if len(context.Context) > 0 {
		requestBody["context"] = context.Context
	}

	// Convert the request body to JSON
	jsonData, err := JsonFromStruct(requestBody)
	if err != nil {
		return nil, err
	}
	jsonBytes := bytes.NewBuffer(jsonData)

	// Send the POST request to the Ollama server
	resp, err := http.Post(context.HostName+"/api/generate", "application/json", jsonBytes)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Fprintln(context.Output, strings.Repeat("*", 80))

	scanner := bufio.NewScanner(resp.Body)
	result := map[string]string{}
	for scanner.Scan() {
		line := scanner.Text()
		response, err := StructFromJson[ResponseFromJson]([]byte(line))
		if err != nil {
			Log.Warn.Printf("Error parsing response line: %v\n", err)
			continue
		}
		fmt.Fprintf(context.Output, WrapText(Codes{ESC_GREEN}, "%s"), response.Response)
		if response.Done {
			if len(response.Context) > 0 {
				jsonData, err := json.Marshal(response.Context)
				if err != nil {
					Log.Warn.Printf("Error marshaling context: %v\n", err)
				} else {
					result["context"] = string(jsonData)
				}
			}
			if context.Verbose > 0 {
				Log.Debug.Printf("%v\n", response)
			}
			fmt.Fprintf(context.Output, "\n\n")
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(context.Error, err.Error())
	}
	return result, nil
}
