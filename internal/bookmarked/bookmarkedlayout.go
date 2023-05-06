package bookmarked

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/tabman"
	"github.com/high-creek-software/weatherlight-workbench/internal/card"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform"
)

type BookmarkedLayout struct {
	widget.BaseWidget
	canvas fyne.Canvas

	registry *platform.Registry

	cardList       *widget.List
	cardAdapter    *card.CardAdapter
	cardTabs       *container.DocTabs
	cardTabManager *tabman.Manager[string]
}

func (bl *BookmarkedLayout) CreateRenderer() fyne.WidgetRenderer {
	split := container.NewHSplit(bl.cardList, bl.cardTabs)
	split.SetOffset(0.18)

	return widget.NewSimpleRenderer(split)
}

func NewBookmarkedLayout(cvs fyne.Canvas, registry *platform.Registry) *BookmarkedLayout {
	bl := &BookmarkedLayout{canvas: cvs, registry: registry}
	bl.ExtendBaseWidget(bl)

	bl.cardAdapter = card.NewCardAdapter(bl.registry)
	bl.cardTabs = container.NewDocTabs()
	bl.cardList = widget.NewList(bl.cardAdapter.Count, bl.cardAdapter.CreateTemplate, bl.cardAdapter.UpdateTemplate)
	bl.cardList.OnSelected = bl.cardSelected
	bl.cardAdapter.SetList(bl.cardList)
	bl.cardTabManager = tabman.NewManager[string]()
	bl.cardTabs.OnClosed = bl.cardTabManager.RemoveTab

	return bl
}

func (bl *BookmarkedLayout) LoadBookmarked() {
	go func() {
		cards, err := bl.registry.Manager.ListBookmarked()
		if err != nil {
			bl.registry.Notifier.ShowError(err)
			return
		}

		bl.cardAdapter.Clear()
		bl.cardAdapter.AppendCards(cards)
		bl.cardList.Refresh()
		bl.cardList.ScrollToTop()
	}()
}

func (bl *BookmarkedLayout) cardSelected(id widget.ListItemID) {
	c := bl.cardAdapter.Item(id)
	bl.cardList.UnselectAll()
	if ti, ok := bl.cardTabManager.GetTabItem(c.Id); ok {
		bl.cardTabs.Select(ti)
		return
	}

	cardLayout := card.NewCardLayout(bl.canvas, &c, bl.registry)
	tab := container.NewTabItem(c.Name, cardLayout)
	bl.cardTabs.Append(tab)
	bl.cardTabs.Select(tab)
	bl.cardTabManager.AddTabItem(c.Id, tab)
}
