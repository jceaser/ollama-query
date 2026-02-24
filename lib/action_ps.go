// **********************************************************************************************100
/*
Ollama Query - A simple command-line tool to interact with the Ollama server API.
Code to issue /api/show requests and display model details.

Created by Thomas.Cherry.gmail.com
*/

package lib

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

/*
{"models":[

	{"name":"qwen3-coder:30b",
		"model":"qwen3-coder:30b",
		"size":19228621312,
		"digest":"06c1097efce0431c2045fe7b2e5108366e43bee1b4603a7aded8f21689e90bca",
		"details":{"parent_model":"",
			"format":"gguf",
			"family":"qwen3moe",
			"families":["qwen3moe"],
			"parameter_size":"30.5B",
			"quantization_level":"Q4_K_M"
		},
		"expires_at":"2026-02-16T15:50:40.325389-05:00",
		"size_vram":12169518592,
		"context_length":8192}]}
*/
func ExecutePS(context AppContext, args ...string) (map[string]string, error) {
	resp, err := http.Get(context.HostName + "/api/ps")
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body) // Use io.ReadAll in Go 1.16+
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(context.Output, strings.Repeat("*", 80))
	fmt.Fprintln(context.Output, "Executing ps command...")
	modelsResponse, err := StructFromJson[ModelsResponse](body)
	if err != nil {
		return nil, err
	}
	if len(modelsResponse.Models) == 0 {
		//return fmt.Errorf("No models found."), nil
		fmt.Fprintln(context.Output, "No models found.")
		return nil, nil
	}

	fmt.Fprintf(context.Output, "%-30s %-30s %15s\n", "NAME", "MODEL", "EXPIRES AT")
	fmt.Fprintf(context.Output, "%-30s %-30s %15s\n", "----", "-----", "----------")
	for _, model := range modelsResponse.Models {

		//2026-02-16T15:50:40.325389-05:00",
		expiresAtTime, _ := time.Parse(time.RFC3339, model.ExpiresAt)
		fmt.Fprintf(context.Output,
			"%-30s %-30s %15s\n",
			model.Name,
			model.Model,
			expiresAtTime.Format("2006-01-02 15:04:05"))
	}
	fmt.Fprintln(context.Output)
	return nil, nil
}
