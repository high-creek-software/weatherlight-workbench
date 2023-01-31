package importing

import (
	"fmt"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	scryfallcards "github.com/high-creek-software/goscryfall/cards"
	"github.com/high-creek-software/goscryfall/decks"
	"github.com/high-creek-software/weatherlight-workbench/internal/card"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
	"golang.org/x/exp/maps"
)

//type currentCardType int
//
//const (
//	cardTypeUnknown currentCardType = iota
//	cardTypeCommander
//	cardTypeCompanion
//	cardTypeSideboard
//	cardTypeDeck
//)

type ImportLayout struct {
	*container.Split

	registry *platform.Registry

	nameEntry *widget.Entry
	deckType  *widget.Select
	deckEntry *widget.Entry
	saveBtn   *widget.Button

	nameLbl     *widget.RichText
	cardNameLbl *widget.RichText

	cardAdapter *card.CardAdapter
	cardList    *widget.List

	cardType     storage.AssociationType
	cardNamePair *decks.NameCountPair
	sectionIndex int
	deck         *decks.Deck

	createdDeck    *storage.Deck
	importComplete func()
}

func NewImportLayout(reg *platform.Registry, importComplete func()) *ImportLayout {
	il := &ImportLayout{registry: reg, importComplete: importComplete}
	il.nameEntry = widget.NewEntry()
	il.nameEntry.PlaceHolder = "Deck Name"
	names := maps.Keys(scryfallcards.LegalitiesNameMap)
	il.deckType = widget.NewSelect(names, func(val string) {})
	il.deckType.PlaceHolder = "Deck type"
	il.deckEntry = widget.NewEntry()
	il.deckEntry.MultiLine = true
	il.saveBtn = widget.NewButtonWithIcon("Import", theme.DownloadIcon(), il.doImport)

	il.nameLbl = widget.NewRichTextWithText("")
	il.cardNameLbl = widget.NewRichTextWithText("")

	il.cardAdapter = card.NewCardAdapter(reg)
	il.cardList = widget.NewList(il.cardAdapter.Count, il.cardAdapter.CreateTemplate, il.cardAdapter.UpdateTemplate)
	il.cardList.OnSelected = il.cardSelected
	il.cardAdapter.SetList(il.cardList)

	importSide := container.NewBorder(il.nameEntry, il.saveBtn, nil, nil, container.NewBorder(il.deckType, nil, nil, nil, il.deckEntry))
	resolveSide := container.NewBorder(il.nameLbl, nil, nil, nil, container.NewBorder(il.cardNameLbl, nil, nil, nil, il.cardList))

	il.Split = container.NewHSplit(importSide, resolveSide)
	il.Split.SetOffset(0.33)
	return il
}

func (il *ImportLayout) doImport() {
	deck, err := il.registry.Manager.ParseDeckDefinition(il.nameEntry.Text, il.deckEntry.Text)
	if err != nil {
		il.registry.Notifier.ShowError(err)
		return
	}
	il.registry.Notifier.ShowDialog("", "Parsed the deck definition")

	il.deck = deck

	typ := scryfallcards.LegalitiesNameMap[il.deckType.Selected]
	il.createdDeck, err = il.registry.Manager.CreateDeck(il.nameEntry.Text, typ)
	if err != nil {
		il.registry.Notifier.ShowError(err)
		return
	}

	il.nameLbl.ParseMarkdown(fmt.Sprintf("# %s", deck.Name))
	il.advanceImport()
}

func (il *ImportLayout) loadPossibles() {
	il.cardNameLbl.ParseMarkdown(fmt.Sprintf("## %s", il.cardNamePair.Name))
	il.cardList.UnselectAll()
	il.cardAdapter.Clear()
	cs, err := il.registry.Manager.FindByName(il.cardNamePair.Name)
	if err != nil {
		il.registry.Notifier.ShowError(err)
		return
	}
	if len(cs) == 0 {
		il.registry.Notifier.ShowDialog("", fmt.Sprintf("Could not find %s", il.cardNamePair.Name))
		il.advanceImport()
	} else if len(cs) == 1 {
		il.addCardToDeck(cs[0])
		il.advanceImport()
	} else {
		il.cardAdapter.AppendCards(cs)
		il.cardList.Refresh()
	}
}

func (il *ImportLayout) cardSelected(id widget.ListItemID) {
	il.addCardToDeck(il.cardAdapter.Item(id))
	il.advanceImport()
}

func (il *ImportLayout) addCardToDeck(card scryfallcards.Card) {
	err := il.registry.Manager.AddCardToDeck(*il.createdDeck, card.Id, card.Name, il.cardNamePair.Count, il.cardType)
	if err != nil {
		il.registry.Notifier.ShowError(err)
	}
}

func (il *ImportLayout) advanceImport() {
	if il.cardType == storage.AssociationUnknown {
		il.cardType = storage.AssociationCommander
		if il.deck.Commander.Name != "" {
			il.cardNamePair = &il.deck.Commander
			il.loadPossibles()
			return
		}
		il.advanceImport()
	} else if il.cardType == storage.AssociationCommander {
		il.cardType = storage.AssociationCompanion
		if il.deck.Companion.Name != "" {
			il.cardNamePair = &il.deck.Companion
			il.loadPossibles()
			return
		}
		il.advanceImport()
	} else if il.cardType == storage.AssociationCompanion {
		il.sectionIndex = 0
		il.cardType = storage.AssociationMain
		il.advanceImport()
	} else if il.cardType == storage.AssociationMain {
		if il.sectionIndex < len(il.deck.Deck) {
			il.cardNamePair = &il.deck.Deck[il.sectionIndex]
			il.sectionIndex++
			il.loadPossibles()
			return
		}
		il.sectionIndex = 0
		il.cardType = storage.AssociationSideboard
		il.advanceImport()
	} else if il.cardType == storage.AssociationSideboard {
		if il.sectionIndex < len(il.deck.Sideboard) {
			il.cardNamePair = &il.deck.Sideboard[il.sectionIndex]
			il.sectionIndex++
			il.loadPossibles()
			return
		}
		il.cardAdapter.Clear()
		il.cardList.Refresh()
		il.cardNameLbl.ParseMarkdown(fmt.Sprintf("## Import complete"))
		il.registry.Notifier.ShowDialog("", "It looks like we're done with the import.")
		il.importComplete()
	}
}
