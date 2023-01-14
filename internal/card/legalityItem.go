package card

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/goscryfall/cards"
	"gitlab.com/kendellfab/mtgstudio/internal/icons"
)

type legalityItem struct {
	widget.BaseWidget

	title    string
	legality cards.Legality
}

func (l *legalityItem) CreateRenderer() fyne.WidgetRenderer {
	titleTxt := widget.NewLabel("")
	image := widget.NewIcon(icons.LegalResource)

	return &legalityItemRenderer{
		ll:       l,
		titleTxt: titleTxt,
		image:    image,
	}
}

func newLegalityItem(title string, legality cards.Legality) *legalityItem {
	ll := &legalityItem{title: title, legality: legality}
	ll.ExtendBaseWidget(ll)

	ll.Refresh()

	return ll
}

type legalityItemRenderer struct {
	ll       *legalityItem
	titleTxt *widget.Label
	image    *widget.Icon
}

func (l legalityItemRenderer) Destroy() {

}

func (l legalityItemRenderer) Layout(size fyne.Size) {
	topLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	titleSize := l.titleTxt.MinSize()
	l.titleTxt.Move(topLeft)

	topLeft = topLeft.Add(fyne.NewPos(0, titleSize.Height-10))
	l.image.Move(topLeft)
	l.image.Resize(fyne.NewSize(100, 50))
}

func (l legalityItemRenderer) MinSize() fyne.Size {
	titleSize := l.titleTxt.MinSize()
	imgSize := l.image.MinSize()

	return fyne.NewSize(fyne.Max(titleSize.Width, imgSize.Width)+2*theme.Padding(), (titleSize.Height+20)+imgSize.Height+3*theme.Padding())
}

func (l legalityItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{l.titleTxt, l.image}
}

func (l legalityItemRenderer) Refresh() {
	l.titleTxt.Text = l.ll.title
	l.titleTxt.Refresh()

	switch l.ll.legality {
	case cards.Legal:
		l.image.SetResource(icons.LegalResource)
	case cards.NotLegal:
		l.image.SetResource(icons.NotLegalResource)
	case cards.Restricted:
		l.image.SetResource(icons.RestrictedResource)
	case cards.Banned:
		l.image.SetResource(icons.BannedResource)
	}
}
