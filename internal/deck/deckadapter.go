package deck

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/kendellfab/mtgstudio/internal/platform"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/storage"
)

type DeckAdapter struct {
	ds       []storage.Deck
	list     *widget.List
	registry *platform.Registry
}

func NewDeckAdapter(ds []storage.Deck, registry *platform.Registry) *DeckAdapter {
	return &DeckAdapter{ds: ds, registry: registry}
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
	return NewDeckListItem(storage.Deck{})
}

func (da *DeckAdapter) UpdateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	deck := da.Item(id)
	li := co.(*DeckListItem)

	li.UpdateDeck(deck)
}
func (da *DeckAdapter) Item(id widget.ListItemID) storage.Deck {
	return da.ds[id]
}
