package deck

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/fynecharts"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
)

type DeckDisplay struct {
	*fyne.Container
	canvas fyne.Canvas

	registry     *platform.Registry
	selectedDeck storage.Deck

	cardList    *widget.List
	cardAdapter *DeckCardAdapter
	cardTabs    *container.DocTabs
	manaCurve   *fynecharts.BarChart
	deckImage   *widget.Icon
	deckNameLbl *widget.Label
}

func NewDeckMetaDisplay(canvas fyne.Canvas, registry *platform.Registry, deck storage.Deck) *DeckDisplay {
	dd := &DeckDisplay{registry: registry, selectedDeck: deck}
	manaCurve := fynecharts.NewBarChart(dd.canvas, "", nil, nil)

	dd.Container = container.NewBorder(nil, nil, manaCurve, nil, widget.NewLabel("Deck Name"))

	return dd
}
