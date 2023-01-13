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
	*container.Scroll

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

	cl.docTabs = container.NewDocTabs()
	cl.docTabs.SetTabLocation(container.TabLocationLeading)

	mainBox := container.NewBorder(cl.topBox, nil, nil, nil, container.NewPadded(container.NewBorder(container.NewHBox(layout.NewSpacer(), image, layout.NewSpacer()), nil, nil, nil, container.NewPadded(cl.docTabs))))

	cl.Scroll = container.NewScroll(container.NewPadded(mainBox))
	cl.Scroll.Direction = container.ScrollVerticalOnly
	cl.Refresh()

	cl.setupDetails()
	cl.setupLegalities()

	go func() {
		if img, err := cl.manager.LoadCardImage(card.ImageUris.Png); err == nil {
			image.Resource = fyne.NewStaticResource(card.ImageUris.Png, cl.resizeImage(img))
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
	table := widget.NewTable(func() (int, int) {
		return 5, 2
	}, func() fyne.CanvasObject {
		return widget.NewLabel("")
	}, func(cell widget.TableCellID, co fyne.CanvasObject) {
		lbl := co.(*widget.Label)
		if cell.Row == 0 {
			if cell.Col == 0 {
				lbl.SetText("Name")
			} else if cell.Col == 1 {
				lbl.SetText(cl.card.Name)
			}
		} else if cell.Row == 1 {
			if cell.Col == 0 {
				lbl.SetText("Released At")
			} else if cell.Col == 1 {
				lbl.SetText(cl.card.ReleasedAt)
			}
		} else if cell.Row == 2 {
			if cell.Col == 0 {
				lbl.SetText("Type Line")
			} else if cell.Col == 1 {
				lbl.SetText(cl.card.TypeLine)
			}
		} else if cell.Row == 3 {
			if cell.Col == 0 {
				lbl.SetText("Oracle Text")
			} else if cell.Col == 1 {
				lbl.SetText(cl.card.OracleText)
			}
		} else if cell.Row == 4 {
			if cell.Col == 0 {
				lbl.SetText("Flavor Text")
			} else if cell.Col == 1 {
				lbl.SetText(cl.card.FlavorText)
			}
		}
	})
	table.SetColumnWidth(0, 100)
	table.SetColumnWidth(1, 500)
	table.SetRowHeight(3, 75)
	table.SetRowHeight(4, 75)

	cl.docTabs.Append(container.NewTabItem("Details", table))
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

	cl.docTabs.Append(container.NewTabItem("Legalities", legalitiesTable))
}
