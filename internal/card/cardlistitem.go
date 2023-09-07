package card

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	//coverImage *widget.Icon
	coverImage *canvas.Image
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
	//cli.ico = resource
	//cli.Refresh()
	//slog.Info("card set resource")
	if resource != cli.coverImage.Resource {
		//slog.Info("	setting resource on card")
		//cli.coverImage.SetResource(resource)
		cli.coverImage.Resource = resource
		cli.coverImage.Refresh()
	}
}

func (cli *CardListItem) CreateRenderer() fyne.WidgetRenderer {
	//icon := widget.NewIcon(nil)
	//icon.Resize(fyne.NewSize(128, 128))
	//// This is breaking the rules of fyne, the widgets are not meant to be stored with the data. But there has been some performance and crashing issues, reloading the manabox takes too much time to empty and reload
	//// by setting the cover image on the cardlistitem an image load callback won't require refreshing the entire renderer
	//cli.coverImage = icon

	icon := canvas.NewImageFromResource(nil)
	icon.Resize(fyne.NewSize(100, 100))
	icon.FillMode = canvas.ImageFillContain
	icon.ScaleMode = canvas.ImageScalePixels
	cli.coverImage = icon

	name := widget.NewRichTextWithText("template")
	name.Wrapping = fyne.TextWrapWord
	typeLine := widget.NewRichTextWithText("template")
	typeLine.Wrapping = fyne.TextWrapWord
	setName := widget.NewLabel("template")
	setName.Wrapping = fyne.TextWrapWord

	renderer := &CardListItemRenderer{listItem: cli, icon: icon, name: name, manaBox: container.NewHBox(), typeLine: typeLine, setIcon: widget.NewIcon(nil), setName: setName, priceLbl: widget.NewLabel("template")}

	for i := 0; i < 4; i++ {
		renderer.manaImages = append(renderer.manaImages, widget.NewIcon(nil))
	}

	return renderer
}

type CardListItemRenderer struct {
	listItem   *CardListItem
	icon       *canvas.Image
	name       *widget.RichText
	manaBox    *fyne.Container
	manaImages []*widget.Icon
	typeLine   *widget.RichText
	setIcon    *widget.Icon
	setName    *widget.Label
	priceLbl   *widget.Label
}

func (c *CardListItemRenderer) Destroy() {

}

func (c *CardListItemRenderer) Layout(size fyne.Size) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Println("Recovered in card list item:", r)
	//	}
	//}()
	iconPos := fyne.NewPos(theme.Padding(), theme.Padding())
	iconSize := c.icon.Size()
	c.icon.Move(iconPos)

	nameSize := c.name.MinSize()
	manaSize := fyne.NewSize(0, 0)
	if c.manaBox != nil && len(c.manaBox.Objects) > 0 {
		manaSize = c.manaBox.MinSize()
	}

	namePos := iconPos.Add(fyne.NewPos(iconSize.Width+theme.Padding(), 0))
	c.name.Move(namePos)
	//c.name.Resize(fyne.NewSize(size.Width-iconSize.Width-2*theme.Padding(), nameSize.Height))

	//manaPos := fyne.NewPos(size.Width-theme.Padding()-manaSize.Width, theme.Padding()+6)
	manaPos := namePos.AddXY(0, nameSize.Height-5)
	c.manaBox.Move(manaPos.AddXY(8, 0))
	c.manaBox.Resize(fyne.NewSize(float32(20*len(c.listItem.manaCost)), manaSize.Height))

	typeSize := c.typeLine.MinSize()
	typeLinePos := manaPos.Add(fyne.NewPos(0, manaSize.Height-4))
	c.typeLine.Move(typeLinePos)
	//c.typeLine.Resize(fyne.NewSize(size.Width-iconSize.Width-2*theme.Padding(), typeSize.Height))

	setSize := c.setName.MinSize()
	setNamePos := typeLinePos.Add(fyne.NewPos(0, typeSize.Height-8))
	c.setName.Move(setNamePos)
	//c.setName.Resize(fyne.NewSize(size.Width-iconSize.Width-2*theme.Padding(), setSize.Height))

	priceSize := c.priceLbl.MinSize()
	pricePos := setNamePos.Add(fyne.NewPos(0, setSize.Height))
	c.priceLbl.Move(pricePos)
	c.priceLbl.Resize(priceSize)

}

func (c *CardListItemRenderer) MinSize() fyne.Size {
	//defer func() {
	//	if r := recover(); r != nil {
	//		log.Println("Recovered in card list item renderer minSize:", r)
	//	}
	//}()
	iconSize := c.icon.Size()
	nameSize := c.name.MinSize()
	manaSize := fyne.NewSize(0, 0)
	if c.manaBox != nil && len(c.manaBox.Objects) > 0 {
		manaSize = c.manaBox.MinSize()
	}
	typeSize := c.typeLine.MinSize()
	setNameSize := c.setName.MinSize()
	priceSize := fyne.NewSize(0, 0)
	if c.priceLbl.Visible() {
		priceSize = c.priceLbl.MinSize()
	}

	height := fyne.Max(iconSize.Height, nameSize.Height+typeSize.Height+setNameSize.Height+manaSize.Height+priceSize.Height)

	return fyne.NewSize(iconSize.Width+fyne.Max(fyne.Max(nameSize.Width+manaSize.Width, 200), typeSize.Width)+4*theme.Padding(), height)
}

func (c *CardListItemRenderer) Objects() []fyne.CanvasObject {
	base := []fyne.CanvasObject{c.icon, c.name, c.manaBox, c.typeLine, c.setName, c.priceLbl}

	return base
}

func (c *CardListItemRenderer) Refresh() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in card list item renderer refresh:", r)
		}
	}()
	//c.icon.SetResource(c.listItem.ico)
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

		if c.manaBox == nil || ico == nil || cost == nil {
			log.Println("Mana box is nil, don't know why !?!")
		} else {
			ico.SetResource(cost)
			c.manaBox.Add(ico)
		}
	}

	c.typeLine.ParseMarkdown(fmt.Sprintf("##### %s", c.listItem.card.TypeLine))
}
