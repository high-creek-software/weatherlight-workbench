package search

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/kendellfab/mtgstudio/internal/card"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/notifier"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/storage"
	"log"
)

type SearchLayout struct {
	*container.Split

	manager    *storage.Manager
	symbolRepo symbol.SymbolRepo
	notifier   notifier.Notifier

	cardTabs    *container.DocTabs
	cardList    *widget.List
	cardAdapter *card.CardAdapter

	name      *widget.Entry
	typeLine  *widget.Entry
	searchBtn *widget.Button

	whiteCheck *widget.Check
	blueCheck  *widget.Check
	blackCheck *widget.Check
	redCheck   *widget.Check
	greenCheck *widget.Check

	brawlCheck *widget.Check
}

func NewSearchLayout(manager *storage.Manager, symbolRepo symbol.SymbolRepo, n notifier.Notifier) *SearchLayout {
	sl := &SearchLayout{manager: manager, symbolRepo: symbolRepo, notifier: n}

	sl.cardAdapter = card.NewCardAdapter(
		ansel.NewAnsel[string](400, ansel.SetLoader[string](sl.manager.LoadCardImage), ansel.SetWorkerCount[string](10), ansel.SetLoadingImage[string](storage.CardLoadingResource), ansel.SetFailedImage[string](storage.CardFailedResource)),
		sl.symbolRepo,
		nil,
	)

	sl.cardTabs = container.NewDocTabs()
	sl.cardList = widget.NewList(sl.cardAdapter.Count, sl.cardAdapter.CreateTemplate, sl.cardAdapter.UpdateTemplate)
	sl.cardList.OnSelected = sl.cardSelected

	sl.name = widget.NewEntry()
	sl.name.SetPlaceHolder("Card Name")
	sl.typeLine = widget.NewEntry()
	sl.typeLine.SetPlaceHolder("Type Line")
	sl.searchBtn = widget.NewButton("Search", sl.doSearch)
	sl.whiteCheck = widget.NewCheck("White", nil)
	sl.blueCheck = widget.NewCheck("Blue", nil)
	sl.blackCheck = widget.NewCheck("Black", nil)
	sl.redCheck = widget.NewCheck("Red", nil)
	sl.greenCheck = widget.NewCheck("Green", nil)

	sl.brawlCheck = widget.NewCheck("Brawl Legal", nil)

	insideSplit := container.NewHSplit(sl.cardList, sl.cardTabs)
	insideSplit.SetOffset(0.20)
	sl.Split = container.NewHSplit(
		container.NewPadded(container.NewVBox(sl.name, sl.typeLine, widget.NewSeparator(), sl.whiteCheck, sl.blueCheck, sl.blackCheck, sl.redCheck, sl.greenCheck, widget.NewSeparator(), sl.brawlCheck, sl.searchBtn)),
		insideSplit,
	)
	sl.Split.SetOffset(0.15)

	return sl
}

func (sl *SearchLayout) doSearch() {
	name := sl.name.Text
	typeLine := sl.typeLine.Text

	sr := storage.SearchRequest{Name: name, TypeLine: typeLine}
	sr.White = sl.whiteCheck.Checked
	sr.Blue = sl.blueCheck.Checked
	sr.Black = sl.blackCheck.Checked
	sr.Red = sl.redCheck.Checked
	sr.Green = sl.greenCheck.Checked
	sr.BrawlLegal = sl.brawlCheck.Checked

	go func() {
		cards, err := sl.manager.Search(sr)
		if err != nil {
			sl.notifier.ShowError(err)
			return
		}

		log.Println("Found", len(cards), "cards.")

		if len(cards) > 0 {
			sl.cardAdapter.Clear()
			sl.cardList.Refresh()
			sl.cardAdapter.AppendCards(cards)
			sl.cardList.Refresh()
			sl.cardList.ScrollToTop()
		} else {
			sl.notifier.ShowDialog("None Found", "No cards were found that match that query.")
		}
	}()
}

func (sl *SearchLayout) cardSelected(id widget.ListItemID) {
	c := sl.cardAdapter.Item(id)

	cardLayout := card.NewCardLayout(&c, sl.symbolRepo, sl.manager, sl.notifier)
	tab := container.NewTabItem(c.Name, cardLayout.Scroll)
	sl.cardTabs.Append(tab)
	sl.cardTabs.Select(tab)
}
