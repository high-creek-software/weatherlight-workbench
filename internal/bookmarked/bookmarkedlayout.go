package bookmarked

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/kendellfab/mtgstudio/internal/card"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/notifier"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/storage"
)

type BookmarkedLayout struct {
	*container.Split

	cardList    *widget.List
	cardAdapter *card.CardAdapter
	cardTabs    *container.DocTabs

	manager    *storage.Manager
	symbolRepo symbol.SymbolRepo
	notifier   notifier.Notifier
}

func NewBookmarkedLayout(manager *storage.Manager, symbolRepo symbol.SymbolRepo, n notifier.Notifier) *BookmarkedLayout {
	bl := &BookmarkedLayout{manager: manager, symbolRepo: symbolRepo, notifier: n}

	bl.cardAdapter = card.NewCardAdapter(
		ansel.NewAnsel[string](400, ansel.SetLoader[string](bl.manager.LoadCardImage), ansel.SetWorkerCount[string](10)),
		bl.symbolRepo,
		nil,
	)
	bl.cardTabs = container.NewDocTabs()
	bl.cardList = widget.NewList(bl.cardAdapter.Count, bl.cardAdapter.CreateTemplate, bl.cardAdapter.UpdateTemplate)
	bl.cardList.OnSelected = bl.cardSelected

	bl.Split = container.NewHSplit(bl.cardList, bl.cardTabs)
	bl.Split.SetOffset(0.18)

	return bl
}

func (bl *BookmarkedLayout) LoadBookmarked() {
	go func() {
		cards, err := bl.manager.ListBookmarked()
		if err != nil {
			bl.notifier.ShowError(err)
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

	cardLayout := card.NewCardLayout(&c, bl.symbolRepo, bl.manager, bl.notifier)
	tab := container.NewTabItem(c.Name, cardLayout.Scroll)
	bl.cardTabs.Append(tab)
	bl.cardTabs.Select(tab)
}
