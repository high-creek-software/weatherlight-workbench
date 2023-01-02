package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/goscryfall"
	"gitlab.com/kendellfab/mtgstudio/internal/browse"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/storage"
	"gitlab.com/kendellfab/mtgstudio/internal/sync"
	"strings"
)

type MtgStudio struct {
	app          fyne.App
	window       fyne.Window
	browseLayout *browse.BrowseLayout

	client        *goscryfall.Client
	manager       *storage.Manager
	importManager *sync.ImportManager

	symbolRepo symbol.SymbolRepo
}

func NewMtgStudio() *MtgStudio {
	//os.Setenv("FYNE_THEME", "light")
	mtgs := &MtgStudio{app: app.New()}
	mtgs.window = mtgs.app.NewWindow("MTG Studio")
	mtgs.window.Resize(fyne.NewSize(1200, 700))
	mtgs.manager = storage.NewManager()
	mtgs.client = goscryfall.NewClient()
	mtgs.importManager = sync.NewImportManager(mtgs.client, mtgs.manager)
	mtgs.symbolRepo = symbol.NewSymbolRepo(mtgs.client, mtgs.manager.LoadSymbolImage)

	mtgs.app.Lifecycle().SetOnStarted(mtgs.appStartedCallback)

	mtgs.setupBody()

	return mtgs
}

func (m *MtgStudio) setupBody() {
	m.browseLayout = browse.NewBrowseLayout(m.manager, m.symbolRepo, m, m.updateSetIcon, m.resizeCardArt)
	appTabs := container.NewAppTabs(container.NewTabItem("Browse", m.browseLayout.Split))
	m.window.SetContent(appTabs)
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
	return bs
	//buff := bytes.NewBuffer(bs)
	//img, err := jpeg.Decode(buff)
	//if err != nil {
	//	log.Println("error parsing image:", err)
	//	return bs
	//}
	//
	//r := resize.Resize(150, 0, img, resize.Lanczos3)
	//
	//var out bytes.Buffer
	//jpeg.Encode(&out, r, nil)
	//return out.Bytes()
}

func (m *MtgStudio) Start() {
	m.window.ShowAndRun()
}

func (m *MtgStudio) appStartedCallback() {
	// TODO: Figure out how to determine if an import is needed.

	setCount := m.manager.SetCount()
	if setCount == 0 {
		progress := widget.NewProgressBar()
		progress.Max = 100
		setName := widget.NewLabel("Set:")
		dialog.ShowCustom("Import progress", "OK", container.NewVBox(setName, progress), m.window)
		resChan, doneChan, err := m.importManager.Import()
		if err != nil {
			m.ShowError(err)
		} else {
			go func() {
				for {
					select {
					case status := <-resChan:
						setName.SetText(status.SetName)
						progress.SetValue(status.Percent)
					case <-doneChan:
						m.browseLayout.LoadSets()
					}
				}
			}()
		}
	} else {
		m.browseLayout.LoadSets()
	}
}

func (m *MtgStudio) ShowDialog(title, message string) {
	dialog.NewInformation(title, message, m.window).Show()
}

func (m *MtgStudio) ShowError(err error) {
	dialog.NewError(err, m.window).Show()
}
