package card

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type cardMetaListItem struct {
	widget.BaseWidget

	cardMeta cardMeta
}

func newCardMetaListItem(cardMeta cardMeta) *cardMetaListItem {
	cmli := &cardMetaListItem{cardMeta: cardMeta}
	cmli.ExtendBaseWidget(cmli)

	return cmli
}

func (c *cardMetaListItem) UpdateData(meta cardMeta) {
	c.cardMeta = meta
	c.Refresh()
}

func (c *cardMetaListItem) CreateRenderer() fyne.WidgetRenderer {
	keyLbl := widget.NewRichTextFromMarkdown(fmt.Sprintf("### %s", c.cardMeta.key))
	valueLbl := widget.NewRichTextFromMarkdown(c.cardMeta.value)
	valueLbl.Wrapping = fyne.TextWrapWord

	return &cardMetaListItemRenderer{
		cmli:     c,
		keyLbl:   keyLbl,
		valueLbl: valueLbl,
	}
}

type cardMetaListItemRenderer struct {
	cmli     *cardMetaListItem
	keyLbl   *widget.RichText
	valueLbl *widget.RichText
}

func (c cardMetaListItemRenderer) Destroy() {

}

func (c cardMetaListItemRenderer) Layout(size fyne.Size) {
	topLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	c.keyLbl.Move(topLeft)
	keySize := c.keyLbl.MinSize()

	topLeft = topLeft.Add(fyne.NewPos(0, keySize.Height+theme.Padding()))
	c.valueLbl.Move(topLeft)
	c.valueLbl.Resize(fyne.NewSize(size.Width-theme.Padding(), c.valueLbl.MinSize().Height))
}

func (c cardMetaListItemRenderer) MinSize() fyne.Size {
	keySize := c.keyLbl.MinSize()
	valSize := c.valueLbl.MinSize()

	size := fyne.NewSize(valSize.Width+2*theme.Padding(), keySize.Height+valSize.Height+(3*theme.Padding()))
	return size
}

func (c cardMetaListItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{c.keyLbl, c.valueLbl}
}

func (c cardMetaListItemRenderer) Refresh() {
	c.keyLbl.ParseMarkdown(fmt.Sprintf("### %s", c.cmli.cardMeta.key))
	c.valueLbl.ParseMarkdown(c.cmli.cardMeta.value)
}
