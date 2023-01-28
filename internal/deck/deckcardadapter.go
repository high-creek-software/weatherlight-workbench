package deck

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
)

type DeckCardAdapter struct {
	cards    []storage.DeckCard
	registry *platform.Registry
	list     *widget.List

	setCover func(card storage.DeckCard)
}

func NewDeckCardAdapter(registry *platform.Registry, setCover func(card storage.DeckCard)) *DeckCardAdapter {
	return &DeckCardAdapter{registry: registry, setCover: setCover}
}

func (dca *DeckCardAdapter) SetList(list *widget.List) {
	dca.list = list
}

func (dca *DeckCardAdapter) AppendCards(cs []storage.DeckCard) {
	dca.cards = append(dca.cards, cs...)
}

func (dca *DeckCardAdapter) Clear() {
	dca.cards = nil
}

func (dca *DeckCardAdapter) Count() int {
	return len(dca.cards)
}

func (dca *DeckCardAdapter) CreateTemplate() fyne.CanvasObject {
	return NewDeckCardListItem(dca.setCover)
}

func (dca *DeckCardAdapter) UpdateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	card := dca.Item(id)
	li := co.(*DeckCardListItem)

	var mc []fyne.Resource
	costs := card.Card.ParseManaCost()
	if len(costs) > 0 {
		cost := costs[0]
		for _, c := range cost {
			func(name string) {
				mc = append(mc, dca.registry.SymbolRepo.Image(name))
			}(c)
		}
	}

	li.Update(card, mc)
	cardImgPath := li.card.Card.ImageUris.ArtCrop
	if cardImgPath == "" && len(li.card.Card.CardFaces) > 0 {
		cardImgPath = li.card.Card.CardFaces[0].ImageUris.ArtCrop
	}
	dca.registry.CardThumbnailLoader.Load(li.card.Card.Id, cardImgPath, li)
	dca.list.SetItemHeight(id, li.MinSize().Height)
}

func (dca *DeckCardAdapter) Item(id widget.ListItemID) storage.DeckCard {
	return dca.cards[id]
}
