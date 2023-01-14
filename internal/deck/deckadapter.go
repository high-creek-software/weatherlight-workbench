package deck

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/kendellfab/mtgstudio/internal/storage"
)

type DeckAdapter struct {
	ds     []storage.Deck
	list   *widget.List
	loader *ansel.Ansel[string]
}

func NewDeckAdapter(ds []storage.Deck, loader *ansel.Ansel[string]) *DeckAdapter {
	return &DeckAdapter{ds: ds, loader: loader}
}

func (da *DeckAdapter) SetList(list *widget.List) {
	da.list = list
}

func (da *DeckAdapter) Update(ds []storage.Deck) {
	da.ds = ds
}

func (da *DeckAdapter) Count() int {
	return len(da.ds)
}

func (da *DeckAdapter) CreateTemplate() fyne.CanvasObject {
	return widget.NewLabel("")
}

func (da *DeckAdapter) UpdateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	deck := da.Item(id)
	li := co.(*widget.Label)

	li.SetText(deck.Name)
}
func (da *DeckAdapter) Item(id widget.ListItemID) storage.Deck {
	return da.ds[id]
}
