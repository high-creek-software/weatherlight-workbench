package deck

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/icons"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/storage"
	"time"
)

type DeckListItem struct {
	widget.BaseWidget

	deck storage.Deck
}

func (dll *DeckListItem) CreateRenderer() fyne.WidgetRenderer {
	nameLbl := widget.NewRichTextWithText("template")
	img := widget.NewIcon(icons.CardLoadingResource)
	img.Resize(fyne.NewSize(128, 128))
	createdAtLbl := widget.NewLabel("")
	deckTypeLbl := widget.NewLabel("")

	return &deckListItemRenderer{
		dll:          dll,
		nameLbl:      nameLbl,
		img:          img,
		createdAtLbl: createdAtLbl,
		deckTypeLbl:  deckTypeLbl,
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
	dll          *DeckListItem
	nameLbl      *widget.RichText
	img          *widget.Icon
	createdAtLbl *widget.Label
	deckTypeLbl  *widget.Label
}

func (d deckListItemRenderer) Destroy() {

}

func (d deckListItemRenderer) Layout(size fyne.Size) {
	imgSize := d.img.Size()
	imgTopLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	d.img.Move(imgTopLeft)

	topLeft := fyne.NewPos(imgSize.Width+2*theme.Padding(), theme.Padding())
	nameSize := d.nameLbl.MinSize()
	d.nameLbl.Move(topLeft)
	d.nameLbl.Resize(fyne.NewSize(size.Width, nameSize.Height))

	topLeft = topLeft.Add(fyne.NewPos(0, 22))
	//createdAtSize := d.createdAtLbl.MinSize()
	d.createdAtLbl.Move(topLeft)

	topLeft = topLeft.Add(fyne.NewPos(0, 22))
	d.deckTypeLbl.Move(topLeft)
}

func (d deckListItemRenderer) MinSize() fyne.Size {
	nameSize := d.nameLbl.MinSize()
	imgSize := d.img.Size()
	height := fyne.Max(nameSize.Height, imgSize.Height) + 2*theme.Padding()
	return fyne.NewSize(imgSize.Width+nameSize.Width+3*theme.Padding(), height)
}

func (d deckListItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.nameLbl, d.img, d.createdAtLbl, d.deckTypeLbl}
}

func (d deckListItemRenderer) Refresh() {
	d.nameLbl.ParseMarkdown(fmt.Sprintf("### %s", d.dll.deck.Name))
	d.createdAtLbl.SetText(d.dll.deck.CreatedAt.Format(time.Stamp))
	d.deckTypeLbl.SetText(d.dll.deck.DeckType)
}
