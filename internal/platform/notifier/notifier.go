package notifier

type Notifier interface {
	ShowDialog(title, message string)
	ShowError(err error)
}
