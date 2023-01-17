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
}

func (li *DeckCardListItem) CreateRenderer() fyne.WidgetRenderer {
	nameLbl := widget.NewRichTextFromMarkdown("Temp")
	countLbl := widget.NewLabel("Temp")
	cardFace := widget.NewIcon(nil)
	cardFace.Resize(fyne.NewSize(128, 128))
	manaBox := container.NewHBox()
	typeLine := widget.NewLabel("")
	setName := widget.NewLabel("")

	dr := &deckCardListItemRenderer{
		li:       li,
		nameLbl:  nameLbl,
		countLbl: countLbl,
		cardFace: cardFace,
		manaBox:  manaBox,
		typeLine: typeLine,
		setName:  setName,
	}

	for i := 0; i < 4; i++ {
		dr.manaImages = append(dr.manaImages, widget.NewIcon(nil))
	}

	return dr
}

func NewDeckCardListItem() *DeckCardListItem {
	li := &DeckCardListItem{}
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
	li         *DeckCardListItem
	nameLbl    *widget.RichText
	countLbl   *widget.Label
	cardFace   *widget.Icon
	manaBox    *fyne.Container
	manaImages []*widget.Icon
	typeLine   *widget.Label
	setName    *widget.Label
}

func (d deckCardListItemRenderer) Destroy() {

}

func (d deckCardListItemRenderer) Layout(size fyne.Size) {
	topLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	cardSize := d.cardFace.Size()
	d.cardFace.Move(topLeft)

	nameSize := d.nameLbl.MinSize()
	nameTopLeft := topLeft.Add(fyne.NewPos(cardSize.Width+theme.Padding(), 10))
	d.nameLbl.Move(nameTopLeft)
	d.nameLbl.Resize(fyne.NewSize(size.Width, nameSize.Height))

	manaPos := nameTopLeft.Add(fyne.NewPos(8, 22))
	manaSize := d.manaBox.MinSize()
	d.manaBox.Move(manaPos)
	d.manaBox.Resize(fyne.NewSize(float32(20*len(d.li.manaCost)), 32))

	countSize := d.countLbl.MinSize()
	topLeft = manaPos.Add(fyne.NewPos(-8, manaSize.Height))
	d.countLbl.Move(topLeft)
	d.countLbl.Resize(countSize)

	//typeSize := d.typeLine.MinSize()
	topLeft = topLeft.Add(fyne.NewPos(0, 22))
	d.typeLine.Move(topLeft)

	topLeft = topLeft.Add(fyne.NewPos(0, 22))
	d.setName.Move(topLeft)
}

func (d deckCardListItemRenderer) MinSize() fyne.Size {
	nameSize := d.nameLbl.MinSize()
	typeSize := d.typeLine.MinSize()
	//countSize := d.countLbl.Size()
	//log.Println("--Name:", nameSize.Width, "X", nameSize.Height, "Count:", countSize.Width, "X", countSize.Height)
	cardSize := d.cardFace.Size()

	size := fyne.NewSize(theme.Padding()+cardSize.Width+theme.Padding()+fyne.Max(fyne.Max(nameSize.Width, 200), typeSize.Width)+theme.Padding(), cardSize.Height+2*theme.Padding())
	//log.Println("Min Size:", size.Width, "X", size.Height)
	return size
}

func (d deckCardListItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{d.nameLbl, d.countLbl, d.cardFace, d.manaBox, d.typeLine, d.setName}
}

func (d deckCardListItemRenderer) Refresh() {
	d.nameLbl.ParseMarkdown(fmt.Sprintf("#### %s", d.li.card.Card.Name))
	d.countLbl.SetText(fmt.Sprintf("X %d", d.li.card.Count))
	d.cardFace.SetResource(d.li.ico)
	d.typeLine.SetText(d.li.card.Card.TypeLine)
	d.setName.SetText(d.li.card.Card.SetName)

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
