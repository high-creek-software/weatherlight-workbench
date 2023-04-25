package deck

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	scryfallcards "github.com/high-creek-software/goscryfall/cards"
	"github.com/high-creek-software/tabman"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
	"golang.org/x/exp/maps"
	"log"
)

type DeckLayout struct {
	widget.BaseWidget
	canvas fyne.Canvas

	toolbar  *widget.Toolbar
	deckList *widget.List
	cardList *widget.List

	deckTabs       *container.DocTabs
	deckTabManager *tabman.Manager[string]

	deckAdapter  *DeckAdapter
	registry     *platform.Registry
	selectedDeck storage.Deck
}

func (dl *DeckLayout) CreateRenderer() fyne.WidgetRenderer {
	hSplit := container.NewHSplit(dl.deckList, dl.deckTabs)
	hSplit.SetOffset(0.15)
	con := container.NewBorder(dl.toolbar, nil, nil, nil, hSplit)
	return widget.NewSimpleRenderer(con)
}

func NewDeckLayout(canvas fyne.Canvas, registry *platform.Registry, showImport func()) *DeckLayout {
	dl := &DeckLayout{canvas: canvas, registry: registry}
	dl.ExtendBaseWidget(dl)
	dl.deckAdapter = NewDeckAdapter(nil, dl.registry, dl)

	dl.deckList = widget.NewList(dl.deckAdapter.Count, dl.deckAdapter.CreateTemplate, dl.deckAdapter.UpdateTemplate)
	dl.deckTabs = container.NewDocTabs()
	dl.deckAdapter.SetList(dl.deckList)
	dl.deckList.OnSelected = dl.deckSelected
	dl.deckTabManager = tabman.NewManager[string]()
	dl.deckTabs.OnClosed = dl.deckTabManager.RemoveTab

	dl.toolbar = widget.NewToolbar(widget.NewToolbarAction(theme.DownloadIcon(), showImport), widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		var popup *widget.PopUp
		nameEntry := widget.NewEntry()
		legalities := maps.Keys(scryfallcards.LegalitiesNameMap)
		legalitySelect := widget.NewSelect(legalities, nil)

		frm := widget.NewForm(widget.NewFormItem("Deck Name", nameEntry), widget.NewFormItem("Deck Type", legalitySelect))

		cancelBtn := widget.NewButton("Cancel", func() {
			popup.Hide()
		})
		saveBtn := widget.NewButton("Save", func() {
			typ := scryfallcards.LegalitiesNameMap[legalitySelect.Selected]
			_, err := dl.registry.Manager.CreateDeck(nameEntry.Text, typ)
			if err != nil {
				dl.registry.Notifier.ShowError(err)
			} else {
				dl.LoadDecks()
			}
			popup.Hide()
		})
		grid := container.NewGridWithColumns(2, cancelBtn, saveBtn)

		popup = widget.NewModalPopUp(container.NewVBox(frm, grid), dl.canvas)
		popup.Show()
	}))

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
	//dl.registry.Notifier.ShowDialog("", fmt.Sprintf("Copy %s; not done yet", d.Name))
	var popup *widget.PopUp
	nameEntry := widget.NewEntry()
	frm := widget.NewForm(widget.NewFormItem("Deck Name:", nameEntry))
	cancelBtn := widget.NewButton("Cancel", func() {

		popup.Hide()
	})
	saveBtn := widget.NewButton("Save", func() {
		name := nameEntry.Text
		popup.Hide()

		go func() {
			_, err := dl.registry.Manager.CopyDeck(d, name)
			if err != nil {
				dl.registry.Notifier.ShowError(err)
				return
			}
			dl.LoadDecks()
		}()
	})
	title := widget.NewLabel(fmt.Sprintf("Copy existing deck %s", d.Name))
	bottom := container.NewGridWithColumns(2, cancelBtn, saveBtn)
	popup = widget.NewModalPopUp(container.NewBorder(container.NewGridWithColumns(3, widget.NewLabel("	"), title, widget.NewLabel("	")), bottom, nil, nil, frm), dl.canvas)
	popup.Show()
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
	dl.deckList.UnselectAll()

	if ti, ok := dl.deckTabManager.GetTabItem(deck.ID); ok {
		dl.deckTabs.Select(ti)
		return
	}

	dl.selectedDeck = deck
	deckDisplay := NewDeckMetaDisplay(dl.canvas, dl.registry, deck, dl.LoadDecks)
	tab := container.NewTabItem(deck.Name, deckDisplay)
	dl.deckTabs.Append(tab)
	dl.deckTabs.Select(tab)
	dl.deckTabManager.AddTabItem(deck.ID, tab)
}
