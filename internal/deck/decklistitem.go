package deck

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/storage"
)

type DeckListItem struct {
	widget.BaseWidget

	deck storage.Deck
}

func (dll *DeckListItem) CreateRenderer() fyne.WidgetRenderer {
	nameLbl := widget.NewLabel("template")

	return &deckListItemRenderer{
		dll:     dll,
		nameLbl: nameLbl,
	}
}

func NewDeckListItem(deck storage.Deck) *DeckListItem {
	dll := &DeckListItem{deck: deck}
	dll.ExtendBaseWidget(dll)

	return dll
}

func (dll *DeckListItem) UpdateDeck(deck storage.Deck) {
	dll.deck = deck
	dll.Refresh()
}

type deckListItemRenderer struct {
	dll     *DeckListItem
	nameLbl *widget.Label
}

func (d deckListItemRenderer) Destroy() {

}

func (d deckListItemRenderer) Layout(size fyne.Size) {
	topLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	nameSize := d.nameLbl.MinSize()
	d.nameLbl.Move(topLeft)
	d.nameLbl.Resize(fyne.NewSize(size.Width, nameSize.Height))
}

func (d deckListItemRenderer) MinSize() fyne.Size {
	nameSize := d.nameLbl.MinSize()
	return fyne.NewSize(nameSize.Width+2*theme.Padding(), nameSize.Height+2*theme.Padding())
}

func (d deckListItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.nameLbl}
}

func (d deckListItemRenderer) Refresh() {
	d.nameLbl.SetText(d.dll.deck.Name)
}
