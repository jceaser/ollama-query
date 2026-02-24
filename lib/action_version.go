// **********************************************************************************************100
/*
Ollama Query - A simple command-line tool to interact with the Ollama server API.
Code to issue /api/show requests and display model details.

Created by Thomas.Cherry.gmail.com
*/

package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// curl http://localhost:11434/api/version -> {"version":"0.1.0"}
func GetVersion(context AppContext, args ...string) (map[string]string, error) {
	resp, err := http.Get(context.HostName + "/api/version")
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var versionResponse VersionResponse

	err = json.Unmarshal(body, &versionResponse)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(context.Output, strings.Repeat("*", 80))
	fmt.Fprintf(context.Output, "Ollama Server Version: %s\n", versionResponse.Version)
	return map[string]string{"version": versionResponse.Version}, nil
}
