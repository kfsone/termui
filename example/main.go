package main

import (
	"flag"
	"github.com/kfsone/termui"
	"os"
)

var useTui bool

func main() {
	flag.BoolVar(&useTui, "tui", false, "Enable the TUI interface")
	flag.Parse()

	var ui termui.UserInterfacer
	if useTui == false {
		ui = termui.NewRawConsoleUI(os.Stdin)
	} else {
		ui = termui.NewTUIUserInterface()
	}
	inputs := ui.Open()

	ui.WriteString("Ready. Enter things to be echoed or 'quit' to quit.")

LOOP:
	for command := range inputs {
		switch command {
		case "quit", "q", "exit":
			{
				ui.WriteString("You asked me to quit.")
				break LOOP
			}
		case "help", "?":
			{
				ui.WriteString("Any line of text you enter will be echoed back at you unless it is 'help', 'quit', or 'exit'.")
			}
		default:
			{
				ui.WriteString("You entered: " + command)
			}
		}
	}


	ui.Close()
}
