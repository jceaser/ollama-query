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
)

/*
{"models":[
{"name":"codellama:7b",
	"model":"codellama:7b",
	"modified_at":"2026-02-15T11:11:25.324211064-05:00",
	"size":3825910662,
	"digest":"8fdf8f752f6e80de33e82f381aba784c025982752cd1ae9377add66449d2225f",
	"details":{"parent_model":"",
		"format":"gguf",
		"family":"llama",
		"families":null,
		"parameter_size":"7B",
		"quantization_level":"Q4_0"}},
{"name":"llama3.1:latest","model":"llama3.1:latest","modified_at":"2026-02-11T19:31:41.775751914-05:00","size":4920753328,"digest":"46e0c10c039e019119339687c3c1757cc81b9da49709a3b3924863ba87ca666e",
	"details":{"parent_model":"","format":"gguf","family":"llama","families":["llama"],"parameter_size":"8.0B","quantization_level":"Q4_K_M"}}]}
*/

func ListModels(context AppContext, args ...string) (map[string]string, error) {
	resp, err := http.Get(context.HostName + "/api/tags")
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(context.Output, strings.Repeat("*", 80))
	fmt.Fprintln(context.Output, "Listing models...")
	modelsResponse, err := StructFromJson[ModelsResponse](body)
	if err != nil {
		return nil, err
	}
	if len(modelsResponse.Models) == 0 {
		fmt.Fprintln(context.Output, "No models available.")
		return nil, nil
	}
	fmt.Fprintln(context.Output, "Models Available:")
	format1 := "%-30s %-7s %-7s %s\n"
	format2 := "%-30s %-7s %-7s %d\n"
	fmt.Fprintf(context.Output, format1, "    ", "Param", "Quan", "Size")
	fmt.Fprintf(context.Output, format1, "Name", "Size", "Level", "(MB)")
	fmt.Fprintf(context.Output, format1, "----", "-----", "----", "----")
	for _, model := range modelsResponse.Models {
		fmt.Fprintf(context.Output,
			format2,
			model.Name,
			model.Details.ParameterSize,
			model.Details.QuantizationLevel,
			model.Size/1024/1024)
	}
	fmt.Fprintln(context.Output)
	return nil, nil
}
