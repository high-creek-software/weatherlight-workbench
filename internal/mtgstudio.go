package internal

import (
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nfnt/resize"
	"gitlab.com/high-creek-software/ansel"
	"gitlab.com/high-creek-software/goscryfall"
	"gitlab.com/kendellfab/mtgstudio/internal/card"
	"gitlab.com/kendellfab/mtgstudio/internal/resources"
	"gitlab.com/kendellfab/mtgstudio/internal/set"
	"image/jpeg"
	"log"
	"strings"
)

type MtgStudio struct {
	app    fyne.App
	window fyne.Window

	setList    *widget.List
	setAdapter *set.SetAdapter

	cardList    *widget.List
	cardAdapter *card.CardAdapter

	client  *goscryfall.Client
	manager *resources.Manager
}

func NewMtgStudio() *MtgStudio {
	//os.Setenv("FYNE_THEME", "light")
	mtgs := &MtgStudio{app: app.New()}
	mtgs.window = mtgs.app.NewWindow("MTG Studio")
	mtgs.window.Resize(fyne.NewSize(1200, 700))
	mtgs.manager = resources.NewManager()
	mtgs.client = goscryfall.NewClient()
	mtgs.setAdapter = set.NewSetAdapter(
		ansel.NewAnsel[string](ansel.SetLoadedCallback[string](mtgs.updateSetIcon), ansel.SetLoader[string](mtgs.manager.LoadSetIcon)),
	)
	mtgs.cardAdapter = card.NewCardAdapter(
		ansel.NewAnsel[string](ansel.SetLoader[string](mtgs.manager.LoadCardImage), ansel.SetLoadedCallback[string](mtgs.resizeCardArt)),
	)

	mtgs.setupBody()

	return mtgs
}

func (m *MtgStudio) setupBody() {
	m.setList = widget.NewList(m.setAdapter.Count, m.setAdapter.CreateTemplate, m.setAdapter.UpdateTemplate)
	m.setList.OnSelected = func(id widget.ListItemID) {
		set := m.setAdapter.Item(id)
		log.Println("Selected:", set.Id)
		m.cardAdapter.Clear()
		m.cardList.Refresh()
		go func() {
			cards, err := m.client.ListCards(set.Code, "")
			if err != nil {
				dialog.NewInformation("Error loading cards", err.Error(), m.window).Show()
				return
			}
			m.cardAdapter.AppendCards(cards.Data)
		}()
	}
	m.cardList = widget.NewList(m.cardAdapter.Count, m.cardAdapter.CreateTemplate, m.cardAdapter.UpdateTemplate)
	insideSplit := container.NewHSplit(m.cardList, container.NewMax())
	insideSplit.SetOffset(0.25)
	split := container.NewHSplit(m.setList, insideSplit)
	split.SetOffset(0.25)
	//insideBorder := container.NewBorder(nil, nil, m.cardList, nil, container.NewMax())
	//outsideBorder := container.NewBorder(nil, nil, m.setList, nil, insideBorder)
	m.window.SetContent(split)
}

func (m *MtgStudio) updateSetIcon(bs []byte) []byte {
	if m.app.Settings().ThemeVariant() == theme.VariantDark {
		strData := string(bs)
		if strings.Contains(strData, `fill="#000"`) {
			strData = strings.Replace(strData, `fill="#000"`, `fill="#999"`, -1)
		} else {
			strData = strings.Replace(strData, "<path d=", `<path style="fill:#999999" d=`, -1)
		}
		return []byte(strData)
	}
	return bs
}

func (m *MtgStudio) resizeCardArt(bs []byte) []byte {

	buff := bytes.NewBuffer(bs)
	img, err := jpeg.Decode(buff)
	if err != nil {
		log.Println("error parsing image:", err)
		return bs
	}

	r := resize.Resize(150, 0, img, resize.Lanczos3)

	var out bytes.Buffer
	jpeg.Encode(&out, r, nil)
	return out.Bytes()
}

func (m *MtgStudio) Start() {
	m.loadSets()
	m.window.ShowAndRun()
}

func (m *MtgStudio) loadSets() {
	sets, err := m.client.ListSets()
	if err != nil {
		dialog.NewInformation("Error loading sets", err.Error(), m.window).Show()
		return
	}

	m.setAdapter.AddSets(sets)
}
