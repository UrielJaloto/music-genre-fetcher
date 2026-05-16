package services

type UserInterface interface {
	ShowMessage(message string)
	ShowError(message string, err error)
	ShowFatalError(message string, err error)
}
