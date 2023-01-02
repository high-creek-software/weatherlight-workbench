package card

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/goscryfall/cards"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/storage"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type CardLayout struct {
	*container.Scroll
	vBox *fyne.Container

	card       *cards.Card
	symbolRepo symbol.SymbolRepo
	manager    *storage.Manager
}

func NewCardLayout(card *cards.Card, symbolRepo symbol.SymbolRepo, manager *storage.Manager) *CardLayout {
	cl := &CardLayout{card: card, symbolRepo: symbolRepo, manager: manager}

	cl.vBox = container.NewVBox()
	cl.Scroll = container.NewScroll(cl.vBox)
	cl.Scroll.Direction = container.ScrollVerticalOnly

	image := canvas.NewImageFromResource(nil)
	image.FillMode = canvas.ImageFillOriginal
	cl.vBox.Add(container.NewHBox(layout.NewSpacer(), image, layout.NewSpacer()))

	go func() {
		if img, err := cl.manager.LoadCardImage(card.ImageUris.Png); err == nil {
			image.Resource = fyne.NewStaticResource(card.ImageUris.Png, img)
			image.Refresh()
		}

	}()
	cl.setupLegalities()

	return cl
}

func (cl *CardLayout) setupLegalities() {
	legalities := cl.card.Legalities
	keys := maps.Keys(legalities)
	slices.Sort(keys)

	legalitiesTable := container.NewGridWithColumns(4)
	for _, key := range keys {
		legalitiesTable.Add(widget.NewLabel(key))
		ico := widget.NewIcon(nil)
		switch legalities[key] {
		case cards.Legal:
			ico.SetResource(storage.LegalResource)
		case cards.NotLegal:
			ico.SetResource(storage.NotLegalResource)
		case cards.Restricted:
			ico.SetResource(storage.RestrictedResource)
		case cards.Banned:
			ico.SetResource(storage.BannedResource)
		}
		legalitiesTable.Add(ico)
	}

	hBox := container.NewHBox(layout.NewSpacer(), legalitiesTable, layout.NewSpacer())
	cl.vBox.Add(widget.NewSeparator())
	cl.vBox.Add(hBox)
}
