package deck

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	scryfallcards "github.com/high-creek-software/goscryfall/cards"
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
	nameLbl := widget.NewRichTextFromMarkdown("Temp")
	nameLbl.Wrapping = fyne.TextWrapWord
	countLbl := widget.NewLabel("Temp")
	cardFace := widget.NewIcon(nil)
	cardFace.Resize(fyne.NewSize(128, 128))
	manaBox := container.NewHBox()
	typeLine := widget.NewLabel("")
	typeLine.Wrapping = fyne.TextWrapWord
	setName := widget.NewLabel("")
	associationLabel := widget.NewLabel("")
	sep := widget.NewSeparator()
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
		nameLbl:          nameLbl,
		countLbl:         countLbl,
		cardFace:         cardFace,
		manaBox:          manaBox,
		typeLine:         typeLine,
		setName:          setName,
		associationLabel: associationLabel,
		separator:        sep,
		setCover:         setCover,
		removeCard:       removeBtn,
		incCard:          incCard,
		decCard:          decCard,
	}

	for i := 0; i < 4; i++ {
		dr.manaImages = append(dr.manaImages, widget.NewIcon(nil))
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
}

type deckCardListItemRenderer struct {
	li               *DeckCardListItem
	nameLbl          *widget.RichText
	countLbl         *widget.Label
	cardFace         *widget.Icon
	manaBox          *fyne.Container
	manaImages       []*widget.Icon
	typeLine         *widget.Label
	setName          *widget.Label
	associationLabel *widget.Label
	separator        *widget.Separator
	setCover         *widget.Button

	removeCard *widget.Button
	incCard    *widget.Button
	decCard    *widget.Button
}

func (d *deckCardListItemRenderer) Destroy() {

}

func (d *deckCardListItemRenderer) Layout(size fyne.Size) {

	/** Get sizes **/
	assocSize := d.associationLabel.MinSize()
	nameSize := d.nameLbl.MinSize()
	manaSize := d.manaBox.MinSize()
	countSize := d.countLbl.MinSize()
	typeSize := d.typeLine.MinSize()
	setSize := d.setName.MinSize()
	incSize := d.incCard.MinSize()
	decSize := d.decCard.MinSize()
	delSize := d.removeCard.MinSize()

	/** Cover Image & Set Cover Button **/
	topLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	cardSize := d.cardFace.Size()
	d.cardFace.Move(topLeft)

	/** Setup Center **/
	nameTopLeft := topLeft.Add(fyne.NewPos(cardSize.Width+theme.Padding(), 0))
	d.nameLbl.Move(nameTopLeft)
	d.nameLbl.Resize(fyne.NewSize(size.Width-cardSize.Width-manaSize.Width-3*theme.Padding(), nameSize.Height))

	topLeft = nameTopLeft.Add(fyne.NewPos(0, nameSize.Height))
	d.typeLine.Move(topLeft)
	d.typeLine.Resize(fyne.NewSize(size.Width-cardSize.Width-2*theme.Padding(), typeSize.Height))

	topLeft = topLeft.Add(fyne.NewPos(0, typeSize.Height-6))
	d.setName.Move(topLeft)
	d.setName.Resize(setSize)

	/*** Setup bottom buttons ***/
	cardOffset := cardSize.Height
	textOffset := topLeft.Y + setSize.Height
	sepPos := fyne.NewPos(theme.Padding(), fyne.Max(cardOffset, textOffset)+theme.Padding())

	setCoverPos := fyne.NewPos(theme.Padding(), sepPos.Y+theme.Padding())
	d.setCover.Move(setCoverPos)
	d.setCover.Resize(d.setCover.MinSize())

	topLeft = fyne.NewPos(topLeft.X+8, setCoverPos.Y)
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

	/** Setup Left Side **/
	manaPos := fyne.NewPos(size.Width-theme.Padding()-manaSize.Width, theme.Padding()+6)
	d.manaBox.Move(manaPos)
	d.manaBox.Resize(fyne.NewSize(float32(20*len(d.li.manaCost)), manaSize.Height))

	assocTop := fyne.NewPos(size.Width-theme.Padding()-assocSize.Width, manaPos.Y+manaSize.Height)
	d.associationLabel.Move(assocTop)
	d.associationLabel.Resize(assocSize)
}

func (d *deckCardListItemRenderer) MinSize() fyne.Size {
	nameSize := d.nameLbl.MinSize()
	assocSize := d.associationLabel.MinSize()
	typeSize := d.typeLine.MinSize()
	countSize := d.countLbl.Size()
	setSize := d.setName.MinSize()
	manaSize := d.manaBox.MinSize()
	cardSize := d.cardFace.Size()
	incSize := d.incCard.MinSize()
	setCoverSize := d.setCover.MinSize()

	height := fyne.Max(cardSize.Height+setCoverSize.Height, nameSize.Height+typeSize.Height+setSize.Height+fyne.Max(countSize.Height, incSize.Height)) + 3*theme.Padding()

	size := fyne.NewSize(cardSize.Width+fyne.Max(fyne.Max(nameSize.Width+manaSize.Width, 200), typeSize.Width+assocSize.Width)+3*theme.Padding(), height)
	//log.Println("Min Size:", size.Width, "X", size.Height)
	return size
}

func (d *deckCardListItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.nameLbl, d.countLbl, d.cardFace, d.manaBox, d.typeLine, d.setName, d.associationLabel, d.separator, d.setCover, d.removeCard, d.incCard, d.decCard}
}

func (d *deckCardListItemRenderer) Refresh() {
	d.nameLbl.ParseMarkdown(fmt.Sprintf("#### %s", d.li.card.Card.Name))
	d.countLbl.SetText(fmt.Sprintf("  %d  ", d.li.card.Count))
	d.cardFace.SetResource(d.li.ico)
	d.typeLine.SetText(d.li.card.Card.TypeLine)
	d.setName.SetText(d.li.card.Card.SetName)

	switch d.li.card.AssociationType {
	case storage.AssociationMain:
		d.associationLabel.SetText("M")
	case storage.AssociationCommander:
		d.associationLabel.SetText("C")
	case storage.AssociationSideboard:
		d.associationLabel.SetText("S")
	default:
		d.associationLabel.SetText("")
	}

	d.handleIncDecVisibility()

	d.manaBox.RemoveAll()
	for i, cost := range d.li.manaCost {
		var ico *widget.Icon
		if i > len(d.manaImages)-1 {
			ico = widget.NewIcon(nil)
			d.manaImages = append(d.manaImages, ico)
		} else {
			ico = d.manaImages[i]
		}
		ico.SetResource(cost)
		d.manaBox.Add(ico)
	}
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
