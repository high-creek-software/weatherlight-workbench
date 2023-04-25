package deck

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	scryfallcards "github.com/high-creek-software/goscryfall/cards"
	"github.com/high-creek-software/weatherlight-workbench/internal/card"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
)

type ManagementCallback interface {
	SetCover(c storage.DeckCard)
	RemoveCard(c storage.DeckCard)
	IncCard(c storage.DeckCard)
	DecCard(c storage.DeckCard)
}

type DeckCardListItem struct {
	widget.BaseWidget

	card     storage.DeckCard
	ico      fyne.Resource
	manaCost []fyne.Resource

	callback ManagementCallback
	deckType string
}

func (li *DeckCardListItem) CreateRenderer() fyne.WidgetRenderer {
	countLbl := widget.NewLabel("Temp")
	cardListItem := card.NewCardListItem(nil)
	associationLabel := widget.NewRichTextWithText("")
	setCover := widget.NewButton("Set As Cover", func() {
		li.callback.SetCover(li.card)
	})
	removeBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		li.callback.RemoveCard(li.card)
	})
	incCard := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		li.callback.IncCard(li.card)
	})
	decCard := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		li.callback.DecCard(li.card)
	})

	dr := &deckCardListItemRenderer{
		li:               li,
		cardListItem:     cardListItem,
		countLbl:         countLbl,
		associationLabel: associationLabel,
		setCover:         setCover,
		removeCard:       removeBtn,
		incCard:          incCard,
		decCard:          decCard,
	}

	return dr
}

func NewDeckCardListItem(callback ManagementCallback, deckType string) *DeckCardListItem {
	li := &DeckCardListItem{callback: callback, deckType: deckType}
	li.ExtendBaseWidget(li)

	return li
}

func (li *DeckCardListItem) Update(card storage.DeckCard, mc []fyne.Resource) {
	li.card = card
	li.manaCost = mc
	li.Refresh()
}

func (li *DeckCardListItem) SetResource(resource fyne.Resource) {
	li.ico = resource
	li.Refresh()
}

type deckCardListItemRenderer struct {
	li               *DeckCardListItem
	cardListItem     *card.CardListItem
	countLbl         *widget.Label
	associationLabel *widget.RichText
	setCover         *widget.Button

	removeCard *widget.Button
	incCard    *widget.Button
	decCard    *widget.Button
}

func (d *deckCardListItemRenderer) Destroy() {

}

func (d *deckCardListItemRenderer) Layout(size fyne.Size) {
	/** Get sizes **/
	cardSize := d.cardListItem.MinSize()
	assocSize := d.associationLabel.MinSize()
	countSize := d.countLbl.MinSize()
	incSize := d.incCard.MinSize()
	decSize := d.decCard.MinSize()
	delSize := d.removeCard.MinSize()

	/*** Setup Association ***/
	topLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	d.associationLabel.Move(topLeft)
	d.associationLabel.Resize(assocSize)

	/*** Setup card ***/
	topLeft = topLeft.AddXY(0, 15)
	d.cardListItem.Move(topLeft)
	d.cardListItem.Resize(size)

	/*** Setup bottom buttons ***/
	cardOffset := cardSize.Height
	topLeft = topLeft.AddXY(5, cardOffset)

	d.setCover.Move(topLeft)
	d.setCover.Resize(d.setCover.MinSize())

	topLeft = topLeft.AddXY(d.setCover.MinSize().Width+2*theme.Padding(), 0)
	d.decCard.Move(topLeft)
	d.decCard.Resize(decSize)

	topLeft = topLeft.AddXY(incSize.Width, 0)
	d.countLbl.Move(topLeft)
	d.countLbl.Resize(countSize)

	topLeft = topLeft.AddXY(countSize.Width, 0)
	d.incCard.Move(topLeft)
	d.incCard.Resize(incSize)

	topLeft = fyne.NewPos(size.Width-theme.Padding()-delSize.Width, topLeft.Y)
	d.removeCard.Move(topLeft)
	d.removeCard.Resize(delSize)
}

func (d *deckCardListItemRenderer) MinSize() fyne.Size {
	cardSize := d.cardListItem.MinSize()
	setCoverSize := d.setCover.MinSize()

	return fyne.NewSize(cardSize.Width, cardSize.Height+setCoverSize.Height+4*theme.Padding()+15)
}

func (d *deckCardListItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.cardListItem, d.countLbl, d.associationLabel, d.setCover, d.removeCard, d.incCard, d.decCard}
}

func (d *deckCardListItemRenderer) Refresh() {
	d.countLbl.SetText(fmt.Sprintf("  %d  ", d.li.card.Count))

	switch d.li.card.AssociationType {
	case storage.AssociationMain:
		d.associationLabel.ParseMarkdown("**Main Deck**")
	case storage.AssociationCommander:
		d.associationLabel.ParseMarkdown("**Commander**")
	case storage.AssociationSideboard:
		d.associationLabel.ParseMarkdown("**Sideboard**")
	default:
		d.associationLabel.ParseMarkdown("")
	}

	d.handleIncDecVisibility()

	d.cardListItem.UpdateCard(&d.li.card.Card, d.li.manaCost)
	d.cardListItem.SetResource(d.li.ico)
}

func (d *deckCardListItemRenderer) handleIncDecVisibility() {
	switch d.li.deckType {
	case scryfallcards.Commander, scryfallcards.PauperCommander, scryfallcards.Brawl, scryfallcards.HistoricBrawl:
		if !d.li.card.Card.IsBasicLand() {
			d.incCard.Hide()
			d.decCard.Hide()
		} else {
			d.incCard.Show()
			d.decCard.Show()
		}
	}
}
