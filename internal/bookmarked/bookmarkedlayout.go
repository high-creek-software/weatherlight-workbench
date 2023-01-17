package bookmarked

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/kendellfab/mtgstudio/internal/card"
	"gitlab.com/kendellfab/mtgstudio/internal/platform"
)

type BookmarkedLayout struct {
	*container.Split

	registry *platform.Registry

	cardList    *widget.List
	cardAdapter *card.CardAdapter
	cardTabs    *container.DocTabs
}

func NewBookmarkedLayout(registry *platform.Registry) *BookmarkedLayout {
	bl := &BookmarkedLayout{registry: registry}

	bl.cardAdapter = card.NewCardAdapter(bl.registry)
	bl.cardTabs = container.NewDocTabs()
	bl.cardList = widget.NewList(bl.cardAdapter.Count, bl.cardAdapter.CreateTemplate, bl.cardAdapter.UpdateTemplate)
	bl.cardList.OnSelected = bl.cardSelected
	bl.cardAdapter.SetList(bl.cardList)

	bl.Split = container.NewHSplit(bl.cardList, bl.cardTabs)
	bl.Split.SetOffset(0.18)

	return bl
}

func (bl *BookmarkedLayout) LoadBookmarked() {
	go func() {
		cards, err := bl.registry.Manager.ListBookmarked()
		if err != nil {
			bl.registry.Notifier.ShowError(err)
			return
		}

		bl.cardAdapter.Clear()
		bl.cardAdapter.AppendCards(cards)
		bl.cardList.Refresh()
		bl.cardList.ScrollToTop()
	}()
}

func (bl *BookmarkedLayout) cardSelected(id widget.ListItemID) {
	c := bl.cardAdapter.Item(id)

	cardLayout := card.NewCardLayout(&c, bl.registry)
	tab := container.NewTabItem(c.Name, cardLayout.Container)
	bl.cardTabs.Append(tab)
	bl.cardTabs.Select(tab)
}
