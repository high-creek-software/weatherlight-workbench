package set

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/goscryfall/sets"
)

var _ fyne.Widget = (*SetListItem)(nil)

type SetListItem struct {
	widget.BaseWidget
	set *sets.Set
	ico fyne.Resource

	icon        *widget.Icon
	name        *widget.Label
	cardSubhead *widget.Label
	cardCount   *widget.Label
	release     *widget.Label
}

func NewSetListItem(set *sets.Set) *SetListItem {
	sli := &SetListItem{set: set}
	sli.ExtendBaseWidget(sli)

	return sli
}

func (sli *SetListItem) SetResource(resource fyne.Resource) {
	//sli.ico = resource
	sli.icon.SetResource(resource)
}

func (sli *SetListItem) UpdateSet(set *sets.Set) {
	sli.set = set
	//sli.Refresh()
	sli.name.SetText(set.Name)
	sli.cardCount.SetText(fmt.Sprintf("%d", set.CardCount))
	sli.release.SetText(set.ReleasedAt)
}

func (sli *SetListItem) CreateRenderer() fyne.WidgetRenderer {
	//sli.ExtendBaseWidget(sli)
	icon := widget.NewIcon(nil)
	name := widget.NewLabel("template")
	name.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	cardSubhead := widget.NewLabel("Card Count:")
	count := widget.NewLabel("template")
	count.TextStyle = fyne.TextStyle{Italic: true}
	releaseSubhead := widget.NewLabel("Release:")
	release := widget.NewLabel("template")
	release.TextStyle = fyne.TextStyle{Italic: true}
	//render := &SetListItemRenderer{listItem: sli, icon: icon, name: name, cardSubhead: cardSubhead, cardCount: count, releaseSubhead: releaseSubhead, release: release}
	//return render

	sli.icon = icon
	sli.name = name
	sli.cardSubhead = cardSubhead
	sli.cardCount = count
	sli.release = release
	cont := container.NewHBox(sli.icon, container.NewVBox(sli.name, container.NewHBox(
		sli.cardSubhead, sli.cardCount, layout.NewSpacer(), releaseSubhead, release,
	)))

	return widget.NewSimpleRenderer(cont)
}

type SetListItemRenderer struct {
	listItem       *SetListItem
	icon           *widget.Icon
	name           *widget.Label
	cardSubhead    *widget.Label
	cardCount      *widget.Label
	releaseSubhead *widget.Label
	release        *widget.Label
}

func (r *SetListItemRenderer) Destroy() {

}

func (r *SetListItemRenderer) Layout(size fyne.Size) {
	namePos := fyne.NewPos(32+12, 6)
	r.name.Move(namePos)

	iconPos := fyne.NewPos(12, 32/2)
	r.icon.Move(iconPos)
	r.icon.Resize(fyne.NewSize(32, 32))

	//countSubheadSize := fyne.MeasureText(r.cardSubhead.Text, theme.TextSize(), r.cardSubhead.TextStyle)
	//countSize := fyne.MeasureText(r.cardCount.Text, theme.TextSize(), r.cardCount.TextStyle)
	countPos := fyne.NewPos(32+12, 12)
	r.cardSubhead.Move(countPos)
	//r.cardCount.Resize(r.cardCount.Size())

}

func (r *SetListItemRenderer) MinSize() fyne.Size {
	nameSize := fyne.MeasureText(r.name.Text, theme.TextSize(), r.name.TextStyle)
	countSize := fyne.MeasureText(r.cardCount.Text, theme.TextSize(), r.cardCount.TextStyle)

	fyne.Max(nameSize.Height+countSize.Height, 32)

	return fyne.NewSize(r.icon.MinSize().Width+nameSize.Width, 64)
}

func (r *SetListItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.icon, r.name, r.cardSubhead, r.cardCount, r.releaseSubhead, r.release}
}

func (r *SetListItemRenderer) Refresh() {
	if r.listItem.ico != nil {
		r.icon.Resource = r.listItem.ico
		r.icon.Refresh()
	}
	if r.listItem.set != nil {
		r.name.SetText(r.listItem.set.Name)
		//r.name.Refresh()

		cardCount := fmt.Sprintf("%d", r.listItem.set.CardCount)
		r.cardCount.SetText(cardCount)
		//r.cardCount.Refresh()

		//r.cardSubhead.Refresh()
		//r.releaseSubhead.Refresh()
		r.release.SetText(r.listItem.set.ReleasedAt)
		//r.release.Refresh()
	}
}
