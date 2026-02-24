// **********************************************************************************************100
// Ollama Query - A simple command-line tool to interact with the Ollama server API.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"

	"github.com/jceaser/ollama-query/app"
	"github.com/jceaser/ollama-query/lib"
)

// ***************************************************************************80

const (
	ollamaServerURL2     = "http://ai.local:11434"
	ollamaServerURL1     = "http://localhost:11434"
	actionableItemFormat = "%10s %-15s %-23s %s"
)

// ***************************************************************************80

type ActionableItems []ActionableItem

type ActionableItem struct {
	Name       string
	Triggers   []string
	Action     app.Action
	Parameters string
	Help       string
}

func (a ActionableItems) Triggers() []string {
	var triggers []string
	for _, item := range a {
		triggers = append(triggers, item.Triggers...)
	}
	return triggers
}

func (a *ActionableItems) UpdateAction(key string, action app.Action) int {
	for i, item := range *a {
		if item.Name == key {
			(*a)[i].Action = action
			return i
		}
	}
	return -1
}

func (a ActionableItems) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(actionableItemFormat+"\n", "Name", "Triggers", "Parameters", "Description"))
	sb.WriteString(fmt.Sprintf(actionableItemFormat+"\n", "----", "--------", "----------", "-----------"))
	for _, item := range a {
		sb.WriteString(item.String() + "\n")
	}
	return sb.String()
}

func (a ActionableItem) Matches(command string) bool {
	for _, trigger := range a.Triggers {
		if strings.HasPrefix(trigger, command) {
			return true
		}
	}
	return false
}

// print out an ActionableItem as a formatted string
func (a ActionableItem) String() string {
	triggers := strings.Join(a.Triggers, ", ")
	return fmt.Sprintf(actionableItemFormat, a.Name, triggers, a.Parameters, a.Help)
}

var actions = ActionableItems{
	{"Chat", []string{"chat"}, app.Chat, "<model> <role> <prompt>", "Chat with model"},
	{"Exit", []string{"exit", "quit"}, Exit, "", "Exit the application"},
	{"Generate", []string{"generate"}, app.GenerateText, "<name> <prompt>", "Converse using context"},
	{"Help", []string{"help", "menu"}, Exit, "", "Display this menu"},
	{"List", []string{"ls", "list", "tags"}, app.ListModels, "", "List Models"},
	{"Processes", []string{"ps", "processes"}, app.ExecutePS, "", "Execute ps command"},
	{"Show", []string{"show", "details"}, app.ShowModelDetails, "<name>", "Show Model Details"},
	{"Version", []string{"version"}, app.GetVersion, "", "Get Version"},
}

// ***************************************************************************80

func isMatch(action string, commands []string) bool {
	for _, command := range commands {
		if strings.HasPrefix(command, action) {
			return true
		}
	}
	return false
}

func Exit(context app.AppContext, args ...string) (map[string]string, error) {
	os.Exit(0)
	return nil, nil
}
func DisplayMenu(context app.AppContext, args ...string) (map[string]string, error) {
	displayMenu()
	return nil, nil
}

func displayMenu() {
	// Display the menu of actionable items
	fmt.Println()
	fmt.Println(strings.Repeat("*", 80))
	fmt.Println(actions)
	fmt.Println("Choose an option:")
}

func DrawLine() {
	fmt.Println(strings.Repeat("-", 80))
}

func askForChoice() string {
	// ask for user input and split it into command and arguments
	fmt.Print(">")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := scanner.Text()
	fmt.Println()
	return strings.TrimSpace(choice)
}

func askForCommand(line *liner.State) string {
	if name, err := line.Prompt(">"); err == nil {
		return strings.TrimSpace(name)
	}
	return ""
}

func jsonToIntArray(jsonStr string) ([]int, error) {
	var intArray []int
	err := json.Unmarshal([]byte(jsonStr), &intArray)
	if err != nil {
		lib.Log.Warn.Printf("Error parsing context as int array: %v\n", err)
	}
	return intArray, err
}

// ***********************************40

func setup_liner(line *liner.State) string {
	//set up liner for command line input with history and tab completion
	history_fn := filepath.Join(os.TempDir(), ".ollama-server_history") //used by liner

	line.SetCtrlCAborts(true)

	line.SetTabCompletionStyle(liner.TabPrints)
	line.SetCompleter(func(line string) (c []string) {
		for _, n := range actions.Triggers() {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
				return
			}
		}
		return
	})
	if f, err := os.Open(history_fn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	return history_fn
}

func saveHistory(line *liner.State, history_fn string) {
	//save the liner history
	if f, err := os.Create(history_fn); err != nil {
		fmt.Print("Error creating history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func main() {
	// Help has to be updated at run time because it references the actions variable which is
	//  initialized in main
	actions.UpdateAction("Help", DisplayMenu)

	context := app.AppContext{
		HostName: ollamaServerURL1,
		Output:   os.Stdout,
		Error:    os.Stderr,
		Context:  nil,
	}

	var initAction string
	flag.StringVar(&context.HostName, "host", ollamaServerURL2, "Ollama server host URL")
	flag.StringVar(&initAction, "action", "", "Initial action to execute. Defaults to 'help'.")
	flag.Parse()

	var rawChoice string //the raw command line input from the user, which may contain multiple commands separated by ";". We will split it up and execute each command in order. If no input is given, we will default to "help" to display the menu.

	//do initial action before asking for user input, if none given, then default to help
	initAction = strings.TrimSpace(initAction)
	if initAction != "" {
		rawChoice = initAction
	} else {
		rawChoice = "help"
	}

	line := liner.NewLiner()
	defer line.Close()
	history := setup_liner(line)

	fmt.Println(lib.WrapText(lib.Codes{lib.ESC_BOLD, lib.ESC_UNDERLINE, lib.ESC_BLUE},
		"Ollama Server Command Line"))
	fmt.Println("By Thomas.Cherry.gmail.com (https://github.com/jceaser)")
	fmt.Println()
	for { //event loop
		if len(rawChoice) > 0 {
			multipleChoices := strings.Split(rawChoice, ";")
			for _, choice := range multipleChoices {
				splitChoice := strings.Fields(strings.TrimSpace(choice))
				action := splitChoice[0]
				params := splitChoice[1:]

				found := false //true when an action has been found in the list
				// Check if the action matches any of the actionable items
				for _, a := range actions {
					if a.Matches(action) {
						metadata, err := a.Action(context, params...)
						if err != nil {
							//action reported an error, print it out
							fmt.Println(lib.WrapText(lib.Codes{lib.ESC_RED}, "Error executing action:"), err)
						}
						// If the action returns metadata, check if it contains a "context" key and
						//  update the AppContext accordingly
						if metadata != nil {
							if sessionContext, okay := metadata["context"]; okay {
								intArray, err := jsonToIntArray(sessionContext)
								if err == nil {
									context.Context = intArray
								}
							}
						}
						found = true
						break
					}
				}
				if !found {
					msg := fmt.Sprintf("Invalid option [%s] with %v.\n", action, params)
					fmt.Println(lib.WrapText(lib.Codes{lib.ESC_RED}, msg))
					displayMenu()
				}
			}
		}

		//set up for the next loop
		//rawChoice = askForChoice()
		rawChoice = askForCommand(line)

		line.AppendHistory(rawChoice)
		saveHistory(line, history)

	}
}
