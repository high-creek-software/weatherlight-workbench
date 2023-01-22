package deck

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/storage"
)

type DeckCardListItem struct {
	widget.BaseWidget

	card     storage.DeckCard
	ico      fyne.Resource
	manaCost []fyne.Resource

	setCover func(card storage.DeckCard)
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
		li.setCover(li.card)
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
	}

	for i := 0; i < 4; i++ {
		dr.manaImages = append(dr.manaImages, widget.NewIcon(nil))
	}

	return dr
}

func NewDeckCardListItem(setCover func(card storage.DeckCard)) *DeckCardListItem {
	li := &DeckCardListItem{setCover: setCover}
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
}

func (d deckCardListItemRenderer) Destroy() {

}

func (d deckCardListItemRenderer) Layout(size fyne.Size) {
	topLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	cardSize := d.cardFace.Size()
	d.cardFace.Move(topLeft)

	assocSize := d.associationLabel.MinSize()
	nameSize := d.nameLbl.MinSize()
	nameTopLeft := topLeft.Add(fyne.NewPos(cardSize.Width+theme.Padding(), 8))
	d.nameLbl.Move(nameTopLeft)
	d.nameLbl.Resize(fyne.NewSize(size.Width-cardSize.Width-assocSize.Width-3*theme.Padding(), nameSize.Height))

	assocTop := fyne.NewPos(size.Width-theme.Padding()-assocSize.Width, theme.Padding())
	d.associationLabel.Move(assocTop)
	d.associationLabel.Resize(assocSize)

	manaPos := nameTopLeft.Add(fyne.NewPos(8, nameSize.Height-6))
	manaSize := d.manaBox.MinSize()
	d.manaBox.Move(manaPos)
	d.manaBox.Resize(fyne.NewSize(float32(20*len(d.li.manaCost)), manaSize.Height))

	countSize := d.countLbl.MinSize()
	topLeft = manaPos.Add(fyne.NewPos(-8, manaSize.Height-6))
	d.countLbl.Move(topLeft)
	d.countLbl.Resize(countSize)

	typeSize := d.typeLine.MinSize()
	topLeft = topLeft.Add(fyne.NewPos(0, countSize.Height-6))
	d.typeLine.Move(topLeft)
	d.typeLine.Resize(fyne.NewSize(size.Width-cardSize.Width-2*theme.Padding(), typeSize.Height))

	topLeft = topLeft.Add(fyne.NewPos(0, typeSize.Height-6))
	d.setName.Move(topLeft)
	//setSize := d.setName.MinSize()

	//topLeft = topLeft.Add(fyne.NewPos(0, setSize.Height+theme.Padding()))
	//d.separator.Move(topLeft)
	//sepSize := d.separator.MinSize()

	setCoverPos := fyne.NewPos(theme.Padding(), cardSize.Height)
	d.setCover.Move(setCoverPos)
	d.setCover.Resize(d.setCover.MinSize())
}

func (d deckCardListItemRenderer) MinSize() fyne.Size {
	nameSize := d.nameLbl.MinSize()
	assocSize := d.associationLabel.MinSize()
	typeSize := d.typeLine.MinSize()
	countSize := d.countLbl.Size()
	setSize := d.setName.MinSize()
	manaSize := d.manaBox.MinSize()
	cardSize := d.cardFace.Size()

	//sepSize := d.separator.MinSize()
	setCoverSize := d.setCover.MinSize()

	height := fyne.Max(cardSize.Height+setCoverSize.Height, nameSize.Height-6+typeSize.Height-6+countSize.Height-6+setSize.Height+manaSize.Height-6) + 2*theme.Padding()

	size := fyne.NewSize(cardSize.Width+fyne.Max(fyne.Max(nameSize.Width+assocSize.Width, 200), typeSize.Width)+3*theme.Padding(), height)
	//log.Println("Min Size:", size.Width, "X", size.Height)
	return size
}

func (d deckCardListItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.nameLbl, d.countLbl, d.cardFace, d.manaBox, d.typeLine, d.setName, d.associationLabel, d.separator, d.setCover}
}

func (d deckCardListItemRenderer) Refresh() {
	d.nameLbl.ParseMarkdown(fmt.Sprintf("#### %s", d.li.card.Card.Name))
	d.countLbl.SetText(fmt.Sprintf("X %d", d.li.card.Count))
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
