package card

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/high-creek-software/goscryfall/cards"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/adapter"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
)

var _ adapter.Adapter[cards.Card] = (*CardAdapter)(nil)

type CardAdapter struct {
	cards      []cards.Card
	loader     *ansel.Ansel[string]
	symbolRepo symbol.SymbolRepo
	setLoader  func(uri string) ([]byte, error)
}

func NewCardAdapter(loader *ansel.Ansel[string], symbolRepo symbol.SymbolRepo, setLoader func(uri string) ([]byte, error)) *CardAdapter {
	return &CardAdapter{loader: loader, symbolRepo: symbolRepo, setLoader: setLoader}
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
				mc = append(mc, ca.symbolRepo.Image(name))
			}(c)
		}
	}

	//var setIcon fyne.Resource
	//if setData, err := ca.setLoader(card.Set); err == nil {
	//	setIcon = fyne.NewStaticResource(set)
	//}

	listItem.UpdateCard(&card, mc)

	cardImgPath := card.ImageUris.ArtCrop
	if cardImgPath == "" && len(card.CardFaces) > 0 {
		cardImgPath = card.CardFaces[0].ImageUris.ArtCrop
	}

	ca.loader.Load(card.Id, cardImgPath, listItem)

}

func (ca *CardAdapter) Item(id widget.ListItemID) cards.Card {
	return ca.cards[id]
}
