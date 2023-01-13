package ruling

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/goscryfall/rulings"
)

type RulingListItem struct {
	widget.BaseWidget

	r rulings.Ruling
}

func (r *RulingListItem) Set(input rulings.Ruling) {
	r.r = input
	r.Refresh()
}

func (r *RulingListItem) CreateRenderer() fyne.WidgetRenderer {
	dateTitle := widget.NewLabel("When:")
	dateLbl := widget.NewLabel("")
	orgTitle := widget.NewLabel("Organization:")
	orgLbl := widget.NewLabel("")
	reasonLbl := widget.NewLabel("")
	reasonLbl.Wrapping = fyne.TextWrapWord

	return &rulingListItemRenderer{
		r:         r,
		dateTitle: dateTitle,
		dateLbl:   dateLbl,
		orgTitle:  orgTitle,
		orgLbl:    orgLbl,
		reasonLbl: reasonLbl,
	}
}

func NewRulingListItem() *RulingListItem {
	rll := &RulingListItem{}
	rll.ExtendBaseWidget(rll)

	return rll
}

type rulingListItemRenderer struct {
	r         *RulingListItem
	dateTitle *widget.Label
	dateLbl   *widget.Label
	orgTitle  *widget.Label
	orgLbl    *widget.Label
	reasonLbl *widget.Label
}

func (r rulingListItemRenderer) Destroy() {

}

func (r rulingListItemRenderer) Layout(size fyne.Size) {
	topLeft := fyne.NewPos(theme.Padding(), theme.Padding())
	r.dateTitle.Move(topLeft)
	dtSize := r.dateTitle.MinSize()
	dlPos := topLeft.Add(fyne.NewPos(dtSize.Width+theme.Padding(), 0))
	r.dateLbl.Move(dlPos)

	orgSize := r.orgLbl.MinSize()
	topRight := fyne.NewPos(size.Width-(orgSize.Width+theme.Padding()*2), theme.Padding())
	r.orgLbl.Move(topRight)
	orgTitleSize := r.orgTitle.MinSize()
	topRight = topRight.Subtract(fyne.NewPos(orgTitleSize.Width+theme.Padding(), 0))
	r.orgTitle.Move(topRight)

	topLeft = topLeft.Add(fyne.NewPos(0, dtSize.Height+theme.Padding()))
	r.reasonLbl.Move(topLeft)
	reasonSize := fyne.NewSize(size.Width-theme.Padding(), r.reasonLbl.MinSize().Height)
	r.reasonLbl.Resize(reasonSize)
}

func (r rulingListItemRenderer) MinSize() fyne.Size {
	topSize := r.dateTitle.MinSize()
	reasonSize := r.reasonLbl.MinSize()

	size := fyne.NewSize(reasonSize.Width, topSize.Height+reasonSize.Height+theme.Padding()*4)
	return size
}

func (r rulingListItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.dateTitle, r.dateLbl, r.orgTitle, r.orgLbl, r.reasonLbl}
}

func (r rulingListItemRenderer) Refresh() {
	r.dateLbl.SetText(r.r.r.PublishedAt)
	r.orgLbl.SetText(r.r.r.Source)
	r.reasonLbl.SetText(r.r.r.Comment)
}
