package card

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/goscryfall/cards"
	"gitlab.com/kendellfab/mtgstudio/internal/platform"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/adapter"
)

var _ adapter.Adapter[cards.Card] = (*CardAdapter)(nil)

type CardAdapter struct {
	cards    []cards.Card
	list     *widget.List
	registry *platform.Registry
}

func NewCardAdapter(registry *platform.Registry) *CardAdapter {
	return &CardAdapter{registry: registry}
}

func (ca *CardAdapter) SetList(list *widget.List) {
	ca.list = list
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
	return NewCardListItem(nil)
}

func (ca *CardAdapter) UpdateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	card := ca.Item(id)
	listItem := co.(*CardListItem)

	var mc []fyne.Resource
	sets := card.ParseManaCost()
	if len(sets) > 0 {
		cost := sets[0]
		for _, c := range cost {
			func(name string) {
				mc = append(mc, ca.registry.SymbolRepo.Image(name))
			}(c)
		}
	}

	listItem.UpdateCard(&card, mc)

	cardImgPath := card.ImageUris.ArtCrop
	if cardImgPath == "" && len(card.CardFaces) > 0 {
		cardImgPath = card.CardFaces[0].ImageUris.ArtCrop
	}

	ca.registry.CardThumbnailLoader.Load(card.Id, cardImgPath, listItem)
	ca.list.SetItemHeight(id, listItem.MinSize().Height)
}

func (ca *CardAdapter) Item(id widget.ListItemID) cards.Card {
	return ca.cards[id]
}
