package card

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/high-creek-software/goscryfall/cards"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/adapter"
)

var _ adapter.Adapter[cards.Card] = (*CardAdapter)(nil)

type CardAdapter struct {
	cards  []cards.Card
	loader *ansel.Ansel[string]
}

func NewCardAdapter(loader *ansel.Ansel[string]) *CardAdapter {
	return &CardAdapter{loader: loader}
}

func (ca *CardAdapter) AppendCards(cs []cards.Card) {
	ca.cards = append(ca.cards, cs...)
}

func (ca *CardAdapter) Clear() {
	ca.cards = nil
}

func (ca *CardAdapter) Count() int {
	if ca.cards == nil {
		return 0
	}

	return len(ca.cards)
}

func (ca *CardAdapter) CreateTemplate() fyne.CanvasObject {
	return NewCardListITem(nil)
}

func (ca *CardAdapter) UpdateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	card := ca.Item(id)
	listItem := co.(*CardListItem)
	listItem.UpdateCard(&card)
	ca.loader.Load(card.Id, card.ImageUris.ArtCrop, listItem)
}

func (ca *CardAdapter) Item(id widget.ListItemID) cards.Card {
	return ca.cards[id]
}
