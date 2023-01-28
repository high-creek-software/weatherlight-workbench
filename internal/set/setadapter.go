package set

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/ansel"
	"github.com/high-creek-software/goscryfall/sets"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/adapter"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sync"
)

var _ adapter.Adapter[sets.Set] = (*SetAdapter)(nil)

type SetAdapter struct {
	sets     []sets.Set
	filtered []sets.Set
	loader   *ansel.Ansel[string]

	locker sync.Locker
}

func NewSetAdapter(loader *ansel.Ansel[string]) *SetAdapter {
	return &SetAdapter{loader: loader, locker: &sync.Mutex{}}
}

func (sa *SetAdapter) AddSets(s []sets.Set) {
	sa.sets = s
}

func (sa *SetAdapter) Count() int {
	if sa.sets == nil {
		return 0
	}

	if len(sa.filtered) > 0 {
		return len(sa.filtered)
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
	if len(sa.filtered) > 0 {
		return sa.filtered[id]
	}
	return sa.sets[id]
}

func (sa *SetAdapter) ExecuteFilter(filter string) {
	sa.locker.Lock()
	defer sa.locker.Unlock()
	sa.filtered = nil
	if filter == "" || len(filter) <= 3 {
		return
	}

	for _, set := range sa.sets {
		if fuzzy.MatchNormalizedFold(filter, set.Name) {
			sa.filtered = append(sa.filtered, set)
		}
	}
}
