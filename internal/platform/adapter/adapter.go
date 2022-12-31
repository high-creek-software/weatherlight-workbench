package adapter

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Adapter[T any] interface {
	Count() int
	CreateTemplate() fyne.CanvasObject
	UpdateTemplate(id widget.ListItemID, co fyne.CanvasObject)
	Item(id widget.ListItemID) T
}
