package card

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/goscryfall/cards"
)

var _ fyne.Widget = (*CardListItem)(nil)

type CardListItem struct {
	widget.BaseWidget
	card *cards.Card
	ico  fyne.Resource

	icon        *widget.Icon
	name        *widget.Label
	cost        *widget.Label
	ManaSymbols []*widget.Icon
}

func NewCardListItem(card *cards.Card) *CardListItem {
	cli := &CardListItem{card: card}
	cli.ExtendBaseWidget(cli)

	return cli
}

func (cli *CardListItem) UpdateCard(card *cards.Card) {
	cli.card = card
	cli.name.SetText(card.Name)
	cli.cost.SetText(card.ManaCost)
}

func (cli *CardListItem) SetResource(resource fyne.Resource) {
	cli.icon.SetResource(resource)
}

func (cli *CardListItem) CreateRenderer() fyne.WidgetRenderer {
	icon := widget.NewIcon(nil)
	name := widget.NewLabel("template")
	cost := widget.NewLabel("template")

	cli.icon = icon
	cli.name = name
	cli.cost = cost

	cont := container.NewGridWithColumns(4, cli.icon, cli.name, layout.NewSpacer(), cli.cost)

	//sets := cli.card.ParseManaCost()
	//if len(sets) > 0 {
	//	for idx := 0; idx < len(sets[0]); idx++ {
	//		mIcon := widget.NewIcon(nil)
	//		cli.ManaSymbols = append(cli.ManaSymbols, mIcon)
	//		cont.Add(mIcon)
	//	}
	//}

	return widget.NewSimpleRenderer(cont)
}
