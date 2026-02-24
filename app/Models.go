// **********************************************************************************************100
/*
Ollama Query - A simple command-line tool to interact with the Ollama server API.
Code to issue /api/show requests and display model details.

Created by Thomas.Cherry.gmail.com
*/

package app

// Model represents an individual model.
type Model struct {
	Name          string  `json:"name"`
	Model         string  `json:"model"`
	Size          int64   `json:"size"`
	Digest        string  `json:"digest"`
	Details       Details `json:"details"`
	ExpiresAt     string  `json:"expires_at"`
	SizeVRAM      int64   `json:"size_vram"`
	ContextLength int64   `json:"context_length"`
}

// Details contains additional information about the model.
type Details struct {
	ParentModel       string   `json:"parent_model"`
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

// ModelsResponse represents the overall structure of the JSON response.
type ModelsResponse struct {
	Models []Model `json:"models"`
}

type VersionResponse struct {
	Version string `json:"version"`
}
