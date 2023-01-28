package ruling

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/goscryfall/rulings"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/adapter"
)

var _ adapter.Adapter[rulings.Ruling] = (*RulingAdapter)(nil)

type RulingAdapter struct {
	r    []rulings.Ruling
	list *widget.List
}

func NewRulingAdapter(r []rulings.Ruling) *RulingAdapter {
	return &RulingAdapter{r: r}
}

func (ra *RulingAdapter) SetList(list *widget.List) {
	ra.list = list
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

	ra.list.SetItemHeight(id, li.MinSize().Height)
}

func (ra *RulingAdapter) Item(id widget.ListItemID) rulings.Ruling {
	return ra.r[id]
}
