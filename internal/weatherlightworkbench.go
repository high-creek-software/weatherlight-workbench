package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/high-creek-software/ansel"
	"github.com/high-creek-software/bento"
	"github.com/high-creek-software/goscryfall"
	"github.com/high-creek-software/weatherlight-workbench/internal/bookmarked"
	"github.com/high-creek-software/weatherlight-workbench/internal/browse"
	"github.com/high-creek-software/weatherlight-workbench/internal/deck"
	"github.com/high-creek-software/weatherlight-workbench/internal/importing"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/icons"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/storage"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/symbol"
	"github.com/high-creek-software/weatherlight-workbench/internal/platform/sync"
	"github.com/high-creek-software/weatherlight-workbench/internal/search"
	"golang.org/x/image/colornames"
	"strings"
	"time"
)

const (
	lastSyncKey = "last_sync_key"
	syncFormat  = "2006-01-02 15:04:05"
)

type WeatherlightWorkbench struct {
	app              fyne.App
	window           fyne.Window
	browseLayout     *browse.BrowseLayout
	searchLayout     *search.SearchLayout
	bookmarkedLayout *bookmarked.BookmarkedLayout
	deckLayout       *deck.DeckLayout

	registry *platform.Registry

	syncBtn      *widget.Button
	syncLastLbl  *widget.Label
	syncProgress *widget.ProgressBar
	syncSetLbl   *widget.Label
	settingsBtn  *widget.Button

	bentoBox *bento.Box
}

func NewWeatherlightWorkbench() *WeatherlightWorkbench {
	//os.Setenv("FYNE_THEME", "light")
	wm := &WeatherlightWorkbench{app: app.NewWithID("github.com/high-creek-software/weatherlight-workbench")}
	wm.app.SetIcon(icons.AppIconResource)
	wm.window = wm.app.NewWindow("Weatherlight Workbench")
	wm.window.SetMaster()
	wm.window.Resize(fyne.NewSize(1920, 1080))
	client := goscryfall.NewClient()
	manager := storage.NewManager(client)
	importManager := sync.NewImportManager(client, manager)
	symbolRepo := symbol.NewSymbolRepo(client, manager.LoadSymbolImage)

	wm.registry = platform.NewRegistry(manager, symbolRepo, client, importManager, wm)
	wm.registry.SymbolRepo = symbolRepo

	wm.registry.SetIconLoader = ansel.NewAnsel[string](200, ansel.SetLoadedCallback[string](wm.updateSetIcon), ansel.SetLoader[string](wm.registry.Manager.LoadSetIcon))
	wm.registry.CardThumbnailLoader = ansel.NewAnsel[string](800, ansel.SetLoader[string](wm.registry.Manager.LoadCardImage), ansel.SetWorkerCount[string](20), ansel.SetLoadingImage[string](icons.CardLoadingResource), ansel.SetFailedImage[string](icons.CardFailedResource))
	wm.registry.CardFullLoader = ansel.NewAnsel[string](200, ansel.SetLoader[string](wm.registry.Manager.LoadCardImage), ansel.SetLoadingImage[string](icons.FullCardLoadingResource), ansel.SetFailedImage[string](icons.FullCardFailedResource))

	wm.app.Lifecycle().SetOnStarted(wm.appStartedCallback)

	wm.setupBody()

	return wm
}

func (m *WeatherlightWorkbench) setupBody() {
	m.browseLayout = browse.NewBrowseLayout(m.window.Canvas(), m.registry, m.updateSetIcon, m.resizeCardArt)
	m.searchLayout = search.NewSearchLayout(m.window.Canvas(), m.registry)
	m.bookmarkedLayout = bookmarked.NewBookmarkedLayout(m.window.Canvas(), m.registry)
	m.deckLayout = deck.NewDeckLayout(m.window.Canvas(), m.registry, m.showImport)
	appTabs := container.NewAppTabs(container.NewTabItem("Browse", m.browseLayout.Split),
		container.NewTabItem("Search", m.searchLayout.Split),
		container.NewTabItem("Bookmarked", m.bookmarkedLayout.Split),
		container.NewTabItem("Decks", m.deckLayout),
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
	mainBody := container.NewBorder(nil, syncBorder, nil, nil, appTabs)

	m.bentoBox = bento.NewBox()
	m.bentoBox.UpdateBottomOffset(45)

	m.window.SetContent(container.NewMax(mainBody, m.bentoBox))

	appTabs.OnSelected = func(ti *container.TabItem) {
		if ti.Text == "Bookmarked" {
			m.bookmarkedLayout.LoadBookmarked()
		} else if ti.Text == "Decks" {
			m.deckLayout.LoadDecks()
		}
	}
}

func (m *WeatherlightWorkbench) showImport() {
	window := m.app.NewWindow("Import Deck")
	window.Resize(fyne.NewSize(1200, 700))

	il := importing.NewImportLayout(m.registry, func() {
		window.Close()
		m.deckLayout.LoadDecks()
	})
	window.SetContent(il.Split)
	window.Show()

}

func (m *WeatherlightWorkbench) updateSetIcon(bs []byte) []byte {
	if m.app.Settings().ThemeVariant() == theme.VariantDark {
		strData := string(bs)
		if strings.Contains(strData, `fill="#000"`) {
			strData = strings.Replace(strData, `fill="#000"`, `fill="#999"`, -1)
		} else {
			strData = strings.Replace(strData, "<path", `<path style="fill:#999999"`, -1)
		}
		return []byte(strData)
	}
	return bs
}

func (m *WeatherlightWorkbench) resizeCardArt(bs []byte) []byte {
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

func (m *WeatherlightWorkbench) Start() {
	m.window.ShowAndRun()
}

func (m *WeatherlightWorkbench) settingsBtnTouched() {
	dialog.ShowInformation("Settings", "Does nothing yet", m.window)
}

func (m *WeatherlightWorkbench) syncBtnTouched() {
	m.runSync()
}

func (m *WeatherlightWorkbench) appStartedCallback() {
	// TODO: Figure out how to determine if an import is needed.
	m.showLastSyncedAt()
	setCount := m.registry.Manager.SetCount()
	if setCount == 0 {
		m.runSync()
	} else {
		m.browseLayout.LoadSets()
	}
}

func (m *WeatherlightWorkbench) runSync() {
	m.syncSetLbl.Show()
	m.syncProgress.Show()
	startCount := m.registry.Manager.SetCount()
	resChan, doneChan, err := m.registry.ImportManager.Import()
	if err != nil {
		m.ShowError(err)
	} else {
		go func() {
		F:
			for {
				select {
				case status := <-resChan:
					m.syncSetLbl.SetText(status.SetName)
					m.syncProgress.SetValue(status.Percent)
					if startCount == 0 {
						m.browseLayout.LoadSets()
					}
				case <-doneChan:
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

func (m *WeatherlightWorkbench) showLastSyncedAt() {
	syncedStr := m.app.Preferences().StringWithFallback(lastSyncKey, "--")
	if syncedStr == "--" {
		m.syncLastLbl.SetText(syncedStr)
		return
	}

	if synced, err := time.Parse(time.RFC3339, syncedStr); err == nil {
		m.syncLastLbl.SetText(synced.Format(syncFormat))
	}

}

func (m *WeatherlightWorkbench) ShowDialog(title, message string) {
	item := bento.NewItemWithMessage(message, bento.LengthLong)
	item.SetBackgroundColor(colornames.Darkslateblue)
	m.bentoBox.AddItem(item)
}

func (m *WeatherlightWorkbench) ShowError(err error) {
	item := bento.NewItemWithMessage(err.Error(), bento.LengthIndefinite)
	item.SetBackgroundColor(colornames.Red)
	m.bentoBox.AddItem(item)
}

func (m *WeatherlightWorkbench) VerifyAction(message, actionTitle string, action func()) {
	item := bento.NewItemWithMessage(message, bento.LengthIndefinite)
	item.SetBackgroundColor(theme.PrimaryColor())
	m.bentoBox.AddItem(item)
	item.AddAction(actionTitle, action)
}
