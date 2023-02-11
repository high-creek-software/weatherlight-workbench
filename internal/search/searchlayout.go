package search

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/tabman"
	"github.com/high-creek-software/weatherlight-workbench/internal/card"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform"
	storage2 "github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
	"log"
)

type SearchLayout struct {
	*container.Split
	canvas fyne.Canvas

	registry *platform.Registry

	cardTabs       *container.DocTabs
	cardList       *widget.List
	cardAdapter    *card.CardAdapter
	cardTabManager *tabman.Manager[string]

	name        *widget.Entry
	typeLine    *widget.Entry
	oracleEntry *widget.Entry
	searchBtn   *widget.Button

	whiteCheck *widget.Check
	blueCheck  *widget.Check
	blackCheck *widget.Check
	redCheck   *widget.Check
	greenCheck *widget.Check

	standardCheck        *widget.Check
	futureCheck          *widget.Check
	historicCheck        *widget.Check
	gladiatorCheck       *widget.Check
	pioneerCheck         *widget.Check
	explorerCheck        *widget.Check
	modernCheck          *widget.Check
	legacyCheck          *widget.Check
	pauperCheck          *widget.Check
	vintageCheck         *widget.Check
	pennyCheck           *widget.Check
	commanderCheck       *widget.Check
	brawlCheck           *widget.Check
	historicBrawlCheck   *widget.Check
	alchemyCheck         *widget.Check
	pauperCommanderCheck *widget.Check
	duelCheck            *widget.Check
	oldschoolCheck       *widget.Check
	premodernCheck       *widget.Check
}

func NewSearchLayout(cvs fyne.Canvas, registry *platform.Registry) *SearchLayout {
	sl := &SearchLayout{canvas: cvs, registry: registry}

	sl.cardAdapter = card.NewCardAdapter(sl.registry)

	sl.cardTabs = container.NewDocTabs()
	sl.cardList = widget.NewList(sl.cardAdapter.Count, sl.cardAdapter.CreateTemplate, sl.cardAdapter.UpdateTemplate)
	sl.cardList.OnSelected = sl.cardSelected
	sl.cardAdapter.SetList(sl.cardList)
	sl.cardTabManager = tabman.NewManager[string]()
	sl.cardTabs.OnClosed = sl.cardTabManager.RemoveTab

	sl.name = widget.NewEntry()
	sl.name.SetPlaceHolder("Card Name")
	sl.typeLine = widget.NewEntry()
	sl.typeLine.SetPlaceHolder("Type Line")
	sl.oracleEntry = widget.NewEntry()
	sl.oracleEntry.SetPlaceHolder("Oracle Text")
	sl.searchBtn = widget.NewButton("Search", sl.doSearch)

	colorsLbl := widget.NewRichTextFromMarkdown("## Colors")
	sl.whiteCheck = widget.NewCheck("White", nil)
	sl.blueCheck = widget.NewCheck("Blue", nil)
	sl.blackCheck = widget.NewCheck("Black", nil)
	sl.redCheck = widget.NewCheck("Red", nil)
	sl.greenCheck = widget.NewCheck("Green", nil)
	colorWrapper := container.NewAdaptiveGrid(2, sl.whiteCheck, sl.blueCheck, sl.blackCheck, sl.redCheck, sl.greenCheck)

	legalLbl := widget.NewRichTextFromMarkdown("## Legalities")
	sl.standardCheck = widget.NewCheck("Standard", nil)
	sl.futureCheck = widget.NewCheck("Future", nil)
	sl.historicCheck = widget.NewCheck("Historic", nil)
	sl.gladiatorCheck = widget.NewCheck("Gladiator", nil)
	sl.pioneerCheck = widget.NewCheck("Pioneer", nil)
	sl.explorerCheck = widget.NewCheck("Explorer", nil)
	sl.modernCheck = widget.NewCheck("Modern", nil)
	sl.legacyCheck = widget.NewCheck("Legacy", nil)
	sl.pauperCheck = widget.NewCheck("Pauper", nil)
	sl.vintageCheck = widget.NewCheck("Vintage", nil)
	sl.pennyCheck = widget.NewCheck("Penny", nil)
	sl.commanderCheck = widget.NewCheck("Commander", nil)
	sl.brawlCheck = widget.NewCheck("Brawl", nil)
	sl.historicBrawlCheck = widget.NewCheck("Historic Brawl", nil)
	sl.alchemyCheck = widget.NewCheck("Alchemy", nil)
	sl.pauperCommanderCheck = widget.NewCheck("Pauper Commander", nil)
	sl.duelCheck = widget.NewCheck("Duel", nil)
	sl.oldschoolCheck = widget.NewCheck("Oldschool", nil)
	sl.premodernCheck = widget.NewCheck("Premodern", nil)

	legalWrapper := container.NewAdaptiveGrid(2, sl.standardCheck,
		sl.futureCheck, sl.historicCheck,
		sl.gladiatorCheck, sl.pioneerCheck,
		sl.explorerCheck, sl.modernCheck,
		sl.legacyCheck, sl.pauperCheck,
		sl.vintageCheck, sl.pennyCheck,
		sl.commanderCheck, sl.brawlCheck,
		sl.historicBrawlCheck, sl.alchemyCheck,
		sl.pauperCommanderCheck, sl.duelCheck,
		sl.oldschoolCheck, sl.premodernCheck)

	insideSplit := container.NewHSplit(sl.cardList, sl.cardTabs)
	insideSplit.SetOffset(0.20)
	scroll := container.NewScroll(container.NewPadded(container.NewVBox(sl.name, sl.typeLine, sl.oracleEntry, colorsLbl, widget.NewSeparator(), colorWrapper, legalLbl, widget.NewSeparator(), legalWrapper, sl.searchBtn)))
	scroll.Direction = container.ScrollVerticalOnly
	sl.Split = container.NewHSplit(
		scroll,
		insideSplit,
	)
	sl.Split.SetOffset(0.15)

	return sl
}

func (sl *SearchLayout) doSearch() {
	name := sl.name.Text
	typeLine := sl.typeLine.Text

	sr := storage2.SearchRequest{Name: name, TypeLine: typeLine, OracleText: sl.oracleEntry.Text}
	sr.White = sl.whiteCheck.Checked
	sr.Blue = sl.blueCheck.Checked
	sr.Black = sl.blackCheck.Checked
	sr.Red = sl.redCheck.Checked
	sr.Green = sl.greenCheck.Checked
	sr.StandardLegal = sl.standardCheck.Checked
	sr.FutureLegal = sl.futureCheck.Checked
	sr.HistoricLegal = sl.historicCheck.Checked
	sr.GladiatorLegal = sl.gladiatorCheck.Checked
	sr.PioneerLegal = sl.pioneerCheck.Checked
	sr.ExplorerLegal = sl.explorerCheck.Checked
	sr.ModernLegal = sl.modernCheck.Checked
	sr.LegacyLegal = sl.legacyCheck.Checked
	sr.PauperLegal = sl.pauperCheck.Checked
	sr.VintageLegal = sl.vintageCheck.Checked
	sr.PennyLegal = sl.pennyCheck.Checked
	sr.CommanderLegal = sl.commanderCheck.Checked
	sr.BrawlLegal = sl.brawlCheck.Checked
	sr.HistoricBrawlLegal = sl.historicBrawlCheck.Checked
	sr.AlchemyLegal = sl.alchemyCheck.Checked
	sr.PauperCommanderLegal = sl.pauperCommanderCheck.Checked
	sr.DuelLegal = sl.duelCheck.Checked
	sr.OldschoolLegal = sl.oldschoolCheck.Checked
	sr.PremodernLegal = sl.premodernCheck.Checked

	if sr.IsEmpty() {
		sl.registry.Notifier.ShowDialog("", "Please select a search term.")
		return
	}

	go func() {
		cards, err := sl.registry.Manager.Search(sr)
		if err != nil {
			sl.registry.Notifier.ShowError(err)
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
			sl.registry.Notifier.ShowDialog("None Found", "No cards were found that match that query.")
		}
	}()
}

func (sl *SearchLayout) cardSelected(id widget.ListItemID) {
	c := sl.cardAdapter.Item(id)
	sl.cardList.UnselectAll()
	if ti, ok := sl.cardTabManager.GetTabItem(c.Id); ok {
		sl.cardTabs.Select(ti)
		return
	}

	cardLayout := card.NewCardLayout(sl.canvas, &c, sl.registry)
	tab := container.NewTabItem(c.Name, cardLayout.Container)
	sl.cardTabs.Append(tab)
	sl.cardTabs.Select(tab)
	sl.cardTabManager.AddTabItem(c.Id, tab)
}
