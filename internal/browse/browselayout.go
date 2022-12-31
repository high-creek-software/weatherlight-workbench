package browse

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/high-creek-software/goscryfall"
	"gitlab.com/kendellfab/mtgstudio/internal/card"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/notifier"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/resources"
	"gitlab.com/kendellfab/mtgstudio/internal/set"
	"log"
)

type BrowseLayout struct {
	*container.Split

	setList    *widget.List
	setAdapter *set.SetAdapter

	cardTabs    *container.DocTabs
	cardList    *widget.List
	cardAdapter *card.CardAdapter

	manager    *resources.Manager
	client     *goscryfall.Client
	notifier   notifier.Notifier
	symbolRepo symbol.SymbolRepo
}

func NewBrowseLayout(manager *resources.Manager, client *goscryfall.Client, symbolRepo symbol.SymbolRepo, n notifier.Notifier, updateSetIcon ansel.LoaderCallback, resizeCardArt ansel.LoaderCallback) *BrowseLayout {
	bl := &BrowseLayout{manager: manager, client: client, notifier: n, symbolRepo: symbolRepo}

	bl.setAdapter = set.NewSetAdapter(
		ansel.NewAnsel[string](100, ansel.SetLoadedCallback[string](updateSetIcon), ansel.SetLoader[string](bl.manager.LoadSetIcon)),
	)
	bl.cardAdapter = card.NewCardAdapter(
		ansel.NewAnsel[string](400, ansel.SetLoader[string](bl.manager.LoadCardImage), ansel.SetLoadedCallback[string](resizeCardArt)),
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
		cards, err := bl.client.ListCards(set.Code, "")
		if err != nil {
			bl.notifier.ShowError(err)
			return
		}
		bl.cardAdapter.AppendCards(cards.Data)
		bl.cardList.Refresh()
	}()
}

func (bl *BrowseLayout) cardSelected(id widget.ListItemID) {
	card := bl.cardAdapter.Item(id)
	go func() {
		if img, err := bl.manager.LoadCardImage(card.ImageUris.Png); err == nil {
			image := canvas.NewImageFromResource(fyne.NewStaticResource(card.ImageUris.Png, img))
			image.FillMode = canvas.ImageFillContain
			tab := container.NewTabItem(card.Name, image)
			bl.cardTabs.Append(tab)
			bl.cardTabs.Select(tab)
		}
	}()
}

func (bl *BrowseLayout) LoadSets() {
	sets, err := bl.client.ListSets()
	if err != nil {
		bl.notifier.ShowError(err)
		return
	}

	bl.setAdapter.AddSets(sets)
}
