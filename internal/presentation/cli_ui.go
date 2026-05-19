package presentation

import (
	"fmt"
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
