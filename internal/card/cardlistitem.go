package card

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/goscryfall/cards"
)

var _ fyne.Widget = (*CardListItem)(nil)

type CardListItem struct {
	widget.BaseWidget
	card *cards.Card
	ico  fyne.Resource

	icon    *widget.Icon
	name    *widget.Label
	manaBox *fyne.Container
}

func NewCardListItem(card *cards.Card) *CardListItem {
	cli := &CardListItem{card: card}
	cli.ExtendBaseWidget(cli)

	return cli
}

func (cli *CardListItem) UpdateCard(card *cards.Card) {
	cli.ClearManaCost()
	cli.card = card
	cli.name.SetText(card.Name)
}

func (cli *CardListItem) SetManaCost(bs []fyne.Resource) {
	cli.manaBox.RemoveAll()
	for _, b := range bs {
		cli.manaBox.Add(widget.NewIcon(b))
	}
}

func (cli *CardListItem) ClearManaCost() {
	cli.manaBox.RemoveAll()
}

func (cli *CardListItem) SetResource(resource fyne.Resource) {
	cli.icon.SetResource(resource)
}

func (cli *CardListItem) CreateRenderer() fyne.WidgetRenderer {
	icon := widget.NewIcon(nil)
	name := widget.NewLabel("template")
	manaBox := container.NewHBox()

	cli.icon = icon
	cli.name = name
	cli.manaBox = manaBox

	cont := container.NewGridWithColumns(3, cli.icon, cli.name, cli.manaBox)

	return widget.NewSimpleRenderer(cont)
}
