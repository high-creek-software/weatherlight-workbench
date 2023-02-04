package deck

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/fynecharts"
	"github.com/high-creek-software/weatherlight-workbench/internal/card"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
	"log"
	"strings"
)

type DeckDisplay struct {
	*fyne.Container
	canvas fyne.Canvas

	registry     *platform.Registry
	selectedDeck storage.Deck

	cardList    *widget.List
	cardAdapter *DeckCardAdapter
	cardTabs    *container.DocTabs
	manaChart   *fynecharts.BarChart
	tpsChart    *fynecharts.BarChart
	deckImage   *widget.Icon
	deckNameLbl *widget.Label

	loadDecks func()
}

func (dd *DeckDisplay) SetCover(c storage.DeckCard) {
	log.Println("Deck:", dd.selectedDeck.Name, dd.selectedDeck.ID, "Card", c.Card.Name)
	if err := dd.registry.Manager.UpdateCover(dd.selectedDeck.ID, c.Card.ImageUris.ArtCrop); err == nil {
		if dd.loadDecks != nil {
			dd.loadDecks()
		}
	} else {
		dd.registry.Notifier.ShowError(err)
	}
}

func (dd *DeckDisplay) RemoveCard(c storage.DeckCard) {
	dd.registry.Notifier.ShowDialog("", fmt.Sprintf("Removing card: %s", c.Card.Name))
}

func (dd *DeckDisplay) IncCard(c storage.DeckCard) {
	dd.registry.Notifier.ShowDialog("", fmt.Sprintf("Increment card: %s", c.Card.Name))
}

func (dd *DeckDisplay) DecCard(c storage.DeckCard) {
	dd.registry.Notifier.ShowDialog("", fmt.Sprintf("Decrement card: %s", c.Card.Name))
}

func NewDeckMetaDisplay(canvas fyne.Canvas, registry *platform.Registry, deck storage.Deck, loadDecks func()) *DeckDisplay {
	dd := &DeckDisplay{canvas: canvas, registry: registry, loadDecks: loadDecks}

	dd.cardAdapter = NewDeckCardAdapter(dd.registry, dd)

	dd.manaChart = fynecharts.NewBarChart(dd.canvas, "Mana Curve", nil, nil)
	dd.manaChart.SetMinHeight(150)
	dd.manaChart.UpdateHoverFormat(func(v float64) string {
		return fmt.Sprintf("%d", int(v))
	})
	dd.manaChart.UpdateOnTouched(func(idx int) {
		log.Println("Index:", idx)
	})
	dd.tpsChart = fynecharts.NewBarChart(dd.canvas, "Spell Types", nil, nil)
	dd.tpsChart.SetMinHeight(150)
	dd.tpsChart.UpdateHoverFormat(func(v float64) string {
		return fmt.Sprintf("%d", int(v))
	})

	dd.cardList = widget.NewList(dd.cardAdapter.Count, dd.cardAdapter.CreateTemplate, dd.cardAdapter.UpdateTemplate)
	dd.cardAdapter.SetList(dd.cardList)
	dd.cardList.OnSelected = dd.cardSelected

	dd.cardAdapter.Clear()
	dd.cardList.Refresh()
	lbls := []string{"0", "1", "2", "3", "4", "5", "6", "7+"}
	data := []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	tps := []string{"Land", "Artifact", "Creature", "Enchantment", "Sorcery", "Instant", "Planeswalker"}
	typData := []float64{0, 0, 0, 0, 0, 0, 0}
	if fullDeck, err := dd.registry.Manager.LoadDeck(deck.ID); err == nil {
		dd.selectedDeck = fullDeck
		dd.cardAdapter.Clear()
		if fullDeck.Commander != nil {
			dd.cardAdapter.AppendCards([]storage.DeckCard{*fullDeck.Commander})
		}
		if fullDeck.Main != nil {
			dd.cardAdapter.AppendCards(fullDeck.Main)
		}
		if fullDeck.Sideboard != nil {
			dd.cardAdapter.AppendCards(fullDeck.Sideboard)
		}
		dd.cardList.Refresh()

		log.Println(dd.cardAdapter.Count())
		for i := 0; i < dd.cardAdapter.Count(); i++ {
			crd := dd.cardAdapter.Item(i)
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
		for i := 0; i < dd.cardAdapter.Count(); i++ {
			crd := dd.cardAdapter.Item(i)
			if strings.Contains(crd.Card.TypeLine, "Land") {
				typData[0] += float64(crd.Count)
			} else if strings.Contains(crd.Card.TypeLine, "Artifact") {
				typData[1] += float64(crd.Count)
			} else if strings.Contains(crd.Card.TypeLine, "Creature") {
				typData[2] += float64(crd.Count)
			} else if strings.Contains(crd.Card.TypeLine, "Enchantment") {
				typData[3] += float64(crd.Count)
			} else if strings.Contains(crd.Card.TypeLine, "Sorcery") {
				typData[4] += float64(crd.Count)
			} else if strings.Contains(crd.Card.TypeLine, "Instant") {
				typData[5] += float64(crd.Count)
			} else if strings.Contains(crd.Card.TypeLine, "Planeswalker") {
				typData[6] += float64(crd.Count)
			} else {
				log.Println("UNKNOWN:", crd.Card.TypeLine)
			}
		}
		dd.manaChart.UpdateData(lbls, data)
		dd.manaChart.Refresh()
		dd.tpsChart.UpdateData(tps, typData)
		dd.tpsChart.Refresh()
	} else {
		dd.registry.Notifier.ShowError(err)
	}

	dd.cardTabs = container.NewDocTabs()
	hSplit := container.NewHSplit(dd.cardList, dd.cardTabs)
	hSplit.SetOffset(0.2)

	dd.Container = container.NewBorder(container.NewGridWithColumns(2, dd.manaChart, dd.tpsChart), nil, nil, nil, hSplit)

	return dd
}

func (dd *DeckDisplay) cardSelected(id widget.ListItemID) {
	c := dd.cardAdapter.Item(id)

	cardLayout := card.NewCardLayout(dd.canvas, &c.Card, dd.registry)
	tab := container.NewTabItem(c.Card.Name, cardLayout.Container)
	dd.cardTabs.Append(tab)
	dd.cardTabs.Select(tab)
}
