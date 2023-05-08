package set

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/goscryfall/sets"
)

var _ fyne.Widget = (*SetListItem)(nil)

type SetListItem struct {
	widget.BaseWidget
	set *sets.Set
	ico fyne.Resource
}

func NewSetListItem(set *sets.Set) *SetListItem {
	sli := &SetListItem{set: set}
	sli.ExtendBaseWidget(sli)

	return sli
}

func (sli *SetListItem) SetResource(resource fyne.Resource) {
	sli.ico = resource
}

func (sli *SetListItem) UpdateSet(set *sets.Set) {
	sli.set = set
	sli.Refresh()
}

func (sli *SetListItem) CreateRenderer() fyne.WidgetRenderer {
	//sli.ExtendBaseWidget(sli)
	icon := widget.NewIcon(nil)
	name := widget.NewLabel("template")
	name.Wrapping = fyne.TextWrapWord
	name.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	cardSubhead := widget.NewLabel("Card Count:")
	count := widget.NewLabel("template")
	count.TextStyle = fyne.TextStyle{Italic: true}
	releaseSubhead := widget.NewLabel("Release:")
	release := widget.NewLabel("template")
	release.TextStyle = fyne.TextStyle{Italic: true}
	render := &SetListItemRenderer{listItem: sli, icon: icon, name: name, cardSubhead: cardSubhead, cardCount: count, releaseSubhead: releaseSubhead, release: release}
	return render
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
	iconSquare := float32(48.0)
	iconPos := fyne.NewPos(12, size.Height/2-iconSquare/2)
	r.icon.Move(iconPos)
	r.icon.Resize(fyne.NewSize(iconSquare, iconSquare))

	namePos := fyne.NewPos(iconSquare+24, 0)
	r.name.Move(namePos)
	r.name.Resize(fyne.NewSize(size.Width-iconSquare-3*theme.Padding(), r.name.MinSize().Height))

	countSubheadSize := r.countSubheadSize()
	countSize := r.countSize()
	releaseSubheadSize := r.releaseSubheadSize()
	releaseSize := r.releaseSize()

	cumulativePos := namePos.Add(fyne.NewPos(0, r.name.MinSize().Height))
	r.cardSubhead.Move(cumulativePos)

	cumulativePos = cumulativePos.Add(fyne.NewPos(countSubheadSize.Width+8, 0))
	r.cardCount.Move(cumulativePos)

	if iconSquare+countSubheadSize.Width+countSize.Width+releaseSubheadSize.Width+releaseSize.Width+5*theme.Padding()+16 > size.Width {
		cumulativePos = fyne.NewPos(namePos.X, cumulativePos.Y+countSubheadSize.Height+theme.Padding())
	} else {
		cumulativePos = cumulativePos.Add(fyne.NewPos(countSize.Width+16, 0))
	}

	r.releaseSubhead.Move(cumulativePos)

	cumulativePos = cumulativePos.Add(fyne.NewPos(releaseSubheadSize.Width+8, 0))
	r.release.Move(cumulativePos)
}

func (r *SetListItemRenderer) nameSize() fyne.Size {
	return fyne.MeasureText(r.name.Text, theme.TextSize(), r.name.TextStyle)
}

func (r *SetListItemRenderer) countSubheadSize() fyne.Size {
	return fyne.MeasureText(r.cardSubhead.Text, theme.TextSize(), r.cardSubhead.TextStyle)
}

func (r *SetListItemRenderer) countSize() fyne.Size {
	return fyne.MeasureText(r.cardCount.Text, theme.TextSize(), r.cardCount.TextStyle)
}

func (r *SetListItemRenderer) releaseSubheadSize() fyne.Size {
	return fyne.MeasureText(r.releaseSubhead.Text, theme.TextSize(), r.releaseSubhead.TextStyle)
}

func (r *SetListItemRenderer) releaseSize() fyne.Size {
	return fyne.MeasureText(r.release.Text, theme.TextSize(), r.release.TextStyle)
}

func (r *SetListItemRenderer) MinSize() fyne.Size {
	iconSize := r.icon.Size()
	//nameSize := r.nameSize()
	//countSubheadSize := r.countSubheadSize()
	//countSize := r.countSize()
	//releaseSubheadSize := r.releaseSubheadSize()
	//releaseSize := r.releaseSize()

	//topRow := iconSize.Width + 12 + nameSize.Width + 32
	//bottomRow := iconSize.Width + 12 + countSubheadSize.Width + 8 + countSize.Width + 16 + releaseSubheadSize.Width + 8 + releaseSize.Width + 32

	nameSize := r.name.MinSize()
	countSubheadSize := r.cardSubhead.MinSize()
	releaseSubheadSize := r.releaseSubhead.MinSize()

	height := nameSize.Height + countSubheadSize.Height

	if r.releaseSubhead.Position().Y > r.cardSubhead.Position().Y {
		height += releaseSubheadSize.Height + 2*theme.Padding()
	}

	return fyne.NewSize(250, fyne.Max(iconSize.Height+24, height))
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

		cardCount := fmt.Sprintf("%d", r.listItem.set.CardCount)
		r.cardCount.SetText(cardCount)

		r.release.SetText(r.listItem.set.ReleasedAt)
	}
}
