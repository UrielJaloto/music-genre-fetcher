package services

type UI interface {
	ShowMessage(message string)
	ShowError(message string, err error)
}
