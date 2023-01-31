package notifier

type Notifier interface {
	VerifyAction(message, actionTitle string, action func())
	ShowDialog(title, message string)
	ShowError(err error)
}
