package deck

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/fynecharts"
)

type DeckMetaDisplay struct {
	*fyne.Container
}

func NewDeckMetaDisplay() *DeckMetaDisplay {
	dmd := &DeckMetaDisplay{}
	manaCurve := fynecharts.NewBarChart("", nil, nil)

	dmd.Container = container.NewBorder(nil, nil, manaCurve, nil, widget.NewLabel("Deck Name"))

	return dmd
}
