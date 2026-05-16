package presentation

import (
	"fmt"
	"os"
)

type CLIUserInterface struct{}

func NewCLIUserInterface() *CLIUserInterface {
	return &CLIUserInterface{}
}

func (ui *CLIUserInterface) ShowMessage(message string) {
	fmt.Println(message)
}

func (ui *CLIUserInterface) ShowError(message string, err error) {
	if err != nil {
		fmt.Printf("ERROR: %s - %v\n", message, err)
		return
	}
	fmt.Printf("ERROR: %s\n", message)
}

func (ui *CLIUserInterface) ShowFatalError(message string, err error) {
	if err != nil {
		fmt.Printf("FATAL: %s - %v\n", message, err)
	} else {
		fmt.Printf("FATAL: %s\n", message)
	}
	os.Exit(1)
}
