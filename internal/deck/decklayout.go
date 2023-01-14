package deck

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/ansel"
	cards2 "gitlab.com/high-creek-software/goscryfall/cards"
	"gitlab.com/kendellfab/mtgstudio/internal/card"
	"gitlab.com/kendellfab/mtgstudio/internal/icons"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/storage"
	"log"
)

type DeckLayout struct {
	*fyne.Container

	manager    *storage.Manager
	symbolRepo symbol.SymbolRepo

	deckList *widget.List
	cardList *widget.List

	cardTab *container.DocTabs

	deckAdapter *DeckAdapter
	cardAdapter *card.CardAdapter
}

func NewDeckLayout(manager *storage.Manager, symbolRepo symbol.SymbolRepo, showImport func()) *DeckLayout {
	dl := &DeckLayout{manager: manager, symbolRepo: symbolRepo}
	anselLoader := ansel.NewAnsel[string](400, ansel.SetLoader[string](manager.LoadCardImage), ansel.SetLoadingImage[string](icons.CardLoadingResource), ansel.SetFailedImage[string](icons.CardFailedResource))
	dl.deckAdapter = NewDeckAdapter(nil, anselLoader)
	dl.cardAdapter = card.NewCardAdapter(anselLoader, dl.symbolRepo, nil)

	dl.deckList = widget.NewList(dl.deckAdapter.Count, dl.deckAdapter.CreateTemplate, dl.deckAdapter.UpdateTemplate)
	dl.deckAdapter.SetList(dl.deckList)
	dl.deckList.OnSelected = dl.deckSelected

	dl.cardList = widget.NewList(dl.cardAdapter.Count, dl.cardAdapter.CreateTemplate, dl.cardAdapter.UpdateTemplate)
	dl.cardList.OnSelected = dl.cardSelected

	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.ContentAddIcon(), showImport))
	dl.cardTab = container.NewDocTabs()

	insideSplit := container.NewHSplit(dl.cardList, dl.cardTab)
	insideSplit.SetOffset(0.25)
	hSplit := container.NewHSplit(dl.deckList, insideSplit)
	hSplit.SetOffset(0.18)
	dl.Container = container.NewBorder(toolbar, nil, nil, nil, hSplit)

	dl.LoadDecks()

	return dl
}

func (dl *DeckLayout) LoadDecks() {
	go func() {
		if decks, err := dl.manager.ListDecks(); err == nil {
			dl.deckAdapter.Update(decks)
			dl.deckList.Refresh()
		} else {
			log.Println("deck load error", err)
		}
	}()
}

func (dl *DeckLayout) deckSelected(id widget.ListItemID) {
	go func() {
		deck := dl.deckAdapter.Item(id)
		dl.cardAdapter.Clear()
		dl.cardList.Refresh()
		if fullDeck, err := dl.manager.LoadDeck(deck.ID); err == nil {
			var cs []cards2.Card
			if fullDeck.Commander != nil {
				cs = append(cs, *fullDeck.Commander)
			}
			cs = append(cs, fullDeck.Main...)
			cs = append(cs, fullDeck.Sideboard...)
			dl.cardAdapter.AppendCards(cs)
		} else {
			log.Println("Error loading deck", err)
		}

	}()
}

func (dl *DeckLayout) cardSelected(id widget.ListItemID) {
	c := dl.cardAdapter.Item(id)

	cardLayout := card.NewCardLayout(&c, dl.symbolRepo, dl.manager, nil)
	tab := container.NewTabItem(c.Name, cardLayout.Container)
	dl.cardTab.Append(tab)
	dl.cardTab.Select(tab)
}
