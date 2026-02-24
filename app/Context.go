// **********************************************************************************************100
/*
Ollama Query - A simple command-line tool to interact with the Ollama server API.
Code to issue /api/show requests and display model details.

Created by Thomas.Cherry.gmail.com
*/

package app

import "os"

/**************************************/
// MARK: - Marshal functions

// action function
type Action func(context AppContext, args ...string) (map[string]string, error)

type AppContext struct {
	HostName string
	Output   *os.File
	Error    *os.File
	Context  []int
	Verbose  int
}
