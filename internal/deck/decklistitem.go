package deck

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/icons"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
	"time"
)

type DeckListItem struct {
	widget.BaseWidget

	deck       storage.Deck
	ico        fyne.Resource
	removeFunc func(d storage.Deck)
}

func (dll *DeckListItem) CreateRenderer() fyne.WidgetRenderer {
	nameLbl := widget.NewRichTextWithText("template")
	img := widget.NewIcon(icons.CardLoadingResource)
	img.Resize(fyne.NewSize(128, 128))
	createdAtLbl := widget.NewLabel("")
	deckTypeLbl := widget.NewLabel("")
	removeBtn := widget.NewButtonWithIcon("Remove", theme.DeleteIcon(), func() {
		dll.removeFunc(dll.deck)
	})

	return &deckListItemRenderer{
		dll:          dll,
		nameLbl:      nameLbl,
		img:          img,
		createdAtLbl: createdAtLbl,
		deckTypeLbl:  deckTypeLbl,
		removeBtn:    removeBtn,
	}
}

func NewDeckListItem(deck storage.Deck, removeFunc func(d storage.Deck)) *DeckListItem {
	dll := &DeckListItem{deck: deck, removeFunc: removeFunc}
	dll.ExtendBaseWidget(dll)

	return dll
}

func (dll *DeckListItem) UpdateDeck(deck storage.Deck) {
	dll.deck = deck
	dll.Refresh()
}

func (dll *DeckListItem) SetResource(resource fyne.Resource) {
	dll.ico = resource
}

type deckListItemRenderer struct {
	dll          *DeckListItem
	nameLbl      *widget.RichText
	img          *widget.Icon
	createdAtLbl *widget.Label
	deckTypeLbl  *widget.Label
	removeBtn    *widget.Button
}

func (d deckListItemRenderer) Destroy() {

}

func (d deckListItemRenderer) Layout(size fyne.Size) {
	imgSize := d.img.Size()
	imgTopLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	d.img.Move(imgTopLeft)

	delSize := d.removeBtn.MinSize()
	delTop := fyne.NewPos(theme.Padding(), imgSize.Height+theme.Padding())
	d.removeBtn.Move(delTop)
	d.removeBtn.Resize(delSize)

	topLeft := fyne.NewPos(imgSize.Width+2*theme.Padding(), 8+theme.Padding())
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
	delSize := d.removeBtn.MinSize()
	height := fyne.Max(nameSize.Height, imgSize.Height+theme.Padding()+delSize.Height+theme.Padding()) + 2*theme.Padding()
	return fyne.NewSize(imgSize.Width+nameSize.Width+3*theme.Padding(), height)
}

func (d deckListItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.nameLbl, d.img, d.createdAtLbl, d.deckTypeLbl, d.removeBtn}
}

func (d deckListItemRenderer) Refresh() {
	d.nameLbl.ParseMarkdown(fmt.Sprintf("### %s", d.dll.deck.Name))
	d.createdAtLbl.SetText(d.dll.deck.CreatedAt.Format(time.Stamp))
	d.deckTypeLbl.SetText(d.dll.deck.DeckType)
	d.img.SetResource(d.dll.ico)
}
