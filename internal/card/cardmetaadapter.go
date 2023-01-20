package card

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type cardMetaAdapter struct {
	meta []cardMeta
	list *widget.List
}

func newCardMetaAdapter(meta []cardMeta) *cardMetaAdapter {
	return &cardMetaAdapter{meta: meta}
}
func (c *cardMetaAdapter) updateMeta(meta []cardMeta) {
	c.meta = meta
}

func (c *cardMetaAdapter) SetList(list *widget.List) {
	c.list = list
}

func (c *cardMetaAdapter) Count() int {
	return len(c.meta)
}

func (c *cardMetaAdapter) CreateTemplate() fyne.CanvasObject {
	return newCardMetaListItem(cardMeta{})
}

func (c *cardMetaAdapter) UpdateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	m := c.Item(id)
	li := co.(*cardMetaListItem)
	li.UpdateData(m)

	c.list.SetItemHeight(id, li.MinSize().Height)
}

func (c *cardMetaAdapter) Item(id widget.ListItemID) cardMeta {
	return c.meta[id]
}

type cardMeta struct {
	key   string
	value string
}
