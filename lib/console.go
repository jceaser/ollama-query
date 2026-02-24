// **********************************************************************************************100
/*
Stuff I copy from project to project to handle console output, like ANSI color codes, cursor
control, and text formatting.

created by Thomas.Cherry.gmail.com
*/
package lib

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// common terminal escape codes for formatting text and controlling the terminal
const (
	ESC_SAVE_SCREEN    = "?47h"
	ESC_RESTORE_SCREEN = "?47l"

	ESC_SAVE_CURSOR    = "s"
	ESC_RESTORE_CURSOR = "u"

	ESC_CURSOR_ON  = "?25h"
	ESC_CURSOR_OFF = "?25l"

	ESC_CLEAR_SCREEN = "2J"
	ESC_CLEAR_LINE   = "2K"
)

// Text formatting codes
type Code string

const (
	ESC_RESET         Code = "0"
	ESC_BOLD          Code = "1"
	ESC_FAINT              = "2"
	ESC_ITALIC             = "3" // shows up as inverse colors in some terminals, like iterm2, but does work in kitty
	ESC_UNDERLINE          = "4"
	ESC_BLINK              = "5"
	ESC_REVERSE            = "7"
	ESC_CONCEAL            = "8"
	ESC_STRIKETHROUGH      = "9"

	// Additional color codes
	ESC_BLACK   = "30"
	ESC_RED     = "31"
	ESC_GREEN   = "32"
	ESC_YELLOW  = "33"
	ESC_BLUE    = "34"
	ESC_MAGENTA = "35"
	ESC_CYAN    = "36"
	ESC_WHITE   = "37"
	ESC_DEFAULT = "39"

	// Background color codes
	ESC_BLACK_BG   = "40"
	ESC_RED_BG     = "41"
	ESC_GREEN_BG   = "42"
	ESC_YELLOW_BG  = "43"
	ESC_BLUE_BG    = "44"
	ESC_MAGENTA_BG = "45"
	ESC_CYAN_BG    = "46"
	ESC_WHITE_BG   = "47"
	ESC_DEFAULT_BG = "49"
)

type Codes []Code

// Strings returns a slice of strings representing the codes.
func (c Codes) Strings() []string {
	var strCodes []string
	for _, code := range c {
		strCodes = append(strCodes, string(code))
	}
	return strCodes
}

// WrapText wraps the given text with the specified ANSI color code.
func WrapText(colorCodes Codes, text string) string {
	codes := strings.Join(colorCodes.Strings(), ";")
	offCodes := []string{}
	for _, code := range colorCodes {
		// convert code to a number, and then add 20 to them
		if num, err := strconv.Atoi(string(code)); err == nil {
			if num >= 30 && num <= 37 {
				offCodes = append(offCodes, ESC_DEFAULT)
			} else if num >= 40 && num <= 47 {
				offCodes = append(offCodes, ESC_DEFAULT_BG)
			} else {
				if num == 1 {
					offCodes = append(offCodes, "22") // Bold off is 22
				} else {
					offCodes = append(offCodes, strconv.Itoa(num+20))
				}
			}
		}
	}
	offCodesStr := strings.Join(offCodes, ";")

	//fmt.Printf("Wrapping text '%s' with codes [%s], off codes [%s]\n", text, codes, offCodesStr)

	return fmt.Sprintf("\033[%sm%s\033[%sm", codes, text, offCodesStr)
}

/*
ESC(0
lqqqqk
x    x
mqqqqj
ESC(B
*/
func boxText(text string) string {
	var sb strings.Builder
	sb.WriteString("\033(0\n")
	lines := strings.Split(text, "\n")
	maxLength := 0
	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = len(line) + 2 // add some padding
		}
	}
	sb.WriteString("l" + strings.Repeat("q", maxLength) + "k\n")
	for _, line := range lines {
		sb.WriteString("x \033(B" + line + strings.Repeat(" ", maxLength-len(line)-2) + "\033(0 x\n")
	}
	sb.WriteString("m" + strings.Repeat("q", maxLength) + "j\n")
	sb.WriteString("\033(B\n")
	return sb.String()
}

func _init() {

	fmt.Println("Before " + WrapText(Codes{ESC_RED, ESC_BOLD, ESC_UNDERLINE}, "multi") + " After")
	fmt.Println("Before " + WrapText(Codes{ESC_RED, ESC_ITALIC, ESC_BLINK}, "multi") + " After")
	fmt.Println("Before " + WrapText(Codes{ESC_GREEN, ESC_REVERSE, ESC_CONCEAL}, "multi") + " After")
	fmt.Println("----")
	fmt.Println("Before " + WrapText(Codes{ESC_RED, ESC_BOLD}, "bold") + " After")
	fmt.Println("Before " + WrapText(Codes{ESC_GREEN, ESC_FAINT}, "faint") + " After")
	fmt.Println("Before " + WrapText(Codes{ESC_GREEN, ESC_ITALIC}, "italic") + " After")
	fmt.Println("Before " + WrapText(Codes{ESC_YELLOW, ESC_UNDERLINE}, "underline") + " After")
	fmt.Println("Before " + WrapText(Codes{ESC_BLUE, ESC_BLINK}, "blink") + " After")
	fmt.Println("Before " + WrapText(Codes{ESC_MAGENTA, ESC_REVERSE}, "reverse") + " After")
	fmt.Println("Before " + WrapText(Codes{ESC_CYAN, ESC_CONCEAL}, "conceal") + " After")
	fmt.Println("Before " + WrapText(Codes{ESC_WHITE, ESC_STRIKETHROUGH}, "strike") + " After")
	fmt.Println("----")
	fmt.Println("Before " + WrapText(Codes{ESC_RED, ESC_BLUE_BG, ESC_STRIKETHROUGH}, "strike") + " After")

	fmt.Println("----")

	fmt.Println(WrapText(Codes{ESC_RED}, "red") + ", " + WrapText(Codes{ESC_BLUE}, "blue") + ".")
	fmt.Println(WrapText(Codes{ESC_MAGENTA}, "MAGENTA") + ", " + WrapText(Codes{ESC_BLUE}, "blue") + ".")
	fmt.Println(WrapText(Codes{ESC_GREEN, ESC_UNDERLINE}, "Green/Underline") + ".")
	fmt.Println("end of line.")

	fmt.Printf("\033(0\n")
	fmt.Printf("lqqqqqqqqk\n")
	fmt.Printf("x        x\n")
	fmt.Printf("x  \033(B%s\033(0  x\n", "text")
	fmt.Printf("x        x\n")
	fmt.Printf("mqqqqqqqqj\n")
	fmt.Printf("\033(B\n")

	fmt.Print(boxText("This is\na box with text inside\nit."))
	//fmt.Print(boxText(WrapText(Codes{ESC_RED, ESC_BLINK}, "This is\na box with text inside\nit.")))

	os.Exit(0)
}
