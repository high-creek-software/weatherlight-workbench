package deck

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
	"log"
)

type DeckLayout struct {
	*fyne.Container
	canvas fyne.Canvas

	deckList *widget.List
	cardList *widget.List

	deckTabs *container.DocTabs

	deckAdapter  *DeckAdapter
	registry     *platform.Registry
	selectedDeck storage.Deck
}

func NewDeckLayout(canvas fyne.Canvas, registry *platform.Registry, showImport func()) *DeckLayout {
	dl := &DeckLayout{canvas: canvas, registry: registry}
	dl.deckAdapter = NewDeckAdapter(nil, dl.registry, dl)

	dl.deckList = widget.NewList(dl.deckAdapter.Count, dl.deckAdapter.CreateTemplate, dl.deckAdapter.UpdateTemplate)
	dl.deckAdapter.SetList(dl.deckList)
	dl.deckList.OnSelected = dl.deckSelected

	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.ContentAddIcon(), showImport))
	dl.deckTabs = container.NewDocTabs()

	hSplit := container.NewHSplit(dl.deckList, dl.deckTabs)
	hSplit.SetOffset(0.15)
	dl.Container = container.NewBorder(toolbar, nil, nil, nil, hSplit)

	dl.LoadDecks()

	return dl
}

func (dl *DeckLayout) Remove(d storage.Deck) {
	dl.registry.Notifier.VerifyAction(fmt.Sprintf("Are you sure you want to remove deck: %s", d.Name), "Remove", func() {
		err := dl.registry.Manager.RemoveDeck(d)
		if err != nil {
			dl.registry.Notifier.ShowError(err)
		} else {
			dl.LoadDecks()
		}
	})
}

func (dl *DeckLayout) Copy(d storage.Deck) {
	dl.registry.Notifier.ShowDialog("", fmt.Sprintf("Copy %s; not done yet", d.Name))
}

func (dl *DeckLayout) LoadDecks() {
	go func() {
		if decks, err := dl.registry.Manager.ListDecks(); err == nil {
			dl.deckAdapter.Update(decks)
			dl.deckList.Refresh()
		} else {
			log.Println("deck load error", err)
		}
	}()
}

func (dl *DeckLayout) deckSelected(id widget.ListItemID) {
	deck := dl.deckAdapter.Item(id)
	dl.selectedDeck = deck
	deckDisplay := NewDeckMetaDisplay(dl.canvas, dl.registry, deck, dl.LoadDecks)
	tab := container.NewTabItem(deck.Name, deckDisplay.Container)
	dl.deckTabs.Append(tab)
	dl.deckTabs.Select(tab)

}
