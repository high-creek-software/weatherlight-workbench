package deck

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/kendellfab/mtgstudio/internal/card"
	"gitlab.com/kendellfab/mtgstudio/internal/platform"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/storage"
	"log"
)

type DeckLayout struct {
	*fyne.Container

	deckList *widget.List
	cardList *widget.List

	cardTab *container.DocTabs

	deckAdapter *DeckAdapter
	cardAdapter *DeckCardAdapter

	registry *platform.Registry

	selectedDeck storage.Deck
}

func NewDeckLayout(registry *platform.Registry, showImport func()) *DeckLayout {
	dl := &DeckLayout{registry: registry}
	dl.deckAdapter = NewDeckAdapter(nil, dl.registry)
	dl.cardAdapter = NewDeckCardAdapter(dl.registry, dl.setCover)

	dl.deckList = widget.NewList(dl.deckAdapter.Count, dl.deckAdapter.CreateTemplate, dl.deckAdapter.UpdateTemplate)
	dl.deckAdapter.SetList(dl.deckList)
	dl.deckList.OnSelected = dl.deckSelected

	dl.cardList = widget.NewList(dl.cardAdapter.Count, dl.cardAdapter.CreateTemplate, dl.cardAdapter.UpdateTemplate)
	dl.cardAdapter.SetList(dl.cardList)
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
		if decks, err := dl.registry.Manager.ListDecks(); err == nil {
			dl.deckAdapter.Update(decks)
			dl.deckList.Refresh()
		} else {
			log.Println("deck load error", err)
		}
	}()
}

func (dl *DeckLayout) setCover(dc storage.DeckCard) {
	log.Println("Deck:", dl.selectedDeck.Name, dl.selectedDeck.ID, "Card", dc.Card.Name)
	if err := dl.registry.Manager.UpdateCover(dl.selectedDeck.ID, dc.Card.ImageUris.ArtCrop); err == nil {
		dl.LoadDecks()
	} else {
		dl.registry.Notifier.ShowError(err)
	}
}

func (dl *DeckLayout) deckSelected(id widget.ListItemID) {
	go func() {
		deck := dl.deckAdapter.Item(id)
		dl.selectedDeck = deck
		dl.cardAdapter.Clear()
		dl.cardList.Refresh()
		if fullDeck, err := dl.registry.Manager.LoadDeck(deck.ID); err == nil {
			dl.cardAdapter.Clear()
			if fullDeck.Commander != nil {
				dl.cardAdapter.AppendCards([]storage.DeckCard{*fullDeck.Commander})
			}
			if fullDeck.Main != nil {
				dl.cardAdapter.AppendCards(fullDeck.Main)
			}
			if fullDeck.Sideboard != nil {
				dl.cardAdapter.AppendCards(fullDeck.Sideboard)
			}
			dl.cardList.Refresh()
		} else {
			log.Println("Error loading deck", err)
		}

	}()
}

func (dl *DeckLayout) cardSelected(id widget.ListItemID) {
	c := dl.cardAdapter.Item(id)

	cardLayout := card.NewCardLayout(&c.Card, dl.registry)
	tab := container.NewTabItem(c.Card.Name, cardLayout.Container)
	dl.cardTab.Append(tab)
	dl.cardTab.Select(tab)
}
