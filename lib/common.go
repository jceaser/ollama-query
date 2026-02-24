// **********************************************************************************************100
/*
Ollama Query - A simple command-line tool to interact with the Ollama server API.
Code to issue /api/show requests and display model details.

Created by Thomas.Cherry.gmail.com
*/

package lib

import (
	"encoding/json"
	"os"
)

// action function
type Action func(context AppContext, args ...string) (map[string]string, error)

type AppContext struct {
	HostName string
	Output   *os.File
	Error    *os.File
	Context  []int
	Verbose  int
}

/**************************************/
// MARK: - Marshal functions

func StructFromJson[K any](someJson []byte) (K, error) {
	var value K
	err := json.Unmarshal(someJson, &value)
	return value, err
}

func PrettyJsonFromStruct[K any](data K, pretty bool) ([]byte, error) {
	result, err := json.MarshalIndent(data, "", "    ")
	return result, err
}

func JsonFromStruct[K any](data K) ([]byte, error) {
	result, err := json.Marshal(data)
	return result, err
}
