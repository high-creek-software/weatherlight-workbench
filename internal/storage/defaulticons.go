package storage

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

//go:embed no_icon.svg
var noIcon []byte

var NoIconResource = fyne.NewStaticResource("no_icon", noIcon)

/*
	Size: 100X25
	Font Size: 9
	Font: FuraCode Nerd Font Mono
*/
//go:embed legal.svg
var legal []byte
var LegalResource = fyne.NewStaticResource("legal", legal)

//go:embed not_legal.svg
var notLegal []byte
var NotLegalResource = fyne.NewStaticResource("not_legal", notLegal)

//go:embed restricted.svg
var restricted []byte
var RestrictedResource = fyne.NewStaticResource("restricted", restricted)

//go:embed banned.svg
var banned []byte
var BannedResource = fyne.NewStaticResource("banned", banned)
