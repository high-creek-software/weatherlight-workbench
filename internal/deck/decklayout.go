package deck

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/fynecharts"
	"github.com/high-creek-software/weatherlight-workbench/internal/card"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
	"log"
	"strings"
)

type DeckLayout struct {
	*fyne.Container
	canvas fyne.Canvas

	deckList  *widget.List
	cardList  *widget.List
	manaChart *fynecharts.BarChart

	cardTab *container.DocTabs

	deckAdapter *DeckAdapter
	cardAdapter *DeckCardAdapter

	registry *platform.Registry

	selectedDeck storage.Deck
}

func NewDeckLayout(canvas fyne.Canvas, registry *platform.Registry, showImport func()) *DeckLayout {
	dl := &DeckLayout{canvas: canvas, registry: registry}
	dl.deckAdapter = NewDeckAdapter(nil, dl.registry)
	dl.cardAdapter = NewDeckCardAdapter(dl.registry, dl.setCover)

	dl.deckList = widget.NewList(dl.deckAdapter.Count, dl.deckAdapter.CreateTemplate, dl.deckAdapter.UpdateTemplate)
	dl.deckAdapter.SetList(dl.deckList)
	dl.deckList.OnSelected = dl.deckSelected

	dl.cardList = widget.NewList(dl.cardAdapter.Count, dl.cardAdapter.CreateTemplate, dl.cardAdapter.UpdateTemplate)
	dl.cardAdapter.SetList(dl.cardList)
	dl.cardList.OnSelected = dl.cardSelected
	dl.manaChart = fynecharts.NewBarChart(dl.canvas, "Mana Curve", nil, nil)
	dl.manaChart.SetMinHeight(150)
	dl.manaChart.UpdateHoverFormat(func(v float64) string {
		return fmt.Sprintf("%d", int(v))
	})

	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.ContentAddIcon(), showImport))
	dl.cardTab = container.NewDocTabs()

	insideSplit := container.NewHSplit(container.NewBorder(dl.manaChart, nil, nil, nil, dl.cardList), dl.cardTab)
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
		lbls := []string{"0", "1", "2", "3", "4", "5", "6", "7+"}
		data := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
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

			for idx := 0; idx < dl.cardAdapter.Count(); idx++ {
				crd := dl.cardAdapter.Item(idx)
				idx := int(crd.Card.Cmc)
				if idx > 7 {
					idx = 7
				}
				if idx == 0 && strings.Contains(crd.Card.TypeLine, "Land") {
					// For now continue, until I determine what is a land and what is not.
					continue
				}
				data[idx] += float64(crd.Count)
			}
			dl.manaChart.UpdateData(lbls, data)
			dl.manaChart.Refresh()
		} else {
			log.Println("Error loading deck", err)
		}

	}()
}

func (dl *DeckLayout) cardSelected(id widget.ListItemID) {
	c := dl.cardAdapter.Item(id)

	cardLayout := card.NewCardLayout(dl.canvas, &c.Card, dl.registry)
	tab := container.NewTabItem(c.Card.Name, cardLayout.Container)
	dl.cardTab.Append(tab)
	dl.cardTab.Select(tab)
}
