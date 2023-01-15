package browse

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/kendellfab/mtgstudio/internal/card"
	"gitlab.com/kendellfab/mtgstudio/internal/platform"
	"gitlab.com/kendellfab/mtgstudio/internal/set"
	"log"
)

type BrowseLayout struct {
	*container.Split

	filterEntry *widget.Entry
	filterClear *widget.Button

	setList    *widget.List
	setAdapter *set.SetAdapter

	cardTabs    *container.DocTabs
	cardList    *widget.List
	cardAdapter *card.CardAdapter

	registry *platform.Registry
}

func NewBrowseLayout(registry *platform.Registry, updateSetIcon ansel.LoaderCallback, resizeCardArt ansel.LoaderCallback) *BrowseLayout {
	bl := &BrowseLayout{registry: registry}

	bl.setAdapter = set.NewSetAdapter(bl.registry.SetIconLoader)
	bl.cardAdapter = card.NewCardAdapter(bl.registry)

	bl.setList = widget.NewList(bl.setAdapter.Count, bl.setAdapter.CreateTemplate, bl.setAdapter.UpdateTemplate)
	bl.setList.OnSelected = bl.setSelected

	bl.cardTabs = container.NewDocTabs()
	bl.cardList = widget.NewList(bl.cardAdapter.Count, bl.cardAdapter.CreateTemplate, bl.cardAdapter.UpdateTemplate)
	bl.cardList.OnSelected = bl.cardSelected

	bl.filterEntry = widget.NewEntry()
	bl.filterEntry.PlaceHolder = "Filter set name..."
	bl.filterEntry.OnChanged = bl.filterChanged
	bl.filterClear = widget.NewButtonWithIcon("", theme.ContentClearIcon(), bl.clearFilter)
	bl.filterClear.Importance = widget.MediumImportance
	filterBorder := container.NewBorder(nil, nil, nil, bl.filterClear, bl.filterEntry)

	insideSplit := container.NewHSplit(bl.cardList, bl.cardTabs)
	insideSplit.SetOffset(0.20)
	bl.Split = container.NewHSplit(container.NewBorder(container.NewPadded(filterBorder), nil, nil, nil, bl.setList), insideSplit)
	bl.Split.SetOffset(0.15)

	return bl
}

func (bl *BrowseLayout) clearFilter() {
	bl.filterEntry.SetText("")
	bl.setAdapter.ExecuteFilter("")
	bl.setList.Refresh()
}

func (bl *BrowseLayout) filterChanged(input string) {
	bl.setAdapter.ExecuteFilter(input)
	bl.setList.Refresh()
}

func (bl *BrowseLayout) setSelected(id widget.ListItemID) {
	set := bl.setAdapter.Item(id)
	log.Println("Selected:", set.Id)
	bl.cardAdapter.Clear()
	bl.cardList.Refresh()
	go func() {
		allCards, err := bl.registry.Manager.ListBySet(set.Code)
		if err != nil {
			bl.registry.Notifier.ShowError(err)
			return
		}

		bl.cardAdapter.AppendCards(allCards)
		bl.cardList.Refresh()
	}()
}

func (bl *BrowseLayout) cardSelected(id widget.ListItemID) {
	c := bl.cardAdapter.Item(id)

	cardLayout := card.NewCardLayout(&c, bl.registry)
	tab := container.NewTabItem(c.Name, cardLayout.Container)
	bl.cardTabs.Append(tab)
	bl.cardTabs.Select(tab)
}

func (bl *BrowseLayout) LoadSets() {
	sets, err := bl.registry.Manager.ListSets()
	if err != nil {
		bl.registry.Notifier.ShowError(err)
		return
	}

	bl.setAdapter.AddSets(sets)
	bl.setList.Refresh()
}
