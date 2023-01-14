package card

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nfnt/resize"
	"gitlab.com/high-creek-software/goscryfall/cards"
	"gitlab.com/kendellfab/mtgstudio/internal/icons"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/notifier"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/ruling"
	"gitlab.com/kendellfab/mtgstudio/internal/storage"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"image/png"
	"log"
)

type CardLayout struct {
	//Container *container.Scroll
	*fyne.Container

	card       *cards.Card
	symbolRepo symbol.SymbolRepo
	manager    *storage.Manager

	topBox            *fyne.Container
	addBookmarkBtn    *widget.Button
	removeBookmarkBtn *widget.Button
	docTabs           *container.DocTabs

	notifier notifier.Notifier
}

func NewCardLayout(card *cards.Card, symbolRepo symbol.SymbolRepo, manager *storage.Manager, n notifier.Notifier) *CardLayout {
	cl := &CardLayout{card: card, symbolRepo: symbolRepo, manager: manager, notifier: n}

	bookmark, _ := cl.manager.FindBookmark(card.Id)
	cl.addBookmarkBtn = widget.NewButtonWithIcon("", icons.BookmarkResource, func() {
		err := cl.manager.AddBookmark(card.Id)
		if err != nil {
			cl.notifier.ShowError(err)
			return
		}
		cl.addBookmarkBtn.Hide()
		cl.removeBookmarkBtn.Show()
		cl.topBox.Refresh()
	})
	cl.addBookmarkBtn.Importance = widget.LowImportance

	cl.removeBookmarkBtn = widget.NewButtonWithIcon("", icons.BookmarkRemoveResource, func() {
		err := cl.manager.RemoveBookmark(card.Id)
		if err != nil {
			cl.notifier.ShowError(err)
			return
		}
		cl.removeBookmarkBtn.Hide()
		cl.addBookmarkBtn.Show()
		cl.topBox.Refresh()
	})
	cl.removeBookmarkBtn.Importance = widget.LowImportance

	if bookmark == nil {
		cl.removeBookmarkBtn.Hide()
	} else {
		cl.addBookmarkBtn.Hide()
	}

	cl.topBox = container.NewBorder(nil, nil, nil, cl.addBookmarkBtn, container.NewBorder(nil, nil, nil, cl.removeBookmarkBtn))

	image := canvas.NewImageFromResource(icons.FullCardLoadingResource)
	image.FillMode = canvas.ImageFillOriginal
	// Setting this image size to keep the image from overlapping.  Why does it work, I don't know ?!?
	image.Resize(fyne.NewSize(450, 500))

	cl.docTabs = container.NewDocTabs()
	cl.docTabs.SetTabLocation(container.TabLocationLeading)

	mainBox := container.NewBorder(cl.topBox, nil, container.NewPadded(image), nil, container.NewPadded(cl.docTabs))
	cl.Container = mainBox

	cl.setupLegalities()
	cl.setupDetails()

	go func() {
		cardImgPath := card.ImageUris.Png
		if cardImgPath == "" && len(card.CardFaces) > 0 {
			cardImgPath = card.CardFaces[0].ImageUris.Png
		}
		if img, err := cl.manager.LoadCardImage(cardImgPath); err == nil {
			image.Resource = fyne.NewStaticResource(cardImgPath, cl.resizeImage(img))
			image.Refresh()
		} else {
			image.Resource = icons.FullCardFailedResource
			image.Refresh()
		}
	}()

	go func() {
		rulings, err := manager.LoadRulings(card)
		if err != nil {
			log.Println("error loading rulings", err)
		} else if len(rulings) > 0 {
			adapter := ruling.NewRulingAdapter(rulings)
			ruleList := widget.NewList(adapter.Count, adapter.CreateTemplate, adapter.UpdateTemplate)
			adapter.SetList(ruleList)
			cl.docTabs.Append(container.NewTabItem("Rulings", ruleList))
		}
	}()

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

func (cl *CardLayout) setupDetails() {

	metaData := []cardMeta{
		{"Name", cl.card.Name},
		{"Cost", cl.card.ManaCost},
		{"Type Line", cl.card.TypeLine},
		{"Oracle Text", cl.card.OracleText},
	}

	if cl.card.FlavorText != "" {
		metaData = append(metaData, cardMeta{"Flavor Text", cl.card.FlavorText})
	}

	if cl.card.Power != "" && cl.card.Toughness != "" {
		metaData = append(metaData, cardMeta{"Power/Toughness", fmt.Sprintf("%s/%s", cl.card.Power, cl.card.Toughness)})
	}

	metaData = append(metaData, cardMeta{"Released At", cl.card.ReleasedAt})

	metaListAdapter := newCardMetaAdapter(metaData)
	metaList := widget.NewList(metaListAdapter.Count, metaListAdapter.CreateTemplate, metaListAdapter.UpdateTemplate)
	metaListAdapter.SetList(metaList)
	cl.docTabs.Append(container.NewTabItem("Details", metaList))
}

func (cl *CardLayout) setupLegalities() {
	legalities := cl.card.Legalities
	keys := maps.Keys(legalities)
	slices.Sort(keys)

	var lls []fyne.CanvasObject
	maxSize := fyne.NewSize(0, 0)
	for _, key := range keys {
		ll := newLegalityItem(key, legalities[key])
		maxSize = maxSize.Max(ll.MinSize())
		lls = append(lls, ll)
	}

	legalitiesTable := container.NewGridWrap(maxSize, lls...)

	scroll := container.NewScroll(legalitiesTable)
	scroll.Direction = container.ScrollVerticalOnly
	padded := container.NewPadded(scroll)
	cl.docTabs.Append(container.NewTabItem("Legalities", padded))
}
