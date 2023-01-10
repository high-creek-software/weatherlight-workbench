package card

import (
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/nfnt/resize"
	"gitlab.com/high-creek-software/goscryfall/cards"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/notifier"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/storage"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"image/png"
	"log"
)

type CardLayout struct {
	*container.Scroll
	vBox *fyne.Container

	card       *cards.Card
	symbolRepo symbol.SymbolRepo
	manager    *storage.Manager

	topBox            *fyne.Container
	addBookmarkBtn    *widget.Button
	removeBookmarkBtn *widget.Button

	notifier notifier.Notifier
}

func NewCardLayout(card *cards.Card, symbolRepo symbol.SymbolRepo, manager *storage.Manager, n notifier.Notifier) *CardLayout {
	cl := &CardLayout{card: card, symbolRepo: symbolRepo, manager: manager, notifier: n}

	cl.vBox = container.NewVBox()
	cl.Scroll = container.NewScroll(cl.vBox)
	cl.Scroll.Direction = container.ScrollVerticalOnly

	bookmark, _ := cl.manager.FindBookmark(card.Id)
	cl.addBookmarkBtn = widget.NewButtonWithIcon("Add", storage.BookmarkResource, func() {
		err := cl.manager.AddBookmark(card.Id)
		if err != nil {
			cl.notifier.ShowError(err)
			return
		}
		cl.addBookmarkBtn.Hide()
		cl.removeBookmarkBtn.Show()
		cl.topBox.Refresh()
	})

	cl.removeBookmarkBtn = widget.NewButtonWithIcon("Remove", storage.BookmarkRemoveResource, func() {
		err := cl.manager.RemoveBookmark(card.Id)
		if err != nil {
			cl.notifier.ShowError(err)
			return
		}
		cl.removeBookmarkBtn.Hide()
		cl.addBookmarkBtn.Show()
		cl.topBox.Refresh()
	})
	if bookmark == nil {
		cl.removeBookmarkBtn.Hide()
	} else {
		cl.addBookmarkBtn.Hide()
	}
	cl.topBox = container.NewHBox(layout.NewSpacer(), cl.addBookmarkBtn, cl.removeBookmarkBtn)
	cl.vBox.Add(cl.topBox)

	image := canvas.NewImageFromResource(nil)
	image.FillMode = canvas.ImageFillOriginal
	cl.vBox.Add(container.NewHBox(layout.NewSpacer(), image, layout.NewSpacer()))

	go func() {
		if img, err := cl.manager.LoadCardImage(card.ImageUris.Png); err == nil {
			image.Resource = fyne.NewStaticResource(card.ImageUris.Png, cl.resizeImage(img))
			image.Refresh()
		}

	}()
	cl.setupLegalities()

	return cl
}

func (cl *CardLayout) resizeImage(bs []byte) []byte {
	buff := bytes.NewBuffer(bs)
	img, err := png.Decode(buff)
	if err != nil {
		log.Println("error parsing image:", err)
		return bs
	}

	r := resize.Resize(450, 0, img, resize.Lanczos3)

	var out bytes.Buffer
	png.Encode(&out, r)
	return out.Bytes()
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
