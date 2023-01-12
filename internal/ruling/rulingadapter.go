package ruling

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/goscryfall/rulings"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/adapter"
)

var _ adapter.Adapter[rulings.Ruling] = (*RulingAdapter)(nil)

type RulingAdapter struct {
	r []rulings.Ruling
}

func NewRulingAdapter(r []rulings.Ruling) *RulingAdapter {
	return &RulingAdapter{r: r}
}

func (ra *RulingAdapter) Count() int {
	return len(ra.r)
}

func (ra *RulingAdapter) CreateTemplate() fyne.CanvasObject {
	return NewRulingListItem()
}

func (ra *RulingAdapter) UpdateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	r := ra.Item(id)
	li := co.(*RulingListItem)
	li.Set(r)
}

func (ra *RulingAdapter) Item(id widget.ListItemID) rulings.Ruling {
	return ra.r[id]
}
