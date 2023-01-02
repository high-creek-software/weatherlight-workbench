package browse

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/kendellfab/mtgstudio/internal/card"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/notifier"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/set"
	"gitlab.com/kendellfab/mtgstudio/internal/storage"
	"log"
)

type BrowseLayout struct {
	*container.Split

	setList    *widget.List
	setAdapter *set.SetAdapter

	cardTabs    *container.DocTabs
	cardList    *widget.List
	cardAdapter *card.CardAdapter

	manager    *storage.Manager
	notifier   notifier.Notifier
	symbolRepo symbol.SymbolRepo
}

func NewBrowseLayout(manager *storage.Manager, symbolRepo symbol.SymbolRepo, n notifier.Notifier, updateSetIcon ansel.LoaderCallback, resizeCardArt ansel.LoaderCallback) *BrowseLayout {
	bl := &BrowseLayout{manager: manager, notifier: n, symbolRepo: symbolRepo}

	bl.setAdapter = set.NewSetAdapter(
		ansel.NewAnsel[string](100, ansel.SetLoadedCallback[string](updateSetIcon), ansel.SetLoader[string](bl.manager.LoadSetIcon)),
	)
	bl.cardAdapter = card.NewCardAdapter(
		ansel.NewAnsel[string](400, ansel.SetLoader[string](bl.manager.LoadCardImage), ansel.SetLoadedCallback[string](resizeCardArt), ansel.SetWorkerCount[string](40)),
		ansel.NewAnsel[string](200, ansel.SetLoader[string](bl.manager.LoadSymbolImage)),
		bl.symbolRepo,
	)

	bl.setList = widget.NewList(bl.setAdapter.Count, bl.setAdapter.CreateTemplate, bl.setAdapter.UpdateTemplate)
	bl.setList.OnSelected = bl.setSelected

	bl.cardTabs = container.NewDocTabs()
	bl.cardList = widget.NewList(bl.cardAdapter.Count, bl.cardAdapter.CreateTemplate, bl.cardAdapter.UpdateTemplate)
	bl.cardList.OnSelected = bl.cardSelected

	insideSplit := container.NewHSplit(bl.cardList, bl.cardTabs)
	insideSplit.SetOffset(0.20)
	bl.Split = container.NewHSplit(bl.setList, insideSplit)
	bl.Split.SetOffset(0.20)

	return bl
}

func (bl *BrowseLayout) setSelected(id widget.ListItemID) {
	set := bl.setAdapter.Item(id)
	log.Println("Selected:", set.Id)
	bl.cardAdapter.Clear()
	bl.cardList.Refresh()
	go func() {
		allCards, err := bl.manager.ListBySet(set.Code)
		if err != nil {
			bl.notifier.ShowError(err)
			return
		}
		//allCards = append(allCards, response.Data...)
		//
		//for response.HasMore {
		//	if response, err = bl.client.ListCards(set.Code, response.NextPage); err == nil {
		//		allCards = append(allCards, response.Data...)
		//	}
		//}

		bl.cardAdapter.AppendCards(allCards)
		bl.cardList.Refresh()
	}()
}

func (bl *BrowseLayout) cardSelected(id widget.ListItemID) {
	c := bl.cardAdapter.Item(id)

	cardLayout := card.NewCardLayout(&c, bl.symbolRepo, bl.manager)
	tab := container.NewTabItem(c.Name, cardLayout.Scroll)
	bl.cardTabs.Append(tab)
	bl.cardTabs.Select(tab)
}

func (bl *BrowseLayout) LoadSets() {
	sets, err := bl.manager.ListSets()
	if err != nil {
		bl.notifier.ShowError(err)
		return
	}

	bl.setAdapter.AddSets(sets)
	bl.setList.Refresh()
}
