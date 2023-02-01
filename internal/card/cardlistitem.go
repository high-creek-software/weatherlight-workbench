package card

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/goscryfall/cards"
	"log"
)

var _ fyne.Widget = (*CardListItem)(nil)

type CardListItem struct {
	widget.BaseWidget
	card    *cards.Card
	ico     fyne.Resource
	setIcon fyne.Resource

	manaCost []fyne.Resource
}

func NewCardListItem(card *cards.Card) *CardListItem {
	cli := &CardListItem{card: card}
	cli.ExtendBaseWidget(cli)

	return cli
}

func (cli *CardListItem) UpdateCard(card *cards.Card, manaCost []fyne.Resource) {
	cli.card = card
	cli.manaCost = manaCost
	cli.Refresh()
}

func (cli *CardListItem) SetResource(resource fyne.Resource) {
	cli.ico = resource
	cli.Refresh()
}

func (cli *CardListItem) CreateRenderer() fyne.WidgetRenderer {
	icon := widget.NewIcon(nil)
	icon.Resize(fyne.NewSize(128, 128))
	name := widget.NewRichTextWithText("template")
	name.Wrapping = fyne.TextWrapWord
	manaBox := container.NewHBox()
	typeLine := widget.NewLabel("template")
	typeLine.Wrapping = fyne.TextWrapWord
	setIcon := widget.NewIcon(nil)
	setName := widget.NewLabel("template")
	priceLbl := widget.NewLabel("template")

	renderer := &CardListItemRenderer{listItem: cli, icon: icon, name: name, manaBox: manaBox, typeLine: typeLine, setIcon: setIcon, setName: setName, priceLbl: priceLbl}

	for i := 0; i < 4; i++ {
		renderer.manaImages = append(renderer.manaImages, widget.NewIcon(nil))
	}

	return renderer
}

type CardListItemRenderer struct {
	listItem   *CardListItem
	icon       *widget.Icon
	name       *widget.RichText
	manaBox    *fyne.Container
	manaImages []*widget.Icon
	typeLine   *widget.Label
	setIcon    *widget.Icon
	setName    *widget.Label
	priceLbl   *widget.Label
}

func (c CardListItemRenderer) Destroy() {

}

func (c CardListItemRenderer) Layout(size fyne.Size) {
	iconPos := fyne.NewPos(theme.Padding(), theme.Padding())
	iconSize := c.icon.Size()
	c.icon.Move(iconPos)

	nameSize := c.name.MinSize()
	manaSize := c.manaBox.MinSize()

	namePos := iconPos.Add(fyne.NewPos(iconSize.Width+theme.Padding(), 0))
	c.name.Move(namePos)
	c.name.Resize(fyne.NewSize(size.Width-iconSize.Width-2*theme.Padding()-manaSize.Width-2*theme.Padding(), nameSize.Height))

	manaPos := fyne.NewPos(size.Width-theme.Padding()-manaSize.Width, theme.Padding()+6)
	c.manaBox.Move(manaPos)
	c.manaBox.Resize(fyne.NewSize(float32(20*len(c.listItem.manaCost)), manaSize.Height))

	typeSize := c.typeLine.MinSize()
	typeLinePos := namePos.Add(fyne.NewPos(0, nameSize.Height+2))
	c.typeLine.Move(typeLinePos)
	c.typeLine.Resize(fyne.NewSize(size.Width-iconSize.Width-2*theme.Padding(), typeSize.Height))

	setSize := c.setName.MinSize()
	setNamePos := typeLinePos.Add(fyne.NewPos(0, typeSize.Height))
	c.setName.Move(setNamePos)
	c.setName.Resize(fyne.NewSize(size.Width-iconSize.Width-2*theme.Padding(), setSize.Height))

	priceSize := c.priceLbl.MinSize()
	pricePos := setNamePos.Add(fyne.NewPos(0, priceSize.Height))
	c.priceLbl.Move(pricePos)
	c.priceLbl.Resize(priceSize)

}

func (c CardListItemRenderer) MinSize() fyne.Size {
	iconSize := c.icon.Size()
	nameSize := c.name.MinSize()
	manaSize := c.manaBox.MinSize()
	typeSize := c.typeLine.MinSize()
	setNameSize := c.setName.MinSize()
	priceSize := fyne.NewSize(0, 0)
	if c.priceLbl.Visible() {
		priceSize = c.priceLbl.MinSize()
	}

	height := fyne.Max(iconSize.Height, nameSize.Height-6+typeSize.Height-6+setNameSize.Height+manaSize.Height-6+priceSize.Height-6) + 2*theme.Padding()

	return fyne.NewSize(iconSize.Width+fyne.Max(fyne.Max(nameSize.Width+manaSize.Width, 200), typeSize.Width)+4*theme.Padding(), height)
}

func (c CardListItemRenderer) Objects() []fyne.CanvasObject {
	base := []fyne.CanvasObject{c.icon, c.name, c.manaBox, c.typeLine, c.setName, c.priceLbl}

	return base
}

func (c CardListItemRenderer) Refresh() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in card list item:", r)
		}
	}()
	c.icon.SetResource(c.listItem.ico)
	c.setIcon.SetResource(c.listItem.setIcon)
	c.setName.SetText(c.listItem.card.SetName)
	c.name.ParseMarkdown(fmt.Sprintf("### %s", c.listItem.card.Name))
	c.manaBox.RemoveAll()

	if c.listItem.card.Prices.Usd != "" {
		c.priceLbl.SetText(fmt.Sprintf("$%s", c.listItem.card.Prices.Usd))
		c.priceLbl.Visible()
	} else {
		c.priceLbl.Hide()
	}

	for i, cost := range c.listItem.manaCost {
		var ico *widget.Icon
		if i > len(c.manaImages)-1 {
			ico = widget.NewIcon(nil)
			c.manaImages = append(c.manaImages, ico)
		} else {
			ico = c.manaImages[i]
		}
		if ico == nil {
			ico = widget.NewIcon(nil)
			c.manaImages = append(c.manaImages, ico)
		}
		ico.SetResource(cost)
		if c.manaBox == nil || ico == nil {
			log.Println("Mana box is nil, don't know why !?!")
		} else {
			c.manaBox.Add(ico)
		}
	}

	c.typeLine.SetText(c.listItem.card.TypeLine)
}
