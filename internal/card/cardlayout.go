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
	"gitlab.com/kendellfab/mtgstudio/internal/platform"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/icons"
	"gitlab.com/kendellfab/mtgstudio/internal/ruling"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"image/png"
	"log"
)

type CardLayout struct {
	//Container *container.Scroll
	*fyne.Container

	card  *cards.Card
	image *canvas.Image

	registry *platform.Registry

	topBox            *fyne.Container
	addBookmarkBtn    *widget.Button
	removeBookmarkBtn *widget.Button
	docTabs           *container.DocTabs

	cardMetaAdapter *cardMetaAdapter
	cardMetaList    *widget.List
}

func NewCardLayout(card *cards.Card, registry *platform.Registry) *CardLayout {
	cl := &CardLayout{card: card, registry: registry}

	bookmark, _ := cl.registry.Manager.FindBookmark(card.Id)
	cl.addBookmarkBtn = widget.NewButtonWithIcon("", icons.BookmarkResource, func() {
		err := cl.registry.Manager.AddBookmark(card.Id)
		if err != nil {
			cl.registry.Notifier.ShowError(err)
			return
		}
		cl.addBookmarkBtn.Hide()
		cl.removeBookmarkBtn.Show()
		cl.topBox.Refresh()
	})
	cl.addBookmarkBtn.Importance = widget.LowImportance

	cl.removeBookmarkBtn = widget.NewButtonWithIcon("", icons.BookmarkRemoveResource, func() {
		err := cl.registry.Manager.RemoveBookmark(card.Id)
		if err != nil {
			cl.registry.Notifier.ShowError(err)
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

	cl.image = canvas.NewImageFromResource(nil)
	cl.image.FillMode = canvas.ImageFillOriginal
	// Setting this image size to keep the image from overlapping.  Why does it work, I don't know ?!?
	cl.image.Resize(fyne.NewSize(450, 500))

	cl.docTabs = container.NewDocTabs()
	cl.docTabs.SetTabLocation(container.TabLocationLeading)

	mainBox := container.NewBorder(cl.topBox, nil, container.NewPadded(cl.image), nil, container.NewPadded(cl.docTabs))
	cl.Container = mainBox

	cl.setupLegalities()

	var cs cardSummary
	if len(card.CardFaces) == 0 {
		cs = cardSummary{
			name:       card.Name,
			cost:       card.ManaCost,
			typeLine:   card.TypeLine,
			oracleText: card.OracleText,
			flavorText: card.FlavorText,
			power:      card.Power,
			toughness:  card.Toughness,
			releasedAt: card.ReleasedAt,
		}
	} else {
		face := card.CardFaces[0]
		cs = cardSummary{
			name:       face.Name,
			cost:       face.ManaCost,
			typeLine:   face.TypeLine,
			oracleText: face.OracleText,
			flavorText: face.FlavorName,
			releasedAt: card.ReleasedAt,
			power:      face.Power,
			toughness:  face.Toughness,
		}
	}

	cl.setupDetails(cs)

	go func() {
		cardImgPath := card.ImageUris.Png
		if cardImgPath == "" && len(card.CardFaces) > 0 {
			cardImgPath = card.CardFaces[0].ImageUris.Png
		}
		if img, err := cl.registry.Manager.LoadCardImage(cardImgPath); err == nil {
			cl.image.Resource = fyne.NewStaticResource(cardImgPath, cl.resizeImage(img))
			cl.image.Refresh()
		} else {
			cl.image.Resource = icons.FullCardFailedResource
			cl.image.Refresh()
		}
	}()

	go func() {
		rulings, err := cl.registry.Manager.LoadRulings(card)
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

func (cl *CardLayout) SetResource(resource fyne.Resource) {
	cl.image.Resource = resource
	cl.image.Refresh()
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

func (cl *CardLayout) setupDetails(cs cardSummary) {

	metaData := []cardMeta{
		{"Name", cs.name},
		{"Cost", cs.cost},
		{"Type Line", cs.typeLine},
		{"Oracle Text", cs.oracleText},
	}

	if cs.flavorText != "" {
		metaData = append(metaData, cardMeta{"Flavor Text", cs.flavorText})
	}

	if cs.power != "" && cs.toughness != "" {
		metaData = append(metaData, cardMeta{"Power/Toughness", fmt.Sprintf("%s/%s", cs.power, cs.toughness)})
	}

	metaData = append(metaData, cardMeta{"Released At", cs.releasedAt})

	if cl.cardMetaAdapter != nil {
		cl.cardMetaAdapter.updateMeta(metaData)
		cl.cardMetaList.Refresh()
	} else {
		cl.cardMetaAdapter = newCardMetaAdapter(metaData)
		cl.cardMetaList = widget.NewList(cl.cardMetaAdapter.Count, cl.cardMetaAdapter.CreateTemplate, cl.cardMetaAdapter.UpdateTemplate)
		cl.cardMetaAdapter.SetList(cl.cardMetaList)
		cl.docTabs.Append(container.NewTabItem("Details", cl.cardMetaList))
	}
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
