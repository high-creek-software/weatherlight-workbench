package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gitlab.com/high-creek-software/goscryfall"
	"gitlab.com/kendellfab/mtgstudio/internal/bookmarked"
	"gitlab.com/kendellfab/mtgstudio/internal/browse"
	"gitlab.com/kendellfab/mtgstudio/internal/deck"
	"gitlab.com/kendellfab/mtgstudio/internal/icons"
	"gitlab.com/kendellfab/mtgstudio/internal/platform/symbol"
	"gitlab.com/kendellfab/mtgstudio/internal/search"
	"gitlab.com/kendellfab/mtgstudio/internal/storage"
	"gitlab.com/kendellfab/mtgstudio/internal/sync"
	"os"
	"strings"
	"time"
)

const (
	lastSyncKey = "last_sync_key"
	syncFormat  = "2006-01-02 15:04:05"
)

type MtgStudio struct {
	app              fyne.App
	window           fyne.Window
	browseLayout     *browse.BrowseLayout
	searchLayout     *search.SearchLayout
	bookmarkedLayout *bookmarked.BookmarkedLayout
	deckLayout       *deck.DeckLayout

	client        *goscryfall.Client
	manager       *storage.Manager
	importManager *sync.ImportManager

	syncBtn      *widget.Button
	syncLastLbl  *widget.Label
	syncProgress *widget.ProgressBar
	syncSetLbl   *widget.Label
	settingsBtn  *widget.Button

	symbolRepo symbol.SymbolRepo
}

func NewMtgStudio() *MtgStudio {
	os.Setenv("FYNE_THEME", "dark")
	mtgs := &MtgStudio{app: app.NewWithID("gitlab.com/kendellfab/mtgstudio")}
	mtgs.app.SetIcon(icons.AppIconResource)
	mtgs.window = mtgs.app.NewWindow("MTG Studio")
	mtgs.window.SetMaster()
	mtgs.window.Resize(fyne.NewSize(1920, 1080))
	mtgs.client = goscryfall.NewClient()
	mtgs.manager = storage.NewManager(mtgs.client)
	mtgs.importManager = sync.NewImportManager(mtgs.client, mtgs.manager)
	mtgs.symbolRepo = symbol.NewSymbolRepo(mtgs.client, mtgs.manager.LoadSymbolImage)

	mtgs.app.Lifecycle().SetOnStarted(mtgs.appStartedCallback)

	mtgs.setupBody()

	return mtgs
}

func (m *MtgStudio) setupBody() {
	m.browseLayout = browse.NewBrowseLayout(m.manager, m.symbolRepo, m, m.updateSetIcon, m.resizeCardArt)
	m.searchLayout = search.NewSearchLayout(m.manager, m.symbolRepo, m)
	m.bookmarkedLayout = bookmarked.NewBookmarkedLayout(m.manager, m.symbolRepo, m)
	m.deckLayout = deck.NewDeckLayout(m.manager, m.symbolRepo, m.showImport)
	appTabs := container.NewAppTabs(container.NewTabItem("Browse", m.browseLayout.Split),
		container.NewTabItem("Search", m.searchLayout.Split),
		container.NewTabItem("Bookmarked", m.bookmarkedLayout.Split),
		container.NewTabItem("Decks", m.deckLayout.Container),
	)

	m.syncBtn = widget.NewButton("Sync", m.syncBtnTouched)
	m.syncLastLbl = widget.NewLabel("Last Synced")
	m.syncProgress = widget.NewProgressBar()
	m.syncProgress.SetValue(0)
	m.syncProgress.Max = 100
	m.syncProgress.Hide()
	m.syncSetLbl = widget.NewLabel("Set Name")
	m.syncSetLbl.Hide()
	m.settingsBtn = widget.NewButtonWithIcon("", theme.SettingsIcon(), m.settingsBtnTouched)

	// container.NewBorder(nil, nil, m.syncSetLbl, nil, m.syncProgress)
	syncBorder := container.NewBorder(nil, nil, container.NewVBox(layout.NewSpacer(), m.settingsBtn, layout.NewSpacer()), nil, container.NewBorder(nil, nil, container.NewVBox(layout.NewSpacer(), m.syncBtn, layout.NewSpacer()), nil, container.NewBorder(nil, nil, container.NewVBox(layout.NewSpacer(), m.syncLastLbl, layout.NewSpacer()), nil, container.NewVBox(m.syncSetLbl, m.syncProgress))))

	m.window.SetContent(container.NewBorder(nil, syncBorder, nil, nil, appTabs))

	appTabs.OnSelected = func(ti *container.TabItem) {
		if ti.Text == "Bookmarked" {
			m.bookmarkedLayout.LoadBookmarked()
		}
	}
}

func (m *MtgStudio) showImport() {
	window := m.app.NewWindow("Import Deck")
	window.Resize(fyne.NewSize(800, 400))

	entry := widget.NewEntry()
	entry.PlaceHolder = "Deck Name"

	data := widget.NewEntry()
	data.MultiLine = true

	save := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		err := m.manager.ImportDeck(entry.Text, data.Text)
		window.Hide()
		window = nil
		if err != nil {
			m.ShowError(err)
			return
		}
		m.deckLayout.LoadDecks()
	})

	border := container.NewBorder(entry, save, nil, nil, data)
	window.SetContent(border)
	window.Show()
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

func (m *MtgStudio) settingsBtnTouched() {
	dialog.ShowInformation("Settings", "Does nothing yet", m.window)
}

func (m *MtgStudio) syncBtnTouched() {
	m.runSync()
}

func (m *MtgStudio) appStartedCallback() {
	// TODO: Figure out how to determine if an import is needed.
	m.showLastSyncedAt()
	setCount := m.manager.SetCount()
	if setCount == 0 {
		m.runSync()
	} else {
		m.browseLayout.LoadSets()
	}
}

func (m *MtgStudio) runSync() {
	/*progress := widget.NewProgressBar()
	progress.Max = 100
	setName := widget.NewLabel("Set:")
	dialog.ShowCustom("Import progress", "OK", container.NewVBox(setName, progress), m.window)*/
	m.syncSetLbl.Show()
	m.syncProgress.Show()
	startCount := m.manager.SetCount()
	resChan, doneChan, err := m.importManager.Import()
	if err != nil {
		m.ShowError(err)
	} else {
		go func() {
		F:
			for {
				select {
				case status := <-resChan:
					/*setName.SetText(status.SetName)
					progress.SetValue(status.Percent)*/
					m.syncSetLbl.SetText(status.SetName)
					m.syncProgress.SetValue(status.Percent)
					if startCount == 0 {
						m.browseLayout.LoadSets()
					}
				case <-doneChan:
					/*setName.SetText("Import Complete")
					progress.SetValue(100)*/
					m.syncSetLbl.Hide()
					m.syncProgress.Hide()
					m.browseLayout.LoadSets()
					break F
				}
			}
			synced := time.Now()
			m.app.Preferences().SetString(lastSyncKey, synced.Format(time.RFC3339))
			m.showLastSyncedAt()
		}()
	}
}

func (m *MtgStudio) showLastSyncedAt() {
	syncedStr := m.app.Preferences().StringWithFallback(lastSyncKey, "--")
	if syncedStr == "--" {
		m.syncLastLbl.SetText(syncedStr)
		return
	}

	if synced, err := time.Parse(time.RFC3339, syncedStr); err == nil {
		m.syncLastLbl.SetText(synced.Format(syncFormat))
	}

}

func (m *MtgStudio) ShowDialog(title, message string) {
	dialog.NewInformation(title, message, m.window).Show()
}

func (m *MtgStudio) ShowError(err error) {
	dialog.NewError(err, m.window).Show()
}
