package resources

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

//go:embed no_icon.svg
var noIcon []byte

var NoIconResource = fyne.NewStaticResource("no_icon", noIcon)
