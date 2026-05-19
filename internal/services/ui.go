package services

type UI interface {
	ShowMessage(message string)
	ShowError(message string, err error)
	ShowFatalError(message string, err error)
}
