package set

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/high-creek-software/goscryfall/sets"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/adapter"
)

var _ adapter.Adapter[sets.Set] = (*SetAdapter)(nil)

type SetAdapter struct {
	sets   []sets.Set
	loader *ansel.Ansel[string]
}

func NewSetAdapter(loader *ansel.Ansel[string]) *SetAdapter {
	return &SetAdapter{loader: loader}
}

func (sa *SetAdapter) AddSets(s []sets.Set) {
	sa.sets = s
}

func (sa *SetAdapter) Count() int {
	if sa.sets == nil {
		return 0
	}

	return len(sa.sets)
}

func (sa *SetAdapter) CreateTemplate() fyne.CanvasObject {
	return NewSetListItem(nil)
}

func (sa *SetAdapter) UpdateTemplate(id widget.ListItemID, co fyne.CanvasObject) {
	set := sa.Item(id)
	listItem := co.(*SetListItem)
	listItem.UpdateSet(&set)
	sa.loader.Load(set.Id, set.IconSvgUri, listItem)
}

func (sa *SetAdapter) Item(id widget.ListItemID) sets.Set {
	return sa.sets[id]
}
